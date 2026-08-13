package xmldom

import "testing"

// mustParse parses src with the secure defaults and fails the test on error.
func mustParse(t *testing.T, src string) *Node {
	t.Helper()
	doc, err := Parse([]byte(src), nil)
	if err != nil {
		t.Fatalf("Parse(%q): %v", src, err)
	}
	return doc
}

// mustParseRoot parses src and returns its document element.
func mustParseRoot(t *testing.T, src string) *Node {
	t.Helper()
	return mustParse(t, src).DocumentElement()
}

// mustSerialize serializes n without an XML declaration.
func mustSerialize(t *testing.T, n *Node) string {
	t.Helper()
	b, err := n.Bytes(&SerializeOptions{})
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	return string(b)
}
