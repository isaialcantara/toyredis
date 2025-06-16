package storage

import (
	"fmt"
)

type MapStore struct {
	data map[string]any
}

var _ StringStorage = (*MapStore)(nil)

func NewMapStore() *MapStore {
	return &MapStore{data: make(map[string]any)}
}

func (s *MapStore) Get(key string) ([]byte, error) {
	val, exists := s.data[key]
	if !exists {
		return nil, fmt.Errorf("not found")
	}

	slice, isByteSlice := val.([]byte)
	if !isByteSlice {
		return nil, fmt.Errorf("wrong type")
	}

	return slice, nil
}

func (s *MapStore) Set(key string, val []byte) error {
	s.data[key] = val
	return nil
}

func (s *MapStore) Del(key string) (bool, error) {
	if v, _ := s.Get(key); v == nil {
		return false, nil
	}

	delete(s.data, key)
	return true, nil
}
