package scan

import (
	"errors"
	"io/fs"
	"path"
)

// Dir scans a directory
func Dir(fs fs.FS, validFn func(de fs.DirEntry) bool) Scanner {
	s := &dirScanner{
		fs:      fs,
		validFn: validFn,
		textCh:  make(chan string),
	}
	go s.walk(".")
	return s
}

type dirScanner struct {
	fs      fs.FS
	validFn func(de fs.DirEntry) bool
	textCh  chan string

	// The following start as empty values
	text string
	err  error
}

func (s *dirScanner) walk(dir string) {
	s.walkAll(dir)
	close(s.textCh)
}

func (s *dirScanner) walkAll(dir string) {
	des, err := fs.ReadDir(s.fs, dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return
		}
		s.err = err
		return
	}
	// Check if the directory contains a valid file first
	for _, de := range des {
		if de.IsDir() || !s.validFn(de) {
			continue
		}
		s.textCh <- dir
		break
	}
	// Walk the subdirectories
	for _, de := range des {
		if !de.IsDir() || !s.validFn(de) {
			continue
		}
		s.walkAll(path.Join(dir, de.Name()))
	}
}

func (s *dirScanner) Scan() bool {
	if text, more := <-s.textCh; more {
		s.text = text
		return true
	}
	return false
}

func (s *dirScanner) Err() error {
	return s.err
}

func (s *dirScanner) Text() string {
	return s.text
}
