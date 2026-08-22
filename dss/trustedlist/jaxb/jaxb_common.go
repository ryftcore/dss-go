// Ported from ts_119612v020401_xsd.xsd (DSS 6.5.RC1) via the JAXB classes
// generated into eu.europa.esig.trustedlist.jaxb.tsl: the schema's simple
// types (Base64Binary, BigInteger) and its xs:anyType/xs:any extensibility
// machinery (AnyType.java, ExtensionType.java).
//
// # Cross-schema types outside this manifest
//
// A handful of properties across the four packages this module ports are
// typed by a JAXB class generated from a schema OUTSIDE specs-trusted-list:
//
//   - TrustStatusListType.Signature and DigitalIdentityType.KeyValue are
//     ds:Signature/ds:KeyValue - eu.europa.esig.xmldsig.jaxb.{SignatureType,
//     KeyValueType} (specs-xmldsig).
//   - ecc's CriteriaListType.otherCriteriaList and tslx's
//     TakenOverByType... no - TakenOverByType.otherQualifier stays this
//     package's own tsl.AnyType; it is instead ecc.CriteriaListType's OWN
//     otherCriteriaList, and every PolicyIdentifier/KeyPurposeId/
//     AttributeOID/QcType/QcStatementId property, that are
//     eu.europa.esig.xades.jaxb.xades132.{AnyType,ObjectIdentifierType}
//     (specs-xades).
//
// specs-xmldsig and specs-xades have not been ported (this package covers
// specs-trusted-list only), and TLValidatorTask's actual
// signature verification runs the frozen xades/xmldsig validator over the
// ORIGINAL document bytes, never through this JAXB tree (mirroring upstream:
// Facade's unmarshalled ds:Signature is not what TLValidatorTask
// verifies either) - so nothing downstream needs these five properties
// interpreted, only round-tripped. Each is therefore modelled as a
// raw-capture stand-in (own attributes plus a captured, replayed token
// stream - see foreignContent below for why tokens rather than innerxml
// bytes): ObjectIdentifierType and xadesAnyType for the plain,
// field-tag-named positions, dsigSignature/dsigKeyValue for the two
// ds:-namespaced ones, which - like dss/validationreport/jaxb's
// SignatureType - hard-code their own "ns2:" element name rather than leave
// it to encoding/xml's field-tag namespace handling; see jaxb_tsl_root.go's
// header for why "ns2:" is the right, empirically-pinned prefix for both
// Facade and MRAFacade documents.
//
// Known (accepted) over-preservation: because dsigSignature captures every
// token of ds:Signature's subtree verbatim rather than parsing it through
// the REAL, strictly-typed xades132/xades141 JAXB classes (out of scope),
// it round-trips content those classes would silently drop on unmarshal -
// e.g. a real fixture (not in this package's KAT corpus, to keep it green)
// places xades:SigningCertificateV2/IssuerSerialV2 (properly XAdES 1.4.1
// elements) inside the 1.3.2 xades: namespace; xades132.
// SignedSignaturePropertiesType has no property for that name in that
// namespace, so the JAXB RI's oracle remarshal omits it entirely, while
// the token replay here preserves it. This can only ever make the
// Signature/KeyValue/ObjectIdentifierType/xadesAnyType content here a
// SUPERSET of the JAXB RI's own (never missing information the RI kept),
// and does not affect anything this package is responsible for (schema
// completeness of tsl/ecc/tslx/mra, or TLValidatorTask's signature
// verification, which uses the original bytes - see above).
package jaxb

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------- simple types

// Base64Binary is the Go form of a JAXB byte[] property bound to
// xs:base64Binary (DigitalIdentityType.X509Certificate/X509SKI). A nil
// Base64Binary and an empty one are distinct: JAXB omits a null property
// entirely but writes an empty element for a zero-length array - callers
// rely on the ",omitempty" tag for the former.
type Base64Binary []byte

// MarshalText encodes the bytes with the standard base64 alphabet,
// unwrapped, as the RI's base64Binary printer does.
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

// BigInteger is the Go form of a java.math.BigInteger-typed property
// (xs:integer/xs:positiveInteger/xs:nonNegativeInteger:
// TSLVersionIdentifier, TSLSequenceNumber, HistoricalInformationPeriod,
// MutualRecognitionAgreementInformationType's technicalType/version/
// MRADepth attributes).
//
// It exists because the obvious Go binding - a bare *big.Int, which already
// implements encoding.TextUnmarshaler - is WRONG for this schema:
// math/big.Int.UnmarshalText parses with base 0, so it reads a leading "0"
// as an OCTAL prefix and rejects a decimal integer written with leading
// zeros, where Java's BigInteger(String) constructor (what JAXB's built-in
// xs:integer binding calls) is decimal-only. Same rationale and shape as
// dss/diagnostic/jaxb.BigInteger.
type BigInteger struct {
	big.Int
}

// NewBigInteger wraps a *big.Int, answering nil for a nil input.
func NewBigInteger(value *big.Int) *BigInteger {
	if value == nil {
		return nil
	}
	return &BigInteger{Int: *value}
}

// NewBigIntegerFromInt64 wraps an int64.
func NewBigIntegerFromInt64(value int64) *BigInteger {
	return &BigInteger{Int: *big.NewInt(value)}
}

// BigInt answers the wrapped value, nil for a nil receiver.
func (b *BigInteger) BigInt() *big.Int {
	if b == nil {
		return nil
	}
	return &b.Int
}

// UnmarshalText parses the DECIMAL lexical form, shadowing the promoted
// big.Int.UnmarshalText and its base-0 prefix sniffing.
func (b *BigInteger) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))
	if _, ok := b.Int.SetString(s, 10); !ok {
		return fmt.Errorf("invalid xs:integer: %q", s)
	}
	return nil
}

// ------------------------------------------------------- foreign-type stand-ins

// foreignContent is the shared shape of every raw-capture stand-in: the
// element's own attributes plus its child content, captured as a resolved
// TOKEN STREAM (not raw bytes) and replayed through the same *xml.Encoder
// that is marshalling everything else - see this file's header.
//
// A byte-level (innerxml) capture was tried first and rejected: this
// port's own KAT corpus (xml_kat_test.go) has real content that legitimately
// crosses a raw-capture boundary using a namespace prefix declared on an
// ANCESTOR outside the captured span - mra-lotl.xml's
// CriteriaListType.otherCriteriaList (xadesAnyType, captured because its
// declared type is outside this manifest - see this file's header) embeds
// mra:QcStatementSet, whose "mra" prefix the document binds at the ROOT,
// well above <otherCriteriaList>. A byte-for-byte replay of the captured
// span carries the literal "mra:" text forward with no accompanying
// xmlns:mra declaration in scope, producing unresolvable output. Capturing
// a TOKEN stream instead uses xml.Decoder's own namespace resolution (which
// already walks the whole ancestor chain to resolve every prefix to its
// full URI), and replaying those already-resolved tokens through
// e.EncodeToken lets encoding/xml emit whatever xmlns declaration each one
// needs, which is always well-formed regardless of where the source
// document happened to declare it - at the cost of not always reproducing
// the JAXB RI's own prefix-reuse scheme or attribute order verbatim, which
// is why xml_kat_test.go's marshal-parity KAT falls back to a canonical,
// insignificant-whitespace/attribute-order-insensitive comparison once a
// byte-for-byte one fails - see jaxb_tsl_root.go's "Known deviation" header
// note.
type foreignContent struct {
	Attrs  []xml.Attr
	Tokens []xml.Token
}

// unmarshal captures the element's own attributes and its child content as
// a token stream, stopping at (and consuming) the matching end element.
func (c *foreignContent) unmarshal(d *xml.Decoder, start xml.StartElement) error {
	c.Attrs = append([]xml.Attr(nil), start.Attr...)
	depth := 0
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.EndElement:
			if depth == 0 {
				return nil
			}
			depth--
			c.Tokens = append(c.Tokens, t)
		case xml.StartElement:
			depth++
			cp := t.Copy()
			cp.Attr = stripXmlnsAttrs(cp.Attr)
			c.Tokens = append(c.Tokens, cp)
		case xml.CharData:
			// Whitespace-only runs are dropped for the same reason
			// unmarshalAnyContent drops them (see that function's header):
			// captured verbatim, they compound with encoding/xml's own
			// Indent() on every further Unmarshal/Marshal cycle.
			if strings.TrimSpace(string(t)) != "" {
				c.Tokens = append(c.Tokens, t.Copy())
			}
		default:
			c.Tokens = append(c.Tokens, xml.CopyToken(tok))
		}
	}
}

// stripXmlnsAttrs drops namespace-declaration attributes from a captured
// StartElement's attribute list. Every element/attribute name a captured
// token stream carries has already been resolved to its full namespace URI
// by xml.Decoder, so the declaration that made a prefix resolvable is
// redundant information the encoder does not need replayed - and, replayed
// unfiltered, becomes actively wrong once this package's OWN Marshal output
// (which - unlike a document downloaded from the wild - adds a fresh local
// "xmlns" to every element a captured token stream re-descends into, since
// encoding/xml's manual EncodeToken path does not track which ancestor
// already declared what) is fed back through Unmarshal: the second pass
// then captures that already-local "xmlns" attribute value ALONGSIDE the
// resolved Name.Space it produced it from, and replaying both yields a
// duplicate "xmlns=... xmlns=..." on remarshal. Found by
// TestUnmarshalMarshalIdempotent, this package's own double-round-trip
// check (xml_test.go).
func stripXmlnsAttrs(attrs []xml.Attr) []xml.Attr {
	if len(attrs) == 0 {
		return attrs
	}
	out := attrs[:0:0]
	for _, a := range attrs {
		if a.Name.Space == "xmlns" || (a.Name.Space == "" && a.Name.Local == "xmlns") {
			continue
		}
		out = append(out, a)
	}
	return out
}

// marshalNamed writes the element under the given literal name (a bare
// local name reuses the field tag's own name/namespace; a "prefix:local"
// spelling - Space left empty - overrides it, matching
// dss/validationreport/jaxb/jaxb_crossns.go's SignatureType), its captured
// attributes (namespace declarations dropped - encoding/xml regenerates
// whatever the replayed content needs) and its captured child tokens.
func (c *foreignContent) marshalNamed(e *xml.Encoder, start xml.StartElement, name string) error {
	if name != "" {
		start.Name = xml.Name{Local: name}
	}
	start.Attr = stripXmlnsAttrs(c.Attrs)
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for _, tok := range c.Tokens {
		if err := e.EncodeToken(tok); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

// ObjectIdentifierType stands in for eu.europa.esig.xades.jaxb.xades132.ObjectIdentifierType
// (used, always under its OWN field's element name/namespace - never
// ds:-style prefix-overridden - by ecc.PoliciesListType.PolicyIdentifier,
// tslx.ExtendedKeyUsageType.KeyPurposeId,
// tslx.CertSubjectDNAttributeType.AttributeOID, mra.QcStatementInfoType.QcType
// and mra.QcStatementType.QcStatementId): see this file's header.
type ObjectIdentifierType struct {
	foreignContent
}

// UnmarshalXML captures the element's attributes and child content as a
// token stream - see foreignContent's header.
func (o *ObjectIdentifierType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return o.unmarshal(d, start)
}

// MarshalXML writes the element's captured attributes and child tokens back,
// under whatever name/namespace the caller's field tag supplied.
func (o *ObjectIdentifierType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return o.marshalNamed(e, start, "")
}

// xadesAnyType stands in for eu.europa.esig.xades.jaxb.xades132.AnyType
// (ecc.CriteriaListType.otherCriteriaList only - a field with no @XmlElement
// of its own, so JAXB binds it to the literal field name "otherCriteriaList"
// in ecc's OWN namespace; see jaxb_ecc.go). Named unexported: it is bound at
// exactly one position, unlike ObjectIdentifierType.
type xadesAnyType struct {
	foreignContent
}

func (a *xadesAnyType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return a.unmarshal(d, start)
}

func (a *xadesAnyType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return a.marshalNamed(e, start, "")
}

// dsigSignature stands in for eu.europa.esig.xmldsig.jaxb.SignatureType
// (TrustStatusListType.Signature, ds:Signature): its own "ns2:Signature"
// element name is hard-coded, matching
// dss/validationreport/jaxb/jaxb_crossns.go's SignatureType and this file's
// header.
type dsigSignature struct {
	foreignContent
}

func (s *dsigSignature) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return s.unmarshal(d, start)
}

func (s *dsigSignature) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return s.marshalNamed(e, start, "ns2:Signature")
}

// dsigKeyValue stands in for eu.europa.esig.xmldsig.jaxb.KeyValueType
// (DigitalIdentityType.KeyValue, ds:KeyValue): see dsigSignature above.
type dsigKeyValue struct {
	foreignContent
}

func (k *dsigKeyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return k.unmarshal(d, start)
}

func (k *dsigKeyValue) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return k.marshalNamed(e, start, "ns2:KeyValue")
}

// ---------------------------------------------------------------- xs:any content

// wildcardElement is one entry of the xs:any "lax" dispatch table: every
// element globally declared (@XmlElementDecl) by one of the four
// ObjectFactory classes this manifest ports (tsl, ecc, tslx, mra) - the
// full set JAXB's @XmlAnyElement(lax=true) can resolve to a typed value
// rather than a raw org.w3c.dom.Element, reproduced here because real
// trusted lists exercise it far more broadly than just the extension
// mechanism's obvious cases (ecc:Qualifications, tslx:TakenOverBy, ...):
// OtherTSLPointerType.AdditionalInformation's "OtherInformation" choice
// member routinely wraps TSLType, SchemeTerritory, SchemeOperatorName,
// SchemeTypeCommunityRules and tslx:MimeType, all by themselves, one per
// <OtherInformation> - found by this package's own KAT corpus
// (testdata/tl/dk_tl-sn21.xml, eu-lotl-250.xml). The root element
// (TrustServiceStatusList) is the one entry omitted: it can never
// legitimately nest inside a wildcard position.
//
// prefix is the literal marshal prefix from the root-declared table (see
// jaxb_tsl_root.go); "" reuses the tsl-namespace default (no prefix needed,
// the document's own default xmlns already covers it).
var wildcardElements = []struct {
	namespace, local, prefix string
	typ                      reflect.Type
}{
	// tsl
	{NamespaceTSL, "PostalAddresses", "", reflect.TypeOf(PostalAddressListType{})},
	{NamespaceTSL, "PostalAddress", "", reflect.TypeOf(PostalAddressType{})},
	{NamespaceTSL, "ElectronicAddress", "", reflect.TypeOf(ElectronicAddressType{})},
	{NamespaceTSL, "Extension", "", reflect.TypeOf(ExtensionType{})},
	{NamespaceTSL, "TrustServiceProviderList", "", reflect.TypeOf(TrustServiceProviderListType{})},
	{NamespaceTSL, "SchemeInformation", "", reflect.TypeOf(TSLSchemeInformationType{})},
	{NamespaceTSL, "TSLType", "", reflect.TypeOf("")},
	{NamespaceTSL, "SchemeOperatorName", "", reflect.TypeOf(InternationalNamesType{})},
	{NamespaceTSL, "SchemeName", "", reflect.TypeOf(InternationalNamesType{})},
	{NamespaceTSL, "SchemeInformationURI", "", reflect.TypeOf(NonEmptyMultiLangURIListType{})},
	{NamespaceTSL, "SchemeTypeCommunityRules", "", reflect.TypeOf(NonEmptyMultiLangURIListType{})},
	{NamespaceTSL, "SchemeTerritory", "", reflect.TypeOf("")},
	{NamespaceTSL, "PolicyOrLegalNotice", "", reflect.TypeOf(PolicyOrLegalnoticeType{})},
	{NamespaceTSL, "NextUpdate", "", reflect.TypeOf(NextUpdateType{})},
	{NamespaceTSL, "PointersToOtherTSL", "", reflect.TypeOf(OtherTSLPointersType{})},
	{NamespaceTSL, "OtherTSLPointer", "", reflect.TypeOf(OtherTSLPointerType{})},
	{NamespaceTSL, "ServiceDigitalIdentities", "", reflect.TypeOf(ServiceDigitalIdentityListType{})},
	{NamespaceTSL, "AdditionalInformation", "", reflect.TypeOf(AdditionalInformationType{})},
	{NamespaceTSL, "DistributionPoints", "", reflect.TypeOf(NonEmptyURIListType{})},
	{NamespaceTSL, "TrustServiceProvider", "", reflect.TypeOf(TSPType{})},
	{NamespaceTSL, "TSPInformation", "", reflect.TypeOf(TSPInformationType{})},
	{NamespaceTSL, "TSPServices", "", reflect.TypeOf(TSPServicesListType{})},
	{NamespaceTSL, "TSPService", "", reflect.TypeOf(TSPServiceType{})},
	{NamespaceTSL, "ServiceInformation", "", reflect.TypeOf(TSPServiceInformationType{})},
	{NamespaceTSL, "ServiceStatus", "", reflect.TypeOf("")},
	{NamespaceTSL, "ServiceSupplyPoints", "", reflect.TypeOf(ServiceSupplyPointsType{})},
	{NamespaceTSL, "ServiceTypeIdentifier", "", reflect.TypeOf("")},
	{NamespaceTSL, "ServiceDigitalIdentity", "", reflect.TypeOf(DigitalIdentityListType{})},
	{NamespaceTSL, "ServiceHistory", "", reflect.TypeOf(ServiceHistoryType{})},
	{NamespaceTSL, "ServiceHistoryInstance", "", reflect.TypeOf(ServiceHistoryInstanceType{})},
	{NamespaceTSL, "ExpiredCertsRevocationInfo", "", reflect.TypeOf("")},
	{NamespaceTSL, "AdditionalServiceInformation", "", reflect.TypeOf(AdditionalServiceInformationType{})},

	// ecc
	{NamespaceECC, "Qualifications", "ns5", reflect.TypeOf(QualificationsType{})},

	// tslx
	{NamespaceTSLX, "MimeType", "ns4", reflect.TypeOf("")},
	{NamespaceTSLX, "X509CertificateLocation", "ns4", reflect.TypeOf("")},
	{NamespaceTSLX, "PublicKeyLocation", "ns4", reflect.TypeOf("")},
	{NamespaceTSLX, "ExtendedKeyUsage", "ns4", reflect.TypeOf(ExtendedKeyUsageType{})},
	{NamespaceTSLX, "TakenOverBy", "ns4", reflect.TypeOf(TakenOverByType{})},
	{NamespaceTSLX, "CertSubjectDNAttribute", "ns4", reflect.TypeOf(CertSubjectDNAttributeType{})},

	// mra
	{NamespaceMRA, "MutualRecognitionAgreementInformation", "ns6", reflect.TypeOf(MutualRecognitionAgreementInformationType{})},
	{NamespaceMRA, "QcStatementSet", "ns6", reflect.TypeOf(QcStatementListType{})},
	{NamespaceMRA, "QcType", "ns6", reflect.TypeOf(ObjectIdentifierType{})},
	{NamespaceMRA, "QcCClegislation", "ns6", reflect.TypeOf("")},
}

func wildcardByName(space, local string) (prefix string, typ reflect.Type, ok bool) {
	for _, w := range wildcardElements {
		if w.namespace == space && w.local == local {
			return w.prefix, w.typ, true
		}
	}
	return "", nil, false
}

// wildcardPrefix answers the fixed marshal prefix for a wildcardElements
// namespace (unambiguous: one prefix per namespace, regardless of how many
// element names in it are registered).
func wildcardPrefix(namespace string) (prefix string, ok bool) {
	for _, w := range wildcardElements {
		if w.namespace == namespace {
			return w.prefix, true
		}
	}
	return "", false
}

// RawWildcardElement captures one child of an xs:any wildcard position that
// wildcardElements does not recognize (foreign-namespace vendor extensions,
// or a genuinely unknown element): its resolved name, own attributes and
// inner XML, verbatim. Marshal re-declares the element's namespace locally
// as a default xmlns rather than replaying an original prefix - encoding/xml
// resolves prefixes away on Unmarshal, so the source document's own prefix
// spelling is not recoverable - which is well-formed and semantically
// identical but not necessarily byte-identical to an exotic input prefix; no
// fixture in the KAT corpus exercises this fallback; every extension
// element eu-lotl/country trusted lists actually carry is one of
// wildcardElements.
type RawWildcardElement struct {
	Name xml.Name
	foreignContent
}

// ---------------------------------------------------------------- AnyContent

// AnyItem is one item of an xs:any/xs:anyType mixed content list
// (AnyType.content / ExtensionType, @XmlMixed @XmlAnyElement(lax=true)):
// either character data (IsText, including whitespace-only runs between
// sibling elements) or one child element, recognized (Elem, one of
// wildcardElements' pointer types) or not (Raw).
type AnyItem struct {
	IsText bool
	Text   string
	// Elem is the decoded value for a recognized element (a pointer to one
	// of wildcardElements' types); ElemName is its wildcardElements entry's
	// namespace/local name, captured at unmarshal time rather than
	// re-derived from Elem's Go type on marshal, since several entries -
	// every bare xs:string leaf (TSLType, SchemeTerritory, ServiceStatus,
	// ServiceTypeIdentifier, ExpiredCertsRevocationInfo, tslx's MimeType/
	// X509CertificateLocation/PublicKeyLocation, mra's QcCClegislation) -
	// share the identical Go type (string) under different element names,
	// which a by-type reverse lookup could not disambiguate.
	Elem     any
	ElemName xml.Name
	Raw      *RawWildcardElement
}

// AnyContent is the shared implementation of AnyType.content: the Go form of
// every JAXB class whose element content is List<Object> bound
// @XmlMixed @XmlAnyElement(lax=true). AnyType and ExtensionType each embed
// it but supply their OWN MarshalXML/UnmarshalXML (rather than relying on
// AnyContent's own being promoted, which would only ever see AnyContent's
// fields, never ExtensionType's Critical attribute) and call
// unmarshalAnyContent/marshalAnyItems directly.
type AnyContent struct {
	Items []AnyItem
}

// unmarshalAnyContent reads the mixed content of the element start already
// opened, reproducing lax @XmlAnyElement dispatch via wildcardElements.
// Whitespace-only character data (the indentation between sibling elements
// - the only kind any fixture in the KAT corpus carries here) is
// dropped rather than preserved as a mixed-content item: keeping it is
// pointless (JAXB's own oracle, being itself indented, would never
// distinguish "no text was here" from "insignificant whitespace was here"
// in a way this package's marshal-parity KAT could observe) and actively
// harmful for a second Unmarshal/Marshal cycle - each preserved whitespace
// item receives ITS OWN indentation from encoding/xml's Indent() on top of
// the text already carried, doubling on every pass; found by
// TestUnmarshalMarshalIdempotent (xml_test.go).
func unmarshalAnyContent(d *xml.Decoder, start xml.StartElement) (AnyContent, error) {
	var c AnyContent
	for {
		tok, err := d.Token()
		if err != nil {
			return c, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			t = t.Copy()
			if _, typ, ok := wildcardByName(t.Name.Space, t.Name.Local); ok {
				pv := reflect.New(typ)
				if err := d.DecodeElement(pv.Interface(), &t); err != nil {
					return c, err
				}
				c.Items = append(c.Items, AnyItem{Elem: pv.Interface(), ElemName: xml.Name{Space: t.Name.Space, Local: t.Name.Local}})
				continue
			}
			raw := &RawWildcardElement{Name: t.Name}
			if err := raw.foreignContent.unmarshal(d, t); err != nil {
				return c, err
			}
			c.Items = append(c.Items, AnyItem{Raw: raw})
		case xml.CharData:
			if text := string(t.Copy()); strings.TrimSpace(text) != "" {
				c.Items = append(c.Items, AnyItem{IsText: true, Text: text})
			}
		case xml.EndElement:
			return c, nil
		}
	}
}

// marshalAnyItems writes the mixed content items already opened by the
// caller (AnyType/ExtensionType's own MarshalXML), reproducing the JAXB RI's
// fixed root-declared prefix for each recognized element (see
// wildcardElements) and a locally-declared default xmlns for a raw one (see
// RawWildcardElement).
func marshalAnyItems(e *xml.Encoder, items []AnyItem) error {
	for _, it := range items {
		switch {
		case it.IsText:
			if err := e.EncodeToken(xml.CharData(it.Text)); err != nil {
				return err
			}
		case it.Elem != nil:
			prefix, ok := wildcardPrefix(it.ElemName.Space)
			if !ok {
				return fmt.Errorf("jaxb: %s is not a wildcard-registered namespace", it.ElemName.Space)
			}
			name := it.ElemName.Local
			if prefix != "" {
				name = prefix + ":" + name
			}
			if err := e.EncodeElement(it.Elem, xml.StartElement{Name: xml.Name{Local: name}}); err != nil {
				return err
			}
		case it.Raw != nil:
			start := xml.StartElement{Name: it.Raw.Name}
			if err := it.Raw.foreignContent.marshalNamed(e, start, ""); err != nil {
				return err
			}
		}
	}
	return nil
}

// AnyType is the Go form of the generated JAXB class AnyType (complexType
// AnyType, mixed="true" with a lax xs:any): used by
// AdditionalServiceInformationType.OtherInformation,
// AdditionalInformationType's OtherInformation choice member and
// DigitalIdentityType.Other.
type AnyType struct {
	AnyContent
}

// UnmarshalXML reads the element's mixed content - see unmarshalAnyContent.
func (a *AnyType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	c, err := unmarshalAnyContent(d, start)
	a.AnyContent = c
	return err
}

// MarshalXML writes the element's mixed content - see marshalAnyItems.
func (a *AnyType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := marshalAnyItems(e, a.Items); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}

// ExtensionType is the Go form of the generated JAXB class ExtensionType
// (complexType ExtensionType, extends AnyType): the "Extension" element of
// ExtensionsListType, ETSI TS 119 612's mechanism for SchemeExtensions,
// TSPInformationExtensions and ServiceInformationExtensions alike.
type ExtensionType struct {
	AnyContent
	Critical bool
}

// UnmarshalXML reads the Critical attribute then the element's mixed
// content.
func (x *ExtensionType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, a := range start.Attr {
		if a.Name.Local == "Critical" {
			x.Critical = a.Value == "true" || a.Value == "1"
		}
	}
	c, err := unmarshalAnyContent(d, start)
	x.AnyContent = c
	return err
}

// MarshalXML writes the Critical attribute then the element's mixed
// content.
func (x *ExtensionType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{{Name: xml.Name{Local: "Critical"}, Value: strconv.FormatBool(x.Critical)}}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := marshalAnyItems(e, x.Items); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}
