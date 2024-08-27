package server

import (
	"fmt"
	"github.com/xitongsys/ethernet-go/header"
	"go-tun/core/network"
	"go-tun/core/transport"
	"go-tun/server/config"
	"go-tun/server/packet"
	"go-tun/server/storage/address"
	"log"
	"os"
	"runtime"
)

//var cAddr = "91.202.27.121:60796"

type Server struct {
	conf            config.Config
	tun             *network.Tun
	conn            *transport.UDPConn
	cAddrKeyFactory address.CAddrKeyFactory
	cAddrStore      address.CAddrStore
	rxModifiers     []packet.Modifier
	txModifiers     []packet.Modifier
	rxCallbacks     []packet.Callback
	txCallbacks     []packet.Callback
}

func CreateServer(options *Options) (*Server, error) {
	conf, err := options.configProvider.GetConfig()
	if err != nil {
		return nil, err
	}

	tunConf := &network.Config{
		Name: conf.TunName,
		Ip:   conf.TunIp,
		Cidr: conf.TunCidr,
		Mtu:  conf.Mtu,
	}
	connConf := &transport.Config{
		Ip:   conf.RemoteIp,
		Port: conf.RemotePort,
		Mtu:  conf.Mtu,
	}

	tun, err := network.CreateTun(tunConf)
	if err != nil {
		return nil, err
	}

	conn, err := transport.CreateConn(connConf)
	if err != nil {
		return nil, err
	}

	log.Println("Server is ready")
	log.Println(fmt.Sprintf("CPUs num: %d", runtime.NumCPU()))
	log.Println(fmt.Sprintf("PID: %d", os.Getpid()))

	return &Server{
		tun:             tun,
		conn:            conn,
		cAddrKeyFactory: options.cAddrKeyFactory,
		cAddrStore:      options.cAddrStore,
		rxModifiers:     options.rxModifiers,
		txModifiers:     options.txModifiers,
		rxCallbacks:     options.rxCallbacks,
		txCallbacks:     options.txCallbacks,
	}, nil
}

func (s *Server) Start() {
	s.listenConn()
	s.listenTun()
}

func (s *Server) listenConn() {
	go func() {
		for {
			n, data, err := s.conn.Receive()
			if err != nil {
				continue
			}
			s.handleConnPacket(n, data)
		}
	}()
}

func (s *Server) listenTun() {
	go func() {
		for {
			n, data, err := s.tun.Receive()
			if err != nil {
				continue
			}
			s.handleTunPacket(n, data)
		}
	}()
}

func (s *Server) handleConnPacket(n int, data *transport.Data) {
	ptc, src, dst, err := header.GetBase(data.Data)
	if err != nil {
		return
	}

	//s.storeCAddr(ptc, src, dst, data.CAddr)
	err = s.tun.Send(data.Data)
	if err != nil {
		return
	}

	s.callCallbacks(ptc, src, dst, n, s.rxCallbacks)
}

// 91.202.27.121:60796
// 91.202.27.121:60796

func (s *Server) handleTunPacket(n int, data []byte) {
	ptc, src, dst, err := header.GetBase(data)
	if err != nil {
		return
	}

	cAddr := s.getCAddr(ptc, src, dst)
	err = s.conn.Send(&transport.Data{
		Data:  data,
		CAddr: cAddr,
	})
	if err != nil {
		return
	}

	s.callCallbacks(ptc, src, dst, n, s.txCallbacks)
}

func (s *Server) storeCAddr(ptc string, src string, dst string, cAddr string) {
	key := s.cAddrKeyFactory.Get(ptc, src, dst)
	s.cAddrStore.Set(key, cAddr)
}

func (s *Server) getCAddr(ptc string, src string, dst string) string {
	key := s.cAddrKeyFactory.Get(ptc, dst, src)
	return s.cAddrStore.Get(key)
}

func (s *Server) callCallbacks(ptc string, src string, dst string, n int, c []packet.Callback) {
	for _, callback := range c {
		callback.Call(ptc, src, dst, n)
	}
}
