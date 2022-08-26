package budhttp

import (
	"fmt"
	"io/fs"

	"github.com/livebud/bud/framework/view/ssr"
)

// Discard client implements Client
type discard struct {
}

var _ Client = discard{}

func (discard) Render(route string, props interface{}) (*ssr.Response, error) {
	return nil, fmt.Errorf("budhttp: discard client does not support render")
}

func (discard) Open(name string) (fs.File, error) {
	return nil, fmt.Errorf("budhttp: discard client does not support open")
}

// Publish nothing
func (discard) Publish(topic string, data []byte) error {
	return nil
}
