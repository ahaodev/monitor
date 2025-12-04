package hardware

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"monitor/pkg"
)

func GetCPUInfo() string {
	return pkg.ProtoDataFmt(CPUPercent(), CPUTemp(), nil, nil)
}

func CPUPercent() string {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		pkg.Log.Error(err)
		return "N/A"
	}
	if len(percent) == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.1f%%", percent[0])
}

func CPUTemp() string {
	temp, err := readCPUTemp()
	if err != nil {
		return "N/A"
	}
	return fmt.Sprintf("%.1fC", temp)
}

func readCPUTemp() (float64, error) {
	files, err := os.ReadDir("/sys/class/thermal/")
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "thermal_zone") {
			continue
		}
		tempFile := fmt.Sprintf("/sys/class/thermal/%s/temp", file.Name())
		data, err := os.ReadFile(tempFile)
		if err != nil {
			continue
		}
		temp, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err != nil {
			continue
		}
		return temp / 1000, nil
	}
	return 0, fmt.Errorf("no thermal zone found")
}
