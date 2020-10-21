package algorithm

import (
	"testing"
	"time"
)

func TestCalMaxContinouiousDays(t *testing.T) {
	type args struct {
		ts  []int64
		loc *time.Location
	}

	t0 := time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC).Unix()
	var day int64 = 24 * 3600
	var hour int64 = 3600

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"3 days",
			args{ts: []int64{t0, t0 + 24*3600, t0 + 48*3600},
				loc: time.UTC},
			3,
		},
		{
			"with other days",
			args{ts: []int64{t0 - 37*hour, t0, t0 + 1*day, t0 + 2*day, t0 + 4*day},
				loc: time.UTC},
			3,
		},
		{
			"with 2 period",
			args{ts: []int64{
				t0,
				t0 + 1*day,
				t0 + 2*day,
				t0 + 5*day,
				t0 + 6*day,
				t0 + 7*day,
				t0 + 8*day,
				t0 + 10*day},
				loc: time.UTC},
			4,
		},
		{
			"with 2 period (test 2)",
			args{ts: []int64{
				t0,
				t0 + 1*day,
				t0 + 2*day,
				t0 + 5*day,
				t0 + 6*day,
				t0 + 8*day,
				t0 + 8*day,
				t0 + 10*day},
				loc: time.UTC},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalMaxContinuousDays(tt.args.ts, tt.args.loc); got != tt.want {
				t.Errorf("CalMaxContinuousDays() = %v, want %v", got, tt.want)
			}
		})
	}
}
