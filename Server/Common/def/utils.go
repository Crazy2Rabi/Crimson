package def

// 注意依赖顺序

func MakeTwoKey(val1, val2 int32) int64 {
	return int64(val1)<<32 | int64(val2)
}
