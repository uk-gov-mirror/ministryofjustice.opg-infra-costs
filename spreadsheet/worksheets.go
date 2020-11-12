package spreadsheet

import (
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"sort"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// ToExcel converts the cost data to mapped and grouped version
// and generate month headers to write to file
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
	err = spreadsheet.SaveAs(spreadsheet.Path)

	return err
}
