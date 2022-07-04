
package goversion

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// Represents the version of golang
type Version struct {
	Major int
	Minor int
	Patch int
}

// Returns a version from the given 3 integers in the format
// <major>.<minor>.<patch>.
func NewVersion(major, minor, patch int) Version {
	return Version{major, minor, patch}
}

// CurrentVersion builds and returns a Version object from the string
// returned by runtime.Version().
func CurrentVersion() Version {
	
	versionString := runtime.Version() // something like "go1.18.3"
	if !strings.HasPrefix(versionString, "go") {
		// version might be a hash or something, can't compare
		return Version{} // 0.0.0
	}
	// trim to semver version like "1.18.3" or possibly "1.18"
	semverString := strings.Replace(versionString, "go", "", 1)
	versionSlice := strings.Split(semverString, ".") // ["1", "18", "3"]
	
	version := Version{}
	
	if len(versionSlice) >= 1 {
		major, err := strconv.Atoi(versionSlice[0])
		if err != nil {
			return Version{}
		}
		version.Major = major // 1
	}
	if len(versionSlice) >= 2 {
		minor, err := strconv.Atoi(versionSlice[1])
		if err != nil {
			return Version{}
		}
		version.Minor = minor // 2
	}
	if len(versionSlice) >= 3 {
		patch, err := strconv.Atoi(versionSlice[2])
		if err != nil {
			return Version{}
		}
		version.Patch = patch
	}
	
	return version
}

// CompareTo returns a negative integer
// if the version is less than the other. It returns a positive integer
// if the version is greater than the other. It returns 0 if both versions
// are the same.
func (v Version) CompareTo(other Version) int {
	
	majorComparison := v.Major - other.Major
	if majorComparison != 0 {
		return majorComparison
	}
	
	minorComparison := v.Minor - other.Minor
	if minorComparison != 0 {
		return minorComparison
	}
	
	patchComparison := v.Patch - other.Patch
	if patchComparison != 0 {
		return patchComparison
	}
	
	return 0
}

// IsValid returns true if the version is > 0.0.0.
func (v Version) IsValid() bool {
	if v.Major <= 0 && v.Minor <= 0 && v.Patch <= 0 {
		return false
	}
	return true
}

// ToString returns a string like "<major>.<minor>.<patch>".
func (v Version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}