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
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type DbAttr struct {
	Key       string
	Preview   string
	Val       []byte
	Args      []string
	Tags      []string
	File      string
	Encrypted bool
	Binary    bool
}

func (attr *DbAttr) ParseKV() error {
	switch len(attr.Args) {
	case 1:
		// 1 arg = key + read from file or stdin
		attr.Key = attr.Args[0]
		if attr.File == "" {
			attr.File = "-"
		}
	case 2:
		attr.Key = attr.Args[0]
		attr.Val = []byte(attr.Args[1])

		if attr.Args[1] == "-" {
			attr.File = "-"
		}
	}

	if attr.File != "" {
		if err := attr.GetFileValue(); err != nil {
			return err
		}
	}

	if attr.Binary {
		attr.Preview = "<encrypted-content>"
	} else {
		if len(attr.Val) > MaxValueWidth {
			attr.Preview = string(attr.Val)[0:MaxValueWidth] + "..."

			if strings.Contains(attr.Preview, "\n") {
				parts := strings.Split(attr.Preview, "\n")
				if len(parts) > 0 {
					attr.Preview = parts[0]
				}
			}
		} else {
			attr.Preview = string(attr.Val)
		}
	}
	if attr.Encrypted {
		attr.Preview = "<encrypted-content>"
	}

	return nil
}

func (attr *DbAttr) GetFileValue() error {
	var fd io.Reader

	if attr.File == "-" {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			fd = os.Stdin
		}
	} else {
		filehandle, err := os.OpenFile(attr.File, os.O_RDONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", attr.File, err)
		}

		fd = filehandle
	}

	if fd != nil {
		// read from file or stdin pipe
		data, err := io.ReadAll(fd)
		if err != nil {
			return fmt.Errorf("failed to read from pipe: %w", err)
		}

		// poor man's text file test
		attr.Val = data

		if utf8.ValidString(string(data)) {
			attr.Binary = false
		} else {
			attr.Binary = true
		}
	} else {
		// read from console stdin
		var input string
		var data string

		for {
			_, err := fmt.Scanln(&input)
			if err != nil {
				break
			}
			data += input + "\n"
		}

		attr.Val = []byte(data)
	}

	return nil
}
