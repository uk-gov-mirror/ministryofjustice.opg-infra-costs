package spreadsheet

import "github.com/360EntSecGroup-Skylar/excelize"

func New() (*excelize.File, error) {
	filename := "costs.xlsx"
	f := excelize.NewFile()
	// add introduction content here?
	i := f.NewSheet("Overview")
	f.SetActiveSheet(i)
	f.DeleteSheet("Sheet1")

	err := f.SaveAs(filename)
	f.Path = filename
	return f, err
}
