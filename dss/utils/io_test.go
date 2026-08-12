// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1) test
// vectors, cross-checked against known org.apache.commons.io.IOUtils
// semantics.

package utils

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type errCloser struct{ closed bool }

func (e *errCloser) Close() error {
	e.closed = true
	return errors.New("boom")
}

func TestToByteArray(t *testing.T) {
	got, err := ToByteArray(strings.NewReader("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("got %q", got)
	}
}

func TestCloseQuietly(t *testing.T) {
	CloseQuietly(nil) // must not panic

	c := &errCloser{}
	CloseQuietly(c) // must swallow the error
	if !c.closed {
		t.Error("expected Close to have been called")
	}
}

func TestCopy(t *testing.T) {
	var buf bytes.Buffer
	if err := Copy(strings.NewReader("hello"), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "hello" {
		t.Errorf("got %q", buf.String())
	}
}

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	if err := Write([]byte("hello"), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "hello" {
		t.Errorf("got %q", buf.String())
	}
}

func TestNullWriter(t *testing.T) {
	n, err := NullWriter().Write([]byte("discard me"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len("discard me") {
		t.Errorf("got %d", n)
	}
}

func TestGetInputStreamSize(t *testing.T) {
	size, err := GetInputStreamSize(strings.NewReader("hello world"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if size != 11 {
		t.Errorf("got %d, want 11", size)
	}
}

func TestCompareInputStreams(t *testing.T) {
	eq, err := CompareInputStreams(strings.NewReader("hello"), strings.NewReader("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !eq {
		t.Error("expected equal")
	}

	eq, err = CompareInputStreams(strings.NewReader("hello"), strings.NewReader("world"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eq {
		t.Error("expected not equal (same length, different content)")
	}

	eq, err = CompareInputStreams(strings.NewReader("hello"), strings.NewReader("hello!"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eq {
		t.Error("expected not equal (different length)")
	}

	// Exercise the chunk boundary.
	big1 := strings.Repeat("a", 8192*2+5)
	big2 := strings.Repeat("a", 8192*2+5)
	eq, err = CompareInputStreams(strings.NewReader(big1), strings.NewReader(big2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !eq {
		t.Error("expected equal for large equal streams")
	}
}

func TestStartsWithStream(t *testing.T) {
	got, err := StartsWithStream(strings.NewReader("hello world"), []byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Error("expected true")
	}

	got, err = StartsWithStream(strings.NewReader("hi"), []byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Error("expected false for short stream")
	}

	got, err = StartsWithStream(nil, []byte("a"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got {
		t.Error("expected false for nil reader")
	}
}
