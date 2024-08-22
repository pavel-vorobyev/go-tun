package network

import (
	"fmt"
	"github.com/songgao/water"
	"go-tun/util"
)

type Tun struct {
	name string
	ip   string
	cidr int
	mtu  int

	iface *water.Interface
	in    chan []byte
	out   chan []byte
}

func CreateTun(c *Config) (*Tun, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = c.Name

	iface, err := water.New(config)
	if err != nil {
		return nil, err
	}

	tun := &Tun{
		name:  iface.Name(),
		ip:    c.Ip,
		cidr:  c.Cidr,
		mtu:   c.Mtu,
		iface: iface,
		in:    make(chan []byte),
		out:   make(chan []byte),
	}

	err = tun.up()
	if err != nil {
		return nil, err
	}

	return tun, nil
}

func (tun *Tun) Start() {
	go func() {
		packet := make([]byte, tun.mtu*2)
		for {
			n, err := tun.iface.Read(packet)
			if err != nil {
				//log.Println(fmt.Sprintf("TUN: failed to read packet: %s", err))
				continue
			}
			tun.out <- packet[:n]
		}
	}()
	go func() {
		for {
			packet := <-tun.in
			_, err := tun.iface.Write(packet)
			if err != nil {
				//log.Println(fmt.Sprintf("TUN: failed to write packet: %s", err))
				continue
			}
		}
	}()
}

func (tun *Tun) Receive() []byte {
	return <-tun.out
}

func (tun *Tun) Send(data []byte) {
	tun.in <- data
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
