package simple

import (
	"fmt"
	"github.com/songgao/water"
	"go-tun/util"
	"runtime"
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
	switch runtime.GOOS {
	case "darwin":
		_, err := util.RunCommand(fmt.Sprintf("sudo ifconfig %s inet %s/8 %s alias", tun.Name, tun.Ip, tun.Ip))
		if err != nil {
			return err
		}
		_, err = util.RunCommand(fmt.Sprintf("sudo ifconfig %s up", tun.Name))
		if err != nil {
			return err
		}
	case "linux":
		_, err := util.RunCommand(fmt.Sprintf("sudo ip addr add %s/24 dev %s", tun.Ip, tun.Name))
		if err != nil {
			return err
		}
		_, err = util.RunCommand(fmt.Sprintf("sudo ip link set dev %s up", tun.Name))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported platform")
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
