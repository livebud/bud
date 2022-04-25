package version

import (
	"context"
	"os"
	"text/tabwriter"

	"github.com/livebud/bud/internal/version"
)

type Command struct {
}

func (c *Command) Run(ctx context.Context) error {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	tw.Write([]byte("bud: \t" + version.Bud + "\n"))
	tw.Write([]byte("svelte: \t" + version.Svelte + "\n"))
	tw.Write([]byte("react: \t" + version.React + "\n"))
	return tw.Flush()
}
