package colontree_test

import (
	"bytes"
	"strings"
	"testing"
	"text/tabwriter"

	"github.com/livebud/bud/internal/colontree"
	"github.com/livebud/bud/internal/is"
)

func heroku() *colontree.Node {
	ct := colontree.New("", "")
	ct.Insert("addons", "lists your add-ons and attachments")
	ct.Insert("addons:attach", "attaches an add-on to an app")
	ct.Insert("addons:create", "creates an add-on")
	ct.Insert("addons:destroy", "removes an add-on from an app")
	ct.Insert("addons:detach", "detaches an add-on from an app")
	ct.Insert("addons:open", "opens an add-on in your browser")
	ct.Insert("addons:rename", "renames an add-on")
	ct.Insert("addons:upgrade", "upgrades an add-on")
	ct.Insert("addons:downgrade", "downgrades an add-on")
	ct.Insert("addons:info", "displays information about an add-on")
	ct.Insert("addons:docs", "opens the documentation for an add-on")
	ct.Insert("addons:plans", "lists the plans for an add-on")
	ct.Insert("addons:wait", "waits for an add-on to be provisioned")
	ct.Insert("addons:services", "lists the services for an add-on")
	ct.Insert("ps", "lists dynos for an app")
	ct.Insert("ps:scale", "changes dyno quantity and size")
	ct.Insert("ps:restart", "restarts dynos")
	ct.Insert("ps:stop", "stops dynos")
	ct.Insert("ps:wait", "waits for dynos to restart")
	ct.Insert("ps:exec", "opens a shell on a dyno")
	ct.Insert("ps:kill", "kills a dyno")
	ct.Insert("ps:resize", "changes dyno size")
	ct.Insert("ps:autoscale", "changes dyno quantity and size based on dyno metrics")
	ct.Insert("ps:autoscale:disable", "disables dyno autoscaling")
	ct.Insert("ps:autoscale:enable", "enables dyno autoscaling")
	return ct
}

var reset = "\033[0m"
var dim = "\033[37m"

func usage(children []*colontree.Node) (string, error) {
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	for _, child := range children {
		tw.Write([]byte(child.Full()))
		if child.Value() != "" {
			tw.Write([]byte("\t" + dim + child.Value() + reset))
		}
		tw.Write([]byte("\n"))
	}
	if err := tw.Flush(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func TestChildren(t *testing.T) {
	is := is.New(t)
	ct := heroku()
	children := ct.Children()
	str, err := usage(children)
	is.NoErr(err)
	is.In(str, "addons")
	is.In(str, "lists your add-ons and attachments")
	is.In(str, "ps")
	is.In(str, "lists dynos for an app")
	is.NotIn(str, "ps:autoscale")
	is.NotIn(str, "addons:attach")

	children = ct.Find("ps").Children()
	str, err = usage(children)
	is.NoErr(err)
	is.In(str, "ps:autoscale")
	is.In(str, "changes dyno quantity and size based on dyno metrics")
	is.In(str, "ps:exec")
	is.In(str, "opens a shell on a dyno")
	is.In(str, "ps:kill")
	is.In(str, "kills a dyno")
	is.In(str, "ps:resize")
	is.In(str, "changes dyno size")
	is.In(str, "ps:restart")
	is.In(str, "restarts dynos")
	is.In(str, "ps:scale")
	is.In(str, "changes dyno quantity and size")
	is.In(str, "ps:stop")
	is.In(str, "stops dynos")
	is.In(str, "ps:wait")
	is.In(str, "waits for dynos to restart")
	is.NotIn(str, "ps:autoscale:disable")
	is.NotIn(str, "addons")

	children = ct.Find("ps:autoscale").Children()
	str, err = usage(children)
	is.NoErr(err)
	is.In(str, "ps:autoscale:disable")
	is.In(str, "disables dyno autoscaling")
	is.In(str, "ps:autoscale:enable")
	is.In(str, "enables dyno autoscaling")
	is.NotIn(str, "ps:wait")
	is.NotIn(str, "addons")
}

func TestNotFound(t *testing.T) {
	is := is.New(t)
	ct := heroku()
	node := ct.Find("notfound")
	is.Equal(node, nil)
}
