// Package urlx is a parser for accepting more server addresses including unix
// domain sockets during development. Eventually, we'll want to switch off of
// peg for better error messages. For now, it's quicker to iterate on the
// grammar using a peg parser.
package urlx

import (
	"errors"
	"fmt"
	"net/url"
)

//go:generate go run github.com/pointlander/peg -switch -inline parse.peg

var ErrParsing = errors.New("urlx: unable to parse")

var defaultHost = "127.0.0.1"
var defaultPort = "3000"

func Parse(input string) (*url.URL, error) {
	parser := &parser{Buffer: input}
	parser.Init()
	err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("%w %q", ErrParsing, input)
	}
	parser.Execute()

	// Fallback to regular url parsing
	if parser.url.uri != "" {
		u, err := url.Parse(input)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	// Create a new url
	u := new(url.URL)

	// Handle the scheme
	if parser.url.scheme != "" {
		u.Scheme = parser.url.scheme
	} else if parser.url.port == "443" {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	hasHost := parser.url.host != ""
	hasPort := parser.url.port != ""
	hasPath := parser.url.path != ""

	// Handle the port
	port := defaultPort
	if !hasPort && u.Scheme == "https" {
		port = "443"
	}

	// Handle the host and port
	if hasHost && hasPort {
		u.Host = parser.url.host + ":" + parser.url.port
	} else if hasPort {
		u.Host = defaultHost + ":" + parser.url.port // default to 127.0.0.1
	} else if hasHost {
		u.Host = parser.url.host + ":" + port
	} else if !hasPath { // Only append default host:port if there's no path
		u.Host = defaultHost + ":" + port
	}

	// Handle the path if there is one
	u.Path = parser.url.path
	return u, nil
}

// Used in the parser
type uri struct {
	port   string
	scheme string
	host   string
	path   string
	uri    string
}
