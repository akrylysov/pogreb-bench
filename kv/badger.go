package kv

import (
	"github.com/dgraph-io/badger/v3"
)

func newBadger(path string) (Store, error) {
	opts := badger.DefaultOptions(path)
	opts.SyncWrites = false
	opts.Dir = path
	opts.ValueDir = path
	db, err := badger.Open(opts)
	return &badgerStore{db: db}, err
}

type badgerStore struct {
	db *badger.DB
}

func (s *badgerStore) Put(key []byte, value []byte) error {
	return s.db.Update(func(tx *badger.Txn) error {
		return tx.Set(key, value)
	})
}

func (s *badgerStore) Get(key []byte) ([]byte, error) {
	var val []byte
	err := s.db.View(func(txn *badger.Txn) error {
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

func (s *badgerStore) Delete(key []byte) error {
	return s.db.Update(func(tx *badger.Txn) error {
		return tx.Delete(key)
	})
}

func (s *badgerStore) Close() error {
	return s.db.Close()
}
