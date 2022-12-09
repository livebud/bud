package vcg

import (
	"github.com/livebud/bud/package/budfs/linkmap"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual"
	"github.com/livebud/bud/package/virtual/vcache"
)

type Cache interface {
	Get(name string) (entry virtual.Entry, ok bool)
	Set(path string, entry virtual.Entry)
	Link(from, to string)
	Check(from string, checker func(path string) (changed bool))
}

func New(log log.Log) *Memory {
	return &Memory{
		lm:  linkmap.New(log),
		log: log,
		vc:  vcache.New(),
	}
}

type Memory struct {
	lm  *linkmap.Map
	log log.Log
	vc  vcache.Cache
}

func (m *Memory) Get(name string) (entry virtual.Entry, ok bool) {
	return m.vc.Get(name)
}

func (m *Memory) Set(path string, entry virtual.Entry) {
	m.vc.Set(path, entry)
}

func (m *Memory) Delete(paths ...string) {
	for i := 0; i < len(paths); i++ {
		path := paths[i]
		if m.vc.Has(path) {
			m.log.Debug("budfs: delete cache path %q", path)
			m.vc.Delete(path)
		}
		m.lm.Range(func(genPath string, fns *linkmap.List) bool {
			if m.vc.Has(genPath) && fns.Check(path) {
				paths = append(paths, genPath)
			}
			return true
		})
	}
}

func (m *Memory) Link(from, to string) {
	m.lm.Scope(from).Link("link", to)
}

func (m *Memory) Check(from string, checker func(path string) (changed bool)) {
	m.lm.Scope(from).Select("check", checker)
}
