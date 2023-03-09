package cli

import "github.com/livebud/bud/framework/generator"

// TODO: consider moving these into bud/cmd/bud/generator.go
var fileGenerators = map[string]generator.Selector{
	"bud/cmd/app/main.go": {
		Import: "github.com/livebud/bud/framework/app",
		Type:   "*Generator",
	},
	"bud/internal/web/web.go": {
		Import: "github.com/livebud/bud/framework/web",
		Type:   "*Generator",
	},
	"bud/internal/web/controller/controller.go": {
		Import: "github.com/livebud/bud/framework/controller",
		Type:   "*Generator",
	},
	"bud/internal/web/view/view.go": {
		Import: "github.com/livebud/bud/framework/view",
		Type:   "*Generator",
	},
	"bud/internal/web/public/public.go": {
		Import: "github.com/livebud/bud/framework/public",
		Type:   "*Generator",
	},
	"bud/view/_ssr.js": {
		Import: "github.com/livebud/bud/framework/view/ssr",
		Type:   "*Generator",
	},
}

var fileServers = map[string]generator.Selector{
	"bud/view": {
		Import: "github.com/livebud/bud/framework/view/dom",
		Type:   "*Generator",
	},
	"bud/node_modules": {
		Import: "github.com/livebud/bud/framework/view/nodemodules",
		Type:   "*Generator",
	},
}
