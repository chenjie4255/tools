package web

import "testing"

func Test_base64Encode(t *testing.T) {
	testStr := []string{
		"1234",
		"gg",
		"",
		"fajsoifdjoasidfjoasijdfioasjfa98341937492y34~!@@$#%^$&*(*%$^&%#@!@~",
	}
	for _, tt := range testStr {
		encoding := Base64Encode(tt)
		decoding := Base64Decode(encoding)
		if decoding != tt {
			t.Fatal("gg")
		}
	}
}
