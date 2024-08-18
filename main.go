package main

import (
	"fmt"
	"github.com/songgao/water"
	"go-tun/util"
	"log"
	"net"
	"os"
)

func CreateTun(ip string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(config)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//out, err := util.RunCommand(fmt.Sprintf("sudo ifconfig %s inet %s/8 %s alias", iface.Name(), ip, ip))
	out, err := util.RunCommand(fmt.Sprintf("sudo ip addr add %s/24 dev %s", ip, iface.Name()))
	if err != nil {
		log.Println(out, err)
		return nil, err
	}

	//out, err = util.RunCommand(fmt.Sprintf("sudo ifconfig %s up", iface.Name()))
	out, err = util.RunCommand(fmt.Sprintf("sudo ip link set dev %s up", iface.Name()))
	if err != nil {
		log.Println(out, err)
		return nil, err
	}

	log.Println(fmt.Sprintf("TUN started: %s", iface.Name()))

	return iface, nil
}

func CreateUdpListener() (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:5995")
	// addr, err := net.ResolveUDPAddr("udp", "159.100.30.164:8933")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return net.ListenUDP("udp", addr)
}

func ListenIface(iface *water.Interface, listener *net.UDPConn) {
	packet := make([]byte, 65535)

	for {
		_, err := iface.Read(packet)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("packet from tun")
	}
}

func ListenUdpConnection(iface *water.Interface, listener *net.UDPConn) {
	packet := make([]byte, 65535)

	for {
		n, _, err := listener.ReadFromUDP(packet)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("packet from client")

		if iface != nil {
			_, err := iface.Write(packet[:n])
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func main() {
	iface, err := CreateTun("10.0.0.2")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	listener, err := CreateUdpListener()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go ListenIface(iface, listener)
	go ListenUdpConnection(iface, listener)

	util.Serve()
}
