package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

type DB struct {
	Debug  bool
	Dbfile string
	DB     *bolt.DB
}

type DbEntry struct {
	Id      string    `json:"id"`
	Key     string    `json:"key"`
	Value   string    `json:"value"`
	Bin     []byte    `json:"bin"`
	Tags    []string  `json:"tags"`
	Created time.Time `json:"created"`
}

type DbEntries []DbEntry

type DbTag struct {
	Keys []string `json:"key"`
}

const BucketData string = "data"
const BucketTags string = "tags"

func New(file string, debug bool) (*DB, error) {
	if _, err := os.Stat(filepath.Dir(file)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(file), 0700)
	}

	return &DB{Debug: debug, Dbfile: file}, nil
}

func (db *DB) Open() error {
	b, err := bolt.Open(db.Dbfile, 0600, nil)
	if err != nil {
		return err
	}

	db.DB = b
	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) List(attr *DbAttr) (DbEntries, error) {
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()

	var entries DbEntries
	var filter *regexp.Regexp

	if len(attr.Args) > 0 {
		filter = regexp.MustCompile(attr.Args[0])
	}

	err := db.DB.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(BucketData))
		if bucket == nil {
			return nil
		}

		err := bucket.ForEach(func(key, jsonentry []byte) error {
			var entry DbEntry
			if err := json.Unmarshal(jsonentry, &entry); err != nil {
				return fmt.Errorf("unable to unmarshal json: %s", err)
			}

			var include bool

			switch {
			case filter != nil:
				if filter.MatchString(entry.Value) ||
					filter.MatchString(entry.Key) ||
					filter.MatchString(strings.Join(entry.Tags, " ")) {
					include = true
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
				entries = append(entries, entry)
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

	if err := attr.ParseKV(); err != nil {
		return err
	}

	entry := DbEntry{
		Key:     attr.Key,
		Value:   attr.Val,
		Bin:     attr.Bin,
		Tags:    attr.Tags,
		Created: time.Now(),
	}

	err := db.DB.Update(func(tx *bolt.Tx) error {
		// insert data
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketData))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		jsonentry, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("json marshalling failure: %s", err)
		}

		err = bucket.Put([]byte(entry.Key), []byte(jsonentry))
		if err != nil {
			return fmt.Errorf("insert data: %s", err)
		}

		// insert tag, if any
		// FIXME: check removed tags
		if len(attr.Tags) > 0 {
			bucket, err := tx.CreateBucketIfNotExists([]byte(BucketTags))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			for _, tag := range entry.Tags {
				dbtag := &DbTag{}

				jsontag := bucket.Get([]byte(tag))
				if jsontag == nil {
					// the tag is empty so far, initialize it
					dbtag.Keys = []string{entry.Key}
				} else {
					if err := json.Unmarshal(jsontag, dbtag); err != nil {
						return fmt.Errorf("unable to unmarshal json: %s", err)
					}

					if !slices.Contains(dbtag.Keys, entry.Key) {
						// current key is not yet assigned to the tag, append it
						dbtag.Keys = append(dbtag.Keys, entry.Key)
					}
				}

				jsontag, err = json.Marshal(dbtag)
				if err != nil {
					return fmt.Errorf("json marshalling failure: %s", err)
				}

				err = bucket.Put([]byte(tag), []byte(jsontag))
				if err != nil {
					return fmt.Errorf("insert data: %s", err)
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
