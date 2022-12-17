package vcache

import "github.com/livebud/bud/package/virtual"

func Mock(fallback Cache) *MockCache {
	return &MockCache{Cache: fallback}
}

type MockCache struct {
	Cache
	MockHas    func(path string) (ok bool)
	MockGet    func(path string) (entry virtual.Entry, ok bool)
	MockSet    func(path string, entry virtual.Entry)
	MockDelete func(path string)
	MockRange  func(fn func(path string, entry virtual.Entry) bool)
	MockClear  func()
}

func (m *MockCache) Has(path string) (ok bool) {
	if m.MockHas != nil {
		return m.MockHas(path)
	}
	return m.Cache.Has(path)
}

func (m *MockCache) Get(path string) (entry virtual.Entry, ok bool) {
	if m.MockGet != nil {
		return m.MockGet(path)
	}
	return m.Cache.Get(path)
}

func (m *MockCache) Set(path string, entry virtual.Entry) {
	if m.MockSet != nil {
		m.MockSet(path, entry)
		return
	}
	m.Cache.Set(path, entry)
}

func (m *MockCache) Delete(path string) {
	if m.MockDelete != nil {
		m.MockDelete(path)
		return
	}
	m.Cache.Delete(path)
}

func (m *MockCache) Range(fn func(path string, entry virtual.Entry) bool) {
	if m.MockRange != nil {
		m.MockRange(fn)
		return
	}
	m.Cache.Range(fn)
}

func (m *MockCache) Clear() {
	if m.MockClear != nil {
		m.MockClear()
	}
	m.Cache.Clear()
}
