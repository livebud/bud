package devclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/urlx"
	hot "github.com/livebud/bud/package/hot2"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/runtime/view/ssr"
)

// FromFile loads the client from a file descriptor
func FromFile() (*Client, error) {
	files := extrafile.Load("BUD")
	if len(files) == 0 {
		return nil, fmt.Errorf("devclient: no listener file passed in")
	}
	listener, err := socket.From(files[0])
	if err != nil {
		return nil, fmt.Errorf("devclient: unable to listen to file. %w", err)
	}
	return Load(listener.Addr().String())
}

// Load a client from an address
func Load(addr string) (*Client, error) {
	url, err := urlx.Parse(addr)
	if err != nil {
		return nil, err
	}
	transport, err := socket.Transport(addr)
	if err != nil {
		return nil, fmt.Errorf("devclient: unable to create transport from listener. %w", err)
	}
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Client{
		baseURL:    url.String(),
		httpClient: client,
	}, nil
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Render a path with props on the dev server
func (c *Client) Render(route string, props interface{}) (*ssr.Response, error) {
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
		return nil, fmt.Errorf("devclient: render returned unexpected %d. %s", res.StatusCode, resBody)
	}
	out := new(ssr.Response)
	if err := json.Unmarshal(resBody, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Proxy a file to the dev server
// func (c *Client) Proxy(urlPath string) (*http.Response, error) {
// 	return c.httpClient.Get(c.baseURL + urlPath)
// }

func (c *Client) Proxy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	res, err := c.httpClient.Get(c.baseURL + r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), res.StatusCode)
		return
	}
	defer res.Body.Close()
	headers := w.Header()
	for key := range res.Header {
		headers.Set(key, res.Header.Get(key))
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func (c *Client) Hot() (*hot.Stream, error) {
	return hot.DialWith(c.httpClient, c.baseURL+"/bud/hot")
}

type Event struct {
	Type string
	Data []byte
}

func (c *Client) Send(event Event) error {
	body, err := json.Marshal(event)
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
		return fmt.Errorf("devclient: send returned unexpected %d. %s", res.StatusCode, resBody)
	}
	return nil
}
