package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
)

func List(writer io.Writer, conf *cfg.Config, entries app.DbEntries) error {
	// FIXME: call sort here
	// FIXME: check output mode switch to subs

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	if conf.Mode == "wide" {
		table.SetHeader([]string{"KEY", "VALUE", "TAGS", "TIMESTAMP"})
	} else {
		table.SetHeader([]string{"KEY", "VALUE"})
	}

	for _, row := range entries {
		if row.Value == "" {
			row.Value = string(row.Bin)[0:60]
		} else if len(row.Value) > 60 {
			row.Value = row.Value[0:60]
		}

		if conf.Mode == "wide" {
			table.Append([]string{row.Key, row.Value, strings.Join(row.Tags, ","), row.Created.Format("02.01.2006T03:04.05")})
		} else {
			table.Append([]string{row.Key, row.Value})
		}
	}

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
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
