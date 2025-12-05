package hardware

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v3/disk"
	"monitor/pkg"
)

const (
	GB = 1024 * 1024 * 1024
	TB = 1024 * GB
)

type DiskInfo struct {
	Mountpoint string
	Total      uint64
	Used       uint64
}

func GetDiskInfoForDisplay() string {
	disks := listDisks()
	if len(disks) == 0 {
		return "N/A,N/A"
	}

	var total, used uint64
	for _, d := range disks {
		total += d.Total
		used += d.Used
	}
	usedPercent := fmt.Sprintf("%.1f%%", float64(used)/float64(total)*100)
	usage := formatDiskUsage(used, total)
	return pkg.ProtoDataFmt(usedPercent, usage, nil, nil)
}

func formatDiskUsage(used, total uint64) string {
	if total >= TB {
		return fmt.Sprintf("%.1f/%.1fT", float64(used)/float64(TB), float64(total)/float64(TB))
	}
	return fmt.Sprintf("%.0f/%.0fG", float64(used)/float64(GB), float64(total)/float64(GB))
}

func findRootDisk(disks []DiskInfo) DiskInfo {
	for _, d := range disks {
		if d.Mountpoint == "/" {
			return d
		}
	}
	return disks[0]
}

func isRealFilesystem(fstype string) bool {
	// 只统计实际的存储文件系统
	realFS := map[string]bool{
		"ext4":  true,
		"ext3":  true,
		"ext2":  true,
		"xfs":   true,
		"btrfs": true,
		"zfs":   true,
		"ntfs":  true,
		"vfat":  true,
		"fat32": true,
		"exfat": true,
		"hfs":   true,
		"apfs":  true,
		"cifs":  true,
		"nfs":   true,
		"nfs4":  true,
	}
	return realFS[fstype]
}

func listDisks() []DiskInfo {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Get disk partitions error: %v", err)
		return nil
	}

	var result []DiskInfo
	for _, p := range partitions {
		if !isRealFilesystem(p.Fstype) {
			continue
		}
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		result = append(result, DiskInfo{
			Mountpoint: p.Mountpoint,
			Total:      usage.Total,
			Used:       usage.Used,
		})
	}
	return result
}
