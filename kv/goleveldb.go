package kv

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type goleveldbStore struct {
	db *leveldb.DB
}

func newGoleveldb(path string) (Store, error) {
	opts := opt.Options{Compression: opt.NoCompression}
	db, err := leveldb.OpenFile(path, &opts)
	if err != nil {
		return nil, err
	}
	return &goleveldbStore{db: db}, err
}

func (s *goleveldbStore) Put(key []byte, value []byte) error {
	return s.db.Put(key, value, nil)
}

func (s *goleveldbStore) Get(key []byte) ([]byte, error) {
	return s.db.Get(key, nil)
}

func (s *goleveldbStore) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

func (s *goleveldbStore) Close() error {
	return s.db.Close()
}
