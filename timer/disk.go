package timer

import (
	"monitor/hardware"
	"monitor/proto"
	"time"

	"github.com/tarm/serial"
)

// DiskTimer 磁盘定时器
type DiskTimer struct {
	*baseTimer
}

// NewDiskTimer 创建磁盘定时器
func NewDiskTimer(port *serial.Port) *DiskTimer {
	return &DiskTimer{
		baseTimer: newBaseTimer(port, time.Minute), // 磁盘信息每分钟更新一次
	}
}

// Start 启动定时器
func (t *DiskTimer) Start() {
	if t.isRunning() {
		return
	}
	t.setRunning(true)

	// 立即发送一次（发送前检查状态）
	if t.isRunning() {
		send(t.port, proto.CmdDisk, hardware.GetDiskInfoForDisplay())
	}

	go func() {
		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()

		for {
			select {
			case <-t.stopCh:
				return
			case <-ticker.C:
				// 采集后再次检查状态，避免 Stop 后仍发送
				if t.isRunning() {
					send(t.port, proto.CmdDisk, hardware.GetDiskInfoForDisplay())
				}
			}
		}
	}()
}

// Stop 停止定时器
func (t *DiskTimer) Stop() {
	t.stop()
}
