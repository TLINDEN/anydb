package app

import (
	"os"
	"path/filepath"
	"time"

	"github.com/asdine/storm/v3"
)

type DB struct {
	Debug bool
	DB    *storm.DB
}

type DbAttr struct {
	Key  string
	Args []string
	Tags []string
	File string
}

type DbEntry struct {
	ID        int       `storm:"id,increment"`
	Key       string    `storm:"unique"`
	Value     string    `storm:"index"` // FIXME: turn info []byte or add blob?
	Tags      []string  `storm:"index"`
	CreatedAt time.Time `storm:"index"`
}

func New(file string, debug bool) (*DB, error) {
	if _, err := os.Stat(filepath.Dir(file)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(file), 0700)
	}

	db, err := storm.Open(file)
	if err != nil {
		return nil, err
	}
	// FIXME: defer db.Close() here leads to: Error: database not open

	return &DB{Debug: debug, DB: db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Set(attr *DbAttr) error {
	entry := DbEntry{Key: attr.Key, Tags: attr.Tags}

	if len(attr.Args) > 0 {
		entry.Value = attr.Args[0]
	}

	// FIXME: check attr.File or STDIN

	return db.DB.Save(&entry)
}
