package server

import (
	"fmt"
	"github.com/xitongsys/ethernet-go/header"
	"go-tun/core/network"
	"go-tun/core/transport"
	"go-tun/server/config"
	"go-tun/server/packet"
	"go-tun/server/storage/address"
	"go-tun/util"
	"log"
	"os"
	"runtime"
)

type Server struct {
	conf                config.Config
	tun                 *network.Tun
	conn                *transport.UDPConn
	cAddrKeyFactory     address.CAddrKeyFactory
	cAddrStore          address.CAddrStore
	rxModifiers         []packet.Modifier
	txModifiers         []packet.Modifier
	rxCallbacks         []packet.Callback
	txCallbacks         []packet.Callback
	rxCallbackCallQueue *util.Queue[packet.CallbackCall]
	txCallbackCallQueue *util.Queue[packet.CallbackCall]
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
		tun:                 tun,
		conn:                conn,
		cAddrKeyFactory:     options.cAddrKeyFactory,
		cAddrStore:          options.cAddrStore,
		rxModifiers:         options.rxModifiers,
		txModifiers:         options.txModifiers,
		rxCallbacks:         options.rxCallbacks,
		txCallbacks:         options.txCallbacks,
		rxCallbackCallQueue: &util.Queue[packet.CallbackCall]{},
		txCallbackCallQueue: &util.Queue[packet.CallbackCall]{},
	}, nil
}

func (s *Server) Start() {
	s.listenConn()
	s.listenTun()
	s.callCallbacks()
}

func (s *Server) listenConn() {
	go func() {
		for {
			n, data, err := s.conn.Receive()
			if err != nil {
				continue
			}

			ptc, src, dst, err := header.GetBase(data.Data)
			if err != nil {
				continue
			}

			s.storeCAddr(ptc, src, dst, data.CAddr)
			_ = s.tun.Send(data.Data)

			s.addRxCallbackCall(ptc, src, dst, n)
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

			ptc, src, dst, err := header.GetBase(data)
			if err != nil {
				continue
			}

			cAddr := s.getCAddr(ptc, src, dst)
			_ = s.conn.Send(&transport.Data{
				Data:  data,
				CAddr: cAddr,
			})

			s.addTxCallbackCall(ptc, src, dst, n)
		}
	}()
}

func (s *Server) callCallbacks() {
	go func() {
		for {
			call := s.rxCallbackCallQueue.Pop()
			if call != nil {
				for _, callback := range s.rxCallbacks {
					callback.Call(call)
				}
			}
			call = s.txCallbackCallQueue.Pop()
			if call != nil {
				for _, callback := range s.txCallbacks {
					callback.Call(call)
				}
			}
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

func (s *Server) addRxCallbackCall(ptc string, src string, dst string, n int) {
	if len(s.rxCallbacks) != 0 {
		s.rxCallbackCallQueue.Put(
			&packet.CallbackCall{
				Ptc: ptc,
				Src: src,
				Dst: dst,
				N:   n,
			},
		)
	}
}

func (s *Server) addTxCallbackCall(ptc string, src string, dst string, n int) {
	if len(s.txCallbacks) != 0 {
		s.txCallbackCallQueue.Put(
			&packet.CallbackCall{
				Ptc: ptc,
				Src: src,
				Dst: dst,
				N:   n,
			},
		)
	}
}
