package budrpc

import (
	"context"
	"encoding/json"
	"io/fs"

	"github.com/keegancsmith/rpc"

	"github.com/livebud/bud/framework/view/ssr"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/socket"
)

type Client interface {
	Render(route string, props interface{}) (*ssr.Response, error)
	Open(name string) (fs.File, error)
	Publish(topic string, data []byte) error
}

func Dial(ctx context.Context, addr string) (Client, error) {
	conn, err := socket.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}
	rpc := rpc.NewClient(conn)
	rfs := remotefs.NewClient(rpc)
	return &budClient{rpc, rfs, context.Background()}, nil
}

type budClient struct {
	rpc *rpc.Client
	rfs *remotefs.Client
	ctx context.Context
}

type RenderRequest struct {
	Route string
	Props json.RawMessage
}

type RenderResponse = ssr.Response

func (c *budClient) WithContext(ctx context.Context) *budClient {
	return &budClient{c.rpc, c.rfs, ctx}
}

func (c *budClient) Render(route string, props interface{}) (*RenderResponse, error) {
	raw, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	var res RenderResponse
	if err := c.rpc.Call(c.ctx, "bud.Render", RenderRequest{route, raw}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *budClient) Open(name string) (fs.File, error) {
	return c.rfs.Open(name)
}

type PublishRequest struct {
	Topic string
	Data  []byte
}

func (c *budClient) Publish(topic string, data []byte) error {
	if err := c.rpc.Call(c.ctx, "bud.Publish", PublishRequest{topic, data}, nil); err != nil {
		return err
	}
	return nil
}
