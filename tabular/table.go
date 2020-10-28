package tabular

import (
	"opg-infra-costs/costs"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func Table(costs costs.CostData, headers []string, rowKeys []string) {
	// convert the string array to interfaces for the table new command
	headersInterfaces := make([]interface{}, len(headers))
	for i := range headers {
		headersInterfaces[i] = headers[i]
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New(headersInterfaces...)
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, row := range costs.Entries {
		//
		inter := make([]interface{}, len(rowKeys))
		for x := range rowKeys {
			inter[x] = row.Get(rowKeys[x])
		}

		tbl.AddRow(inter...)
	}
	tbl.Print()
}
