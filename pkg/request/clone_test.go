package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"sync"
	"testing"

	"github.com/livebud/bud/pkg/request"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func TestCloneRequest(t *testing.T) {
	is := is.New(t)
	r1 := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewBufferString(`{"foo":"bar"}`))
	r1.Header.Set("Content-Type", "application/json")
	r1.Header.Set("X-Test", "R1")
	r2 := request.Clone(r1)
	r2.Header.Set("X-Test", "R2")
	r3 := request.Clone(r2)
	r3.Header.Set("X-Test", "R3")
	o1, err := httputil.DumpRequestOut(r1, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(o1), `
		POST / HTTP/1.1
		Host: example.com
		User-Agent: Go-http-client/1.1
		Content-Length: 13
		Content-Type: application/json
		X-Test: R1
		Accept-Encoding: gzip

		{"foo":"bar"}
	`)
	o2, err := httputil.DumpRequestOut(r2, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(o2), `
		POST / HTTP/1.1
		Host: example.com
		User-Agent: Go-http-client/1.1
		Content-Length: 13
		Content-Type: application/json
		X-Test: R2
		Accept-Encoding: gzip

		{"foo":"bar"}
	`)
	o3, err := httputil.DumpRequestOut(r3, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(o3), `
		POST / HTTP/1.1
		Host: example.com
		User-Agent: Go-http-client/1.1
		Content-Length: 13
		Content-Type: application/json
		X-Test: R3
		Accept-Encoding: gzip

		{"foo":"bar"}
	`)
}

func TestAsyncCloneRequest(t *testing.T) {
	is := is.New(t)
	r1 := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewBufferString(`{"foo":"bar"}`))
	r1.Header.Set("Content-Type", "application/json")
	r1.Header.Set("X-Test", "R1")
	r2 := request.Clone(r1)
	r2.Header.Set("X-Test", "R2")
	r3 := request.Clone(r2)
	r3.Header.Set("X-Test", "R3")
	r4 := request.Clone(r3)
	r4.Header.Set("X-Test", "R4")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		o1, err := httputil.DumpRequestOut(r1, true)
		is.NoErr(err)
		diff.TestHTTP(t, string(o1), `
			POST / HTTP/1.1
			Host: example.com
			User-Agent: Go-http-client/1.1
			Content-Length: 13
			Content-Type: application/json
			X-Test: R1
			Accept-Encoding: gzip

			{"foo":"bar"}
		`)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		o2, err := httputil.DumpRequestOut(r2, true)
		is.NoErr(err)
		diff.TestHTTP(t, string(o2), `
			POST / HTTP/1.1
			Host: example.com
			User-Agent: Go-http-client/1.1
			Content-Length: 13
			Content-Type: application/json
			X-Test: R2
			Accept-Encoding: gzip

			{"foo":"bar"}
		`)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		o3, err := httputil.DumpRequestOut(r3, true)
		is.NoErr(err)
		diff.TestHTTP(t, string(o3), `
			POST / HTTP/1.1
			Host: example.com
			User-Agent: Go-http-client/1.1
			Content-Length: 13
			Content-Type: application/json
			X-Test: R3
			Accept-Encoding: gzip

			{"foo":"bar"}
		`)
	}()
	wg.Wait()
	// Laggard
	o4, err := httputil.DumpRequestOut(r4, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(o4), `
		POST / HTTP/1.1
		Host: example.com
		User-Agent: Go-http-client/1.1
		Content-Length: 13
		Content-Type: application/json
		X-Test: R4
		Accept-Encoding: gzip

		{"foo":"bar"}
	`)
}

const sherlock = `
To Sherlock Holmes she is always THE woman. I have seldom heard
him mention her under any other name. In his eyes she eclipses
and predominates the whole of her sex. It was not that he felt
any emotion akin to love for Irene Adler. All emotions, and that
one particularly, were abhorrent to his cold, precise but
admirably balanced mind. He was, I take it, the most perfect
reasoning and observing machine that the world has seen, but as a
lover he would have placed himself in a false position. He never
spoke of the softer passions, save with a gibe and a sneer. They
were admirable things for the observer--excellent for drawing the
veil from men's motives and actions. But for the trained reasoner
to admit such intrusions into his own delicate and finely
adjusted temperament was to introduce a distracting factor which
might throw a doubt upon all his mental results. Grit in a
sensitive instrument, or a crack in one of his own high-power
lenses, would not be more disturbing than a strong emotion in a
nature such as his. And yet there was but one woman to him, and
that woman was the late Irene Adler, of dubious and questionable
memory.

I had seen little of Holmes lately. My marriage had drifted us
away from each other. My own complete happiness, and the
home-centred interests which rise up around the man who first
finds himself master of his own establishment, were sufficient to
absorb all my attention, while Holmes, who loathed every form of
society with his whole Bohemian soul, remained in our lodgings in
Baker Street, buried among his old books, and alternating from
week to week between cocaine and ambition, the drowsiness of the
drug, and the fierce energy of his own keen nature. He was still,
as ever, deeply attracted by the study of crime, and occupied his
immense faculties and extraordinary powers of observation in
following out those clues, and clearing up those mysteries which
had been abandoned as hopeless by the official police. From time
to time I heard some vague account of his doings: of his summons
to Odessa in the case of the Trepoff murder, of his clearing up
of the singular tragedy of the Atkinson brothers at Trincomalee,
and finally of the mission which he had accomplished so
delicately and successfully for the reigning family of Holland.
Beyond these signs of his activity, however, which I merely
shared with all the readers of the daily press, I knew little of
my former friend and companion.
`

const expect1 = `
POST / HTTP/1.1
Host: example.com
User-Agent: Go-http-client/1.1
Content-Length: 2450
Content-Type: text/plaintext
X-Test: R1
Accept-Encoding: gzip

` + sherlock

const expect2 = `
POST / HTTP/1.1
Host: example.com
User-Agent: Go-http-client/1.1
Content-Length: 2450
Content-Type: text/plaintext
X-Test: R2
Accept-Encoding: gzip

` + sherlock

const expect3 = `
POST / HTTP/1.1
Host: example.com
User-Agent: Go-http-client/1.1
Content-Length: 2450
Content-Type: text/plaintext
X-Test: R3
Accept-Encoding: gzip

` + sherlock

const expect4 = `
POST / HTTP/1.1
Host: example.com
User-Agent: Go-http-client/1.1
Content-Length: 2450
Accepts: text/plaintext
Content-Type: text/plaintext
X-Test: R4
Accept-Encoding: gzip

` + sherlock

func TestLargeCloneRequest(t *testing.T) {
	is := is.New(t)
	r1 := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewBufferString(sherlock))
	r1.Header.Set("Content-Type", "text/plaintext")
	r1.Header.Set("X-Test", "R1")
	r2 := request.Clone(r1)
	r2.Header.Set("X-Test", "R2")
	r3 := request.Clone(r2)
	r3.Header.Set("X-Test", "R3")
	r4 := request.Clone(r3)
	r4.Header.Set("X-Test", "R4")
	r4.Header.Set("Accepts", "text/plaintext")
	o1, err := httputil.DumpRequestOut(r1, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(o1), expect1)
	o2, err := httputil.DumpRequestOut(r2, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(o2), expect2)
	o3, err := httputil.DumpRequestOut(r3, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(o3), expect3)
	o4, err := httputil.DumpRequestOut(r4, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(o4), expect4)
}

func TestLargeAsyncCloneRequest(t *testing.T) {
	is := is.New(t)
	r1 := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewBufferString(sherlock))
	r1.Header.Set("Content-Type", "text/plaintext")
	r1.Header.Set("X-Test", "R1")
	r2 := request.Clone(r1)
	r2.Header.Set("X-Test", "R2")
	r3 := request.Clone(r2)
	r3.Header.Set("X-Test", "R3")
	r4 := request.Clone(r3)
	r4.Header.Set("X-Test", "R4")
	r4.Header.Set("Accepts", "text/plaintext")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		o1, err := httputil.DumpRequestOut(r1, true)
		is.NoErr(err)
		diff.TestHTTP(t, string(o1), expect1)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		o2, err := httputil.DumpRequestOut(r2, true)
		is.NoErr(err)
		diff.TestHTTP(t, string(o2), expect2)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		o3, err := httputil.DumpRequestOut(r3, true)
		is.NoErr(err)
		diff.TestHTTP(t, string(o3), expect3)
	}()
	wg.Wait()
	// Laggard
	o4, err := httputil.DumpRequestOut(r4, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(o4), expect4)
}
