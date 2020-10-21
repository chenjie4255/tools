package toolkit

func CutWords(str string, count int) string {
	if count <= 0 {
		return ""
	}
	runeStr := []rune(str)

	if len(runeStr) <= count {
		return str
	}

	return string(runeStr[:count])
}
