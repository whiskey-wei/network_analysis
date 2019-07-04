package handlepacket

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func GetInfo(packet gopacket.Packet) {
	//ethernet
	info := PacketInfo{}

	info.DataSize = len(packet.Data())

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetInfo := ethernetLayer.(*layers.Ethernet)
		info.SrcMAC = ethernetInfo.SrcMAC
		info.DstMAC = ethernetInfo.DstMAC
	}

	//IP
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ipInfo, _ := ipLayer.(*layers.IPv4)
		//fmt.Println("From ", ipInfo.SrcIP, " To ", ipInfo.DstIP)
		//fmt.Println("Protocol: ", ipInfo.Protocol)
		info.SrcIP = ipInfo.SrcIP
		info.DstIP = ipInfo.DstIP
	}

	//TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcpInfo, _ := tcpLayer.(*layers.TCP)
		//fmt.Println("From ", tcpInfo.SrcPort, " To ", tcpInfo.DstPort)
		info.SrcPort = uint16(tcpInfo.SrcPort)
		info.DstPort = uint16(tcpInfo.DstPort)
		info.Protocol = "tcp"
		if tcpInfo.SYN == true { //SYN-SENT
			info.State = "SYN-SENT"
		} else if tcpInfo.SYN == true && tcpInfo.ACK == true {
			info.State = "SYN-RECV"
		} else if tcpInfo.FIN == true {
			//将连接从链表清除
			return
		} else {
			info.State = "ESTABLISH"
		}
		info.TimeStamp = time.Now().Unix()

		//add to tcp hashlist
	}

	//UDP
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udpInfo, _ := udpLayer.(*layers.UDP)
		info.SrcPort = uint16(udpInfo.SrcPort)
		info.DstPort = uint16(udpInfo.DstPort)
		info.Protocol = "udp"
		info.TimeStamp = time.Now().Unix()

		//add to udp hashlist

	}

	fmt.Println(info)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	time.Sleep(2 * time.Second)
	cmd.Run()
}
