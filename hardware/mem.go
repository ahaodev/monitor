package hardware

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
	"monitor/pkg"
)

func GetMemInfo() string {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		pkg.Log.Error(err)
		return "N/A,N/A"
	}

	usedPercent := fmt.Sprintf("%.1f%%", vmStat.UsedPercent)
	usedGB := float64(vmStat.Used) / GB
	totalGB := float64(vmStat.Total) / GB
	usage := fmt.Sprintf("%.1f/%.0fG", usedGB, totalGB)

	return pkg.ProtoDataFmt(usedPercent, usage, nil, nil)
}
