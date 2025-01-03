package genTable

import "html/template"

type Table struct {
	TableConfigs []TableConfig
}

type TableConfig struct {
	TableName  string        // 表名
	StructName string        // 类型名
	KeyName    []string      // key字段
	KeyType    string        // key类型
	KeyDesc    template.HTML // key的描述 "+"会被转义，所以不用string类型
	Fields     []Field       // 成员

	// 导入库
	Utils bool
	Def   bool
}

type Field struct {
	Name     string
	Type     string
	Desc     string
	Col      int32  // 列数
	Package  string // 自定义结构体所在的包名，方便生成模板
	IsNumber bool   // 是否是数字
}
