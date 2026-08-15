// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule.
package jaxb

// SignatureConstraints is the Go form of the generated JAXB class
// SignatureConstraints (complexType SignatureConstraints).
type SignatureConstraints struct {
	StructuralValidation        *LevelConstraint               `xml:"StructuralValidation,omitempty"`
	AcceptablePolicies          *MultiValuesConstraint         `xml:"AcceptablePolicies,omitempty"`
	PolicyAvailable             *LevelConstraint               `xml:"PolicyAvailable,omitempty"`
	SignaturePolicyStorePresent *LevelConstraint               `xml:"SignaturePolicyStorePresent,omitempty"`
	PolicyHashMatch             *LevelConstraint               `xml:"PolicyHashMatch,omitempty"`
	AcceptableFormats           *MultiValuesConstraint         `xml:"AcceptableFormats,omitempty"`
	FullScope                   *LevelConstraint               `xml:"FullScope,omitempty"`
	BasicSignatureConstraints   *BasicSignatureConstraints     `xml:"BasicSignatureConstraints,omitempty"`
	SignedAttributes            *SignedAttributesConstraints   `xml:"SignedAttributes,omitempty"`
	UnsignedAttributes          *UnsignedAttributesConstraints `xml:"UnsignedAttributes,omitempty"`
}

// SignedAttributesConstraints is the Go form of the generated JAXB class
// SignedAttributesConstraints (complexType SignedAttributesConstraints).
type SignedAttributesConstraints struct {
	SigningCertificatePresent                *LevelConstraint       `xml:"SigningCertificatePresent,omitempty"`
	UnicitySigningCertificate                *LevelConstraint       `xml:"UnicitySigningCertificate,omitempty"`
	SigningCertificateRefersCertificateChain *LevelConstraint       `xml:"SigningCertificateRefersCertificateChain,omitempty"`
	ReferencesToAllCertificateChainPresent   *LevelConstraint       `xml:"ReferencesToAllCertificateChainPresent,omitempty"`
	SigningCertificateDigestAlgorithm        *LevelConstraint       `xml:"SigningCertificateDigestAlgorithm,omitempty"`
	CertDigestPresent                        *LevelConstraint       `xml:"CertDigestPresent,omitempty"`
	CertDigestMatch                          *LevelConstraint       `xml:"CertDigestMatch,omitempty"`
	IssuerSerialMatch                        *LevelConstraint       `xml:"IssuerSerialMatch,omitempty"`
	KeyIdentifierPresent                     *LevelConstraint       `xml:"KeyIdentifierPresent,omitempty"`
	KeyIdentifierMatch                       *LevelConstraint       `xml:"KeyIdentifierMatch,omitempty"`
	X509UrlPresent                           *LevelConstraint       `xml:"X509UrlPresent,omitempty"`
	X509UrlMatch                             *LevelConstraint       `xml:"X509UrlMatch,omitempty"`
	SigningTime                              *LevelConstraint       `xml:"SigningTime,omitempty"`
	SigningTimeInCertRange                   *LevelConstraint       `xml:"SigningTimeInCertRange,omitempty"`
	SignatureType                            *MultiValuesConstraint `xml:"SignatureType,omitempty"`
	ContentType                              *MultiValuesConstraint `xml:"ContentType,omitempty"`
	ContentHints                             *MultiValuesConstraint `xml:"ContentHints,omitempty"`
	ContentIdentifier                        *MultiValuesConstraint `xml:"ContentIdentifier,omitempty"`
	MessageDigestOrSignedPropertiesPresent   *LevelConstraint       `xml:"MessageDigestOrSignedPropertiesPresent,omitempty"`
	EllipticCurveKeySize                     *LevelConstraint       `xml:"EllipticCurveKeySize,omitempty"`
	CommitmentTypeIndication                 *MultiValuesConstraint `xml:"CommitmentTypeIndication,omitempty"`
	SignerLocation                           *LevelConstraint       `xml:"SignerLocation,omitempty"`
	ClaimedRoles                             *MultiValuesConstraint `xml:"ClaimedRoles,omitempty"`
	CertifiedRoles                           *MultiValuesConstraint `xml:"CertifiedRoles,omitempty"`
	ContentTimeStamp                         *LevelConstraint       `xml:"ContentTimeStamp,omitempty"`
	ContentTimeStampMessageImprint           *LevelConstraint       `xml:"ContentTimeStampMessageImprint,omitempty"`
}

// UnsignedAttributesConstraints is the Go form of the generated JAXB class
// UnsignedAttributesConstraints (complexType UnsignedAttributesConstraints).
type UnsignedAttributesConstraints struct {
	CounterSignature                *LevelConstraint `xml:"CounterSignature,omitempty"`
	SignatureTimeStamp              *LevelConstraint `xml:"SignatureTimeStamp,omitempty"`
	ValidationDataTimeStamp         *LevelConstraint `xml:"ValidationDataTimeStamp,omitempty"`
	ValidationDataRefsOnlyTimeStamp *LevelConstraint `xml:"ValidationDataRefsOnlyTimeStamp,omitempty"`
	ArchiveTimeStamp                *LevelConstraint `xml:"ArchiveTimeStamp,omitempty"`
	DocumentTimeStamp               *LevelConstraint `xml:"DocumentTimeStamp,omitempty"`
	TLevelTimeStamp                 *LevelConstraint `xml:"TLevelTimeStamp,omitempty"`
	LTALevelTimeStamp               *LevelConstraint `xml:"LTALevelTimeStamp,omitempty"`
}

// BasicSignatureConstraints is the Go form of the generated JAXB class
// BasicSignatureConstraints (complexType BasicSignatureConstraints): group of
// common checks for any kind of signed token (signature, timestamp or
// revocation data).
type BasicSignatureConstraints struct {
	ReferenceDataExistence       *LevelConstraint         `xml:"ReferenceDataExistence,omitempty"`
	ReferenceDataIntact          *LevelConstraint         `xml:"ReferenceDataIntact,omitempty"`
	ReferenceDataNameMatch       *LevelConstraint         `xml:"ReferenceDataNameMatch,omitempty"`
	ManifestEntryObjectExistence *LevelConstraint         `xml:"ManifestEntryObjectExistence,omitempty"`
	ManifestEntryObjectGroup     *LevelConstraint         `xml:"ManifestEntryObjectGroup,omitempty"`
	ManifestEntryObjectIntact    *LevelConstraint         `xml:"ManifestEntryObjectIntact,omitempty"`
	ManifestEntryNameMatch       *LevelConstraint         `xml:"ManifestEntryNameMatch,omitempty"`
	SignatureIntact              *LevelConstraint         `xml:"SignatureIntact,omitempty"`
	SignatureValid               *LevelConstraint         `xml:"SignatureValid,omitempty"`
	SignatureDuplicated          *LevelConstraint         `xml:"SignatureDuplicated,omitempty"`
	ProspectiveCertificateChain  *LevelConstraint         `xml:"ProspectiveCertificateChain,omitempty"`
	SignerInformationStore       *LevelConstraint         `xml:"SignerInformationStore,omitempty"`
	ByteRange                    *LevelConstraint         `xml:"ByteRange,omitempty"`
	ByteRangeCollision           *LevelConstraint         `xml:"ByteRangeCollision,omitempty"`
	ByteRangeAllDocument         *LevelConstraint         `xml:"ByteRangeAllDocument,omitempty"`
	PdfSignatureDictionary       *LevelConstraint         `xml:"PdfSignatureDictionary,omitempty"`
	PdfPageDifference            *LevelConstraint         `xml:"PdfPageDifference,omitempty"`
	PdfAnnotationOverlap         *LevelConstraint         `xml:"PdfAnnotationOverlap,omitempty"`
	PdfVisualDifference          *LevelConstraint         `xml:"PdfVisualDifference,omitempty"`
	DocMDP                       *LevelConstraint         `xml:"DocMDP,omitempty"`
	FieldMDP                     *LevelConstraint         `xml:"FieldMDP,omitempty"`
	SigFieldLock                 *LevelConstraint         `xml:"SigFieldLock,omitempty"`
	FormFillChanges              *LevelConstraint         `xml:"FormFillChanges,omitempty"`
	AnnotationChanges            *LevelConstraint         `xml:"AnnotationChanges,omitempty"`
	UndefinedChanges             *LevelConstraint         `xml:"UndefinedChanges,omitempty"`
	TrustServiceTypeIdentifier   *MultiValuesConstraint   `xml:"TrustServiceTypeIdentifier,omitempty"`
	TrustServiceStatus           *MultiValuesConstraint   `xml:"TrustServiceStatus,omitempty"`
	SigningCertificate           *CertificateConstraints  `xml:"SigningCertificate,omitempty"`
	CACertificate                *CertificateConstraints  `xml:"CACertificate,omitempty"`
	Cryptographic                *CryptographicConstraint `xml:"Cryptographic,omitempty"`
}
