package common

import (
	"reflect"
)

func in(needle interface{}, haystack interface{}) int32 {
	nType := reflect.TypeOf(needle)
	hType := reflect.TypeOf(haystack)
	if hType.Kind() != reflect.Slice ||
		hType.Elem() != nType {
		return -1
	}

	hValue := reflect.ValueOf(haystack)
	hValueLength := hValue.Len()
	for i := 0; i < hValueLength; i++ {
		if reflect.DeepEqual(hValue.Index(i).Interface(), needle) {
			return int32(i)
		}
	}

	return -1
}