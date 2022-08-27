package hot

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/livebud/bud/internal/errs"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/package/log"
)

// Dial creates a server-sent event (SSE) stream. This stream has been adapted
// from the following minimal eventsource stream:
// - https://github.com/neelance/eventsource/blob/master/client/client.go
// Thanks Richard!
func Dial(log log.Interface, url string) (*Stream, error) {
	return DialWith(http.DefaultClient, log, url)
}

// DialWith creates a server-sent event (SSE) stream with a custom HTTP client.
func DialWith(client *http.Client, log log.Interface, url string) (*Stream, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Close = true
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	stream := &Stream{
		log:     log,
		res:     res,
		eventCh: make(chan *Event, 1),
		errorCh: make(chan error),
		closeCh: make(chan struct{}),
	}
	go stream.loop()
	return stream, nil
}

type Stream struct {
	log     log.Interface
	res     *http.Response
	eventCh chan *Event
	errorCh chan error
	closeCh chan struct{}
	once    once.Error
}

func (s *Stream) loop() {
	sc := bufio.NewScanner(s.res.Body)
	sc.Split(scanLines)
	event := &Event{}
	data := [][]byte{}
	// Return the error
	defer func() { s.errorCh <- sc.Err() }()
	// Scan line by line
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			event.Data = bytes.Join(data, []byte{'\n'})
			// Don't let pending events block the client from closing
			select {
			case s.eventCh <- event:
			case <-s.closeCh:
				break
			}
			event = &Event{}
			data = [][]byte{}
		}
		key, value := parseLine(line)
		switch string(key) {
		case "event":
			event.Type = string(value)
		case "data":
			data = append(data, value)
		case "id":
			event.ID = string(value)
		case "retry":
			// TODO
		default:
			// ignored
		}
	}
}

func (s *Stream) Next(ctx context.Context) (*Event, error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case evt := <-s.eventCh:
			return evt, nil
		case err := <-s.errorCh:
			return nil, err
		case <-ticker.C:
			s.log.Debug("hot: client waiting for next event")
		}
	}
}

func (s *Stream) Close() error {
	return s.once.Do(s.close)
}

func (s *Stream) close() (err error) {
	err = errs.Join(err, s.res.Body.Close())
	close(s.closeCh)
	// Drain event channel
	if e := <-s.errorCh; e != nil {
		// Closed errors are expected since we closed the body
		if !isExpectedCloseError(e) {
			err = errs.Join(err, e)
		}
	}
	close(s.errorCh)
	close(s.eventCh)
	return err
}

// isExpectedCloseError returns true if the close error is expected
func isExpectedCloseError(err error) bool {
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	// Unfortunately this error is not exported
	// https://github.com/golang/go/blob/f4274e64aac99aaa9af05988f2f8c36c47554889/src/net/http/transport.go#L2636
	if err.Error() == "http: read on closed response body" {
		return true
	}
	return false
}

func parseLine(line []byte) (key []byte, value []byte) {
	key = line
	colon := bytes.IndexByte(line, ':')
	if colon == -1 {
		return nil, nil
	}
	// Handle comments
	if colon == 0 {
		key = []byte("comment")
		value = line[colon+1:]
		return key, value
	}
	// Parse into key-value
	key = line[:colon]
	value = line[colon+1:]
	if value[0] == ' ' {
		value = value[1:]
	}
	return key, value
}

// Scan each line of data
func scanLines(data []byte, eof bool) (advance int, token []byte, err error) {
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		if data[i] == '\r' {
			if i == len(data)-1 {
				if eof {
					// final line
					return len(data), data[:len(data)-1], io.EOF
				}
				return 0, nil, nil // LF may follow, request more data
			}
			if data[i+1] == '\n' {
				return i + 2, data[:i], nil
			}
			return i + 1, data[:i], nil
		}
		// data[i] == '\n'
		return i + 1, data[:i], nil
	}
	if eof {
		// final line
		return len(data), data, io.EOF
	}
	// request more data
	return 0, nil, nil
}
