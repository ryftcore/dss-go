// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. Per the phase-8a generated-JAXB rule the
// generated classes are grouped into schema-area files rather than one file per
// class; every Java class keeps its name and its exact field order so that
// encoding/xml reproduces the JAXB element sequence byte for byte.

package jaxb

// XmlCertForPID is the Go form of the generated JAXB class XmlCertForPID
// (complexType CertForPID).
type XmlCertForPID struct {
	Present bool `xml:"present,attr"`
}

// XmlCertForWallet is the Go form of the generated JAXB class XmlCertForWallet
// (complexType CertForWallet).
type XmlCertForWallet struct {
	Present bool `xml:"present,attr"`
}

// XmlLangAndValue is the Go form of the generated JAXB class XmlLangAndValue
// (complexType LangAndValue).
type XmlLangAndValue struct {
	Value string  `xml:",chardata"`
	Lang  *string `xml:"lang,attr,omitempty"`
}

// XmlMRACertificateMapping is the Go form of the generated JAXB class XmlMRACertificateMapping
// (complexType MRACertificateMapping).
type XmlMRACertificateMapping struct {
	TrustServiceEquivalenceInformation *XmlTrustServiceEquivalenceInformation      `xml:"TrustServiceEquivalenceInformation,omitempty"`
	OriginalThirdCountryMapping        *XmlOriginalThirdCountryQcStatementsMapping `xml:"OriginalThirdCountryMapping,omitempty"`
}

// XmlOriginalThirdCountryQcStatementsMapping is the Go form of the generated JAXB class XmlOriginalThirdCountryQcStatementsMapping
// (complexType OriginalThirdCountryQcStatementsMapping).
type XmlOriginalThirdCountryQcStatementsMapping struct {
	QcCompliance      *XmlQcCompliance          `xml:"QcCompliance,omitempty"`
	QcSSCD            *XmlQcSSCD                `xml:"QcSSCD,omitempty"`
	QcTypes           *QcTypesWrapper           `xml:"QcTypes"`
	QcCClegislation   *QcCClegislationWrapper   `xml:"QcCClegislation"`
	QcQSCDlegislation *QcQSCDlegislationWrapper `xml:"QcQSCDlegislation"`
	OtherOIDs         *OtherOIDsWrapper         `xml:"OtherOIDs"`
}

// XmlPSD2QcInfo is the Go form of the generated JAXB class XmlPSD2QcInfo
// (complexType PSD2QcInfo).
type XmlPSD2QcInfo struct {
	RolesOfPSP *RolesOfPSPWrapper `xml:"RolesOfPSP"`
	NcaName    *string            `xml:"ncaName,omitempty"`
	NcaId      *string            `xml:"ncaId,omitempty"`
}

// XmlQcCompliance is the Go form of the generated JAXB class XmlQcCompliance
// (complexType QcCompliance).
type XmlQcCompliance struct {
	Present bool `xml:"present,attr"`
}

// XmlQcEuLimitValue is the Go form of the generated JAXB class XmlQcEuLimitValue
// (complexType QcEuLimitValue).
type XmlQcEuLimitValue struct {
	Currency *string `xml:"currency,omitempty"`
	Amount   int     `xml:"amount"`
	Exponent int     `xml:"exponent"`
}

// XmlQcPSB is the Go form of the generated JAXB class XmlQcPSB
// (complexType QcPSB).
type XmlQcPSB struct {
	CountryOfLegislation      *string `xml:"countryOfLegislation,omitempty"`
	AuthSourceIdentification  *string `xml:"authSourceIdentification,omitempty"`
	LegislationIdentification *string `xml:"legislationIdentification,omitempty"`
}

// XmlQcSSCD is the Go form of the generated JAXB class XmlQcSSCD
// (complexType QcSSCD).
type XmlQcSSCD struct {
	Present bool `xml:"present,attr"`
}

// XmlRoleOfPSP is the Go form of the generated JAXB class XmlRoleOfPSP
// (complexType RoleOfPSP).
type XmlRoleOfPSP struct {
	Oid  *XmlOID `xml:"oid,omitempty"`
	Name *string `xml:"name,omitempty"`
}

// XmlQcStatements is the Go form of the generated JAXB class XmlQcStatements
// (complexType QcStatements).
type XmlQcStatements struct {
	XmlCertificateExtensionContent
	QcCompliance          *XmlQcCompliance          `xml:"QcCompliance,omitempty"`
	QcEuLimitValue        *XmlQcEuLimitValue        `xml:"QcEuLimitValue,omitempty"`
	QcEuRetentionPeriod   *int                      `xml:"QcEuRetentionPeriod,omitempty"`
	QcSSCD                *XmlQcSSCD                `xml:"QcSSCD,omitempty"`
	QcEuPDS               *QcEuPDSWrapper           `xml:"QcEuPDS"`
	QcTypes               *QcTypesWrapper           `xml:"QcTypes"`
	QcCClegislation       *QcCClegislationWrapper   `xml:"QcCClegislation"`
	SemanticsIdentifier   *XmlOID                   `xml:"SemanticsIdentifier,omitempty"`
	PSD2QcInfo            *XmlPSD2QcInfo            `xml:"PSD2QcInfo,omitempty"`
	QcQSCDlegislation     *QcQSCDlegislationWrapper `xml:"QcQSCDlegislation"`
	QcIdentMethod         *XmlOID                   `xml:"QcIdentMethod,omitempty"`
	QcPSB                 *XmlQcPSB                 `xml:"QcPSB,omitempty"`
	OtherOIDs             *OtherOIDsWrapper         `xml:"OtherOIDs"`
	MRACertificateMapping *XmlMRACertificateMapping `xml:"MRACertificateMapping,omitempty"`
	EnactedMRA            *bool                     `xml:"enactedMRA,attr,omitempty"`
	XmlCertificateExtensionAttrs
}
