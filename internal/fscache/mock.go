package fscache

import "github.com/livebud/bud/package/virtual"

type Mock struct {
	MockGet   func(name string) (entry virtual.Entry, ok bool)
	MockSet   func(path string, entry virtual.Entry)
	MockLink  func(from string, toPatterns ...string) error
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

func (m *Mock) Link(from string, toPatterns ...string) error {
	if m.MockLink != nil {
		return m.MockLink(from, toPatterns...)
	}
	return nil
}

func (m *Mock) Check(from string, checker func(path string) (changed bool)) {
	if m.MockCheck != nil {
		m.MockCheck(from, checker)
	}
}
