package v8server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/livebud/bud/package/js"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/js/v8client"
)

// empty
var emptyCloser = func() error { return nil }

// Pipe creates a V8 server from a file pipe
func Pipe() (*Server, error) {
	r1, w2, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	r2, w1, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	server := &Server{
		reader: r1,
		writer: w1,
		files:  []*os.File{r2, w2},
		closer: func() error {
			if w1.Close(); err != nil {
				return err
			}
			if w2.Close(); err != nil {
				return err
			}
			return err
		},
	}
	return server, nil
}

func New(reader io.ReadCloser, writer io.WriteCloser) *Server {
	return &Server{
		reader: reader,
		writer: writer,
		closer: emptyCloser,
	}
}

type Server struct {
	reader io.ReadCloser
	writer io.WriteCloser
	files  []*os.File
	closer func() error
}

func (s *Server) Serve() error {
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	dec := json.NewDecoder(s.reader)
	enc := json.NewEncoder(s.writer)
	for {
		// Decode messages into input
		var in v8client.Input
		if err := dec.Decode(&in); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			// Return an error avoid potential infinite loops even though it will kill
			// the V8 server.
			return fmt.Errorf("v8server: error decoding: %w", err)
		}
		// Handle eval
		if in.Type == "eval" {
			if err := eval(vm, enc, in.Path, in.Code); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		// Handle scripts
		if in.Type == "script" {
			if err := script(vm, enc, in.Path, in.Code); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
	}
}

// Files returns a list of files that can be used to connect to the server
func (s *Server) Files() []*os.File {
	return s.files
}

// Close the server
func (s *Server) Close() error {
	return s.closer()
}

func script(vm js.VM, enc *json.Encoder, path, code string) error {
	var out v8client.Output
	err := vm.Script(path, code)
	if err != nil {
		out.Error = err.Error()
	}
	return enc.Encode(out)
}

func eval(vm js.VM, enc *json.Encoder, path, code string) error {
	var out v8client.Output
	result, err := vm.Eval(path, code)
	if err != nil {
		out.Error = err.Error()
		return enc.Encode(out)
	}
	out.Result = result
	return enc.Encode(out)
}
