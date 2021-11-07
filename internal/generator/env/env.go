package env

import (
	"gitlab.com/mnm/bud/gen"
)

func Generator() gen.Generator {
	return gen.GenerateFile(func(f gen.F, file *gen.File) error {
		return nil
	})
}
