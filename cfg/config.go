package cfg

import "github.com/tlinden/anydb/app"

var Version string = "v0.0.4"

type Config struct {
	Debug      bool
	Dbfile     string
	Dbbucket     string
	Template   string
	Mode       string // wide, table, yaml, json
	NoHeaders  bool
	NoHumanize bool
	Encrypt    bool
	DB         *app.DB
	File       string
	Tags       []string
	Listen     string
}
