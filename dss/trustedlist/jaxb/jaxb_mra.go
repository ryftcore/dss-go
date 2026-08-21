// Ported from mra_schema_v2_19612v020401.xsd (DSS 6.5.RC1) via the JAXB
// classes generated into eu.europa.esig.trustedlist.jaxb.mra (Mutual
// Recognition Agreement extension). Per the phase-8a generated-JAXB rule
// the generated classes are grouped into one file per Java package; every
// Java class keeps its name and its exact field order so that encoding/xml
// reproduces the JAXB element sequence byte for byte.
//
// mra's XmlAdapter classes (Adapter1/Adapter3, both MRAStatus; Adapter2,
// MRAEquivalenceContext) have no file of their own - see doc.go's header -
// and are folded into mraStatusAttr/mraEquivalenceContextAttr below, the
// same pattern as jaxb_ecc.go's keyUsageBitAttr/assertAttr. They print/parse
// through the already-ported dss/enumerations.MRAStatus/
// MRAEquivalenceContext's own .URI() (matching Java's MRAStatus.getUri()/
// MRAEquivalenceContext.getUri()) rather than importing dss/trustedlist's
// MRAStatusParser/MRAEquivalenceContextParser: this package sits BELOW
// dss/trustedlist (which wraps it into the TrustedListFacade/MRAFacade
// API), so importing back up would cycle - the parse/print logic itself is
// a two-line loop, cheaply duplicated on both sides of that boundary
// exactly as jaxb_ecc.go already does for dss-jaxb-parsers'
// KeyUsageBitParser/AssertParser.
package jaxb

import (
	"encoding/xml"

	"github.com/utain/esig/dss/enumerations"
)

// mraStatusAttr adapts enumerations.MRAStatus to the xs:anyURI attribute
// lexical form mra.Adapter1/Adapter3 (MRAStatusParser) use.
type mraStatusAttr enumerations.MRAStatus

// MarshalXMLAttr writes the MRAStatus's URI.
func (s mraStatusAttr) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: enumerations.MRAStatus(s).URI()}, nil
}

// UnmarshalXMLAttr resolves the URI, mirroring MRAStatusParser.parse.
func (s *mraStatusAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	for _, v := range enumerations.MRAStatusValues() {
		if v.URI() == attr.Value {
			*s = mraStatusAttr(v)
			return nil
		}
	}
	*s = ""
	return nil
}

// mraEquivalenceContextAttr adapts enumerations.MRAEquivalenceContext to the
// xs:anyURI attribute lexical form mra.Adapter2
// (MRAEquivalenceContextParser) uses.
type mraEquivalenceContextAttr enumerations.MRAEquivalenceContext

// MarshalXMLAttr writes the MRAEquivalenceContext's URI.
func (c mraEquivalenceContextAttr) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: enumerations.MRAEquivalenceContext(c).URI()}, nil
}

// UnmarshalXMLAttr resolves the URI, mirroring MRAEquivalenceContextParser.parse.
func (c *mraEquivalenceContextAttr) UnmarshalXMLAttr(attr xml.Attr) error {
	for _, v := range enumerations.MRAEquivalenceContextValues() {
		if v.URI() == attr.Value {
			*c = mraEquivalenceContextAttr(v)
			return nil
		}
	}
	*c = ""
	return nil
}

// CertificateContentReferenceEquivalenceType is the Go form of the
// generated JAXB class CertificateContentReferenceEquivalenceType
// (complexType CertificateContentReferenceEquivalenceType). The
// CertificateContentReferenceEquivalenceContext element is bound
// String-typed with an XmlAdapter in Java (@XmlElement(type = String.class)
// @XmlJavaTypeAdapter(Adapter2.class)); here it is the enum itself, with
// MarshalXML/UnmarshalXML doing what the adapter did.
type CertificateContentReferenceEquivalenceType struct {
	CertificateContentReferenceEquivalenceContext mraEquivalenceContextElement `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# CertificateContentReferenceEquivalenceContext"`
	CertificateContentDeclarationPointingParty    *CriteriaListType            `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# CertificateContentDeclarationPointingParty"`
	CertificateContentDeclarationPointedParty     *CriteriaListType            `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# CertificateContentDeclarationPointedParty"`
}

// mraEquivalenceContextElement is CertificateContentReferenceEquivalenceContext's
// element-position counterpart to mraEquivalenceContextAttr (same enum,
// bound as element text rather than an attribute).
type mraEquivalenceContextElement enumerations.MRAEquivalenceContext

// MarshalText writes the MRAEquivalenceContext's URI.
func (c mraEquivalenceContextElement) MarshalText() ([]byte, error) {
	return []byte(enumerations.MRAEquivalenceContext(c).URI()), nil
}

// UnmarshalText resolves the URI, mirroring MRAEquivalenceContextParser.parse.
func (c *mraEquivalenceContextElement) UnmarshalText(text []byte) error {
	for _, v := range enumerations.MRAEquivalenceContextValues() {
		if v.URI() == string(text) {
			*c = mraEquivalenceContextElement(v)
			return nil
		}
	}
	*c = ""
	return nil
}

// CertificateContentReferencesEquivalenceListType is the Go form of the
// generated JAXB class CertificateContentReferencesEquivalenceListType
// (complexType CertificateContentReferencesEquivalenceListType).
type CertificateContentReferencesEquivalenceListType struct {
	CertificateContentReferenceEquivalence []*CertificateContentReferenceEquivalenceType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# CertificateContentReferenceEquivalence"`
}

// MutualRecognitionAgreementInformationType is the Go form of the generated
// JAXB class MutualRecognitionAgreementInformationType (complexType
// MutualRecognitionAgreementInformationType, the root of the
// mra:MutualRecognitionAgreementInformation extension).
type MutualRecognitionAgreementInformationType struct {
	TrustServiceEquivalenceInformation  []*TrustServiceEquivalenceInformationType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceEquivalenceInformation"`
	TechnicalType                       *BigInteger                               `xml:"technicalType,attr"`
	Version                             *BigInteger                               `xml:"version,attr"`
	PointingContractingPartyLegislation string                                    `xml:"pointingContractingPartyLegislation,attr"`
	PointedContractingPartyLegislation  string                                    `xml:"pointedContractingPartyLegislation,attr"`
	MRADepth                            *BigInteger                               `xml:"MRADepth,attr"`
}

// QcStatementInfoType is the Go form of the generated JAXB class
// QcStatementInfoType (complexType QcStatementInfoType).
type QcStatementInfoType struct {
	QcType          *ObjectIdentifierType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# QcType,omitempty"`
	QcCClegislation *string               `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# QcCClegislation,omitempty"`
}

// QcStatementListType is the Go form of the generated JAXB class
// QcStatementListType (complexType QcStatementListType).
type QcStatementListType struct {
	QcStatement []*QcStatementType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# QcStatement"`
}

// QcStatementType is the Go form of the generated JAXB class QcStatementType
// (complexType QcStatementType).
type QcStatementType struct {
	QcStatementId   *ObjectIdentifierType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# QcStatementId"`
	QcStatementInfo *QcStatementInfoType  `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# QcStatementInfo,omitempty"`
}

// QualifierEquivalenceListType is the Go form of the generated JAXB class
// QualifierEquivalenceListType (complexType QualifierEquivalenceListType).
type QualifierEquivalenceListType struct {
	QualifierEquivalence []*QualifierEquivalenceType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# QualifierEquivalence"`
}

// QualifierEquivalenceType is the Go form of the generated JAXB class
// QualifierEquivalenceType (complexType QualifierEquivalenceType).
type QualifierEquivalenceType struct {
	QualifierPointingParty *QualifierType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# QualifierPointingParty"`
	QualifierPointedParty  *QualifierType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# QualifierPointedParty"`
}

// TrustServiceEquivalenceHistoryInstanceType is the Go form of the
// generated JAXB class TrustServiceEquivalenceHistoryInstanceType
// (complexType TrustServiceEquivalenceHistoryInstanceType).
// TrustServiceEquivalenceStatusStartingTime is the raw xs:dateTime lexical
// form, round-tripped verbatim - see jaxb_tsl_root.go's header.
type TrustServiceEquivalenceHistoryInstanceType struct {
	TrustServiceTSLTypeEquivalenceList                   *TrustServiceTSLTypeEquivalenceListType                   `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLTypeEquivalenceList"`
	TrustServiceEquivalenceStatus                        *mraStatusElement                                         `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceEquivalenceStatus"`
	TrustServiceEquivalenceStatusStartingTime            *string                                                   `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceEquivalenceStatusStartingTime,omitempty"`
	TrustServiceTSLStatusEquivalenceList                 *TrustServiceTSLStatusEquivalenceListType                 `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLStatusEquivalenceList"`
	CertificateContentReferencesEquivalenceList          *CertificateContentReferencesEquivalenceListType          `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# CertificateContentReferencesEquivalenceList,omitempty"`
	TrustServiceTSLQualificationExtensionEquivalenceList *TrustServiceTSLQualificationExtensionEquivalenceListType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLQualificationExtensionEquivalenceList,omitempty"`
}

// mraStatusElement is TrustServiceEquivalenceStatus's element-position
// counterpart to mraStatusAttr (same enum, bound as element text rather
// than an attribute).
type mraStatusElement enumerations.MRAStatus

// MarshalText writes the MRAStatus's URI.
func (s mraStatusElement) MarshalText() ([]byte, error) {
	return []byte(enumerations.MRAStatus(s).URI()), nil
}

// UnmarshalText resolves the URI, mirroring MRAStatusParser.parse.
func (s *mraStatusElement) UnmarshalText(text []byte) error {
	for _, v := range enumerations.MRAStatusValues() {
		if v.URI() == string(text) {
			*s = mraStatusElement(v)
			return nil
		}
	}
	*s = ""
	return nil
}

// TrustServiceEquivalenceHistoryType is the Go form of the generated JAXB
// class TrustServiceEquivalenceHistoryType (complexType
// TrustServiceEquivalenceHistoryType).
type TrustServiceEquivalenceHistoryType struct {
	TrustServiceEquivalenceHistoryInstance []*TrustServiceEquivalenceHistoryInstanceType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceEquivalenceHistoryInstance"`
}

// TrustServiceEquivalenceInformationType is the Go form of the generated
// JAXB class TrustServiceEquivalenceInformationType (complexType
// TrustServiceEquivalenceInformationType).
// TrustServiceEquivalenceStatusStartingTime is the raw xs:dateTime lexical
// form, round-tripped verbatim - see jaxb_tsl_root.go's header.
type TrustServiceEquivalenceInformationType struct {
	TrustServiceLegalIdentifier                          string                                                    `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceLegalIdentifier"`
	TrustServiceTSLTypeEquivalenceList                   *TrustServiceTSLTypeEquivalenceListType                   `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLTypeEquivalenceList"`
	TrustServiceEquivalenceStatus                        *mraStatusElement                                         `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceEquivalenceStatus"`
	TrustServiceEquivalenceStatusStartingTime            *string                                                   `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceEquivalenceStatusStartingTime,omitempty"`
	TrustServiceTSLStatusEquivalenceList                 *TrustServiceTSLStatusEquivalenceListType                 `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLStatusEquivalenceList"`
	CertificateContentReferencesEquivalenceList          *CertificateContentReferencesEquivalenceListType          `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# CertificateContentReferencesEquivalenceList,omitempty"`
	TrustServiceTSLQualificationExtensionEquivalenceList *TrustServiceTSLQualificationExtensionEquivalenceListType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLQualificationExtensionEquivalenceList,omitempty"`
	TrustServiceEquivalenceHistory                       []*TrustServiceEquivalenceHistoryType                     `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceEquivalenceHistory,omitempty"`
}

// TrustServiceTSLQualificationExtensionEquivalenceListType is the Go form of
// the generated JAXB class
// TrustServiceTSLQualificationExtensionEquivalenceListType (complexType
// TrustServiceTSLQualificationExtensionEquivalenceListType).
type TrustServiceTSLQualificationExtensionEquivalenceListType struct {
	TrustServiceTSLQualificationExtensionName *TrustServiceTSLQualificationExtensionNameType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLQualificationExtensionName"`
	QualifierEquivalenceList                  []*QualifierEquivalenceListType                `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# QualifierEquivalenceList"`
}

// TrustServiceTSLQualificationExtensionNameType is the Go form of the
// generated JAXB class TrustServiceTSLQualificationExtensionNameType
// (complexType TrustServiceTSLQualificationExtensionNameType).
type TrustServiceTSLQualificationExtensionNameType struct {
	TrustServiceTSLQualificationExtensionNamePointingParty string `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLQualificationExtensionNamePointingParty"`
	TrustServiceTSLQualificationExtensionNamePointedParty  string `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLQualificationExtensionNamePointedParty"`
}

// TrustServiceTSLStatusEquivalenceListType is the Go form of the generated
// JAXB class TrustServiceTSLStatusEquivalenceListType (complexType
// TrustServiceTSLStatusEquivalenceListType).
type TrustServiceTSLStatusEquivalenceListType struct {
	TrustServiceTSLStatusValidEquivalence   *TrustServiceTSLStatusEquivalenceType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLStatusValidEquivalence"`
	TrustServiceTSLStatusInvalidEquivalence *TrustServiceTSLStatusEquivalenceType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLStatusInvalidEquivalence"`
}

// TrustServiceTSLStatusEquivalenceType is the Go form of the generated JAXB
// class TrustServiceTSLStatusEquivalenceType (complexType
// TrustServiceTSLStatusEquivalenceType).
type TrustServiceTSLStatusEquivalenceType struct {
	TrustServiceTSLStatusListPointingParty *TrustServiceTSLStatusList `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLStatusListPointingParty"`
	TrustServiceTSLStatusListPointedParty  *TrustServiceTSLStatusList `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLStatusListPointedParty"`
}

// TrustServiceTSLStatusList is the Go form of the generated JAXB class
// TrustServiceTSLStatusList (complexType TrustServiceTSLStatusList). Its
// ServiceStatus is explicitly bound into the tsl namespace rather than
// mra's own.
type TrustServiceTSLStatusList struct {
	ServiceStatus []string `xml:"http://uri.etsi.org/02231/v2# ServiceStatus"`
}

// TrustServiceTSLTypeEquivalenceListType is the Go form of the generated
// JAXB class TrustServiceTSLTypeEquivalenceListType (complexType
// TrustServiceTSLTypeEquivalenceListType).
type TrustServiceTSLTypeEquivalenceListType struct {
	TrustServiceTSLTypeListPointingParty *TrustServiceTSLTypeListType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLTypeListPointingParty"`
	TrustServiceTSLTypeListPointedParty  *TrustServiceTSLTypeListType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLTypeListPointedParty"`
}

// TrustServiceTSLTypeListType is the Go form of the generated JAXB class
// TrustServiceTSLTypeListType (complexType TrustServiceTSLTypeListType).
type TrustServiceTSLTypeListType struct {
	TrustServiceTSLType []*TrustServiceTSLTypeType `xml:"http://ec.europa.eu/tools/lotl/mra/schema/v2# TrustServiceTSLType"`
}

// TrustServiceTSLTypeType is the Go form of the generated JAXB class
// TrustServiceTSLTypeType (complexType TrustServiceTSLTypeType). Both
// properties are explicitly bound into the tsl namespace rather than mra's
// own, reusing this module's own tsl.AdditionalServiceInformationType.
type TrustServiceTSLTypeType struct {
	ServiceTypeIdentifier        string                            `xml:"http://uri.etsi.org/02231/v2# ServiceTypeIdentifier"`
	AdditionalServiceInformation *AdditionalServiceInformationType `xml:"http://uri.etsi.org/02231/v2# AdditionalServiceInformation,omitempty"`
}
