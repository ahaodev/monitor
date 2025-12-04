package proto

const (
	STX byte   = 0xAA
	Ack string = "ACK"
	Err string = "ERR"
)

const (
	CmdCPU   byte = 0x10 //CPU
	CmdMem   byte = 0x11 //内存
	CmdNet   byte = 0x12 //网络
	CmdDisk  byte = 0x13 //磁盘
	CmdClock byte = 0x14 //时钟
)
