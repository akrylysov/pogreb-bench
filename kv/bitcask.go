package kv

import (
	"git.mills.io/prologic/bitcask"
)

func newBitcask(path string) (Store, error) {
	db, err := bitcask.Open(path)
	if err != nil {
		return nil, err
	}
	return db, nil
}
