package def

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// 特殊物品Id
type ItemId = int32

const (
	GoldId    ItemId = 1001 // 金币
	DiamondId        = 1002 // 钻石
)

// 物品类型
type ItemType int32

const (
	Materials ItemType = iota // 材料
	Resources                 // 资源
	Equip                     // 装备
	MaxItemType
)

// Item 物品
type Item struct {
	Id  int32 // 物品id
	Num int32 // 物品数量
}

func (p Item) IsValid() bool {
	return p.Id > 0 && p.Num > 0
}

func NewItem(id, num int32) *Item {
	var item Item
	item.Id = id
	item.Num = num
	return &item
}

type Items []*Item

func (p Items) FindById(id int32) *Item {
	n, found := slices.BinarySearchFunc(p,
		&Item{
			Id: id,
		},
		func(a *Item, b *Item) int {
			return cmp.Compare(a.Id, b.Id)
		})

	if found {
		return p[n]
	}

	return nil
}

func (p *Items) Append(v *Item) {
	*p = append(*p, v)
}

func (p *Items) Delete(index int) {
	*p = append((*p)[:index], (*p)[index+1:]...)
}

func (p *Items) Clear() {
	*p = make(Items, 0)
}

func (p *Items) PrintInfo() {
	for _, i := range *p {
		fmt.Printf("%d_%d\n", i.Id, i.Num)
	}
}

func IsItemTypeValid(itemType ItemType) bool {
	return itemType >= 0 && itemType < MaxItemType
}

func StrToItemType(s string) (t ItemType, err error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return
	}
	if !IsItemTypeValid(t) {
		err = fmt.Errorf("strToItemType():%s not valid", s)
	}
	return ItemType(i), nil
}

func StrToItems(s string) (items Items, err error) {
	if s == "0" || s == "" {
		return
	}

	var (
		id  int64
		num int64
	)

	parts := strings.Split(s, "|")
	for _, part := range parts {
		subParts := strings.Split(part, "_")
		if len(subParts) != 2 {
			err = fmt.Errorf("strToItem():%s not valid", s)
			return
		}

		id, err = strconv.ParseInt(subParts[0], 10, 32)
		if err != nil {
			err = fmt.Errorf("strToItem():%s Id not valid", s)
			return
		}

		num, err = strconv.ParseInt(subParts[1], 10, 32)
		if err != nil {
			err = fmt.Errorf("strToItem():%s Num not valid", s)
			return
		}

		item := NewItem(int32(id), int32(num))
		items = append(items, item)
	}

	return
}
