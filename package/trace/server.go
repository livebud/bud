package trace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/xlab/treeprint"
	"gitlab.com/mnm/bud/pkg/router"

	"gitlab.com/mnm/bud/pkg/socket"
)

// TODO: align with OpenTelemetry's data structure
type spanData struct {
	ID       string
	Name     string
	ParentID string
	Duration string
}

func Serve(path string) (*http.Server, error) {
	server := newServer()
	listener, err := socket.Listen(path)
	if err != nil {
		return nil, err
	}
	go server.Serve(listener)
	return server, nil
}

func Handler() http.Handler {
	router := router.New()
	api := &api{
		Router: router,
		traces: map[string][]*spanData{},
	}
	router.Post("/", http.HandlerFunc(api.Create))
	router.Get("/", http.HandlerFunc(api.Index))
	return api
}

func newServer() *http.Server {
	return &http.Server{Handler: Handler()}
}

type api struct {
	*router.Router
	mu     sync.RWMutex
	traces map[string][]*spanData
}

func (a *api) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var spans []*spanData
	if err := json.Unmarshal(body, &spans); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, span := range spans {
		sid := span.ParentID
		a.traces[sid] = append(a.traces[sid], span)
	}
	w.WriteHeader(200)
}

func (a *api) Index(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	roots := a.traces["0000000000000000"]
	if len(roots) == 0 {
		w.Write([]byte(""))
		return
	}
	root := roots[0]
	tree := treeprint.NewWithRoot(fmt.Sprintf("%s (%s)", root.Name, root.Duration))
	children := a.traces[root.ID]
	for _, child := range children {
		a.print(tree, child)
	}
	w.Write([]byte(tree.String()))
}

func (a *api) print(tree treeprint.Tree, span *spanData) {
	tree = tree.AddBranch(fmt.Sprintf("%s (%s)", span.Name, span.Duration))
	children := a.traces[span.ID]
	for _, child := range children {
		a.print(tree, child)
	}
}
