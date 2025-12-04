package hardware

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v3/disk"
	"monitor/pkg"
)

const GB = 1000 * 1000 * 1000

type DiskInfo struct {
	Mountpoint  string
	TotalGB     float64
	UsedGB      float64
	UsedPercent float64
}

func GetDiskInfoForDisplay() string {
	disks := listDisks()
	if len(disks) == 0 {
		return "N/A,N/A"
	}

	target := findRootDisk(disks)
	usedPercent := fmt.Sprintf("%.1f%%", target.UsedPercent)
	usage := fmt.Sprintf("%.0f/%.0fG", target.UsedGB, target.TotalGB)
	return pkg.ProtoDataFmt(usedPercent, usage, nil, nil)
}

func findRootDisk(disks []DiskInfo) DiskInfo {
	for _, d := range disks {
		if d.Mountpoint == "/" {
			return d
		}
	}
	return disks[0]
}

func listDisks() []DiskInfo {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Get disk partitions error: %v", err)
		return nil
	}

	var result []DiskInfo
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		result = append(result, DiskInfo{
			Mountpoint:  p.Mountpoint,
			TotalGB:     float64(usage.Total) / GB,
			UsedGB:      float64(usage.Used) / GB,
			UsedPercent: usage.UsedPercent,
		})
	}
	return result
}
