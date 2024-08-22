package transport

import (
	"fmt"
	"log"
	"net"
)

type UDPConn struct {
	conn *net.UDPConn
	mtu  int
	in   chan *Data
	out  chan *Data
}

func CreateConn(c *Config) (*UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.Ip, c.Port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &UDPConn{
		conn: conn,
		mtu:  c.Mtu,
		in:   make(chan *Data),
		out:  make(chan *Data),
	}, nil
}

func (conn *UDPConn) Start() {
	go func() {
		packet := make([]byte, conn.mtu)
		for {
			n, addr, err := conn.conn.ReadFromUDP(packet)
			if err != nil {
				log.Println(fmt.Sprintf("UDP: failed to read packet: %s", err))
				continue
			}

			cAddr := fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port)
			conn.out <- &Data{
				Data:  packet[:n],
				CAddr: cAddr,
			}
		}
	}()
	go func() {
		for {
			data := <-conn.in

			addr, err := net.ResolveUDPAddr("udp", data.CAddr)
			if err != nil {
				log.Println(fmt.Sprintf("UDP: failed to resolve address: %s", err))
				continue
			}

			_, err = conn.conn.WriteToUDP(data.Data, addr)
			if err != nil {
				log.Println(fmt.Sprintf("UDP: failed to write packet: %s", err))
				continue
			}
		}
	}()
}

func (conn *UDPConn) Receive() *Data {
	return <-conn.out
}

func (conn *UDPConn) Send(data *Data) {
	conn.in <- data
}
