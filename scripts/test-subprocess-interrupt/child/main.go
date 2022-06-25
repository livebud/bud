package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/livebud/bud/internal/sig"
)

func main() {
	fmt.Println("child: started")
	ctx := sig.Trap(context.Background(), os.Interrupt)
	select {
	case <-ctx.Done():
		fmt.Println("child: interrupted!")
		time.Sleep(time.Second)
		fmt.Println("child: exiting")
	}
}
