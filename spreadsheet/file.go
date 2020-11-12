package spreadsheet

import "github.com/360EntSecGroup-Skylar/excelize/v2"

func New() (*excelize.File, error) {
	filename := "costs.xlsx"
	f := excelize.NewFile()
	err := f.SaveAs(filename)
	f.Path = filename
	return f, err
}
