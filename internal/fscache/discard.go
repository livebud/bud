package fscache

import (
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/virtual"
)

var Discard genfs.Cache = discard{}

type discard struct{}

func (discard) Get(name string) (entry virtual.Entry, ok bool) {
	return
}

func (discard) Set(path string, entry virtual.Entry) {
}

func (discard) Link(from, to string) {
}

func (discard) Check(from string, checker func(path string) (changed bool)) {
}
