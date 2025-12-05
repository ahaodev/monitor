package main

import (
	"fmt"
	"monitor/hardware"
	"monitor/proto"
	"time"

	"github.com/tarm/serial"
)

// 配置常量
const (
	SerialPort = "/dev/ttyACM0"
	SerialBaud = 115200
)

// ANSI颜色码
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

var serialPort *serial.Port

// 命令名称映射
func getCmdName(cmd byte) string {
	switch cmd {
	case proto.CmdCPU:
		return "CPU"
	case proto.CmdMem:
		return "MEM"
	case proto.CmdNet:
		return "NET"
	case proto.CmdDisk:
		return "DISK"
	case proto.CmdClock:
		return "CLOCK"
	default:
		return fmt.Sprintf("0x%02X", cmd)
	}
}

// 日志输出
func logInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[INFO]%s  %s\n", colorGreen, colorReset, msg)
}

func logRecv(cmd byte) {
	name := getCmdName(cmd)
	fmt.Printf("%s[RECV]%s  ← %s%s%s request\n", colorCyan, colorReset, colorYellow, name, colorReset)
}

func logSend(cmd byte, data string) {
	name := getCmdName(cmd)
	// 截断过长的数据
	displayData := data
	if len(data) > 30 {
		displayData = data[:30] + "..."
	}
	fmt.Printf("%s[SEND]%s  → %s%s%s: %s%s%s\n", colorPurple, colorReset, colorYellow, name, colorReset, colorGray, displayData, colorReset)
}

func logError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[ERR]%s   %s\n", colorRed, colorReset, msg)
}

func main() {
	var err error

	// 打印启动横幅
	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║%s   LCD Monitor - Go Backend v1.3.3  %s║%s\n", colorCyan, colorReset, colorCyan, colorReset)
	fmt.Printf("%s╚════════════════════════════════════╝%s\n", colorCyan, colorReset)
	fmt.Println()

	logInfo("Opening serial port: %s @ %d baud", SerialPort, SerialBaud)

	// 串口打开失败时每3秒重试
	for {
		serialPort, err = serial.OpenPort(&serial.Config{Name: SerialPort, Baud: SerialBaud})
		if err == nil {
			break
		}
		logError("Failed to open serial port: %v", err)
		logInfo("Retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}

	logInfo("Serial port opened ✓")
	logInfo("Waiting for LCD requests...")
	fmt.Println()

	// 启动串口监听，响应LCD请求
	serialListener(serialPort)
}

func serialListener(s *serial.Port) {
	buf := make([]byte, 256)

	for {
		n, err := s.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		for i := 0; i < n; i++ {
			if buf[i] != proto.STX {
				continue
			}

			remaining := buf[i:]
			if len(remaining) < 4 {
				continue
			}

			cmd := remaining[1]
			dataLen := int(remaining[2])
			expectedLen := 4 + dataLen

			if len(remaining) < expectedLen {
				continue
			}

			data := remaining[3 : 3+dataLen]
			checksum := remaining[3+dataLen]

			if checksum == proto.CalculateChecksum(data) {
				logRecv(cmd)
				handleRequest(s, cmd)
			}
			i += expectedLen - 1
		}
	}
}

// handleRequest 处理LCD请求并响应数据
func handleRequest(s *serial.Port, cmd byte) {
	var response string

	switch cmd {
	case proto.CmdCPU:
		response = hardware.GetCPUInfo()
	case proto.CmdMem:
		response = hardware.GetMemInfo()
	case proto.CmdNet:
		if data, err := hardware.Net(); err == nil {
			response = data
		} else {
			response = "0,0"
		}
	case proto.CmdDisk:
		response = hardware.GetDiskInfoForDisplay()
	case proto.CmdClock:
		now := time.Now()
		weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
		response = fmt.Sprintf("%d,%d,%d,%s,%s", now.Hour(), now.Minute(), now.Second(), now.Format("2006-01-02"), weekdays[now.Weekday()])
	default:
		logError("Unknown command: 0x%02X", cmd)
		return
	}

	// 发送响应
	sendResponse(s, cmd, response)
}

// sendResponse 发送响应数据
func sendResponse(s *serial.Port, cmd byte, data string) {
	if s == nil {
		return
	}
	msg := proto.BuildMsg(cmd, data)
	s.Write(msg)
	logSend(cmd, data)
}
