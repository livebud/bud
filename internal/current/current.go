package current

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// Filename gets the current filename of the caller
func Filename() (string, error) {
	return filename(2)
}

func filename(skip int) (string, error) {
	_, filename, _, ok := runtime.Caller(skip)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

// Directory gets the current directory of the caller
func Directory() (string, error) {
	name, err := filename(2)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(name)
	// When we use --trimpath, attempt to find the absolute path anyway
	if !filepath.IsAbs(dir) {
		// Hail mary attempt to find it within $GOPATH/src
		if gopath := os.Getenv("GOPATH"); gopath != "" {
			dir = filepath.Join(gopath, "src", dir)
		}
	}
	return dir, nil
}
