package address

import (
	"go-tun/core/transport"
	"sync"
)

type CAddrStore interface {
	Get(key string) *transport.CAddr
	Set(key string, address *transport.CAddr)
}

type DefaultCAddrStore struct {
	storage map[string]*transport.CAddr
	rwMutex sync.RWMutex
}

func NewDefaultCAddrStore() *DefaultCAddrStore {
	return &DefaultCAddrStore{
		storage: make(map[string]*transport.CAddr),
	}
}

func (s *DefaultCAddrStore) Get(key string) *transport.CAddr {
	s.rwMutex.RLock()
	value := s.storage[key]
	s.rwMutex.RUnlock()
	return value
}

func (s *DefaultCAddrStore) Set(key string, value *transport.CAddr) {
	s.rwMutex.Lock()
	s.storage[key] = value
	s.rwMutex.Unlock()
}
