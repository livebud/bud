package snapshot

import (
	"encoding/base64"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing/fstest"

	"gitlab.com/mnm/bud/internal/targz"

	"github.com/cespare/xxhash"
)

func cachePath(key string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cacheDir, "bud", "snapshot")
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
	mapfs := fstest.MapFS{}
	err := fs.WalkDir(fsys, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := de.Info()
		if err != nil {
			return err
		}
		mapfs[path] = &fstest.MapFile{
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			Sys:     info.Sys(),
		}
		if de.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		mapfs[path].Data = data
		return nil
	})
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(mapfs)
	if err != nil {
		return "", err
	}
	return hash(string(buf)), nil
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
