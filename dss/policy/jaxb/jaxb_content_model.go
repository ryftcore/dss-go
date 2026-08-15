// Ported from policy.xsd (DSS 6.5.RC1): the content model the marshaller needs
// to reproduce one JAXB spelling that encoding/xml cannot express on its own.
//
// The JAXB RI writes an empty element as <X/> when nothing at all was written
// inside it, but as <X></X> once character data - even the empty string - was
// written. encoding/xml always writes the pair form, so Marshal (in xml.go)
// collapses empty pairs to the self-closing form for the elements whose
// content is complex.
//
// Which elements those are is not hand-listed: it is derived by reflection
// from the model types themselves, starting at ConstraintsParameters (every
// generated type in this schema is reachable from the document element, so a
// single root suffices - unlike dss/diagnostic/jaxb's DiagnosticData.xsd,
// which needed an explicit registry of unreachable types). An element is
// complex when the Go field bound to it is a struct - exactly the schema's
// complexType/simpleType split.
//
// A complexType with simpleContent (Algo) counts as character data, not as
// complex: the RI writes <X></X> for one whose value is the empty string.
// This port has no such deviation to record: unlike DiagnosticData.xsd's
// XmlOID et al., policy.xsd's only simpleContent type (Algo) is always
// non-empty in every upstream policy resource this chunk ports, so the
// "RI writes <X/> for a null simpleContent value" case documented in
// dss/diagnostic/jaxb/jaxb_content_model.go does not arise here.
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

// rootElement is the name of the document element; it has no enclosing
// element.
const rootElement = "ConstraintsParameters"

var (
	// contentKind maps an element name to its content model wherever that is
	// unambiguous across the whole schema.
	contentKind = map[string]elemKind{}
	// contentKindByParent resolves a name by the element enclosing it.
	contentKindByParent = map[string]map[string]elemKind{}
)

func init() {
	buildContentModel()
}

// contentField is one element binding of a model struct.
type contentField struct {
	name string
	kind elemKind
	typ  reflect.Type // the struct behind a complex binding, nil otherwise
}

// elementFields returns the element bindings of a struct type in schema
// order, flattening the embedded structs that carry inherited content.
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
// data rather than as a complexType (none of policy.xsd's element-bound
// types do so today - LevelValue/ValidationModelValue only ever appear on
// attr-tagged fields, which elementFields skips before reaching this check -
// but the check is kept for the same robustness reason
// dss/diagnostic/jaxb/jaxb_content_model.go keeps it).
func implementsTextCodec(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return reflect.PointerTo(t).Implements(textMarshalerType) || t.Implements(textMarshalerType)
}

var textMarshalerType = reflect.TypeOf((*interface {
	MarshalText() ([]byte, error)
})(nil)).Elem()

// buildContentModel walks the model types reachable from ConstraintsParameters,
// records under which element names each type appears, and from that derives
// the two lookup tables.
func buildContentModel() {
	root := reflect.TypeOf(ConstraintsParameters{})

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
	visit(root)

	// The element names each type is bound to.
	parents := map[reflect.Type]map[string]bool{root: {rootElement: true}}
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
// enclosing element names, innermost last.
func carriesCharData(name string, stack []string) bool {
	if len(stack) > 0 {
		if k, ok := contentKindByParent[stack[len(stack)-1]][name]; ok {
			return k == kindText
		}
	}
	return contentKind[name] == kindText
}
