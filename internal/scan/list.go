package scan

import (
	"io/fs"
)

// List returns all the matches within a directory.
func List(fs fs.FS, dir string, matcher func(de fs.DirEntry) bool) (list []string, err error) {
	scanner := Dir(fs, dir, matcher)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		list = append(list, scanner.Text())
	}
	return list, scanner.Err()
}
