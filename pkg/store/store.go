package store

import (
	"errors"

	"github.com/ameyw07/go-kv-store/pkg"
)

type Store struct {
	data map[string]*pkg.Entry
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]*pkg.Entry, 0),
	}
}

func (s *Store) Get(key string) (*pkg.Entry, error) {

	if _, ok := s.data[key]; !ok {
		return nil, errors.New("Key not found")
	}

	return s.data[key], nil
}

func (s *Store) Put(key string, entry *pkg.Entry) (*pkg.Entry, error) {

	if _, ok := s.data[key]; !ok {
		return nil, errors.New("Key not found")
	}

	return s.data[key], nil
}
