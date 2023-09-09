package injector

import (
	"github.com/livebud/bud/pkg/di"
	"github.com/livebud/bud/pkg/injector"
)

func New() di.Injector {
	in := injector.New()
	return in
}
