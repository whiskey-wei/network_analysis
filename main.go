package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/whiskey-wei/network_analysis/conf"
	"github.com/whiskey-wei/network_analysis/handlepacket"

	"github.com/google/gopacket/pcap"
)

//命令行解析
func mainInit() {
	Protocol := flag.String("p", "tcp or udp", "protocol type tcp or udp")
	Port := flag.Int("i", 0, "port")
	CatchTime := flag.Int64("t", 10, "time of catch packet")
	DeviceName := flag.String("d", "en0", "device name")
	SnapShotLen := flag.Int("s", 65536, "max length of packet")
	UdpFilePath := flag.String("up", "./info/udpinfo", "udp file path")
	TcpFilePath := flag.String("tp", "./info/tcpinfo", "tcp file path")
	flag.Parse()
	conf.Protocol = *Protocol
	conf.Port = *Port
	conf.CatchTime = *CatchTime
	conf.DeviceName = *DeviceName
	conf.SnapShotLen = *SnapShotLen
	conf.UdpFilePath = *UdpFilePath
	conf.TcpFilePath = *TcpFilePath
}

func main() {
	mainInit()

	handle, err := pcap.OpenLive(conf.DeviceName, int32(conf.SnapShotLen), true, time.Second*30)
	defer handle.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = handle.SetBPFFilter(getFilter(conf.Protocol, conf.Port))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	//开启线程打印相关信息
	go handlepacket.ShowInfo()

	//统计固定时间内数据量,起一个线程每隔一段时间去重置

	//实时获取数据包，并存放只哈希表
	for packet := range packetSource.Packets() {
		handlepacket.GetInfo(packet)
	}

}

//设置过滤器
func getFilter(proto string, port int) string {
	if proto != "tcp" && proto != "udp" && proto != "tcp or udp" {
		fmt.Println("do not support this protocol")
		os.Exit(0)
	}
	if port != 0 {
		return fmt.Sprintf("%v and ((src port %v) or (dst port %v))", proto, port, port)
	} else {
		return fmt.Sprintf("%v", proto)
	}
}
