package kv

import (
	"github.com/akrylysov/pogreb"
	"github.com/akrylysov/pogreb/fs"
)

func newPogreb(path string) (Store, error) {
	db, err := pogreb.Open(path, &pogreb.Options{FileSystem: fs.OS})
	if err != nil {
		return nil, err
	}
	return db, nil
}
