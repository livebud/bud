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

//go:embed gitignore.txt
var gitignore []byte

// Gitignore returns the gitignore data
func Gitignore() []byte {
	return gitignore
}
