package main

import (
	"github.com/recoilme/slowpoke"
)

func newSlowpoke(path string) (kvEngine, error) {
	_, err := slowpoke.Open(path)
	return &slowpokeEngine{Path: path}, err
}

type slowpokeEngine struct {
	Path string
}

func (db *slowpokeEngine) Put(key []byte, value []byte) error {
	return slowpoke.Set(db.Path, key, value)
}

func (db *slowpokeEngine) Get(key []byte) ([]byte, error) {
	return slowpoke.Get(db.Path, key)
}

func (db *slowpokeEngine) Close() error {
	return slowpoke.Close(db.Path)
}

func (db *slowpokeEngine) FileSize() (int64, error) {
	return 0, nil
}
