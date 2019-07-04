package main

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/whiskey-wei/network_analysis/handlepacket"

	"github.com/google/gopacket/pcap"
)

var (
	DeviceName  *string
	SnapShotLen *int
	TimeOut     *int64
	Type        *string
	Port        *string
)

//命令行解析
func mainInit() {
	DeviceName = flag.String("d", "en0", "device name")
	SnapShotLen = flag.Int("s", 65536, "max length of packet")
	TimeOut = flag.Int64("t", int64(30*time.Second), "timeout")

}

//判断输入网卡名是否正确
func checkDeviceName(deviceName *string) error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}

	for _, d := range devices {
		if d.Name == *deviceName {
			return nil
		}
	}

	return errors.New("Device Name Not Exist")
}

func main() {
	mainInit()
	flag.Parse()
	checkDeviceName(DeviceName)

	handle, err := pcap.OpenLive(*DeviceName, int32(*SnapShotLen), true, time.Duration(*TimeOut))
	defer handle.Close()
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		handlepacket.GetInfo(packet)
	}
}
