//go:build cgo
// +build cgo

// Package js does this
package v8

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"

	"gitlab.com/mnm/bud/js"
	"github.com/jackc/puddle"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"rogchap.com/v8go"
)

var errLocked = errors.New("v8: script can't be added after evaluating")

// TODO: Figure out a reasonable pool size. It's probably much higher than the
// number of CPUs.
var defaultPool = New()

func Eval(path, expr string) (string, error) {
	engine, err := defaultPool.Acquire(context.Background())
	if err != nil {
		return "", err
	}
	defer engine.Release()
	return engine.Eval(path, expr)
}

// New pool. The pool is primarily here to ensure that one thread is accessing
// an isolate at a given time.
func New() *Pool {
	return NewSize(int32(runtime.NumCPU()))
}

var _ js.VM = (*Pool)(nil)

func NewSize(maxSize int32) *Pool {
	pool := &Pool{
		locked: new(atomicLock),
	}
	pool.puddle = puddle.NewPool(pool.constructor, pool.destructor, maxSize)
	return pool
}

func loadV8() (*V8Context, error) {
	iso, err := v8go.NewIsolate()
	if err != nil {
		return nil, err
	}
	global, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, err
	}
	if err := fetch.InjectTo(iso, global); err != nil {
		return nil, err
	}
	context, err := v8go.NewContext(iso, global)
	if err != nil {
		iso.TerminateExecution()
		iso.Dispose()
		return nil, err
	}
	v8 := &V8Context{
		Isolate: iso,
		Context: context,
	}
	if err := addConsole(v8); err != nil {
		return nil, err
	}
	return v8, nil
}

func addConsole(v8 *V8Context) error {
	global := v8.Context.Global()
	console, err := v8go.NewObjectTemplate(v8.Isolate)
	if err != nil {
		return err
	}
	logfn, err := v8go.NewFunctionTemplate(v8.Isolate, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		for i, arg := range info.Args() {
			if i > 0 {
				fmt.Print(" ")
			}
			// TODO: support traversing objects
			fmt.Printf("%v", arg)
		}
		fmt.Print("\n")
		return nil
	})
	if err != nil {
		return err
	}
	if err := console.Set("log", logfn); err != nil {
		return err
	}
	consoleObj, err := console.NewInstance(v8.Context)
	if err != nil {
		return err
	}
	err = global.Set("console", consoleObj)
	if err != nil {
		return err
	}
	return nil
}

type V8Context struct {
	Isolate *v8go.Isolate
	Context *v8go.Context
}

func (c *V8Context) Eval(path, expr string) (value *v8go.Value, err error) {
	value, err = c.Context.RunScript(expr, path)
	if err != nil {
		if jsErr, ok := err.(*v8go.JSError); ok {
			// fmt.Println(jsErr.Location)
			return nil, jsErr
		}
		return nil, err
	}
	return value, nil
}

func (c *V8Context) Close() error {
	c.Context.Close()
	c.Isolate.TerminateExecution()
	c.Isolate.Dispose()
	return nil
}

// // Construct a V8 isolate with an initial script
// func scriptConstructor(path, source string) func(context.Context) (interface{}, error) {
// 	return func(context.Context) (interface{}, error) {
// 		v8, err := loadV8()
// 		if err != nil {
// 			return nil, err
// 		}
// 		// Initialize with a script
// 		if _, err := v8.Eval(path, source); err != nil {
// 			return nil, err
// 		}
// 		return v8, nil
// 	}
// }

// func constructor(context.Context) (interface{}, error) {
// 	return loadV8()
// }

// func destructor(value interface{}) {
// 	value.(*V8Context).Close()
// }

type Pool struct {
	puddle  *puddle.Pool
	locked  *atomicLock
	scripts []*script
}

func (p *Pool) constructor(context.Context) (interface{}, error) {
	p.locked.Lock()
	v8, err := loadV8()
	if err != nil {
		return nil, err
	}
	for _, script := range p.scripts {
		// Initialize with script
		if _, err := v8.Eval(script.Path, script.Code); err != nil {
			return nil, err
		}
	}
	return v8, nil
}

func (p *Pool) destructor(value interface{}) {
	value.(*V8Context).Close()
}

// Script adds a script before
func (p *Pool) Script(path, code string) error {
	if p.locked.Locked() {
		return errLocked
	}
	p.scripts = append(p.scripts, &script{path, code})
	return nil
}

func (p *Pool) Eval(path, expr string) (string, error) {
	engine, err := p.Acquire(context.Background())
	if err != nil {
		return "", err
	}
	defer engine.Release()
	return engine.Eval(path, expr)
}

func (p *Pool) Acquire(ctx context.Context) (*Engine, error) {
	resource, err := p.puddle.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &Engine{
		res: resource,
	}, nil
}

// Close the pool down
func (p *Pool) Close() {
	p.puddle.Close()
}

// Engine worker
type Engine struct {
	res *puddle.Resource
}

func (e *Engine) Eval(path, expr string) (string, error) {
	v8c := e.res.Value().(*V8Context)
	value, err := v8c.Eval(path, expr)
	if err != nil {
		return "", err
	}
	if value.IsPromise() {
		prom, err := value.AsPromise()
		if err != nil {
			return "", err
		}
		// TODO: this could run forever
		for prom.State() == v8go.Pending {
			continue
		}
		return prom.Result().String(), nil
	}
	return value.String(), nil
}

func (e *Engine) Release() {
	e.res.Release()
}

type script struct {
	Path string
	Code string
}

type atomicLock int32

func (al *atomicLock) Lock() {
	atomic.StoreInt32((*int32)(al), 1)
}

func (al *atomicLock) Locked() bool {
	return atomic.LoadInt32((*int32)(al))&1 == 1
}
