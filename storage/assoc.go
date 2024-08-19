package storage

import "sync"

type Storage struct {
	storage map[string]string
	mutex   sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]string),
	}
}

func (st *Storage) Get(key string) (string, bool) {
	st.mutex.RLock()
	value, exists := st.storage[key]
	st.mutex.RUnlock()
	return value, exists
}

func (st *Storage) Put(key string, value string) {
	st.mutex.Lock()
	st.storage[key] = value
	st.mutex.Unlock()
}
