package kv

import (
	"go.etcd.io/bbolt"
)

var boltBucketName = []byte("benchmark")

type bboltStore struct {
	db *bbolt.DB
}

func newBbolt(path string) (Store, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.NoSync = true
	_ = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket(boltBucketName)
		return err
	})
	return &bboltStore{db: db}, err
}

func (s *bboltStore) Put(key []byte, value []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(boltBucketName)
		return b.Put(key, value)
	})
}

func (s *bboltStore) Get(key []byte) ([]byte, error) {
	var val []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(boltBucketName)
		val = b.Get(key)
		return nil
	})
	return val, err
}

func (s *bboltStore) Delete(key []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(boltBucketName)
		return b.Delete(key)
	})
}

func (s *bboltStore) Close() error {
	return s.db.Close()
}
