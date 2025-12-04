package timer

import (
	"monitor/hardware"
	"monitor/proto"
	"time"

	"github.com/tarm/serial"
)

// CPUTimer CPU定时器
type CPUTimer struct {
	*baseTimer
}

// NewCPUTimer 创建CPU定时器
func NewCPUTimer(port *serial.Port) *CPUTimer {
	return &CPUTimer{
		baseTimer: newBaseTimer(port, time.Second),
	}
}

// Start 启动定时器
func (t *CPUTimer) Start() {
	if t.isRunning() {
		return
	}
	t.setRunning(true)

	// 立即发送一次（发送前检查状态）
	if t.isRunning() {
		send(t.port, proto.CmdCPU, hardware.GetCPUInfo())
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
					send(t.port, proto.CmdCPU, hardware.GetCPUInfo())
				}
			}
		}
	}()
}

// Stop 停止定时器
func (t *CPUTimer) Stop() {
	t.stop()
}
