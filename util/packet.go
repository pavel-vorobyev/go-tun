package util

import (
	"fmt"
	"log"
)

const (
	ProtoUdp = 6
	ProtoTcp = 17
)

func GetPacketBaseInfo(d []byte) (int, string, string) {
	return GetPacketProtocol(d), GetPacketSrc(d), GetPacketDst(d)
}

func GetPacketProtocol(d []byte) int {
	log.Println(len(d))
	return int(d[9])
}

func GetPacketSrc(d []byte) string {
	addr := d[12:16]
	return fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3])
}

func GetPacketDst(d []byte) string {
	addr := d[16:20]
	return fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3])
}
