package versions

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

const minGoVersion = "v1.17"

// ErrMinGoVersion error is returned when Bud needs a newer version of Go
var ErrMinGoVersion = fmt.Errorf("bud requires Go %s or later", minGoVersion)

// CheckGo checks if the current version of Go is greater than the
// minimum required Go version.
func CheckGo(currentVersion string) error {
	currentVersion = "v" + strings.TrimPrefix(currentVersion, "go")
	// If we encounter an invalid version, it's probably a development version of
	// Go. We'll let those pass through. Reference:
	// https://github.com/golang/go/blob/3cf79d96105d890d7097d274804644b2a2093df1/src/runtime/extern.go#L273-L275
	if !semver.IsValid(currentVersion) {
		return nil
	}
	if semver.Compare(currentVersion, minGoVersion) < 0 {
		return ErrMinGoVersion
	}
	return nil
}
