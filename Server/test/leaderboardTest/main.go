package main

import (
	"Common/Utils/def"
	"Common/Utils/random"
	"fmt"
	"github.com/jinzhu/now"
	"time"
)

func main() {
	Test2()
}

// 基础测试
func Test1() {
	/*rankInfos := []def.RankInfo{
		{Uid: 1001, Value: 400, Time: 1},
		{Uid: 1002, Value: 100},
		{Uid: 1003, Value: 500},
		{Uid: 1004, Value: 200},
		{Uid: 1005, Value: 400, Time: 5},
		{Uid: 1006, Value: 400, Time: 3},
	}*/

	var rankList def.RankingList
	//rankList.Init(10000)

	// 初始化
	/*for _, info := range rankInfos {
		rankList.Put(&info)
	}*/
	fmt.Println("初始化：")
	rankList.PrintInfo()

	fmt.Println("size:", rankList.RankInfos.Size())
	fmt.Println("left:", rankList.RankInfos.Left())
	fmt.Println("right:", rankList.RankInfos.Right())
	fmt.Println("Get 1005:", rankList.GetNode(1005))
	fmt.Println("Get 0:", rankList.GetNode(0))

	/*fmt.Println("Remove 1006:")
	rankList.Remove(&def.RankInfo{Uid: 1006})
	rankList.PrintInfo()*/

	fmt.Println("update 1005:")
	//rankList.Put(&def.RankInfo{Uid: 1005, Value: 500, Time: 6})
	rankList.PrintInfo()
}

func Test2() {
	var rankList def.RankingList
	//rankList.Init(10000)

	// 10w次入榜操作
	startTime := time.Now()
	for i := 1; i <= 100000; i++ {
		nowTime := now.New(time.Now())

		var tmp def.RankInfo
		tmp.Uid = uint64(random.GetRandLess(50000))
		tmp.Time = nowTime.Unix()
		tmp.Value = int64(random.GetRandLess(int32(i)))
		rankList.Put(&tmp)
	}
	endTime := time.Now()
	fmt.Println("插入10w次，共耗时：", endTime.Sub(startTime))
	rankList.PrintInfo()
}
