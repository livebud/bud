package snapshot

import (
	"encoding/base64"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/dirhash"
	"github.com/livebud/bud/internal/targz"

	"github.com/cespare/xxhash"
)

func cachePath(key string) (string, error) {
	dir := filepath.Join(os.TempDir(), "bud", "snapshot")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, key+".tar.gz"), nil
}

func hash(input string) string {
	hash := xxhash.New()
	hash.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

// Hash a filesystem
func Hash(fsys fs.FS) (string, error) {
	return dirhash.Hash(fsys)
}

// Backup a filesystem
func Backup(hash string, fsys fs.FS) error {
	path, err := cachePath(hash)
	if err != nil {
		return err
	}
	data, err := targz.Zip(fsys)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Restore a filesystem from a hash
func Restore(hash string) (fs.FS, error) {
	path, err := cachePath(hash)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return targz.Unzip(data)
}
