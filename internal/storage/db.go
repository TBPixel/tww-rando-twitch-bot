package storage

import (
	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
	"github.com/dgraph-io/badger"
	"github.com/timshannon/badgerhold"
)

// DB
type DB struct {
	store *badgerhold.Store
}

// Open
func Open(conf config.DB) (*DB, error) {
	options := badger.DefaultOptions(conf.Path)
	if conf.EnableLogging == false {
		options.Logger = nil
	}

	store, err := badgerhold.Open(badgerhold.Options{
		Options:          options,
		Encoder:          badgerhold.DefaultEncode,
		Decoder:          badgerhold.DefaultDecode,
		SequenceBandwith: 100,
	})
	if err != nil {
		return nil, err
	}

	return &DB{
		store,
	}, nil
}

// Close
func (db *DB) Close() error {
	return db.store.Close()
}
