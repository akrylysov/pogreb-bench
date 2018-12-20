package main

import (
	"github.com/recoilme/pudge"
)

func newPudge(path string) (kvEngine, error) {
	db, err := pudge.Open(path, nil)
	return &pudgeEngine{Db: db, Path: path}, err
}

type pudgeEngine struct {
	Db   *pudge.Db
	Path string
}

func (en *pudgeEngine) Put(key []byte, value []byte) error {
	return en.Db.Set(key, value)
}

func (en *pudgeEngine) Get(key []byte) ([]byte, error) {
	var b []byte
	err := en.Db.Get(key, &b)
	return b, err
}

func (en *pudgeEngine) Close() error {
	return en.Db.Close()
}

func (en *pudgeEngine) FileSize() (int64, error) {
	return en.Db.FileSize()
}
