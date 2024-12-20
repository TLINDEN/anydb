package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
)

func WriteFile(attr *app.DbAttr, conf *cfg.Config, entries app.DbEntries) error {
	jsonentries, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("failed to marshall json: %w", err)
	}

	if attr.File == "-" {
		fmt.Println(string(jsonentries))
	} else {
		fd, err := os.OpenFile(attr.File, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return fmt.Errorf("failed to open file %s for writing: %w", attr.File, err)
		}

		if _, err := fd.Write(jsonentries); err != nil {
			return fmt.Errorf("failed writing to file %s: %w", attr.File, err)
		}

		fmt.Printf("database contents exported to %s\n", attr.File)
	}

	return nil
}
