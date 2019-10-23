package football

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func createTable(header []string, buf *bytes.Buffer) *tablewriter.Table {

	table := tablewriter.NewWriter(buf)
	if len(header) != 0 {
		table.SetHeader(header)
	}
	table.SetBorder(false)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	return table
}

func parseDateString(d string) string {

	return fmt.Sprintf("🕙 %s", strings.Replace(d, "T", " ", 1))
}
