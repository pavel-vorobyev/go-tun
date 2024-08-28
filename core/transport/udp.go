package transport

import (
	"fmt"
	"net"
)

type UDPConn struct {
	conn       *net.UDPConn
	mtu        int
	in         chan *Data
	out        chan *Data
	readPacket []byte
}

func NewConn(c *Config) (*UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.Ip, c.Port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &UDPConn{
		conn:       conn,
		mtu:        c.Mtu,
		in:         make(chan *Data),
		out:        make(chan *Data),
		readPacket: make([]byte, c.Mtu*2),
	}, nil
}

func (conn *UDPConn) Receive() (int, *Data, error) {
	n, addr, err := conn.conn.ReadFromUDP(conn.readPacket)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read from udp: %s", err)
	}
	return n, NewData(conn.readPacket[:n], addr.String()), nil
}

func (conn *UDPConn) Send(data *Data) error {
	addr, err := net.ResolveUDPAddr("udp", data.GetCAddr())
	if err != nil {
		return fmt.Errorf("failed to resolove udp addr: %s", err)
	}
	_, err = conn.conn.WriteToUDP(data.GetData(), addr)
	if err != nil {
		return fmt.Errorf("failed to write udp: %s", err)
	}
	return nil
}
