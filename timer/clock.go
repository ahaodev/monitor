package timer

import (
	"fmt"
	"monitor/proto"
	"time"

	"github.com/tarm/serial"
)

// ClockTimer 时钟定时器
type ClockTimer struct {
	*baseTimer
}

// NewClockTimer 创建时钟定时器
func NewClockTimer(port *serial.Port) *ClockTimer {
	return &ClockTimer{
		baseTimer: newBaseTimer(port, 15*time.Minute), // 时钟每15分钟同步一次
	}
}

// Start 启动定时器
func (t *ClockTimer) Start() {
	if t.isRunning() {
		return
	}
	t.setRunning(true)

	// 立即发送一次（发送前检查状态）
	if t.isRunning() {
		t.sendClock()
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
					t.sendClock()
				}
			}
		}
	}()
}

// Stop 停止定时器
func (t *ClockTimer) Stop() {
	t.stop()
}

func (t *ClockTimer) sendClock() {
	now := time.Now()
	send(t.port, proto.CmdClock, fmt.Sprintf("%d,%d,%d", now.Hour(), now.Minute(), now.Second()))
}
