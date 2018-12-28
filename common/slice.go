package common

import (
	"fmt"
	"reflect"
)

func in(needle interface{}, haystack interface{}) int {
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
			return i
		}
	}

	return -1
}

func column(haystack interface{}, column string) (interface{}, error) {
	aType := reflect.TypeOf(haystack)
	aTypeKind := aType.Kind()
	if aTypeKind != reflect.Slice && aTypeKind != reflect.Map {
		return nil, fmt.Errorf("slice/map get column error:haystack is not a slice or map")
	}

	return nil, nil
}

func slice_column(haystack []interface{}, column string) (interface{}, error) {
	return nil, nil
}

func map_column(haystack []interface{}, column string) (interface{}, error) {
	return nil, nil
}

func value(object interface{}, key string) (interface{}, error) {
	return nil, nil
}

func map_value(object map[string]interface{}, key string) (interface{}, error) {
	return nil, nil
}

func struct_value(object interface{}, key string) (interface{}, error) {
	return nil, nil
}
