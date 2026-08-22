// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades132/XAdES132Path.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES132Path holds the XAdES 132 paths.
type XAdES132Path struct{}

// NewXAdES132Path creates a XAdES132Path.
func NewXAdES132Path() *XAdES132Path {
	return &XAdES132Path{}
}

// Namespace implements XAdESPath. Ports getNamespace().
func (p *XAdES132Path) Namespace() *common.DSSNamespace {
	return XAdESNamespace_XADES_132
}

// SignedPropertiesUri implements XAdESPath. Ports getSignedPropertiesUri().
func (p *XAdES132Path) SignedPropertiesUri() string {
	return "http://uri.etsi.org/01903#SignedProperties"
}

// CounterSignatureUri implements XAdESPath. Ports getCounterSignatureUri().
func (p *XAdES132Path) CounterSignatureUri() string {
	return "http://uri.etsi.org/01903#CountersignedSignature"
}

// QualifyingPropertiesPath implements XAdESPath. Ports getQualifyingPropertiesPath().
func (p *XAdES132Path) QualifyingPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES)
}

// SignedPropertiesPath implements XAdESPath. Ports getSignedPropertiesPath().
func (p *XAdES132Path) SignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES)
}

// SignedSignaturePropertiesPath implements XAdESPath. Ports getSignedSignaturePropertiesPath().
func (p *XAdES132Path) SignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES)
}

// SigningTimePath implements XAdESPath. Ports getSigningTimePath().
func (p *XAdES132Path) SigningTimePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNING_TIME)
}

// SigningCertificatePath implements XAdESPath. Ports getSigningCertificatePath().
func (p *XAdES132Path) SigningCertificatePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNING_CERTIFICATE)
}

// SigningCertificateChildren implements XAdESPath. Ports getSigningCertificateChildren().
func (p *XAdES132Path) SigningCertificateChildren() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNING_CERTIFICATE, XAdES132Element_CERT)
}

// SigningCertificateV2Path implements XAdESPath. Ports getSigningCertificateV2Path().
func (p *XAdES132Path) SigningCertificateV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNING_CERTIFICATE_V2)
}

// SigningCertificateV2Children implements XAdESPath. Ports getSigningCertificateV2Children().
func (p *XAdES132Path) SigningCertificateV2Children() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNING_CERTIFICATE_V2, XAdES132Element_CERT)
}

// SignatureProductionPlacePath implements XAdESPath. Ports getSignatureProductionPlacePath().
func (p *XAdES132Path) SignatureProductionPlacePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNATURE_PRODUCTION_PLACE)
}

// SignatureProductionPlaceV2Path implements XAdESPath. Ports getSignatureProductionPlaceV2Path().
func (p *XAdES132Path) SignatureProductionPlaceV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNATURE_PRODUCTION_PLACE_V2)
}

// SignaturePolicyIdentifierPath implements XAdESPath. Ports getSignaturePolicyIdentifierPath().
func (p *XAdES132Path) SignaturePolicyIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNATURE_POLICY_IDENTIFIER)
}

// SignerRolePath implements XAdESPath. Ports getSignerRolePath().
func (p *XAdES132Path) SignerRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNER_ROLE)
}

// ClaimedRolePath implements XAdESPath. Ports getClaimedRolePath().
func (p *XAdES132Path) ClaimedRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNER_ROLE, XAdES132Element_CLAIMED_ROLES, XAdES132Element_CLAIMED_ROLE)
}

// SignedAssertionPath implements XAdESPath. Ports getSignedAssertionPath().
func (p *XAdES132Path) SignedAssertionPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNER_ROLE_V2, XAdES132Element_SIGNED_ASSERTIONS, XAdES132Element_SIGNED_ASSERTION)
}

// SignerRoleV2Path implements XAdESPath. Ports getSignerRoleV2Path().
func (p *XAdES132Path) SignerRoleV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNER_ROLE_V2)
}

// ClaimedRoleV2Path implements XAdESPath. Ports getClaimedRoleV2Path().
func (p *XAdES132Path) ClaimedRoleV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNER_ROLE_V2, XAdES132Element_CLAIMED_ROLES, XAdES132Element_CLAIMED_ROLE)
}

// CertifiedRolePath implements XAdESPath. Ports getCertifiedRolePath().
func (p *XAdES132Path) CertifiedRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNER_ROLE, XAdES132Element_CERTIFIED_ROLES, XAdES132Element_CERTIFIED_ROLE)
}

// CertifiedRoleV2Path implements XAdESPath. Ports getCertifiedRoleV2Path().
func (p *XAdES132Path) CertifiedRoleV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNER_ROLE_V2, XAdES132Element_CERTIFIED_ROLES_V2, XAdES132Element_CERTIFIED_ROLE)
}

// SignedDataObjectPropertiesPath implements XAdESPath. Ports getSignedDataObjectPropertiesPath().
func (p *XAdES132Path) SignedDataObjectPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES)
}

// AllDataObjectsTimestampPath implements XAdESPath. Ports getAllDataObjectsTimestampPath().
func (p *XAdES132Path) AllDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP)
}

// IndividualDataObjectsTimestampPath implements XAdESPath. Ports getIndividualDataObjectsTimestampPath().
func (p *XAdES132Path) IndividualDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES132Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP)
}

// DataObjectFormat implements XAdESPath. Ports getDataObjectFormat().
func (p *XAdES132Path) DataObjectFormat() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES132Element_DATA_OBJECT_FORMAT)
}

// DataObjectFormatMimeType implements XAdESPath. Ports getDataObjectFormatMimeType().
func (p *XAdES132Path) DataObjectFormatMimeType() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES132Element_DATA_OBJECT_FORMAT, XAdES132Element_MIME_TYPE)
}

// DataObjectFormatObjectIdentifier implements XAdESPath. Ports getDataObjectFormatObjectIdentifier().
func (p *XAdES132Path) DataObjectFormatObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES132Element_DATA_OBJECT_FORMAT, XAdES132Element_OBJECT_IDENTIFIER)
}

// CommitmentTypeIndicationPath implements XAdESPath. Ports getCommitmentTypeIndicationPath().
func (p *XAdES132Path) CommitmentTypeIndicationPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_SIGNED_PROPERTIES, XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES132Element_COMMITMENT_TYPE_INDICATION)
}

// UnsignedPropertiesPath implements XAdESPath. Ports getUnsignedPropertiesPath().
func (p *XAdES132Path) UnsignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES)
}

// UnsignedSignaturePropertiesPath implements XAdESPath. Ports getUnsignedSignaturePropertiesPath().
func (p *XAdES132Path) UnsignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES)
}

// CounterSignaturePath implements XAdESPath. Ports getCounterSignaturePath().
func (p *XAdES132Path) CounterSignaturePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_COUNTER_SIGNATURE)
}

// AttributeRevocationRefsPath implements XAdESPath. Ports getAttributeRevocationRefsPath().
func (p *XAdES132Path) AttributeRevocationRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_ATTRIBUTE_REVOCATION_REFS)
}

// CompleteRevocationRefsPath implements XAdESPath. Ports getCompleteRevocationRefsPath().
func (p *XAdES132Path) CompleteRevocationRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_COMPLETE_REVOCATION_REFS)
}

// CompleteCertificateRefsPath implements XAdESPath. Ports getCompleteCertificateRefsPath().
func (p *XAdES132Path) CompleteCertificateRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_COMPLETE_CERTIFICATE_REFS)
}

// CompleteCertificateRefsCertPath implements XAdESPath. Ports getCompleteCertificateRefsCertPath().
func (p *XAdES132Path) CompleteCertificateRefsCertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_COMPLETE_CERTIFICATE_REFS, XAdES132Element_CERT_REFS, XAdES132Element_CERT)
}

// CompleteCertificateRefsV2Path implements XAdESPath. Ports getCompleteCertificateRefsV2Path().
func (p *XAdES132Path) CompleteCertificateRefsV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_COMPLETE_CERTIFICATE_REFS_V2)
}

// CompleteCertificateRefsV2CertPath implements XAdESPath. Ports getCompleteCertificateRefsV2CertPath().
func (p *XAdES132Path) CompleteCertificateRefsV2CertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_COMPLETE_CERTIFICATE_REFS_V2, XAdES141Element_CERT_REFS, XAdES132Element_CERT)
}

// AttributeCertificateRefsPath implements XAdESPath. Ports getAttributeCertificateRefsPath().
func (p *XAdES132Path) AttributeCertificateRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS)
}

// AttributeCertificateRefsCertPath implements XAdESPath. Ports getAttributeCertificateRefsCertPath().
func (p *XAdES132Path) AttributeCertificateRefsCertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS, XAdES132Element_CERT_REFS, XAdES132Element_CERT)
}

// AttributeCertificateRefsV2Path implements XAdESPath. Ports getAttributeCertificateRefsV2Path().
func (p *XAdES132Path) AttributeCertificateRefsV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_ATTRIBUTE_CERTIFICATE_REFS_V2)
}

// AttributeCertificateRefsV2CertPath implements XAdESPath. Ports getAttributeCertificateRefsV2CertPath().
func (p *XAdES132Path) AttributeCertificateRefsV2CertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_ATTRIBUTE_CERTIFICATE_REFS_V2, XAdES141Element_CERT_REFS, XAdES132Element_CERT)
}

// CertificateValuesPath implements XAdESPath. Ports getCertificateValuesPath().
func (p *XAdES132Path) CertificateValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_CERTIFICATE_VALUES)
}

// RevocationValuesPath implements XAdESPath. Ports getRevocationValuesPath().
func (p *XAdES132Path) RevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_REVOCATION_VALUES)
}

// AttributeRevocationValuesPath implements XAdESPath. Ports getAttributeRevocationValuesPath().
func (p *XAdES132Path) AttributeRevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_ATTRIBUTE_REVOCATION_VALUES)
}

// EncapsulatedCertificateValuesPath implements XAdESPath. Ports getEncapsulatedCertificateValuesPath().
func (p *XAdES132Path) EncapsulatedCertificateValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_CERTIFICATE_VALUES, XAdES132Element_ENCAPSULATED_X509_CERTIFICATE)
}

// AttrAuthoritiesCertValuesPath implements XAdESPath. Ports getAttrAuthoritiesCertValuesPath().
func (p *XAdES132Path) AttrAuthoritiesCertValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_ATTR_AUTHORITIES_CERT_VALUES)
}

// EncapsulatedAttrAuthoritiesCertValuesPath implements XAdESPath. Ports getEncapsulatedAttrAuthoritiesCertValuesPath().
func (p *XAdES132Path) EncapsulatedAttrAuthoritiesCertValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_ATTR_AUTHORITIES_CERT_VALUES, XAdES132Element_ENCAPSULATED_X509_CERTIFICATE)
}

// EncapsulatedTimeStampValidationDataCertValuesPath implements XAdESPath. Ports getEncapsulatedTimeStampValidationDataCertValuesPath().
func (p *XAdES132Path) EncapsulatedTimeStampValidationDataCertValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_TIMESTAMP_VALIDATION_DATA, XAdES132Element_CERTIFICATE_VALUES, XAdES132Element_ENCAPSULATED_X509_CERTIFICATE)
}

// TimeStampValidationDataRevocationValuesPath implements XAdESPath. Ports getTimeStampValidationDataRevocationValuesPath().
func (p *XAdES132Path) TimeStampValidationDataRevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_TIMESTAMP_VALIDATION_DATA, XAdES132Element_REVOCATION_VALUES)
}

// AnyValidationDataPath implements XAdESPath. Ports getAnyValidationDataPath().
func (p *XAdES132Path) AnyValidationDataPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_ANY_VALIDATION_DATA)
}

// EncapsulatedAnyValidationDataCertValuesPath implements XAdESPath. Ports getEncapsulatedAnyValidationDataCertValuesPath().
func (p *XAdES132Path) EncapsulatedAnyValidationDataCertValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_ANY_VALIDATION_DATA, XAdES132Element_CERTIFICATE_VALUES, XAdES132Element_ENCAPSULATED_X509_CERTIFICATE)
}

// AnyValidationDataRevocationValuesPath implements XAdESPath. Ports getAnyValidationDataRevocationValuesPath().
func (p *XAdES132Path) AnyValidationDataRevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_ANY_VALIDATION_DATA, XAdES132Element_REVOCATION_VALUES)
}

// SignatureTimestampPath implements XAdESPath. Ports getSignatureTimestampPath().
func (p *XAdES132Path) SignatureTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIGNATURE_TIMESTAMP)
}

// SigAndRefsTimestampPath implements XAdESPath. Ports getSigAndRefsTimestampPath().
func (p *XAdES132Path) SigAndRefsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_SIG_AND_REFS_TIMESTAMP)
}

// SigAndRefsTimestampV2Path implements XAdESPath. Ports getSigAndRefsTimestampV2Path().
func (p *XAdES132Path) SigAndRefsTimestampV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_SIG_AND_REFS_TIMESTAMP_V2)
}

// RefsOnlyTimestampPath implements XAdESPath. Ports getRefsOnlyTimestampPath().
func (p *XAdES132Path) RefsOnlyTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES132Element_REFS_ONLY_TIMESTAMP)
}

// RefsOnlyTimestampV2Path implements XAdESPath. Ports getRefsOnlyTimestampV2Path().
func (p *XAdES132Path) RefsOnlyTimestampV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_REFS_ONLY_TIMESTAMP_V2)
}

// ArchiveTimestampPath implements XAdESPath. Ports getArchiveTimestampPath().
func (p *XAdES132Path) ArchiveTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_ARCHIVE_TIMESTAMP)
}

// TimestampValidationDataPath implements XAdESPath. Ports getTimestampValidationDataPath().
func (p *XAdES132Path) TimestampValidationDataPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_TIMESTAMP_VALIDATION_DATA)
}

// SignaturePolicyStorePath implements XAdESPath. Ports getSignaturePolicyStorePath().
func (p *XAdES132Path) SignaturePolicyStorePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES141Element_SIGNATURE_POLICY_STORE)
}

// SealingEvidenceRecordsPath implements XAdESPath. Ports getSealingEvidenceRecordsPath().
func (p *XAdES132Path) SealingEvidenceRecordsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES132Element_QUALIFYING_PROPERTIES, XAdES132Element_UNSIGNED_PROPERTIES, XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdESEvidencerecordNamespaceElement_SEALING_EVIDENCE_RECORDS)
}

// CurrentCRLValuesChildren implements XAdESPath. Ports getCurrentCRLValuesChildren().
func (p *XAdES132Path) CurrentCRLValuesChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CRL_VALUES, XAdES132Element_ENCAPSULATED_CRL_VALUE)
}

// CurrentCRLRefsChildren implements XAdESPath. Ports getCurrentCRLRefsChildren().
func (p *XAdES132Path) CurrentCRLRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CRL_REFS, XAdES132Element_CRL_REF)
}

// CurrentCRLRefCRLIdentifier implements XAdESPath. Ports getCurrentCRLRefCRLIdentifier().
func (p *XAdES132Path) CurrentCRLRefCRLIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CRL_IDENTIFIER)
}

// CurrentCRLRefCRLIdentifierIssuer implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssuer().
func (p *XAdES132Path) CurrentCRLRefCRLIdentifierIssuer() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CRL_IDENTIFIER, XAdES132Element_ISSUER)
}

// CurrentCRLRefCRLIdentifierIssueTime implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssueTime().
func (p *XAdES132Path) CurrentCRLRefCRLIdentifierIssueTime() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CRL_IDENTIFIER, XAdES132Element_ISSUE_TIME)
}

// CurrentCRLRefCRLIdentifierNumber implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierNumber().
func (p *XAdES132Path) CurrentCRLRefCRLIdentifierNumber() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CRL_IDENTIFIER, XAdES132Element_NUMBER)
}

// CurrentOCSPValuesChildren implements XAdESPath. Ports getCurrentOCSPValuesChildren().
func (p *XAdES132Path) CurrentOCSPValuesChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_OCSP_VALUES, XAdES132Element_ENCAPSULATED_OCSP_VALUE)
}

// CurrentOCSPRefsChildren implements XAdESPath. Ports getCurrentOCSPRefsChildren().
func (p *XAdES132Path) CurrentOCSPRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_OCSP_REFS, XAdES132Element_OCSP_REF)
}

// CurrentOCSPRefResponderID implements XAdESPath. Ports getCurrentOCSPRefResponderID().
func (p *XAdES132Path) CurrentOCSPRefResponderID() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_OCSP_IDENTIFIER, XAdES132Element_RESPONDER_ID)
}

// CurrentOCSPRefResponderIDByName implements XAdESPath. Ports getCurrentOCSPRefResponderIDByName().
func (p *XAdES132Path) CurrentOCSPRefResponderIDByName() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_OCSP_IDENTIFIER, XAdES132Element_RESPONDER_ID, XAdES132Element_BY_NAME)
}

// CurrentOCSPRefResponderIDByKey implements XAdESPath. Ports getCurrentOCSPRefResponderIDByKey().
func (p *XAdES132Path) CurrentOCSPRefResponderIDByKey() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_OCSP_IDENTIFIER, XAdES132Element_RESPONDER_ID, XAdES132Element_BY_KEY)
}

// CurrentOCSPRefProducedAt implements XAdESPath. Ports getCurrentOCSPRefProducedAt().
func (p *XAdES132Path) CurrentOCSPRefProducedAt() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_OCSP_IDENTIFIER, XAdES132Element_PRODUCED_AT)
}

// CurrentDigestAlgAndValue implements XAdESPath. Ports getCurrentDigestAlgAndValue().
func (p *XAdES132Path) CurrentDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_DIGEST_ALG_AND_VALUE)
}

// CurrentCertRefsCertChildren implements XAdESPath. Ports getCurrentCertRefsCertChildren().
func (p *XAdES132Path) CurrentCertRefsCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CERT_REFS, XAdES132Element_CERT)
}

// CurrentCertRefs141CertChildren implements XAdESPath. Ports getCurrentCertRefs141CertChildren().
func (p *XAdES132Path) CurrentCertRefs141CertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141Element_CERT_REFS, XAdES132Element_CERT)
}

// CurrentCertChildren implements XAdESPath. Ports getCurrentCertChildren().
func (p *XAdES132Path) CurrentCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CERT)
}

// CurrentCertDigest implements XAdESPath. Ports getCurrentCertDigest().
func (p *XAdES132Path) CurrentCertDigest() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CERT_DIGEST)
}

// CurrentEncapsulatedTimestamp implements XAdESPath. Ports getCurrentEncapsulatedTimestamp().
func (p *XAdES132Path) CurrentEncapsulatedTimestamp() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_ENCAPSULATED_TIMESTAMP)
}

// CurrentEncapsulatedCertificate implements XAdESPath. Ports getCurrentEncapsulatedCertificate().
func (p *XAdES132Path) CurrentEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_ENCAPSULATED_X509_CERTIFICATE)
}

// CurrentCertificateValuesEncapsulatedCertificate implements XAdESPath. Ports getCurrentCertificateValuesEncapsulatedCertificate().
func (p *XAdES132Path) CurrentCertificateValuesEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CERTIFICATE_VALUES, XAdES132Element_ENCAPSULATED_X509_CERTIFICATE)
}

// CurrentRevocationValuesEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedOCSPValue().
func (p *XAdES132Path) CurrentRevocationValuesEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_REVOCATION_VALUES, XAdES132Element_OCSP_VALUES, XAdES132Element_ENCAPSULATED_OCSP_VALUE)
}

// CurrentEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentEncapsulatedOCSPValue().
func (p *XAdES132Path) CurrentEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_OCSP_VALUES, XAdES132Element_ENCAPSULATED_OCSP_VALUE)
}

// CurrentRevocationValuesEncapsulatedCRLValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedCRLValue().
func (p *XAdES132Path) CurrentRevocationValuesEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_REVOCATION_VALUES, XAdES132Element_CRL_VALUES, XAdES132Element_ENCAPSULATED_CRL_VALUE)
}

// CurrentEncapsulatedCRLValue implements XAdESPath. Ports getCurrentEncapsulatedCRLValue().
func (p *XAdES132Path) CurrentEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_CRL_VALUES, XAdES132Element_ENCAPSULATED_CRL_VALUE)
}

// CurrentIssuerSerialIssuerNamePath implements XAdESPath. Ports getCurrentIssuerSerialIssuerNamePath().
func (p *XAdES132Path) CurrentIssuerSerialIssuerNamePath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_ISSUER_SERIAL, common.XMLDSigElement_X509_ISSUER_NAME)
}

// CurrentIssuerSerialSerialNumberPath implements XAdESPath. Ports getCurrentIssuerSerialSerialNumberPath().
func (p *XAdES132Path) CurrentIssuerSerialSerialNumberPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_ISSUER_SERIAL, common.XMLDSigElement_X509_SERIAL_NUMBER)
}

// CurrentIssuerSerialV2Path implements XAdESPath. Ports getCurrentIssuerSerialV2Path().
func (p *XAdES132Path) CurrentIssuerSerialV2Path() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_ISSUER_SERIAL_V2)
}

// CurrentCommitmentIdentifierPath implements XAdESPath. Ports getCurrentCommitmentIdentifierPath().
func (p *XAdES132Path) CurrentCommitmentIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_COMMITMENT_TYPE_ID, XAdES132Element_IDENTIFIER)
}

// CurrentCommitmentDescriptionPath implements XAdESPath. Ports getCurrentCommitmentDescriptionPath().
func (p *XAdES132Path) CurrentCommitmentDescriptionPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_COMMITMENT_TYPE_ID, XAdES132Element_DESCRIPTION)
}

// CurrentCommitmentDocumentationReferencesPath implements XAdESPath. Ports getCurrentCommitmentDocumentationReferencesPath().
func (p *XAdES132Path) CurrentCommitmentDocumentationReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_COMMITMENT_TYPE_ID, XAdES132Element_DOCUMENTATION_REFERENCES)
}

// CurrentDocumentationReference implements XAdESPath. Ports getCurrentDocumentationReference().
func (p *XAdES132Path) CurrentDocumentationReference() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_DOCUMENTATION_REFERENCE)
}

// CurrentDescription implements XAdESPath. Ports getCurrentDescription().
func (p *XAdES132Path) CurrentDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_DESCRIPTION)
}

// CurrentObjectIdentifier implements XAdESPath. Ports getCurrentObjectIdentifier().
func (p *XAdES132Path) CurrentObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_OBJECT_IDENTIFIER)
}

// CurrentCommitmentObjectReferencesPath implements XAdESPath. Ports getCurrentCommitmentObjectReferencesPath().
func (p *XAdES132Path) CurrentCommitmentObjectReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_OBJECT_REFERENCE)
}

// CurrentCommitmentAllSignedDataObjectsPath implements XAdESPath. Ports getCurrentCommitmentAllSignedDataObjectsPath().
func (p *XAdES132Path) CurrentCommitmentAllSignedDataObjectsPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_ALL_SIGNED_DATA_OBJECTS)
}

// CurrentMimeType implements XAdESPath. Ports getCurrentMimeType().
func (p *XAdES132Path) CurrentMimeType() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_MIME_TYPE)
}

// CurrentEncoding implements XAdESPath. Ports getCurrentEncoding().
func (p *XAdES132Path) CurrentEncoding() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_ENCODING)
}

// CurrentSignaturePolicyId implements XAdESPath. Ports getCurrentSignaturePolicyId().
func (p *XAdES132Path) CurrentSignaturePolicyId() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, XAdES132Element_SIG_POLICY_ID, XAdES132Element_IDENTIFIER)
}

// CurrentSignaturePolicyDigestAlgAndValue implements XAdESPath. Ports getCurrentSignaturePolicyDigestAlgAndValue().
func (p *XAdES132Path) CurrentSignaturePolicyDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, XAdES132Element_SIG_POLICY_HASH)
}

// CurrentSignaturePolicySPURI implements XAdESPath. Ports getCurrentSignaturePolicySPURI().
func (p *XAdES132Path) CurrentSignaturePolicySPURI() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, XAdES132Element_SIG_POLICY_QUALIFIERS, XAdES132Element_SIG_POLICY_QUALIFIER, XAdES132Element_SP_URI)
}

// CurrentSignaturePolicySPUserNotice implements XAdESPath. Ports getCurrentSignaturePolicySPUserNotice().
func (p *XAdES132Path) CurrentSignaturePolicySPUserNotice() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, XAdES132Element_SIG_POLICY_QUALIFIERS, XAdES132Element_SIG_POLICY_QUALIFIER, XAdES132Element_SP_USER_NOTICE)
}

// CurrentSPUserNoticeNoticeRefOrganization implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefOrganization().
func (p *XAdES132Path) CurrentSPUserNoticeNoticeRefOrganization() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_NOTICE_REF, XAdES132Element_ORGANIZATION)
}

// CurrentSPUserNoticeNoticeRefNoticeNumbers implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefNoticeNumbers().
func (p *XAdES132Path) CurrentSPUserNoticeNoticeRefNoticeNumbers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_NOTICE_REF, XAdES132Element_NOTICE_NUMBERS)
}

// CurrentSPUserNoticeExplicitText implements XAdESPath. Ports getCurrentSPUserNoticeExplicitText().
func (p *XAdES132Path) CurrentSPUserNoticeExplicitText() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_EXPLICIT_TEXT)
}

// CurrentSignaturePolicySPDocSpecification implements XAdESPath. Ports getCurrentSignaturePolicySPDocSpecification().
func (p *XAdES132Path) CurrentSignaturePolicySPDocSpecification() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, XAdES132Element_SIG_POLICY_QUALIFIERS, XAdES132Element_SIG_POLICY_QUALIFIER, XAdES141Element_SP_DOC_SPECIFICATION)
}

// CurrentSignaturePolicySPDocSpecificationIdentifier implements XAdESPath. Ports getCurrentSignaturePolicySPDocSpecificationIdentifier().
func (p *XAdES132Path) CurrentSignaturePolicySPDocSpecificationIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, XAdES132Element_SIG_POLICY_QUALIFIERS, XAdES132Element_SIG_POLICY_QUALIFIER, XAdES141Element_SP_DOC_SPECIFICATION, XAdES132Element_IDENTIFIER)
}

// CurrentSignaturePolicyDescription implements XAdESPath. Ports getCurrentSignaturePolicyDescription().
func (p *XAdES132Path) CurrentSignaturePolicyDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, XAdES132Element_SIG_POLICY_ID, XAdES132Element_DESCRIPTION)
}

// CurrentSignaturePolicyDocumentationReferences implements XAdESPath. Ports getCurrentSignaturePolicyDocumentationReferences().
func (p *XAdES132Path) CurrentSignaturePolicyDocumentationReferences() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, XAdES132Element_SIG_POLICY_ID, XAdES132Element_DOCUMENTATION_REFERENCES)
}

// CurrentSignaturePolicyImplied implements XAdESPath. Ports getCurrentSignaturePolicyImplied().
func (p *XAdES132Path) CurrentSignaturePolicyImplied() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_IMPLIED)
}

// CurrentSignaturePolicyTransforms implements XAdESPath. Ports getCurrentSignaturePolicyTransforms().
func (p *XAdES132Path) CurrentSignaturePolicyTransforms() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, common.XMLDSigElement_TRANSFORMS)
}

// CurrentSignaturePolicyQualifiers implements XAdESPath. Ports getCurrentSignaturePolicyQualifiers().
func (p *XAdES132Path) CurrentSignaturePolicyQualifiers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_SIGNATURE_POLICY_ID, XAdES132Element_SIG_POLICY_QUALIFIERS)
}

// CurrentInclude implements XAdESPath. Ports getCurrentInclude().
func (p *XAdES132Path) CurrentInclude() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_INCLUDE)
}

// CurrentQualifyingPropertiesPath implements XAdESPath. Ports getCurrentQualifyingPropertiesPath().
func (p *XAdES132Path) CurrentQualifyingPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_QUALIFYING_PROPERTIES)
}

// CurrentSPDocSpecification implements XAdESPath. Ports getCurrentSPDocSpecification().
func (p *XAdES132Path) CurrentSPDocSpecification() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141Element_SP_DOC_SPECIFICATION)
}

// CurrentIdentifier implements XAdESPath. Ports getCurrentIdentifier().
func (p *XAdES132Path) CurrentIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_IDENTIFIER)
}

// CurrentSPDocSpecificationIdentifier implements XAdESPath. Ports getCurrentSPDocSpecificationIdentifier().
func (p *XAdES132Path) CurrentSPDocSpecificationIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141Element_SP_DOC_SPECIFICATION, XAdES132Element_IDENTIFIER)
}

// CurrentSPDocSpecificationDescription implements XAdESPath. Ports getCurrentSPDocSpecificationDescription().
func (p *XAdES132Path) CurrentSPDocSpecificationDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141Element_SP_DOC_SPECIFICATION, XAdES132Element_DESCRIPTION)
}

// CurrentDocumentationReferenceElements implements XAdESPath. Ports getCurrentDocumentationReferenceElements().
func (p *XAdES132Path) CurrentDocumentationReferenceElements() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132Element_DOCUMENTATION_REFERENCES, XAdES132Element_DOCUMENTATION_REFERENCE)
}

// CurrentSPDocSpecificationDocumentationReferenceElements implements XAdESPath. Ports getCurrentSPDocSpecificationDocumentationReferenceElements().
func (p *XAdES132Path) CurrentSPDocSpecificationDocumentationReferenceElements() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141Element_SP_DOC_SPECIFICATION, XAdES132Element_DOCUMENTATION_REFERENCES, XAdES132Element_DOCUMENTATION_REFERENCE)
}

// CurrentSignaturePolicyDocument implements XAdESPath. Ports getCurrentSignaturePolicyDocument().
func (p *XAdES132Path) CurrentSignaturePolicyDocument() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141Element_SIGNATURE_POLICY_DOCUMENT)
}

// CurrentSigPolDocLocalURI implements XAdESPath. Ports getCurrentSigPolDocLocalURI().
func (p *XAdES132Path) CurrentSigPolDocLocalURI() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141Element_SIG_POL_DOC_LOCAL_URI)
}

var _ XAdESPath = (*XAdES132Path)(nil)
