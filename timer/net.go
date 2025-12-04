package timer

import (
	"monitor/hardware"
	"monitor/proto"
	"time"

	"github.com/tarm/serial"
)

// NetTimer 网络定时器
type NetTimer struct {
	*baseTimer
}

// NewNetTimer 创建网络定时器
func NewNetTimer(port *serial.Port) *NetTimer {
	return &NetTimer{
		baseTimer: newBaseTimer(port, 2*time.Second), // 网络统计需要1秒采样，间隔2秒
	}
}

// Start 启动定时器
func (t *NetTimer) Start() {
	if t.isRunning() {
		return
	}
	t.setRunning(true)

	// 立即发送一次（发送前检查状态）
	if data, err := hardware.Net(); err == nil && t.isRunning() {
		send(t.port, proto.CmdNet, data)
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
				if data, err := hardware.Net(); err == nil && t.isRunning() {
					send(t.port, proto.CmdNet, data)
				}
			}
		}
	}()
}

// Stop 停止定时器
func (t *NetTimer) Stop() {
	t.stop()
}
