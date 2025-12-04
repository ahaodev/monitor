package timer

import (
	"monitor/proto"
	"sync"
	"time"

	"github.com/tarm/serial"
)

// Mode 显示模式
type Mode int

const (
	ModeCPU Mode = iota
	ModeMem
	ModeNet
	ModeDisk
	ModeClock
	ModeCount
)

// Timer 定时器接口
type Timer interface {
	Start()
	Stop()
}

// Manager 定时器管理器
type Manager struct {
	port        *serial.Port
	currentMode Mode
	timers      map[Mode]Timer
	mu          sync.Mutex
}

// NewManager 创建定时器管理器
func NewManager(port *serial.Port) *Manager {
	m := &Manager{
		port:        port,
		currentMode: ModeClock,
		timers:      make(map[Mode]Timer),
	}

	// 初始化各类型定时器
	m.timers[ModeCPU] = NewCPUTimer(port)
	m.timers[ModeMem] = NewMemTimer(port)
	m.timers[ModeNet] = NewNetTimer(port)
	m.timers[ModeDisk] = NewDiskTimer(port)
	m.timers[ModeClock] = NewClockTimer(port)

	return m
}

// Start 启动当前模式定时器
func (m *Manager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t, ok := m.timers[m.currentMode]; ok {
		t.Start()
	}
}

// Stop 停止当前模式定时器
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t, ok := m.timers[m.currentMode]; ok {
		t.Stop()
	}
}

// HandleCommand 处理串口命令，返回是否处理成功
func (m *Manager) HandleCommand(cmd byte) bool {
	switch cmd {
	case proto.CmdPageNext:
		m.switchMode(1)
		return true
	case proto.CmdPagePrev:
		m.switchMode(-1)
		return true
	default:
		return false
	}
}

// NextPage 切换到下一页
func (m *Manager) NextPage() {
	m.switchMode(1)
}

// PrevPage 切换到上一页
func (m *Manager) PrevPage() {
	m.switchMode(-1)
}

// switchMode 内部切换模式方法
func (m *Manager) switchMode(delta int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 停止当前定时器
	if t, ok := m.timers[m.currentMode]; ok {
		t.Stop()
	}

	// 切换模式
	m.currentMode = Mode((int(m.currentMode) + delta + int(ModeCount)) % int(ModeCount))

	// 启动新定时器
	if t, ok := m.timers[m.currentMode]; ok {
		t.Start()
	}
}

// SwitchMode 切换模式（保留兼容）
func (m *Manager) SwitchMode(delta int) {
	m.switchMode(delta)
}

// CurrentMode 获取当前模式
func (m *Manager) CurrentMode() Mode {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.currentMode
}

// send 发送数据到串口
func send(port *serial.Port, cmd byte, data string) {
	if port == nil {
		return
	}
	port.Write(proto.BuildMsg(cmd, data))
}

// baseTimer 基础定时器
type baseTimer struct {
	port     *serial.Port
	interval time.Duration
	stopCh   chan struct{}
	running  bool
	mu       sync.Mutex
}

func newBaseTimer(port *serial.Port, interval time.Duration) *baseTimer {
	return &baseTimer{
		port:     port,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (t *baseTimer) stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		close(t.stopCh)
		t.running = false
		t.stopCh = make(chan struct{})
	}
}

func (t *baseTimer) isRunning() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running
}

func (t *baseTimer) setRunning(r bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.running = r
}
