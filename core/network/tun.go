package network

import (
	"github.com/songgao/water"
)

type Tun struct {
	name string
	ip   string
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

	return &Tun{
		name:  iface.Name(),
		ip:    c.Ip,
		mtu:   c.Mtu,
		iface: iface,
		in:    make(chan []byte),
		out:   make(chan []byte),
	}, nil
}

func (tun *Tun) Start() {
	go func() {
		packet := make([]byte, tun.mtu)
		for {
			n, err := tun.iface.Read(packet)
			if err != nil {
				continue
			}
			tun.out <- packet[n:]
		}
	}()
	go func() {
		for {
			packet := <-tun.in
			_, err := tun.iface.Write(packet)
			if err != nil {
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
