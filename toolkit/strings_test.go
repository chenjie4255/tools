package toolkit

import "testing"

func TestCutWords(t *testing.T) {
	type args struct {
		str   string
		count int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"t1",
			args{"123", 1},
			"1",
		},
		{
			"t2",
			args{"看心情123", 4},
			"看心情1",
		},
		{
			"t3",
			args{"看心情123", 41},
			"看心情123",
		},
		{
			"t4",
			args{"看心情123", 0},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CutWords(tt.args.str, tt.args.count); got != tt.want {
				t.Errorf("CutWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
