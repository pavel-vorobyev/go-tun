package util

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"go-tun/constant"
)

// TODO Create the struct wrapper for the Packet entity

const (
	TcpId = 6
	UdpId = 17
)

func ParsePacket(data []byte) (int, string, string, error) {
	packet := gopacket.NewPacket(data, layers.LayerTypeIPv4, gopacket.Default)
	network := packet.NetworkLayer().NetworkFlow()
	transport := packet.TransportLayer().TransportFlow()

	src := fmt.Sprintf("%s:%s", network.Src().String(), transport.Src().String())
	dst := fmt.Sprintf("%s:%s", network.Dst().String(), transport.Dst().String())

	if layer := packet.Layer(layers.LayerTypeTCP); layer != nil {
		return TcpId, src, dst, nil
	} else if layer = packet.Layer(layers.LayerTypeUDP); layer != nil {
		return UdpId, src, dst, nil
	} else {
		return 0, constant.StringEmpty, constant.StringEmpty, fmt.Errorf("unknown protocol")
	}
}

// TODO Create separate methods that return Src and Dst
// TODO Implement them the way Flows are extracted once
