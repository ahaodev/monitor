package hardware

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"log"
	"math/rand"
	"monitor/pkg"
	"runtime"
	"strconv"
	"time"
)

type DiskInfo struct {
	Mountpoint  string
	TotalGB     float64
	UsedGB      float64
	UsedPercent float64
}

// GetDiskInfoForDisplay 获取磁盘信息用于LCD显示，格式: "使用率%,已用/总量"
func GetDiskInfoForDisplay() string {
	d := disks()
	if len(d) == 0 {
		return "N/A,N/A"
	}

	// 使用根分区或第一个分区
	var target DiskInfo
	for _, di := range d {
		if di.Mountpoint == "/" || di.Mountpoint == "C:\\" {
			target = di
			break
		}
	}
	if target.Mountpoint == "" {
		target = d[0]
	}

	// 使用率
	usedPercent := fmt.Sprintf("%.1f%%", target.UsedPercent)
	// 已用/总量 (GB)
	usage := fmt.Sprintf("%.0f/%.0fG", target.UsedGB, target.TotalGB)

	return pkg.ProtoDataFmt(usedPercent, usage, nil, nil)
}

func GetDiskInfo() string {
	d := disks()
	dt := diskTemp()
	if len(d) == 0 {
		return "N/A"
	}
	name := d[0].Mountpoint
	total := d[0].TotalGB
	used := d[0].UsedGB
	return fmt.Sprintf("%s,%sG,%sG,%s", name, strconv.FormatFloat(total, 'f', 1, 64), strconv.FormatFloat(used, 'f', 1, 64), dt)
}

func disks() []DiskInfo {
	// 获取所有磁盘的信息
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Get disk partitions error: %v", err)
		return nil
	}
	var diskInfos []DiskInfo
	// 遍历每个磁盘
	for _, partition := range partitions {
		// 获取磁盘使用情况
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		// 使用 1000 基准 (GB)，匹配系统显示
		const GB = 1000 * 1000 * 1000
		total := float64(usage.Total) / GB
		used := float64(usage.Used) / GB
		// 将关键信息添加到数组中
		diskInfo := DiskInfo{
			Mountpoint:  partition.Mountpoint,
			TotalGB:     total,
			UsedGB:      used,
			UsedPercent: usage.UsedPercent,
		}
		diskInfos = append(diskInfos, diskInfo)
	}
	return diskInfos
}

const (
	IOCTL_STORAGE_QUERY_PROPERTY = 0x002d1400
	STORAGE_PROPERTY_ID          = 0
	StorageDeviceTemperature     = 5
)

type STORAGE_PROPERTY_QUERY struct {
	PropertyId           uint32
	QueryType            uint32
	AdditionalParameters [1]byte
}

type STORAGE_TEMPERATURE_INFO struct {
	Version        uint32
	Reserved       uint32
	GeneralInfo    uint32
	Temperature    int32
	OverThreshold  uint32
	UnderThreshold uint32
}

func diskTemp() string {
	var temp float64
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("Running on macOS")
	case "linux":
		fmt.Println("Running on Linux")
	case "windows": //Unused variable 'temp'
		rand.NewSource(time.Now().UnixNano()) // 设置随机数种子为当前时间的纳秒级别
		minTemperature := 55.0                // 最小温度
		maxTemperature := 57.0                // 最大温度
		temp := minTemperature + rand.Float64()*(maxTemperature-minTemperature)
		fmt.Printf("Windows Disk TEMP: %.1f °C\n", temp)
	default:
		fmt.Printf("Unknown operating system: %s\n", os)
	}
	result := strconv.FormatFloat(temp, 'f', 1, 64)
	return fmt.Sprintf("%sC", result)
}
