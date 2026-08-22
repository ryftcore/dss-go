// Ported from ts_119612v020401_sie_xsd.xsd (DSS 6.5.RC1) via the JAXB
// classes generated into eu.europa.esig.trustedlist.jaxb.ecc (Qualified
// Certificate qualifiers / "SvcInfoExt" extension). The generated classes
// are grouped into one file per Java package; every Java class keeps its
// name and its exact field order so that encoding/xml reproduces the JAXB
// element sequence byte for byte.
//
// ecc's XmlAdapter classes (Adapter1/Adapter2, generated from the
// KeyUsageBitType.name/CriteriaListType.assert attributes' xs:string ->
// enum bindings) have no file of their own - see doc.go's header - and are
// folded into keyUsageBitAttr/assertAttr below: small unexported wrappers
// around the already-ported dss/enumerations.KeyUsageBit/Assert (this
// package cannot add methods to a type it does not own), implementing
// MarshalXMLAttr/UnmarshalXMLAttr the way KeyUsageBitParser/AssertParser
// (dss-jaxb-parsers, not in this manifest) print/parse them: by the
// enum's .Value() lexical form, with an unresolvable value silently
// producing the zero KeyUsageBit/Assert - the RI's XmlAdapter.unmarshal
// returning null is likewise non-fatal (see
// dss/diagnostic/jaxb/xml.go's XSDateTime for the same
// unmarshaller-tolerance rationale).
package jaxb

import (
	"encoding/xml"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// keyUsageBitAttr adapts enumerations.KeyUsageBit to KeyUsageBitType's
// "name" attribute lexical form (ecc.Adapter1 / KeyUsageBitParser).
type keyUsageBitAttr enumerations.KeyUsageBit

// MarshalXMLAttr writes the KeyUsageBit's XSD lexical value.
func (k keyUsageBitAttr) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: enumerations.KeyUsageBit(k).Value()}, nil
}

// UnmarshalXMLAttr resolves the lexical value, mirroring KeyUsageBitParser.parse.
func (k *keyUsageBitAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	for _, v := range enumerations.KeyUsageBitValues() {
		if v.Value() == attr.Value {
			*k = keyUsageBitAttr(v)
			return nil
		}
	}
	*k = ""
	return nil
}

// assertAttr adapts enumerations.Assert to CriteriaListType's "assert"
// attribute lexical form (ecc.Adapter2 / AssertParser).
type assertAttr enumerations.Assert

// MarshalXMLAttr writes the Assert's XSD lexical value.
func (a assertAttr) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: enumerations.Assert(a).Value()}, nil
}

// UnmarshalXMLAttr resolves the lexical value, mirroring AssertParser.parse.
func (a *assertAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	for _, v := range enumerations.AssertValues() {
		if v.Value() == attr.Value {
			*a = assertAttr(v)
			return nil
		}
	}
	*a = ""
	return nil
}

// QualifierType is the Go form of the generated JAXB class QualifierType
// (complexType QualifierType).
type QualifierType struct {
	URI *string `xml:"uri,attr,omitempty"`
}

// QualifiersType is the Go form of the generated JAXB class QualifiersType
// (complexType QualifiersType).
type QualifiersType struct {
	Qualifier []*QualifierType `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# Qualifier"`
}

// KeyUsageBitType is the Go form of the generated JAXB class
// KeyUsageBitType (complexType KeyUsageBitType).
type KeyUsageBitType struct {
	Value bool             `xml:",chardata"`
	Name  *keyUsageBitAttr `xml:"name,attr,omitempty"`
}

// KeyUsageType is the Go form of the generated JAXB class KeyUsageType
// (complexType KeyUsageType).
type KeyUsageType struct {
	KeyUsageBit []*KeyUsageBitType `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# KeyUsageBit"`
}

// PoliciesListType is the Go form of the generated JAXB class
// PoliciesListType (complexType PoliciesListType).
type PoliciesListType struct {
	PolicyIdentifier []*ObjectIdentifierType `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# PolicyIdentifier"`
}

// CriteriaListType is the Go form of the generated JAXB class
// CriteriaListType (complexType CriteriaListType). OtherCriteriaList is
// eu.europa.esig.xades.jaxb.xades132.AnyType, outside this manifest - see
// jaxb_common.go's header - bound, since the field carries no @XmlElement
// of its own, to the literal field name "otherCriteriaList" in ecc's own
// namespace (JAXB's XmlAccessorType.FIELD default).
type CriteriaListType struct {
	KeyUsage          []*KeyUsageType     `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# KeyUsage"`
	PolicySet         []*PoliciesListType `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# PolicySet"`
	CriteriaList      []*CriteriaListType `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# CriteriaList"`
	Description       *string             `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# Description,omitempty"`
	OtherCriteriaList *xadesAnyType       `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# otherCriteriaList,omitempty"`
	Assert            *assertAttr         `xml:"assert,attr,omitempty"`
}

// QualificationElementType is the Go form of the generated JAXB class
// QualificationElementType (complexType QualificationElementType).
type QualificationElementType struct {
	Qualifiers   *QualifiersType   `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# Qualifiers"`
	CriteriaList *CriteriaListType `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# CriteriaList"`
}

// QualificationsType is the Go form of the generated JAXB class
// QualificationsType (complexType QualificationsType, the root of the
// ecc:Qualifications extension).
type QualificationsType struct {
	QualificationElement []*QualificationElementType `xml:"http://uri.etsi.org/TrstSvc/SvcInfoExt/eSigDir-1999-93-EC-TrustedList/# QualificationElement"`
}
