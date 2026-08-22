// Ported from 1910202xmlSchema.xsd (DSS 6.5.RC1): the content model the
// marshaller needs to reproduce one JAXB spelling that encoding/xml cannot
// express on its own. This is the same mechanism dss/diagnostic/jaxb and the
// rest of the generated-JAXB packages use; see
// dss/diagnostic/jaxb/jaxb_content_model.go for the full rationale.
//
// Three schema elements bind to a plain xs:anyURI/xs:base64Binary inside the
// hand-driven <xs:choice> of ValidationObjectRepresentationType
// (jaxb_validation_object.go's choiceWrapper): "base64" and "URI" carry
// character data (JAXB writes <base64></base64> for a present-but-empty
// value, never <base64/>), unlike every other choice member in this schema,
// which is always a complexType and so always defaults correctly to the
// self-closing form. Because the choice wrapper's MarshalXML/UnmarshalXML
// bypasses the struct-tag reflection this file's buildContentModel walks,
// those two names are seeded by hand rather than discovered.
package jaxb

import (
	"encoding/xml"
	"reflect"
	"strings"
)

// elemKind is how the RI writes an element of this content model when it is
// empty.
type elemKind uint8

const (
	// kindComplex elements are written <X/>.
	kindComplex elemKind = iota
	// kindText elements are written <X></X>.
	kindText
)

// rootElement is the name of the document element; it has no enclosing element.
const rootElement = "ValidationReport"

var (
	// contentKind maps an element name to its content model wherever that is
	// unambiguous across the whole schema.
	contentKind = map[string]elemKind{}
	// contentKindByParent resolves a name by the element enclosing it.
	contentKindByParent = map[string]map[string]elemKind{}
)

func init() {
	buildContentModel()
	// See this file's header: bypasses struct-tag reflection via a hand-driven choice.
	contentKind["base64"] = kindText
	contentKind["URI"] = kindText
}

// contentField is one element binding of a model struct.
type contentField struct {
	name string
	kind elemKind
	typ  reflect.Type // the struct behind a complex binding, nil otherwise
}

// elementFields returns the element bindings of a struct type in schema
// order, flattening the embedded structs that carry inherited content
// (AttributeBaseType's SA* subtypes).
func elementFields(t reflect.Type) []contentField {
	var out []contentField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			out = append(out, elementFields(f.Type)...)
			continue
		}
		if f.Type == reflect.TypeOf(xml.Name{}) {
			// XMLName names the document element, it is not content.
			continue
		}
		tag := f.Tag.Get("xml")
		parts := strings.Split(tag, ",")
		name := parts[0]
		if name == "" || name == "-" {
			continue
		}
		skip := false
		for _, opt := range parts[1:] {
			if opt == "attr" || opt == "chardata" || opt == "innerxml" || opt == "comment" {
				skip = true
			}
		}
		if skip {
			continue
		}
		if idx := strings.LastIndex(name, " "); idx >= 0 {
			name = name[idx+1:] // a namespaced tag: "<ns> <local>"
		}
		et := f.Type
		for et.Kind() == reflect.Pointer || et.Kind() == reflect.Slice {
			et = et.Elem()
		}
		if et.Kind() == reflect.Struct && !implementsTextCodec(f.Type) {
			kind := kindComplex
			if hasCharData(et) {
				kind = kindText
			}
			out = append(out, contentField{name: name, kind: kind, typ: et})
		} else {
			out = append(out, contentField{name: name, kind: kindText})
		}
	}
	return out
}

// hasCharData reports whether a struct binds a simpleContent value, i.e.
// writes character data of its own.
func hasCharData(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct && hasCharData(f.Type) {
			return true
		}
		for _, opt := range strings.Split(f.Tag.Get("xml"), ",")[1:] {
			if opt == "chardata" {
				return true
			}
		}
	}
	return false
}

// implementsTextCodec reports whether a type marshals itself as character
// data (XSDateTime, Base64Binary, the URI-adapted enum wrappers are structs
// or defined types but bind to a simpleType, not to a complexType).
func implementsTextCodec(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return reflect.PointerTo(t).Implements(textMarshalerType) || t.Implements(textMarshalerType)
}

var textMarshalerType = reflect.TypeOf((*interface {
	MarshalText() ([]byte, error)
})(nil)).Elem()

// buildContentModel walks the model types, records under which element names
// each type appears, and from that derives the two lookup tables.
func buildContentModel() {
	fields := map[reflect.Type][]contentField{}
	var visit func(t reflect.Type)
	visit = func(t reflect.Type) {
		if _, seen := fields[t]; seen {
			return
		}
		fs := elementFields(t)
		fields[t] = fs
		for _, f := range fs {
			if f.typ != nil {
				visit(f.typ)
			}
		}
	}
	for _, t := range modelTypes {
		visit(t)
	}

	parents := map[reflect.Type]map[string]bool{}
	parents[reflect.TypeOf(ValidationReportType{})] = map[string]bool{rootElement: true}
	for _, fs := range fields {
		for _, f := range fs {
			if f.typ == nil {
				continue
			}
			if parents[f.typ] == nil {
				parents[f.typ] = map[string]bool{}
			}
			parents[f.typ][f.name] = true
		}
	}

	seen := map[string]map[elemKind]bool{}
	for t, fs := range fields {
		for _, f := range fs {
			if seen[f.name] == nil {
				seen[f.name] = map[elemKind]bool{}
			}
			seen[f.name][f.kind] = true
			for parent := range parents[t] {
				if contentKindByParent[parent] == nil {
					contentKindByParent[parent] = map[string]elemKind{}
				}
				contentKindByParent[parent][f.name] = f.kind
			}
		}
	}
	for name, kinds := range seen {
		if len(kinds) == 1 {
			for k := range kinds {
				contentKind[name] = k
			}
		}
	}
	contentKind[rootElement] = kindComplex
}

// carriesCharData reports whether an empty <name> element would have been
// written by JAXB as <name></name> rather than <name/>. stack holds the
// enclosing element names, innermost last. name and the top of stack may
// carry a literal "ns2:"-style prefix (see jaxb_crossns.go's header on
// hard-coded ds: element names): the model's contentKind/contentKindByParent
// tables are keyed by local name only, as derived from Go struct tags which
// never carry one, so both are stripped before lookup.
func carriesCharData(name string, stack []string) bool {
	local := localName(name)
	// ds:SignatureValue is bound to a byte[]; when the signature value is
	// absent the JAXB RI writes <ns2:SignatureValue/>, not the <x></x> an
	// empty text element would otherwise get here. Verified against the
	// upstream ETSI-VR oracle corpus, in which every empty element - this one
	// included - is self-closed. Not exercised by the KAT corpus: no fixture
	// produces an empty signature value.
	if local == "SignatureValue" {
		return false
	}
	if len(stack) > 0 {
		if k, ok := contentKindByParent[localName(stack[len(stack)-1])][local]; ok {
			return k == kindText
		}
	}
	return contentKind[local] == kindText
}

// localName strips a literal "prefix:" from a tag name written by hand
// (jaxb_crossns.go's ns2:-prefixed elements); tag names encoding/xml derives
// from Go struct tags never carry one.
func localName(name string) string {
	if i := strings.IndexByte(name, ':'); i >= 0 {
		return name[i+1:]
	}
	return name
}
