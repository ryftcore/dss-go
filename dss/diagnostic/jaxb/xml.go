// Ported from DiagnosticData.xsd (DSS 6.5.RC1) and the marshalling behaviour of
// eu.europa.esig.dss.diagnostic.DiagnosticDataFacade (dss-jaxb-common's
// AbstractJaxbFacade with JAXB_FORMATTED_OUTPUT=true).
//
// This file holds the runtime the generated model needs: the XML Schema simple
// types JAXB binds by hand (base64Binary, dateTime, xs:list, the collapsed-string
// ID), the IDREF object graph, and the Marshal/Unmarshal entry points.
//
// # JAXB quirks encoding/xml cannot reproduce
//
// Two behaviours of the JAXB reference implementation are outside what
// encoding/xml can express, so Marshal post-processes the encoder output with
// jaxbCanonical. Both are pure XML-syntax normalisations - no information is
// added or dropped - and they are applied to the Go output only, so the Java
// bytes stay the untouched reference in the marshal-parity KAT:
//
//  1. Self-closing tags. The RI writes <X/> for an element with no children but
//     <X></X> once characters (even the empty string) have been written.
//     encoding/xml always writes <X></X>. jaxbCanonical collapses empty pairs to
//     the self-closing form, keeping <X></X> for the element names the schema
//     binds to a possibly-empty simple type (see jaxb_content_model.go, which
//     derives those names by reflection over the model rather than listing
//     them).
//
//  2. Character escaping. encoding/xml escapes " and ' as &#34;/&#39; everywhere
//     and writes \t \n \r as numeric references; the RI leaves " and ' alone in
//     character data, writes " as &quot; in attribute values, and uses lowercase
//     hexadecimal references. jaxbCanonical rewrites the affected references
//     according to whether they sit in character data or in an attribute value.

package jaxb

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"
)

// Namespace is the target namespace of DiagnosticData.xsd. The schema declares
// elementFormDefault="qualified" and the RI binds it to the default prefix, so
// only the document element carries an xmlns declaration.
const Namespace = "http://dss.esig.europa.eu/validation/diagnostic"

// xmlDeclaration is the declaration the RI's marshaller emits.
const xmlDeclaration = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

// dateTimeFormat is the pattern of eu.europa.esig.dss.jaxb.parsers.DateParser,
// "yyyy-MM-dd'T'HH:mm:ss'Z'" evaluated in UTC.
const dateTimeFormat = "2006-01-02T15:04:05Z"

// ---------------------------------------------------------------- simple types

// Base64Binary is the Go form of a JAXB byte[] property, bound to
// xs:base64Binary. A nil Base64Binary and an empty one are distinct: JAXB omits
// a null property entirely but writes an empty element for a zero-length array.
type Base64Binary []byte

// MarshalText encodes the bytes with the standard base64 alphabet, unwrapped,
// as the RI's base64Binary printer does.
func (b Base64Binary) MarshalText() ([]byte, error) {
	out := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(out, b)
	return out, nil
}

// UnmarshalText decodes a base64Binary lexical form, ignoring the whitespace
// the schema type allows.
func (b *Base64Binary) UnmarshalText(text []byte) error {
	s := strings.Map(func(r rune) rune {
		switch r {
		case ' ', '\t', '\n', '\r':
			return -1
		}
		return r
	}, string(text))
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("invalid base64Binary value: %w", err)
	}
	*b = decoded
	return nil
}

// Bytes returns the raw bytes, tolerating a nil receiver.
func (b *Base64Binary) Bytes() []byte {
	if b == nil {
		return nil
	}
	return *b
}

// NewBase64Binary wraps bytes for a Base64Binary-typed field.
func NewBase64Binary(v []byte) *Base64Binary {
	b := Base64Binary(v)
	return &b
}

// XSDateTime is the Go form of a JAXB Date property routed through DateParser.
type XSDateTime time.Time

// MarshalText prints the date in UTC with DateParser's pattern.
func (d XSDateTime) MarshalText() ([]byte, error) {
	return []byte(time.Time(d).UTC().Format(dateTimeFormat)), nil
}

// UnmarshalText parses DateParser's pattern; any other lexical form is an error,
// mirroring the IllegalArgumentException the parser throws.
func (d *XSDateTime) UnmarshalText(text []byte) error {
	t, err := time.ParseInLocation(dateTimeFormat, string(text), time.UTC)
	if err != nil {
		return fmt.Errorf("string '%s' doesn't follow the pattern 'yyyy-MM-dd'T'HH:mm:ss'Z''", text)
	}
	*d = XSDateTime(t)
	return nil
}

// Time returns the instant, tolerating a nil receiver.
func (d *XSDateTime) Time() time.Time {
	if d == nil {
		return time.Time{}
	}
	return time.Time(*d)
}

// NewXSDateTime wraps an instant for an XSDateTime-typed field.
func NewXSDateTime(t time.Time) *XSDateTime {
	v := XSDateTime(t)
	return &v
}

// CollapsedString is the Go form of a property bound through JAXB's
// CollapsedStringAdapter: leading and trailing whitespace is stripped and inner
// whitespace runs collapse to a single space on the way in.
type CollapsedString string

// MarshalText writes the value unchanged, as the adapter's marshal() does.
func (c CollapsedString) MarshalText() ([]byte, error) { return []byte(c), nil }

// UnmarshalText applies the whitespace collapsing of CollapsedStringAdapter.
func (c *CollapsedString) UnmarshalText(text []byte) error {
	*c = CollapsedString(collapseWhitespace(string(text)))
	return nil
}

// String returns the value, tolerating a nil receiver.
func (c *CollapsedString) String() string {
	if c == nil {
		return ""
	}
	return string(*c)
}

// NewCollapsedString wraps a string for a CollapsedString-typed field.
func NewCollapsedString(v string) *CollapsedString {
	c := CollapsedString(v)
	return &c
}

func collapseWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// BigIntegerList is the Go form of a List<BigInteger> property bound to an
// xs:list of integers, printed as a whitespace-separated lexical form.
type BigIntegerList []*big.Int

// MarshalText joins the members with a single space.
func (l BigIntegerList) MarshalText() ([]byte, error) {
	parts := make([]string, 0, len(l))
	for _, v := range l {
		if v == nil {
			continue
		}
		parts = append(parts, v.String())
	}
	return []byte(strings.Join(parts, " ")), nil
}

// UnmarshalText splits the lexical form on whitespace.
func (l *BigIntegerList) UnmarshalText(text []byte) error {
	fields := strings.Fields(string(text))
	out := make(BigIntegerList, 0, len(fields))
	for _, f := range fields {
		v, ok := new(big.Int).SetString(f, 10)
		if !ok {
			return fmt.Errorf("invalid integer in xs:list: %q", f)
		}
		out = append(out, v)
	}
	*l = out
	return nil
}

// XmlEncapsulationType is the Go form of the generated JAXB enum
// XmlEncapsulationType (simpleType EncapsulationType).
type XmlEncapsulationType string

// The EncapsulationType values, with the exact lexical form of the schema.
const (
	XmlEncapsulationType_BINARIES  XmlEncapsulationType = "BINARIES"
	XmlEncapsulationType_REFERENCE XmlEncapsulationType = "REFERENCE"
)

// MarshalText writes the enum value, as the generated value() method does.
func (e XmlEncapsulationType) MarshalText() ([]byte, error) { return []byte(e), nil }

// UnmarshalText resolves the lexical form, mirroring fromValue().
func (e *XmlEncapsulationType) UnmarshalText(text []byte) error {
	switch v := XmlEncapsulationType(text); v {
	case XmlEncapsulationType_BINARIES, XmlEncapsulationType_REFERENCE:
		*e = v
		return nil
	default:
		return fmt.Errorf("no enum constant XmlEncapsulationType.%s", text)
	}
}

// ------------------------------------------------------------------ IDREF graph

// XmlToken is the Go stand-in for the abstract XmlAbstractToken/XmlTrustSourceList
// bases in an IDREF position: the set of types that carry an xs:ID and may be
// pointed at by an IDREF attribute.
type XmlToken interface {
	// TokenID returns the value of the Id attribute.
	TokenID() string
}

// TokenID returns the Id attribute of the token.
func (a *XmlAbstractTokenAttrs) TokenID() string { return a.Id.String() }

func (a *XmlAbstractTokenAttrs) setTokenID(id string) { a.Id = NewCollapsedString(id) }

// TokenID returns the Id attribute of the trust source list.
func (a *XmlTrustSourceListAttrs) TokenID() string { return a.Id.String() }

func (a *XmlTrustSourceListAttrs) setTokenID(id string) { a.Id = NewCollapsedString(id) }

// identifiable is satisfied by every type carrying an @XmlID attribute.
type identifiable interface {
	TokenID() string
	setTokenID(string)
}

// XmlTokenRef is the Go form of an IDREF attribute whose Java type is the
// abstract XmlAbstractToken: Go has no upcast, so the reference is held as an
// interface next to the raw ID read from the document.
type XmlTokenRef struct {
	// ID is the referenced xs:ID, always populated.
	ID string
	// Token is the referenced object once Unmarshal has linked the graph.
	Token XmlToken
}

// NewXmlTokenRef builds a reference to an already-built token.
func NewXmlTokenRef(t XmlToken) *XmlTokenRef {
	if t == nil {
		return nil
	}
	return &XmlTokenRef{ID: t.TokenID(), Token: t}
}

// MarshalXMLAttr writes the referenced xs:ID, as JAXB does for an IDREF.
func (r *XmlTokenRef) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if r == nil {
		return xml.Attr{}, nil
	}
	id := r.ID
	if r.Token != nil {
		id = r.Token.TokenID()
	}
	if id == "" {
		return xml.Attr{}, nil
	}
	return xml.Attr{Name: name, Value: id}, nil
}

// UnmarshalXMLAttr records the referenced xs:ID; Unmarshal resolves it to the
// referenced object afterwards.
func (r *XmlTokenRef) UnmarshalXMLAttr(attr xml.Attr) error {
	r.ID = collapseWhitespace(attr.Value)
	return nil
}

func idRefAttr(name xml.Name, id string) (xml.Attr, error) {
	if id == "" {
		return xml.Attr{}, nil
	}
	return xml.Attr{Name: name, Value: id}, nil
}

// MarshalXMLAttr writes the Id of the referenced certificate.
func (x *XmlCertificate) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlCertificate) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced revocation.
func (x *XmlRevocation) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlRevocation) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced signature.
func (x *XmlSignature) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlSignature) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced timestamp.
func (x *XmlTimestamp) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlTimestamp) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced evidence record.
func (x *XmlEvidenceRecord) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlEvidenceRecord) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced signer data.
func (x *XmlSignerData) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlSignerData) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced EAA.
func (x *XmlEAA) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlEAA) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced EAA revocation token.
func (x *XmlEAARevocationToken) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlEAARevocationToken) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced orphan certificate token.
func (x *XmlOrphanCertificateToken) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlOrphanCertificateToken) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced orphan revocation token.
func (x *XmlOrphanRevocationToken) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlOrphanRevocationToken) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced trusted list.
func (x *XmlTrustedList) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlTrustedList) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced list of trusted entities.
func (x *XmlListOfTrustedEntities) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlListOfTrustedEntities) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

// MarshalXMLAttr writes the Id of the referenced trust source list.
func (x *XmlTrustSourceList) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if x == nil {
		return xml.Attr{}, nil
	}
	return idRefAttr(name, x.TokenID())
}

// UnmarshalXMLAttr records the referenced Id; Unmarshal links it afterwards.
func (x *XmlTrustSourceList) UnmarshalXMLAttr(attr xml.Attr) error {
	x.setTokenID(collapseWhitespace(attr.Value))
	return nil
}

var (
	tokenRefType    = reflect.TypeOf((*XmlTokenRef)(nil))
	identifiableTyp = reflect.TypeOf((*identifiable)(nil)).Elem()
)

// Link resolves every IDREF attribute in the tree to the object carrying the
// matching xs:ID, reproducing what the JAXB unmarshaller does natively. It is
// called by Unmarshal; call it by hand after building a tree programmatically.
// References whose target is missing keep the stub carrying the raw ID.
func Link(root any) {
	v := reflect.ValueOf(root)
	ids := map[string]reflect.Value{}
	walk(v, map[uintptr]bool{}, func(sv reflect.Value) {
		if !sv.CanAddr() {
			return
		}
		pv := sv.Addr()
		if !pv.Type().Implements(identifiableTyp) {
			return
		}
		if id := pv.Interface().(identifiable).TokenID(); id != "" {
			if _, dup := ids[id]; !dup {
				ids[id] = pv
			}
		}
	})
	walk(v, map[uintptr]bool{}, func(sv reflect.Value) { resolveRefs(sv, ids) })
}

// walk visits every owned struct in the tree. IDREF attributes are not followed:
// they are the only back-edges, so the owned graph is a tree.
func walk(v reflect.Value, seen map[uintptr]bool, visit func(reflect.Value)) {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		if v.Elem().Kind() == reflect.Struct {
			p := v.Pointer()
			if seen[p] {
				return
			}
			seen[p] = true
		}
		walk(v.Elem(), seen, visit)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			walk(v.Index(i), seen, visit)
		}
	case reflect.Struct:
		visit(v)
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" && !f.Anonymous {
				continue
			}
			if isAttrField(f) {
				continue
			}
			walk(v.Field(i), seen, visit)
		}
	}
}

func isAttrField(f reflect.StructField) bool {
	tag := f.Tag.Get("xml")
	for _, part := range strings.Split(tag, ",")[1:] {
		if part == "attr" {
			return true
		}
	}
	return false
}

func resolveRefs(sv reflect.Value, ids map[string]reflect.Value) {
	t := sv.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !isAttrField(f) {
			continue
		}
		fv := sv.Field(i)
		if fv.Kind() != reflect.Pointer || fv.IsNil() {
			continue
		}
		if f.Type == tokenRefType {
			ref := fv.Interface().(*XmlTokenRef)
			if target, ok := ids[ref.ID]; ok {
				if tok, ok := target.Interface().(XmlToken); ok {
					ref.Token = tok
				}
			}
			continue
		}
		if fv.Elem().Kind() != reflect.Struct || !f.Type.Implements(identifiableTyp) {
			continue
		}
		id := fv.Interface().(identifiable).TokenID()
		if target, ok := ids[id]; ok && target.Type() == f.Type && target.Pointer() != fv.Pointer() {
			fv.Set(target)
		}
	}
}

// ---------------------------------------------------------- certificate extensions

// XmlCertificateExtensionItem is the Go stand-in for the CertificateExtension
// substitution the schema expresses with a choice of named elements: every
// certificate-extension type that may appear inside <CertificateExtensions>
// satisfies it.
type XmlCertificateExtensionItem interface {
	// ExtensionOctets returns the DER of the extension (base type content).
	ExtensionOctets() *Base64Binary
	// ExtensionOID returns the OID attribute of the base type.
	ExtensionOID() *string
	// ExtensionDescription returns the description attribute of the base type.
	ExtensionDescription() *string
	// ExtensionCritical returns the critical attribute of the base type.
	ExtensionCritical() *bool
}

// ExtensionOctets returns the DER of the extension.
func (c *XmlCertificateExtensionContent) ExtensionOctets() *Base64Binary { return c.Octets }

// ExtensionOID returns the OID attribute.
func (a *XmlCertificateExtensionAttrs) ExtensionOID() *string { return a.OID }

// ExtensionDescription returns the description attribute.
func (a *XmlCertificateExtensionAttrs) ExtensionDescription() *string { return a.Description }

// ExtensionCritical returns the critical attribute.
func (a *XmlCertificateExtensionAttrs) ExtensionCritical() *bool { return a.Critical }

// certificateExtensionElements is the @XmlElements table of
// XmlCertificate.certificateExtensions, in schema order.
var certificateExtensionElements = []struct {
	name string
	typ  reflect.Type
}{
	{"KeyUsages", reflect.TypeOf(XmlKeyUsages{})},
	{"ExtendedKeyUsages", reflect.TypeOf(XmlExtendedKeyUsages{})},
	{"CertificatePolicies", reflect.TypeOf(XmlCertificatePolicies{})},
	{"SubjectAlternativeNames", reflect.TypeOf(XmlSubjectAlternativeNames{})},
	{"BasicConstraints", reflect.TypeOf(XmlBasicConstraints{})},
	{"PolicyConstraints", reflect.TypeOf(XmlPolicyConstraints{})},
	{"InhibitAnyPolicy", reflect.TypeOf(XmlInhibitAnyPolicy{})},
	{"NameConstraints", reflect.TypeOf(XmlNameConstraints{})},
	{"CRLDistributionPoints", reflect.TypeOf(XmlCRLDistributionPoints{})},
	{"FreshestCRL", reflect.TypeOf(XmlFreshestCRL{})},
	{"AuthorityKeyIdentifier", reflect.TypeOf(XmlAuthorityKeyIdentifier{})},
	{"SubjectKeyIdentifier", reflect.TypeOf(XmlSubjectKeyIdentifier{})},
	{"AuthorityInformationAccess", reflect.TypeOf(XmlAuthorityInformationAccess{})},
	{"IdPkixOcspNoCheck", reflect.TypeOf(XmlIdPkixOcspNoCheck{})},
	{"ValAssuredShortTermCertificate", reflect.TypeOf(XmlValAssuredShortTermCertificate{})},
	{"NoRevAvail", reflect.TypeOf(XmlNoRevAvail{})},
	{"QcStatements", reflect.TypeOf(XmlQcStatements{})},
	{"OtherExtension", reflect.TypeOf(XmlCertificateExtension{})},
}

func certificateExtensionType(name string) reflect.Type {
	for _, e := range certificateExtensionElements {
		if e.name == name {
			return e.typ
		}
	}
	return nil
}

func certificateExtensionName(v XmlCertificateExtensionItem) string {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	for _, e := range certificateExtensionElements {
		if e.typ == t {
			return e.name
		}
	}
	return ""
}

// CertificateExtensionsWrapper wraps the CertificateExtensions property of
// XmlCertificate. Unlike the other wrappers its children are a choice of named
// element types, so it drives encoding/xml by hand.
type CertificateExtensionsWrapper struct {
	Items []XmlCertificateExtensionItem
}

// All returns the wrapped extensions, tolerating a nil wrapper.
func (w *CertificateExtensionsWrapper) All() []XmlCertificateExtensionItem {
	if w == nil {
		return nil
	}
	return w.Items
}

// MarshalXML writes each extension under the element name its type is bound to.
func (w *CertificateExtensionsWrapper) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for _, item := range w.Items {
		name := certificateExtensionName(item)
		if name == "" {
			return fmt.Errorf("jaxb: %T is not a certificate extension element", item)
		}
		if err := e.EncodeElement(item, xml.StartElement{Name: xml.Name{Local: name}}); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

// UnmarshalXML builds each child from the type its element name is bound to.
func (w *CertificateExtensionsWrapper) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			typ := certificateExtensionType(t.Name.Local)
			if typ == nil {
				if err := d.Skip(); err != nil {
					return err
				}
				continue
			}
			pv := reflect.New(typ)
			if err := d.DecodeElement(pv.Interface(), &t); err != nil {
				return err
			}
			w.Items = append(w.Items, pv.Interface().(XmlCertificateExtensionItem))
		case xml.EndElement:
			return nil
		}
	}
}

// --------------------------------------------------------------- entry points

// Unmarshal parses a diagnostic-data document and links its IDREF graph, the way
// DiagnosticDataFacade.unmarshall does.
func Unmarshal(data []byte) (*XmlDiagnosticData, error) {
	dd := &XmlDiagnosticData{}
	if err := xml.Unmarshal(data, dd); err != nil {
		return nil, err
	}
	Link(dd)
	return dd, nil
}

// Marshal writes a diagnostic-data document byte-for-byte the way
// DiagnosticDataFacade.marshall does: the XML declaration, four-space indented
// output, a trailing newline, and the JAXB spellings jaxbCanonical restores.
func Marshal(dd *XmlDiagnosticData) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "    ")
	if err := enc.Encode(dd); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	body := jaxbCanonical(buf.Bytes())
	out := make([]byte, 0, len(xmlDeclaration)+len(body)+1)
	out = append(out, xmlDeclaration...)
	out = append(out, body...)
	out = append(out, '\n')
	return out, nil
}

// ------------------------------------------------------------- canonicalising

// jaxbCanonical rewrites encoding/xml output into the spelling the JAXB RI
// produces; see the package-level notes at the top of this file.
func jaxbCanonical(in []byte) []byte {
	var out bytes.Buffer
	out.Grow(len(in))
	var stack []string
	for i := 0; i < len(in); {
		if in[i] != '<' {
			j := bytes.IndexByte(in[i:], '<')
			if j < 0 {
				j = len(in) - i
			}
			writeCharData(&out, in[i:i+j])
			i += j
			continue
		}
		end := tagEnd(in, i)
		tag := in[i:end]
		if len(tag) > 1 && tag[1] == '/' {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			out.Write(tag)
			i = end
			continue
		}
		name := tagName(tag)
		closing := append(append([]byte("</"), name...), '>')
		if bytes.HasPrefix(in[end:], closing) && !carriesCharData(name, stack) {
			writeStartTag(&out, tag[:len(tag)-1])
			out.WriteString("/>")
			i = end + len(closing)
			continue
		}
		writeStartTag(&out, tag)
		stack = append(stack, name)
		i = end
	}
	return out.Bytes()
}

// tagEnd returns the index just past the '>' closing the tag starting at i.
// encoding/xml escapes '<' and '>' inside attribute values, so no quoting-aware
// scan is needed.
func tagEnd(in []byte, i int) int {
	j := bytes.IndexByte(in[i:], '>')
	if j < 0 {
		return len(in)
	}
	return i + j + 1
}

func tagName(tag []byte) string {
	s := tag[1:]
	for k := 0; k < len(s); k++ {
		switch s[k] {
		case ' ', '\t', '\n', '\r', '>', '/':
			return string(s[:k])
		}
	}
	return string(s)
}

// writeStartTag copies a start tag, rewriting the escapes inside attribute
// values to the RI's spelling.
func writeStartTag(out *bytes.Buffer, tag []byte) {
	inValue := false
	start := 0
	for k := 0; k < len(tag); k++ {
		if tag[k] != '"' {
			continue
		}
		if inValue {
			writeAttrValue(out, tag[start:k])
		} else {
			out.Write(tag[start:k])
		}
		out.WriteByte('"')
		inValue = !inValue
		start = k + 1
	}
	if start < len(tag) {
		out.Write(tag[start:])
	}
}

// charDataEscapes maps the numeric references encoding/xml emits in character
// data to the spelling the RI uses there.
var charDataEscapes = []struct{ from, to string }{
	{"&#34;", `"`},
	{"&#39;", `'`},
	{"&#x9;", "\t"},
	{"&#xA;", "\n"},
	{"&#xD;", "&#13;"},
}

// attrValueEscapes maps the same references to the spelling the RI uses inside
// attribute values, where a double quote must stay escaped.
var attrValueEscapes = []struct{ from, to string }{
	{"&#34;", "&quot;"},
	{"&#39;", `'`},
	{"&#x9;", "\t"},
	{"&#xA;", "&#10;"},
	{"&#xD;", "&#13;"},
}

func writeCharData(out *bytes.Buffer, b []byte) {
	writeRewritten(out, b, charDataEscapes)
}

func writeAttrValue(out *bytes.Buffer, b []byte) {
	writeRewritten(out, b, attrValueEscapes)
}

func writeRewritten(out *bytes.Buffer, b []byte, table []struct{ from, to string }) {
	if bytes.IndexByte(b, '&') < 0 {
		out.Write(b)
		return
	}
	s := string(b)
	for _, e := range table {
		if e.from != e.to {
			s = strings.ReplaceAll(s, e.from, e.to)
		}
	}
	out.WriteString(s)
}
