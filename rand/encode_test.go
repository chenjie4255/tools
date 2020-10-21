package rand

import "testing"

func TestEncode2NormalText(t *testing.T) {
	data, _:= Bytes(16)

	str := Encode2NormalText(data)
	t.Logf("result:%s", str)


}
