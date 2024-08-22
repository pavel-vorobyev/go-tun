package simple

import (
	"fmt"
	"github.com/songgao/water"
	"go-tun/util"
)

type Tun struct {
	Ip   string
	Name string
	Mtu  int

	Interface *water.Interface

	In  chan []byte
	Out chan []byte
}

func CreateTun(ip string, name string, mtu int) (*Tun, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name

	iface, err := water.New(config)
	if err != nil {
		return nil, err
	}

	tun := &Tun{
		Ip:        ip,
		Name:      iface.Name(),
		Interface: iface,
		Mtu:       mtu,
	}

	err = tun.up()
	if err != nil {
		return nil, err
	}
	return tun, nil
}

func (tun *Tun) up() error {
	_, err := util.RunCommand(fmt.Sprintf("sudo ip link set dev %s mtu %d", tun.Name, tun.Mtu))
	if err != nil {
		return err
	}
	_, err = util.RunCommand(fmt.Sprintf("sudo ip addr add %s/24 dev %s", tun.Ip, tun.Name))
	if err != nil {
		return err
	}
	_, err = util.RunCommand(fmt.Sprintf("sudo ip link set dev %s up", tun.Name))
	if err != nil {
		return err
	}
	return nil
}

//func (network *NTun) listen() {
//	go func() {
//		data := make([]byte, network.mtu)
//		for {
//
//		}
//	}()
//}
