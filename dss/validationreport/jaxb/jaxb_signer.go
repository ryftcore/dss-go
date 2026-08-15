// Ported from the generated JAXB classes:
//   - SignerInformationType.java
//   - SignatureQualityType.java
//   - SignatureValidationProcessType.java
//   - ValidationTimeInfoType.java
//   - POEType.java
//   - POEProvisioningType.java
//   - CertificateChainType.java
//   - RevocationStatusInformationType.java
//   - CryptoInformationType.java
//   - ValidationReportDataType.java
//
// (specs-validation-report, DSS 6.5.RC1).
package jaxb

// SignerInformationType is the Go form of the generated JAXB class
// SignerInformationType (complexType SignerInformationType).
type SignerInformationType struct {
	SignerCertificate VOReferenceType `xml:"SignerCertificate"`
	Signer            *string         `xml:"Signer,omitempty"`
	OtherInformation  *AnyType        `xml:"OtherInformation,omitempty"`
	Pseudonym         *bool           `xml:"Pseudonym,attr,omitempty"`
}

// SignatureQualityType is the Go form of the generated JAXB class
// SignatureQualityType (complexType SignatureQualityType).
type SignatureQualityType struct {
	SignatureQualityInformation []string `xml:"SignatureQualityInformation"`
}

// SignatureValidationProcessType is the Go form of the generated JAXB class
// SignatureValidationProcessType (complexType SignatureValidationProcessType).
type SignatureValidationProcessType struct {
	SignatureValidationProcessID         *SignatureValidationProcessID `xml:"SignatureValidationProcessID,omitempty"`
	SignatureValidationServicePolicy     *string                       `xml:"SignatureValidationServicePolicy,omitempty"`
	SignatureValidationPracticeStatement *string                       `xml:"SignatureValidationPracticeStatement,omitempty"`
	OtherInformation                     *AnyType                      `xml:"OtherInformation,omitempty"`
}

// ValidationTimeInfoType is the Go form of the generated JAXB class
// ValidationTimeInfoType (complexType ValidationTimeInfoType).
type ValidationTimeInfoType struct {
	ValidationTime    XSDateTime `xml:"ValidationTime"`
	BestSignatureTime POEType    `xml:"BestSignatureTime"`
}

// POEType is the Go form of the generated JAXB class POEType (complexType
// POEType).
type POEType struct {
	POETime     XSDateTime       `xml:"POETime"`
	TypeOfProof TypeOfProof      `xml:"TypeOfProof"`
	POEObject   *VOReferenceType `xml:"POEObject,omitempty"`
}

// POEProvisioningType is the Go form of the generated JAXB class
// POEProvisioningType (complexType POEProvisioningType).
type POEProvisioningType struct {
	POETime            XSDateTime                `xml:"POETime"`
	ValidationObject   []*VOReferenceType        `xml:"ValidationObject,omitempty"`
	SignatureReference []*SignatureReferenceType `xml:"SignatureReference,omitempty"`
}

// CertificateChainType is the Go form of the generated JAXB class
// CertificateChainType (complexType CertificateChainType).
type CertificateChainType struct {
	SigningCertificate      VOReferenceType    `xml:"SigningCertificate"`
	IntermediateCertificate []*VOReferenceType `xml:"IntermediateCertificate,omitempty"`
	TrustAnchor             *VOReferenceType   `xml:"TrustAnchor,omitempty"`
	OtherInformation        *AnyType           `xml:"OtherInformation,omitempty"`
}

// RevocationStatusInformationType is the Go form of the generated JAXB
// class RevocationStatusInformationType (complexType
// RevocationStatusInformationType).
type RevocationStatusInformationType struct {
	ValidationObjectId VOReferenceType      `xml:"ValidationObjectId"`
	RevocationTime     XSDateTime           `xml:"RevocationTime"`
	RevocationReason   *URIRevocationReason `xml:"RevocationReason,omitempty"`
	RevocationObject   *VOReferenceType     `xml:"RevocationObject,omitempty"`
	OtherInformation   *AnyType             `xml:"OtherInformation,omitempty"`
}

// CryptoInformationType is the Go form of the generated JAXB class
// CryptoInformationType (complexType CryptoInformationType).
type CryptoInformationType struct {
	ValidationObjectId  VOReferenceType `xml:"ValidationObjectId"`
	Algorithm           string          `xml:"Algorithm"`
	AlgorithmParameters *TypedDataType  `xml:"AlgorithmParameters,omitempty"`
	SecureAlgorithm     bool            `xml:"SecureAlgorithm"`
	NotAfter            *XSDateTime     `xml:"NotAfter,omitempty"`
	OtherInformation    *AnyType        `xml:"OtherInformation,omitempty"`
}

// ValidationReportDataType is the Go form of the generated JAXB class
// ValidationReportDataType (complexType ValidationReportDataType).
type ValidationReportDataType struct {
	TrustAnchor                    *VOReferenceType                    `xml:"TrustAnchor,omitempty"`
	CertificateChain               *CertificateChainType               `xml:"CertificateChain,omitempty"`
	RelatedValidationObject        []*VOReferenceType                  `xml:"RelatedValidationObject,omitempty"`
	RevocationStatusInformation    *RevocationStatusInformationType    `xml:"RevocationStatusInformation,omitempty"`
	CryptoInformation              *CryptoInformationType              `xml:"CryptoInformation,omitempty"`
	AdditionalValidationReportData *AdditionalValidationReportDataType `xml:"AdditionalValidationReportData,omitempty"`
}
