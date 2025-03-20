package server

import (
	"sync"
	"time"
)

type KVStore struct {
	data   map[string]string
	expiry map[string]time.Time
	mu     sync.RWMutex
}

func NewStore() *KVStore {
	return &KVStore{
		data:   make(map[string]string),
		expiry: make(map[string]time.Time),
	}
}

func (s *KVStore) Set(key, value string, px int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	delete(s.expiry, key)

	if px > 0 {
		expTime := time.Now().Add(time.Duration(px) * time.Millisecond)
		s.expiry[key] = expTime

		go func(k string, expiry time.Time) {
			time.Sleep(time.Until(expiry))
			s.mu.Lock()
			if time.Now().After(s.expiry[k]) {
				delete(s.data, k)
				delete(s.expiry, k)
			}
			s.mu.Unlock()
		}(key, expTime)
	}
}

func (s *KVStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if expiry, exists := s.expiry[key]; exists && time.Now().After(expiry) {
		s.mu.RUnlock()
		s.mu.Lock()
		delete(s.data, key)
		delete(s.expiry, key)
		s.mu.Unlock()
		return "", false
	}

	val, exists := s.data[key]
	return val, exists
}

func (s *KVStore) GetAllKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}
