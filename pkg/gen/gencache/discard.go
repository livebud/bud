package gencache

import (
	"errors"

	"github.com/livebud/bud/pkg/gen"
	"github.com/livebud/bud/pkg/virt"
)

func Discard() gen.Cache {
	return discard{}
}

type discard struct{}

var _ gen.Cache = (*discard)(nil)

func (discard) Get(path string) (*virt.File, error) {
	return nil, errors.New("not found")
}

func (discard) Set(path string, file *virt.File) error {
	return nil
}

func (discard) Link(from string, toPatterns ...string) error {
	return nil
}

func (discard) Delete(paths ...string) error {
	return nil
}

func (discard) Reset() error {
	return nil
}

func (discard) Close() error {
	return nil
}