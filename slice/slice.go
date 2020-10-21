package slice

import (
	"math/rand"
	"strings"
)

func Contains(sl []interface{}, v interface{}) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func ContainsInt(sl []int, v int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func ContainsInt64(sl []int64, v int64) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func ContainsString(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func ContainsStringCaseInsensitive(sl []string, v string) bool {
	for _, vv := range sl {
		if strings.ToUpper(vv) == strings.ToUpper(v) {
			return true
		}
	}
	return false
}

func CompareStringCaseInsensitive(sl string, v string) bool {
	return strings.ToUpper(sl) == strings.ToUpper(v)
}

func CompareStrings(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func RandomIndexes(lenght, randomCount int) []int {
	if lenght == 0 {
		return nil
	}

	if randomCount > lenght {
		randomCount = lenght
	}

	// rand.Seed(time.Now().Unix())

	values := make([]int, lenght)
	for i := 0; i < lenght; i++ {
		values[i] = i
	}

	ret := []int{}

	for len(ret) < randomCount {
		idx := rand.Intn(len(values))
		ret = append(ret, values[idx])

		values = append(values[:idx], values[idx+1:]...)
	}

	return ret
}

func RandomIndexByString(lenght, count int, str string) []int {
	if lenght == 0 {
		return nil
	}

	if count > lenght {
		count = lenght
	}

	bs := []byte(str)
	lastB := 0
	if len(bs) > 0 {
		lastB = int(bs[len(bs)-1])
	}

	idx := lastB % lenght
	ret := []int{idx}
	for i := 0; i < count-1; i++ {
		idx = idx + 1
		if idx >= lenght {
			idx = 0
		}

		ret = append(ret, idx)
	}

	return ret
}
