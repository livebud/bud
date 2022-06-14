package embedded

import (
	_ "embed"
)

//go:embed favicon.ico
var favicon []byte

// Favicon returns the favicon data
func Favicon() []byte {
	return favicon
}
