package helper

import (
	"math"
)

func SafeAdd(a, b int32) int32 {
	if uint64(a)+uint64(b) > uint64(math.MaxInt32) {
		return math.MaxInt32
	}
	return a + b
}

func SafeSub(a, b int32) int32 {
	if a >= b {
		return a - b
	}
	return 0
}
