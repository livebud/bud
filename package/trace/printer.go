package trace

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

func Printer(ctx context.Context, w io.Writer) (c context.Context, print func()) {
	tracer, ctx, err := Serve(ctx)
	if err != nil {
		return ctx, func() {
			fmt.Fprintf(w, "\n%s", err)
		}
	}
	return ctx, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		trace, err := tracer.Print(ctx)
		if err != nil {
			fmt.Fprintf(w, "\n%s", err)
			return
		}
		fmt.Fprintf(os.Stderr, "\n%s", trace)
		if err := tracer.Shutdown(ctx); err != nil {
			fmt.Fprintf(w, "\n%s", err)
			return
		}
	}
}
