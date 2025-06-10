/*
Copyright © 2025 Thomas von Dein

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
	"regexp"
	"strconv"
	"strings"
	tpl "text/template"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
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

	table := tablewriter.NewTable(tableString,
		tablewriter.WithRenderer(
			renderer.NewBlueprint(tw.Rendition{
				Borders: tw.BorderNone,
				Symbols: tw.NewSymbols(tw.StyleNone),
				Settings: tw.Settings{
					Separators: tw.Separators{BetweenRows: tw.Off, BetweenColumns: tw.On},
					Lines:      tw.Lines{ShowFooterLine: tw.Off, ShowHeaderLine: tw.Off},
				},
			})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoFormat: tw.Off,
				},
				Padding: tw.CellPadding{
					Global: tw.Padding{Left: "", Right: ""},
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoWrap:  tw.WrapNone,
					Alignment: tw.AlignLeft,
				},
				Padding: tw.CellPadding{
					Global: tw.Padding{Left: "", Right: ""},
				},
			},
		}),
		tablewriter.WithPadding(tw.PaddingDefault),
	)

	if !conf.NoHeaders {
		if conf.Mode == "wide" {
			table.Header([]string{"KEY", "TAGS", "SIZE", "UPDATED", "VALUE"})
		} else {
			table.Header([]string{"KEY", "VALUE"})
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

	table.Render()

	trimmer := regexp.MustCompile(`(?m)^\s*`)

	fmt.Fprint(writer, trimmer.ReplaceAllString(tableString.String(), ""))

	return nil
}
