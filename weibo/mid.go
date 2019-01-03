package weibo

import (
	"fmt"
	"math"
	"strconv"
)

const (
	encodeCharGroupLength = 7
	decodeCharGroupLength = 4
)

func MidDecode(s string) (int, error) {
	rs := []rune(s)
	l := len(rs)
	t := int(math.Ceil(float64(l) / decodeCharGroupLength))
	r := ""

	for i := 1; i <= t; i++ {
		s := l - i*decodeCharGroupLength
		e := s + decodeCharGroupLength
		if s < 0 {
			s = 0
		}
		if e > l {
			e = l
		}

		m := string(rs[s:e])
		r = fmt.Sprintf("%0*s", 7, strconv.Itoa(Base62Decode(m))) + r
	}

	return strconv.Atoi(r)
}

func MidEncode(mid int) (string, error) {
	rMid := []rune(strconv.Itoa(mid))
	l := len(rMid)
	t := int(math.Ceil(float64(l) / encodeCharGroupLength))
	r := ""

	for i := 1; i <= t; i++ {
		s := l - i*encodeCharGroupLength
		e := s + encodeCharGroupLength
		if s < 0 {
			s = 0
		}
		if e > l {
			e = l
		}

		m := string(rMid[s:e])
		mD, err := strconv.Atoi(m)
		if err != nil {
			return "", err
		}

		ir := fmt.Sprintf("%0*s", 4, Base62Encode(mD))
		if i == t {
			ir = Base62Encode(mD)
		}
		r = ir + r
	}

	return r, nil
}
