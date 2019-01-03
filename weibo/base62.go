package weibo

import "math"

var bc = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func Base62Encode(n int) string {
	r := ""
	for {
		if n <= 0 {
			break
		}
		i := n % 62
		r = bc[i] + r
		n = n / 62
	}

	return r
}

func Base62Decode(s string) int {
	r := 0
	len := len(s)
	f := Base62Flip(bc)
	for i := 0; i < len; i++ {
		r += f[string(s[i])] * int(math.Pow(62, float64(len-i-1)))
	}
	return r
}

func Base62Flip(s []string) map[string]int {
	f := make(map[string]int)
	for idx, v := range s {
		f[v] = idx
	}
	return f
}
