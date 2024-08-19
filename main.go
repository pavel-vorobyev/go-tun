package main

import (
	"fmt"
	"github.com/songgao/water"
	"go-tun/network"
	"go-tun/storage"
	"go-tun/util"
	"log"
	"net"

	"github.com/xitongsys/ethernet-go/header"
)

// TODO Constant peers count fix idea: set client's TUN IP to the packet on the client
// TODO Since client knows it's own TUN's IP address we can set it on the client
// TODO And replace on the server with real IP address provided by UDP server

func CreateUdpListener() (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:5995")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return net.ListenUDP("udp", addr)
}

func ListenUdp(iface *water.Interface, listener *net.UDPConn) {
	go func() {
		packet := make([]byte, 1500*2)
		for {
			n, uaddr, err := listener.ReadFromUDP(packet)
			if err != nil {
				log.Println(err)
				continue
			}

			protocol, src, dst, err := header.GetBase(packet[:n])
			if err != nil {
				log.Println(err)
				continue
			}

			_, err = iface.Write(packet[:n])
			if err != nil {
				log.Println(err)
				continue
			}

			key := fmt.Sprintf("%s/%s@%s", protocol, src, dst)
			saddr := fmt.Sprintf("%s:%d", uaddr.IP.String(), uaddr.Port)
			connections.Put(key, saddr)

			log.Println(fmt.Sprintf("i: %s ← %s", key, saddr))
		}
	}()
}

func ListenTun(iface *water.Interface, listener *net.UDPConn) {
	go func() {
		packet := make([]byte, 1500*2)

		for {
			n, err := iface.Read(packet)
			if err != nil {
				log.Println(err)
				continue
			}

			protocol, src, dst, err := header.GetBase(packet[:n])
			if err != nil {
				log.Println(err)
				continue
			}

			key := fmt.Sprintf("%s/%s@%s", protocol, dst, src)
			saddr, exists := connections.Get(key)
			if !exists {
				continue
			}

			uaddr, err := net.ResolveUDPAddr("udp", saddr)
			if err != nil {
				log.Println(err)
				continue
			}

			//_, err = listener.WriteToUDP(packet, uaddr)
			//if err != nil {
			//	log.Println(err)
			//} else {
			//	log.Println(fmt.Sprintf("o: %s → %s", key, saddr))
			//}
		}
	}()
}

var connections = storage.NewStorage()

func main() {
	tun, err := network.CreateTun("10.8.0.2", "tun0", 1500)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to create TUN: %s", err))
	} else {
		log.Println(fmt.Sprintf("TUN device started; IP: %s; name: %s", tun.Ip, tun.Name))
	}

	listener, err := CreateUdpListener()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to start UDP listener: %s", err))
	} else {
		log.Println(fmt.Sprintf("UDP listener started; port: %s", "8933"))
	}

	ListenUdp(tun.Interface, listener)
	ListenTun(tun.Interface, listener)

	util.Serve()
}
