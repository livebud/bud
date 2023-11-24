//go:build embed

package asset

import (
	"embed"
	"io/fs"

	"github.com/livebud/bud/_example/jack/env"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
)

//go:embed .embed/**
var embedded embed.FS

func Load(env *env.Env, log logs.Log, module *mod.Module) fs.FS {
	fsys, err := fs.Sub(embedded, ".embed")
	if err != nil {
		panic(err)
	}
	return fsys
}
