package cli

import (
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/livebud/bud/internal/versions"
)

type Version struct {
	Key string
}

func (c *CLI) Version(ctx context.Context, in *Version) error {
	switch in.Key {
	case "bud":
		fmt.Fprintln(c.Stdout, versions.Bud)
		return nil
	case "svelte":
		fmt.Fprintln(c.Stdout, versions.Svelte)
		return nil
	case "react":
		fmt.Fprintln(c.Stdout, versions.React)
		return nil
	default:
		tw := tabwriter.NewWriter(c.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
		tw.Write([]byte("bud: \t" + versions.Bud + "\n"))
		tw.Write([]byte("svelte: \t" + versions.Svelte + "\n"))
		tw.Write([]byte("react: \t" + versions.React + "\n"))
		return tw.Flush()
	}
}
