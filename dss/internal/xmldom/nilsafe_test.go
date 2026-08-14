package xmldom

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// TestNilReceiversAreSafe covers the read-only accessors, which callers reach through
// chains such as doc.DocumentElement().FirstElementChild() that can legitimately end
// in nil. Only the mutating entry points panic, and those are covered by
// TestMutationPanics.
func TestNilReceiversAreSafe(t *testing.T) {
	var n *Node

	if got := n.Attr("", "a"); got != nil {
		t.Error("Attr")
	}
	if got := n.AttrValue("", "a"); got != "" {
		t.Error("AttrValue")
	}
	if n.RemoveAttr("", "a") {
		t.Error("RemoveAttr")
	}
	if n.Document() != nil {
		t.Error("Document")
	}
	if n.DocumentElement() != nil {
		t.Error("DocumentElement")
	}
	if n.Children() != nil {
		t.Error("Children")
	}
	if n.Elements() != nil {
		t.Error("Elements")
	}
	if n.FirstElementChild() != nil {
		t.Error("FirstElementChild")
	}
	if n.NextElementSibling() != nil {
		t.Error("NextElementSibling")
	}
	if n.Ancestors() != nil {
		t.Error("Ancestors")
	}
	if n.Depth() != 0 {
		t.Error("Depth")
	}
	if n.Contains(n) {
		t.Error("Contains")
	}
	if n.TextContent() != "" {
		t.Error("TextContent")
	}
	n.Walk(func(*Node) bool { t.Error("Walk visited a nil node"); return true })
	if uri, ok := n.LookupNamespaceURI("p"); ok || uri != "" {
		t.Error("LookupNamespaceURI")
	}
	if p, ok := n.LookupPrefix("urn:1"); ok || p != "" {
		t.Error("LookupPrefix")
	}
	if n.ElementByID("x") != nil {
		t.Error("ElementByID")
	}
	if n.IDAttrs(n) != nil {
		t.Error("IDAttrs")
	}
	if n.DuplicateIDs() != nil {
		t.Error("DuplicateIDs")
	}
	if n.Clone(true) != nil {
		t.Error("Clone")
	}
}

// TestAccessorsOnLeafKinds checks that the element-oriented accessors are inert on
// nodes that cannot have element children.
func TestAccessorsOnLeafKinds(t *testing.T) {
	for _, n := range []*Node{NewText("t"), NewComment("c"), NewProcInst("p", "d"), NewAttr(Name{Local: "a"}, "v")} {
		if n.Elements() != nil || n.FirstElementChild() != nil || n.Children() != nil {
			t.Errorf("%s should have no children", n.Kind)
		}
		if n.NextElementSibling() != nil {
			t.Errorf("%s detached node should have no next element sibling", n.Kind)
		}
		if n.DocumentElement() != nil {
			t.Errorf("%s should have no document element", n.Kind)
		}
	}
	// An element whose only children are non-elements.
	el := mustParseRoot(t, `<r>text<!--c--><?pi d?></r>`)
	if el.Elements() != nil {
		t.Error("Elements should be nil when there are no element children")
	}
	if el.FirstElementChild() != nil {
		t.Error("FirstElementChild should be nil when there are no element children")
	}
	if el.FirstChild.NextElementSibling() != nil {
		t.Error("NextElementSibling should be nil when no element follows")
	}
}

func TestParseReaderPropagatesReadErrors(t *testing.T) {
	want := errors.New("read boom")
	if _, err := ParseReader(failingReader{want}, nil); !errors.Is(err, want) {
		t.Errorf("ParseReader error = %v, want %v", err, want)
	}
	doc, err := ParseReader(strings.NewReader(`<r a="1"/>`), nil)
	if err != nil {
		t.Fatalf("ParseReader: %v", err)
	}
	if got := doc.DocumentElement().AttrValue("", "a"); got != "1" {
		t.Errorf("a = %q, want %q", got, "1")
	}
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

func TestLineColCountsRunesNotBytes(t *testing.T) {
	buf := []byte("aé€\U0001F600x\nyz")
	for _, tc := range []struct {
		off       int64
		line, col int
	}{
		{0, 1, 1},
		{1, 1, 2},  // after "a"
		{3, 1, 3},  // after the 2-byte é
		{6, 1, 4},  // after the 3-byte €
		{10, 1, 5}, // after the 4-byte emoji
		{11, 1, 6}, // after "x"
		{12, 2, 1}, // after the newline
		{14, 2, 3},
		{1000, 2, 3}, // clamped to the end
		{-5, 1, 1},   // clamped to the start
	} {
		line, col := lineCol(buf, tc.off)
		if line != tc.line || col != tc.col {
			t.Errorf("lineCol(%d) = (%d, %d), want (%d, %d)", tc.off, line, col, tc.line, tc.col)
		}
	}
}

// TestUnsupportedEncodingIgnoresACharsetReaderForKnownNames checks that the built-in
// table wins: a caller-supplied CharsetReader must not be able to redefine UTF-8.
func TestUnsupportedEncodingIgnoresACharsetReaderForKnownNames(t *testing.T) {
	called := false
	doc, err := Parse([]byte(`<?xml version="1.0" encoding="UTF-8"?><r>a</r>`), &ParseOptions{
		CharsetReader: func(string, io.Reader) (io.Reader, error) {
			called = true
			return strings.NewReader("<x/>"), nil
		},
	})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if called {
		t.Error("CharsetReader was consulted for a built-in encoding")
	}
	if got := doc.DocumentElement().Name.Local; got != "r" {
		t.Errorf("root = %q, want %q", got, "r")
	}
}

func TestSerializeToWriter(t *testing.T) {
	var b strings.Builder
	if err := mustParse(t, `<r a="1"/>`).Serialize(&b, &SerializeOptions{}); err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	if got, want := b.String(), `<r a="1"/>`; got != want {
		t.Errorf("Serialize = %q, want %q", got, want)
	}
}
