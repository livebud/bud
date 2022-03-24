package v8client

import "github.com/jackc/puddle"

// ClientResource worker
type ClientResource struct {
	res *puddle.Resource
}

func (e *ClientResource) Eval(path, expr string) (string, error) {
	command := e.res.Value().(*Command)
	value, err := command.Eval(path, expr)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (e *ClientResource) Release() {
	e.res.Release()
}
