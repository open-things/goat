package goat

import (
	"path/filepath"
	"strings"
)

func ReplaceAll(s string, old, new []string) string {
	var r string = s
	for o, n := 0, 0; o < len(old); o++ {
		n = o
		if n >= len(new) {
			n = len(new) - 1
		}
		r = strings.Replace(r, old[o], new[n], -1)
	}
	return r
}

var TemplatePatterns []string = []string{
	filepath.Join("{type}", "{group}", "{name}.{extension}"),
	filepath.Join("{type}", "{group}", "{name}_*.{extension}"),
}

type TemplateFieldAccessors interface {
	Group() string
	Name() string
	Type() string
	Extension() string
}

func GetPatterns(template TemplateFieldAccessors) []string {
	var result []string = make([]string, len(TemplatePatterns))
	var old []string = []string{"{group}", "{name}", "{type}", "{extension}"}
	var new []string = []string{
		template.Group(),
		template.Name(),
		template.Type(),
		template.Extension(),
	}
	for idx, pat := range TemplatePatterns {
		result[idx] = ReplaceAll(pat, old, new)
	}
	return result
}
