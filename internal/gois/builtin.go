package gois

import "strings"

// Builtin checks if the dataType is built-in.
// TODO handle more complex types e.g. map[string]LocalControllerStruct
func Builtin(dataType string) bool {
	dataType = strings.TrimLeft(dataType, "[]*")
	if _, ok := builtin[dataType]; ok {
		return true
	}
	return false
}

// builtin types
var builtin = map[string]struct{}{
	"string":     {},
	"bool":       {},
	"error":      {},
	"int8":       {},
	"uint8":      {},
	"byte":       {},
	"int16":      {},
	"uint16":     {},
	"int32":      {},
	"rune":       {},
	"uint32":     {},
	"int64":      {},
	"uint64":     {},
	"int":        {},
	"uint":       {},
	"uintptr":    {},
	"float32":    {},
	"float64":    {},
	"complex64":  {},
	"complex128": {},
}
