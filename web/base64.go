package web

import "encoding/base64"

func Base64Encode(origin string) string {
	return base64.URLEncoding.EncodeToString([]byte(origin))
}

func Base64Decode(code string) string {
	result, err := base64.URLEncoding.DecodeString(code)
	if err != nil {
		return ""
	}
	return string(result)
}
