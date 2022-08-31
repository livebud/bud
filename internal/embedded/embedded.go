package embedded

import (
	_ "embed"
)

//go:embed favicon.ico
var favicon []byte

//go:embed normalize.css
var normalize []byte

// Favicon returns the favicon data
func Favicon() []byte {
	return favicon
}

// NormalizeCss returns the normalize css data
func NormalizeCss() []byte {
	return normalize
}

// EmptyCss reset the default css data
func EmptyCss() []byte {
	return []byte("/* No Default CSS Loaded */")
}
