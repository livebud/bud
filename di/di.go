package di

import (
	"fmt"
	"sync"

	"github.com/livebud/bud/internal/reflector"
)

var ErrNoProvider = fmt.Errorf("di: no provider")

func New() Injector {
	return &injector{
		fns:        make(map[string]any),
		cache:      make(map[string]any),
		registered: make(map[string][]any),
	}
}

// Clone an injector
func Clone(in Injector) Injector {
	return in.clone()
}

type Injector interface {
	clone() Injector
	setProvider(name string, provider any) error
	register(name string, provider any) error
	registrants(name string) (registrants []any, ok bool)
	getProvider(name string) (provider any, ok bool)
	setCache(name string, dep any)
	getCache(name string) (dep any, ok bool)
}

type injector struct {
	mu         sync.RWMutex
	fns        map[string]any
	cache      map[string]any
	registered map[string][]any
}

var _ Injector = (*injector)(nil)

func (in *injector) clone() Injector {
	clone := New()
	for name, provider := range in.fns {
		clone.setProvider(name, provider)
	}
	for name, dep := range in.cache {
		clone.setCache(name, dep)
	}
	for name, registrants := range in.registered {
		for _, registrant := range registrants {
			clone.register(name, registrant)
		}
	}
	return clone
}

func (in *injector) setProvider(name string, provider any) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.fns[name] = provider
	return nil
}

func (in *injector) register(name string, provider any) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	if _, ok := in.registered[name]; !ok {
		in.registered[name] = []any{}
	}
	in.registered[name] = append(in.registered[name], provider)
	return nil
}

func (in *injector) registrants(name string) (registrants []any, ok bool) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	registrants, ok = in.registered[name]
	return registrants, ok
}

func (in *injector) getProvider(name string) (provider any, ok bool) {
	in.mu.RLock()
	defer in.mu.RUnlock()
	if fn, ok := in.fns[name]; ok {
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
	in.setProvider(name, fn)
	return nil
}

func Extend[Dep any](in Injector, with func(in Injector, dep Dep) error) error {
	var dep Dep
	name, err := reflector.TypeOf(dep)
	if err != nil {
		return err
	}
	return in.register(name, with)
}

func Has[Dep any](in Injector) bool {
	var dep Dep
	name, err := reflector.TypeOf(dep)
	if err != nil {
		return false
	}
	if _, ok := in.getProvider(name); ok {
		return true
	}
	return false
}

func Load[Dep any](in Injector) (dep Dep, err error) {
	name, err := reflector.TypeOf(dep)
	if err != nil {
		return dep, err
	}
	if dep, ok := in.getCache(name); ok {
		return dep.(Dep), nil
	}
	v, ok := in.getProvider(name)
	if !ok {
		return dep, fmt.Errorf("%w for %s", ErrNoProvider, name)
	}
	fn, ok := v.(func(in Injector) (Dep, error))
	if !ok {
		return dep, fmt.Errorf("di: invalid provider for %s", name)
	}
	d, err := fn(in)
	if err != nil {
		return dep, fmt.Errorf("di: unable to load %q: %w", name, err)
	}
	registrants, ok := in.registrants(name)
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
