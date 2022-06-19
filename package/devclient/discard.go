package devclient

import (
	"fmt"
	"net/http"

	hot "github.com/livebud/bud/package/hot2"
	"github.com/livebud/bud/runtime/view/ssr"
)

// Discard client implements Client
type discard struct {
}

var _ Client = discard{}

func (discard) Render(route string, props interface{}) (*ssr.Response, error) {
	return nil, fmt.Errorf("devclient: discard client does not support render")
}

func (discard) Proxy(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "devclient: discard client does not support proxy", http.StatusInternalServerError)
}

func (discard) Hot() (*hot.Stream, error) {
	return nil, fmt.Errorf("devclient: discard client does not support hot")
}

// Publish nothing
func (discard) Publish(topic string, data []byte) error {
	return nil
}
