package main

import (
	"log"
	"monitor/proto"
	"monitor/timer"

	"github.com/tarm/serial"
)

// 配置常量
const (
	SerialPort = "/dev/ttyACM0"
	SerialBaud = 115200
)

var manager *timer.Manager

func main() {
	s, err := serial.OpenPort(&serial.Config{Name: SerialPort, Baud: SerialBaud})
	if err != nil {
		log.Fatal(err)
	}

	// 创建定时器管理器
	manager = timer.NewManager(s)

	// 启动串口监听
	go serialListener(s)

	// 启动当前模式定时器
	manager.Start()

	// 保持主goroutine运行
	select {}
}

func serialListener(s *serial.Port) {
	buf := make([]byte, 256)

	for {
		n, err := s.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		// 打印收到的原始数据
		log.Printf("recv<- %x (%s)", buf[:n], string(buf[:n]))

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
				log.Printf("parsed-> cmd: %x, data: %s", cmd, string(data))
				handleCommand(cmd)
			}
			i += expectedLen - 1
		}
	}
}

func handleCommand(cmd byte) {
	switch cmd {
	case proto.CmdPageNext:
		manager.SwitchMode(1)
	case proto.CmdPagePrev:
		manager.SwitchMode(-1)
	}
}
