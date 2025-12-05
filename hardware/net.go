package hardware

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

// formatSpeed 根据速度大小自动选择合适的单位
func formatSpeed(bytesPerSec float64) string {
	if bytesPerSec >= 1024*1024 {
		return fmt.Sprintf("%.1fMB/s", bytesPerSec/1024/1024)
	} else if bytesPerSec >= 1024 {
		return fmt.Sprintf("%.1fKB/s", bytesPerSec/1024)
	}
	return fmt.Sprintf("%.0fB/s", bytesPerSec)
}

func isPhysicalInterface(name string) bool {
	// 排除虚拟接口
	excludes := []string{"lo", "docker", "veth", "br-", "virbr"}
	for _, ex := range excludes {
		if strings.HasPrefix(name, ex) {
			return false
		}
	}
	return true
}

func Net() (string, error) {
	// 使用 true 获取每个网卡的数据
	prevNetStat, err := net.IOCounters(true)
	if err != nil || len(prevNetStat) == 0 {
		fmt.Printf("获取网络信息时出错: %s", err)
		return "", err
	}

	time.Sleep(time.Second)

	currNetStat, err := net.IOCounters(true)
	if err != nil || len(currNetStat) == 0 {
		fmt.Printf("获取网络信息时出错: %s", err)
		return "", err
	}

	// 构建当前统计的 map
	currMap := make(map[string]net.IOCountersStat)
	for _, stat := range currNetStat {
		currMap[stat.Name] = stat
	}

	// 计算所有物理网卡的总流量
	var totalRecv, totalSent uint64
	for _, prev := range prevNetStat {
		if !isPhysicalInterface(prev.Name) {
			continue
		}
		if curr, ok := currMap[prev.Name]; ok {
			totalRecv += curr.BytesRecv - prev.BytesRecv
			totalSent += curr.BytesSent - prev.BytesSent
		}
	}

	inStr := formatSpeed(float64(totalRecv))
	outStr := formatSpeed(float64(totalSent))

	fmt.Printf("进站数据: %s, 出站数据: %s\n", inStr, outStr)
	return fmt.Sprintf("%s,%s", inStr, outStr), nil
}
