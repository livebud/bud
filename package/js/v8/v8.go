package v8

import (
	"github.com/livebud/bud/package/js"
	"go.kuoruan.net/v8go-polyfills/console"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"go.kuoruan.net/v8go-polyfills/url"
	"rogchap.com/v8go"
)

type Value = v8go.Value
type Error = v8go.JSError

func Eval(path, code string) (string, error) {
	vm, err := Load()
	if err != nil {
		return "", err
	}
	return vm.Eval(path, code)
}

func load() (*v8go.Isolate, *v8go.Context, error) {
	isolate := v8go.NewIsolate()
	global := v8go.NewObjectTemplate(isolate)
	// Fetch support
	if err := fetch.InjectTo(isolate, global); err != nil {
		isolate.TerminateExecution()
		isolate.Dispose()
		return nil, nil, err
	}
	// Create the context
	context := v8go.NewContext(isolate, global)
	// URL support
	if err := url.InjectTo(context); err != nil {
		context.Close()
		isolate.TerminateExecution()
		isolate.Dispose()
		return nil, nil, err
	}
	// Console support
	// TODO: this dependency looks like it can be improved
	if err := console.InjectTo(context); err != nil {
		context.Close()
		isolate.TerminateExecution()
		isolate.Dispose()
		return nil, nil, err
	}
	return isolate, context, nil
}

func Load() (*VM, error) {
	isolate, context, err := load()
	if err != nil {
		return nil, err
	}
	return &VM{
		isolate: isolate,
		context: context,
	}, nil
}

func Compile(path, code string) (*VM, error) {
	isolate, context, err := load()
	if err != nil {
		return nil, err
	}
	script, err := isolate.CompileUnboundScript(code, path, v8go.CompileOptions{})
	if err != nil {
		return nil, err
	}
	// Bind to the context
	if _, err := script.Run(context); err != nil {
		return nil, err
	}
	return &VM{
		isolate: isolate,
		context: context,
	}, nil
}

type VM struct {
	isolate *v8go.Isolate
	context *v8go.Context
}

var _ js.VM = (*VM)(nil)

// Compile a script into the context
func (vm *VM) Script(path, code string) error {
	script, err := vm.isolate.CompileUnboundScript(code, path, v8go.CompileOptions{})
	if err != nil {
		return err
	}
	// Bind to the context
	if _, err := script.Run(vm.context); err != nil {
		return err
	}
	return nil
}

func (vm *VM) Eval(path, expr string) (string, error) {
	value, err := vm.context.RunScript(expr, path)
	if err != nil {
		return "", err
	}
	// Handle promises
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

func (vm *VM) Close() {
	vm.context.Close()
	vm.isolate.TerminateExecution()
	vm.isolate.Dispose()
}
