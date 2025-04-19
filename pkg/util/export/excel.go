package export

import (
	"path"

	"github.com/oklog/ulid/v2"
	"github.com/xuri/excelize/v2"
)

func saveExcelFile(conf *Conf, content [][]string) (string, error) {
	file := excelize.NewFile()
	defer file.Close()

	sheetName := "Sheet1"
	index, e := file.NewSheet(sheetName)
	if e != nil {
		return "", e
	}

	content = append([][]string{conf.getAliases()}, content...)
	for rowIdx, rowData := range content {
		for colIdx, cellData := range rowData {
			axis, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			if e = file.SetCellValue(sheetName, axis, cellData); e != nil {
				return "", e
			}
		}
	}

	file.SetActiveSheet(index)
	filePath := path.Join(conf.FileBasePath, ulid.Make().String()+".xlsx")
	e = file.SaveAs(filePath)
	if e != nil {
		return "", e
	}
	return filePath, nil
}
