package v8client

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"

	"github.com/go-duo/bud/js"
	"github.com/jackc/puddle"
)

var errLocked = errors.New("v8client: script can't be added after evaluating")

// Launch V8 as a sidecar process. We use this in development to speed up
// rebuilds by not linking V8 every time. It may also be used in environments
// without cgo.
func Launch(command string, args ...string) *Client {
	return LaunchWithSize(int32(runtime.NumCPU()), command, args...)
}

var _ js.VM = (*Client)(nil)

func LaunchWithSize(maxSize int32, command string, args ...string) *Client {
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

func launchV8(ctx context.Context, command string, args ...string) (*Command, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Command{
		cmd:    cmd,
		stdin:  json.NewEncoder(stdin),
		stdout: json.NewDecoder(stdout),
	}, nil
}

type Command struct {
	cmd    *exec.Cmd
	stdin  *json.Encoder
	stdout *json.Decoder
}

func (c *Command) Eval(path, expr string) (value string, err error) {
	if err := c.stdin.Encode(expr); err != nil {
		return "", err
	}
	var raw string
	if err := c.stdout.Decode(&raw); err != nil {
		return "", err
	}
	return string(raw), nil
}

func (c *Command) Close() error {
	if c.cmd.Process == nil {
		return nil
	}
	if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}
	if err := c.cmd.Wait(); err != nil && err.Error() != "signal: interrupt" {
		return err
	}
	return nil
}

func (c *Client) constructor(ctx context.Context) (interface{}, error) {
	c.locked.Lock()
	command, err := launchV8(ctx, c.command, c.args...)
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

// ClientResource worker
type ClientResource struct {
	res *puddle.Resource
}

func (e *ClientResource) Eval(path, expr string) (string, error) {
	command := e.res.Value().(*Command)
	value, err := command.Eval(path, expr)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (e *ClientResource) Release() {
	e.res.Release()
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
