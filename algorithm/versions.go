package algorithm

import (
	"math"
	"strconv"
	"strings"
)

// CompareVersion 比较2个版本的大小，-1表示小于，1表示大于，0表示相等
func CompareVersion(v1, v2 string, bit int) int {
	if bit <= 0 || bit > 4 {
		panic("bit should be in range [1,4]")
	}

	v1Slice := strings.Split(v1, ".")
	v2Slice := strings.Split(v2, ".")

	var v1Val int64
	var v2Val int64

	for i := range v1Slice {
		val, _ := strconv.ParseInt(v1Slice[i], 10, 64)
		if val > 999 {
			val = 999
		}
		if i >= bit {
			val = 0
		}
		v1Val += val * int64(math.Pow10(3*(3-i)))
	}

	for i := range v2Slice {
		val, _ := strconv.ParseInt(v2Slice[i], 10, 64)
		if val > 999 {
			val = 999
		}
		if i >= bit {
			val = 0
		}
		v2Val += val * int64(math.Pow10(3*(3-i)))
	}

	if v1Val == v2Val {
		return 0
	}

	if v1Val > v2Val {
		return 1
	}

	return -1

}
