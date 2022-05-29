package version

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/livebud/bud/internal/version"
)

type Command struct {
	Key string
}

func (c *Command) Run(ctx context.Context) error {
	switch c.Key {
	case "bud":
		fmt.Println(version.Bud)
		return nil
	case "svelte":
		fmt.Println(version.Svelte)
		return nil
	case "react":
		fmt.Println(version.React)
		return nil
	default:
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
		tw.Write([]byte("bud: \t" + version.Bud + "\n"))
		tw.Write([]byte("svelte: \t" + version.Svelte + "\n"))
		tw.Write([]byte("react: \t" + version.React + "\n"))
		return tw.Flush()
	}
}
