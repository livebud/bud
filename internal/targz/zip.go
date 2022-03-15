package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/fs"
)

func Zip(fsys fs.FS) ([]byte, error) {
	b := new(bytes.Buffer)
	zw := gzip.NewWriter(b)
	tw := tar.NewWriter(zw)
	// Walk the directory
	err := fs.WalkDir(fsys, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fi, err := de.Info()
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fi, path)
		if err != nil {
			return err
		}
		// Override to use the full path name
		header.Name = path
		if err = tw.WriteHeader(header); err != nil {
			return err
		}
		// Don't try reading from directories or symlinks
		if de.IsDir() || fi.Mode()&fs.ModeSymlink != 0 {
			return nil
		}
		file, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Close the tar writer
	if err := tw.Close(); err != nil {
		return nil, err
	}
	// Close the gzip writer
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
