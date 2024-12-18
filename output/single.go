package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
	"golang.org/x/term"
)

func Print(writer io.Writer, conf *cfg.Config, entry *app.DbEntry) error {
	if conf.Mode != "" {
		// consider this to be a file
		fd, err := os.OpenFile(conf.Mode, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
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
			return err
		}

		return nil
	}

	isatty := term.IsTerminal(int(os.Stdout.Fd()))
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

	return nil
}
