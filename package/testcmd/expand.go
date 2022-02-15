package testcmd

import (
	"bytes"
	"context"
	"os"

	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/expander"
)

// Run expand with test configuration
func Expand(dir string, args ...string) (stdout string, err error) {
	expander, err := expander.Load(dir)
	if err != nil {
		return "", err
	}
	expander.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + os.Getenv("PATH"),
		"GOPATH=" + os.Getenv("GOPATH"),
		"GOCACHE=" + os.Getenv("GOCACHE"),
		"GOMODCACHE=" + testdir.ModCache(dir).Directory(),
		"NO_COLOR=1",
		// TODO: remove once we can write a sum file to the modcache
		"GOPRIVATE=*",
	}
	out := new(bytes.Buffer)
	expander.Stdout = out
	expander.Stderr = os.Stderr
	ctx := context.Background()
	if err := expander.Expand(ctx, args...); err != nil {
		return "", err
	}
	return out.String(), nil
}
