package viewrt

import (
	"net/http"

	"github.com/livebud/bud/package/router"
)

type Viewer interface {
	// TODO: remove handler
	Handler(route string, props interface{}) http.Handler
	Render(w http.ResponseWriter, viewPath string, props interface{})
	Serve(router *router.Router)
}
