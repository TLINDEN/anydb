/*
Copyright © 2024 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const MaxValueWidth int = 60

type DB struct {
	Debug  bool
	Dbfile string
	Bucket string
	DB     *bolt.DB
}

type BucketInfo struct {
	Name     string
	Keys     int
	Size     int
	Sequence uint64
	Stats    bolt.BucketStats
}

type DbInfo struct {
	Buckets []BucketInfo
	Path    string
}

type DbEntries []*DbEntry

type DbTag struct {
	Keys []string `json:"key"`
}

const BucketData string = "data"

func GetDbFile(file string) string {
	if file != "" {
		return file
	}

	file = os.Getenv("ANYDB_DB")
	if file != "" {
		return file
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, ".config", "anydb", "default.db")
}

func New(file string, bucket string, debug bool) (*DB, error) {
	return &DB{Debug: debug, Dbfile: file, Bucket: bucket}, nil
}

func (db *DB) Open() error {
	slog.Debug("opening DB", "dbfile", db.Dbfile)

	if _, err := os.Stat(filepath.Dir(db.Dbfile)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(db.Dbfile), 0700); err != nil {
			return err
		}
	}

	b, err := bolt.Open(db.Dbfile, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open DB %s: %w", db.Dbfile, err)
	}

	db.DB = b
	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) List(attr *DbAttr, fulltext bool) (DbEntries, error) {
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()

	var entries DbEntries
	var filter *regexp.Regexp

	if len(attr.Args) > 0 {
		// via cli
		filter = regexp.MustCompile(attr.Args[0])
	}

	if len(attr.Key) > 0 {
		// via api
		filter = regexp.MustCompile(attr.Key)
	}

	err := db.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(db.Bucket))
		if root == nil {
			return nil
		}

		slog.Debug("opened root bucket", "root", root)

		bucket := root.Bucket([]byte("meta"))
		if bucket == nil {
			return nil
		}

		slog.Debug("opened buckets", "root", root, "data", bucket)

		databucket := root.Bucket([]byte("data"))
		if databucket == nil {
			return fmt.Errorf("failed to retrieve data sub bucket")
		}

		err := bucket.ForEach(func(key, pbentry []byte) error {
			var entry DbEntry
			if err := proto.Unmarshal(pbentry, &entry); err != nil {
				return fmt.Errorf("failed to unmarshal from protobuf: %w", err)
			}

			if fulltext {
				// avoid crash due to access fault
				value := databucket.Get([]byte(entry.Key)) // empty is ok
				vc := make([]byte, len(value))
				copy(vc, value)
				entry.Value = string(vc)
			}

			var include bool

			switch {
			case filter != nil:
				if filter.MatchString(entry.Key) ||
					filter.MatchString(strings.Join(entry.Tags, " ")) {
					include = true
				}

				if !entry.Binary && !include && fulltext {
					if filter.MatchString(string(entry.Value)) {
						include = true
					}
				}
			case len(attr.Tags) > 0:
				for _, search := range attr.Tags {
					for _, tag := range entry.Tags {
						if tag == search {
							include = true
							break
						}
					}

					if include {
						break
					}
				}
			default:
				include = true
			}

			if include {
				entries = append(entries, &entry)
			}

			return nil
		})

		return err
	})

	return entries, err
}

func (db *DB) Set(attr *DbAttr) error {
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	entry := DbEntry{
		Key:       attr.Key,
		Binary:    attr.Binary,
		Tags:      attr.Tags,
		Encrypted: attr.Encrypted,
		Created:   timestamppb.Now(),
		Size:      uint64(len(attr.Val)),
		Preview:   attr.Preview,
	}

	// check if the  entry already exists and if yes,  check if it has
	// any  tags. if so,  we initialize  our update struct  with these
	// tags unless it has new tags configured.
	// FIXME: use Get()
	err := db.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(db.Bucket))
		if root == nil {
			return nil
		}

		bucket := root.Bucket([]byte("meta"))
		if bucket == nil {
			return nil
		}

		slog.Debug("opened buckets", "root", root, "data", bucket)

		pbentry := bucket.Get([]byte(entry.Key))
		if pbentry == nil {
			return nil
		}

		var oldentry DbEntry
		if err := proto.Unmarshal(pbentry, &oldentry); err != nil {
			return fmt.Errorf("failed to unmarshal from protobuf: %w", err)
		}

		if len(oldentry.Tags) > 0 && len(entry.Tags) == 0 {
			// initialize update entry with tags from old entry
			entry.Tags = oldentry.Tags
		}

		return nil
	})

	if err != nil {
		return err
	}

	// marshall our data
	pbentry, err := proto.Marshal(&entry)
	if err != nil {
		return fmt.Errorf("failed to marshall protobuf: %w", err)
	}

	err = db.DB.Update(func(tx *bolt.Tx) error {
		// create root bucket
		root, err := tx.CreateBucketIfNotExists([]byte(db.Bucket))
		if err != nil {
			return fmt.Errorf("failed to create DB bucket: %w", err)
		}

		// create meta bucket
		bucket, err := root.CreateBucketIfNotExists([]byte("meta"))
		if err != nil {
			return fmt.Errorf("failed to create DB meta sub bucket: %w", err)
		}

		slog.Debug("opened/created buckets", "root", root, "data", bucket)

		// write meta data
		err = bucket.Put([]byte(entry.Key), []byte(pbentry))
		if err != nil {
			return fmt.Errorf("failed to insert data: %w", err)
		}

		// create data bucket
		databucket, err := root.CreateBucketIfNotExists([]byte("data"))
		if err != nil {
			return fmt.Errorf("failed to create DB data sub bucket: %w", err)
		}

		// write value
		err = databucket.Put([]byte(entry.Key), attr.Val)
		if err != nil {
			return fmt.Errorf("failed to insert data: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Get(attr *DbAttr) (*DbEntry, error) {
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()

	entry := DbEntry{}

	err := db.DB.View(func(tx *bolt.Tx) error {
		// root bucket
		root := tx.Bucket([]byte(db.Bucket))
		if root == nil {
			return nil
		}

		// get meta sub bucket
		bucket := root.Bucket([]byte("meta"))
		if bucket == nil {
			return nil
		}

		slog.Debug("opened buckets", "root", root, "data", bucket)

		// retrieve meta data
		pbentry := bucket.Get([]byte(attr.Key))
		if pbentry == nil {
			return fmt.Errorf("no such key: %s", attr.Key)
		}

		// put into struct
		if err := proto.Unmarshal(pbentry, &entry); err != nil {
			return fmt.Errorf("failed to unmarshal from protobuf: %w", err)
		}

		// get data sub bucket
		databucket := root.Bucket([]byte("data"))
		if databucket == nil {
			return fmt.Errorf("failed to retrieve data sub bucket")
		}

		// retrieve actual data value
		value := databucket.Get([]byte(attr.Key))
		if len(value) == 0 {
			return fmt.Errorf("no such key: %s", attr.Key)
		}

		// we  need to make a  copy of it, otherwise  we'll get an
		// "unexpected fault address" error
		vc := make([]byte, len(value))
		copy(vc, value)

		entry.Value = string(vc)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read from DB: %w", err)
	}

	return &entry, nil
}

func (db *DB) Del(attr *DbAttr) error {
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	err := db.DB.Update(func(tx *bolt.Tx) error {
		// root bucket
		root := tx.Bucket([]byte(db.Bucket))
		if root == nil {
			return nil
		}

		// get data sub bucket
		bucket := root.Bucket([]byte("meta"))
		if bucket == nil {
			return nil
		}

		slog.Debug("opened buckets", "data", bucket)

		return bucket.Delete([]byte(attr.Key))
	})

	return err
}

func (db *DB) Import(attr *DbAttr) (string, error) {
	// open json file into attr.Val
	if err := attr.GetFileValue(); err != nil {
		return "", err
	}

	if len(attr.Val) == 0 {
		return "", errors.New("empty json file")
	}

	var entries DbEntries
	now := time.Now()
	newfile := db.Dbfile + now.Format("-02.01.2006T03:04.05")

	if err := json.Unmarshal([]byte(attr.Val), &entries); err != nil {
		return "", cleanError(newfile, fmt.Errorf("failed to unmarshal json: %w", err))
	}

	if fileExists(db.Dbfile) {
		// backup the old file
		err := os.Rename(db.Dbfile, newfile)
		if err != nil {
			return "", fmt.Errorf("failed to rename file %s to %s: %w", db.Dbfile, newfile, err)
		}

	}

	// should now be a new db file
	if err := db.Open(); err != nil {
		return "", cleanError(newfile, err)
	}
	defer db.Close()

	err := db.DB.Update(func(tx *bolt.Tx) error {
		// create root bucket
		root, err := tx.CreateBucketIfNotExists([]byte(db.Bucket))
		if err != nil {
			return fmt.Errorf("failed to create DB bucket: %w", err)
		}

		// create meta bucket
		bucket, err := root.CreateBucketIfNotExists([]byte("meta"))
		if err != nil {
			return fmt.Errorf("failed to create DB meta sub bucket: %w", err)
		}

		slog.Debug("opened buckets", "root", root, "data", bucket)

		for _, entry := range entries {
			pbentry, err := proto.Marshal(entry)
			if err != nil {
				return fmt.Errorf("failed to marshall protobuf: %w", err)
			}

			// write meta data
			err = bucket.Put([]byte(entry.Key), []byte(pbentry))
			if err != nil {
				return fmt.Errorf("failed to insert data into DB: %w", err)
			}

			// create data bucket
			databucket, err := root.CreateBucketIfNotExists([]byte("data"))
			if err != nil {
				return fmt.Errorf("failed to create DB data sub bucket: %w", err)
			}

			// write value
			err = databucket.Put([]byte(entry.Key), []byte(entry.Value))
			if err != nil {
				return fmt.Errorf("failed to insert data: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return "", cleanError(newfile, err)
	}

	return fmt.Sprintf("backed up database file to %s\nimported %d database entries\n",
		newfile, len(entries)), nil
}

func (db *DB) Info() (*DbInfo, error) {
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()

	info := &DbInfo{Path: db.Dbfile}

	err := db.DB.View(func(tx *bolt.Tx) error {
		err := tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			stats := bucket.Stats()

			binfo := BucketInfo{
				Name:     string(name),
				Sequence: bucket.Sequence(),
				Keys:     stats.KeyN,
				Stats:    bucket.Stats(),
			}

			err := bucket.ForEach(func(key, entry []byte) error {
				binfo.Size += len(entry) + len(key)

				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to read keys: %w", err)
			}

			info.Buckets = append(info.Buckets, binfo)
			return nil

		})

		if err != nil {
			return fmt.Errorf("failed to read from DB: %w", err)
		}

		return nil

	})

	return info, err
}

func (db *DB) Getall(attr *DbAttr) (DbEntries, error) {
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()

	var entries DbEntries

	err := db.DB.View(func(tx *bolt.Tx) error {
		// root bucket
		root := tx.Bucket([]byte(db.Bucket))
		if root == nil {
			return nil
		}

		// get meta sub bucket
		bucket := root.Bucket([]byte("meta"))
		if bucket == nil {
			return nil
		}

		// get data sub bucket
		databucket := root.Bucket([]byte("data"))
		if databucket == nil {
			return fmt.Errorf("failed to retrieve data sub bucket")
		}

		slog.Debug("opened buckets", "root", root, "data", bucket)

		// iterate over all db entries in meta sub bucket
		err := bucket.ForEach(func(key, pbentry []byte) error {
			var entry DbEntry
			if err := proto.Unmarshal(pbentry, &entry); err != nil {
				return fmt.Errorf("failed to unmarshal from protobuf: %w", err)
			}

			// retrieve the value from the data sub bucket
			value := databucket.Get([]byte(entry.Key))

			// we  need to make a  copy of it, otherwise  we'll get an
			// "unexpected fault address" error
			vc := make([]byte, len(value))
			copy(vc, value)

			entry.Value = string(vc)
			entries = append(entries, &entry)

			return nil
		})

		return err
	})
	return entries, err
}
