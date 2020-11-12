package spreadsheet

import (
	"fmt"
	"opg-infra-costs/costs"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func writeDataToSheet(
	spreadsheet *excelize.File,
	sheet string,
	headers []string,
	dates []string,
	data map[string]costs.CostRow) error {

	// create the new sheet
	i := spreadsheet.NewSheet(sheet)
	if sheet == "Totals" {
		spreadsheet.SetActiveSheet(i)
		//remove sheet1 default
		spreadsheet.DeleteSheet("Sheet1")
	}

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
	maxcol := 'A'

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
	// add sparkline col
	spreadsheet.SetCellValue(
		sheet,
		fmt.Sprintf("%s%v", string(col), row),
		"Trend")

	maxcol = col
	// sparkline data
	sparks := map[string][]string{
		"Location": []string{},
		"Range":    []string{}}

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
			v, _ := strconv.ParseFloat(dataRow[m], 64)
			spreadsheet.SetCellValue(sheet, cell, v)
			col++
		}

		// add sparkline setup
		sparks["Location"] = append(
			sparks["Location"],
			fmt.Sprintf("%v%v",
				string(col),
				row))

		sparks["Range"] = append(
			sparks["Range"],
			fmt.Sprintf("%v!A%v:%v%v",
				sheet,
				row,
				string(col-1),
				row))

	}

	// add spark lines to the sheet
	spreadsheet.AddSparkline(sheet, &excelize.SparklineOption{
		Location: sparks["Location"],
		Range:    sparks["Range"],
		Markers:  true,
		Type:     "column",
		Style:    3})

	// add the table and filtering options
	max := fmt.Sprintf("%v%v", string(maxcol), row)
	spreadsheet.AddTable(
		sheet,
		"A1",
		max,
		`{
			"table_style": "TableStyleMedium9"
		}`)
	// ignore the sparkline col
	spreadsheet.AutoFilter(
		sheet,
		"A1",
		fmt.Sprintf("%v%v", string(maxcol-1), row),
		"")

	return nil
}
