package proto

import (
	"log"
)

// Pack struct
// | 起始位 | 命令位 | 数据长度位 | 数据块 | 校验位 |
//
// - **起始位**: 固定值，用于标识消息的开始。
// - **命令位**: 标识消息类型的一个字节。
// - **数据长度位**: 标识数据块长度的一个字节。
// - **数据块**: 实际传输的数据，长度可变。
// - **校验位**: 校验和或其他校验机制，用于确保数据完整性。
type Pack struct {
	Header byte
	Cmd    byte
	Length uint16
	Data   []byte
	CRC    uint32
}

// BuildMsg 构建协议报文
func BuildMsg(command byte, data string) []byte {
	log.Printf("send->%x %s", command, data)
	dataLength := len(data)
	message := []byte{STX, command, byte(dataLength)}
	message = append(message, []byte(data)...)
	checksum := CalculateChecksum([]byte(data))
	message = append(message, checksum)
	log.Printf(" %x %x %x %x %x", STX, command, dataLength, data, checksum)
	return message
}

// CalculateChecksum 异或校验
func CalculateChecksum(data []byte) byte {
	var checksum byte
	for _, b := range data {
		checksum ^= b
	}
	return checksum
}
