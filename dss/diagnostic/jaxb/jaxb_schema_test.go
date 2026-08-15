package jaxb

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// node is a generic XML element, enough to walk DiagnosticData.xsd.
type node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Children []node     `xml:",any"`
}

func (n node) attr(name string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

func loadSchema(t *testing.T) node {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "xsd", "DiagnosticData.xsd"))
	if err != nil {
		t.Fatal(err)
	}
	var root node
	if err := xml.Unmarshal(data, &root); err != nil {
		t.Fatal(err)
	}
	return root
}

func walkSchema(n node, visit func(node)) {
	visit(n)
	for _, c := range n.Children {
		walkSchema(c, visit)
	}
}

// modelNames returns every element name and every attribute name the Go model
// binds, walking the registry and the wrapper structs it reaches.
func modelNames() (elements, attributes map[string]bool) {
	elements = map[string]bool{}
	attributes = map[string]bool{}
	seen := map[reflect.Type]bool{}
	var visit func(t reflect.Type)
	visit = func(t reflect.Type) {
		if seen[t] {
			return
		}
		seen[t] = true
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Anonymous && f.Type.Kind() == reflect.Struct {
				visit(f.Type)
				continue
			}
			ift := f.Type
			for ift.Kind() == reflect.Pointer || ift.Kind() == reflect.Slice {
				ift = ift.Elem()
			}
			if ift.Kind() == reflect.Interface {
				for _, e := range certificateExtensionElements {
					elements[e.name] = true
					visit(e.typ)
				}
				continue
			}
			parts := strings.Split(f.Tag.Get("xml"), ",")
			name := parts[0]
			if idx := strings.LastIndex(name, " "); idx >= 0 {
				name = name[idx+1:]
			}
			isAttr := false
			for _, opt := range parts[1:] {
				if opt == "attr" {
					isAttr = true
				}
			}
			et := f.Type
			for et.Kind() == reflect.Pointer || et.Kind() == reflect.Slice {
				et = et.Elem()
			}
			if name != "" && name != "-" {
				if isAttr {
					attributes[name] = true
				} else {
					elements[name] = true
				}
			}
			if et.Kind() == reflect.Struct && !implementsTextCodec(f.Type) {
				visit(et)
			}
		}
	}
	for _, t := range modelTypes {
		visit(t)
	}
	return elements, attributes
}

// TestSchemaNamesCovered is the XSD-completeness sweep: every element and every
// attribute declared by DiagnosticData.xsd must be bound somewhere in the model,
// and the model must not bind a name the schema does not declare.
func TestSchemaNamesCovered(t *testing.T) {
	root := loadSchema(t)
	xsdElements := map[string]bool{}
	xsdAttributes := map[string]bool{}
	walkSchema(root, func(n node) {
		name := n.attr("name")
		if name == "" {
			return
		}
		switch n.XMLName.Local {
		case "element":
			xsdElements[name] = true
		case "attribute":
			xsdAttributes[name] = true
		}
	})
	if len(xsdElements) == 0 || len(xsdAttributes) == 0 {
		t.Fatal("schema parse produced no declarations")
	}

	goElements, goAttributes := modelNames()
	for _, tc := range []struct {
		kind     string
		xsd, mdl map[string]bool
	}{
		{"element", xsdElements, goElements},
		{"attribute", xsdAttributes, goAttributes},
	} {
		var missing, extra []string
		for name := range tc.xsd {
			if !tc.mdl[name] {
				missing = append(missing, name)
			}
		}
		for name := range tc.mdl {
			if !tc.xsd[name] {
				extra = append(extra, name)
			}
		}
		sort.Strings(missing)
		sort.Strings(extra)
		if len(missing) > 0 {
			t.Errorf("%d %s names declared by DiagnosticData.xsd are not bound by the model: %v",
				len(missing), tc.kind, missing)
		}
		if len(extra) > 0 {
			t.Errorf("%d %s names bound by the model are not declared by DiagnosticData.xsd: %v",
				len(extra), tc.kind, extra)
		}
	}
}

// schemaChildren returns the child element names a complexType declares, in
// schema order, following xs:extension into the base type. Nested inline
// complexTypes are not descended into: they belong to the child element.
func schemaChildren(types map[string]node, n node) []string {
	var out []string
	var rec func(n node, top bool)
	rec = func(n node, top bool) {
		for _, c := range n.Children {
			switch c.XMLName.Local {
			case "element":
				out = append(out, c.attr("name"))
			case "sequence", "choice", "all", "complexContent":
				rec(c, false)
			case "extension":
				if base, ok := types[c.attr("base")]; ok {
					rec(base, true)
				}
				rec(c, false)
			}
		}
	}
	rec(n, true)
	return out
}

// schemaAttributes returns the attribute names a complexType declares,
// following xs:extension and xs:simpleContent into the base type.
func schemaAttributes(types map[string]node, n node) []string {
	var out []string
	var rec func(n node)
	rec = func(n node) {
		for _, c := range n.Children {
			switch c.XMLName.Local {
			case "attribute":
				out = append(out, c.attr("name"))
			case "sequence", "choice", "all", "complexContent", "simpleContent":
				rec(c)
			case "extension":
				if base, ok := types[c.attr("base")]; ok {
					rec(base)
				}
				rec(c)
			}
		}
	}
	rec(n)
	sort.Strings(out)
	return out
}

// TestSchemaComplexTypesMatchModel checks each named complexType of the schema
// against the model type that mirrors it: same child elements in the same
// order (which is what makes the marshalled sequence match), same attributes.
func TestSchemaComplexTypesMatchModel(t *testing.T) {
	root := loadSchema(t)
	types := map[string]node{}
	walkSchema(root, func(n node) {
		if n.XMLName.Local == "complexType" && n.attr("name") != "" {
			types[n.attr("name")] = n
		}
	})
	byName := map[string]reflect.Type{}
	for _, mt := range modelTypes {
		byName[mt.Name()] = mt
	}
	checked := 0
	for name, ct := range types {
		gt, ok := byName["Xml"+name]
		if !ok {
			continue // a wrapper type xjc folded into an @XmlElementWrapper
		}
		checked++
		want := schemaChildren(types, ct)
		var got []string
		for _, f := range elementFields(gt) {
			got = append(got, f.name)
		}
		if !equalStrings(want, got) {
			t.Errorf("%s: schema children %v, model children %v", gt.Name(), want, got)
		}
		wantAttrs := schemaAttributes(types, ct)
		gotAttrs := modelAttributes(gt)
		if !equalStrings(wantAttrs, gotAttrs) {
			t.Errorf("%s: schema attributes %v, model attributes %v", gt.Name(), wantAttrs, gotAttrs)
		}
	}
	if checked < 130 {
		t.Fatalf("only %d complexTypes matched a model type", checked)
	}
}

// modelAttributes returns the sorted attribute names a model struct binds.
func modelAttributes(t reflect.Type) []string {
	var out []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			out = append(out, modelAttributes(f.Type)...)
			continue
		}
		parts := strings.Split(f.Tag.Get("xml"), ",")
		for _, opt := range parts[1:] {
			if opt == "attr" && parts[0] != "" && parts[0] != "-" {
				out = append(out, parts[0])
			}
		}
	}
	sort.Strings(out)
	return out
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// schemaChildNodes returns the child element declarations of a complexType in
// schema order, following xs:extension into the base type.
func schemaChildNodes(types map[string]node, n node) []node {
	var out []node
	var rec func(n node)
	rec = func(n node) {
		for _, c := range n.Children {
			switch c.XMLName.Local {
			case "element":
				out = append(out, c)
			case "sequence", "choice", "all", "complexContent":
				rec(c)
			case "extension":
				if base, ok := types[c.attr("base")]; ok {
					rec(base)
				}
				rec(c)
			}
		}
	}
	rec(n)
	return out
}

// TestSchemaTreeMatchesModel walks the schema from the document element down,
// pairing every element declaration with the Go field bound to it. Unlike the
// per-complexType check this also covers the anonymous and wrapper types xjc
// folds away - the element wrappers, whose child element name is what ends up
// in the document.
func TestSchemaTreeMatchesModel(t *testing.T) {
	root := loadSchema(t)
	types := map[string]node{}
	var docElement node
	walkSchema(root, func(n node) {
		if n.XMLName.Local == "complexType" && n.attr("name") != "" {
			types[n.attr("name")] = n
		}
	})
	for _, c := range root.Children {
		if c.XMLName.Local == "element" && c.attr("name") == rootElement {
			docElement = c
		}
	}
	if docElement.attr("name") == "" {
		t.Fatal("no document element in the schema")
	}

	// complexTypeOf resolves the type of an element declaration: a named type,
	// an inline one, or none for a simple type.
	// The second result is the type name, empty for an inline complexType.
	complexTypeOf := func(el node) (node, string, bool) {
		if name := el.attr("type"); name != "" {
			ct, ok := types[name]
			return ct, name, ok
		}
		for _, c := range el.Children {
			if c.XMLName.Local == "complexType" {
				return c, "", true
			}
		}
		return node{}, "", false
	}

	visited := map[string]bool{}
	pairs := 0
	var check func(path, typeName string, ct node, gt reflect.Type)
	check = func(path, typeName string, ct node, gt reflect.Type) {
		// A named type is visited once per model type; an inline one exists at a
		// single place in the schema, so its path identifies it. Keying on the
		// type rather than the path is what terminates on a recursive schema.
		key := typeName
		if key == "" {
			key = path
		}
		key += "|" + gt.String()
		if visited[key] {
			return
		}
		visited[key] = true
		children := schemaChildNodes(types, ct)
		fields := elementFields(gt)
		if len(children) != len(fields) {
			var want, got []string
			for _, c := range children {
				want = append(want, c.attr("name"))
			}
			for _, f := range fields {
				got = append(got, f.name)
			}
			t.Errorf("%s (%s): schema declares %v, model binds %v", path, gt.Name(), want, got)
			return
		}
		for i, c := range children {
			name := c.attr("name")
			f := fields[i]
			if f.name != name {
				t.Errorf("%s (%s): child %d is %s in the schema and %s in the model",
					path, gt.Name(), i, name, f.name)
				continue
			}
			pairs++
			sub, subName, complex := complexTypeOf(c)
			if !complex {
				if f.typ != nil && !hasCharData(f.typ) {
					t.Errorf("%s/%s: simple type in the schema, %s in the model", path, name, f.typ)
				}
				continue
			}
			if f.typ == nil {
				t.Errorf("%s/%s: complexType in the schema, simple binding in the model", path, name)
				continue
			}
			check(path+"/"+name, subName, sub, f.typ)
		}
	}
	ct, ctName, ok := complexTypeOf(docElement)
	if !ok {
		t.Fatal("the document element has no complexType")
	}
	check(rootElement, ctName, ct, reflect.TypeOf(XmlDiagnosticData{}))
	if pairs < 400 {
		t.Fatalf("only %d element declarations paired with a model field", pairs)
	}
}

// TestOracleCorpusExercisesModel is the counterpart of the schema sweeps: they
// prove the model binds exactly the names DiagnosticData.xsd declares, this one
// proves the marshal-parity corpus actually puts every one of those names
// through the round trip. Without it a binding could be wrong in a way no dump
// would ever reveal - the schema sweep only reads the struct tags, it never
// marshals them.
func TestOracleCorpusExercisesModel(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(oracleDir(), "*.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatalf("no oracle dumps under %s", oracleDir())
	}
	seenElements := map[string]bool{}
	seenAttributes := map[string]bool{}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		dec := xml.NewDecoder(bytes.NewReader(data))
		for {
			tok, err := dec.Token()
			if err != nil {
				break
			}
			start, ok := tok.(xml.StartElement)
			if !ok {
				continue
			}
			seenElements[start.Name.Local] = true
			for _, a := range start.Attr {
				if a.Name.Space == "xmlns" || a.Name.Local == "xmlns" {
					continue
				}
				seenAttributes[a.Name.Local] = true
			}
		}
	}

	elements, attributes := modelNames()
	for _, tc := range []struct {
		kind      string
		mdl, seen map[string]bool
	}{
		{"element", elements, seenElements},
		{"attribute", attributes, seenAttributes},
	} {
		var missing []string
		for name := range tc.mdl {
			if !tc.seen[name] {
				missing = append(missing, name)
			}
		}
		sort.Strings(missing)
		if len(missing) > 0 {
			t.Errorf("%d %s names the model binds never appear in the oracle corpus, "+
				"so the round trip never exercises them: %v", len(missing), tc.kind, missing)
		}
	}
}
