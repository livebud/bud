package current

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// Filename gets the current filename of the caller
func Filename() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

// Directory gets the current directory of the caller
func Directory() (string, error) {
	filename, err := Filename()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(filename)
	// When we use --trimpath, attempt to find the absolute path anyway
	if !filepath.IsAbs(dir) {
		// Attempt to find it within $GOPATH/src
		if gopath := os.Getenv("GOPATH"); gopath != "" {
			dir = filepath.Join(gopath, "src", dir)
		}
	}
	return dir, nil
}
