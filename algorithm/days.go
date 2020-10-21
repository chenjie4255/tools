package algorithm

import (
	"sort"
	"time"
)

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func CalMaxContinuousDays(ts []int64, loc *time.Location) int {
	if len(ts) == 0 {
		return 0
	}

	sortTs := make([]int64, len(ts))
	copy(sortTs, ts)
	sort.Sort(Int64Slice(sortTs))

	max := 1
	current := 1
	lastTs := sortTs[0]

	for i := 1; i < len(sortTs); i++ {
		if isContinuousDays(lastTs, sortTs[i], loc) {
			current++
			if current > max {
				max = current
			}
		} else if !isSameDay(lastTs, sortTs[i], loc) {
			// reset current if not a same day
			current = 1
		}
		lastTs = sortTs[i]
	}

	return max
}

func isSameDay(t1 int64, t2 int64, loc *time.Location) bool {
	t1Time := time.Unix(t1, 0).In(loc)
	t2Time := time.Unix(t2, 0).In(loc)

	return t1Time.Year() == t2Time.Year() && t1Time.YearDay() == t2Time.YearDay()
}

func isContinuousDays(t1 int64, t2 int64, loc *time.Location) bool {
	if t1 > t2 {
		t1, t2 = t2, t1
	}
	// make sure t2 > t1
	if t2-t1 > 3600*24*2 {
		return false
	}

	t1Time := time.Unix(t1, 0).In(loc)
	t2Time := time.Unix(t2, 0).In(loc)

	t1Day := time.Date(t1Time.Year(), t1Time.Month(), t1Time.Day(), 0, 0, 0, 0, loc).Unix()
	t2Day := time.Date(t2Time.Year(), t2Time.Month(), t2Time.Day(), 0, 0, 0, 0, loc).Unix()

	return t2Day-t1Day == 3600*24
}
