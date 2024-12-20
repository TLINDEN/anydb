package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
	"golang.org/x/term"
)

func Print(writer io.Writer, conf *cfg.Config, attr *app.DbAttr, entry *app.DbEntry) error {
	if attr.File != "" {
		return WriteFile(writer, conf, attr, entry)
	}

	isatty := term.IsTerminal(int(os.Stdout.Fd()))

	switch conf.Mode {
	case "simple", "":
		if len(entry.Bin) > 0 {
			if isatty {
				fmt.Println("binary data omitted")
			} else {
				os.Stdout.Write(entry.Bin)
			}
		} else {
			fmt.Print(entry.Value)

			if !strings.HasSuffix(entry.Value, "\n") {
				// always add a terminal newline
				fmt.Println()
			}
		}
	case "json":
		jsonentry, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to marshall json: %s", err)
		}

		fmt.Println(string(jsonentry))
	case "wide":
		return ListTable(writer, conf, app.DbEntries{*entry})
	case "template":
		return ListTemplate(writer, conf, app.DbEntries{*entry})
	}

	return nil
}

func WriteFile(writer io.Writer, conf *cfg.Config, attr *app.DbAttr, entry *app.DbEntry) error {
	fd, err := os.OpenFile(attr.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file %s for writing: %w", attr.File, err)
	}
	defer fd.Close()

	if len(entry.Bin) > 0 {
		// binary file content
		_, err = fd.Write(entry.Bin)
	} else {
		val := entry.Value
		if !strings.HasSuffix(val, "\n") {
			// always add a terminal newline
			val += "\n"
		}

		_, err = fd.Write([]byte(val))
	}

	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", attr.File, err)
	}

	return nil
}
