package Table

import (
	"Common/Utils/def"
	"encoding/json"
	"fmt"
	"strconv"
)

// Code generated by genTable/gen.
// DO NOT EDIT.

type TbItem struct {
	Id      int32        // Id
	Name    string       // 名字
	Type    def.ItemType // 大类 0=普通道具 1=资源
	MaxPile int32        // 最大堆叠
	Prize   def.Items    // 复合字符串
}

type TbConfigItem struct {
	dataMap  map[int32]*TbItem // key = Id
	dataList []*TbItem
}

func (t *TbConfigItem) Clear() {
	t.dataMap = make(map[int32]*TbItem)
	t.dataList = make([]*TbItem, 0)
}

func (t *TbConfigItem) ParseRow(row []string) (buf map[string]interface{}, err error) {
	buf = make(map[string]interface{})

	buf["Id"], err = strconv.ParseInt(row[1], 10, 32)
	if err != nil {
		return nil, err
	}

	buf["Name"] = row[2]

	buf["Type"] = row[3]

	buf["MaxPile"], err = strconv.ParseInt(row[4], 10, 32)
	if err != nil {
		return nil, err
	}

	buf["Prize"] = row[5]
	return
}

func (t *TbConfigItem) Add(buf map[string]interface{}) (err error) {
	v := &TbItem{}
	{
		var val string
		val = buf["Id"].(json.Number).String()
		var n int64
		n, err = strconv.ParseInt(val, 10, 32)
		v.Id = int32(n)
	}
	{
		var val string
		var ok bool
		if val, ok = buf["Name"].(string); !ok {
			err = fmt.Errorf("Name error")
			return
		}
		v.Name = val

	}
	{
		var val string
		var ok bool
		if val, ok = buf["Type"].(string); !ok {
			err = fmt.Errorf("Type error")
			return
		}
		v.Type, err = def.StrToItemType(val)
		if err != nil {
			err = fmt.Errorf("Type error: %v", err)
			return
		}

	}
	{
		var val string
		val = buf["MaxPile"].(json.Number).String()
		var n int64
		n, err = strconv.ParseInt(val, 10, 32)
		v.MaxPile = int32(n)
	}
	{
		var val string
		var ok bool
		if val, ok = buf["Prize"].(string); !ok {
			err = fmt.Errorf("Prize error")
			return
		}
		v.Prize, err = def.StrToItems(val)
		if err != nil {
			err = fmt.Errorf("Prize error: %v", err)
			return
		}

	}
	key := v.Id
	if _, ok := t.dataMap[key]; ok {
		err = fmt.Errorf("Duplicate elements")
		return
	}

	t.dataMap[key] = v
	t.dataList = append(t.dataList, v)
	return
}

func (t *TbConfigItem) FindById(id int32) *TbItem {
	return t.dataMap[id]
}

func NewItemConfig(bufs []map[string]interface{}) (t *TbConfigItem, err error) {
	t = &TbConfigItem{}
	t.Clear()

	for _, buf := range bufs {
		if err = t.Add(buf); err != nil {
			return
		}
	}
	return
}