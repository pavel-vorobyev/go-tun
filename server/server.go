package server

import (
	"fmt"
	"go-tun/core/network"
	"go-tun/core/transport"
	"go-tun/server/config"
	"go-tun/server/storage"
	"go-tun/util"
	"os"
)

type Server struct {
	conf        config.Config
	tun         *network.Tun
	conn        *transport.UDPConn
	cAddrStore  *storage.CAddrStore
	rxModifiers []PacketModifier
	txModifiers []PacketModifier
	rxCallbacks []PacketCallback
	txCallbacks []PacketCallback
}

func NewServer(options *Options) (*Server, error) {
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

	tun, err := network.NewTun(tunConf)
	if err != nil {
		return nil, err
	}
	conn, err := transport.NewConn(connConf)
	if err != nil {
		return nil, err
	}

	cAddStore := storage.NewCAddrStore()

	return &Server{
		tun:         tun,
		conn:        conn,
		cAddrStore:  cAddStore,
		rxModifiers: options.rxModifiers,
		txModifiers: options.txModifiers,
		rxCallbacks: options.rxCallbacks,
		txCallbacks: options.txCallbacks,
	}, nil
}

func (s *Server) Start() {
	s.printMessages()
	s.listenConn()
	s.listenTun()
}

func (s *Server) listenConn() {
	go func() {
		for {
			n, data, cAddr, err := s.conn.Receive()
			if err != nil {
				continue
			}
			s.handleConnPacket(n, data, cAddr)
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

func (s *Server) handleConnPacket(n int, data []byte, cAddr string) {
	dMod, err := s.callModifiers(data, s.rxModifiers)
	if err != nil {
		return
	}

	ptc, src, dst, err := util.GetPacketBaseInfo(dMod, n)
	if err != nil {
		return
	}

	s.cAddrStore.Set(src, cAddr)

	err = s.tun.Send(dMod)
	if err != nil {
		return
	}

	s.callCallbacks(ptc, src, dst, n, s.rxCallbacks)
}

func (s *Server) handleTunPacket(n int, data []byte) {
	dMod, err := s.callModifiers(data, s.txModifiers)
	if err != nil {
		return
	}

	ptc, src, dst, err := util.GetPacketBaseInfo(dMod, n)
	if err != nil {
		return
	}

	cAddr := s.cAddrStore.Get(dst)
	err = s.conn.Send(dMod, cAddr)
	if err != nil {
		return
	}

	s.callCallbacks(ptc, src, dst, n, s.txCallbacks)
}

func (s *Server) callModifiers(data []byte, m []PacketModifier) ([]byte, error) {
	dMod := data
	for _, modifier := range m {
		res, err := modifier.Process(data)
		if err != nil {
			return nil, err
		}
		dMod = res
	}
	return dMod, nil
}

func (s *Server) callCallbacks(ptc int, src string, dst string, n int, c []PacketCallback) {
	for _, callback := range c {
		callback.Call(ptc, src, dst, n)
	}
}

func (s *Server) Sum() {
	s.cAddrStore.Summary()
}

func (s *Server) printMessages() {
	util.PrintHelloKiesp()
	util.PrintMessage("Server successfully started")
	util.PrintMessage(fmt.Sprintf("PID: %d", os.Getpid()))
}
