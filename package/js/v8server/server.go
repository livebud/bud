package v8server

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/livebud/bud/package/js"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/js/v8client"
)

func Serve() error {
	server := &Server{os.Stdin, os.Stdout}
	return server.Serve()
}

func New(reader io.ReadCloser, writer io.WriteCloser) *Server {
	return &Server{reader, writer}
}

type Server struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (s *Server) Serve() error {
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(s.reader)
	enc := gob.NewEncoder(s.writer)
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

func script(vm js.VM, enc *gob.Encoder, path, code string) error {
	var out v8client.Output
	err := vm.Script(path, code)
	if err != nil {
		out.Error = err.Error()
	}
	return enc.Encode(out)
}

func eval(vm js.VM, enc *gob.Encoder, path, code string) error {
	var out v8client.Output
	result, err := vm.Eval(path, code)
	if err != nil {
		out.Error = err.Error()
		return enc.Encode(out)
	}
	out.Result = result
	return enc.Encode(out)
}
