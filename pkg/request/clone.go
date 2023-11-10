package request

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/livebud/bud/pkg/once"
)

func Clone(r *http.Request) *http.Request {
	mr, ok := r.Body.(*multiReader)
	if !ok {
		mr = &multiReader{
			&internalState{
				source: r.Body,
				buffer: new(bytes.Buffer),
			},
			0,
		}
		r.Body = mr
	}
	req := r.Clone(r.Context())
	req.Body = &multiReader{mr.state, 0}
	return req
}

type multiReader struct {
	state *internalState
	pos   int64
}

type internalState struct {
	once   once.Error
	mu     sync.Mutex
	source io.ReadCloser
	buffer *bytes.Buffer
}

func (mr *multiReader) Read(p []byte) (n int, err error) {
	mr.state.mu.Lock()
	defer mr.state.mu.Unlock()

	bytesNeeded := len(p)
	// Read from the buffer first
	bytesReadFromBuffer := copy(p, mr.state.buffer.Bytes()[mr.pos:])
	// Check out many bytes we will need to read from the source
	bytesRemaining := bytesNeeded - bytesReadFromBuffer
	if bytesRemaining <= 0 {
		n = bytesReadFromBuffer
		mr.pos += int64(bytesReadFromBuffer)
		return n, nil
	}
	// Copy remaining needed bytes from source into the buffer
	bytesReadFromSource, err := io.CopyN(mr.state.buffer, mr.state.source, int64(bytesRemaining))
	// Take the bytes read into the buffer and fill the remaining bytes in the output buffer
	copy(p[bytesReadFromBuffer:], mr.state.buffer.Bytes()[mr.pos+int64(bytesReadFromBuffer):])
	n = bytesReadFromBuffer + int(bytesReadFromSource)
	// Update the position of our reader
	mr.pos += int64(n)
	return n, err
}

func (mr *multiReader) Close() error {
	return mr.state.once.Do(func() error {
		return mr.state.source.Close()
	})
}
