package rand

const maskTable = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

const maskLen = len(maskTable)

func Encode2NormalText(src []byte) string {
	dst := make([]byte, len(src))
	for i := range src {
		dst[i] = maskTable[int(src[i]) % maskLen]
	}

	return string(dst)
}
