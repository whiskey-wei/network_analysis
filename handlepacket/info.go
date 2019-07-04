package handlepacket

import (
	"fmt"
	"os"
	"time"

	"github.com/whiskey-wei/network_analysis/conf"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func GetInfo(packet gopacket.Packet) {

	info := PacketInfo{}
	//fmt.Println(packet)
	info.DataSize = len(packet.Data())
	//ethernet
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
			TcpMap.Del(info)
			return
		} else {
			info.State = "ESTABLISH"
		}
		info.TimeStamp = time.Now().Unix()
		info.Ack = tcpInfo.Ack
		info.Seq = tcpInfo.Seq
		//add to tcp hashlist
		TcpMap.Add(&info)
		SumMap.Add(&info)
		SaveInfo(&info, conf.TcpFilePath)
		return
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
		UdpMap.Add(&info)
		SumMap.Add(&info)
		SaveInfo(&info, conf.UdpFilePath)
		return
	}

}

//将包信息存放到磁盘文件
func SaveInfo(info *PacketInfo, path string) {

	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	if info.Protocol == "tcp" {
		f.WriteString(fmt.Sprintf("protocol:%v\tSrcIP:%v\tDstIP:%v\tSrcPort:%v\tDstPort:%v\tState:%v\tDataSize:%vB\tSeq:%v\tAck:%v\tTime:%v\n",
			info.Protocol, info.SrcIP, info.DstIP, info.SrcPort, info.DstPort, info.State, info.DataSize, info.Seq, info.Ack, info.TimeStamp))
	} else {
		f.WriteString(fmt.Sprintf("protocol:%v\tSrcIP:%v\tDstIP:%v\tSrcPort:%v\tDstPort:%v\tDataSize:%vB\tTime:%v\n",
			info.Protocol, info.SrcIP, info.DstIP, info.SrcPort, info.DstPort, info.DataSize, info.TimeStamp))
	}

}
