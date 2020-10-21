package filters

import (
	"strconv"
	"strings"
)

func isReverse(filter string) bool {
	if len(filter) < 1 {
		return false
	}
	return filter[:1] == "!"
}

func isAllReverseFilters(filters []string) bool {
	for _, filter := range filters {
		if !isReverse(filter) {
			return false
		}
	}

	return true
}

// FilterString 过滤简单的String类型，如!CN, !US等
func FilterString(filters []string, value string) bool {
	if len(filters) == 0 || value == "" {
		return true
	}

	if isAllReverseFilters(filters) {
		// 如果全部都是返向过滤器，则
		for _, filter := range filters {
			if strings.ToUpper(filter[1:]) == strings.ToUpper(value) {
				return false
			}
		}

		return true
	}

	for _, filter := range filters {
		if isReverse(filter) {
			if strings.ToUpper(filter[1:]) != strings.ToUpper(value) {
				return true
			}
		} else {
			if strings.ToUpper(value) == strings.ToUpper(filter) {
				return true
			}
		}
	}

	return false
}

func parseVersion2Int(version string) (int64, error) {
	versions := strings.Split(version, ".")
	var ret int64

	// 如果没有逗号分隔，则直接比较
	if len(versions) == 1 {
		return strconv.ParseInt(version, 10, 64)
	}

	if len(versions) == 2 {
		versions = append(versions, "0")
	}

	for _, v := range versions {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}

		ret = ret*100 + val
	}

	return ret, nil
}

// CheckVersionInRanges 检查是否在指定的范围内
func CheckVersionInRanges(ranges []string, version string) bool {
	if len(ranges) == 0 {
		return true
	}
	for _, rg := range ranges {
		if checkVersionInRange(rg, version) {
			return true
		}
	}

	return false
}

func checkVersionInRange(rg string, version string) bool {
	if version == rg {
		return true
	}

	dbVal, err := parseVersion2Int(version)
	if err != nil {
		return false
	}

	if rg[len(rg)-1:] == "-" {
		end, err := parseVersion2Int(rg[:len(rg)-1])
		if err != nil {
			return false
		}

		return dbVal <= end
	}

	if rg[len(rg)-1:] == "+" {
		start, err := parseVersion2Int(rg[:len(rg)-1])
		if err != nil {
			return false
		}

		return dbVal >= start
	}

	rgs := strings.Split(rg, "-")
	if len(rgs) == 2 {
		s, _ := parseVersion2Int(rgs[0])
		e, _ := parseVersion2Int(rgs[1])

		return dbVal >= s && dbVal <= e
	}

	return false
}
