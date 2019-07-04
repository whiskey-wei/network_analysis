package main

import (
	"fmt"

	"github.com/google/gopacket/pcap"
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, device := range devices {
		for _, addr := range device.Addresses {
			fmt.Println("device: ", device.Name)
			fmt.Println("ip addr: ", addr.IP)
			fmt.Println("mask: ", addr.Netmask)
		}
	}
}
