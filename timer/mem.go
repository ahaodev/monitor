package timer

import (
	"monitor/hardware"
	"monitor/proto"
	"time"

	"github.com/tarm/serial"
)

// MemTimer 内存定时器
type MemTimer struct {
	*baseTimer
}

// NewMemTimer 创建内存定时器
func NewMemTimer(port *serial.Port) *MemTimer {
	return &MemTimer{
		baseTimer: newBaseTimer(port, time.Second),
	}
}

// Start 启动定时器
func (t *MemTimer) Start() {
	if t.isRunning() {
		return
	}
	t.setRunning(true)

	// 立即发送一次（发送前检查状态）
	if t.isRunning() {
		send(t.port, proto.CmdMem, hardware.GetMemInfo())
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
					send(t.port, proto.CmdMem, hardware.GetMemInfo())
				}
			}
		}
	}()
}

// Stop 停止定时器
func (t *MemTimer) Stop() {
	t.stop()
}
