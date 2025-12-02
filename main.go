package main

import (
	"fmt"
	"log"
	"monitor/hardware"
	"monitor/pkg"
	"monitor/proto"
	"runtime"
	"sync"
	"time"

	"github.com/tarm/serial"
)

const (
	LogLineRegex = `^\[(\d+/\d+/\d+\s+\d+:\d+:\d+)\]\s+(.*)$`
) // Arduino日志行的正则表达式

// 显示模式
const (
	ModeCPU = iota
	ModeMem
	ModeNet
	ModeCount // 模式总数
)

var (
	currentMode = ModeCPU
	modeMutex   sync.Mutex
	serialPort  *serial.Port
	refreshChan = make(chan struct{}, 1) // 立即刷新通道
)

func main() {
	var path string
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("Running on macOS")
	case "linux":
		path = "/dev/ttyACM0"
	case "windows":
		path = "COM14"
	default:
		fmt.Printf("Unknown operating system: %s\n", os)
	}
	// 打开串口
	c := &serial.Config{Name: path, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	serialPort = s
	go serialListener(s)

	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			sendCurrentModeData(err, s)
		case <-refreshChan:
			sendCurrentModeData(err, s)
		}
	}
}

// 发送当前模式数据
func sendCurrentModeData(err error, s *serial.Port) {
	modeMutex.Lock()
	mode := currentMode
	modeMutex.Unlock()

	switch mode {
	case ModeCPU:
		cpuInfo2(err, s)
	case ModeMem:
		memInfo(err, s)
	case ModeNet:
		netInfo(err, s)
	}
}

// 切换到下一个模式
func nextMode() {
	modeMutex.Lock()
	currentMode = (currentMode + 1) % ModeCount
	modeMutex.Unlock()
	pkg.Log.Printf("Mode changed to: %d", currentMode)

	// 立即刷新显示
	select {
	case refreshChan <- struct{}{}:
	default:
	}
}

// 切换到上一个模式
func prevMode() {
	modeMutex.Lock()
	currentMode = (currentMode - 1 + ModeCount) % ModeCount
	modeMutex.Unlock()
	pkg.Log.Printf("Mode changed to: %d", currentMode)

	// 立即刷新显示
	select {
	case refreshChan <- struct{}{}:
	default:
	}
}

func serialListener(s *serial.Port) {
	buf := make([]byte, 256)
	lineBuf := make([]byte, 0, 256)

	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Printf("Serial read error: %v", err)
			continue
		}
		if n == 0 {
			continue
		}

		for i := 0; i < n; i++ {
			b := buf[i]

			// 检测到协议起始位
			if b == proto.STX {
				// 先处理之前的日志数据
				if len(lineBuf) > 0 {
					pkg.Log.Printf("<-arduino %s", string(lineBuf))
					lineBuf = lineBuf[:0]
				}

				// 读取完整数据包: STX已读取，还需要 CMD + LEN
				// 等待足够的数据
				remaining := buf[i:]
				if len(remaining) >= 4 {
					command := remaining[1]
					dataLength := int(remaining[2])
					expectedLen := 3 + dataLength + 1 // STX + CMD + LEN + DATA + CHECKSUM

					if len(remaining) >= expectedLen {
						dataBlock := remaining[3 : 3+dataLength]
						receivedChecksum := remaining[3+dataLength]
						expectedChecksum := proto.CalculateChecksum(dataBlock)

						pkg.Log.Printf("Packet: %02X,%02X,%02X,%s,%02X", proto.STX, command, dataLength, string(dataBlock), receivedChecksum)

						if receivedChecksum == expectedChecksum {
							switch command {
							case proto.CmdPageNext:
								pkg.Log.Printf("Received PageNext command")
								nextMode()
							case proto.CmdPagePrev:
								pkg.Log.Printf("Received PagePrev command")
								prevMode()
							}
						}
						i += expectedLen - 1 // 跳过已处理的字节
						continue
					}
				}
			}

			// 普通字符，累积到行缓冲
			if b == '\n' {
				if len(lineBuf) > 0 {
					pkg.Log.Printf("<-arduino %s", string(lineBuf))
					lineBuf = lineBuf[:0]
				}
			} else if b != '\r' {
				lineBuf = append(lineBuf, b)
			}
		}
	}
}

func cpuInfo2(err error, s *serial.Port) {
	cpu := hardware.GetCPUInfo()
	bytes := proto.BuildMsg(proto.CmdCPU, cpu)
	_, err = s.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func netInfo(err error, s *serial.Port) {
	net, err := hardware.Net()
	bytes := proto.BuildMsg(proto.CmdNet, net)
	// 发送数据到串口
	_, err = s.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func memInfo(err error, s *serial.Port) {
	mem := hardware.GetMemInfo()
	bytes := proto.BuildMsg(proto.CmdMem, mem)
	_, err = s.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
