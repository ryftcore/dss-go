package qwac

import "testing"

func TestLinkHeaderParser_Parse_SingleLinkWithParams(t *testing.T) {
	got, err := NewLinkHeaderParser().Parse(`<https://example.com/binding>; rel="tls-certificate-binding"; type="application/pkcs7-signature"`)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 link header, got %d", len(got))
	}
	if got[0].URL() != "https://example.com/binding" {
		t.Fatalf("URL = %q, want %q", got[0].URL(), "https://example.com/binding")
	}
	attrs := got[0].Attributes()
	if rel, ok := attrs["rel"]; !ok || rel == nil || *rel != "tls-certificate-binding" {
		t.Fatalf("rel attribute = %v, want tls-certificate-binding", attrs["rel"])
	}
	if typ, ok := attrs["type"]; !ok || typ == nil || *typ != "application/pkcs7-signature" {
		t.Fatalf("type attribute = %v, want application/pkcs7-signature", attrs["type"])
	}
}

func TestLinkHeaderParser_Parse_MultipleLinksCommaSeparated(t *testing.T) {
	got, err := NewLinkHeaderParser().Parse(`<https://example.com/a>; rel="first", <https://example.com/b>; rel="second"`)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 link headers, got %d", len(got))
	}
	if got[0].URL() != "https://example.com/a" || got[1].URL() != "https://example.com/b" {
		t.Fatalf("unexpected URLs: %q, %q", got[0].URL(), got[1].URL())
	}
}

func TestLinkHeaderParser_Parse_CommaInsideQuotesNotSplit(t *testing.T) {
	got, err := NewLinkHeaderParser().Parse(`<https://example.com/a>; title="hello, world"`)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 link header (comma inside quotes must not split), got %d", len(got))
	}
	if title := got[0].Attributes()["title"]; title == nil || *title != "hello, world" {
		t.Fatalf("title attribute = %v, want %q", title, "hello, world")
	}
}

func TestLinkHeaderParser_Parse_FlagAttributeWithoutEquals(t *testing.T) {
	got, err := NewLinkHeaderParser().Parse(`<https://example.com/a>; noopt`)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	attrs := got[0].Attributes()
	v, ok := attrs["noopt"]
	if !ok {
		t.Fatal("expected the 'noopt' flag key to be present")
	}
	if v != nil {
		t.Fatalf("expected the 'noopt' flag value to be nil, got %q", *v)
	}
}

func TestLinkHeaderParser_Parse_NoAttributes(t *testing.T) {
	got, err := NewLinkHeaderParser().Parse(`<https://example.com/a>`)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(got) != 1 || got[0].URL() != "https://example.com/a" {
		t.Fatalf("unexpected result: %+v", got)
	}
	if len(got[0].Attributes()) != 0 {
		t.Fatalf("expected no attributes, got %v", got[0].Attributes())
	}
}

func TestLinkHeaderParser_Parse_EmptyOrBlankErrors(t *testing.T) {
	for _, in := range []string{"", "   "} {
		if _, err := NewLinkHeaderParser().Parse(in); err == nil {
			t.Fatalf("Parse(%q): expected an error", in)
		}
	}
}

func TestLinkHeaderParser_Parse_MissingOpenAngleBracketErrors(t *testing.T) {
	if _, err := NewLinkHeaderParser().Parse(`https://example.com/a>; rel="x"`); err == nil {
		t.Fatal("expected an error for a missing leading '<'")
	}
}

func TestLinkHeaderParser_Parse_MissingCloseAngleBracketErrors(t *testing.T) {
	if _, err := NewLinkHeaderParser().Parse(`<https://example.com/a; rel="x"`); err == nil {
		t.Fatal("expected an error for a missing '>'")
	}
}
