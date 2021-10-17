package session

import (
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/log"
)

type Config struct {
	Modfile mod.File
	Logger  log.Log
	Embed   bool
	Hot     bool
	Minify  bool
}
