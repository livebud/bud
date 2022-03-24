package dirhash

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"

	"github.com/cespare/xxhash"
)

type skipFunc = func(name string, isDir bool) bool

type option struct {
	Skip skipFunc
}

type Option func(o *option)

// Provide a skip function
//
// Note: try to skip as high up in the tree as possible.
// E.g. if the source filesystem doesn't have bud, it will
// delete bud, even if you're skipping bud/generate.
func WithSkip(skips ...skipFunc) Option {
	return func(o *option) {
		o.Skip = composeSkips(skips)
	}
}

func composeSkips(skips []skipFunc) skipFunc {
	return func(name string, isDir bool) bool {
		for _, skip := range skips {
			if skip(name, isDir) {
				return true
			}
		}
		return false
	}
}

// Hash a filesystem. Based on the sumdb/dirhash.
func Hash(fsys fs.FS, options ...Option) (string, error) {
	opt := &option{
		Skip: func(string, bool) bool { return false },
	}
	for _, option := range options {
		option(opt)
	}
	h := xxhash.New()
	err := fs.WalkDir(fsys, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if opt.Skip(path, de.IsDir()) {
			return fs.SkipDir
		} else if de.IsDir() {
			return nil
		}
		f, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		hf := xxhash.New()
		_, err = io.Copy(hf, f)
		if err != nil {
			return err
		}
		fmt.Fprintf(h, "%x  %s\n", hf.Sum(nil), path)
		return nil
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil)), nil
}
