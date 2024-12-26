package excelizeExt

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadDataRows(file *excelize.File, sheetName string, serverArr []string) error {
	// 获取数据行开始的行数
	dataLine, err := FindFirstCellValueInCol(file, sheetName, 0, "")
	if err != nil {
		return err
	}
	if dataLine == -1 {
		return nil
	}

	// 读取每一行
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return err
	}
	for i := dataLine; i < len(rows); i++ {
		buf := make(map[string]interface{})
		for j := 1; j < len(serverArr); j++ {
			if j >= len(rows[i]) || rows[i][j] == "" {
				return fmt.Errorf("表%s %d行 有空单元格", sheetName, i)
			}
			buf[serverArr[j]] = rows[i][j]
		}
	}

	return nil
}

func FindFristCellValue(file *excelize.File, sheetName, value string) (row, col int, err error) {
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return -1, -1, err
	}

	for i, row := range rows {
		for j, cell := range row {
			if cell == value {
				return i, j, nil
			}
		}
	}

	return -1, -1, nil
}

// 若查找到，返回行数（0开始）；未查找到，返回-1
func FindFirstCellValueInCol(file *excelize.File, sheetName string, colName int, value string) (row int, err error) {
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return -1, err
	}

	for i, row := range rows {
		if row == nil {
			err = fmt.Errorf("表%s 有空行 第%d行", sheetName, i+1)
			return -1, err
		}
		if row[colName] == value {
			return i, nil
		}
	}

	return -1, nil
}

func GetRow(file *excelize.File, sheetName string, rowIndex int) (row []string, err error) {
	if rowIndex < 0 {
		err = fmt.Errorf("rowIndex < 0")
		return
	}

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return
	}

	if rowIndex >= len(rows) {
		err = fmt.Errorf("rowIndex out of range")
		return
	}
	return rows[rowIndex], nil
}
