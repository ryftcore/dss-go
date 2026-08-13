package xmldom

import (
	"bytes"
	"encoding/base64"
	"testing"
)

// The ownerDocument contract: org.w3c.dom assigns a node's owning document when the document
// factory creates it and NEVER clears it, so removeChild leaves getOwnerDocument() intact.
// DomUtils.serializeNode reads exactly that to choose its output encoding, which means the
// answer is inside the digest of a no-transform ds:Reference: an element lifted out of an
// ISO-8859-1 document and serialized on its own is written in ISO-8859-1 by Java, not UTF-8.
//
// The goldens below were produced on OpenJDK 21 by the DSS bodies verbatim - the same
// serializeNode/getNodeBytes pair xml/utils/testdata/gen/SerializeOracle.java copies -
// against a document built with DSS's secure DocumentBuilderFactory:
//
//	Document d = parse(the UTF-8 bytes of
//	                   "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><r xmlns=\"urn:D\"><a>café</a></r>");
//	// the declaration wins, so Xerces decodes those bytes as Latin-1 and the text is "cafÃ©"
//	Element a = (Element) d.getDocumentElement().getFirstChild();
//	d.getDocumentElement().removeChild(a);
//	serializeNode(a) -> <?xml version="1.0" encoding="ISO-8859-1"?><a xmlns="urn:D">caf..</a>
//	getNodeBytes(a)  -> <a xmlns="urn:D">caf..</a>
//
// where the two bytes written for the last character are C3 A9 - ISO-8859-1 for the two
// characters the source's UTF-8 bytes decode to under that declaration.
func TestDetachedNodeKeepsItsOwnerDocumentEncoding(t *testing.T) {
	// The two bytes C3 A9 are the UTF-8 encoding of "é"; under the ISO-8859-1 the
	// declaration insists on they decode to two characters and encode straight back.
	const src = "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><r xmlns=\"urn:D\"><a>caf\xc3\xa9</a></r>"
	doc := mustParse(t, src)
	root := doc.DocumentElement()
	a := root.FirstChild

	if got := a.OwnerDocument(); got != doc {
		t.Fatalf("an attached node's OwnerDocument = %v, want the document", got)
	}
	if doc.OwnerDocument() != nil {
		t.Error("a Document's OwnerDocument must be nil, as getOwnerDocument() is")
	}

	root.RemoveChild(a)

	if a.Document() != nil {
		t.Error("Document() must still answer nil for a detached node - it is positional")
	}
	if got := a.OwnerDocument(); got != doc {
		t.Fatalf("a detached node's OwnerDocument = %v, want the document it was removed from", got)
	}
	if got := a.XMLEncoding(); got != "ISO-8859-1" {
		t.Fatalf("XMLEncoding after detaching = %q, want %q", got, "ISO-8859-1")
	}

	want := "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><a xmlns=\"urn:D\">caf\xc3\xa9</a>"
	got, err := a.Bytes(nil)
	if err != nil {
		t.Fatalf("Bytes: %v", err)
	}
	if string(got) != want {
		t.Fatalf("serializing a detached subtree\nwant %q\ngot  %q", want, got)
	}

	// A descendant of the detached root answers through it.
	text := a.FirstChild
	if got := text.OwnerDocument(); got != doc {
		t.Errorf("a detached node's child OwnerDocument = %v, want the document", got)
	}

	// Re-attaching somewhere else makes ownership positional again.
	other := NewDocument()
	otherRoot := NewElement(Name{Local: "x"})
	other.AppendChild(otherRoot)
	otherRoot.AppendChild(a)
	if got := a.OwnerDocument(); got != other {
		t.Errorf("after re-attaching, OwnerDocument = %v, want the new document", got)
	}
	if got := a.XMLEncoding(); got != "" {
		t.Errorf("after re-attaching, XMLEncoding = %q, want the new document's (none)", got)
	}
}

// TestOwnerDocumentSurvivesEveryDetachPath covers the other four ways a node leaves a tree.
// Each has to behave like removeChild, because org.w3c.dom clears ownerDocument in none of
// them.
func TestOwnerDocumentSurvivesEveryDetachPath(t *testing.T) {
	newDoc := func() *Node {
		return mustParse(t, `<?xml version="1.0" encoding="ISO-8859-1"?><r a="1"><b/><c/></r>`)
	}

	t.Run("ReplaceChild", func(t *testing.T) {
		doc := newDoc()
		root := doc.DocumentElement()
		old := root.FirstChild
		root.ReplaceChild(NewElement(Name{Local: "z"}), old)
		if got := old.OwnerDocument(); got != doc {
			t.Errorf("OwnerDocument = %v, want the document", got)
		}
	})

	t.Run("RemoveAttr", func(t *testing.T) {
		doc := newDoc()
		root := doc.DocumentElement()
		attr := root.Attr("", "a")
		if !root.RemoveAttr("", "a") {
			t.Fatal("RemoveAttr did not remove anything")
		}
		if got := attr.OwnerDocument(); got != doc {
			t.Errorf("OwnerDocument = %v, want the document", got)
		}
	})

	t.Run("SetTextContent", func(t *testing.T) {
		doc := newDoc()
		root := doc.DocumentElement()
		b := root.FirstChild
		root.SetTextContent("replaced")
		if got := b.OwnerDocument(); got != doc {
			t.Errorf("OwnerDocument = %v, want the document", got)
		}
	})

	t.Run("Clone", func(t *testing.T) {
		doc := newDoc()
		// cloneNode keeps the source's ownerDocument until the copy is inserted.
		clone := doc.DocumentElement().Clone(true)
		if got := clone.OwnerDocument(); got != doc {
			t.Errorf("OwnerDocument = %v, want the source document", got)
		}
		if got := clone.XMLEncoding(); got != "ISO-8859-1" {
			t.Errorf("XMLEncoding = %q, want the source document's", got)
		}
	})

	t.Run("Import", func(t *testing.T) {
		doc := newDoc()
		target := NewDocument()
		// importNode assigns the IMPORTING document, not the source one.
		imported := target.Import(doc.DocumentElement(), true)
		if got := imported.OwnerDocument(); got != target {
			t.Errorf("OwnerDocument = %v, want the importing document", got)
		}
		if got := imported.XMLEncoding(); got != "" {
			t.Errorf("XMLEncoding = %q, want the importing document's (none)", got)
		}
	})

	t.Run("NeverAttached", func(t *testing.T) {
		// Nothing in this package is created from a document factory, so a fresh node has
		// no owner until it is inserted. Recorded rather than pinned to Java, which has no
		// equivalent state.
		if got := NewElement(Name{Local: "fresh"}).OwnerDocument(); got != nil {
			t.Errorf("OwnerDocument = %v, want nil", got)
		}
	})
}

// TestLatin2DecodesToTheSameCharactersAsJava pins the ISO-8859-2 table against the JDK's
// own charset rather than against itself.
//
// The serializer goldens cannot see an error here: decoding and encoding read the SAME
// table, so a wrong entry round-trips to the same byte and the output is unchanged. What
// differs is the CHARACTER the document is understood to contain - which is what
// canonicalization writes out as UTF-8, what XPath compares strings against, and therefore
// what a digest is ultimately taken over.
//
// The golden is the base64 of new String(bytes A0..FF, "ISO-8859-2").getBytes("UTF-8") on
// OpenJDK 21.
func TestLatin2DecodesToTheSameCharactersAsJava(t *testing.T) {
	const javaUTF8Base64 = "wqDEhMuYxYHCpMS9xZrCp8KoxaDFnsWkxbnCrcW9xbvCsMSFy5vFgsK0xL7Fm8uHwrjFocWf" +
		"xaXFusudxb7FvMWUw4HDgsSCw4TEucSGw4fEjMOJxJjDi8Saw43DjsSOxJDFg8WHw5PDlMWQ" +
		"w5bDl8WYxa7DmsWww5zDncWiw5/FlcOhw6LEg8OkxLrEh8OnxI3DqcSZw6vEm8Otw67Ej8SR" +
		"xYTFiMOzw7TFkcO2w7fFmcWvw7rFscO8w73Fo8uZ"

	want, err := base64.StdEncoding.DecodeString(javaUTF8Base64)
	if err != nil {
		t.Fatal(err)
	}
	high := make([]byte, 96)
	for i := range high {
		high[i] = byte(0xA0 + i)
	}
	if got := decodeLatin2(high); !bytes.Equal(got, want) {
		t.Fatalf("decodeLatin2 of A0..FF\nwant %q\ngot  %q", want, got)
	}

	// And the reverse map has to be the exact inverse, or writeLatin2 would substitute '?'
	// for a character the charset does carry.
	for i, r := range latin2High {
		b, ok := latin2Byte(r)
		if !ok || b != byte(0xA0+i) {
			t.Errorf("latin2Byte(U+%04X) = %#x, %t; want %#x, true", r, b, ok, 0xA0+i)
		}
	}
	// Below A0 the charset is Unicode, and above the table it has nothing.
	if b, ok := latin2Byte(0x41); !ok || b != 0x41 {
		t.Errorf("latin2Byte('A') = %#x, %t", b, ok)
	}
	if _, ok := latin2Byte('€'); ok {
		t.Error("latin2Byte(U+20AC) reported a byte; ISO-8859-2 has no euro sign")
	}
}
