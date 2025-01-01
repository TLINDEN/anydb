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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	tpl "text/template"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
)

func List(writer io.Writer, conf *cfg.Config, entries app.DbEntries) error {
	switch conf.Mode {
	case "wide", "", "table":
		return ListTable(writer, conf, entries)
	case "json":
		return ListJson(writer, conf, entries)
	case "template":
		return ListTemplate(writer, conf, entries)
	default:
		return errors.New("unsupported mode")
	}
}

func ListJson(writer io.Writer, conf *cfg.Config, entries app.DbEntries) error {
	jsonentries, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("failed marshall json: %s", err)
	}

	fmt.Println(string(jsonentries))
	return nil
}

func ListTemplate(writer io.Writer, conf *cfg.Config, entries app.DbEntries) error {
	tmpl, err := tpl.New("list").Parse(conf.Template)
	if err != nil {
		return fmt.Errorf("failed to parse output template: %w", err)
	}

	buf := bytes.Buffer{}

	for _, row := range entries {
		buf.Reset()
		err = tmpl.Execute(&buf, row)
		if err != nil {
			return fmt.Errorf("failed to execute output template: %w", err)
		}

		if buf.Len() > 0 {
			fmt.Fprintln(writer, buf.String())
		}
	}

	return nil
}

func ListTable(writer io.Writer, conf *cfg.Config, entries app.DbEntries) error {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	if !conf.NoHeaders {
		if conf.Mode == "wide" {
			table.SetHeader([]string{"KEY", "TAGS", "SIZE", "UPDATED", "VALUE"})
		} else {
			table.SetHeader([]string{"KEY", "VALUE"})
		}
	}

	for _, row := range entries {
		if conf.Mode == "wide" {
			switch conf.NoHumanize {
			case true:
				table.Append([]string{
					row.Key,
					strings.Join(row.Tags, ","),
					strconv.FormatUint(row.Size, 10),
					row.Created.AsTime().Format("02.01.2006T03:04.05"),
					row.Preview,
				})
			default:
				table.Append([]string{
					row.Key,
					strings.Join(row.Tags, ","),
					humanize.Bytes(uint64(row.Size)),
					humanize.Time(row.Created.AsTime()),
					row.Preview,
				})
			}

		} else {
			table.Append([]string{row.Key, row.Preview})
		}
	}

	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetNoWhiteSpace(true)

	table.SetTablePadding("\t") // pad with tabs

	table.Render()

	fmt.Fprint(writer, tableString.String())

	return nil
}
