package cfg

import "github.com/tlinden/anydb/app"

var Version string = "v0.0.1"

type Config struct {
	Debug     bool
	Dbfile    string
	Mode      string // wide, table, yaml, json
	NoHeaders bool
	Encrypt   bool
	DB        *app.DB
	File      string
	Tags      []string
}
