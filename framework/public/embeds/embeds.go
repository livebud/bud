package embeds

import (
	_ "embed"
)

//go:embed favicon.ico
var favicon []byte

// Favicon returns the favicon data
func Favicon() []byte {
	return favicon
}

// Default CSS is modern-normalize by Sindre Sorhus
// https://raw.githubusercontent.com/sindresorhus/modern-normalize/v1.1.0/modern-normalize.css
//go:embed default.css
var stylesheet []byte

// Stylesheet returns the default stylesheet
func Stylesheet() []byte {
	return stylesheet
}
