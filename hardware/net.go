package hardware

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"time"
)

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

	// 计算每秒的进站和出站数据
	incomingKbps := float64(currNetStat[0].BytesRecv-prevNetStat[0].BytesRecv) / 1024
	outgoingKbps := float64(currNetStat[0].BytesSent-prevNetStat[0].BytesSent) / 1024

	fmt.Printf("每秒进站数据: %.2f kbps, 每秒出站数据: %.2f kbps\n", incomingKbps, outgoingKbps)
	return fmt.Sprintf("%.2fkbps,%.2fkbps", incomingKbps, outgoingKbps), nil
}
