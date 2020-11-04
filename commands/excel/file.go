package excel

import "github.com/360EntSecGroup-Skylar/excelize"

func spreadsheet() (*excelize.File, error) {
	filename := "costs.xlsx"
	f := excelize.NewFile()
	err := f.SaveAs(filename)
	f.Path = filename
	return f, err
}
