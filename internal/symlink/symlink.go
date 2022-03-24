package symlink

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Link creates a symlink. It's different than os.Symlink in that it can
// override an existing symlink. To do this correctly, you need to create a
// temporary symlink then move it over the old symlink.
func Link(from, to string) error {
	tmpPath := tmp(from)
	if err := os.Symlink(from, tmpPath); err != nil {
		return err
	}
	return os.Rename(tmpPath, to)
}

func tmp(from string) string {
	return from + "." + randomString(6) + ".tmp"
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
