package handlepacket

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/whiskey-wei/network_analysis/conf"
)

func ShowInfo() {
	for {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		time.Sleep(time.Second * 2)
		cmd.Run()
		fmt.Printf("%v%18v%20v%20v%20v%20v%20v%vs)\n", "Proto", "SrcIP", "DstIP", "SrcPort", "DstPort", "State", "DataSize(", conf.CatchTime)
		for i := 0; i < 130; i++ {
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
		fmt.Printf("%v%20v%20v%20v%20v%20v%20vB\n", nodeVal.Protocol, nodeVal.SrcIP, nodeVal.DstIP, nodeVal.SrcPort, nodeVal.DstPort, nodeVal.State, mapVal.NowSize)
	}
}
