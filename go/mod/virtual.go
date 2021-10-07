package mod

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/go-duo/bud/go/is"
)

// Virtual module file
func Virtual(modulePath, dir string) File {
	return &virtual{modulePath, dir}
}

type virtual struct {
	modulePath string
	dir        string
}

func (v *virtual) Directory() string {
	return v.dir
}

func (v *virtual) ModulePath() string {
	return v.modulePath
}

func (v *virtual) ResolveImport(dir string) (importPath string, err error) {
	return resolveImport(v, dir)
}

func (v *virtual) ResolveDirectory(importPath string) (dir string, err error) {
	if is.StdLib(importPath) {
		return filepath.Join(stdDir, importPath), nil
	}
	dir = filepath.Join(build.Default.GOPATH, "src", importPath)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("%q doesn't exist. Unable to resolve import path %q", dir, importPath)
		}
		return "", err
	}
	return dir, nil
}
