package main

import (
	"Common/Utils/def"
	"Common/Utils/random"
	"fmt"
)

func main() {
	randomPool := def.Items{
		def.NewItem(1001, 100),
		def.NewItem(1002, 200),
		def.NewItem(1003, 300),
	}

	var results map[int32]int
	results = make(map[int32]int)

	Test2(randomPool, results)

	for i, result := range results {
		fmt.Printf("%d:%d\n", i, result)
	}
}

// 测试GetRandomValueFromItems
// 概率 1:2:3
func Test1(randomPool def.Items, results map[int32]int) {
	for i := 0; i < 60000; i++ {
		result := random.GetRandomValueFromItems(randomPool)
		results[result]++
	}
}

// 测试GetRandomValuesFromItems 不放回
// 概率 25:44:51
func Test2(randomPool def.Items, results map[int32]int) {
	for i := 0; i < 60000; i++ {
		result := random.GetRandomValuesFromItems(randomPool, 2, false)
		for _, v := range result {
			results[v]++
		}
	}
}
