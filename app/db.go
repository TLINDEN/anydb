package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

const MaxValueWidth int = 60

type DB struct {
	Debug  bool
	Dbfile string
	DB     *bolt.DB
}

type DbEntry struct {
	Id        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Encrypted bool      `json:"encrypted"`
	Bin       []byte    `json:"bin"`
	Tags      []string  `json:"tags"`
	Created   time.Time `json:"created"`
	Size      int
}

// Post process an entry for list output.
// Do NOT call it during write processing!
func (entry *DbEntry) Normalize() {
	entry.Size = len(entry.Value)

	if entry.Encrypted {
		entry.Value = "<encrypted-content>"
	}

	if len(entry.Bin) > 0 {
		entry.Value = "<binary-content>"
		entry.Size = len(entry.Bin)
	}

	if len(entry.Value) > MaxValueWidth {
		entry.Value = entry.Value[0:MaxValueWidth] + "..."
	}
}

type DbEntries []DbEntry

type DbTag struct {
	Keys []string `json:"key"`
}

const BucketData string = "data"

func New(file string, debug bool) (*DB, error) {
	return &DB{Debug: debug, Dbfile: file}, nil
}

func (db *DB) Open() error {
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
				return fmt.Errorf("failed to unmarshal from json: %w", err)
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

	entry := DbEntry{
		Key:       attr.Key,
		Value:     attr.Val,
		Bin:       attr.Bin,
		Tags:      attr.Tags,
		Encrypted: attr.Encrypted,
		Created:   time.Now(),
	}

	// check if the  entry already exists and if yes,  check if it has
	// any  tags. if so,  we initialize  our update struct  with these
	// tags unless it has new tags configured.
	err := db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketData))
		if bucket == nil {
			return nil
		}

		jsonentry := bucket.Get([]byte(entry.Key))
		if jsonentry == nil {
			return nil
		}

		var oldentry DbEntry
		if err := json.Unmarshal(jsonentry, &oldentry); err != nil {
			return fmt.Errorf("failed to unmarshal from json: %w", err)
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

	err = db.DB.Update(func(tx *bolt.Tx) error {
		// insert data
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketData))
		if err != nil {
			return fmt.Errorf("failed to create DB bucket: %w", err)
		}

		jsonentry, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to marshall json: %w", err)
		}

		err = bucket.Put([]byte(entry.Key), []byte(jsonentry))
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
		bucket := tx.Bucket([]byte(BucketData))
		if bucket == nil {
			return nil
		}

		jsonentry := bucket.Get([]byte(attr.Key))
		if jsonentry == nil {
			return nil
		}

		if err := json.Unmarshal(jsonentry, &entry); err != nil {
			return fmt.Errorf("failed to unmarshal from json: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read from DB: %w", err)
	}

	return &entry, nil
}

func (db *DB) Del(attr *DbAttr) error {
	// FIXME: check if it exists prior to just call bucket.Delete()?
	if err := db.Open(); err != nil {
		return err
	}
	defer db.Close()

	err := db.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketData))

		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(attr.Key))
	})

	return err
}

func (db *DB) Import(attr *DbAttr) error {
	// open json file into attr.Val
	if err := attr.GetFileValue(); err != nil {
		return err
	}

	if attr.Val == "" {
		return errors.New("empty json file")
	}

	var entries DbEntries
	now := time.Now()
	newfile := db.Dbfile + now.Format("-02.01.2006T03:04.05")

	if err := json.Unmarshal([]byte(attr.Val), &entries); err != nil {
		return cleanError(newfile, fmt.Errorf("failed to unmarshal json: %w", err))
	}

	if fileExists(db.Dbfile) {
		// backup the old file
		err := os.Rename(db.Dbfile, newfile)
		if err != nil {
			return fmt.Errorf("failed to rename file %s to %s: %w", db.Dbfile, newfile, err)
		}

	}

	// should now be a new db file
	if err := db.Open(); err != nil {
		return cleanError(newfile, err)
	}
	defer db.Close()

	err := db.DB.Update(func(tx *bolt.Tx) error {
		// insert data
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketData))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		for _, entry := range entries {
			jsonentry, err := json.Marshal(entry)
			if err != nil {
				return fmt.Errorf("failed to marshall json: %w", err)
			}

			err = bucket.Put([]byte(entry.Key), []byte(jsonentry))
			if err != nil {
				return fmt.Errorf("failed to insert data into DB: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return cleanError(newfile, err)
	}

	fmt.Printf("backed up database file to %s\n", newfile)
	fmt.Printf("imported %d database entries\n", len(entries))

	return nil
}

func cleanError(file string, err error) error {
	// remove given [backup] file and forward the given error
	os.Remove(file)
	return err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if err != nil {
		// return false on any error
		return false
	}

	return !info.IsDir()
}
