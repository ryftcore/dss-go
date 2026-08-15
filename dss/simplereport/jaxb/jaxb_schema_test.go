package jaxb

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestSchemaElementsAndAttributesAreBound is the XSD-completeness sweep:
// every xs:element/xs:attribute name declared in SimpleReport.xsd must
// appear as an xml struct tag somewhere in this package's *.go sources (the
// generated model), so no schema addition can silently go unbound. It is a
// text sweep over the package's own source, not a schema validator: it
// would not catch a binding to the wrong element under the wrong parent,
// but it does catch a name the model never mentions at all - the class of
// drift most likely to happen when SimpleReport.xsd changes upstream.
func TestSchemaElementsAndAttributesAreBound(t *testing.T) {
	xsd, err := os.ReadFile(filepath.Join("testdata", "xsd", "SimpleReport.xsd"))
	if err != nil {
		t.Fatal(err)
	}
	names := schemaNames(t, xsd)

	src := packageSource(t)

	var missing []string
	for _, name := range names {
		if !strings.Contains(src, `"`+name+`,`) && !strings.Contains(src, `"`+name+`"`) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		t.Errorf("SimpleReport.xsd names with no xml struct tag in the package: %v", missing)
	}
}

// schemaNames extracts every xs:element/xs:attribute name attribute from
// the schema, deduplicated.
func schemaNames(t *testing.T, xsd []byte) []string {
	t.Helper()
	elemRe := regexp.MustCompile(`<xs:element\s+name="([^"]+)"`)
	attrRe := regexp.MustCompile(`<xs:attribute\s+name="([^"]+)"`)
	seen := map[string]bool{}
	var names []string
	for _, m := range elemRe.FindAllSubmatch(xsd, -1) {
		n := string(m[1])
		if !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}
	for _, m := range attrRe.FindAllSubmatch(xsd, -1) {
		n := string(m[1])
		if !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}
	return names
}

// packageSource concatenates every non-test *.go file of this package, the
// text the sweep searches.
func packageSource(t *testing.T) string {
	t.Helper()
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		sb.Write(data)
		sb.WriteByte('\n')
	}
	return sb.String()
}
