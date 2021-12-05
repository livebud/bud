package di_test

import (
	"testing"

	"gitlab.com/mnm/bud/internal/di"
)

func TestWire(t *testing.T) {
	stct := &di.Struct2{
		Name: "Command",
		Fields: []*di.StructField{
			{
				Name:   "Web",
				Import: "myapp.com/web",
				Type:   "*Server",
			},
		},
	}
	fn := &di.Function2{
		Name:    "Load",
		ModFile: nil,
		Results: []di.Result{
			stct,
		},
	}
	provider, err := fn.Generate(target)
}
