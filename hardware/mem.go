package hardware

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"monitor/pkg"
)

// 单位换算常量 (使用 1000 基准，匹配系统显示)
const (
	GB = 1000 * 1000 * 1000 // 1 GB = 10^9 bytes
)

// GetMemInfo 获取内存信息，格式: "使用率%,已用/总量"
func GetMemInfo() string {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		pkg.Log.Error(err)
		return "N/A,N/A"
	}

	// 使用率
	usedPercent := fmt.Sprintf("%.1f%%", vmStat.UsedPercent)

	// 已用/总量 (GB, 1000基准)
	usedGB := float64(vmStat.Used) / GB
	totalGB := float64(vmStat.Total) / GB
	usage := fmt.Sprintf("%.1f/%.0fG", usedGB, totalGB)

	return pkg.ProtoDataFmt(usedPercent, usage, nil, nil)
}

func Mem() (string, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("获取内存信息时出错: %s", err)
		return "获取内存错误", err
	}

	totalGB := float64(vmStat.Total) / GB
	usedGB := float64(vmStat.Used) / GB

	fmt.Printf("内存总量: %.2f GB, 内存使用量: %.2f GB\n", totalGB, usedGB)
	return fmt.Sprintf("内存总量: %.2f GB, 内存使用量: %.2f GB", totalGB, usedGB), nil
}
