package handlepacket

import (
	"container/list"
	"net"
)

type PacketInfo struct {
	Protocol  string
	State     string
	SrcMAC    net.HardwareAddr
	DstMAC    net.HardwareAddr
	SrcIP     net.IP
	DstIP     net.IP
	SrcPort   uint16
	DstPort   uint16
	TimeStamp int64 //时间戳
	DataSize  int   //这个包的大小
}

type MapValue struct {
	List       list.List
	DataSize   int
	RecordSize int //在一定时间内的包的总大小记录，会被其他线程修改
}
