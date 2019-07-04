package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket/layers"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	device       string        = "en0"
	snapshot_len int32         = 65536
	promiscuous  bool          = true
	timeout      time.Duration = 30 * time.Second
)

//输出包信息
func printPacketInfo(packet gopacket.Packet) {

	fmt.Print("\n\n-------------------------\n")
	//IP
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ipInfo, _ := ipLayer.(*layers.IPv4)
		fmt.Println("From ", ipInfo.SrcIP, " To ", ipInfo.DstIP)
		fmt.Println("Protocol: ", ipInfo.Protocol)
	}

	//TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcpInfo, _ := tcpLayer.(*layers.TCP)
		fmt.Println("From ", tcpInfo.SrcPort, " To ", tcpInfo.DstPort)
	}

	//UDP
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {

	}
}

func main() {
	handle, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	defer handle.Close()
	if err != nil {
		log.Fatal(err)
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		printPacketInfo(packet)
	}
}
