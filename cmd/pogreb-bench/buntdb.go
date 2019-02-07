package main

import (
	"github.com/tidwall/buntdb"
)

type buntEngine struct {
	db *buntdb.DB
}

func newBunt(path string) (kvEngine, error) {
	db, err := buntdb.Open(path)

	if err != nil {
		return nil, err
	}

	return &buntEngine{db: db}, err
}

func (db *buntEngine) Put(key []byte, value []byte) error {
	return db.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(string(key), string(value), nil)
		return err
	})
}

func (db *buntEngine) Get(key []byte) ([]byte, error) {
	var val []byte
	err := db.db.View(func(tx *buntdb.Tx) error {
		valstr, err := tx.Get(string(key))
		if err != nil {
			return err
		}
		val = []byte(valstr)
		return nil
	})
	return val, err
}

func (db *buntEngine) Close() error {
	return db.db.Close()
}
func (db *buntEngine) FileSize() (int64, error) {
	return 0, nil
}
