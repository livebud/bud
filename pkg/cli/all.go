package cli

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// All returns a subcommand that runs all the given subcommands in parallel
func All(subs ...Subcommand) Subcommand {
	return (all)(subs)
}

type all []Subcommand

func (subs all) Usage(cmd Command) {
	for _, sub := range subs {
		sub.Usage(cmd)
	}
}

func (subs all) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, sub := range subs {
		sub := sub
		eg.Go(func() error {
			return sub.Run(ctx)
		})
	}
	return eg.Wait()
}
