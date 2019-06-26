package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type goleveldbEngine struct {
	db    *leveldb.DB
	path  string
	batch *leveldb.Batch
}

func newGolevelDB(path string) (kvEngine, error) {
	opts := opt.Options{Compression: opt.NoCompression}
	db, err := leveldb.OpenFile(path, &opts)
	if err != nil {
		return nil, err
	}
	err = db.CompactRange(util.Range{})
	batch := new(leveldb.Batch)
	return &goleveldbEngine{db: db, path: path, batch: batch}, err
}

func (db *goleveldbEngine) Put(key []byte, value []byte) error {
	db.batch.Put(key, value)
	return nil
}

func (db *goleveldbEngine) Get(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

func (db *goleveldbEngine) Close() error {
	if err := db.db.Write(db.batch, nil); err != nil {
		return err
	}
	return db.db.Close()
}

func (db *goleveldbEngine) FileSize() (int64, error) {
	return dirSize(db.path)
}
