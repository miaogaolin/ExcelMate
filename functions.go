package main

import (
	"html/template"

	"github.com/Masterminds/sprig/v3"
)

var functions = map[string]interface{}{
	"substr": substring,
}

func FuncMap() template.FuncMap {
	maps := sprig.FuncMap()
	for k, v := range functions {
		maps[k] = v
	}
	return maps
}

// substring creates a substring of the given string.
//
// If start is < 0, this calls string[:end].
//
// If start is >= 0 and end < 0 or end bigger than s length, this calls string[start:]
//
// Otherwise, this calls string[start, end].
func substring(start, end int, s string) string {

	r := []rune(s)
	if start < 0 {
		return string(r[:end])
	}
	if end < 0 || end > len(r) {
		return string(r[start:])
	}
	return string(r[start:end])
}
