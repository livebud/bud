package fscache

import "github.com/livebud/bud/package/virtual"

type Mock struct {
	MockGet   func(name string) (entry virtual.Entry, ok bool)
	MockSet   func(path string, entry virtual.Entry)
	MockLink  func(from, to string)
	MockCheck func(from string, checker func(path string) (changed bool))
}

func (m *Mock) Get(name string) (entry virtual.Entry, ok bool) {
	if m.MockGet != nil {
		return m.MockGet(name)
	}
	return
}

func (m *Mock) Set(path string, entry virtual.Entry) {
	if m.MockSet != nil {
		m.MockSet(path, entry)
	}
}

func (m *Mock) Link(from, to string) {
	if m.MockLink != nil {
		m.MockLink(from, to)
	}
}

func (m *Mock) Check(from string, checker func(path string) (changed bool)) {
	if m.MockCheck != nil {
		m.MockCheck(from, checker)
	}
}
