package cli

import (
	"context"
	"strings"

	"github.com/livebud/bud/package/budhttp"
)

type ToolBSPublish struct {
	DialDev string
	Topic   string
	Data    []string
}

func (c *CLI) ToolBSPublish(ctx context.Context, in *ToolBSPublish) error {
	log, err := c.loadLog()
	if err != nil {
		return err
	}
	client, err := budhttp.Load(log, in.DialDev)
	if err != nil {
		return err
	}
	data := strings.Join(in.Data, " ")
	if err := client.Publish(in.Topic, []byte(data)); err != nil {
		return err
	}
	log.Info("published %q with data %s", in.Topic, data)
	return nil
}
