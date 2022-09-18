package vcache

import (
	"github.com/livebud/bud/package/virtual"
)

var Discard Cache = discard{}

type discard struct{}

func (discard) Has(path string) (ok bool)                            { return false }
func (discard) Get(path string) (entry virtual.Entry, ok bool)       { return nil, false }
func (discard) Set(path string, entry virtual.Entry)                 {}
func (discard) Delete(path string)                                   {}
func (discard) Range(fn func(path string, entry virtual.Entry) bool) {}
func (discard) Clear()                                               {}
