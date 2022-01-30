package v8client

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"

	"github.com/jackc/puddle"
	"gitlab.com/mnm/bud/pkg/js"
)

var errLocked = errors.New("v8client: script can't be added after evaluating")

// New Client, using v8 as a sidecar process. We use this in development to
// speed up rebuilds by not linking V8 every time. It may also be used in
// environments without cgo.
//
// This requires the bud CLI to be in $PATH and is not recommended for
// production.
func New(command string, args ...string) *Client {
	return NewWithSize(int32(runtime.NumCPU()), command, args...)
}

// Default launches the default client
func Default() *Client {
	return New("bud", "tool", "v8", "client")
}

var _ js.VM = (*Client)(nil)

func NewWithSize(maxSize int32, command string, args ...string) *Client {
	pool := &Client{
		locked:  new(atomicLock),
		command: command,
		args:    args,
	}
	pool.puddle = puddle.NewPool(pool.constructor, pool.destructor, maxSize)
	return pool
}

type Client struct {
	locked  *atomicLock
	command string
	args    []string
	puddle  *puddle.Pool
	scripts []*script
}

func (c *Client) constructor(ctx context.Context) (interface{}, error) {
	c.locked.Lock()
	command, err := launchBudToolV8(ctx, c.command, c.args...)
	if err != nil {
		return nil, err
	}
	for _, script := range c.scripts {
		// Initialize with script
		if _, err := command.Eval(script.Path, script.Code); err != nil {
			return nil, err
		}
	}
	return command, nil
}

func (c *Client) destructor(value interface{}) {
	value.(*Command).Close()
}

// Script adds a script before
func (c *Client) Script(path, code string) error {
	if c.locked.Locked() {
		return errLocked
	}
	c.scripts = append(c.scripts, &script{path, code})
	return nil
}

func (c *Client) Eval(path, expr string) (string, error) {
	engine, err := c.Acquire(context.Background())
	if err != nil {
		return "", err
	}
	defer engine.Release()
	return engine.Eval(path, expr)
}

func (c *Client) Acquire(ctx context.Context) (*ClientResource, error) {
	resource, err := c.puddle.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &ClientResource{
		res: resource,
	}, nil
}

// Close the pool down
func (c *Client) Close() {
	c.puddle.Close()
}

type script struct {
	Path string
	Code string
}

type atomicLock int32

func (al *atomicLock) Lock() {
	atomic.StoreInt32((*int32)(al), 1)
}

func (al *atomicLock) Locked() bool {
	return atomic.LoadInt32((*int32)(al))&1 == 1
}
