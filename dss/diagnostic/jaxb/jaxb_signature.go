// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. The generated classes are grouped into
// schema-area files rather than one file per class; every Java class keeps its
// name and its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte.

package jaxb

// XmlCOSESignatureType is the Go form of the generated JAXB class XmlCOSESignatureType
// (complexType COSESignatureType).
type XmlCOSESignatureType struct {
	Value  COSESignatureTypeValue `xml:",chardata"`
	Tagged *bool                  `xml:"tagged,attr,omitempty"`
}

// XmlCommitmentTypeIndication is the Go form of the generated JAXB class XmlCommitmentTypeIndication
// (complexType CommitmentTypeIndication).
type XmlCommitmentTypeIndication struct {
	Identifier              *string                         `xml:"Identifier,omitempty"`
	Description             *string                         `xml:"Description,omitempty"`
	DocumentationReferences *DocumentationReferencesWrapper `xml:"DocumentationReferences"`
	ObjectReferences        *ObjectReferencesWrapper        `xml:"ObjectReferences"`
	AllDataSignedObjects    *bool                           `xml:"AllDataSignedObjects,omitempty"`
}

// XmlPolicy is the Go form of the generated JAXB class XmlPolicy
// (complexType Policy).
type XmlPolicy struct {
	Id                      *string                         `xml:"Id,omitempty"`
	Url                     *string                         `xml:"Url,omitempty"`
	UserNotice              *XmlUserNotice                  `xml:"UserNotice,omitempty"`
	DocSpecification        *XmlSPDocSpecification          `xml:"DocSpecification,omitempty"`
	Description             *string                         `xml:"Description,omitempty"`
	Identified              *bool                           `xml:"Identified,omitempty"`
	Asn1Processable         *bool                           `xml:"Asn1Processable,omitempty"`
	Transformations         *TransformationsWrapper         `xml:"Transformations"`
	DocumentationReferences *DocumentationReferencesWrapper `xml:"DocumentationReferences"`
	ProcessingError         *string                         `xml:"ProcessingError,omitempty"`
	DigestAlgoAndValue      *XmlPolicyDigestAlgoAndValue    `xml:"DigestAlgoAndValue,omitempty"`
}

// XmlSPDocSpecificationContent carries the element content of the XmlSPDocSpecification base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlSPDocSpecificationContent struct {
	Id                      *string                         `xml:"Id,omitempty"`
	Description             *string                         `xml:"Description,omitempty"`
	DocumentationReferences *DocumentationReferencesWrapper `xml:"DocumentationReferences"`
}

// XmlSPDocSpecification is the Go form of the generated JAXB class XmlSPDocSpecification
// (complexType SPDocSpecification).
type XmlSPDocSpecification struct {
	XmlSPDocSpecificationContent
}

// XmlSignatureDigestReference is the Go form of the generated JAXB class XmlSignatureDigestReference
// (complexType DigestReference).
type XmlSignatureDigestReference struct {
	CanonicalizationMethod *string               `xml:"CanonicalizationMethod,omitempty"`
	DigestMethod           *DigestAlgorithmValue `xml:"DigestMethod,omitempty"`
	DigestValue            *Base64Binary         `xml:"DigestValue,omitempty"`
}

// XmlSignatureProductionPlace is the Go form of the generated JAXB class XmlSignatureProductionPlace
// (complexType ProductionPlace).
type XmlSignatureProductionPlace struct {
	PostalAddress       []string `xml:"PostalAddress"`
	City                *string  `xml:"City,omitempty"`
	StateOrProvince     *string  `xml:"StateOrProvince,omitempty"`
	PostOfficeBoxNumber *string  `xml:"PostOfficeBoxNumber,omitempty"`
	PostalCode          *string  `xml:"PostalCode,omitempty"`
	CountryName         *string  `xml:"CountryName,omitempty"`
	StreetAddress       *string  `xml:"StreetAddress,omitempty"`
}

// XmlSignatureScope is the Go form of the generated JAXB class XmlSignatureScope
// (complexType SignatureScope).
type XmlSignatureScope struct {
	Scope           *SignatureScopeTypeValue `xml:"Scope,omitempty"`
	Name            *string                  `xml:"Name,omitempty"`
	Description     *string                  `xml:"Description,omitempty"`
	Transformations *TransformationsWrapper  `xml:"Transformations"`
	SignerData      *XmlSignerData           `xml:"SignerData,attr,omitempty"`
}

// XmlSignerDocumentRepresentations is the Go form of the generated JAXB class XmlSignerDocumentRepresentations
// (complexType anonymous).
type XmlSignerDocumentRepresentations struct {
	HashOnly    bool `xml:"HashOnly,attr"`
	DocHashOnly bool `xml:"DocHashOnly,attr"`
}

// XmlSignerRole is the Go form of the generated JAXB class XmlSignerRole
// (complexType SignerRole).
type XmlSignerRole struct {
	Role      *string               `xml:"Role,omitempty"`
	NotAfter  *XSDateTime           `xml:"NotAfter,omitempty"`
	NotBefore *XSDateTime           `xml:"NotBefore,omitempty"`
	Category  *EndorsementTypeValue `xml:"Category,attr,omitempty"`
}

// XmlUserNotice is the Go form of the generated JAXB class XmlUserNotice
// (complexType UserNotice).
type XmlUserNotice struct {
	Organization  *string         `xml:"Organization,omitempty"`
	NoticeNumbers *BigIntegerList `xml:"NoticeNumbers"`
	ExplicitText  *string         `xml:"ExplicitText,omitempty"`
}

// XmlPolicyDigestAlgoAndValue is the Go form of the generated JAXB class XmlPolicyDigestAlgoAndValue
// (complexType PolicyDigestAlgoAndValue).
type XmlPolicyDigestAlgoAndValue struct {
	XmlDigestAlgoAndValueContent
	DigestAlgorithmsEqual *bool `xml:"digestAlgorithmsEqual,attr,omitempty"`
	ZeroHash              *bool `xml:"zeroHash,attr,omitempty"`
	XmlDigestAlgoAndValueAttrs
}

// XmlSignature is the Go form of the generated JAXB class XmlSignature
// (complexType Signature).
type XmlSignature struct {
	DAIdentifier                  *string                           `xml:"DAIdentifier,omitempty"`
	SignatureFilename             *string                           `xml:"SignatureFilename,omitempty"`
	ErrorMessage                  *string                           `xml:"ErrorMessage,omitempty"`
	ClaimedSigningTime            *XSDateTime                       `xml:"ClaimedSigningTime,omitempty"`
	ExpirationTime                *XSDateTime                       `xml:"ExpirationTime,omitempty"`
	SignatureFormat               *SignatureLevelValue              `xml:"SignatureFormat,omitempty"`
	JWSSerializationType          *JWSSerializationTypeValue        `xml:"JWSSerializationType,omitempty"`
	SignatureType                 *string                           `xml:"SignatureType,omitempty"`
	COSESignatureType             *XmlCOSESignatureType             `xml:"COSESignatureType,omitempty"`
	StructuralValidation          *XmlStructuralValidation          `xml:"StructuralValidation,omitempty"`
	DigestMatchers                *DigestMatchersWrapper            `xml:"DigestMatchers"`
	BasicSignature                *XmlBasicSignature                `xml:"BasicSignature,omitempty"`
	SigningCertificate            *XmlSigningCertificate            `xml:"SigningCertificate,omitempty"`
	CertificateChain              *CertificateChainWrapper          `xml:"CertificateChain"`
	ContentType                   *string                           `xml:"ContentType,omitempty"`
	MimeType                      *string                           `xml:"MimeType,omitempty"`
	ContentIdentifier             *string                           `xml:"ContentIdentifier,omitempty"`
	ContentHints                  *string                           `xml:"ContentHints,omitempty"`
	SignatureProductionPlace      *XmlSignatureProductionPlace      `xml:"SignatureProductionPlace,omitempty"`
	CommitmentTypeIndications     *CommitmentTypeIndicationsWrapper `xml:"CommitmentTypeIndications"`
	SignerRole                    []*XmlSignerRole                  `xml:"SignerRole"`
	Policy                        *XmlPolicy                        `xml:"Policy,omitempty"`
	SignaturePolicyStore          *XmlSignaturePolicyStore          `xml:"SignaturePolicyStore,omitempty"`
	SignerInformationStore        *SignerInformationStoreWrapper    `xml:"SignerInformationStore"`
	PDFRevision                   *XmlPDFRevision                   `xml:"PDFRevision,omitempty"`
	VRIDictionaryCreationTime     *XSDateTime                       `xml:"VRIDictionaryCreationTime,omitempty"`
	SignerDocumentRepresentations *XmlSignerDocumentRepresentations `xml:"SignerDocumentRepresentations,omitempty"`
	FoundCertificates             *XmlFoundCertificates             `xml:"FoundCertificates,omitempty"`
	FoundRevocations              *XmlFoundRevocations              `xml:"FoundRevocations,omitempty"`
	FoundTimestamps               *FoundTimestampsWrapper           `xml:"FoundTimestamps"`
	FoundEvidenceRecords          *FoundEvidenceRecordsWrapper      `xml:"FoundEvidenceRecords"`
	SignatureScopes               *SignatureScopesWrapper           `xml:"SignatureScopes"`
	SignatureDigestReference      *XmlSignatureDigestReference      `xml:"SignatureDigestReference,omitempty"`
	DataToBeSignedRepresentation  *XmlDigestAlgoAndValue            `xml:"DataToBeSignedRepresentation,omitempty"`
	SignatureValue                *Base64Binary                     `xml:"SignatureValue,omitempty"`
	CounterSignature              *bool                             `xml:"CounterSignature,attr,omitempty"`
	KeyBindingSignature           *bool                             `xml:"KeyBindingSignature,attr,omitempty"`
	Parent                        *XmlSignature                     `xml:"Parent,attr,omitempty"`
	XmlAbstractTokenAttrs
}

// XmlSignaturePolicyStore is the Go form of the generated JAXB class XmlSignaturePolicyStore
// (complexType SignaturePolicyStore).
type XmlSignaturePolicyStore struct {
	XmlSPDocSpecificationContent
	SigPolDocLocalURI  *string                `xml:"SigPolDocLocalURI,omitempty"`
	DigestAlgoAndValue *XmlDigestAlgoAndValue `xml:"DigestAlgoAndValue,omitempty"`
}
