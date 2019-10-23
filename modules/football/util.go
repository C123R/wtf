package football

import (
	"bytes"
	"fmt"
	"strings"
	"time"

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

func getDateFrame(offset int) (string, string) {

	today := time.Now()

	if offset >= 1 {
		return today.Format("2006-01-02"), today.AddDate(0, 0, offset).Format("2006-01-02")
	} else if offset < 0 {
		return today.AddDate(0, 0, offset).Format("2006-01-02"), today.Format("2006-01-02")
	}

	return today.Format("2006-01-02"), today.Format("2006-01-02")
}
