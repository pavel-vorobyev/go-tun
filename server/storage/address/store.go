package address

import (
	"sync"
)

type CAddrStore interface {
	Get(key string) string
	Set(key string, value string)
}

type DefaultCAddrStore struct {
	storage map[string]string
	rwMutex sync.RWMutex
}

func NewDefaultCAddrStore() *DefaultCAddrStore {
	return &DefaultCAddrStore{
		storage: make(map[string]string),
	}
}

func (s *DefaultCAddrStore) Get(key string) string {
	s.rwMutex.RLock()
	value := s.storage[key]
	s.rwMutex.RUnlock()
	return value
}

func (s *DefaultCAddrStore) Set(key string, value string) {
	s.rwMutex.Lock()
	s.storage[key] = value
	s.rwMutex.Unlock()
}
