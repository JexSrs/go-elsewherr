package utils

import (
	"regexp"
	"strings"
)

func Get[T any](m map[string]interface{}, key string) T {
	dt, exists := m[key]
	if !exists || dt == nil {
		var zero T
		return zero
	}
	return dt.(T)
}

func ToIntArray(arr []interface{}) []int {
	intArray := make([]int, len(arr))
	for i, v := range arr {
		intArray[i] = int(v.(float64))
	}

	return intArray
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func FindIndex[T any](ss []T, test func(T) bool) int {
	for i, s := range ss {
		if test(s) {
			return i
		}
	}
	return -1
}

func CleanString(str string) string {
	reSpaces := regexp.MustCompile(`\s+`)
	str = reSpaces.ReplaceAllString(str, " ")
	str = strings.ReplaceAll(str, " ", "-")
	re := regexp.MustCompile(`[^A-Za-z0-9-]+`)
	return strings.ToLower(re.ReplaceAllString(str, ""))
}
