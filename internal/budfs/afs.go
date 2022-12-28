package budfs

import (
	"context"
	"fmt"

	"github.com/livebud/bud/package/virt"
)

type FS interface {
	Open(path string, file *virt.File) error
}

func Serve(ctx context.Context) (FS, error) {
	// files := extrafile.Load("BUD_REMOTEFS")
	// if len(files) > 0 {
	// log.Debug("afs: serving from BUD_REMOTEFS file listener passed in from the parent process")
	// return socket.From(files[0])
	// }
	return nil, fmt.Errorf("serve not implemented")
}
