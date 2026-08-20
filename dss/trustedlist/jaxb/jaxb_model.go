// Ported from ts_119612v020401_xsd.xsd, ts_119612v020401_sie_xsd.xsd,
// ts_119612v020401_additionaltypes_xsd.xsd and
// mra_schema_v2_19612v020401.xsd (DSS 6.5.RC1): the registry of every
// generated model type across the four Java packages this module ports,
// and the content model the marshaller needs to reproduce one JAXB
// spelling encoding/xml cannot express on its own - see
// dss/diagnostic/jaxb/jaxb_content_model.go for the fuller rationale and
// dss/diagnostic/jaxb/jaxb_model.go for the registry's twin purpose (which
// this file's TestSchemaCompleteness sweep, in jaxb_schema_test.go, also
// relies on).
//
// The JAXB RI writes an empty element as <X/> when nothing at all was
// written inside it, but as <X></X> once character data - even the empty
// string - was written. encoding/xml always writes the pair form, so
// Marshal collapses empty pairs to the self-closing form for the elements
// whose content is complex - derived, as in the diagnostic-data precedent,
// by reflection over the tagged struct fields in modelTypes.
//
// # Custom-marshaled elements
//
// AnyType, ExtensionType, the raw-capture stand-ins (ObjectIdentifierType,
// dsigSignature, dsigKeyValue, xadesAnyType - see jaxb_common.go) and every
// element wildcardElements dispatches to by a bare literal name rather than
// a tagged struct field (Qualifications, ExpiredCertsRevocationInfo,
// ExtendedKeyUsage, TakenOverBy, CertSubjectDNAttribute,
// MutualRecognitionAgreementInformation) are invisible to the tagged-field
// reflection walk: their content is driven dynamically, not by a fixed set
// of xml-tagged Go fields, and the wildcard names are never spelled in a
// struct tag at all. Each has an unambiguous, fixed content kind regardless
// (kindComplex for the six structured ones - none has simpleContent in its
// XSD type; kindText for ExpiredCertsRevocationInfo, a bare xs:dateTime
// leaf; kindText for TextualInformation, AdditionalInformationType's other
// choice member, MultiLangStringType), so manualContentKind supplies them
// directly rather than teaching elementFields to trace dynamic dispatch.
package jaxb

import (
	"encoding/xml"
	"reflect"
	"strings"
)

// modelTypes is the set of generated model types, one per schema
// complexType (or, for AdditionalInformationItem, the choice-wrapper
// pattern dss/diagnostic/jaxb's CertificateExtensionsWrapper also uses),
// across all four packages.
var modelTypes = []reflect.Type{
	// tsl
	reflect.TypeOf(MultiLangStringType{}),
	reflect.TypeOf(MultiLangNormStringType{}),
	reflect.TypeOf(InternationalNamesType{}),
	reflect.TypeOf(AttributedNonEmptyURIType{}),
	reflect.TypeOf(NonEmptyMultiLangURIType{}),
	reflect.TypeOf(NonEmptyMultiLangURIListType{}),
	reflect.TypeOf(NonEmptyURIListType{}),
	reflect.TypeOf(PostalAddressType{}),
	reflect.TypeOf(PostalAddressListType{}),
	reflect.TypeOf(ElectronicAddressType{}),
	reflect.TypeOf(AddressType{}),
	reflect.TypeOf(PolicyOrLegalnoticeType{}),
	reflect.TypeOf(NextUpdateType{}),
	reflect.TypeOf(OtherTSLPointerType{}),
	reflect.TypeOf(OtherTSLPointersType{}),
	reflect.TypeOf(AdditionalInformationType{}),
	reflect.TypeOf(AdditionalServiceInformationType{}),
	reflect.TypeOf(ExtensionsListType{}),
	reflect.TypeOf(ExtensionType{}),
	reflect.TypeOf(AnyType{}),
	reflect.TypeOf(ServiceSupplyPointsType{}),
	reflect.TypeOf(DigitalIdentityType{}),
	reflect.TypeOf(DigitalIdentityListType{}),
	reflect.TypeOf(ServiceDigitalIdentityListType{}),
	reflect.TypeOf(TSPServiceInformationType{}),
	reflect.TypeOf(TSPServiceType{}),
	reflect.TypeOf(TSPServicesListType{}),
	reflect.TypeOf(ServiceHistoryInstanceType{}),
	reflect.TypeOf(ServiceHistoryType{}),
	reflect.TypeOf(TSPInformationType{}),
	reflect.TypeOf(TSPType{}),
	reflect.TypeOf(TrustServiceProviderListType{}),
	reflect.TypeOf(TSLSchemeInformationType{}),
	reflect.TypeOf(TrustStatusListType{}),

	// ecc
	reflect.TypeOf(CriteriaListType{}),
	reflect.TypeOf(KeyUsageBitType{}),
	reflect.TypeOf(KeyUsageType{}),
	reflect.TypeOf(PoliciesListType{}),
	reflect.TypeOf(QualificationElementType{}),
	reflect.TypeOf(QualificationsType{}),
	reflect.TypeOf(QualifierType{}),
	reflect.TypeOf(QualifiersType{}),

	// tslx
	reflect.TypeOf(CertSubjectDNAttributeType{}),
	reflect.TypeOf(ExtendedKeyUsageType{}),
	reflect.TypeOf(TakenOverByType{}),

	// mra
	reflect.TypeOf(CertificateContentReferenceEquivalenceType{}),
	reflect.TypeOf(CertificateContentReferencesEquivalenceListType{}),
	reflect.TypeOf(MutualRecognitionAgreementInformationType{}),
	reflect.TypeOf(QcStatementInfoType{}),
	reflect.TypeOf(QcStatementListType{}),
	reflect.TypeOf(QcStatementType{}),
	reflect.TypeOf(QualifierEquivalenceListType{}),
	reflect.TypeOf(QualifierEquivalenceType{}),
	reflect.TypeOf(TrustServiceEquivalenceHistoryInstanceType{}),
	reflect.TypeOf(TrustServiceEquivalenceHistoryType{}),
	reflect.TypeOf(TrustServiceEquivalenceInformationType{}),
	reflect.TypeOf(TrustServiceTSLQualificationExtensionEquivalenceListType{}),
	reflect.TypeOf(TrustServiceTSLQualificationExtensionNameType{}),
	reflect.TypeOf(TrustServiceTSLStatusEquivalenceListType{}),
	reflect.TypeOf(TrustServiceTSLStatusEquivalenceType{}),
	reflect.TypeOf(TrustServiceTSLStatusList{}),
	reflect.TypeOf(TrustServiceTSLTypeEquivalenceListType{}),
	reflect.TypeOf(TrustServiceTSLTypeListType{}),
	reflect.TypeOf(TrustServiceTSLTypeType{}),

	// cross-schema raw-capture stand-ins (jaxb_common.go)
	reflect.TypeOf(ObjectIdentifierType{}),
}

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
const rootElement = "TrustServiceStatusList"

var (
	// contentKind maps an element name to its content model wherever that is
	// unambiguous across the whole schema.
	contentKind = map[string]elemKind{}
	// contentKindByParent resolves a name by the element enclosing it.
	contentKindByParent = map[string]map[string]elemKind{}
)

// manualContentKind supplies the content kind of every element this
// package's tagged-field reflection walk cannot see - see this file's
// header.
var manualContentKind = map[string]elemKind{
	"TextualInformation":                    kindText,
	"Qualifications":                        kindComplex,
	"ExpiredCertsRevocationInfo":            kindText,
	"ExtendedKeyUsage":                      kindComplex,
	"TakenOverBy":                           kindComplex,
	"CertSubjectDNAttribute":                kindComplex,
	"MutualRecognitionAgreementInformation": kindComplex,
	"MimeType":                              kindText,
	"X509CertificateLocation":               kindText,
	"PublicKeyLocation":                     kindText,
	"QcStatementSet":                        kindComplex,
}

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
// data via encoding.TextMarshaler.
func implementsTextCodec(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return reflect.PointerTo(t).Implements(textMarshalerType) || t.Implements(textMarshalerType)
}

var textMarshalerType = reflect.TypeOf((*interface {
	MarshalText() ([]byte, error)
})(nil)).Elem()

// buildContentModel walks the model types, records under which element
// names each type appears, and from that derives the two lookup tables.
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
	parents[reflect.TypeOf(TrustStatusListType{})] = map[string]bool{rootElement: true}
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
	for name, kind := range manualContentKind {
		contentKind[name] = kind
	}
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
