package hardware

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"monitor/pkg"
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

	// 已用/总量 (GB)
	usedGB := float64(vmStat.Used) / 1024 / 1024 / 1024
	totalGB := float64(vmStat.Total) / 1024 / 1024 / 1024
	usage := fmt.Sprintf("%.1f/%.0fG", usedGB, totalGB)

	return pkg.ProtoDataFmt(usedPercent, usage, nil, nil)
}

func Mem() (string, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("获取内存信息时出错: %s", err)
		return "获取内存错误", err
	}

	totalGB := float64(vmStat.Total) / 1024 / 1024 / 1024
	usedGB := float64(vmStat.Used) / 1024 / 1024 / 1024

	fmt.Printf("内存总量: %.2f GB, 内存使用量: %.2f GB\n", totalGB, usedGB)
	return fmt.Sprintf("内存总量: %.2f GB, 内存使用量: %.2f GB", totalGB, usedGB), nil
}
