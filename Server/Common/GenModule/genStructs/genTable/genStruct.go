package genTable

import (
	"Common/Framework/config"
	"Common/GenModule/GenFile"
	"Common/Utils/excelizeExt"
	"fmt"
	"github.com/xuri/excelize/v2"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var Tables Table
var tablePath string
var genTablePath string

func LoadTableConfigs() error {
	startTime := time.Now()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	tablePath = filepath.Join(dir, config.Instance().GenModuleConfig.TablesPath)
	genTablePath = filepath.Join(dir, config.Instance().GenModuleConfig.TablesStructPath)

	// 遍历文件夹，直接读取
	err = filepath.Walk(tablePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".xlsx" {
			err = ReadExcel(tablePath, info.Name())
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// 生成Tables结构体
	GenTables(genTablePath)

	tc := time.Since(startTime)
	fmt.Printf("gen table structs ok! Total cost time %v\n\n", tc)
	return nil
}

// ReadExcel 生成表中的结构体，不读数据
// keyCol 作为键的列，如果没有则为第一列
func ReadExcel(dirPath, fileName string) error {
	startTime := time.Now()

	if fileName == "" {
		return fmt.Errorf("genTable: ReadExcel() filename is empty")
	}

	filePath := filepath.Join(dirPath, fileName)

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("genTable: ReadExcel() open file error: %w", err)
	}

	for _, sheet := range file.GetSheetList() {
		err := ReadSheet(file, sheet)
		if err != nil {
			return fmt.Errorf("genTable: ReadExcel() 表格%s读取失败 %w\n", fileName, err)
		}
	}
	tc := time.Since(startTime)
	fmt.Printf("Parse <%s> structs. Time cost %v\n", fileName, tc)
	return nil
}

func ReadSheet(file *excelize.File, sheet string) (err error) {
	// A1单元格为表名，如果为空则跳过该sheet
	tableName, err := file.GetCellValue(sheet, "A1")
	if err != nil {
		return err
	}
	if tableName == "" {
		return
	}

	// 第一列没有SERVER关键字，跳过
	serverLine, err := excelizeExt.FindFirstCellValueInCol(file, sheet, 0, "SERVER")
	if err != nil {
		return err
	}
	if serverLine == -1 {
		return
	}

	// 第一列没有TYPE关键字，跳过
	typeLine, err := excelizeExt.FindFirstCellValueInCol(file, sheet, 0, "TYPE")
	if err != nil {
		return err
	}
	if typeLine == -1 {
		return
	}

	keyLine, err := excelizeExt.FindFirstCellValueInCol(file, sheet, 0, "KEY")
	if err != nil {
		return err
	}

	serverRow, err := excelizeExt.GetRow(file, sheet, serverLine)
	if err != nil {
		return err
	}
	typeRow, err := excelizeExt.GetRow(file, sheet, typeLine)
	if err != nil {
		return err
	}
	if len(serverRow) > len(typeRow) {
		return fmt.Errorf("表格%s SERVER行数量 > TYPE行数量，检查是否少定义TYPE行或多定义SERVER行\n", sheet)
	}

	var keyRow []string
	if keyLine >= 0 {
		keyRow, err = excelizeExt.GetRow(file, sheet, keyLine)
		if err != nil {
			return err
		}
		if len(keyRow) > len(typeRow) {
			return fmt.Errorf("表格%s KEY行数 > TYPE行数，检查是否有多余KEY行", sheet)
		}
	}

	var descRow []string
	descRow, err = excelizeExt.GetRow(file, sheet, 0)
	if err != nil {
		return err
	}

	// 生成结构体
	if err := GenTableStruct(tableName, serverRow, typeRow, descRow, keyRow); err != nil {
		return err
	}

	return
}

func GenTableStruct(tableName string, serverArr, typeArr, descArr, keyArr []string) error {
	var table TableConfig
	// 首字母大写，加上前缀Tb
	table.TableName = tableName
	table.StructName = strings.Title(tableName)
	for i := 1; i < len(serverArr); i++ {
		if serverArr[i] == "" {
			continue
		}
		if typeArr[i] == "" {
			return fmt.Errorf("表格%s %d列 SERVER行有定义，而TYPE无定义", tableName, i+1)
		}

		isNumber := true
		name := strings.Title(serverArr[i])
		// 下面这个函数，根据语言，单词中的首字母大写，但单词其他字母会小写
		// name := cases.Title(language.English).String(serverArr[i])
		desc := ""
		if i < len(descArr) {
			desc = descArr[i]
			desc = strings.Replace(desc, "\n", " ", -1)
		}

		// 把类型转换一下
		convertType := typeArr[i]
		packageName := ""
		if v, ok := GenFile.TypeConverter[typeArr[i]]; ok {
			convertType = v
			if v == "string" {
				isNumber = false
			}
		} else {
			// 不在转换容器中的，统一当成自定义类型，修改为 def.类型
			packageName = "def."
			isNumber = false

			if !table.Def {
				table.Def = true
			}
		}

		field := Field{
			Name:     name,
			Type:     convertType,
			Desc:     desc,
			Col:      int32(i),
			Package:  packageName,
			IsNumber: isNumber,
		}

		table.Fields = append(table.Fields, field)
	}

	// 处理表的key类型
	table.KeyType = "int32"
	if len(keyArr) > 0 {
		keyNum := 0
		for i, key := range keyArr {
			if key == "1" {
				keyNum++
				if i > 1 {
					table.KeyDesc += " + "
				}
				table.KeyName = append(table.KeyName, serverArr[i])
				table.KeyDesc += template.HTML(descArr[i])
			}
		}

		if keyNum > 1 {
			table.KeyType = "int64"

			if !table.Def {
				table.Def = true
			}
		}
	} else {
		// 默认表的第2列作为key
		table.KeyName = append(table.KeyName, serverArr[1])
		table.KeyDesc = template.HTML(descArr[1])
	}

	path := filepath.Join(genTablePath, fmt.Sprintf("%s.go", tableName))
	GenFile.GenFile(path, TableConfigText, table, template.FuncMap{})

	Tables.TableConfigs = append(Tables.TableConfigs, TableConfig{
		TableName:  table.TableName,
		StructName: table.StructName,
	})

	return nil
}

func GenTables(genDirPath string) {
	genFilePath := filepath.Join(genDirPath, "table.go")
	GenFile.GenFile(genFilePath, TableText, Tables, template.FuncMap{})
}
