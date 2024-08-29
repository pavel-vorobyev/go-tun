package server

import (
	"fmt"
	"go-tun/core/network"
	"go-tun/core/transport"
	"go-tun/server/config"
	"go-tun/server/storage"
	"go-tun/util"
	"os"
	"runtime"
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
	//dMod, err := s.callModifiers(data.GetData(), s.rxModifiers)
	//if err != nil {
	//	return
	//}

	ptc, src, dst := util.GetPacketBaseInfo(data.GetData())
	s.cAddrStore.Set(src, data.GetCAddr())

	if ptc != 6 && ptc != 17 {
		return
	}

	err := s.tun.Send(data.GetData())
	if err != nil {
		return
	}

	s.callCallbacks(ptc, src, dst, n, s.rxCallbacks)
}

func (s *Server) handleTunPacket(n int, data []byte) {
	//dMod, err := s.callModifiers(data, s.txModifiers)
	//if err != nil {
	//	return
	//}

	ptc, src, dst := util.GetPacketBaseInfo(data)
	cAddr := s.cAddrStore.Get(dst)

	if ptc != 6 && ptc != 17 {
		return
	}

	err := s.conn.Send(transport.NewData(data, cAddr))
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
	if runtime.NumCPU() < 2 {
		util.PrintWarning(
			"\nThis instance running on 1 CPU core.\n" +
				"A minimum of 2 CPU cores is recommended for best performance.\n" +
				"Niddle will run on 1 CPU core anyway, but significantly slower.",
		)
	}
	util.PrintHelloNidde()
	util.PrintMessage("Niddle successfully started")
	util.PrintMessage(fmt.Sprintf("PID: %d", os.Getpid()))
}
