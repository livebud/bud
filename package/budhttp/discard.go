package budhttp

import (
	"fmt"

	"github.com/livebud/bud/framework/view/ssr"
)

// Discard client implements Client
type discard struct {
}

var _ Client = discard{}

func (discard) Render(route string, props interface{}) (*ssr.Response, error) {
	return nil, fmt.Errorf("budhttp: discard client does not support render")
}

func (discard) Script(path, script string) error {
	return fmt.Errorf("budhttp: discard client does not support script")
}

func (discard) Eval(path, expression string) (string, error) {
	return "", fmt.Errorf("budhttp: discard client does not support eval")
}

// Publish nothing
func (discard) Publish(topic string, data []byte) error {
	return nil
}
