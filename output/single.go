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
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/dustin/go-humanize"
	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
	"golang.org/x/term"
	//"github.com/alecthomas/repr"
)

func Print(writer io.Writer, conf *cfg.Config, attr *app.DbAttr, entry *app.DbEntry) error {
	if attr.File != "" {
		return WriteFile(writer, conf, attr, entry)
	}

	isatty := term.IsTerminal(int(os.Stdout.Fd()))

	switch conf.Mode {
	case "simple", "":
		if entry.Binary {
			if isatty {
				fmt.Println("binary data omitted")
			} else {
				if _, err := os.Stdout.WriteString(entry.Value); err != nil {
					return err
				}
			}
		} else {
			fmt.Print(string(entry.Value))
			if entry.Value[entry.Size-1] != '\n' {
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
		return ListTable(writer, conf, app.DbEntries{entry})
	case "template":
		return ListTemplate(writer, conf, app.DbEntries{entry})
	}

	return nil
}

func WriteFile(writer io.Writer, conf *cfg.Config, attr *app.DbAttr, entry *app.DbEntry) error {
	var fileHandle *os.File
	var err error

	if attr.File == "-" {
		fileHandle = os.Stdout
	} else {
		fd, err := os.OpenFile(attr.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to open file %s for writing: %w", attr.File, err)
		}
		defer func() {
			if err := fd.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		fileHandle = fd
	}

	// actually write file content
	_, err = fileHandle.WriteString(entry.Value)

	if !entry.Binary {
		if entry.Value[entry.Size-1] != '\n' {
			// always add a terminal newline
			_, err = fileHandle.Write([]byte{'\n'})
		}
	}

	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", attr.File, err)
	}

	return nil
}

func Info(writer io.Writer, conf *cfg.Config, info *app.DbInfo) error {
	if _, err := fmt.Fprintf(writer, "Database: %s\n", info.Path); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	for _, bucket := range info.Buckets {
		if conf.NoHumanize {
			if _, err := fmt.Fprintf(
				writer,
				"%19s: %s\n%19s: %d\n%19s: %d\n%19s: %t\n",
				"Bucket", bucket.Name,
				"Size", bucket.Size,
				"Keys", bucket.Keys,
				"Encrypted", conf.Encrypt); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		} else {
			if _, err := fmt.Fprintf(
				writer,
				"%19s: %s\n%19s: %s\n%19s: %d\n",
				"Bucket", bucket.Name,
				"Size", humanize.Bytes(uint64(bucket.Size)),
				"Keys", bucket.Keys); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		}

		if conf.Debug {
			val := reflect.ValueOf(&bucket.Stats).Elem()
			for i := 0; i < val.NumField(); i++ {
				if _, err := fmt.Fprintf(writer, "%19s: %v\n", val.Type().Field(i).Name, val.Field(i)); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
			}
		}

		if _, err := fmt.Fprintln(writer); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	return nil
}
