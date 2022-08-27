package embedded

import (
	_ "embed"
)

//go:embed favicon.ico
var favicon []byte

//go:embed normalize.css
var normalize string

// Favicon returns the favicon data
func Favicon() []byte {
	return favicon
}

// DefaultCss returns the default css data
func DefaultCss(css string) string {
	switch css {
	case "normalize":
		return normalize
	}
	return "/* No Default CSS Loaded */"
}
