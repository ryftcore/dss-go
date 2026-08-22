package jaxb

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/corpustest"
)

// node is a generic XML element, enough to walk 1910202xmlSchema.xsd.
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

func loadValidationReportSchema(t *testing.T) node {
	t.Helper()
	data, err := os.ReadFile(corpustest.Path(t, filepath.Join("xsd", "1910202xmlSchema.xsd")))
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

// crossNamespaceNames are element and attribute names 1910202xmlSchema.xsd
// references via xs:element ref="ds:..." (or, for Algorithm, an attribute
// of such a referenced type) without declaring locally - they live in the
// XMLDSig/XAdES/trusted-list schemas this phase's manifest does not include
// (see doc.go's "Cross-namespace types"), so the model binds them (as
// DigestMethodType, DSDigestValue, SignatureValueType, SignatureType and
// their fields) without them ever appearing as a top-level xs:element/
// xs:attribute declaration in this schema file.
var crossNamespaceNames = map[string]bool{
	"DigestMethod":   true,
	"DigestValue":    true,
	"SignatureValue": true,
	"Signature":      true,
	"Algorithm":      true,
	"Id":             true, // ds:SignatureValueType's own optional Id attribute
}

// choiceTables are the xs:choice groups this package marshals by hand (see
// jaxb_sa_attributes.go's header); modelNames walks their element name
// tables the way it walks ordinary struct tags.
var choiceTables = [][]choiceElement{
	signatureAttributeElements,
	signersDocumentElements,
	validationObjectRepresentationElements,
}

// modelNames returns every element name and every attribute name the Go
// model binds, walking the registry, the choice tables and SARevIDListType's
// hand-written CRLID/OCSPID choice.
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
		if t == reflect.TypeOf(SARevIDListType{}) {
			elements["CRLID"] = true
			elements["OCSPID"] = true
			visit(reflect.TypeOf(SACRLIDType{}))
			visit(reflect.TypeOf(SAOCSPIDType{}))
		}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" && !f.Anonymous {
				continue // unexported (e.g. ValidationReportType.extraNamespaces)
			}
			if f.Anonymous && f.Type.Kind() == reflect.Struct {
				visit(f.Type)
				continue
			}
			if f.Type == reflect.TypeOf([]choiceItem(nil)) {
				for _, table := range choiceTables {
					for _, e := range table {
						elements[e.name] = true
						if e.typ.Kind() == reflect.Struct {
							visit(e.typ)
						}
					}
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
	// The XMLDSig-backed leaf types hard-code their own element names rather
	// than a field tag (see jaxb_crossns.go's header); modelNames must record
	// them too, matched against crossNamespaceNames rather than a schema
	// declaration.
	elements["DigestMethod"] = true
	elements["DigestValue"] = true
	elements["SignatureValue"] = true
	elements["Signature"] = true
	attributes["Algorithm"] = true
	return elements, attributes
}

// TestSchemaNamesCovered is the XSD-completeness sweep: every element and
// every attribute 1910202xmlSchema.xsd declares locally must be bound
// somewhere in the model, and the model must not bind a name the schema
// does not declare - except the handful of cross-namespace names in
// crossNamespaceNames, which the schema references but does not declare.
func TestSchemaNamesCovered(t *testing.T) {
	root := loadValidationReportSchema(t)
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
			if !tc.xsd[name] && !crossNamespaceNames[name] {
				extra = append(extra, name)
			}
		}
		sort.Strings(missing)
		sort.Strings(extra)
		if len(missing) > 0 {
			t.Errorf("%d %s names declared by 1910202xmlSchema.xsd are not bound by the model: %v",
				len(missing), tc.kind, missing)
		}
		if len(extra) > 0 {
			t.Errorf("%d %s names bound by the model are not declared by 1910202xmlSchema.xsd: %v",
				len(extra), tc.kind, extra)
		}
	}
}
