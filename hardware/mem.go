package hardware

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
	"monitor/pkg"
)

const GiB = 1024 * 1024 * 1024

func GetMemInfo() string {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		pkg.Log.Error(err)
		return "N/A,N/A"
	}

	usedPercent := fmt.Sprintf("%.1f%%", vmStat.UsedPercent)
	usedGB := float64(vmStat.Used) / GiB
	totalGB := float64(vmStat.Total) / GiB
	usage := fmt.Sprintf("%.1f/%.0fG", usedGB, totalGB)

	return pkg.ProtoDataFmt(usedPercent, usage, nil, nil)
}
