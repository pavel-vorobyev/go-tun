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

	out, err := util.RunCommand(fmt.Sprintf("sudo ip addr add %s/24 dev %s", ip, iface.Name()))
	if err != nil {
		log.Println(out, err)
		return nil, err
	}

	out, err = util.RunCommand(fmt.Sprintf("sudo ip link set dev %s up", iface.Name()))
	if err != nil {
		log.Println(out, err)
		return nil, err
	}

	return iface, nil
}

func CreateListener() (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", "192.168.50.27:8933")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return net.ListenUDP("udp", addr)
}

func ListenIface(iface *water.Interface, listener *net.UDPConn) {
	packet := make([]byte, 65535)

	for {
		n, err := iface.Read(packet)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("packet from iface")

		if udpAddr != nil {
			//_, err := conn.Write(packet[:n])
			_, err := listener.WriteToUDP(packet[:n], udpAddr)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func Listen(iface *water.Interface, listener *net.UDPConn) {
	for {
		packet := make([]byte, 65535)

		n, addr, err := listener.ReadFromUDP(packet)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("packet from client")

		if udpAddr == nil {
			udpAddr = addr
		}

		if iface != nil {
			_, err := iface.Write(packet[:n])
			if err != nil {
				log.Println(err)
			}
		}
	}
}

var udpAddr *net.UDPAddr

func main() {
	iface, err := CreateTun("192.168.9.11")
	if err != nil {
		log.Println(err)
		os.Exit(5)
	}

	listener, err := CreateListener()
	if err != nil {
		log.Println(err)
		os.Exit(5)
	}

	go ListenIface(iface, listener)
	go Listen(iface, listener)

	util.Serve()
}
