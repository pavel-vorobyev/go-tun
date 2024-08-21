package transport

import "net"

type CAddr = net.UDPAddr

type Data struct {
	Data  []byte
	CAddr *CAddr
}
