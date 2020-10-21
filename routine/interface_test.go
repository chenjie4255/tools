package routine

import "testing"

func TestSchedulerPos_offset(t *testing.T) {
	tests := []struct {
		name string
		p    SchedulerPos
		want uint32
	}{
		{"t1",
			SchedulerPos(6717496672496124032),
			371874,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.offset(); got != tt.want {
				t.Errorf("SchedulerPos.offset() = %v, want %v", got, tt.want)
			}
		})
	}
}
