package tabular

import (
	"opg-infra-costs/costs"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func Table(costs []costs.CostRow) {

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Account", "Project", "Environment", "Service", "Date", "Cost")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, row := range costs {
		tbl.AddRow(
			row.Account.Id,
			row.Account.Name,
			row.Account.Environment,
			row.Service,
			row.Date,
			row.Cost)
	}
	tbl.Print()
}
