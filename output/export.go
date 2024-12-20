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
	"os"

	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
)

func WriteJSON(attr *app.DbAttr, conf *cfg.Config, entries app.DbEntries) error {
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
