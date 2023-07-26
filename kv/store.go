package kv

import (
	"errors"
)

type Store interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Close() error
}

var stores = map[string]func(string) (Store, error){
	"pogreb":    newPogreb,
	"goleveldb": newGoleveldb,
	"bbolt":     newBbolt,
	"badger":    newBadger,
	"bitcask":   newBitcask,
}

func NewStore(name string, path string) (Store, error) {
	ctr, ok := stores[name]
	if !ok {
		return nil, errors.New("unknown kv store")
	}
	return ctr(path)
}
