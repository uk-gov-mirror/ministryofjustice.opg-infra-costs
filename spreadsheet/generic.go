package spreadsheet

import (
	"fmt"
	"opg-infra-costs/costs"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func writeDataToSheet(
	spreadsheet *excelize.File,
	sheet string,
	headers []string,
	dates []string,
	data map[string]costs.CostRow) error {

	// create the new sheet
	i := spreadsheet.NewSheet(sheet)
	spreadsheet.SetActiveSheet(i)

	// restructure data to map to excel layout
	excel := map[string]map[string]string{}
	for _, d := range data {
		key := costs.GenerateGroupKey(headers, &d)

		m := make(map[string]string)
		// if we have made it, set m to be equal to it
		// else add header values to it
		if f, ok := excel[key]; ok {
			m = f
		} else {
			for _, h := range headers {
				m[h] = d.Get(h)
			}
		}
		// AWS result to month (2020-01-01 => 2020-01)
		df := d.Date[:len(d.Date)-3]
		m[df] = fmt.Sprintf("%.2f", d.Cost)
		excel[key] = m
	}

	// now write headers to spreadsheet
	row := 1
	col := 'A'
	for _, h := range headers {
		c := string(col)
		cell := fmt.Sprintf("%s%v", c, row)
		spreadsheet.SetCellValue(sheet, cell, h)
		col++
	}
	for _, m := range dates {
		c := string(col)
		cell := fmt.Sprintf("%s%v", c, row)
		spreadsheet.SetCellValue(sheet, cell, m)
		col++
	}

	// now write data to spreadsheet
	for _, dataRow := range excel {
		// set col and row
		col = 'A'
		row++
		// headers
		for _, h := range headers {
			c := string(col)
			cell := fmt.Sprintf("%s%v", c, row)
			spreadsheet.SetCellValue(sheet, cell, dataRow[h])
			col++
		}
		// dates (and costs)
		for _, m := range dates {
			c := string(col)
			cell := fmt.Sprintf("%s%v", c, row)
			spreadsheet.SetCellValue(sheet, cell, dataRow[m])
			col++
		}

	}

	return nil
}
