package is

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

type mockT struct {
	failed bool
}

var _ T = (*mockT)(nil)

func (t *mockT) FailNow() {
	t.failed = true
}

func TestEqual(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = true
	is.Equal(1, 2)
	if !m.failed {
		t.Fatalf("expected is.Equal to fail")
	}
	expect := "\x1b[90mis_test.go:26: \x1b[39m1 != 2\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestEqualNoColor(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.Equal(1, 2)
	if !m.failed {
		t.Fatalf("expected is.Equal to fail")
	}
	expect := "is_test.go:42: 1 != 2\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestEqualIgnoreComment(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.Equal(1, 2) // Ignore comment
	if !m.failed {
		t.Fatalf("expected is.Equal to fail")
	}
	expect := "is_test.go:58: 1 != 2\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestEqualMessageArgsColor(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = true
	is.Equal(1, 2, "%d does not equal %d", 1, 2)
	if !m.failed {
		t.Fatalf("expected is.Equal to fail")
	}
	expect := "\x1b[90mis_test.go:74: \x1b[39m1 != 2. \x1b[31m1 does not equal 2\x1b[39m\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestEqualMessageArgs(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.Equal(1, 2, "%d does not equal %d", 1, 2)
	if !m.failed {
		t.Fatalf("expected is.Equal to fail")
	}
	expect := "is_test.go:90: 1 != 2. 1 does not equal 2\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestEqualMessageOnly(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.Equal(1, 2, "incorrect value")
	if !m.failed {
		t.Fatalf("expected is.Equal to fail")
	}
	expect := "is_test.go:106: 1 != 2. incorrect value\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestEqualMessageOnlyIgnoreFormatting(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.Equal(1, 2, "incorrect value %s")
	if !m.failed {
		t.Fatalf("expected is.Equal to fail")
	}
	expect := "is_test.go:122: 1 != 2. incorrect value %s\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestNoErr(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = true
	is.NoErr(errors.New("some error"))
	if !m.failed {
		t.Fatalf(`expected is.NoErr to fail`)
	}
	expect := "\x1b[90mis_test.go:138: \x1b[39merr: some error\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestNoErrFormatf(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.NoErr(fmt.Errorf("some error: %w", errors.New("inner error")))
	if !m.failed {
		t.Fatalf(`expected is.NoErr(errors.New("some error")) to fail`)
	}
	expect := "is_test.go:154: err: some error: inner error\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestNoErrMessageArgs(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.NoErr(fmt.Errorf("some error: %w", errors.New("inner error")), "details: %s", "deets")
	if !m.failed {
		t.Fatalf(`expected is.NoErr to fail`)
	}
	expect := "is_test.go:170: err: some error: inner error. details: deets\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestEqualOkInt(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = true
	is.Equal(1, 1)
	if m.failed {
		t.Fatalf(`expected no failure`)
	}
	if out.String() != "" {
		t.Fatalf("expected no buffer")
	}
}

func TestEqualOkString(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = true
	is.Equal("a", "a")
	if m.failed {
		t.Fatalf(`expected no failure`)
	}
	if out.String() != "" {
		t.Fatalf("expected no buffer")
	}
}

func TestNoErrOk(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = true
	is.NoErr(nil)
	if m.failed {
		t.Fatalf(`expected no failure`)
	}
	if out.String() != "" {
		t.Fatalf("expected no buffer")
	}
}

func one(is *I, a, b int) {
	is.Helper()
	two(is, a, b)
}

func two(is *I, a, b int) {
	is.Helper()
	is.Equal(a, b)
}

func TestHelper(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	one(is, 1, 2)
	if !m.failed {
		t.Fatalf("expected is.Equal to fail")
	}
	expect := "is_test.go:241: 1 != 2\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestInOk(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = true
	is.In("hello", "lo")
	if m.failed {
		t.Fatalf(`expected no failure`)
	}
	if out.String() != "" {
		t.Fatalf("expected no buffer")
	}
}

func TestInNotList(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.In(3, 3)
	if !m.failed {
		t.Fatalf(`expected is.In to fail`)
	}
	expect := "is_test.go:272: int(3) is not a list\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestIn(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.In("hello", "hi")
	if !m.failed {
		t.Fatalf(`expected is.In to fail`)
	}
	expect := "is_test.go:288: hi not in hello\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestInSliceOk(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.In([]string{"a", "b", "c"}, "c")
	if m.failed {
		t.Fatalf(`expected no failure`)
	}
	if out.String() != "" {
		t.Fatalf("expected no buffer")
	}
}

func TestInSlice(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.In([]string{"a", "b", "c"}, "d")
	if !m.failed {
		t.Fatalf(`expected is.In to fail`)
	}
	expect := "is_test.go:319: d not in [a b c]\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestTrue(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.True(true && false)
	if !m.failed {
		t.Fatalf(`expected is.True to fail`)
	}
	expect := "is_test.go:335: not true: true && false\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestTrueOk(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.True(1 == 1)
	if m.failed {
		t.Fatalf(`expected no failure`)
	}
	if out.String() != "" {
		t.Fatalf("expected no buffer")
	}
}

func TestTrueMessage(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.True(1 != 1, "expected %d and %d to be equal", 1, 1)
	if !m.failed {
		t.Fatalf(`expected is.True to fail`)
	}
	expect := "is_test.go:366: not true: 1 != 1. expected 1 and 1 to be equal\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestInMessage(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.In("hello", "hi", "%v doesn't contain %v", "hello", "hi")
	if !m.failed {
		t.Fatalf(`expected is.In to fail`)
	}
	expect := "is_test.go:382: hi not in hello. hello doesn't contain hi\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestFail(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.Fail()
	if !m.failed {
		t.Fatalf(`expected is.Fail to fail`)
	}
	expect := "is_test.go:398: failed\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}

func TestFailMessage(t *testing.T) {
	m := &mockT{}
	is := New(m)
	out := new(bytes.Buffer)
	is.writer = out
	is.colorful = false
	is.Fail("context from %q should have exited", "component")
	if !m.failed {
		t.Fatalf(`expected is.Fail to fail`)
	}
	expect := "is_test.go:414: failed. context from \"component\" should have exited\n"
	if out.String() != expect {
		t.Fatalf("expected %q, got %q", expect, out)
	}
}
