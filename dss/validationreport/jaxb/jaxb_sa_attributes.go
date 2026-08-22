// Ported from the generated JAXB classes of Annex A ("Signature Attributes")
// of ETSI TS 119 102-2: AttributeBaseType.java, SACertIDListType.java,
// SACertIDType.java, SACommitmentTypeIndicationType.java, SAContactInfoType.java,
// SACounterSignatureType.java, SACRLIDType.java, SADataObjectFormatType.java,
// SADSSType.java, SAFilterType.java, SAMessageDigestType.java, SANameType.java,
// SAOCSPIDType.java, SAOneSignerRoleType.java, SAReasonType.java,
// SARevIDListType.java, SASigningTimeType.java, SASignatureProductionPlaceType.java,
// SASignerRoleType.java, SASigPolicyIdentifierType.java, SASubFilterType.java,
// SATimestampType.java, SAVRIType.java and SignatureAttributesType.java
// (specs-validation-report, DSS 6.5.RC1).
//
// # SignatureAttributesType and SARevIDListType.crlidOrOCSPID
//
// SignatureAttributesType's only content is a 35-member xs:choice
// maxOccurs="unbounded" (every kind of signature attribute a validation
// report can carry): xjc bound it to List<JAXBElement<?>>, choosing the
// concrete type by element name at unmarshal time. SARevIDListType's
// crlidOrOCSPID property is the same idea at a smaller scale - a two-member
// choice bound through @XmlElements rather than JAXBElement. Both are
// modelled the same way here: a name-tagged item list plus hand-written
// MarshalXML/UnmarshalXML that dispatches by element name against a small
// name->type table, exactly the technique dss/diagnostic/jaxb's
// CertificateExtensionsWrapper uses for DiagnosticData.xsd's
// CertificateExtensions choice. Order is preserved item-for-item, which is
// what the marshal-parity KAT depends on.
//
// Several element names bind to the *same* Go type at different positions
// in the choice (SATimestampType backs seven of them, AttributeBaseType
// backs five, SACertIDListType three, SARevIDListType and SAListOfIntegers
// two each) - the reverse cannot be inferred from a value's runtime type
// alone, so each item records its element name explicitly rather than
// deriving it, unlike CertificateExtensionsWrapper's 1:1 table.
package jaxb

import (
	"encoding/xml"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

// AttributeBaseType is the Go form of the generated JAXB class
// AttributeBaseType (complexType AttributeBaseType), the common base of
// every SA* signature-attribute type.
type AttributeBaseType struct {
	AttributeObject []*VOReferenceType `xml:"AttributeObject,omitempty"`
	Signed          *bool              `xml:"Signed,attr,omitempty"`
}

// SASigningTimeType is the Go form of the generated JAXB class SASigningTimeType.
type SASigningTimeType struct {
	AttributeBaseType
	Time XSDateTime `xml:"Time"`
}

// SACertIDType is the Go form of the generated JAXB class SACertIDType.
type SACertIDType struct {
	X509IssuerSerial Base64Binary     `xml:"X509IssuerSerial,omitempty"`
	DigestMethod     DigestMethodType `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod"`
	DigestValue      DSDigestValue    `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue"`
}

// SACertIDListType is the Go form of the generated JAXB class SACertIDListType.
type SACertIDListType struct {
	AttributeBaseType
	CertID []*SACertIDType `xml:"CertID,omitempty"`
}

// SADataObjectFormatType is the Go form of the generated JAXB class SADataObjectFormatType.
type SADataObjectFormatType struct {
	AttributeBaseType
	ContentType *string `xml:"ContentType,omitempty"`
	MimeType    *string `xml:"MimeType,omitempty"`
}

// SACommitmentTypeIndicationType is the Go form of the generated JAXB class
// SACommitmentTypeIndicationType.
type SACommitmentTypeIndicationType struct {
	AttributeBaseType
	CommitmentTypeIdentifier string `xml:"CommitmentTypeIdentifier"`
}

// SASigPolicyIdentifierType is the Go form of the generated JAXB class SASigPolicyIdentifierType.
type SASigPolicyIdentifierType struct {
	AttributeBaseType
	SigPolicyId string `xml:"SigPolicyId"`
}

// SASignatureProductionPlaceType is the Go form of the generated JAXB class
// SASignatureProductionPlaceType.
type SASignatureProductionPlaceType struct {
	AttributeBaseType
	AddressString []string `xml:"AddressString"`
}

// SAOneSignerRoleType is the Go form of the generated JAXB class SAOneSignerRoleType.
type SAOneSignerRoleType struct {
	Role            string               `xml:"Role"`
	EndorsementType ValueEndorsementType `xml:"EndorsementType"`
}

// SASignerRoleType is the Go form of the generated JAXB class SASignerRoleType.
type SASignerRoleType struct {
	AttributeBaseType
	RoleDetails []*SAOneSignerRoleType `xml:"RoleDetails"`
}

// SACounterSignatureType is the Go form of the generated JAXB class SACounterSignatureType.
type SACounterSignatureType struct {
	AttributeBaseType
	CounterSignature SignatureReferenceType `xml:"CounterSignature"`
}

// SACRLIDType is the Go form of the generated JAXB class SACRLIDType.
type SACRLIDType struct {
	DigestMethod DigestMethodType `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod"`
	DigestValue  DSDigestValue    `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue"`
}

// SAOCSPIDType is the Go form of the generated JAXB class SAOCSPIDType.
type SAOCSPIDType struct {
	ProducedAt        XSDateTime   `xml:"ProducedAt"`
	ResponderIDByName *string      `xml:"ResponderIDByName,omitempty"`
	ResponderIDByKey  Base64Binary `xml:"ResponderIDByKey,omitempty"`
}

// SARevIDItem is one member of SARevIDListType's CRLID/OCSPID choice: exactly
// one of CRLID/OCSPID is set.
type SARevIDItem struct {
	CRLID  *SACRLIDType
	OCSPID *SAOCSPIDType
}

// SARevIDListType is the Go form of the generated JAXB class SARevIDListType.
// Unlike the other AttributeBaseType extensions its choice property has no
// dedicated wrapping element of its own (its members sit directly among
// AttributeBaseType's AttributeObject elements), so this type implements
// MarshalXML/UnmarshalXML by hand rather than relying on struct-tag-driven
// encoding for its whole content, not only for the choice list.
type SARevIDListType struct {
	AttributeBaseType
	CRLIDOrOCSPID []SARevIDItem
}

// MarshalXML writes Signed, AttributeObject and the CRLID/OCSPID choice
// items, in that schema order.
func (t *SARevIDListType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if t.Signed != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Signed"}, Value: strconv.FormatBool(*t.Signed)})
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for _, ao := range t.AttributeObject {
		if err := e.EncodeElement(ao, xml.StartElement{Name: xml.Name{Local: "AttributeObject"}}); err != nil {
			return err
		}
	}
	for _, item := range t.CRLIDOrOCSPID {
		switch {
		case item.CRLID != nil:
			if err := e.EncodeElement(item.CRLID, xml.StartElement{Name: xml.Name{Local: "CRLID"}}); err != nil {
				return err
			}
		case item.OCSPID != nil:
			if err := e.EncodeElement(item.OCSPID, xml.StartElement{Name: xml.Name{Local: "OCSPID"}}); err != nil {
				return err
			}
		}
	}
	return e.EncodeToken(start.End())
}

// UnmarshalXML reads Signed, AttributeObject and the CRLID/OCSPID choice items.
func (t *SARevIDListType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, a := range start.Attr {
		if a.Name.Local == "Signed" {
			v, err := strconv.ParseBool(a.Value)
			if err != nil {
				return fmt.Errorf("invalid Signed attribute: %w", err)
			}
			t.Signed = &v
		}
	}
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := tok.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "AttributeObject":
				var v VOReferenceType
				if err := d.DecodeElement(&v, &tt); err != nil {
					return err
				}
				t.AttributeObject = append(t.AttributeObject, &v)
			case "CRLID":
				var v SACRLIDType
				if err := d.DecodeElement(&v, &tt); err != nil {
					return err
				}
				t.CRLIDOrOCSPID = append(t.CRLIDOrOCSPID, SARevIDItem{CRLID: &v})
			case "OCSPID":
				var v SAOCSPIDType
				if err := d.DecodeElement(&v, &tt); err != nil {
					return err
				}
				t.CRLIDOrOCSPID = append(t.CRLIDOrOCSPID, SARevIDItem{OCSPID: &v})
			default:
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			return nil
		}
	}
}

// SADSSType is the Go form of the generated JAXB class SADSSType.
type SADSSType struct {
	AttributeBaseType
	Certs *VOReferenceType `xml:"Certs,omitempty"`
	CRLs  *VOReferenceType `xml:"CRLs,omitempty"`
	OCSPs *VOReferenceType `xml:"OCSPs,omitempty"`
}

// SATimestampType is the Go form of the generated JAXB class SATimestampType.
type SATimestampType struct {
	AttributeBaseType
	TimeStampValue XSDateTime `xml:"TimeStampValue"`
}

// SAVRIType is the Go form of the generated JAXB class SAVRIType.
type SAVRIType struct {
	AttributeBaseType
	Certs *VOReferenceType `xml:"Certs,omitempty"`
	CRLs  *VOReferenceType `xml:"CRLs,omitempty"`
	OCSPs *VOReferenceType `xml:"OCSPs,omitempty"`
	TU    *string          `xml:"TU,omitempty"`
	TS    *SATimestampType `xml:"TS,omitempty"`
}

// SAReasonType is the Go form of the generated JAXB class SAReasonType.
type SAReasonType struct {
	AttributeBaseType
	ReasonElement string `xml:"ReasonElement"`
}

// SANameType is the Go form of the generated JAXB class SANameType.
type SANameType struct {
	AttributeBaseType
	NameElement string `xml:"NameElement"`
}

// SAContactInfoType is the Go form of the generated JAXB class SAContactInfoType.
type SAContactInfoType struct {
	AttributeBaseType
	ContactInfoElement string `xml:"ContactInfoElement"`
}

// SASubFilterType is the Go form of the generated JAXB class SASubFilterType.
type SASubFilterType struct {
	AttributeBaseType
	SubFilterElement string `xml:"SubFilterElement"`
}

// SAFilterType is the Go form of the generated JAXB class SAFilterType.
type SAFilterType struct {
	Filter string `xml:"Filter"`
}

// SAMessageDigestType is the Go form of the generated JAXB class SAMessageDigestType.
type SAMessageDigestType struct {
	AttributeBaseType
	Digest Base64Binary `xml:"Digest"`
}

// SAListOfIntegers is the Go form of the schema simpleType SAListOfIntegersType
// (xs:list itemType="xs:integer"), bound directly to List<BigInteger> - it has
// no generated wrapper class of its own (RenewedDigests, ByteRange).
type SAListOfIntegers []*big.Int

// MarshalText joins the members with a single space.
func (l SAListOfIntegers) MarshalText() ([]byte, error) {
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
func (l *SAListOfIntegers) UnmarshalText(text []byte) error {
	fields := strings.Fields(string(text))
	out := make(SAListOfIntegers, 0, len(fields))
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

// ------------------------------------------------------- xs:choice machinery

// choiceItem is one member of a schema xs:choice this package marshals by
// hand: the element name selected, paired with a pointer to the concrete Go
// value decoded (or to be encoded) under that name. See this file's header.
type choiceItem struct {
	Name  string
	Value any
}

// ChoiceItem exposes choiceItem to the report builders that populate these
// models (dss/validation/executor's ETSIValidationReportBuilder, the Go
// counterpart of Java's ETSIValidationReportBuilder). Java's generated classes
// expose the choice as a plain List<Object> the builder fills directly; this
// port needs the element name recorded alongside the value, so the builder
// constructs items through NewChoiceItem instead. Marshalling behaviour is
// unchanged - this is an alias plus a constructor.
type ChoiceItem = choiceItem

// NewChoiceItem records value as the choice member selected under the schema
// element name.
func NewChoiceItem(name string, value any) ChoiceItem {
	return choiceItem{Name: name, Value: value}
}

// choiceElement is one name->type binding of an xs:choice's name table.
type choiceElement struct {
	name string
	typ  reflect.Type
}

func choiceType(table []choiceElement, name string) reflect.Type {
	for _, e := range table {
		if e.name == name {
			return e.typ
		}
	}
	return nil
}

// encodeChoiceItems writes each item under its recorded element name.
func encodeChoiceItems(e *xml.Encoder, items []choiceItem) error {
	for _, it := range items {
		if err := e.EncodeElement(it.Value, xml.StartElement{Name: xml.Name{Local: it.Name}}); err != nil {
			return err
		}
	}
	return nil
}

// decodeChoiceItems reads choice members until the enclosing element's end
// tag, building each item from the type its element name is bound to in
// table. Unrecognised element names are skipped, mirroring xjc's generated
// unmarshaller ignoring content it has no binding for.
func decodeChoiceItems(d *xml.Decoder, table []choiceElement) ([]choiceItem, error) {
	var items []choiceItem
	for {
		tok, err := d.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			typ := choiceType(table, t.Name.Local)
			if typ == nil {
				if err := d.Skip(); err != nil {
					return nil, err
				}
				continue
			}
			pv := reflect.New(typ)
			if err := d.DecodeElement(pv.Interface(), &t); err != nil {
				return nil, err
			}
			items = append(items, choiceItem{Name: t.Name.Local, Value: pv.Interface()})
		case xml.EndElement:
			return items, nil
		}
	}
}

// signatureAttributeElements is the name->type table of
// SignatureAttributesType's 35-member choice, in schema order.
var signatureAttributeElements = []choiceElement{
	{"SigningTime", reflect.TypeOf(SASigningTimeType{})},
	{"SigningCertificate", reflect.TypeOf(SACertIDListType{})},
	{"DataObjectFormat", reflect.TypeOf(SADataObjectFormatType{})},
	{"CommitmentTypeIndication", reflect.TypeOf(SACommitmentTypeIndicationType{})},
	{"AllDataObjectsTimeStamp", reflect.TypeOf(SATimestampType{})},
	{"IndividualDataObjectsTimeStamp", reflect.TypeOf(SATimestampType{})},
	{"SigPolicyIdentifier", reflect.TypeOf(SASigPolicyIdentifierType{})},
	{"SignatureProductionPlace", reflect.TypeOf(SASignatureProductionPlaceType{})},
	{"SignerRole", reflect.TypeOf(SASignerRoleType{})},
	{"CounterSignature", reflect.TypeOf(SACounterSignatureType{})},
	{"SignatureTimeStamp", reflect.TypeOf(SATimestampType{})},
	{"CompleteCertificateRefs", reflect.TypeOf(SACertIDListType{})},
	{"CompleteRevocationRefs", reflect.TypeOf(SARevIDListType{})},
	{"AttributeCertificateRefs", reflect.TypeOf(SACertIDListType{})},
	{"AttributeRevocationRefs", reflect.TypeOf(SARevIDListType{})},
	{"SigAndRefsTimeStamp", reflect.TypeOf(SATimestampType{})},
	{"RefsOnlyTimeStamp", reflect.TypeOf(SATimestampType{})},
	{"CertificateValues", reflect.TypeOf(AttributeBaseType{})},
	{"RevocationValues", reflect.TypeOf(AttributeBaseType{})},
	{"AttrAuthoritiesCertValues", reflect.TypeOf(AttributeBaseType{})},
	{"AttributeRevocationValues", reflect.TypeOf(AttributeBaseType{})},
	{"TimeStampValidationData", reflect.TypeOf(AttributeBaseType{})},
	{"ArchiveTimeStamp", reflect.TypeOf(SATimestampType{})},
	{"RenewedDigests", reflect.TypeOf(SAListOfIntegers{})},
	{"MessageDigest", reflect.TypeOf(SAMessageDigestType{})},
	{"DSS", reflect.TypeOf(SADSSType{})},
	{"VRI", reflect.TypeOf(SAVRIType{})},
	{"DocTimeStamp", reflect.TypeOf(SATimestampType{})},
	{"Reason", reflect.TypeOf(SAReasonType{})},
	{"Name", reflect.TypeOf(SANameType{})},
	{"ContactInfo", reflect.TypeOf(SAContactInfoType{})},
	{"SubFilter", reflect.TypeOf(SASubFilterType{})},
	{"ByteRange", reflect.TypeOf(SAListOfIntegers{})},
	{"Filter", reflect.TypeOf(SAFilterType{})},
	{"OtherInformation", reflect.TypeOf(AnyType{})},
}

// SignatureAttributesType is the Go form of the generated JAXB class
// SignatureAttributesType (complexType SignatureAttributesType, the
// 35-member choice described in this file's header).
type SignatureAttributesType struct {
	Items []choiceItem
}

// MarshalXML writes each attribute item under the element name it was
// recorded with.
func (t *SignatureAttributesType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := encodeChoiceItems(e, t.Items); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}

// UnmarshalXML builds each attribute item from the type its element name is
// bound to.
func (t *SignatureAttributesType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	items, err := decodeChoiceItems(d, signatureAttributeElements)
	if err != nil {
		return err
	}
	t.Items = items
	return nil
}
