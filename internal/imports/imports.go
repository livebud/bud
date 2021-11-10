// Package imports is a small package for dealing with import paths. imports
// picks unique names if needed and is able to determine the assumed name of the
// of an import path. It also orders the imports properly for `gofmt`.
package imports

import (
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// New import set
func New() *Set {
	return &Set{
		names:    map[string]int{},
		paths:    map[string]string{},
		reserved: map[string]string{},
	}
}

// Set of imports
type Set struct {
	names    map[string]int
	paths    map[string]string
	reserved map[string]string
}

// AddStd is a convenience function for adding standard library packages
func (s *Set) AddStd(pkgs ...string) {
	for _, pkg := range pkgs {
		// Split to support subpackages like net/http
		_, name := path.Split(pkg)
		s.AddNamed(name, pkg)
	}
}

// Add to the set
func (s *Set) Add(path string) string {
	if name, ok := s.paths[path]; ok {
		return name
	}
	if reserved, ok := s.reserved[path]; ok {
		delete(s.reserved, path)
		s.paths[path] = reserved
		return reserved
	}
	name := AssumedName(path)
	ith := s.names[name]
	uniqueName := name
	if ith > 0 {
		uniqueName += strconv.Itoa(ith)
	}
	s.paths[path] = uniqueName
	s.names[name]++
	return uniqueName
}

// AddNamed brings a preferred name
func (s *Set) AddNamed(name, path string) string {
	if name, ok := s.paths[path]; ok {
		return name
	}
	if reserved, ok := s.reserved[path]; ok {
		delete(s.reserved, path)
		s.paths[path] = reserved
		return reserved
	}
	ith := s.names[name]
	uniqueName := name
	if ith > 0 {
		uniqueName += strconv.Itoa(ith)
	}
	s.paths[path] = uniqueName
	s.names[name]++
	return uniqueName
}

// Reserve a name, but don't take import it. If Add or AddNamed come later,
// we'll use the same name but add it to the list of paths
func (s *Set) Reserve(path string) string {
	if name, ok := s.paths[path]; ok {
		return name
	}
	name := AssumedName(path)
	ith := s.names[name]
	uniqueName := name
	if ith > 0 {
		uniqueName += strconv.Itoa(ith)
	}
	s.reserved[path] = uniqueName
	s.names[name]++
	return uniqueName
}

// List imports by path first, then by name
func (s *Set) List() (imports []*Import) {
	imports = make([]*Import, len(s.paths))
	i := 0
	for path, name := range s.paths {
		imports[i] = &Import{name, path}
		i++
	}
	sort.Slice(imports, func(i int, j int) bool {
		if imports[i].Path == imports[j].Path {
			return imports[i].Name < imports[j].Name
		}
		return imports[i].Path < imports[j].Path
	})
	return imports
}

// Import struct returned by list
type Import struct {
	// Package identifier name
	Name string `json:"name,omitempty"`
	// Path to the import
	Path string `json:"path,omitempty"`
}

// AssumedName returns the assumed package name of an import path.
// It does this using only string parsing of the import path.
// It picks the last element of the path that does not look like a major
// version, and then picks the valid identifier off the start of that element.
// It is used to determine if a local rename should be added to an import for
// clarity.
// This function could be moved to a standard package and exported if we want
// for use in other tools.
func AssumedName(importPath string) string {
	base := path.Base(importPath)
	if strings.HasPrefix(base, "v") {
		if _, err := strconv.Atoi(base[1:]); err == nil {
			dir := path.Dir(importPath)
			if dir != "." {
				base = path.Base(dir)
			}
		}
	}
	base = strings.TrimPrefix(base, "go-")
	if i := strings.IndexFunc(base, notIdentifier); i >= 0 {
		base = base[:i]
	}
	return base
}

// notIdentifier reports whether ch is an invalid identifier character.
func notIdentifier(ch rune) bool {
	return !('a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
		'0' <= ch && ch <= '9' ||
		ch == '_' ||
		ch >= utf8.RuneSelf && (unicode.IsLetter(ch) || unicode.IsDigit(ch)))
}
