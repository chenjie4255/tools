package rand

import "testing"

func TestRandomStr(t *testing.T) {
	str, _ := Str(12)
	if len(str) != 12 {
		t.Fatalf("unexcepted length:%d", len(str))
	}
}

func TestDigital(t *testing.T) {
	str, _ := Digital(4)
	if len(str) != 4 {
		t.Fatalf("unexcepted length:%d", len(str))
	}
	t.Log(str)
}
