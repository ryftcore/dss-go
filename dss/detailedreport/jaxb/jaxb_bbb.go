// Ported from DetailedReport.xsd (DSS 6.5.RC1) via the JAXB classes generated
// into eu.europa.esig.dss.detailedreport.jaxb. See jaxb_common.go's header for
// the file-grouping rule and jaxb_process.go's header for the
// Content/Attrs extension-embedding pattern. This file holds the ETSI
// EN 319 102-1 Basic Building Blocks: BasicBuildingBlocks itself and every
// sub-block (FC, ISC, VCI, XCV, SubXCV, CV, SAV, AOV, PSV, CRS, PCV, VTS, RAC,
// RFC, CC) plus the certificate-chain types they share.

package jaxb

// XmlCertificateChain is the Go form of the generated JAXB class
// XmlCertificateChain (complexType CertificateChain).
type XmlCertificateChain struct {
	ChainItem []*XmlChainItem `xml:"ChainItem,omitempty"`
}

// XmlChainItem is the Go form of the generated JAXB class XmlChainItem (the
// anonymous complexType of CertificateChain/ChainItem).
type XmlChainItem struct {
	Source CertificateSourceTypeValue `xml:"Source"`
	Id     string                     `xml:"Id,attr"`
}

// XmlFC is the Go form of the generated JAXB class XmlFC (complexType FC,
// extends ConstraintsConclusion with no additions).
type XmlFC struct {
	XmlConstraintsConclusionContent
	XmlConstraintsConclusionAttrs
}

// XmlISC is the Go form of the generated JAXB class XmlISC (complexType ISC,
// extends ConstraintsConclusion).
type XmlISC struct {
	XmlConstraintsConclusionContent
	CertificateChain *XmlCertificateChain `xml:"CertificateChain"`
	XmlConstraintsConclusionAttrs
}

// XmlVCI is the Go form of the generated JAXB class XmlVCI (complexType VCI,
// extends ConstraintsConclusion with no additions).
type XmlVCI struct {
	XmlConstraintsConclusionContent
	XmlConstraintsConclusionAttrs
}

// XmlRFC is the Go form of the generated JAXB class XmlRFC (complexType RFC,
// extends ConstraintsConclusion).
type XmlRFC struct {
	XmlConstraintsConclusionContent
	Id *string `xml:"Id,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlSubXCV is the Go form of the generated JAXB class XmlSubXCV (complexType
// SubXCV, extends ConstraintsConclusion).
type XmlSubXCV struct {
	XmlConstraintsConclusionContent
	CrossCertificate      *StringList               `xml:"CrossCertificate,omitempty"`
	EquivalentCertificate *StringList               `xml:"EquivalentCertificate,omitempty"`
	CRS                   *XmlCRS                   `xml:"CRS,omitempty"`
	RFC                   *XmlRFC                   `xml:"RFC,omitempty"`
	RevocationInfo        *XmlRevocationInformation `xml:"RevocationInfo,omitempty"`
	Id                    string                    `xml:"Id,attr"`
	TrustAnchor           *bool                     `xml:"TrustAnchor,attr,omitempty"`
	SelfSigned            *bool                     `xml:"SelfSigned,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlXCV is the Go form of the generated JAXB class XmlXCV (complexType XCV,
// extends ConstraintsConclusion).
type XmlXCV struct {
	XmlConstraintsConclusionContent
	SubXCV []*XmlSubXCV `xml:"SubXCV,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlCV is the Go form of the generated JAXB class XmlCV (complexType CV,
// extends ConstraintsConclusion with no additions).
type XmlCV struct {
	XmlConstraintsConclusionContent
	XmlConstraintsConclusionAttrs
}

// XmlSAV is the Go form of the generated JAXB class XmlSAV (complexType SAV,
// extends ConstraintsConclusion with no additions).
type XmlSAV struct {
	XmlConstraintsConclusionContent
	XmlConstraintsConclusionAttrs
}

// XmlCertificateChainCryptographicValidation is the Go form of the generated
// JAXB class XmlCertificateChainCryptographicValidation (the anonymous
// complexType of AOV/CertificateChainCryptographicValidation).
type XmlCertificateChainCryptographicValidation struct {
	CertificateCryptographicValidation []*XmlCryptographicValidation `xml:"CertificateCryptographicValidation"`
}

// XmlAOV is the Go form of the generated JAXB class XmlAOV (complexType AOV,
// extends ConstraintsConclusion).
type XmlAOV struct {
	XmlConstraintsConclusionContent
	SignatureCryptographicValidation        *XmlCryptographicValidation                 `xml:"SignatureCryptographicValidation,omitempty"`
	SignedAttributesValidation              *XmlCryptographicValidation                 `xml:"SignedAttributesValidation,omitempty"`
	DigestMatchersValidation                *XmlCryptographicValidation                 `xml:"DigestMatchersValidation,omitempty"`
	CertificateChainCryptographicValidation *XmlCertificateChainCryptographicValidation `xml:"CertificateChainCryptographicValidation,omitempty"`
	ValidationTime                          *XSDateTime                                 `xml:"ValidationTime,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlPSV is the Go form of the generated JAXB class XmlPSV (complexType PSV,
// extends ConstraintsConclusionWithControlTime with no additions).
type XmlPSV struct {
	XmlConstraintsConclusionWithControlTimeContent
	XmlConstraintsConclusionAttrs
}

// XmlRAC is the Go form of the generated JAXB class XmlRAC (complexType RAC,
// extends ConstraintsConclusion).
type XmlRAC struct {
	XmlConstraintsConclusionContent
	RevocationThisUpdate     XSDateTime `xml:"RevocationThisUpdate"`
	RevocationProductionDate XSDateTime `xml:"RevocationProductionDate"`
	CRS                      *XmlCRS    `xml:"CRS,omitempty"`
	Id                       *string    `xml:"Id,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlCRS is the Go form of the generated JAXB class XmlCRS (complexType CRS,
// "CertificateRevocationSelector", extends ConstraintsConclusion).
type XmlCRS struct {
	XmlConstraintsConclusionContent
	RAC                          []*XmlRAC   `xml:"RAC,omitempty"`
	AcceptableRevocationId       *StringList `xml:"AcceptableRevocationId,omitempty"`
	Id                           *string     `xml:"Id,attr,omitempty"`
	LatestAcceptableRevocationId *string     `xml:"LatestAcceptableRevocationId,attr,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlPCV is the Go form of the generated JAXB class XmlPCV (complexType PCV,
// "Past Certificate Validation", extends ConstraintsConclusionWithControlTime
// with no additions).
type XmlPCV struct {
	XmlConstraintsConclusionWithControlTimeContent
	XmlConstraintsConclusionAttrs
}

// XmlVTS is the Go form of the generated JAXB class XmlVTS (complexType VTS,
// "Validation Time Sliding", extends ConstraintsConclusionWithControlTime).
type XmlVTS struct {
	XmlConstraintsConclusionWithControlTimeContent
	TrustAnchor *string   `xml:"TrustAnchor,omitempty"`
	CRS         []*XmlCRS `xml:"CRS,omitempty"`
	XmlConstraintsConclusionAttrs
}

// XmlCC is the Go form of the generated JAXB class XmlCC (complexType CC,
// "Cryptographic Checker", extends ConstraintsConclusion).
type XmlCC struct {
	XmlConstraintsConclusionContent
	CryptographicValidation *XmlCryptographicValidation `xml:"CryptographicValidation"`
	XmlConstraintsConclusionAttrs
}

// XmlBasicBuildingBlocks is the Go form of the generated JAXB class
// XmlBasicBuildingBlocks (complexType BasicBuildingBlocks).
type XmlBasicBuildingBlocks struct {
	FC               *XmlFC               `xml:"FC,omitempty"`
	ISC              *XmlISC              `xml:"ISC,omitempty"`
	VCI              *XmlVCI              `xml:"VCI,omitempty"`
	XCV              *XmlXCV              `xml:"XCV,omitempty"`
	CV               *XmlCV               `xml:"CV,omitempty"`
	SAV              *XmlSAV              `xml:"SAV,omitempty"`
	AOV              *XmlAOV              `xml:"AOV,omitempty"`
	PSV              *XmlPSV              `xml:"PSV,omitempty"`
	PSVCRS           *XmlCRS              `xml:"PSV_CRS,omitempty"`
	PCV              *XmlPCV              `xml:"PCV,omitempty"`
	VTS              *XmlVTS              `xml:"VTS,omitempty"`
	CertificateChain *XmlCertificateChain `xml:"CertificateChain,omitempty"`
	Conclusion       *XmlConclusion       `xml:"Conclusion"`
	Id               string               `xml:"Id,attr"`
	Type             ContextValue         `xml:"Type,attr"`
}
