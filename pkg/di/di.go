package di

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/livebud/bud/pkg/di/internal/reflector"
)

var ErrNoProvider = fmt.Errorf("di: no provider")

type dependency struct {
	ID       string
	ArgIDs   []string
	provider reflect.Value
}

type Injector interface {
	load(depId string) (val reflect.Value, err error)
	provide(depId string, dep *dependency)
	when(depId string, dep *dependency)
	clone() Injector
}

type injector struct {
	mu    sync.RWMutex
	deps  map[string]*dependency
	whens map[string][]*dependency
	cache map[string]reflect.Value
}

func (in *injector) provide(depId string, dep *dependency) {
	in.mu.Lock()
	in.deps[depId] = dep
	in.mu.Unlock()
}

func (in *injector) when(depId string, dep *dependency) {
	in.mu.Lock()
	in.whens[depId] = append(in.whens[depId], dep)
	in.mu.Unlock()
}

func (in *injector) load(depId string) (val reflect.Value, err error) {
	if val, ok := in.cache[depId]; ok {
		return val, nil
	}
	dep, ok := in.deps[depId]
	if !ok {
		return val, fmt.Errorf("%w for %s", ErrNoProvider, depId)
	}
	args := make([]reflect.Value, len(dep.ArgIDs))
	for i, argId := range dep.ArgIDs {
		arg, err := in.load(argId)
		if err != nil {
			return val, err
		}
		args[i] = arg
	}
	results := dep.provider.Call(args)
	numResults := len(results)
	if numResults == 0 || numResults > 2 {
		return val, fmt.Errorf("di: wrong number of results while loading %s", depId)
	} else if numResults == 2 {
		if err, ok := results[1].Interface().(error); ok && err != nil {
			return val, err
		}
	}
	val = results[0]
	in.cache[depId] = val
	if err := in.loadWhens(depId); err != nil {
		return val, err
	}
	return val, nil
}

func (in *injector) loadWhens(depId string) error {
	for _, when := range in.whens[depId] {
		args := make([]reflect.Value, len(when.ArgIDs))
		for i, argId := range when.ArgIDs {
			arg, err := in.load(argId)
			if err != nil {
				return err
			}
			args[i] = arg
		}
		results := when.provider.Call(args)
		numResults := len(results)
		if numResults > 1 {
			return fmt.Errorf("di: wrong number of results while loading %s", depId)
		} else if numResults == 1 {
			if err, ok := results[0].Interface().(error); ok && err != nil {
				return err
			}
		}
	}
	return nil
}

func (in *injector) clone() Injector {
	in.mu.RLock()
	defer in.mu.RUnlock()
	deps := make(map[string]*dependency, len(in.deps))
	for k, v := range in.deps {
		deps[k] = v
	}
	whens := make(map[string][]*dependency, len(in.whens))
	for k, v := range in.whens {
		whens[k] = v
	}
	cache := make(map[string]reflect.Value, len(in.cache))
	for k, v := range in.cache {
		cache[k] = v
	}
	return &injector{
		deps:  deps,
		whens: whens,
		cache: cache,
	}
}

func New() Injector {
	return &injector{
		deps:  make(map[string]*dependency),
		whens: make(map[string][]*dependency),
		cache: make(map[string]reflect.Value),
	}
}

// Load a dependency. Requires a corresponding provider.
func Load[Dep any](in Injector) (d Dep, err error) {
	depType := reflect.TypeOf(d)
	depId, err := reflector.ID(depType)
	if err != nil {
		return d, err
	}
	val, err := load(in, depId)
	if err != nil {
		return d, err
	}
	dep, ok := val.Interface().(Dep)
	if !ok {
		return d, fmt.Errorf("di: invalid provider for %s", depType)
	}

	return dep, nil
}

func load(in Injector, depId string) (val reflect.Value, err error) {
	return in.load(depId)
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// Provide a function for initializing a dependency
func Provide[Dep, Func any](in Injector, fn Func) error {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("di: %s must be a function", fnType)
	}
	var dep Dep
	depType := reflect.TypeOf(dep)
	depID, err := reflector.ID(depType)
	if err != nil {
		return err
	}
	numResults := fnType.NumOut()
	if numResults == 0 || numResults > 2 {
		return fmt.Errorf("di: invalid provider for %s", fnType)
	}
	if numResults == 1 || numResults == 2 {
		resultType := fnType.Out(0)
		resultID, err := reflector.ID(resultType)
		if err != nil {
			return err
		}
		if resultID != depID && !depType.Implements(resultType) {
			return fmt.Errorf("di: invalid provider for %s", fnType)
		}
	}
	if numResults == 2 {
		errType := fnType.Out(1)
		if !errType.Implements(errorType) {
			return fmt.Errorf("di: invalid provider for %s", fnType)
		}
	}
	argIDs := make([]string, fnType.NumIn())
	for i := range argIDs {
		argType := fnType.In(i)
		argID, err := reflector.ID(argType)
		if err != nil {
			return err
		}
		argIDs[i] = argID
	}
	fnValue := reflect.ValueOf(fn)
	in.provide(depID, &dependency{
		ID:       depID,
		ArgIDs:   argIDs,
		provider: fnValue,
	})
	return nil
}

// When a dependency is loaded, call this function as well
func When[Dep, Func any](in Injector, fn Func) error {
	var dep Dep
	depType := reflect.TypeOf(dep)
	depID, err := reflector.ID(depType)
	if err != nil {
		return err
	}
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("di: %s must be a function", fnType)
	}
	numOut := fnType.NumOut()
	if numOut > 1 {
		return fmt.Errorf("di: invalid function signature for %s", fnType)
	} else if numOut == 1 {
		errType := fnType.Out(0)
		if !errType.Implements(errorType) {
			return fmt.Errorf("di: expected an error result type but got %s", fnType)
		}
	}
	args := make([]string, fnType.NumIn())
	for i := range args {
		argType := fnType.In(i)
		argID, err := reflector.ID(argType)
		if err != nil {
			return err
		}
		args[i] = argID
	}
	fnValue := reflect.ValueOf(fn)
	in.when(depID, &dependency{
		ID:       depID,
		ArgIDs:   args,
		provider: fnValue,
	})
	return nil
}

func Clone(in Injector) Injector {
	return in.clone()
}
