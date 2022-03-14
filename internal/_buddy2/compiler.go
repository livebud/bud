package buddy

import "context"

type Option struct {
	Trace  bool
	Embed  bool
	Minify bool
	Hot    bool
}

type Compiler struct {
	*Option
	Dir string
}

func (c *Compiler) Compile(ctx context.Context) (*CLI, error) {
	return nil, nil
}

type Builder struct {
	*Compiler
}

type Runner struct {
	*Compiler
	Port string
	Hot  bool
}
