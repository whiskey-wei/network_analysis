package handlepacket

import (
	"fmt"
	"os"
	"os/exec"
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
		if tcpInfo.SYN == true && tcpInfo.ACK == true { //SYN-SENT
			info.State = "SYN-SENT"
		} else if tcpInfo.SYN == true {
			info.State = "SYN-RECV"
		} else if tcpInfo.FIN == true {
			//将连接从链表清除
			TcpMap.Del(info)
			SaveInfo(&info, conf.TcpFilePath)
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

//信息打印
func ShowInfo() {
	for {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		time.Sleep(time.Second * 2)
		cmd.Run()
		fmt.Printf("%v%18v%20v%15v%15v%15v%15v%vs)%15v%15v%vs)%15v\n", "Proto", "SrcIP", "DstIP", "SrcPort", "DstPort", "State", "PreSize(", conf.CatchTime, "NowSize", "PreCount(", conf.CatchTime, "NowCount")
		for i := 0; i < 160; i++ {
			fmt.Printf("-")
		}
		fmt.Printf("\n")
		//fmt.Println("Proto\tSrcIP\tDstIP\tSrcPort\tDstPort\tState")
		SumMap.Range(printInfo)
	}
}

func printInfo(k, v interface{}) bool {
	mapVal, ok := v.(*MapValue)
	if !ok {
		return false
	}
	printVal(mapVal)
	return true
}

func printVal(mapVal *MapValue) {
	nodeVal, ok := mapVal.List.Front().Value.(*PacketInfo)
	if ok {
		HashMux.Lock()
		fmt.Printf("%v%20v%20v%15v%15v%15v%15vB%15vB%15v%15v\n",
			nodeVal.Protocol, nodeVal.SrcIP, nodeVal.DstIP, nodeVal.SrcPort,
			nodeVal.DstPort, nodeVal.State, mapVal.PreSize, mapVal.NowSize,
			mapVal.PreCount, mapVal.NowCount)
		HashMux.Unlock()
	}
}

func ResetDataSize() {
	for {
		time.Sleep(time.Duration(conf.CatchTime) * time.Second)
		SumMap.Range(resetDataSize)
	}
}

func resetDataSize(k, v interface{}) bool {
	mapVal, ok := v.(*MapValue)
	if !ok {
		return false
	}
	HashMux.Lock()
	mapVal.PreSize = mapVal.NowSize
	mapVal.NowSize = 0
	mapVal.PreCount = mapVal.NowCount
	mapVal.NowCount = 0
	HashMux.Unlock()
	return true
}
