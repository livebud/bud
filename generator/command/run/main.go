package run

import "gitlab.com/mnm/bud/bfs"

func Generator() bfs.Generator {
	return bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
		return nil
	})
}
