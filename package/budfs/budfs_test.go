package budfs_test

import (
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/log"
)

type tailwind struct {
}

func (t *tailwind) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	file.Data = []byte("/* tailwind */")
	return nil
}

type svelte struct {
}

func (s *svelte) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	file.Data = []byte("/* svelte */")
	return nil
}

func TestFS(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{}
	log := log.Discard
	budfs := budfs.New(fsys, log)
	budfs.FileGenerator("bud/public/tailwind/tailwind.css", &tailwind{})
	budfs.FileGenerator("bud/view/index.svelte", &svelte{})
	is.NoErr(fstest.TestFS(budfs, "bud/public/tailwind/tailwind.css", "bud/view/index.svelte"))
}
