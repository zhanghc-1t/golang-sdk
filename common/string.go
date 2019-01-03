package common

func Substr(s string, start int, length int) string {
	if len(s) == 0 {
		return ""
	}
	rs := []rune(s)
	rl := len(s)

	end := 0

	if start < 0 {
		start = rl + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}

	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
