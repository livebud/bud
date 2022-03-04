package tracer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/mnm/bud/pkg/socket"
)

func NewClient(addr string) (*Client, error) {
	rt, err := socket.Transport(addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		c: &http.Client{
			Transport: rt,
		},
	}, nil
}

type Client struct {
	c *http.Client
}

var _ Exporter = (*Client)(nil)

func (c *Client) Print(ctx context.Context) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "http://traceserver/", nil)
	if err != nil {
		return "", err
	}
	res, err := c.c.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ExportSpans exports the spans to the server
func (c *Client) ExportSpans(ctx context.Context, spans []ReadOnlySpan) error {
	var scs []*SpanData
	for _, span := range spans {
		scs = append(scs, ToSpanData(span))
	}
	body, err := json.Marshal(scs)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, "http://traceserver/", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	res, err := c.c.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("trace: unexpected response status %s", res.Status)
	}
	return nil
}

// Nothing to shutdown
func (c *Client) Shutdown(ctx context.Context) error {
	return nil
}
