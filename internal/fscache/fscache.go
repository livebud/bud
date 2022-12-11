package fscache

import (
	"github.com/livebud/bud/package/virtual"
)

type Cache interface {
	Get(name string) (entry virtual.Entry, ok bool)
	Set(path string, entry virtual.Entry)
	Link(from, to string)
	Check(from string, checker func(path string) (changed bool))
}
