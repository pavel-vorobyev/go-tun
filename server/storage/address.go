package storage

import (
	"fmt"
	"log"
	"sync"
)

type CAddrStore struct {
	storage map[string]string
	rwMutex sync.RWMutex
}

func NewCAddrStore() *CAddrStore {
	return &CAddrStore{
		storage: make(map[string]string),
	}
}

func (s *CAddrStore) Get(key string) string {
	s.rwMutex.RLock()
	value := s.storage[key]
	s.rwMutex.RUnlock()
	return value
}

func (s *CAddrStore) Set(key string, value string) {
	s.rwMutex.Lock()
	s.storage[key] = value
	s.rwMutex.Unlock()
}

func (s *CAddrStore) Summary() {
	count := 0

	for k, v := range s.storage {
		log.Println(fmt.Sprintf("%s : %s", k, v))
		count++
	}
	log.Println(fmt.Sprintf("Total: %d", count))
}
