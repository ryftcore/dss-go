package xpath10_test

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/corpustest"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xpath10"
)

// Known answers for the XML-DSig transform expressions: every ds:XPath and xpf:XPath
// expression in the upstream dss-xades corpus, evaluated both ways Apache Santuario evaluates
// them - a boolean per candidate node, and a node-set against the document.
//
// The answers in testdata/transform-kat.txt come from the JDK's own XPath engine, which is
// what Santuario 3.0 uses (JDKXPathAPI); see gen/generate-transform.sh. The corpus is
// internal/xmldsig/testdata/corpus, shared rather than duplicated, because these are the same
// documents whose reference digests that package pins.

// transformCorpusRoot is the module-root-relative path (inside the external
// corpus/ tree, see internal/corpustest) of internal/xmldsig's corpus.
const transformCorpusRoot = "internal/xmldsig/testdata/corpus"

type transformRow struct {
	fixture string
	index   int
	mode    string
	expr    string
	answer  string
}

func loadTransformKAT(t *testing.T, katFile string) []transformRow {
	t.Helper()
	f, err := os.Open(filepath.Join(corpustest.Path(t, "."), katFile))
	if err != nil {
		t.Fatalf("open %s: %v", katFile, err)
	}
	defer f.Close()

	var rows []transformRow
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 1<<20)
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		p := strings.Split(line, "\t")
		if len(p) != 5 {
			t.Fatalf("malformed transform-kat line (%d fields): %q", len(p), line)
		}
		i, err := strconv.Atoi(p[1])
		if err != nil {
			t.Fatalf("bad xpath index %q", p[1])
		}
		rows = append(rows, transformRow{p[0], i, p[2], p[3], p[4]})
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("read %s: %v", katFile, err)
	}
	if len(rows) == 0 {
		t.Fatalf("%s is empty", katFile)
	}
	return rows
}

func TestTransformKnownAnswers(t *testing.T) {
	runTransformKAT(t, "transform-kat.txt", corpustest.RootPath(t, transformCorpusRoot))
}

// TestTransformAdversarialKnownAnswers runs the same comparison over testdata/adversarial,
// a set of documents written for this test rather than taken from upstream: here(), which
// Xalan provided and Santuario 3.0 does not; ancestor-or-self over a 300-deep tree; a
// ds:Signature nested in another signature's ds:Object; id() over duplicated Ids; prefixes
// the ds:XPath element does not declare; and expressions the JDK refuses outright. The
// answers still come only from the JDK XPath - see gen/make-adversarial.py and
// gen/TransformXPathOracle.java.
func TestTransformAdversarialKnownAnswers(t *testing.T) {
	runTransformKAT(t, "transform-adversarial-kat.txt", filepath.Join(corpustest.Path(t, "."), "adversarial"))
}

// loadUnsupported reads testdata/transform-adversarial-unsupported.txt: the expressions the
// JDK answers and this package deliberately refuses. See the file's own header.
func loadUnsupported(t *testing.T) map[string]bool {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(corpustest.Path(t, "."), "transform-adversarial-unsupported.txt"))
	if err != nil {
		t.Fatalf("open the unsupported list: %v", err)
	}
	out := map[string]bool{}
	for _, line := range strings.Split(string(b), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out[line] = true
	}
	if len(out) == 0 {
		t.Fatal("the unsupported list is empty")
	}
	return out
}

func runTransformKAT(t *testing.T, katFile, corpus string) {
	rows := loadTransformKAT(t, katFile)

	// An expression the JDK refuses in BOTH modes is refused outright rather than merely
	// mistyped, and this port has to refuse it too - at compile time or at evaluation time,
	// but somewhere. (A "!" on the node-set row alone is the different, narrower case of a
	// boolean-valued expression being asked for a node-set; see below.)
	refused := map[string]bool{}
	for _, r := range rows {
		if r.mode == "bool" && strings.HasPrefix(r.answer, "!") {
			refused[r.fixture+"\x00"+strconv.Itoa(r.index)] = true
		}
	}

	// Expressions outside this package's transform subset. Only the adversarial corpus has
	// any: every expression the upstream fixtures contain is inside it.
	unsupported := map[string]bool{}
	seenUnsupported := map[string]bool{}
	if katFile != "transform-kat.txt" {
		unsupported = loadUnsupported(t)
		t.Cleanup(func() {
			for expr := range unsupported {
				if !seenUnsupported[expr] {
					t.Errorf("%q is on the unsupported list but no adversarial fixture uses it, "+
						"or it now compiles - delete the line", expr)
				}
			}
		})
	}

	docs := map[string]*xmldom.Node{}
	load := func(fixture string) *xmldom.Node {
		if doc, ok := docs[fixture]; ok {
			return doc
		}
		src, err := os.ReadFile(filepath.Join(corpus, filepath.FromSlash(fixture)))
		if err != nil {
			t.Fatalf("read %s: %v", fixture, err)
		}
		doc, err := xmldom.Parse(src, nil)
		if err != nil {
			t.Fatalf("parse %s: %v", fixture, err)
		}
		doc.RegisterIDs()
		docs[fixture] = doc
		return doc
	}

	for _, row := range rows {
		row := row
		t.Run(row.fixture+"/"+strconv.Itoa(row.index)+"/"+row.mode, func(t *testing.T) {
			doc := load(row.fixture)
			el := xpathElements(doc)[row.index]
			expr, err := xpath10.CompileTransform(row.expr, xpath10.NamespaceContextOf(el))
			if unsupported[row.expr] {
				var ue *xpath10.UnsupportedError
				if !errors.As(err, &ue) {
					t.Fatalf("%q is on the unsupported list but CompileTransform returned %v - "+
						"if the subset really was widened, delete the line and let the known "+
						"answer be compared instead", row.expr, err)
				}
				seenUnsupported[row.expr] = true
				return
			}
			if refused[row.fixture+"\x00"+strconv.Itoa(row.index)] {
				if err != nil {
					return
				}
				var evalErr error
				if row.mode == "bool" {
					_, evalErr = expr.EvaluateBoolean(doc)
				} else {
					_, evalErr = expr.EvaluateNodeSet(doc)
				}
				if evalErr == nil {
					t.Fatalf("the JDK refused %q; Go answered it", row.expr)
				}
				return
			}
			if err != nil {
				t.Fatalf("CompileTransform(%q): %v", row.expr, err)
			}

			switch row.mode {
			case "bool":
				var keys []string
				for _, n := range candidateNodes(doc) {
					ok, err := expr.EvaluateBoolean(n)
					if err != nil {
						t.Fatalf("EvaluateBoolean: %v", err)
					}
					if ok {
						keys = append(keys, nodeKey(n))
					}
				}
				// The boolean answer is a set; it is compared sorted, because Xerces sorts an
				// element's attributes by node name while xmldom keeps them in source order.
				sort.Strings(keys)
				if got := strings.Join(keys, " "); got != row.answer {
					t.Fatalf("boolean filtering differs from the JDK\nexpr: %s\nwant: %s\ngot:  %s",
						row.expr, row.answer, got)
				}

			case "nodeset":
				nodes, err := expr.EvaluateNodeSet(doc)
				if err != nil {
					t.Fatalf("EvaluateNodeSet: %v", err)
				}
				var keys []string
				for _, n := range nodes {
					keys = append(keys, nodeKey(n))
				}
				got := strings.Join(keys, " ")
				if strings.HasPrefix(row.answer, "!") {
					// The JDK raises when an expression whose value is a boolean is asked for
					// a NODESET. This package has no type-error channel and yields the empty
					// set instead - documented on EvaluateNodeSet, and pinned here so the
					// divergence cannot drift into "some nodes".
					if got != "" {
						t.Fatalf("the JDK refused %q as a node-set; Go selected %s", row.expr, got)
					}
					return
				}
				if got != row.answer {
					t.Fatalf("node-set differs from the JDK\nexpr: %s\nwant: %s\ngot:  %s",
						row.expr, row.answer, got)
				}

			default:
				t.Fatalf("unknown mode %q", row.mode)
			}
		})
	}
}

// xpathElements returns every ds:XPath and xpf:XPath element in document order - the same list,
// in the same order, that the oracle enumerates.
func xpathElements(doc *xmldom.Node) []*xmldom.Node {
	const (
		nsDSig = "http://www.w3.org/2000/09/xmldsig#"
		nsF2   = "http://www.w3.org/2002/06/xmldsig-filter2"
	)
	var out []*xmldom.Node
	doc.Walk(func(n *xmldom.Node) bool {
		if n.Kind == xmldom.Element && n.Name.Local == "XPath" &&
			(n.Name.Space == nsDSig || n.Name.Space == nsF2) {
			out = append(out, n)
		}
		return true
	})
	return out
}

// candidateNodes is every node the canonicalizer may ask a filter about: the document,
// elements, their attributes (namespace declarations included) and all character-data,
// comment and processing-instruction nodes.
func candidateNodes(doc *xmldom.Node) []*xmldom.Node {
	var out []*xmldom.Node
	var rec func(n *xmldom.Node)
	rec = func(n *xmldom.Node) {
		out = append(out, n)
		if n.Kind == xmldom.Element {
			out = append(out, n.Attrs...)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rec(c)
		}
	}
	rec(doc)
	return out
}

// nodeKey is the structural key the oracle emits: the path of child indices from the document
// node, with "@QName" appended for an attribute. It deliberately avoids any global numbering,
// because Xerces sorts an element's attributes by node name while xmldom keeps them in source
// order, and an index-based key would compare those orders instead of the XPath answers.
func nodeKey(n *xmldom.Node) string {
	if n.Kind == xmldom.Attribute {
		return nodeKey(n.Parent) + "@" + n.Name.QName()
	}
	if n.Kind == xmldom.Document {
		return "/"
	}
	index := 0
	for c := n.Parent.FirstChild; c != nil && c != n; c = c.NextSibling {
		index++
	}
	return nodeKey(n.Parent) + "/" + strconv.Itoa(index)
}
