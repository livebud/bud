package version

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/livebud/bud/internal/versions"
)

func New() *Command {
	return &Command{}
}

type Command struct {
	Key string
}

func (c *Command) Run(ctx context.Context) error {
	switch c.Key {
	case "bud":
		fmt.Println(versions.Bud)
		return nil
	case "svelte":
		fmt.Println(versions.Svelte)
		return nil
	case "react":
		fmt.Println(versions.React)
		return nil
	default:
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
		tw.Write([]byte("bud: \t" + versions.Bud + "\n"))
		tw.Write([]byte("svelte: \t" + versions.Svelte + "\n"))
		tw.Write([]byte("react: \t" + versions.React + "\n"))
		return tw.Flush()
	}
}
