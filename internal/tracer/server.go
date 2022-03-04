package tracer

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xlab/treeprint"
	"gitlab.com/mnm/bud/pkg/router"

	"gitlab.com/mnm/bud/pkg/socket"
)

// TODO: align with OpenTelemetry's data structure
type SpanData struct {
	ID       string
	ParentID string
	Name     string
	Attrs    map[string]string
	Start    int64
	End      int64
	Error    string
}

func (s *SpanData) Duration() time.Duration {
	startTime := time.Unix(0, s.Start)
	endTime := time.Unix(0, s.End)
	return endTime.Sub(startTime)
}

type SpanFields []*SpanField

func (f SpanFields) Len() int {
	return len(f)
}

func (f SpanFields) Less(i, j int) bool {
	return f[i].Key < f[j].Key
}

func (f SpanFields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type SpanField struct {
	Key   string
	Value string
}

func (s *SpanData) Fields() (fields SpanFields) {
	for key, val := range s.Attrs {
		fields = append(fields, &SpanField{key, val})
	}
	sort.Sort(fields)
	return fields
}

func (s *SpanData) String() string {
	out := new(strings.Builder)
	out.WriteString(s.Name)
	dur := s.Duration()
	if dur > 0 {
		out.WriteString(" (" + dur.String() + ")")
	}
	if s.Error != "" {
		out.WriteString(" error=" + strconv.Quote(s.Error))
	}
	for _, field := range s.Fields() {
		out.WriteString(" " + field.Key + "=" + field.Value)
	}
	return out.String()
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
		traces: map[string][]*SpanData{},
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
	traces map[string][]*SpanData
}

func (a *api) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var spans []*SpanData
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
	tree := treeprint.NewWithRoot(root.String())
	children := a.traces[root.ID]
	sortByStart(children)
	for _, child := range children {
		a.print(tree, child)
	}
	w.Write([]byte(tree.String()))
}

func (a *api) print(tree treeprint.Tree, span *SpanData) {
	tree = tree.AddBranch(span.String())
	children := a.traces[span.ID]
	sortByStart(children)
	for _, child := range children {
		a.print(tree, child)
	}
}

func sortByStart(spans []*SpanData) {
	sort.Slice(spans, func(i, j int) bool {
		return spans[i].Start < spans[j].Start
	})
}
