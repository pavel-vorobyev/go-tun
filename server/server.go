package server

import (
	"github.com/xitongsys/ethernet-go/header"
	"go-tun/core/network"
	"go-tun/core/transport"
	"go-tun/server/config"
	"go-tun/server/storage/address"
)

type Server struct {
	conf            config.Config
	tun             *network.Tun
	conn            *transport.UDPConn
	cAddrKeyFactory address.CAddrKeyFactory
	cAddrStore      address.CAddrStore
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

	return &Server{
		tun:             tun,
		conn:            conn,
		cAddrKeyFactory: options.cAddrKeyFactory,
		cAddrStore:      options.cAddrStore,
	}, nil
}

func (s *Server) Start() {
	s.listenConn()
	s.listenTun()
}

func (s *Server) listenConn() {
	go func() {
		for {
			data, err := s.conn.Receive()
			if err != nil {
				//log.Println(fmt.Sprintf("SERVER: failed to read from con: %s", err))
				continue
			}

			ptc, src, dst, err := header.GetBase(data.Data)
			if err != nil {
				//log.Println(fmt.Sprintf("SERVER: failed to parse packet from con: %s", err))
				continue
			}

			//log.Println(fmt.Sprintf("in: %s %s %s %s", ptc, src, dst, data.CAddr))

			s.storeCAddr(ptc, src, dst, data.CAddr)
			_ = s.tun.Send(data.Data)
		}
	}()
}

func (s *Server) listenTun() {
	go func() {
		for {
			data, err := s.tun.Receive()
			if err != nil {
				//log.Println(fmt.Sprintf("SERVER: failed to read from tun: %s", err))
				continue
			}

			ptc, src, dst, err := header.GetBase(data)
			if err != nil {
				//log.Println(fmt.Sprintf("SERVER: failed to parse packet from tun: %s", err))
				continue
			}

			cAddr := s.getCAddr(ptc, src, dst)
			_ = s.conn.Send(&transport.Data{
				Data:  data,
				CAddr: cAddr,
			})

			//log.Println(fmt.Sprintf("out: %s %s %s %s", ptc, src, dst, cAddr))
		}
	}()
}

func (s *Server) storeCAddr(ptc string, src string, dst string, cAddr string) {
	key := s.cAddrKeyFactory.Get(ptc, src, dst)
	s.cAddrStore.Set(key, cAddr)
}

func (s *Server) getCAddr(ptc string, src string, dst string) string {
	key := s.cAddrKeyFactory.Get(ptc, dst, src)
	return s.cAddrStore.Get(key)
}
