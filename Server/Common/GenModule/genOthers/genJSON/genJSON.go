package genJSON

import (
	"Common/Framework/config"
	"Common/Table"
	"Common/Utils/excelizeExt"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

var tables Table.Tables
var tablePath string
var genJSONDirPath string

// GenJSON 需要table.go生成好后才能执行
func GenJSON() error {
	startTime := time.Now()

	// 相对于build.bat中的路径
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// fmt.Println("当前目录：", dir)
	tablePath = filepath.Join(dir, config.Instance().GenModuleConfig.TablesPath)
	genJSONDirPath = filepath.Join(dir, config.Instance().GenModuleConfig.TablesJSONPath)

	// 遍历文件夹下所有表格
	err = filepath.Walk(tablePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".xlsx" {
			err = ReadExcel(tablePath, info.Name())
			if err != nil {
				return fmt.Errorf("表%s 生成json失败 %w", info.Name(), err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("gen JSON ok! Total cost time %v\n\n", time.Since(startTime))
	return nil
}

func ReadExcel(dirPath, fileName string) error {
	filePath := filepath.Join(dirPath, fileName)

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开表格失败 %w", err)
	}

	for _, sheet := range file.GetSheetList() {
		err := ReadSheet(file, fileName, sheet)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadSheet(file *excelize.File, fileName, sheet string) error {
	// A1单元格为表名，如果为空则跳过该sheet
	tableName, err := file.GetCellValue(sheet, "A1")
	if err != nil {
		return err
	}
	if tableName == "" {
		return nil
	}

	// 第一列没有SERVER关键字，跳过
	serverLine, err := excelizeExt.FindFirstCellValueInCol(file, sheet, 0, "SERVER")
	if err != nil {
		return err
	}
	if serverLine == -1 {
		return nil
	}

	// 查找Tables结构体中是否有该类型
	rt := reflect.TypeOf(tables)
	fieldName := fmt.Sprintf("TbConfig" + strings.Title(tableName))
	field, found := rt.FieldByName(fieldName)
	if !found {
		return fmt.Errorf("tables中未找到反射 %s 请重新生成表结构", fieldName)
	}

	// 初始化成员
	fieldValue := reflect.ValueOf(&tables).Elem().FieldByName(fieldName)
	if !fieldValue.IsValid() || fieldValue.IsNil() {
		fieldValue.Set(reflect.New(field.Type.Elem()))
	}

	// 检查是否有实现ParseRow接口
	method := fieldValue.MethodByName("ParseRow")
	if !method.IsValid() {
		return fmt.Errorf("%s 未实现ParseRow", fieldName)
	}

	// 调用ParseRow方法，读取每一行
	err = ReadRows(file, fileName, sheet, method)
	if err != nil {
		return err
	}

	return nil
}

func ReadRows(file *excelize.File, fileName, sheet string, method reflect.Value) error {
	startTime := time.Now()
	// 找到数据开始行
	dataLine, err := excelizeExt.FindFirstCellValueInCol(file, sheet, 0, "")
	if err != nil {
		return err
	}
	if dataLine == -1 {
		return nil
	}

	rows, err := file.GetRows(sheet)
	if err != nil {
		return err
	}

	// 调用ParseRow 解析数据
	var bufs []map[string]interface{}
	bufs = make([]map[string]interface{}, len(rows)-dataLine)
	for i := dataLine; i < len(rows); i++ {
		bufs[i-dataLine] = make(map[string]interface{})
		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(rows[i])
		rs := method.Call(params)
		if !rs[1].IsNil() {
			return fmt.Errorf("行数%d\n%w", i+1, rs[1].Interface().(error))
		}
		bufs[i-dataLine] = rs[0].Interface().(map[string]interface{})
	}

	// 生成JSON文件
	jsonData, err := json.MarshalIndent(bufs, "", "	")
	if err != nil {
		return err
	}
	ext := filepath.Ext(fileName)
	filePath := filepath.Join(genJSONDirPath, fileName[:len(fileName)-len(ext)]+".json")
	// 生成中间的文件夹
	if err = os.MkdirAll(genJSONDirPath, 0755); err != nil {
		debug.PrintStack()
		panic(err)
	}
	err = os.WriteFile(filePath, jsonData, 0644)

	fmt.Printf("Gen <%s> JSON. Time cost %v\n", fileName, time.Since(startTime))
	return nil
}
