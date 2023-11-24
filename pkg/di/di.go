package di

import (
	"fmt"
	"sync"

	"github.com/livebud/bud/internal/reflector"
)

var ErrNoLoader = fmt.Errorf("di: no loader")

func New() Injector {
	return &injector{
		loaders: make(map[string]any),
		cache:   make(map[string]any),
		appends: make(map[string][]any),
	}
}

// Clone an injector
func Clone(in Injector) Injector {
	return in.clone()
}

type Injector interface {
	clone() Injector
	setLoader(name string, loader any) error
	append(name string, loader any) error
	getAppends(name string) (registrants []any, ok bool)
	getLoader(name string) (loader any, ok bool)
	setCache(name string, dep any)
	getCache(name string) (dep any, ok bool)
}

type injector struct {
	mu      sync.RWMutex
	loaders map[string]any
	cache   map[string]any
	appends map[string][]any
}

var _ Injector = (*injector)(nil)

func (in *injector) clone() Injector {
	clone := New()
	for name, provider := range in.loaders {
		clone.setLoader(name, provider)
	}
	for name, dep := range in.cache {
		clone.setCache(name, dep)
	}
	for name, registrants := range in.appends {
		for _, registrant := range registrants {
			clone.append(name, registrant)
		}
	}
	return clone
}

func (in *injector) setLoader(name string, provider any) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.loaders[name] = provider
	return nil
}

func (in *injector) append(name string, provider any) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	if _, ok := in.appends[name]; !ok {
		in.appends[name] = []any{}
	}
	in.appends[name] = append(in.appends[name], provider)
	return nil
}

func (in *injector) getAppends(name string) (registrants []any, ok bool) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	registrants, ok = in.appends[name]
	return registrants, ok
}

func (in *injector) getLoader(name string) (provider any, ok bool) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	if fn, ok := in.loaders[name]; ok {
		return fn, ok
	}
	return nil, false
}

func (in *injector) setCache(name string, dep any) {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.cache[name] = dep
}

func (in *injector) getCache(name string) (dep any, ok bool) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	dep, ok = in.cache[name]
	return dep, ok
}

func Loader[Dep any](in Injector, fn func(in Injector) (d Dep, err error)) error {
	var dep Dep
	name, err := reflector.TypeOf(dep)
	if err != nil {
		return err
	}
	in.setLoader(name, fn)
	return nil
}

// Append a loader to a dependency. These loaders will be called after the
// dependency is loaded.
func Append[To any](in Injector, loader func(in Injector, to To) error) error {
	var to To
	name, err := reflector.TypeOf(to)
	if err != nil {
		return err
	}
	return in.append(name, loader)
}

// Load a dependency from an injector and cache the result for future loads
func Load[Dep any](in Injector) (dep Dep, err error) {
	name, err := reflector.TypeOf(dep)
	if err != nil {
		return dep, err
	}
	if dep, ok := in.getCache(name); ok {
		return dep.(Dep), nil
	}
	v, ok := in.getLoader(name)
	if !ok {
		return dep, fmt.Errorf("%w for %s", ErrNoLoader, name)
	}
	fn, ok := v.(func(in Injector) (Dep, error))
	if !ok {
		return dep, fmt.Errorf("di: invalid provider for %s", name)
	}
	d, err := fn(in)
	if err != nil {
		return dep, fmt.Errorf("di: unable to load %q: %w", name, err)
	}
	registrants, ok := in.getAppends(name)
	if ok {
		for _, registrant := range registrants {
			fn, ok := registrant.(func(in Injector, dep Dep) error)
			if !ok {
				return dep, fmt.Errorf("di: invalid registrant for %s", name)
			}
			if err := fn(in, d); err != nil {
				return dep, err
			}
		}
	}
	in.setCache(name, d)
	return d, nil
}
