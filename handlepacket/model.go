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
	Seq       uint32
	Ack       uint32
}

type MapValue struct {
	List    list.List
	PreSize int
	NowSize int //在一定时间内的包的总大小记录，会被其他线程修改
}
