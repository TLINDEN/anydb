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
package cfg

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/common"
)

var Version string = "v0.1.0"

type BucketConfig struct {
	Encrypt bool
}

type Config struct {
	Debug           bool
	Dbfile          string
	Dbbucket        string
	Template        string
	Mode            string // wide, table, yaml, json
	NoHeaders       bool
	NoHumanize      bool
	Encrypt         bool // one entry
	CaseInsensitive bool
	Fulltext        bool
	Listen          string
	Buckets         map[string]BucketConfig // config file only

	Tags []string // internal
	DB   *app.DB  // internal
	File string   // internal
}

func (conf *Config) GetConfig(files []string) error {
	for _, file := range files {
		if err := conf.ParseConfigFile(file); err != nil {
			return err
		}
	}

	return nil
}

func (conf *Config) ParseConfigFile(file string) error {
	if !common.FileExists(file) {
		return nil
	}

	fd, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %w", file, err)
	}

	data, err := io.ReadAll(fd)
	if err != nil {
		return fmt.Errorf("failed to read from config file: %w", err)
	}

	add := Config{}
	err = toml.Unmarshal(data, &add)
	if err != nil {
		return fmt.Errorf("failed to unmarshall toml: %w", err)
	}

	// merge new values into existing config
	switch {
	case add.Debug != conf.Debug:
		conf.Debug = add.Debug
	case add.Dbfile != "":
		conf.Dbfile = add.Dbfile
	case add.Dbbucket != "":
		conf.Dbbucket = add.Dbbucket
	case add.Template != "":
		conf.Template = add.Template
	case add.NoHeaders != conf.NoHeaders:
		conf.NoHeaders = add.NoHeaders
	case add.NoHumanize != conf.NoHumanize:
		conf.NoHumanize = add.NoHumanize
	case add.Encrypt != conf.Encrypt:
		conf.Encrypt = add.Encrypt
	case add.Listen != "":
		conf.Listen = add.Listen
	}

	// only supported in config files
	conf.Buckets = add.Buckets

	// determine bucket encryption mode
	for name, bucket := range conf.Buckets {
		if name == conf.Dbbucket {
			conf.Encrypt = bucket.Encrypt
		}
	}

	return nil
}
