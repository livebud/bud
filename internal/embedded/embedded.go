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

// Default CSS is modern-normalize by Sindre Sorhus
// https://raw.githubusercontent.com/sindresorhus/modern-normalize/v1.1.0/modern-normalize.css
//go:embed layout.css
var layout []byte

// Layout returns the default stylesheet
func Layout() []byte {
	return layout
}
