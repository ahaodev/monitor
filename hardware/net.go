package hardware

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"time"
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

func Net() (string, error) {
	// 使用 false 获取所有网卡的汇总数据
	prevNetStat, err := net.IOCounters(false)
	if err != nil || len(prevNetStat) == 0 {
		fmt.Printf("获取网络信息时出错: %s", err)
		return "", err
	}

	time.Sleep(time.Second) // 等待一秒钟

	currNetStat, err := net.IOCounters(false)
	if err != nil || len(currNetStat) == 0 {
		fmt.Printf("获取网络信息时出错: %s", err)
		return "", err
	}

	// 计算每秒的进站和出站字节数
	incomingBps := float64(currNetStat[0].BytesRecv - prevNetStat[0].BytesRecv)
	outgoingBps := float64(currNetStat[0].BytesSent - prevNetStat[0].BytesSent)

	inStr := formatSpeed(incomingBps)
	outStr := formatSpeed(outgoingBps)

	fmt.Printf("每秒进站数据: %s, 每秒出站数据: %s\n", inStr, outStr)
	return fmt.Sprintf("%s,%s", inStr, outStr), nil
}
