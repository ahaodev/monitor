package proto

//CPU 温度、总负载
//SYSTEM 系统信息 32GB DDR5,RAM 12%,Flash 4% ,LOG 1%,Docker 56%
//INTERFACE In 485kbps ,Out 88.9kbps
//DISK ,DISK1 ,39 C

const (
	STX          byte   = 0xAA
	StartByteErr byte   = 0xAB
	CmdData      byte   = 0x01
	CmdAck       byte   = 0x02
	Ack          string = "ACK"
	Err          string = "ERR"
)

// 前16个命令为系统命令
const (
	CmdInit  byte = 0x01 //初始化
	CmdSleep byte = 0x02 //休眠
	CmdACK   byte = 0x03 //确认
	CmdErr   byte = 0x04 //错误
	CmdReset byte = 0x05 //重置

	CmdPageNext byte = 0x0A //下一页 (手势触发)
	CmdPagePrev byte = 0x0B //上一页 (手势触发)

	CmdCPU  byte = 0x10 //CPU
	CmdMem  byte = 0x11 //内存
	CmdNet  byte = 0x12 //网络
	CmdDisk byte = 0x13 //磁盘
)
