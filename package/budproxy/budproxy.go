package budproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/livebud/bud/package/socket"

	"github.com/livebud/bud/internal/extrafile"
)

func Load() (*Proxy, error) {
	files := extrafile.Load("BUD")
	if len(files) == 0 {
		return nil, fmt.Errorf("budproxy: no listener file passed in")
	}
	listener, err := socket.From(files[0])
	if err != nil {
		return nil, fmt.Errorf("budproxy: unable to listen to file. %w", err)
	}
	transport, err := socket.Transport(listener.Addr().String())
	if err != nil {
		return nil, fmt.Errorf("budproxy: unable to create transport from listener. %w", err)
	}
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Proxy{"http://" + listener.Addr().String(), client}, nil
}

type Proxy struct {
	baseURL string
	client  *http.Client
}

func (p *Proxy) Render(path string, props interface{}) (*http.Response, error) {
	body, err := json.Marshal(struct {
		Path  string      `json:"path"`
		Props interface{} `json:"props"`
	}{path, props})
	if err != nil {
		return nil, err
	}
	return p.client.Post(p.baseURL+"/bud/view/_ssr.js", "application/json", bytes.NewReader(body))
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Proxy the request to the bud server
	res, err := p.client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := res.Write(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
