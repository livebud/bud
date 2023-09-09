package injector

import (
	"blog.com/controller"
	"github.com/livebud/bud/pkg/di"
	"github.com/livebud/bud/pkg/injector"
	"github.com/livebud/bud/pkg/web"
)

func New() di.Injector {
	in := injector.New()
	di.Provide[*controller.Controller](in, controller.New)
	di.Register[web.Router](in, controller.Register)
	return in
}
