package buddy

import (
	"context"
	"fmt"
	"os/exec"
)

type App struct {
}

type startOption interface {
	start(o *startConfig)
}

type startConfig struct {
	Port string
}

func (a *App) Command(ctx context.Context, args ...string) *exec.Cmd {
	return nil
}

func (a *App) Start(ctx context.Context, options ...startOption) (*Process, error) {
	return nil, fmt.Errorf("not implemented yet")
}
