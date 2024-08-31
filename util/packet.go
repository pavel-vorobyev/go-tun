package util

import (
	"fmt"
)

const (
	ProtoUdp = 6
	ProtoTcp = 17
)

func GetPacketBaseInfo(d []byte, n int) (int, string, string, error) {
	if n < 20 {
		return 0, "", "", fmt.Errorf("too short")
	}
	return GetPacketProtocol(d), GetPacketSrc(d), GetPacketDst(d), nil
}

func GetPacketProtocol(d []byte) int {
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
