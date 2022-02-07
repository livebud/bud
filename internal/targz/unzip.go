package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/fs"
	"testing/fstest"
)

func Unzip(b []byte) (fs.FS, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gr)
	fsys := fstest.MapFS{}
	for {
		header, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		fi := header.FileInfo()
		fsys[header.Name] = &fstest.MapFile{
			Mode:    fi.Mode(),
			ModTime: fi.ModTime(),
			Sys:     fi.Sys(),
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		fsys[header.Name].Data = data
	}
	// Close the gzip reader
	if err := gr.Close(); err != nil {
		return nil, err
	}
	return fsys, nil
}
