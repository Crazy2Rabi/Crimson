package Utils

import (
	"Common/Table"
	"Common/Utils/def"
	"log/slog"
	"math"
)

// 注意依赖顺序

// 根据堆叠限制组织道具
func OrganizeItems(items *def.Items) {
	var (
		mapItems = make(map[int32]int64)
		data     *Table.TbItem
		maxPile  int32
	)

	for _, item := range *items {
		if item.Num > 0 {
			mapItems[item.Id] += int64(item.Num)
		}
	}
	items.Clear()

	for id, _ := range mapItems {
		data = Table.Table.TbConfigItem.FindById(id)
		if data == nil {
			slog.Warn("organizeItems: item Id not found", slog.Any("Id", id))
			continue
		}

		if data.MaxPile > 0 {
			maxPile = data.MaxPile
		} else {
			maxPile = math.MaxInt32
		}

		// 拆分超过堆叠上限的物品
		for {
			if mapItems[id] <= 0 {
				break
			}

			var item def.Item
			item.Id = id
			if mapItems[id] > int64(maxPile) {
				item.Num = maxPile
				mapItems[id] -= int64(maxPile)
			} else {
				item.Num = int32(mapItems[id])
				mapItems[id] = 0
			}
			items.Append(&item)
		}
	}

	return
}
