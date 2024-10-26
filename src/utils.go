package main

import (
	"fmt"
	"github.com/JexSrs/go-elsewherr/src/environment"
	"regexp"
	"strings"
)

func debug(format string, args ...any) {
	if environment.Env.Debug {
		fmt.Printf(format+"\n", args...)
	}
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

func cleanString(str string) string {
	re := regexp.MustCompile(`[^A-Za-z0-9]+`)
	return strings.ToLower(re.ReplaceAllString(str, ""))
}
