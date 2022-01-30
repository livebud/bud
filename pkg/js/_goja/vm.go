package goja

import (
	"github.com/dop251/goja"
	"gitlab.com/mnm/bud/pkg/js"
)

func New() *VM {
	return &VM{goja.New()}
}

var _ js.VM = (*VM)(nil)

type VM struct {
	vm *goja.Runtime
}

func (v *VM) Script(path, script string) error {
	_, err := v.vm.RunScript(path, script)
	return err
}

func (v *VM) Eval(path, expression string) (string, error) {
	value, err := v.vm.RunString(expression)
	if err != nil {
		return "", err
	}
	return value.String(), nil
}
