package mod

import (
	"errors"
	"fmt"
	"go/build"
	"path/filepath"
	"strings"
)

// gopathToModulePath tries inferring the module path of directory. This only works if you're
func modulePathFromGoPath(path string) string {
	src := filepath.Join(build.Default.GOPATH, "src") + "/"
	if !strings.HasPrefix(path, src) {
		return ""
	}
	modulePath := strings.TrimPrefix(path, src)
	return modulePath
}

var ErrCantInfer = errors.New("mod: unable to infer the module path")

// Infer the module path. This only works if you work inside $GOPATH.
// Expects dir to be an absolute path
func Infer(dir string) (modulePath string, err error) {
	modulePath = modulePathFromGoPath(dir)
	if modulePath == "" {
		return "", fmt.Errorf("%w for %q, run `go mod init` to fix", ErrCantInfer, dir)
	}
	return modulePath, nil
}
