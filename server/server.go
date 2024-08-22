package server

import (
	"fmt"
	"github.com/xitongsys/ethernet-go/header"
	"go-tun/core/network"
	"go-tun/core/transport"
	"go-tun/server/config"
	"go-tun/server/storage/address"
	"log"
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
	s.tun.Start()
	s.conn.Start()
	s.listen()
}

func (s *Server) listen() {
	s.listenConn()
	s.listenTun()
}

func (s *Server) listenConn() {
	go func() {
		for {
			data := s.conn.Receive()

			ptc, src, dst, err := header.GetBase(data.Data)
			if err != nil {
				continue
			}

			log.Println(fmt.Sprintf("in: %s %s %s %s:%d", ptc, src, dst, data.CAddr.IP.String(), data.CAddr.Port))

			s.storeCAddr(ptc, src, dst, data.CAddr)
			s.tun.Send(data.Data)
		}
	}()
}

func (s *Server) listenTun() {
	go func() {
		for {
			data := s.tun.Receive()

			ptc, src, dst, err := header.GetBase(data)
			if err != nil {
				continue
			}

			cAddr := s.getCAddr(ptc, src, dst)
			log.Println(fmt.Sprintf("out: %s %s %s %s:%d", ptc, src, dst, cAddr.IP.String(), cAddr.Port))
			
			s.conn.Send(&transport.Data{
				Data:  data,
				CAddr: cAddr,
			})
		}
	}()
}

func (s *Server) storeCAddr(ptc string, src string, dst string, cAddr *transport.CAddr) {
	key := s.cAddrKeyFactory.Get(ptc, src, dst)
	s.cAddrStore.Set(key, cAddr)
}

func (s *Server) getCAddr(ptc string, src string, dst string) *transport.CAddr {
	key := s.cAddrKeyFactory.Get(ptc, dst, src)
	return s.cAddrStore.Get(key)
}