// Ported from DetailedReport.xsd (DSS 6.5.RC1): the content model the
// marshaller needs to reproduce one JAXB spelling that encoding/xml cannot
// express on its own. This is dss/diagnostic/jaxb/jaxb_content_model.go's
// mechanism, replicated for this schema: the same pattern per package,
// rather than inventing a new marshalling approach.
//
// The JAXB RI writes an empty element as <X/> when nothing at all was written
// inside it, but as <X></X> once character data - even the empty string - was
// written. encoding/xml always writes the pair form, so Marshal collapses empty
// pairs to the self-closing form for the elements whose content is complex.
//
// Which elements those are is not hand-listed: it is derived by reflection from
// the model types themselves (see modelTypes in jaxb_model.go). An element is
// complex when the Go field bound to it is a struct - exactly the schema's
// complexType/simpleType split.
//
// Two constructs need a hand-written entry rather than a struct-tag walk,
// because the generated Java classes route them through custom code the
// reflection can't see:
//
//   - DetailedReport's top-level choice group
//     (Signature|Timestamp|EvidenceRecord|EAA|Certificate) is bound to the
//     XmlReportItem interface (see jaxb_report.go); reportItemTypes supplies
//     the element-name/type pairs the way certificateExtensionElements does in
//     dss/diagnostic/jaxb.
//   - XmlCertificateApprovalStatus (jaxb_qualification.go) implements
//     xml.Marshaler/xml.Unmarshaler by hand because its three elements are
//     bound through loader-backed interface adapters with no defined
//     underlying type; certificateApprovalStatusFields supplies its content
//     model directly since its fields carry no xml struct tags to walk.

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
const rootElement = "DetailedReport"

// reportItemTypes is the @XmlElements table of
// DetailedReport.signatureOrTimestampOrEvidenceRecord (see jaxb_report.go's
// reportItemElements, which the marshaller itself consults), restated here as
// reflect.Type so elementFields can resolve the interface field.
var reportItemTypes = []struct {
	name string
	typ  reflect.Type
}{
	{"Signature", reflect.TypeOf(XmlSignature{})},
	{"Timestamp", reflect.TypeOf(XmlTimestamp{})},
	{"EvidenceRecord", reflect.TypeOf(XmlEvidenceRecord{})},
	{"EAA", reflect.TypeOf(XmlEAA{})},
	{"Certificate", reflect.TypeOf(XmlCertificate{})},
}

// certificateApprovalStatusFields is XmlCertificateApprovalStatus's content
// model, hand-supplied because its fields carry no xml struct tags (see the
// file header).
var certificateApprovalStatusFields = []contentField{
	{name: "ListType", kind: kindText},
	{name: "ServiceTypeIdentifier", kind: kindText},
	{name: "ServiceStatus", kind: kindText},
}

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

// elementFields returns the element bindings of a struct type in schema order,
// flattening the embedded structs that carry inherited content.
func elementFields(t reflect.Type) []contentField {
	if t == reflect.TypeOf(XmlCertificateApprovalStatus{}) {
		return certificateApprovalStatusFields
	}
	if t == reflect.TypeOf(XmlLoTEAnalysisEntry{}) {
		// The LoTEAnalysis slot's content is XmlLoTEAnalysis's, whether or not
		// this entry carries an xsi:type override (see jaxb_report.go).
		return elementFields(reflect.TypeOf(XmlLoTEAnalysis{}))
	}
	var out []contentField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			out = append(out, elementFields(f.Type)...)
			continue
		}
		ift := f.Type
		for ift.Kind() == reflect.Pointer || ift.Kind() == reflect.Slice {
			ift = ift.Elem()
		}
		if ift.Kind() == reflect.Interface {
			// The only interface-typed property is DetailedReport's top-level
			// choice group; every alternative is a complexType.
			for _, e := range reportItemTypes {
				out = append(out, contentField{name: e.name, kind: kindComplex, typ: e.typ})
			}
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
			if hasCharData(et) || et == reflect.TypeOf(XmlCertificateApprovalStatus{}) {
				kind = kindText
			}
			out = append(out, contentField{name: name, kind: kind, typ: et})
		} else {
			out = append(out, contentField{name: name, kind: kindText})
		}
	}
	return out
}

// hasCharData reports whether a struct binds a simpleContent value, i.e. writes
// character data of its own.
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

// implementsTextCodec reports whether a type marshals itself as character data
// (XSDateTime is a struct but binds to xs:dateTime, not to a complexType).
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
	// Reachability: the registry plus every struct a registered type refers to.
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

	// The element names each type is bound to.
	parents := map[reflect.Type]map[string]bool{}
	parents[reflect.TypeOf(XmlDetailedReport{})] = map[string]bool{rootElement: true}
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

// compile-time check that xml.Marshaler/Unmarshaler is what XmlDetailedReport
// and XmlCertificateApprovalStatus implement, matching the assumption
// elementFields makes about them above.
var (
	_ xml.Marshaler   = (*XmlDetailedReport)(nil)
	_ xml.Unmarshaler = (*XmlDetailedReport)(nil)
	_ xml.Marshaler   = (*XmlCertificateApprovalStatus)(nil)
	_ xml.Unmarshaler = (*XmlCertificateApprovalStatus)(nil)
)
