package budclient

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/livebud/bud/framework/view/ssr"
)

// Discard client implements Client
type discard struct {
}

var _ Client = discard{}

func (discard) Script(path, script string) error {
	return fmt.Errorf("budclient: discard script not implemented yet")
}

func (discard) Eval(path, expression string) (string, error) {
	return "", fmt.Errorf("budclient: discard eval not implemented yet")
}

func (discard) Open(name string) (fs.File, error) {
	return nil, fmt.Errorf("budclient: discard open not implemented yet")
}

func (discard) Render(route string, props interface{}) (*ssr.Response, error) {
	return nil, fmt.Errorf("budclient: discard client does not support render")
}

func (discard) Proxy(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "budclient: discard client does not support proxy", http.StatusInternalServerError)
}

// Publish nothing
func (discard) Publish(topic string, data []byte) error {
	return nil
}
