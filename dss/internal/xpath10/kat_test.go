package xpath10

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/utain/esig/dss/internal/corpustest"

	"github.com/utain/esig/dss/internal/xmldom"
)

// TestKnownAnswers rebuilds testdata/kat.txt from the fixtures and requires it to come out
// byte for byte identical to what gen/XPathOracle.java produced with javax.xml.xpath.
//
// The vectors cover every expression of testdata/expressions.txt - the deduplicated inventory
// of every XPath expression the upstream tree can hand to that engine - evaluated from every
// context node of every fixture, so a disagreement anywhere in the subset fails here.
//
// Regenerating rather than parsing-and-comparing means the node table at the head of each
// fixture is checked too, which pins xmldom's tree against Xerces's: text node splitting,
// CDATA sections, comments and processing instructions all have to land in the same places
// for the index paths in the answers to mean the same thing.
func TestKnownAnswers(t *testing.T) {
	// The two answer sets the oracle writes. "" is the upstream inventory - expressions.txt
	// against fixtures/, answers in kat.txt. "semantics" is the grammar-reachable conversion
	// rules the inventory does not happen to exercise; see testdata/semantics.txt.
	for _, set := range []string{"", "semantics"} {
		name := set
		if name == "" {
			name = "inventory"
		}
		t.Run(name, func(t *testing.T) {
			exprFile, fixtureDir, katFile := setFiles(set)
			want, err := os.ReadFile(filepath.Join(corpustest.Path(t, "."), katFile))
			if err != nil {
				t.Fatal(err)
			}
			got := buildKAT(t, exprFile, fixtureDir, katFile)
			if bytes.Equal(got, want) {
				return
			}
			t.Errorf("regenerated vectors differ from the Java oracle:\n%s", firstDiff(got, want))
			if os.Getenv("XPATH10_WRITE_KAT") != "" {
				if err := os.WriteFile(filepath.Join(corpustest.Path(t, "."), katFile+".got"), got, 0o644); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// setFiles names the three files of an answer set, matching the convention in the header of
// gen/XPathOracle.java: the empty set name is the default inventory triple.
func setFiles(set string) (exprFile, fixtureDir, katFile string) {
	if set == "" {
		return "expressions.txt", "fixtures", "kat.txt"
	}
	return set + ".txt", set + "-fixtures", set + "-kat.txt"
}

func buildKAT(t *testing.T, exprFile, fixtureDir, katFile string) []byte {
	t.Helper()
	ns := loadNamespaces(t)
	exprs := loadExpressions(t, exprFile)

	var buf bytes.Buffer
	for _, line := range katHeader(t, katFile) {
		buf.WriteString(line)
		buf.WriteByte('\n')
	}

	compiled := make([]*Expr, len(exprs))
	for i, e := range exprs {
		x, err := Compile(e, ns)
		if err != nil {
			t.Fatalf("Compile(%q): %v", e, err)
		}
		compiled[i] = x
	}

	for _, fixture := range fixtureFiles(t, fixtureDir) {
		src, err := os.ReadFile(fixture)
		if err != nil {
			t.Fatal(err)
		}
		doc, err := xmldom.Parse(src, nil)
		if err != nil {
			t.Fatalf("%s: %v", fixture, err)
		}

		var all []*xmldom.Node
		var collect func(n *xmldom.Node)
		collect = func(n *xmldom.Node) {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				all = append(all, c)
				collect(c)
			}
		}
		collect(doc)

		contexts := []*xmldom.Node{doc}
		for _, n := range all {
			if n.Kind == xmldom.Element {
				contexts = append(contexts, n)
			}
		}

		rel, err := filepath.Rel(filepath.Join(corpustest.Path(t, "."), fixtureDir), fixture)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&buf, "F %s nodes=%d contexts=%d\n",
			filepath.ToSlash(rel), len(all), len(contexts))
		for _, n := range all {
			fmt.Fprintf(&buf, "N %s %s\n", katPath(n), katLabel(n))
		}

		for i, e := range exprs {
			var body bytes.Buffer
			hits := 0
			for _, ctx := range contexts {
				nodes, err := compiled[i].Evaluate(ctx)
				if err != nil {
					t.Fatalf("%s: %q at %s: %v", rel, e, katPath(ctx), err)
				}
				if len(nodes) == 0 {
					continue
				}
				hits++
				fmt.Fprintf(&body, "= %s -> %s\n", katPath(ctx), katPaths(nodes))
			}
			fmt.Fprintf(&buf, "X hits=%d %s\n", hits, e)
			buf.Write(body.Bytes())
		}
	}
	return buf.Bytes()
}

func katPaths(nodes []*xmldom.Node) string {
	var sb strings.Builder
	for i, n := range nodes {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(katPath(n))
	}
	return sb.String()
}

// katPath is the oracle's node path: the document is "/", any other node is its parent's path
// followed by its 1-based index among ALL of the parent's child nodes, and an attribute is its
// owner element's path followed by "/@" and the qualified name.
func katPath(n *xmldom.Node) string {
	switch n.Kind {
	case xmldom.Document:
		return "/"
	case xmldom.Attribute:
		return katPath(n.Parent) + "/@" + n.Name.QName()
	}
	prefix := ""
	if n.Parent.Kind != xmldom.Document {
		prefix = katPath(n.Parent)
	}
	index := 1
	for c := n.Parent.FirstChild; c != nil && c != n; c = c.NextSibling {
		index++
	}
	return prefix + "/" + strconv.Itoa(index)
}

func katLabel(n *xmldom.Node) string {
	switch n.Kind {
	case xmldom.Element:
		return n.Name.QName()
	case xmldom.Text:
		return "#text"
	case xmldom.CDATA:
		return "#cdata"
	case xmldom.Comment:
		return "#comment"
	case xmldom.ProcInst:
		return "?" + n.Name.Local
	}
	return "?"
}

// ------------------------------------------------------------------ testdata loading

// loadExpressions reads the inventory: every line that is neither blank nor a comment, in
// file order, which is the order the oracle compiled them in.
func loadExpressions(t *testing.T, file string) []string {
	t.Helper()
	var out []string
	for _, line := range readLines(t, filepath.Join(corpustest.Path(t, "."), file)) {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}
	if len(out) == 0 {
		t.Fatalf("%s is empty", file)
	}
	return out
}

func loadNamespaces(t *testing.T) NamespaceContext {
	t.Helper()
	ns := NamespaceContext{}
	for _, line := range readLines(t, filepath.Join(corpustest.Path(t, "."), "namespaces.txt")) {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		prefix, uri, ok := strings.Cut(line, "\t")
		if !ok {
			t.Fatalf("namespaces.txt: no tab in %q", line)
		}
		ns.RegisterNamespace(prefix, uri)
	}
	return ns
}

// fixtureFiles lists the fixtures in the order the oracle walked them, which is Java's
// Path.compareTo over the paths - byte order over the path strings.
func fixtureFiles(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	err := filepath.Walk(filepath.Join(corpustest.Path(t, "."), dir),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".xml") {
				out = append(out, path)
			}
			return nil
		})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(out)
	if len(out) == 0 {
		t.Fatalf("no fixtures under testdata/%s", dir)
	}
	return out
}

// katHeader returns the two comment lines the oracle writes, taken from the golden file so the
// JDK version it records does not have to be duplicated here.
func katHeader(t *testing.T, katFile string) []string {
	t.Helper()
	var out []string
	for _, line := range readLines(t, filepath.Join(corpustest.Path(t, "."), katFile)) {
		if !strings.HasPrefix(line, "#") {
			break
		}
		out = append(out, line)
	}
	return out
}

func readLines(t *testing.T, path string) []string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out []string
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for s.Scan() {
		out = append(out, s.Text())
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

// firstDiff reports the first differing line with a little context, which is far more useful
// than a three-megabyte unified diff.
func firstDiff(got, want []byte) string {
	g := strings.Split(string(got), "\n")
	w := strings.Split(string(want), "\n")
	for i := 0; i < len(g) && i < len(w); i++ {
		if g[i] == w[i] {
			continue
		}
		var sb strings.Builder
		for j := max(0, i-3); j < i; j++ {
			fmt.Fprintf(&sb, "  %s\n", w[j])
		}
		fmt.Fprintf(&sb, "line %d\n  java: %s\n  go:   %s\n", i+1, w[i], g[i])
		return sb.String()
	}
	return fmt.Sprintf("line counts differ: go %d, java %d", len(g), len(w))
}
