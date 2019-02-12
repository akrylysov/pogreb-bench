package main

import (
	"github.com/etcd-io/bbolt"
)

var boltBucketName = []byte("benchmark")

type bboltEngine struct {
	db   *bbolt.DB
	path string
}

func newBbolt(path string) (kvEngine, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.NoSync = true
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket(boltBucketName)
		return err
	})
	return &bboltEngine{db: db, path: path}, err
}

func (db *bboltEngine) Put(key []byte, value []byte) error {
	return db.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(boltBucketName)
		return b.Put(key, value)
	})
}

func (db *bboltEngine) Get(key []byte) ([]byte, error) {
	var val []byte
	err := db.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(boltBucketName)
		val = b.Get(key)
		return nil
	})
	return val, err
}

func (db *bboltEngine) Close() error {
	return db.db.Close()
}

func (db *bboltEngine) FileSize() (int64, error) {
	return dirSize(db.path)
}
