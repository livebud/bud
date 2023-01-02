package dag

import "github.com/livebud/bud/package/virt"

type Cache interface {
	Get(path string) (*virt.File, error)
	Set(path string, file *virt.File) error
	Link(from string, toPatterns ...string) error
	Delete(paths ...string) error
	Reset() error
	Close() error
}
