package manager

import (
	"Common/Table"
	"Common/Utils"
	"Common/Utils/helper"
	"Common/def"
	"Game/player"
	"cmp"
	"fmt"
	"slices"
)

func AddItems(p *player.Player, items *def.Items) (err error) {
	Utils.OrganizeItems(items)

	// todo 背包容量判断

	for _, item := range *items {
		err = addItem(p, item.Id, item.Num)
		if err != nil {
			return err
		}
	}

	slices.SortFunc(p.Items, func(a, b *def.Item) int {
		if a.Id == b.Id {
			return cmp.Compare(b.Num, a.Num)
		}
		return cmp.Compare(a.Id, b.Id)
	})

	// todo 同步消息

	return
}

func SubItems(p *player.Player, items *def.Items) (err error) {
	if !HasEnoughItems(p, *items) {
		err = fmt.Errorf("SubItems: no enough items")
		return
	}

	Utils.OrganizeItems(items)

	for _, item := range *items {
		err = subItem(p, item.Id, item.Num)
		if err != nil {
			return
		}
	}

	// todo 同步消息

	return
}

func addItem(p *player.Player, id, num int32) (err error) {
	if num <= 0 {
		err = fmt.Errorf("AddItem: num <= 0")
		return
	}

	data := Table.Table.TbConfigItem.FindById(id)
	if data == nil {
		err = fmt.Errorf("AddItem: item id not found")
		return
	}

	switch id {
	case def.GoldId:
		p.Gold = helper.SafeAdd(p.Gold, num)
	case def.DiamondId:
		p.Diamond = helper.SafeAdd(p.Diamond, num)
	default:
		// 道具合并
		if data.MaxPile > 1 {
			for _, item := range p.Items {
				if item.Num+num <= data.MaxPile {
					item.Num = item.Num + num
					return
				} else if num < data.MaxPile {
					num -= data.MaxPile - item.Num
					item.Num = data.MaxPile
				}
			}
		}

		// 道具堆叠已经在外面判断
		if num > 0 {
			p.Items.Append(def.NewItem(id, num))
		}
	}

	return
}

func subItem(p *player.Player, id, num int32) (err error) {
	if num <= 0 {
		err = fmt.Errorf("SubItem: num <= 0")
		return
	}

	data := Table.Table.TbConfigItem.FindById(id)
	if data == nil {
		err = fmt.Errorf("SubItem: item id not found")
		return
	}

	// 道具是否充足已经在外面判断
	switch id {
	case def.GoldId:
		p.Gold = helper.SafeSub(p.Gold, num)
	case def.DiamondId:
		p.Diamond = helper.SafeSub(p.Diamond, num)
	default:
		for i := len(p.Items) - 1; i >= 0; i-- {
			if p.Items[i].Id != id {
				continue
			}

			if p.Items[i].Num-num > 0 {
				p.Items[i].Num -= num
				return
			} else {
				num -= p.Items[i].Num
				p.Items.Delete(i)
			}

			if num <= 0 {
				break
			}
		}
	}

	return
}

func HasEnoughItems(p *player.Player, items def.Items) bool {
	var mapItems = make(map[int32]int32)

	for _, item := range items {
		if item.Num > 0 {
			mapItems[item.Id] += item.Num
		}
	}

	for id, num := range mapItems {
		if !HasEnoughItem(p, id, num) {
			return false
		}
	}

	return true
}

func HasEnoughItem(p *player.Player, Id, num int32) bool {
	switch Id {
	case def.GoldId:
		if p.Gold >= num {
			return true
		}
	case def.DiamondId:
		if p.Diamond >= num {
			return true
		}
	default:
		for _, i := range p.Items {
			if i.Id == Id {
				num -= i.Num
				if num <= 0 {
					return true
				}
			}
		}
	}
	return false
}

func PrintItemsInfo(p player.Player) {
	fmt.Println("当前拥有物品：")
	p.Items.PrintInfo()
	fmt.Println("Gold:", p.Gold)
	fmt.Println("Diamond:", p.Diamond)
}
