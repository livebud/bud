package view_test

import (
	"context"
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/framework"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/svelte"

	"github.com/livebud/bud/runtime/transform"

	"github.com/livebud/bud/framework/view"
	"github.com/livebud/bud/internal/golden"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/gomod"
)

func load(ctx context.Context, fsys fs.FS, module *gomod.Module, flag *framework.Flag) (*view.State, error) {
	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	svelteCompiler, err := svelte.Load(vm)
	if err != nil {
		return nil, err
	}
	transform, err := transform.Load(svelte.NewTransformable(svelteCompiler))
	if err != nil {
		return nil, err
	}
	return view.Load(ctx, fsys, module, transform, flag)
}

func TestEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	state, err := load(ctx, module, module, &framework.Flag{})
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(state, nil)
}

func TestEmptyViewDir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.MapFiles["view"] = &fstest.MapFile{Mode: fs.ModeDir}
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	state, err := load(ctx, module, module, &framework.Flag{})
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(state, nil)
}

func TestIndex(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	state, err := load(ctx, module, module, &framework.Flag{})
	is.NoErr(err)
	code, err := view.Generate(state)
	is.NoErr(err)
	is.NoErr(parser.Check(code))
	golden.TestGenerator(t, state, code)
}
