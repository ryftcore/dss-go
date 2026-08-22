// Ported from ts_119612v020401_additionaltypes_xsd.xsd (DSS 6.5.RC1) via
// the JAXB classes generated into eu.europa.esig.trustedlist.jaxb.tslx
// ("additional types" extension). The generated classes are grouped into
// one file per Java package; every Java class keeps its name and its exact
// field order so that encoding/xml reproduces the JAXB element sequence
// byte for byte.
package jaxb

// CertSubjectDNAttributeType is the Go form of the generated JAXB class
// CertSubjectDNAttributeType (complexType CertSubjectDNAttributeType).
type CertSubjectDNAttributeType struct {
	AttributeOID []*ObjectIdentifierType `xml:"http://uri.etsi.org/02231/v2/additionaltypes# AttributeOID"`
}

// ExtendedKeyUsageType is the Go form of the generated JAXB class
// ExtendedKeyUsageType (complexType ExtendedKeyUsageType).
type ExtendedKeyUsageType struct {
	KeyPurposeId []*ObjectIdentifierType `xml:"http://uri.etsi.org/02231/v2/additionaltypes# KeyPurposeId"`
}

// TakenOverByType is the Go form of the generated JAXB class TakenOverByType
// (complexType TakenOverByType). URI, TSPName and OtherQualifier's element
// type are this module's own tsl package (tsl.NonEmptyMultiLangURIType,
// tsl.InternationalNamesType, tsl.AnyType) but, per JAXB's
// XmlAccessorType.FIELD/elementFormDefault="qualified" default, their
// ELEMENT is qualified into tslx's own namespace, not tsl's; only
// SchemeOperatorName and SchemeTerritory are explicitly bound into the tsl
// namespace (an explicit @XmlElement(namespace=...) override in the
// generated Java).
type TakenOverByType struct {
	URI                *NonEmptyMultiLangURIType `xml:"http://uri.etsi.org/02231/v2/additionaltypes# URI"`
	TSPName            *InternationalNamesType   `xml:"http://uri.etsi.org/02231/v2/additionaltypes# TSPName"`
	SchemeOperatorName *InternationalNamesType   `xml:"http://uri.etsi.org/02231/v2# SchemeOperatorName"`
	SchemeTerritory    string                    `xml:"http://uri.etsi.org/02231/v2# SchemeTerritory"`
	OtherQualifier     []*AnyType                `xml:"http://uri.etsi.org/02231/v2/additionaltypes# OtherQualifier,omitempty"`
}
