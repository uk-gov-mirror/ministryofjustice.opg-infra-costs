package spreadsheet

import (
	"fmt"
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"sort"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func ToExcel(
	spreadsheet *excelize.File,
	sheet string,
	costData *costs.CostData,
	months []time.Time,
	groupBy []string,
	headers []string,
) error {

	groupMap := costData.GroupByKeysMap(groupBy)

	df := dates.AWSDateFormatYM()
	m := []string{}
	for _, k := range months {
		m = append(m, k.Format(df))
	}
	//order the months
	sort.Strings(m)

	err := writeDataToSheet(
		spreadsheet,
		sheet,
		headers,
		m,
		groupMap)

	// save!
	if err = spreadsheet.SaveAs(spreadsheet.Path); err != nil {
		fmt.Printf("FAILED: %v\n", err)
	}
	return err
}
