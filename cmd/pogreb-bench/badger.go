package main

import (
	"github.com/dgraph-io/badger"
)

func newBadger(path string) (kvEngine, error) {
	opts := badger.DefaultOptions
	opts.SyncWrites = false
	opts.Dir = path
	opts.ValueDir = path
	db, err := badger.Open(opts)
	return &badgerEngine{db: db, path: path}, err
}

type badgerEngine struct {
	path string
	db   *badger.DB
}

func (db *badgerEngine) Put(key []byte, value []byte) error {
	return db.db.Update(func(tx *badger.Txn) error {
		return tx.Set(key, value)
	})
}

func (db *badgerEngine) Get(key []byte) ([]byte, error) {
	var val []byte
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		v, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		val = v
		return nil
	})
	return val, err
}

func (db *badgerEngine) Close() error {
	return db.db.Close()
}

func (db *badgerEngine) FileSize() (int64, error) {
	return dirSize(db.path)
}
