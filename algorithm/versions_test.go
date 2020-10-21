package algorithm

import "testing"

func TestCompareVersion(t *testing.T) {
	type args struct {
		v1  string
		v2  string
		bit int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"equal",
			args{"1.0.0", "1.0.0", 3},
			0,
		},
		{
			"equal_different bit",
			args{"1.0.0.0", "1.0.0", 3},
			0,
		},
		{
			"equal_different bit2",
			args{"1.0.0.0", "1.0.0", 4},
			0,
		},
		{
			"equal_3",
			args{"1.0.0.0", "1.0.0.4", 3},
			0,
		},
		{
			"equal_different badformat",
			args{"1.", "1.0.0", 4},
			0,
		},
		{
			"greater",
			args{"1.1.0", "1.0.999", 3},
			1,
		},
		{
			"greater bit diff",
			args{"1.1.0.1", "1.0.999", 3},
			1,
		},
		{
			"greater bit diff2",
			args{"1.1.0", "1.0.999", 2},
			1,
		},
		{
			"greater diff at 3",
			args{"1.0.999", "1.0.998", 3},
			1,
		},
		{
			"greater_1",
			args{"2", "1", 3},
			1,
		},
		{
			"greater_2",
			args{"2.2.1.3", "2.2.0.9", 3},
			1,
		},
		{
			"greater_3",
			args{"2.2.2.4", "2.2.2.3", 4},
			1,
		},
		//
		{
			"less_1",
			args{"2.2.2.4", "2.2.2.9", 4},
			-1,
		},
		{
			"less_2",
			args{"2.2.1", "2.2.2.9", 3},
			-1,
		},
		{
			"less_3",
			args{"1.999.1", "2.0.2.9", 3},
			-1,
		},
		{
			"less_4",
			args{"3.5.9.1", "3.6.0", 4},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareVersion(tt.args.v1, tt.args.v2, tt.args.bit); got != tt.want {
				t.Errorf("CompareVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
