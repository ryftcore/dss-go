// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/XAdESPath.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// XAdESPath contains a list of useful XAdES XPaths, one implementation per XAdES schema
// version (XAdES111Path, XAdES122Path, XAdES132Path). A path a given schema version does
// not define returns nil, mirroring Java's null-returning override.
type XAdESPath interface {
	// Namespace gets the current namespace.
	Namespace() *common.DSSNamespace

	// SignedPropertiesUri gets signed properties reference URI.
	SignedPropertiesUri() string

	// CounterSignatureUri gets counter signature reference URI.
	CounterSignatureUri() string

	// QualifyingPropertiesPath gets path "./ds:Object/xades:QualifyingProperties"
	QualifyingPropertiesPath() common.XPathQuery

	// SignedPropertiesPath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties"
	SignedPropertiesPath() common.XPathQuery

	// SignedSignaturePropertiesPath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties"
	SignedSignaturePropertiesPath() common.XPathQuery

	// SigningTimePath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SigningTime"
	SigningTimePath() common.XPathQuery

	// SigningCertificatePath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SigningCertificate"
	SigningCertificatePath() common.XPathQuery

	// SigningCertificateChildren gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SigningCertificate/xades:Cert"
	SigningCertificateChildren() common.XPathQuery

	// SigningCertificateV2Path gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SigningCertificateV2"
	SigningCertificateV2Path() common.XPathQuery

	// SigningCertificateV2Children gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SigningCertificateV2/xades:Cert"
	SigningCertificateV2Children() common.XPathQuery

	// SignatureProductionPlacePath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignatureProductionPlace"
	SignatureProductionPlacePath() common.XPathQuery

	// SignatureProductionPlaceV2Path gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignatureProductionPlaceV2"
	SignatureProductionPlaceV2Path() common.XPathQuery

	// SignaturePolicyIdentifierPath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignaturePolicyIdentifier"
	SignaturePolicyIdentifierPath() common.XPathQuery

	// SignerRolePath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignerRole"
	SignerRolePath() common.XPathQuery

	// ClaimedRolePath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignerRole/xades:ClaimedRoles/xades:ClaimedRole"
	ClaimedRolePath() common.XPathQuery

	// SignedAssertionPath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignerRole/xades:SignedAssertions/xades:SignedAssertion"
	SignedAssertionPath() common.XPathQuery

	// SignerRoleV2Path gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignerRoleV2"
	SignerRoleV2Path() common.XPathQuery

	// ClaimedRoleV2Path gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignerRoleV2/xades:ClaimedRoles/xades:ClaimedRole"
	ClaimedRoleV2Path() common.XPathQuery

	// CertifiedRolePath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignerRole/xades:CertifiedRoles/xades:CertifiedRole"
	CertifiedRolePath() common.XPathQuery

	// CertifiedRoleV2Path gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedSignatureProperties/xades:SignerRoleV2/xades:CertifiedRoles/xades:CertifiedRole"
	CertifiedRoleV2Path() common.XPathQuery

	// SignedDataObjectPropertiesPath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedDataObjectProperties"
	SignedDataObjectPropertiesPath() common.XPathQuery

	// AllDataObjectsTimestampPath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedDataObjectProperties/xades:AllDataObjectsTimeStamp"
	AllDataObjectsTimestampPath() common.XPathQuery

	// IndividualDataObjectsTimestampPath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedDataObjectProperties/xades:IndividualDataObjectsTimeStamp"
	IndividualDataObjectsTimestampPath() common.XPathQuery

	// DataObjectFormat gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedDataObjectProperties/xades:DataObjectFormat"
	DataObjectFormat() common.XPathQuery

	// DataObjectFormatMimeType gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedDataObjectProperties/xades:DataObjectFormat/xades:MimeType"
	DataObjectFormatMimeType() common.XPathQuery

	// DataObjectFormatObjectIdentifier gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedDataObjectProperties/xades:DataObjectFormat/xades:ObjectIdentifier"
	DataObjectFormatObjectIdentifier() common.XPathQuery

	// CommitmentTypeIndicationPath gets path "./ds:Object/xades:QualifyingProperties/xades:SignedProperties/xades:SignedDataObjectProperties/xades:CommitmentTypeIndication"
	CommitmentTypeIndicationPath() common.XPathQuery

	// UnsignedPropertiesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties"
	UnsignedPropertiesPath() common.XPathQuery

	// UnsignedSignaturePropertiesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties"
	UnsignedSignaturePropertiesPath() common.XPathQuery

	// CounterSignaturePath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:CounterSignature"
	CounterSignaturePath() common.XPathQuery

	// AttributeRevocationRefsPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:AttributeRevocationRefs"
	AttributeRevocationRefsPath() common.XPathQuery

	// CompleteRevocationRefsPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:CompleteRevocationRefs"
	CompleteRevocationRefsPath() common.XPathQuery

	// CompleteCertificateRefsPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:CompleteCertificateRefs"
	CompleteCertificateRefsPath() common.XPathQuery

	// CompleteCertificateRefsCertPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:CompleteCertificateRefs/xades:CertRefs/xades:Cert"
	CompleteCertificateRefsCertPath() common.XPathQuery

	// CompleteCertificateRefsV2Path gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:CompleteCertificateRefsV2"
	CompleteCertificateRefsV2Path() common.XPathQuery

	// CompleteCertificateRefsV2CertPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:CompleteCertificateRefsV2/xades141:CertRefs/xades:Cert"
	CompleteCertificateRefsV2CertPath() common.XPathQuery

	// AttributeCertificateRefsPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:AttributeCertificateRefs"
	AttributeCertificateRefsPath() common.XPathQuery

	// AttributeCertificateRefsCertPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:AttributeCertificateRefs/xades:CertRefs/xades:Cert"
	AttributeCertificateRefsCertPath() common.XPathQuery

	// AttributeCertificateRefsV2Path gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:AttributeCertificateRefsV2"
	AttributeCertificateRefsV2Path() common.XPathQuery

	// AttributeCertificateRefsV2CertPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:AttributeCertificateRefsV2/xades141:CertRefs/xades:Cert"
	AttributeCertificateRefsV2CertPath() common.XPathQuery

	// CertificateValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:CertificateValues"
	CertificateValuesPath() common.XPathQuery

	// RevocationValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:RevocationValues"
	RevocationValuesPath() common.XPathQuery

	// AttributeRevocationValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:AttributeRevocationValues"
	AttributeRevocationValuesPath() common.XPathQuery

	// EncapsulatedCertificateValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:SigAndRefsTimeStampV2"
	EncapsulatedCertificateValuesPath() common.XPathQuery

	// AttrAuthoritiesCertValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:AttrAuthoritiesCertValues"
	AttrAuthoritiesCertValuesPath() common.XPathQuery

	// EncapsulatedAttrAuthoritiesCertValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:AttrAuthoritiesCertValues/xades:EncapsulatedX509Certificate"
	EncapsulatedAttrAuthoritiesCertValuesPath() common.XPathQuery

	// EncapsulatedTimeStampValidationDataCertValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:TimeStampValidationData/xades:CertificateValues/xades:EncapsulatedX509Certificate"
	EncapsulatedTimeStampValidationDataCertValuesPath() common.XPathQuery

	// TimeStampValidationDataRevocationValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:TimeStampValidationData/xades:RevocationValues"
	TimeStampValidationDataRevocationValuesPath() common.XPathQuery

	// AnyValidationDataPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:AnyValidationData"
	AnyValidationDataPath() common.XPathQuery

	// EncapsulatedAnyValidationDataCertValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:AnyValidationData/xades:CertificateValues/xades:EncapsulatedX509Certificate"
	EncapsulatedAnyValidationDataCertValuesPath() common.XPathQuery

	// AnyValidationDataRevocationValuesPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:AnyValidationData/xades:RevocationValues"
	AnyValidationDataRevocationValuesPath() common.XPathQuery

	// SignatureTimestampPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:SignatureTimeStamp"
	SignatureTimestampPath() common.XPathQuery

	// SigAndRefsTimestampPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:SigAndRefsTimeStamp"
	SigAndRefsTimestampPath() common.XPathQuery

	// SigAndRefsTimestampV2Path gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:SigAndRefsTimeStampV2"
	SigAndRefsTimestampV2Path() common.XPathQuery

	// RefsOnlyTimestampPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades:RefsOnlyTimeStamp"
	RefsOnlyTimestampPath() common.XPathQuery

	// RefsOnlyTimestampV2Path gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:RefsOnlyTimeStampV2"
	RefsOnlyTimestampV2Path() common.XPathQuery

	// ArchiveTimestampPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:ArchiveTimeStamp"
	ArchiveTimestampPath() common.XPathQuery

	// TimestampValidationDataPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:TimeStampValidationData"
	TimestampValidationDataPath() common.XPathQuery

	// SignaturePolicyStorePath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xades141:SignaturePolicyStore"
	SignaturePolicyStorePath() common.XPathQuery

	// SealingEvidenceRecordsPath gets path "./ds:Object/xades:QualifyingProperties/xades:UnsignedProperties/xades:UnsignedSignatureProperties/xadesen:SealingEvidenceRecords"
	SealingEvidenceRecordsPath() common.XPathQuery

	// CurrentCRLValuesChildren gets path "./xades:CRLValues/xades:EncapsulatedCRLValue"
	CurrentCRLValuesChildren() common.XPathQuery

	// CurrentCRLRefsChildren gets path "./xades:CRLRefs/xades:CRLRef"
	CurrentCRLRefsChildren() common.XPathQuery

	// CurrentCRLRefCRLIdentifier gets path "./xades:CRLIdentifier"
	CurrentCRLRefCRLIdentifier() common.XPathQuery

	// CurrentCRLRefCRLIdentifierIssuer gets path "./xades:CRLIdentifier/xades:Issuer"
	CurrentCRLRefCRLIdentifierIssuer() common.XPathQuery

	// CurrentCRLRefCRLIdentifierIssueTime gets path "./xades:CRLIdentifier/xades:IssueTime"
	CurrentCRLRefCRLIdentifierIssueTime() common.XPathQuery

	// CurrentCRLRefCRLIdentifierNumber gets path "./xades:CRLIdentifier/xades:Number"
	CurrentCRLRefCRLIdentifierNumber() common.XPathQuery

	// CurrentOCSPValuesChildren gets path "./xades:OCSPValues/xades:EncapsulatedOCSPValue"
	CurrentOCSPValuesChildren() common.XPathQuery

	// CurrentOCSPRefsChildren gets path "./xades:OCSPRefs/xades:OCSPRef"
	CurrentOCSPRefsChildren() common.XPathQuery

	// CurrentOCSPRefResponderID gets path "./xades:OCSPIdentifier/xades:ResponderID"
	CurrentOCSPRefResponderID() common.XPathQuery

	// CurrentOCSPRefResponderIDByName gets path "./xades:OCSPIdentifier/xades:ResponderID/xades:ByName"
	CurrentOCSPRefResponderIDByName() common.XPathQuery

	// CurrentOCSPRefResponderIDByKey gets path "./xades:OCSPIdentifier/xades:ResponderID/xades:ByKey"
	CurrentOCSPRefResponderIDByKey() common.XPathQuery

	// CurrentOCSPRefProducedAt gets path "./xades:OCSPIdentifier/xades:ProducedAt"
	CurrentOCSPRefProducedAt() common.XPathQuery

	// CurrentDigestAlgAndValue gets path "./xades:DigestAlgAndValue"
	CurrentDigestAlgAndValue() common.XPathQuery

	// CurrentCertRefsCertChildren gets path "./xades:CertRefs/xades:Cert"
	CurrentCertRefsCertChildren() common.XPathQuery

	// CurrentCertRefs141CertChildren gets path "./xades141:CertRefs/xades:Cert"
	CurrentCertRefs141CertChildren() common.XPathQuery

	// CurrentCertChildren gets path "./xades:Cert"
	CurrentCertChildren() common.XPathQuery

	// CurrentCertDigest gets path "./xades:CertDigest"
	CurrentCertDigest() common.XPathQuery

	// CurrentEncapsulatedTimestamp gets path "./xades:EncapsulatedTimeStamp"
	CurrentEncapsulatedTimestamp() common.XPathQuery

	// CurrentEncapsulatedCertificate gets path "./xades:EncapsulatedX509Certificate"
	CurrentEncapsulatedCertificate() common.XPathQuery

	// CurrentCertificateValuesEncapsulatedCertificate gets path "./xades:CertificateValues/xades:EncapsulatedX509Certificate"
	CurrentCertificateValuesEncapsulatedCertificate() common.XPathQuery

	// CurrentRevocationValuesEncapsulatedOCSPValue gets path "./xades:RevocationValues/xades:OCSPValues/xades:EncapsulatedOCSPValue"
	CurrentRevocationValuesEncapsulatedOCSPValue() common.XPathQuery

	// CurrentEncapsulatedOCSPValue gets path "./xades:OCSPValues/xades:EncapsulatedOCSPValue"
	CurrentEncapsulatedOCSPValue() common.XPathQuery

	// CurrentRevocationValuesEncapsulatedCRLValue gets path "./xades:RevocationValues/xades:CRLValues/xades:EncapsulatedCRLValue"
	CurrentRevocationValuesEncapsulatedCRLValue() common.XPathQuery

	// CurrentEncapsulatedCRLValue gets path "./xades:CRLValues/xades:EncapsulatedCRLValue"
	CurrentEncapsulatedCRLValue() common.XPathQuery

	// CurrentIssuerSerialIssuerNamePath gets path "./xades:IssuerSerial/xades:X509IssuerName"
	CurrentIssuerSerialIssuerNamePath() common.XPathQuery

	// CurrentIssuerSerialSerialNumberPath gets path "./xades:IssuerSerial/xades:X509SerialNumber"
	CurrentIssuerSerialSerialNumberPath() common.XPathQuery

	// CurrentIssuerSerialV2Path gets path "./xades:IssuerSerialV2"
	CurrentIssuerSerialV2Path() common.XPathQuery

	// CurrentCommitmentIdentifierPath gets path "./xades:CommitmentTypeId/xades:Identifier"
	CurrentCommitmentIdentifierPath() common.XPathQuery

	// CurrentCommitmentDescriptionPath gets path "./xades:CommitmentTypeId/xades:Description"
	CurrentCommitmentDescriptionPath() common.XPathQuery

	// CurrentCommitmentDocumentationReferencesPath gets path "./xades:CommitmentTypeId/xades:DocumentationReferences"
	CurrentCommitmentDocumentationReferencesPath() common.XPathQuery

	// CurrentDocumentationReference gets path "./xades:DocumentationReferences"
	CurrentDocumentationReference() common.XPathQuery

	// CurrentDescription gets path "./xades:Description"
	CurrentDescription() common.XPathQuery

	// CurrentObjectIdentifier gets path "./xades:ObjectIdentifier"
	CurrentObjectIdentifier() common.XPathQuery

	// CurrentCommitmentObjectReferencesPath gets path "./xades:ObjectReference"
	CurrentCommitmentObjectReferencesPath() common.XPathQuery

	// CurrentCommitmentAllSignedDataObjectsPath gets path "./xades:AllSignedDataObjects"
	CurrentCommitmentAllSignedDataObjectsPath() common.XPathQuery

	// CurrentMimeType gets path "./xades:MimeType"
	CurrentMimeType() common.XPathQuery

	// CurrentEncoding gets path "./xades:Encoding"
	CurrentEncoding() common.XPathQuery

	// CurrentSignaturePolicyId gets path "./xades:SignaturePolicyId/xades:SigPolicyId/xades:Identifier"
	CurrentSignaturePolicyId() common.XPathQuery

	// CurrentSignaturePolicyDigestAlgAndValue gets path "./xades:SignaturePolicyId/xades:SigPolicyHash"
	CurrentSignaturePolicyDigestAlgAndValue() common.XPathQuery

	// CurrentSignaturePolicySPURI gets path "./xades:SignaturePolicyId/xades:SigPolicyQualifiers/xades:SigPolicyQualifier/xades:SPURI"
	CurrentSignaturePolicySPURI() common.XPathQuery

	// CurrentSignaturePolicySPUserNotice gets path "./xades:SignaturePolicyId/xades:SigPolicyQualifiers/xades:SigPolicyQualifier/xades:SPUserNotice"
	CurrentSignaturePolicySPUserNotice() common.XPathQuery

	// CurrentSPUserNoticeNoticeRefOrganization gets path "./xades:NoticeRef/xades:Organization"
	CurrentSPUserNoticeNoticeRefOrganization() common.XPathQuery

	// CurrentSPUserNoticeNoticeRefNoticeNumbers gets path "./xades:NoticeRef/xades:NoticeNumbers"
	CurrentSPUserNoticeNoticeRefNoticeNumbers() common.XPathQuery

	// CurrentSPUserNoticeExplicitText gets path "./xades:NoticeRef/xades:NoticeNumbers"
	CurrentSPUserNoticeExplicitText() common.XPathQuery

	// CurrentSignaturePolicySPDocSpecification gets path "./xades:SignaturePolicyId/xades:SigPolicyQualifiers/xades:SigPolicyQualifier/xades141:SPDocSpecification"
	CurrentSignaturePolicySPDocSpecification() common.XPathQuery

	// CurrentSignaturePolicySPDocSpecificationIdentifier gets path "./xades:SignaturePolicyId/xades:SigPolicyQualifiers/xades:SigPolicyQualifier/xades141:SPDocSpecification/xades:Identifier"
	CurrentSignaturePolicySPDocSpecificationIdentifier() common.XPathQuery

	// CurrentSignaturePolicyDescription gets path "./xades:SignaturePolicyId/xades:SigPolicyId/xades:Description"
	CurrentSignaturePolicyDescription() common.XPathQuery

	// CurrentSignaturePolicyDocumentationReferences gets path "./xades:SignaturePolicyId/xades:SigPolicyId/xades:DocumentationReferences"
	CurrentSignaturePolicyDocumentationReferences() common.XPathQuery

	// CurrentSignaturePolicyImplied gets path "./xades:SignaturePolicyImplied"
	CurrentSignaturePolicyImplied() common.XPathQuery

	// CurrentSignaturePolicyTransforms gets path "./xades:SignaturePolicyId/ds:Transforms"
	CurrentSignaturePolicyTransforms() common.XPathQuery

	// CurrentSignaturePolicyQualifiers gets path "./xades:SignaturePolicyId/ds:SigPolicyQualifiers"
	CurrentSignaturePolicyQualifiers() common.XPathQuery

	// CurrentInclude gets path "./xades:Include"
	CurrentInclude() common.XPathQuery

	// CurrentQualifyingPropertiesPath gets path "./xades:QualifyingProperties"
	CurrentQualifyingPropertiesPath() common.XPathQuery

	// CurrentSPDocSpecification gets path "./xades141:SPDocSpecification"
	CurrentSPDocSpecification() common.XPathQuery

	// CurrentIdentifier gets path "./xades:Identifier"
	CurrentIdentifier() common.XPathQuery

	// CurrentSPDocSpecificationIdentifier gets path "./xades141:SPDocSpecification/xades:Identifier"
	CurrentSPDocSpecificationIdentifier() common.XPathQuery

	// CurrentSPDocSpecificationDescription gets path "./xades141:SPDocSpecification/xades:Description"
	CurrentSPDocSpecificationDescription() common.XPathQuery

	// CurrentDocumentationReferenceElements gets path ".xades:DocumentationReferences/xades:DocumentationReference"
	CurrentDocumentationReferenceElements() common.XPathQuery

	// CurrentSPDocSpecificationDocumentationReferenceElements gets path "./xades141:SPDocSpecification/xades:DocumentationReferences/xades:DocumentationReference"
	CurrentSPDocSpecificationDocumentationReferenceElements() common.XPathQuery

	// CurrentSignaturePolicyDocument gets path "./xades141:SignaturePolicyDocument"
	CurrentSignaturePolicyDocument() common.XPathQuery

	// CurrentSigPolDocLocalURI gets path "./xades141:SigPolDocLocalURI"
	CurrentSigPolDocLocalURI() common.XPathQuery
}
