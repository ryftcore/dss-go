package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/utain/esig/dss/internal/xmldom"
)

// TestSerializeAgainstJavaTransformerOracle is the byte-parity gate for
// DomUtilsSerializeNode and DomUtilsGetNodeBytes.
//
// testdata/serialize/corpus.txt is an adversarial corpus of documents and of nodes inside
// them; testdata/serialize/goldens.txt is what DSS 6.5.RC1's DomUtils.serializeNode and
// DomUtils.getNodeBytes returned for every one of them on OpenJDK 21, produced by
// testdata/gen/SerializeOracle.java, which copies the Java bodies verbatim - including
// the fact that neither sets OMIT_XML_DECLARATION.
//
// This matters because DSSXMLUtils.applyTransforms(node, emptyTransforms) returns
// getNodeBytes(node) for a ds:Reference that carries no ds:Transforms, and those bytes are
// digested into the DigestValue. Every byte of divergence here is a signature that
// verifies against one implementation and not the other.
//
// Regenerate the goldens with:
//
//	CC=$HOME/.m2/repository/commons-codec/commons-codec/1.18.0/commons-codec-1.18.0.jar
//	cd testdata/gen && javac -cp "$CC" -d /tmp/serOracle SerializeOracle.java &&
//	  java -cp "$CC:/tmp/serOracle" SerializeOracle ../serialize/corpus.txt ../serialize/goldens.txt
func TestSerializeAgainstJavaTransformerOracle(t *testing.T) {
	corpus := readCorpus(t, "testdata/serialize/corpus.txt")
	goldens := readGoldens(t, "testdata/serialize/goldens.txt")

	if len(corpus) == 0 || len(corpus) != len(goldens) {
		t.Fatalf("corpus has %d cases, goldens %d", len(corpus), len(goldens))
	}

	for _, c := range corpus {
		g, ok := goldens[c.name]
		if !ok {
			t.Fatalf("case %q has no golden; regenerate goldens.txt", c.name)
		}
		t.Run(c.name, func(t *testing.T) {
			doc, target, err := c.build()

			// A case the corpus records as a known divergence is not compared. It must
			// still be refused somewhere on this side - at parse time, at serialization
			// time, or by returning the error Java's exception maps to.
			if c.divergence != "" {
				if err == nil {
					_, serr := DomUtilsSerializeNode(target)
					err = serr
				}
				if err == nil {
					t.Fatalf("recorded as a divergence but accepted outright: %s", c.divergence)
				}
				t.Logf("known divergence: %s (Go: %v)", c.divergence, err)
				return
			}

			// A case Java could not even build - a document Xerces rejects - has to be
			// rejected here too, and there is nothing further to compare.
			if want := g["tree"]; strings.HasPrefix(want, "ERROR") {
				if err == nil {
					t.Fatalf("built fine here, but Java refused it: %s", want)
				}
				return
			}
			if err != nil {
				t.Fatalf("building the case: %v", err)
			}

			// The structural dump comes next: comparing serializations of two different
			// trees would be meaningless, and a mismatch here is a parser or DOM
			// divergence rather than a serializer one.
			if want, ok := g["tree"]; ok {
				if got := dumpTree(doc); got != decode(t, want) {
					t.Fatalf("DOM structure differs from Xerces'\n got:\n%s\nwant:\n%s",
						got, decode(t, want))
				}
			}

			checkBytes(t, "DomUtilsSerializeNode", g["serialize"], func() ([]byte, error) {
				return DomUtilsSerializeNode(target)
			})
			checkBytes(t, "DomUtilsGetNodeBytes", g["bytes"], func() ([]byte, error) {
				return DomUtilsGetNodeBytes(target)
			})
		})
	}
}

// checkBytes compares one call against its golden. "NULL" is a Java null return, and
// "ERROR ..." records that Java threw - which this port has to turn into an error too,
// since DomUtils propagates it as a DSSException.
func checkBytes(t *testing.T, what, golden string, call func() ([]byte, error)) {
	t.Helper()
	got, err := call()
	switch {
	case strings.HasPrefix(golden, "ERROR"):
		if err == nil {
			t.Errorf("%s = %q, want an error (Java: %s)", what, got, golden)
		}
		return
	case golden == "NULL":
		if err != nil {
			t.Errorf("%s errored (%v), want nil", what, err)
		} else if got != nil {
			t.Errorf("%s = %q, want nil", what, got)
		}
		return
	}
	if err != nil {
		t.Fatalf("%s errored: %v", what, err)
	}
	want, decErr := base64.StdEncoding.DecodeString(golden)
	if decErr != nil {
		t.Fatalf("undecodable golden: %v", decErr)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s mismatch\n got: %s\nwant: %s\n got hex: %x\nwant hex: %x",
			what, printable(got), printable(want), got, want)
	}
}

func printable(b []byte) string {
	var sb strings.Builder
	for _, c := range b {
		if c >= 0x20 && c < 0x7F {
			sb.WriteByte(c)
		} else {
			fmt.Fprintf(&sb, "\\x%02x", c)
		}
	}
	return sb.String()
}

func decode(t *testing.T, b64 string) string {
	t.Helper()
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		t.Fatalf("undecodable golden: %v", err)
	}
	return string(b)
}

// ------------------------------------------------------------------- corpus records

type serializeCase struct {
	name       string
	parse      string
	parseBytes []byte
	dsl        string
	sel        string
	divergence string
}

func (c serializeCase) build() (doc, target *xmldom.Node, err error) {
	switch {
	case c.parseBytes != nil:
		doc, err = xmldom.Parse(c.parseBytes, nil)
	case c.parse != "":
		doc, err = xmldom.Parse([]byte(c.parse), nil)
	default:
		doc = xmldom.NewDocument()
	}
	if err != nil {
		return nil, nil, err
	}
	if c.dsl != "" {
		// A build script alongside a parse source is DOM surgery on the parsed document -
		// parse, graft, serialize, which is the shape XAdES signing has - and starts at
		// the document element.
		start := doc
		if el := doc.DocumentElement(); el != nil {
			start = el
		}
		if err = runDSL(doc, c.dsl, start); err != nil {
			return nil, nil, err
		}
	}
	target, err = selectNode(doc, c.sel)
	return doc, target, err
}

func selectNode(doc *xmldom.Node, sel string) (*xmldom.Node, error) {
	if sel == "" || sel == "." {
		return doc, nil
	}
	attr := ""
	if at := strings.Index(sel, "@"); at >= 0 {
		attr = sel[at+1:]
		sel = sel[:at]
	}
	n := doc
	if sel != "" {
		for _, part := range strings.Split(sel, ".") {
			idx, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("bad selector %q: %w", sel, err)
			}
			kids := n.Children()
			if idx >= len(kids) {
				return nil, fmt.Errorf("selector %q does not resolve", sel)
			}
			n = kids[idx]
		}
	}
	if attr != "" {
		for _, a := range n.Attrs {
			if a.Name.QName() == attr {
				return a, nil
			}
		}
		return nil, fmt.Errorf("no attribute %q", attr)
	}
	return n, nil
}

// runDSL mirrors SerializeOracle.runDSL. The ops map onto the org.w3c.dom calls named in
// corpus.txt, so "a" is setAttributeNS (a namespace-aware attribute) and "ap" is
// setAttribute (node name only, null namespace URI) - the distinction
// DomUtilsAddNamespaceAttribute turns on.
func runDSL(doc *xmldom.Node, script string, start *xmldom.Node) error {
	cur := start
	for _, op := range strings.Split(script, ";") {
		op = strings.TrimSpace(op)
		if op == "" || op == "~" {
			continue
		}
		tok := strings.Split(op, " ")
		switch tok[0] {
		case "e":
			el := xmldom.NewElement(nameFromNS(unescape(tok[1]), unescape(tok[2])))
			cur.AppendChild(el)
			cur = el
		case "a":
			cur.SetAttr(nameFromNS(unescape(tok[1]), unescape(tok[2])), unescape(tok[3]))
		case "ap":
			// Element.setAttribute(name, value): the whole thing is the node name and
			// the namespace URI is null, so a colon in it is not a prefix.
			cur.SetAttr(xmldom.Name{Local: unescape(tok[1])}, unescape(tok[2]))
		case "t":
			cur.AppendChild(xmldom.NewText(unescape(tok[1])))
		case "c":
			cur.AppendChild(xmldom.NewComment(unescape(tok[1])))
		case "d":
			cur.AppendChild(xmldom.NewCDATA(unescape(tok[1])))
		case "p":
			cur.AppendChild(xmldom.NewProcInst(unescape(tok[1]), unescape(tok[2])))
		case "/":
			cur = cur.Parent
		default:
			return fmt.Errorf("unknown DSL op %q", tok[0])
		}
	}
	return nil
}

// nameFromNS is createElementNS/createAttributeNS: the qualified name is split on its
// colon and the namespace URI is attached. "-" is Java's null.
func nameFromNS(ns, qname string) xmldom.Name {
	if ns == "-" {
		ns = ""
	}
	if colon := strings.Index(qname, ":"); colon > 0 {
		return xmldom.Name{Space: ns, Prefix: qname[:colon], Local: qname[colon+1:]}
	}
	return xmldom.Name{Space: ns, Local: qname}
}

func readCorpus(t *testing.T, path string) []serializeCase {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open corpus: %v", err)
	}
	defer f.Close()

	var out []serializeCase
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, _ := strings.Cut(line, " ")
		switch key {
		case "case":
			out = append(out, serializeCase{name: val, sel: "."})
		case "parse":
			out[len(out)-1].parse = unescape(val)
		case "parseb64":
			b, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				t.Fatalf("case %s: bad parseb64: %v", out[len(out)-1].name, err)
			}
			out[len(out)-1].parseBytes = b
		case "build":
			out[len(out)-1].dsl = val
		case "sel":
			out[len(out)-1].sel = val
		case "divergence":
			c := &out[len(out)-1]
			if c.divergence != "" {
				c.divergence += " "
			}
			c.divergence += val
		default:
			t.Fatalf("unknown corpus key %q", key)
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	return out
}

func readGoldens(t *testing.T, path string) map[string]map[string]string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open goldens: %v", err)
	}
	defer f.Close()

	out := map[string]map[string]string{}
	cur := ""
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, _ := strings.Cut(line, " ")
		if key == "case" {
			cur = val
			out[cur] = map[string]string{}
			continue
		}
		out[cur][key] = val
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("read goldens: %v", err)
	}
	return out
}

// unescape mirrors SerializeOracle.unescape.
func unescape(s string) string {
	if s == "~" {
		return ""
	}
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' || i+1 >= len(s) {
			sb.WriteByte(c)
			continue
		}
		i++
		switch s[i] {
		case 'n':
			sb.WriteByte('\n')
		case 'r':
			sb.WriteByte('\r')
		case 't':
			sb.WriteByte('\t')
		case 's':
			sb.WriteByte(' ')
		case '\\':
			sb.WriteByte('\\')
		case 'u':
			u, err := strconv.ParseUint(s[i+1:i+5], 16, 16)
			if err != nil {
				break
			}
			i += 4
			// A Java String is UTF-16, so a supplementary character is written as its
			// surrogate pair; a Go string has no way to hold the halves separately, so
			// the pair is recombined here.
			if utf16.IsSurrogate(rune(u)) && i+6 < len(s) && s[i+1] == '\\' && s[i+2] == 'u' {
				if lo, err := strconv.ParseUint(s[i+3:i+7], 16, 16); err == nil {
					if r := utf16.DecodeRune(rune(u), rune(lo)); r != 0xFFFD {
						sb.WriteRune(r)
						i += 6
						break
					}
				}
			}
			sb.WriteRune(rune(u))
		default:
			sb.WriteByte('\\')
			sb.WriteByte(s[i])
		}
	}
	return sb.String()
}

// ---------------------------------------------------------------- structural dump

// dumpTree reproduces SerializeOracle.dumpTree: one line per node in document order,
// with an element's attributes listed in Xerces' NamedNodeMap order, which is sorted by
// node name. This port keeps attributes in source order instead - a difference confined
// to iteration, since both canonicalization and this serializer impose their own order -
// so the dump sorts a copy to compare like with like.
func dumpTree(doc *xmldom.Node) string {
	var sb strings.Builder
	var walk func(path string, n *xmldom.Node)
	walk = func(path string, n *xmldom.Node) {
		p := path
		if p == "" {
			p = "."
		}
		sb.WriteString(dumpNode(p, n))
		sb.WriteByte('\n')
		i := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if path == "" {
				walk(strconv.Itoa(i), c)
			} else {
				walk(path+"."+strconv.Itoa(i), c)
			}
			i++
		}
	}
	walk("", doc)
	return sb.String()
}

func dumpNode(path string, n *xmldom.Node) string {
	var sb strings.Builder
	sb.WriteString(path)
	sb.WriteByte(' ')
	sb.WriteString(dumpKind(n))
	sb.WriteByte(' ')
	sb.WriteString(dumpName(n))
	if n.Kind == xmldom.Element {
		names := make([]string, 0, len(n.Attrs))
		for _, a := range n.Attrs {
			names = append(names, a.Name.QName())
		}
		sort.SliceStable(names, func(i, j int) bool { return lessUTF16(names[i], names[j]) })
		for _, q := range names {
			sb.WriteByte(' ')
			sb.WriteString(q)
		}
	}
	return sb.String()
}

func lessUTF16(a, b string) bool {
	au, bu := utf16.Encode([]rune(a)), utf16.Encode([]rune(b))
	for i := 0; i < len(au) && i < len(bu); i++ {
		if au[i] != bu[i] {
			return au[i] < bu[i]
		}
	}
	return len(au) < len(bu)
}

func dumpKind(n *xmldom.Node) string {
	switch n.Kind {
	case xmldom.Document:
		return "document"
	case xmldom.Element:
		return "element"
	case xmldom.Attribute:
		return "attribute"
	case xmldom.Text:
		return "text"
	case xmldom.CDATA:
		return "cdata"
	case xmldom.Comment:
		return "comment"
	case xmldom.ProcInst:
		return "pi"
	}
	return "?"
}

// dumpName is org.w3c.dom.Node.getNodeName(), which is the qualified name for an element
// or attribute, the target for a processing instruction and a fixed "#name" otherwise.
func dumpName(n *xmldom.Node) string {
	switch n.Kind {
	case xmldom.Document:
		return "#document"
	case xmldom.Element, xmldom.Attribute:
		return n.Name.QName()
	case xmldom.Text:
		return "#text"
	case xmldom.CDATA:
		return "#cdata-section"
	case xmldom.Comment:
		return "#comment"
	case xmldom.ProcInst:
		return n.Name.Local
	}
	return "?"
}
