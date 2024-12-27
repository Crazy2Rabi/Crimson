package random

import (
	"Common/def"
	"math/rand"
	"time"
)

var s = rand.New(rand.NewSource(time.Now().UnixNano()))

/*func NewRand() *rand.Rand {
	source := rand.NewSource(time.Now().UnixNano())
	return rand.New(source)
}*/

// [0, n)
func GetRandLess(n int32) int32 {
	return s.Int31n(n)
}

// [A, B]
func GetRandBetween(a, b int32) int32 {
	if a > b {
		return 0
	}

	if a == b {
		return a
	}

	return a + s.Int31n(b-a+1)
}

// 根据权重随机一个
func GetRandomValueFromItems(items def.Items) int32 {
	var (
		itemsTemp   def.Items
		totalWeight int32
	)

	for _, item := range items {
		if item.IsValid() {
			totalWeight += item.Num
			itemsTemp = append(itemsTemp, def.NewItem(item.Id, totalWeight))
		}
	}

	if totalWeight <= 0 {
		return 0
	}

	randVal := GetRandLess(totalWeight)
	for _, item := range itemsTemp {
		if randVal < item.Num {
			return item.Id
		}
	}

	return 0
}

// 根据权重随机多个，数量不足直接返回
// replaceBack 放回与不放回
func GetRandomValuesFromItems(items def.Items, count int32, replaceBack bool) (results []int32) {
	if count <= 0 {
		return
	}

	var (
		itemsTemp   def.Items
		totalWeight int32
		num         int32
	)

	for _, item := range items {
		if item.IsValid() {
			totalWeight += item.Num
			itemsTemp = append(itemsTemp, def.NewItem(item.Id, totalWeight))
			//fmt.Printf("%d_%d", item.Id, item.Num)
		}
	}

	// 不放回 数量不足，直接返回
	if !replaceBack && count >= int32(len(itemsTemp)) {
		for _, item := range itemsTemp {
			results = append(results, item.Id)
		}
		return
	}

	for {
		var (
			index   int
			randVal int32
		)
		randVal = GetRandLess(totalWeight)
		for i, item := range itemsTemp {
			if randVal < item.Num {
				results = append(results, item.Id)
				num++
				index = i
				break
			}
		}

		if num == count {
			return
		}

		// 不放回
		if !replaceBack {
			for i := index + 1; i < len(itemsTemp); i++ {
				itemsTemp[i].Num -= itemsTemp[index].Num
				if index > 0 {
					itemsTemp[i].Num += itemsTemp[index-1].Num
				}
			}
			itemsTemp.Delete(index)
			totalWeight = itemsTemp[len(itemsTemp)-1].Num
		}
	}
}
