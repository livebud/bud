package env

import "gitlab.com/mnm/bud/pkg/gen"

func Generator() gen.Generator {
	return gen.GenerateFile(func(f gen.F, file *gen.File) error {
		return nil
	})
}
