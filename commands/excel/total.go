package excel

import (
	"fmt"
	"opg-infra-costs/accounts"
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"opg-infra-costs/tabular"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func totalsByMonth(spreadsheet *excelize.File,
	accounts *[]accounts.Account,
	startDate time.Time,
	endDate time.Time,
	months []time.Time,
	period string) {

	//get costs
	costData, _ := costs.AsyncCosts(
		accounts,
		startDate,
		endDate,
		period,
		"")

	headers := []string{"Date", "Cost"}
	row := []string{"Date", "Cost"}
	group := []string{"Date"}

	grouped := costData.GroupByKeys(group)
	tabular.Table(grouped, headers, row)

	groupMap := costData.GroupByKeysMap(group)
	// setup the spreadsheet
	sheet := "Totals"
	i := spreadsheet.NewSheet(sheet)
	spreadsheet.SetActiveSheet(i)

	df := dates.AWSDateFormatYM()
	col := 'A'
	// add to sheet
	for _, d := range months {
		c := string(col)
		headerCell := fmt.Sprintf("%s1", c)
		valuesCell := fmt.Sprintf("%s2", c)
		month := d.Format(df)
		key := month + "-01"
		v := groupMap[key]
		spreadsheet.SetCellValue(sheet, headerCell, month)
		spreadsheet.SetCellValue(sheet, valuesCell, v.Cost)

		col++
	}

	// save!
	if err := spreadsheet.SaveAs(spreadsheet.Path); err != nil {
		fmt.Printf("FAILED: %v\n", err)
	}
}
