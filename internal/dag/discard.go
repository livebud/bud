package dag

import (
	"errors"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/virt"
)

var Discard = discard{}

type discard struct{}

var _ genfs.Cache = (*discard)(nil)

func (discard) Get(path string) (*virt.File, error) {
	return nil, errors.New("not found")
}
func (discard) Set(path string, file *virt.File) error {
	return nil
}
func (discard) Link(from string, toPatterns ...string) error {
	return nil
}

func (discard) Reset() error {
	return nil
}
