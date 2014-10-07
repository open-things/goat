package goat

import (
	"strings"
)

func ReplaceAll(s string, old, new []string) string {
	r := s
	for o, n := 0, 0; o < len(old); o++ {
		n = o
		if n >= len(new) {
			n = len(new) - 1
		}
		r = strings.Replace(r, old[o], new[n], -1)
	}
	return r
}

func CopyMap(source map[string]interface{}) map[string]interface{} {
	destination := make(map[string]interface{})
	if source != nil {
		for key, val := range source {
			destination[key] = val
		}
	}
	return destination
}
