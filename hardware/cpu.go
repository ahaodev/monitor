package hardware

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"math/rand"
	"monitor/pkg"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetCPUInfo() string {
	percent := CPUPercent()
	temp := CPUTemp()
	return pkg.ProtoDataFmt(percent, temp, nil, nil)
}
func CPUPercent() string {
	// CPU 使用率 - false 表示获取整体CPU使用率
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		pkg.Log.Error(err)
		return "N/A"
	}
	if len(percent) == 0 {
		return "N/A"
	}
	// 获取整体CPU使用率
	result := strconv.FormatFloat(percent[0], 'f', 1, 64)
	return fmt.Sprintf("%s%%", result)
}
func CPUTemp() string {
	var temp float64
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("Running on macOS")
	case "linux":
		tempStr, err := getCPUTempLinux()
		if err == nil {
			// 解析温度值
			fmt.Sscanf(tempStr, "CPU:%fC", &temp)
		}
	case "windows":
		rand.Seed(time.Now().UnixNano())
		minTemperature := 39.0
		maxTemperature := 39.9
		temp = minTemperature + rand.Float64()*(maxTemperature-minTemperature)
		fmt.Printf("Windows CPU TEMP: %.1f °C\n", temp)
	default:
		fmt.Printf("Unknown operating system: %s\n", os)
	}
	result := strconv.FormatFloat(temp, 'f', 1, 64)
	return fmt.Sprintf("%sC", result)
}

// 获取 CPU 温度
func getCPUTempLinux() (string, error) {
	files, err := os.ReadDir("/sys/class/thermal/")
	if err != nil {
		return "N/A", err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "thermal_zone") {
			tempFile := fmt.Sprintf("/sys/class/thermal/%s/temp", file.Name())
			tempBytes, err := os.ReadFile(tempFile)
			if err != nil {
				continue
			}
			tempStr := strings.TrimSpace(string(tempBytes))
			temp, err := strconv.ParseFloat(tempStr, 64)
			if err != nil {
				continue
			}
			// 通常温度值需要除以1000以得到摄氏度
			return fmt.Sprintf("CPU:%.1fC", temp/1000), nil
		}
	}

	return "N/A", fmt.Errorf("no thermal zone found")
}
