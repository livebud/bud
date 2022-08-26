package budhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/livebud/bud/framework/view/ssr"
	"github.com/livebud/bud/internal/urlx"
	"github.com/livebud/bud/internal/virtual"
	"github.com/livebud/bud/package/socket"
)

type Client interface {
	Render(route string, props interface{}) (*ssr.Response, error)
	Publish(topic string, data []byte) error
	Open(name string) (fs.File, error)
}

// Try tries loading a dev client from an environment variable or returns an
// empty client if no environment variable is set
func Try(addr string) (Client, error) {
	if addr == "" {
		return discard{}, nil
	}
	return Load(addr)
}

// Load a client from an address
func Load(addr string) (Client, error) {
	url, err := urlx.Parse(addr)
	if err != nil {
		return nil, err
	}
	transport, err := socket.Transport(addr)
	if err != nil {
		return nil, fmt.Errorf("budhttp: unable to create transport from listener. %w", err)
	}
	httpClient := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &client{
		baseURL:    url.String(),
		httpClient: httpClient,
	}, nil
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

var _ Client = (*client)(nil)

// Render a path with props on the dev server
func (c *client) Render(route string, props interface{}) (*ssr.Response, error) {
	body, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	url := strings.TrimSuffix(c.baseURL+"/bud/view"+route, "/")
	res, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("budhttp: render returned unexpected %d. %s", res.StatusCode, resBody)
	}
	out := new(ssr.Response)
	if err := json.Unmarshal(resBody, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *client) Open(name string) (fs.File, error) {
	res, err := c.httpClient.Get(c.baseURL + "/" + name)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("budhttp: open returned unexpected %d. %s", res.StatusCode, body)
	}
	return virtual.UnmarshalJSON(body)
}

type Event struct {
	Topic string `json:"topic,omitempty"`
	Data  []byte `json:"data,omitempty"`
}

func (c *client) Publish(topic string, data []byte) error {
	body, err := json.Marshal(Event{topic, data})
	if err != nil {
		return err
	}
	url := c.baseURL + "/bud/events"
	res, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("budhttp: send returned unexpected %d. %s", res.StatusCode, resBody)
	}
	return nil
}
