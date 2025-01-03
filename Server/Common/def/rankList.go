package def

import (
	"fmt"
	"github.com/emirpasic/gods/trees/redblacktree"
)

// 排行榜类型
type RankListType int32

const (
	CombatEffectiveness RankListType = iota // 战力
	MaxRankListType
)

type RankInfo struct {
	Uid   uint64
	Value int64
	Time  int64 // 上榜时间
}

// 排行榜
type RankingList struct {
	RankInfos *redblacktree.Tree
	uids      map[uint64]*RankInfo
	MaxSize   int
}

func IsRankListTypeValid(listType RankListType) bool {
	return listType >= 0 && listType < MaxRankListType
}

// 降序比较
func CompareRankInfoDesc(a, b interface{}) int {
	rankA, okA := a.(*RankInfo)
	rankB, okB := b.(*RankInfo)
	if !okA || !okB {
		return 0
	}

	if rankA.Value == rankB.Value {
		switch {
		case rankA.Time > rankB.Time:
			return -1
		case rankA.Time < rankB.Time:
			return 1
		default:
			return 0
		}
	}

	switch {
	case rankA.Value > rankB.Value:
		return 1
	case rankA.Value < rankB.Value:
		return -1
	default:
		return 0
	}
}

func (l *RankingList) init(size int) (err error) {
	if size <= 0 {
		err = fmt.Errorf("MaxSize <= 0")
	}

	l.MaxSize = size
	l.RankInfos = redblacktree.NewWith(CompareRankInfoDesc)
	l.uids = make(map[uint64]*RankInfo)

	return
}

// 查找
func (l *RankingList) GetNode(uid uint64) *RankInfo {
	if node, ok := l.uids[uid]; ok {
		return node
	}

	return nil
}

// 入榜
func (l *RankingList) Put(info *RankInfo) {
	// 在榜单中，更新排名
	if node, ok := l.uids[info.Uid]; ok {
		l.remove(node)
		l.put(info)
		return
	}

	// 排行榜已满
	if l.RankInfos.Size() >= l.MaxSize {
		if CompareRankInfoDesc(info, l.RankInfos.Left()) < 0 {
			return
		}

		l.put(info)
		return
	}

	l.put(info)
}

func (l *RankingList) reSize() {
	for l.RankInfos.Size() > l.MaxSize {
		l.remove(l.RankInfos.Left().Key.(*RankInfo))
	}
}

func (l *RankingList) put(info *RankInfo) {
	l.RankInfos.Put(info, nil)
	l.uids[info.Uid] = info
	l.reSize()
}

func (l *RankingList) remove(info *RankInfo) {
	l.RankInfos.Remove(l.uids[info.Uid])
	delete(l.uids, info.Uid)
}

func (l *RankingList) PrintInfo() {
	fmt.Println(l.RankInfos.String())
	fmt.Println("maxsize:", l.MaxSize)
	fmt.Println("size:", l.RankInfos.Size())
}

// 所有排行榜
type RankingLists map[RankListType]*RankingList

func (l *RankingLists) GetList(listType int) *RankingList {
	if !IsRankListTypeValid(RankListType(listType)) {
		return nil
	}

	// todo 排行榜是否开启

	if list, ok := (*l)[RankListType(listType)]; ok {
		return list
	}
	return nil
}
