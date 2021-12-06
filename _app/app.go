package app

import (
	"gitlab.com/mnm/bud/vfs"

	"gitlab.com/mnm/bud/go/mod"
)

func Load(dir string) (*Config, error) {
	return nil, nil
}

type Config struct {
	Mod mod.File
	FS  vfs.ReadWritable
}
