package network

import (
	"fmt"
	"go-tun/util"
)

type Tun struct {
	name       string
	ip         string
	cidr       int
	mtu        int
	iface      *Interface
	in         chan []byte
	out        chan []byte
	readPacket []byte
}

func CreateTun(c *Config) (*Tun, error) {
	iface, err := NewInterface(c.Name, c.Mtu)
	if err != nil {
		return nil, err
	}

	tun := &Tun{
		name:       c.Name,
		ip:         c.Ip,
		cidr:       c.Cidr,
		mtu:        c.Mtu,
		iface:      iface,
		in:         make(chan []byte),
		out:        make(chan []byte),
		readPacket: make([]byte, c.Mtu*2),
	}

	err = tun.up()
	if err != nil {
		return nil, err
	}

	return tun, nil
}

func (tun *Tun) Start() {
	//go func() {
	//	packet := make([]byte, tun.mtu*2)
	//	for {
	//		n, err := tun.iface.Read(packet)
	//		if err != nil {
	//			//log.Println(fmt.Sprintf("TUN: failed to read packet: %s", err))
	//			continue
	//		}
	//		tun.out <- packet[:n]
	//	}
	//}()
	//go func() {
	//	for {
	//		packet := <-tun.in
	//		_, err := tun.iface.Write(packet)
	//		if err != nil {
	//			//log.Println(fmt.Sprintf("TUN: failed to write packet: %s", err))
	//			continue
	//		}
	//	}
	//}()
}

func (tun *Tun) Receive() ([]byte, error) {
	n, err := tun.iface.Read(tun.readPacket)
	if err != nil {
		//log.Println(fmt.Sprintf("TUN: failed to read packet: %s", err))
		return nil, err
	}
	return tun.readPacket[:n], err
}

func (tun *Tun) Send(data []byte) {
	_, err := tun.iface.Write(data)
	if err != nil {
		//log.Println(fmt.Sprintf("TUN: failed to write packet: %s", err))
	}
}

func (tun *Tun) up() error {
	_, err := util.RunCommand(fmt.Sprintf("sudo ip link set dev %s mtu %d", tun.name, tun.mtu))
	if err != nil {
		return err
	}
	_, err = util.RunCommand(fmt.Sprintf("sudo ip addr add %s/%d dev %s", tun.ip, tun.cidr, tun.name))
	if err != nil {
		return err
	}
	_, err = util.RunCommand(fmt.Sprintf("sudo ip link set dev %s up", tun.name))
	if err != nil {
		return err
	}
	return nil
}
