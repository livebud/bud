// Copyright (C) 2017 ScyllaDB
// Use of this source code is governed by a ALv2-style
// license that can be found at https://github.com/scylladb/go-set/LICENSE.

package set

import (
	"fmt"
	"io/fs"
	"strings"
)

var (
	// helpful to not write everywhere struct{}{}
	keyExists   = struct{}{}
	nonExistent fs.DirEntry
)

const maxInt = int(^uint(0) >> 1)

// Set is the main set structure that holds all the data
// and methods used to working with the set.
type Set struct {
	m map[string]fs.DirEntry
}

// New creates and initializes a new Set.
func New(ts ...fs.DirEntry) *Set {
	s := NewWithSize(len(ts))
	s.Add(ts...)
	return s
}

// NewWithSize creates a new Set and gives make map a size hint.
func NewWithSize(size int) *Set {
	return &Set{make(map[string]fs.DirEntry, size)}
}

// Add includes the specified items (one or more) to the Set. The underlying
// Set s is modified. If passed nothing it silently returns.
func (s *Set) Add(items ...fs.DirEntry) {
	for _, item := range items {
		s.m[item.Name()] = item
	}
}

// Remove deletes the specified items from the Set. The underlying Set s is
// modified. If passed nothing it silently returns.
func (s *Set) Remove(items ...fs.DirEntry) {
	for _, item := range items {
		delete(s.m, item.Name())
	}
}

// Pop deletes and returns an item from the Set. The underlying Set s is
// modified. If Set is empty, the zero value is returned.
func (s *Set) Pop() fs.DirEntry {
	for name, item := range s.m {
		delete(s.m, name)
		return item
	}
	return nil
}

// Pop2 tries to delete and return an item from the Set. The underlying Set s
// is modified. The second value is a bool that is true if the item existed in
// the set, and false if not. If Set is empty, the zero value and false are
// returned.
func (s *Set) Pop2() (fs.DirEntry, bool) {
	for name, item := range s.m {
		delete(s.m, name)
		return item, true
	}
	return nonExistent, false
}

// Has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of  the items exist.
func (s *Set) Has(items ...fs.DirEntry) bool {
	has := false
	for _, item := range items {
		if _, has = s.m[item.Name()]; !has {
			break
		}
	}
	return has
}

// HasAny looks for the existence of any of the items passed.
// It returns false if nothing is passed.
// For multiple items it returns true if any of the items exist.
func (s *Set) HasAny(items ...fs.DirEntry) bool {
	has := false
	for _, item := range items {
		if _, has = s.m[item.Name()]; has {
			break
		}
	}
	return has
}

// Size returns the number of items in a Set.
func (s *Set) Size() int {
	return len(s.m)
}

// Clear removes all items from the Set.
func (s *Set) Clear() {
	s.m = make(map[string]fs.DirEntry)
}

// IsEmpty reports whether the Set is empty.
func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

// IsEqual test whether s and t are the same in size and have the same items.
func (s *Set) IsEqual(t *Set) bool {
	// return false if they are no the same size
	if s.Size() != t.Size() {
		return false
	}

	equal := true
	t.Each(func(item fs.DirEntry) bool {
		_, equal = s.m[item.Name()]
		return equal // if false, Each() will end
	})

	return equal
}

// IsSubset tests whether t is a subset of s.
func (s *Set) IsSubset(t *Set) bool {
	if s.Size() < t.Size() {
		return false
	}

	subset := true

	t.Each(func(item fs.DirEntry) bool {
		_, subset = s.m[item.Name()]
		return subset
	})

	return subset
}

// IsSuperset tests whether t is a superset of s.
func (s *Set) IsSuperset(t *Set) bool {
	return t.IsSubset(s)
}

// Each traverses the items in the Set, calling the provided function for each
// Set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s *Set) Each(f func(item fs.DirEntry) bool) {
	for _, item := range s.m {
		if !f(item) {
			break
		}
	}
}

// Copy returns a new Set with a copy of s.
func (s *Set) Copy() *Set {
	u := NewWithSize(s.Size())
	for name, item := range s.m {
		u.m[name] = item
	}
	return u
}

// String returns a string representation of s
func (s *Set) String() string {
	v := make([]string, 0, s.Size())
	for item := range s.m {
		v = append(v, fmt.Sprintf("%v", item))
	}
	return fmt.Sprintf("[%s]", strings.Join(v, ", "))
}

// List returns a slice of all items. There is also StringSlice() and
// IntSlice() methods for returning slices of type string or int.
func (s *Set) List() []fs.DirEntry {
	v := make([]fs.DirEntry, 0, s.Size())
	for _, item := range s.m {
		v = append(v, item)
	}
	return v
}

// Merge is like Union, however it modifies the current Set it's applied on
// with the given t Set.
func (s *Set) Merge(t *Set) {
	for name, item := range t.m {
		s.m[name] = item
	}
}

// Separate removes the Set items containing in t from Set s. Please aware that
// it's not the opposite of Merge.
func (s *Set) Separate(t *Set) {
	for item := range t.m {
		delete(s.m, item)
	}
}

// Union is the merger of multiple sets. It returns a new set with all the
// elements present in all the sets that are passed.
func Union(sets ...*Set) *Set {
	maxPos := -1
	maxSize := 0
	for i, set := range sets {
		if l := set.Size(); l > maxSize {
			maxSize = l
			maxPos = i
		}
	}
	if maxSize == 0 {
		return New()
	}

	u := sets[maxPos].Copy()
	for i, set := range sets {
		if i == maxPos {
			continue
		}
		for name, item := range set.m {
			u.m[name] = item
		}
	}
	return u
}

// Difference returns a new set which contains items which are in in the first
// set but not in the others.
func Difference(set1 *Set, sets ...*Set) *Set {
	s := set1.Copy()
	for _, set := range sets {
		s.Separate(set)
	}
	return s
}

// Intersection returns a new set which contains items that only exist in all
// given sets.
func Intersection(sets ...*Set) *Set {
	minPos := -1
	minSize := maxInt
	for i, set := range sets {
		if l := set.Size(); l < minSize {
			minSize = l
			minPos = i
		}
	}
	if minSize == maxInt || minSize == 0 {
		return New()
	}

	t := sets[minPos].Copy()
	for i, set := range sets {
		if i == minPos {
			continue
		}
		for item := range t.m {
			if _, has := set.m[item]; !has {
				delete(t.m, item)
			}
		}
	}
	return t
}

// SymmetricDifference returns a new set which s is the difference of items
// which are in one of either, but not in both.
func SymmetricDifference(s *Set, t *Set) *Set {
	u := Difference(s, t)
	v := Difference(t, s)
	return Union(u, v)
}
