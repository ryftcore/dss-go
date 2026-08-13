package utils

import (
	"strings"
	"testing"
	"time"
)

func TestDomUtilsStartsWithXmlPreamble(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"plain xml", []byte("<r/>"), true},
		{"utf8 bom then xml", append([]byte{0xEF, 0xBB, 0xBF}, []byte("<r/>")...), true},
		{"not xml", []byte("not xml"), false},
		{"empty", []byte{}, false},
		{"nil", nil, false},
		{"bom only, no xml", []byte{0xEF, 0xBB, 0xBF}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DomUtilsStartsWithXmlPreamble(tt.data); got != tt.want {
				t.Errorf("DomUtilsStartsWithXmlPreamble(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestDomUtilsBuildDOMFromBytesNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil bytes")
		}
	}()
	_, _ = DomUtilsBuildDOMFromBytes(nil)
}

func TestDomUtilsBuildDOMFromBytesEmptySliceDoesNotPanic(t *testing.T) {
	_, err := DomUtilsBuildDOMFromBytes([]byte{})
	if err == nil {
		t.Fatal("expected a parse error for empty (non-nil) input, got nil")
	}
}

func TestDomUtilsBuildDOMFromStringRoundTrip(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(`<r a="1"><c/></r>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	root := doc.DocumentElement()
	if root == nil || root.Name.Local != "r" {
		t.Fatalf("unexpected root: %+v", root)
	}
	if root.AttrValue("", "a") != "1" {
		t.Fatalf("unexpected attribute value: %q", root.AttrValue("", "a"))
	}
}

func TestDomUtilsBuildDOMFromStringInvalidXML(t *testing.T) {
	if _, err := DomUtilsBuildDOMFromString(`<r>`); err == nil {
		t.Fatal("expected error for malformed XML")
	}
}

func TestDomUtilsIsDOMBytes(t *testing.T) {
	if !DomUtilsIsDOMBytes([]byte(`<r/>`)) {
		t.Error("expected <r/> to be recognized as DOM")
	}
	if DomUtilsIsDOMBytes([]byte(`not xml`)) {
		t.Error("expected non-XML to not be recognized as DOM")
	}
	if DomUtilsIsDOMBytes([]byte(`<r>`)) {
		t.Error("expected malformed XML (preamble matches, parse fails) to not be recognized as DOM")
	}
	if DomUtilsIsDOMBytes(nil) {
		t.Error("expected nil to not be recognized as DOM")
	}
}

func TestDomUtilsGetId(t *testing.T) {
	tests := []struct {
		uri  string
		want string
	}{
		{"#signature", "signature"},
		{"signature", "signature"},
		{"#r-signature-1", "r-signature-1"},
		{`#xpointer(id('foo'))`, "foo"},
		{`#xpointer(id("bar"))`, "bar"},
	}
	for _, tt := range tests {
		if got := DomUtilsGetId(tt.uri); got != tt.want {
			t.Errorf("DomUtilsGetId(%q) = %q, want %q", tt.uri, got, tt.want)
		}
	}
}

func TestDomUtilsIsXPointerQuery(t *testing.T) {
	tests := []struct {
		uri  string
		want bool
	}{
		{"", false},
		{"#signature", false},
		{`#xpointer(id('foo'))`, true},
		{`#xmlns(p=urn:x) xpointer(id('foo'))`, true},
		{`#xpointer(/)`, true},
		{"not-an-xpointer", false},
	}
	for _, tt := range tests {
		if got := DomUtilsIsXPointerQuery(tt.uri); got != tt.want {
			t.Errorf("DomUtilsIsXPointerQuery(%q) = %v, want %v", tt.uri, got, tt.want)
		}
	}
}

func TestDomUtilsGetXPointerId(t *testing.T) {
	tests := []struct {
		uri    string
		wantID string
		wantOk bool
	}{
		{`#xpointer(id('foo'))`, "foo", true},
		{`#xpointer(id("bar"))`, "bar", true},
		{`#xpointer(/)`, "", false},
		{`not-xpointer`, "", false},
		{`#xpointer(id(foo))`, "", false}, // no quotes
	}
	for _, tt := range tests {
		id, ok := DomUtilsGetXPointerId(tt.uri)
		if id != tt.wantID || ok != tt.wantOk {
			t.Errorf("DomUtilsGetXPointerId(%q) = (%q, %v), want (%q, %v)", tt.uri, id, ok, tt.wantID, tt.wantOk)
		}
	}
}

func TestDomUtilsIsRootXPointer(t *testing.T) {
	if !DomUtilsIsRootXPointer("#xpointer(/)") {
		t.Error("expected #xpointer(/) to be the root xpointer")
	}
	if DomUtilsIsRootXPointer("#xpointer(id('x'))") {
		t.Error("expected an id xpointer to not be the root xpointer")
	}
}

func TestDomUtilsStartsFromHashAndToElementReference(t *testing.T) {
	if !DomUtilsStartsFromHash("#foo") {
		t.Error("expected #foo to start from hash")
	}
	if DomUtilsStartsFromHash("foo") {
		t.Error("expected foo to not start from hash")
	}
	if DomUtilsStartsFromHash("") {
		t.Error("expected blank to not start from hash")
	}
	if got := DomUtilsToElementReference("foo"); got != "#foo" {
		t.Errorf("DomUtilsToElementReference(foo) = %q, want #foo", got)
	}
	if got := DomUtilsToElementReference("#foo"); got != "#foo" {
		t.Errorf("DomUtilsToElementReference(#foo) = %q, want #foo", got)
	}
}

func TestDomUtilsIsElementReference(t *testing.T) {
	if !DomUtilsIsElementReference("#foo") {
		t.Error("expected #foo to be an element reference")
	}
	if DomUtilsIsElementReference(`#xpointer(id('foo'))`) {
		t.Error("expected an xpointer query to not be a plain element reference")
	}
	if DomUtilsIsElementReference("foo") {
		t.Error("expected foo (no hash) to not be an element reference")
	}
}

func TestDomUtilsCreateXMLGregorianCalendarAndGetDateRoundTrip(t *testing.T) {
	date := time.Date(2013, 11, 23, 11, 22, 52, 0, time.UTC)
	s := DomUtilsCreateXMLGregorianCalendar(date)
	if s != "2013-11-23T11:22:52Z" {
		t.Fatalf("unexpected xml format: %q", s)
	}
	got := DomUtilsGetDate(s)
	if !got.Equal(date) {
		t.Fatalf("round trip mismatch: got %v, want %v", got, date)
	}
}

func TestDomUtilsCreateXMLGregorianCalendarZero(t *testing.T) {
	if got := DomUtilsCreateXMLGregorianCalendar(time.Time{}); got != "" {
		t.Fatalf("expected \"\" for the zero Time, got %q", got)
	}
}

func TestDomUtilsGetDateVariants(t *testing.T) {
	tests := []struct {
		text string
		want time.Time
	}{
		{"2013-11-23T11:22:52Z", time.Date(2013, 11, 23, 11, 22, 52, 0, time.UTC)},
		{"2013-11-23T11:22:52.123Z", time.Date(2013, 11, 23, 11, 22, 52, 123000000, time.UTC)},
		{"2013-11-23T11:22:52+01:00", time.Date(2013, 11, 23, 10, 22, 52, 0, time.UTC)},
		{"not a date", time.Time{}},
	}
	for _, tt := range tests {
		got := DomUtilsGetDate(tt.text)
		if !got.Equal(tt.want) {
			t.Errorf("DomUtilsGetDate(%q) = %v, want %v", tt.text, got, tt.want)
		}
	}
}

func TestDomUtilsGetNodeListAndGetElement(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(`<r><a>1</a><a>2</a><b>3</b></r>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	nodes, err := DomUtilsGetNodeList(doc, "//a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 <a> nodes, got %d", len(nodes))
	}

	if _, err := DomUtilsGetElement(doc, "//a"); err == nil {
		t.Fatal("expected an error for more than one match")
	}

	el, err := DomUtilsGetElement(doc, "//b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if el == nil || el.Name.Local != "b" {
		t.Fatalf("unexpected element: %+v", el)
	}

	missing, err := DomUtilsGetElement(doc, "//missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if missing != nil {
		t.Fatalf("expected nil for no match, got %+v", missing)
	}
}

func TestDomUtilsGetValue(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(`<r><a>  hello  </a></r>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := DomUtilsGetValue(doc, "//a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "hello" {
		t.Fatalf("expected trimmed value 'hello', got %q", v)
	}

	v, err = DomUtilsGetValue(doc, "//missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "" {
		t.Fatalf("expected \"\" for no match, got %q", v)
	}
}

func TestDomUtilsGetChildrenNames(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(`<r>text<a/><!--c--><b/></r>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names, err := DomUtilsGetChildrenNames(doc, "/r")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// DomUtilsGetChildrenNames (deprecated) appends "" for the text/comment children,
	// mirroring Java's null getLocalName() there - see the function's doc comment.
	want := []string{"", "a", "", "b"}
	if len(names) != len(want) {
		t.Fatalf("unexpected names: %#v", names)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("names[%d] = %q, want %q", i, names[i], want[i])
		}
	}
}

func TestDomUtilsSerializeNodeAndXmlToString(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(`<r a="1"><c/></r>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, err := DomUtilsXmlToString(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(s, `<r a="1">`) || !strings.Contains(s, "<c></c>") {
		t.Fatalf("unexpected serialization: %q", s)
	}
}

func TestDomUtilsExcludeComments(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(`<r><!--x--><a><!--y--></a></r>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := DomUtilsExcludeComments(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, err := DomUtilsXmlToString(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(s, "<!--") {
		t.Fatalf("expected comments to be excluded, got %q", s)
	}
}

func TestDomUtilsCreateDeepCopy(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(`<r><a><b/></a></r>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a, err := DomUtilsGetElement(doc, "//a")
	if err != nil || a == nil {
		t.Fatalf("could not find <a>: %v", err)
	}
	copyA := DomUtilsCreateDeepCopy(a)
	if copyA == nil {
		t.Fatal("expected a non-nil copy")
	}
	if copyA.Name.Local != "a" {
		t.Fatalf("unexpected copy element: %+v", copyA)
	}
	// The temporary Id attribute must have been cleaned up on both sides.
	if a.AttrValue("", "Id") != "" {
		t.Errorf("expected the temporary Id attribute to be removed from the original")
	}
	if copyA.AttrValue("", "Id") != "" {
		t.Errorf("expected the temporary Id attribute to be removed from the copy")
	}
	if copyA.FirstChild == nil || copyA.FirstChild.Name.Local != "b" {
		t.Errorf("expected the copy to carry its child <b>")
	}
}

func TestDomUtilsCreateDeepCopyNil(t *testing.T) {
	if got := DomUtilsCreateDeepCopy(nil); got != nil {
		t.Fatalf("expected nil for nil input, got %+v", got)
	}
}
