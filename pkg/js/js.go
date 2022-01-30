package js

// VM for evaluating javascript
type VM interface {
	Script(path, script string) error
	Eval(path, expression string) (string, error)
}
