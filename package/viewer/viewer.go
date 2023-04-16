package viewer

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"

	"github.com/livebud/bud/framework/controller/controllerrt/request"
	"github.com/livebud/bud/internal/errs"

	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/virtual"
)

var ErrViewerNotFound = errors.New("viewer not found")
var ErrPageNotFound = errors.New("page not found")

type Key = string
type Ext = string
type PropMap = map[Key]interface{}

// FS can be used to either use the real filesystem or an embedded one,
// depending on how Bud was built.
type FS = fs.FS

// Interface for viewers
type Interface interface {
	Mount(r *router.Router) error
	Render(ctx context.Context, key string, propMap PropMap) ([]byte, error)
	RenderError(ctx context.Context, key string, propMap PropMap) []byte
}

type View struct {
	Key    Key
	Path   string
	Ext    string
	Client *Client // View client
}

type Page struct {
	*View // Entry
	// Frames are the views that are rendered inside the layout. Frames start with
	// the innermost view first and end with the outermost view.
	Frames []*View
	Layout *View
	Error  *View
	Route  string
	Client *Client // Entry client
}

type Client struct {
	Path  string
	Route string
}

type Embed = virtual.File
type Embeds = map[string]*Embed
type Pages map[Key]*Page

type Viewer interface {
	Mount(r *router.Router) error
	Render(ctx context.Context, key string, propMap PropMap) ([]byte, error)
	RenderError(ctx context.Context, key string, propMap PropMap, err error) []byte
	Bundle(ctx context.Context, embed virtual.Tree) error
}

// StaticPropMap returns a prop map for static views based on the request data.
func StaticPropMap(page *Page, r *http.Request) (PropMap, error) {
	props := map[string]interface{}{}
	if err := request.Unmarshal(r, &props); err != nil {
		return nil, err
	}
	propMap := PropMap{}
	propMap[page.Key] = props
	if page.Layout != nil {
		propMap[page.Layout.Key] = props
	}
	for _, frame := range page.Frames {
		propMap[frame.Key] = props
	}
	if page.Error != nil {
		propMap[page.Error.Key] = props
	}
	return propMap, nil
}

func Error(err error) error {
	ve := &viewerError{
		original: err,
		Message:  err.Error(),
	}
	// TODO: add a stack trace, if we have one
	if errs, ok := err.(errs.Errors); ok {
		for _, err := range errs.Errors() {
			ve.Errors = append(ve.Errors, Error(err))
		}
	}
	return ve
}

type viewerError struct {
	original error
	Message  string
	Stack    []*StackFrame
	Errors   []error
}

func (v *viewerError) Error() string {
	return v.Message
}

func (v *viewerError) MarshalJSON() ([]byte, error) {
	if marshaler, ok := v.original.(json.Marshaler); ok {
		return marshaler.MarshalJSON()
	}
	return json.Marshal(map[string]interface{}{
		"message": v.Message,
		"stack":   v.Stack,
		"errors":  v.Errors,
	})
}

type StackFrame struct {
	Path   string
	Line   int
	Column int
}
