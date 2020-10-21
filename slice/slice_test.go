package slice

import (
	"github.com/ahmetb/go-linq"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestRandomIndexes(t *testing.T) {
	type args struct {
		lenght      int
		randomCount int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"t1",
			args{10, 4},
		},
		{
			"t2",
			args{10, 4},
		},
		{
			"t3",
			args{10, 8},
		},
		{
			"t4",
			args{10, 9},
		},
		{
			"t5",
			args{10, 0},
		},
		{
			"t5",
			args{10, 10},
		},
	}
	rand.Seed(time.Now().Unix())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomIndexes(tt.args.lenght, tt.args.randomCount)
			t.Logf("result:%v", got)
			if len(got) != tt.args.randomCount {
				t.Errorf("RandomIndexes() = %v, want %v", got, tt.args.randomCount)
			}

			if linq.From(got).Distinct().Count() != len(got) {
				t.Errorf("index duplicate: %v", got)
			}
		})
	}
}

func TestRandomIndexByString(t *testing.T) {
	type args struct {
		lenght int
		count  int
		str    string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"t1",
			args{10, 1, ""},
			[]int{0},
		},
		{
			"t2",
			args{10, 10, ""},
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			"t3",
			args{10, 3, "1"},
			[]int{9, 0, 1},
		},
		{
			"t4",
			args{10, 3, "2"},
			[]int{0, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomIndexByString(tt.args.lenght, tt.args.count, tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RandomIndexByString() = %v, want %v", got, tt.want)
			}
		})
	}
}
