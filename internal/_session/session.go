package session

import (
	"github.com/go-duo/bud/go/mod"
	"github.com/go-duo/bud/log"
)

type Config struct {
	Modfile mod.File
	Logger  log.Log
	Embed   bool
	Hot     bool
	Minify  bool
}
