package cli

import (
	"context"
	"fmt"

	"github.com/livebud/bud/internal/once"
)

func newBudFS(provider *provider) *budfs {
	return &budfs{
		provider,
		&once.Closer{},
	}
}

type budfs struct {
	p *provider
	c *once.Closer
}

func (f *budfs) Sync(ctx context.Context, dirs ...string) error {
	fmt.Println("syncing....", dirs)
	return nil
}

func (f *budfs) Refresh(ctx context.Context, paths ...string) error {
	return nil
}

func (f *budfs) Close() error {
	return f.p.Close()
}
