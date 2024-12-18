package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
)

func List(writer io.Writer, conf *cfg.Config, entries app.DbEntries) error {
	// FIXME: call sort here
	switch conf.Mode {
	case "wide":
		fallthrough
	case "":
		fallthrough
	case "table":
		return ListTable(writer, conf, entries)
	case "json":
		return ListJson(writer, conf, entries)
	default:
		return errors.New("unsupported mode")
	}

	return nil
}

func ListJson(writer io.Writer, conf *cfg.Config, entries app.DbEntries) error {
	jsonentries, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("json marshalling failure: %s", err)
	}

	fmt.Println(string(jsonentries))
	return nil
}

func ListTable(writer io.Writer, conf *cfg.Config, entries app.DbEntries) error {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	if !conf.NoHeaders {
		if conf.Mode == "wide" {
			table.SetHeader([]string{"KEY", "TAGS", "SIZE", "AGE", "VALUE"})
		} else {
			table.SetHeader([]string{"KEY", "VALUE"})
		}
	}

	for _, row := range entries {
		size := len(row.Value)

		if len(row.Bin) > 0 {
			row.Value = "binary-content"
			size = len(row.Bin)
		}

		if len(row.Value) > 60 {
			row.Value = row.Value[0:60] + "..."
		}

		if conf.Mode == "wide" {
			table.Append([]string{
				row.Key,
				strings.Join(row.Tags, ","),
				humanize.Bytes(uint64(size)),
				//row.Created.Format("02.01.2006T03:04.05"),
				humanize.Time(row.Created),
				row.Value,
			})
		} else {
			table.Append([]string{row.Key, row.Value})
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
