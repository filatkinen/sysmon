package main

import (
	"fmt"
	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	// The same default as tcpdump.
	defaultSnapLen = 262144
)

func main() {
	handle, err := pcap.OpenLive("wlp4s0:", defaultSnapLen, true,
		pcap.BlockForever)
	if err != nil {
		panic(err)
	}
	defer handle.Close()

	//if err := handle.SetBPFFilter("port 3030"); err != nil {
	//	panic(err)
	//}

	packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	for pkt := range packets {
		fmt.Println(pkt.Metadata().CaptureInfo)
	}
}
