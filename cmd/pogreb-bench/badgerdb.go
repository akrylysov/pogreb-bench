package main

import (
	"github.com/dgraph-io/badger/v2"
)

func newBadgerdb(path string) (kvEngine, error) {
	opts := badger.DefaultOptions
	opts.SyncWrites = false
	opts.Dir = path
	opts.ValueDir = path
	db, err := badger.Open(opts)
	wb := db.NewWriteBatch()
	return &badgerdbEngine{db: db, path: path, wb: wb}, err
}

type badgerdbEngine struct {
	path string
	db   *badger.DB
	wb   *badger.WriteBatch
}

func (db *badgerdbEngine) Put(key []byte, value []byte) error {
	return db.wb.Set(key, value)
}

func (db *badgerdbEngine) Get(key []byte) ([]byte, error) {
	var val []byte
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	return val, err
}

func (db *badgerdbEngine) Close() error {
	if err := db.wb.Flush(); err != nil {
		return err
	}
	return db.db.Close()
}

func (db *badgerdbEngine) FileSize() (int64, error) {
	return dirSize(db.path)
}
