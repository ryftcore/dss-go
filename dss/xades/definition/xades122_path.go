// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades122/XAdES122Path.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES122Path holds the XAdES 122 paths.
type XAdES122Path struct{}

// NewXAdES122Path creates a XAdES122Path.
func NewXAdES122Path() *XAdES122Path {
	return &XAdES122Path{}
}

// Namespace implements XAdESPath. Ports getNamespace().
func (p *XAdES122Path) Namespace() *common.DSSNamespace {
	return XAdESNamespace_XADES_122
}

// SignedPropertiesUri implements XAdESPath. Ports getSignedPropertiesUri().
func (p *XAdES122Path) SignedPropertiesUri() string {
	return "http://uri.etsi.org/01903/v1.2.2#SignedProperties"
}

// CounterSignatureUri implements XAdESPath. Ports getCounterSignatureUri().
func (p *XAdES122Path) CounterSignatureUri() string {
	return "http://uri.etsi.org/01903#CountersignedSignature"
}

// QualifyingPropertiesPath implements XAdESPath. Ports getQualifyingPropertiesPath().
func (p *XAdES122Path) QualifyingPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES)
}

// SignedPropertiesPath implements XAdESPath. Ports getSignedPropertiesPath().
func (p *XAdES122Path) SignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES)
}

// SignedSignaturePropertiesPath implements XAdESPath. Ports getSignedSignaturePropertiesPath().
func (p *XAdES122Path) SignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_SIGNATURE_PROPERTIES)
}

// SigningTimePath implements XAdESPath. Ports getSigningTimePath().
func (p *XAdES122Path) SigningTimePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIGNING_TIME)
}

// SigningCertificatePath implements XAdESPath. Ports getSigningCertificatePath().
func (p *XAdES122Path) SigningCertificatePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIGNING_CERTIFICATE)
}

// SigningCertificateChildren implements XAdESPath. Ports getSigningCertificateChildren().
func (p *XAdES122Path) SigningCertificateChildren() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIGNING_CERTIFICATE, XAdES122Element_CERT)
}

// SigningCertificateV2Path implements XAdESPath. Ports getSigningCertificateV2Path().
func (p *XAdES122Path) SigningCertificateV2Path() common.XPathQuery {
	return nil
}

// SigningCertificateV2Children implements XAdESPath. Ports getSigningCertificateV2Children().
func (p *XAdES122Path) SigningCertificateV2Children() common.XPathQuery {
	return nil
}

// SignatureProductionPlacePath implements XAdESPath. Ports getSignatureProductionPlacePath().
func (p *XAdES122Path) SignatureProductionPlacePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIGNATURE_PRODUCTION_PLACE)
}

// SignatureProductionPlaceV2Path implements XAdESPath. Ports getSignatureProductionPlaceV2Path().
func (p *XAdES122Path) SignatureProductionPlaceV2Path() common.XPathQuery {
	return nil
}

// SignaturePolicyIdentifierPath implements XAdESPath. Ports getSignaturePolicyIdentifierPath().
func (p *XAdES122Path) SignaturePolicyIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIGNATURE_POLICY_IDENTIFIER)
}

// SignerRolePath implements XAdESPath. Ports getSignerRolePath().
func (p *XAdES122Path) SignerRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIGNER_ROLE)
}

// ClaimedRolePath implements XAdESPath. Ports getClaimedRolePath().
func (p *XAdES122Path) ClaimedRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIGNER_ROLE, XAdES122Element_CLAIMED_ROLES, XAdES122Element_CLAIMED_ROLE)
}

// SignedAssertionPath implements XAdESPath. Ports getSignedAssertionPath().
func (p *XAdES122Path) SignedAssertionPath() common.XPathQuery {
	return nil
}

// SignerRoleV2Path implements XAdESPath. Ports getSignerRoleV2Path().
func (p *XAdES122Path) SignerRoleV2Path() common.XPathQuery {
	return nil
}

// ClaimedRoleV2Path implements XAdESPath. Ports getClaimedRoleV2Path().
func (p *XAdES122Path) ClaimedRoleV2Path() common.XPathQuery {
	return nil
}

// CertifiedRolePath implements XAdESPath. Ports getCertifiedRolePath().
func (p *XAdES122Path) CertifiedRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIGNER_ROLE, XAdES122Element_CERTIFIED_ROLES, XAdES122Element_CERTIFIED_ROLE)
}

// CertifiedRoleV2Path implements XAdESPath. Ports getCertifiedRoleV2Path().
func (p *XAdES122Path) CertifiedRoleV2Path() common.XPathQuery {
	return nil
}

// SignedDataObjectPropertiesPath implements XAdESPath. Ports getSignedDataObjectPropertiesPath().
func (p *XAdES122Path) SignedDataObjectPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES)
}

// AllDataObjectsTimestampPath implements XAdESPath. Ports getAllDataObjectsTimestampPath().
func (p *XAdES122Path) AllDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP)
}

// IndividualDataObjectsTimestampPath implements XAdESPath. Ports getIndividualDataObjectsTimestampPath().
func (p *XAdES122Path) IndividualDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES122Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP)
}

// DataObjectFormat implements XAdESPath. Ports getDataObjectFormat().
func (p *XAdES122Path) DataObjectFormat() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES122Element_DATA_OBJECT_FORMAT)
}

// DataObjectFormatMimeType implements XAdESPath. Ports getDataObjectFormatMimeType().
func (p *XAdES122Path) DataObjectFormatMimeType() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES122Element_DATA_OBJECT_FORMAT, XAdES122Element_MIME_TYPE)
}

// DataObjectFormatObjectIdentifier implements XAdESPath. Ports getDataObjectFormatObjectIdentifier().
func (p *XAdES122Path) DataObjectFormatObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES122Element_DATA_OBJECT_FORMAT, XAdES122Element_OBJECT_IDENTIFIER)
}

// CommitmentTypeIndicationPath implements XAdESPath. Ports getCommitmentTypeIndicationPath().
func (p *XAdES122Path) CommitmentTypeIndicationPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_SIGNED_PROPERTIES, XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES122Element_COMMITMENT_TYPE_INDICATION)
}

// UnsignedPropertiesPath implements XAdESPath. Ports getUnsignedPropertiesPath().
func (p *XAdES122Path) UnsignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES)
}

// UnsignedSignaturePropertiesPath implements XAdESPath. Ports getUnsignedSignaturePropertiesPath().
func (p *XAdES122Path) UnsignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES)
}

// CounterSignaturePath implements XAdESPath. Ports getCounterSignaturePath().
func (p *XAdES122Path) CounterSignaturePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_COUNTER_SIGNATURE)
}

// AttributeRevocationRefsPath implements XAdESPath. Ports getAttributeRevocationRefsPath().
func (p *XAdES122Path) AttributeRevocationRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_ATTRIBUTE_REVOCATION_REFS)
}

// CompleteRevocationRefsPath implements XAdESPath. Ports getCompleteRevocationRefsPath().
func (p *XAdES122Path) CompleteRevocationRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_COMPLETE_REVOCATION_REFS)
}

// CompleteCertificateRefsPath implements XAdESPath. Ports getCompleteCertificateRefsPath().
func (p *XAdES122Path) CompleteCertificateRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_COMPLETE_CERTIFICATE_REFS)
}

// CompleteCertificateRefsCertPath implements XAdESPath. Ports getCompleteCertificateRefsCertPath().
func (p *XAdES122Path) CompleteCertificateRefsCertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_COMPLETE_CERTIFICATE_REFS, XAdES122Element_CERT_REFS, XAdES122Element_CERT)
}

// CompleteCertificateRefsV2Path implements XAdESPath. Ports getCompleteCertificateRefsV2Path().
func (p *XAdES122Path) CompleteCertificateRefsV2Path() common.XPathQuery {
	return nil
}

// CompleteCertificateRefsV2CertPath implements XAdESPath. Ports getCompleteCertificateRefsV2CertPath().
func (p *XAdES122Path) CompleteCertificateRefsV2CertPath() common.XPathQuery {
	return nil
}

// AttributeCertificateRefsPath implements XAdESPath. Ports getAttributeCertificateRefsPath().
func (p *XAdES122Path) AttributeCertificateRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_ATTRIBUTE_CERTIFICATE_REFS)
}

// AttributeCertificateRefsCertPath implements XAdESPath. Ports getAttributeCertificateRefsCertPath().
func (p *XAdES122Path) AttributeCertificateRefsCertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_ATTRIBUTE_CERTIFICATE_REFS, XAdES122Element_CERT_REFS, XAdES122Element_CERT)
}

// AttributeCertificateRefsV2Path implements XAdESPath. Ports getAttributeCertificateRefsV2Path().
func (p *XAdES122Path) AttributeCertificateRefsV2Path() common.XPathQuery {
	return nil
}

// AttributeCertificateRefsV2CertPath implements XAdESPath. Ports getAttributeCertificateRefsV2CertPath().
func (p *XAdES122Path) AttributeCertificateRefsV2CertPath() common.XPathQuery {
	return nil
}

// CertificateValuesPath implements XAdESPath. Ports getCertificateValuesPath().
func (p *XAdES122Path) CertificateValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_CERTIFICATE_VALUES)
}

// RevocationValuesPath implements XAdESPath. Ports getRevocationValuesPath().
func (p *XAdES122Path) RevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_REVOCATION_VALUES)
}

// AttributeRevocationValuesPath implements XAdESPath. Ports getAttributeRevocationValuesPath().
func (p *XAdES122Path) AttributeRevocationValuesPath() common.XPathQuery {
	return nil
}

// EncapsulatedCertificateValuesPath implements XAdESPath. Ports getEncapsulatedCertificateValuesPath().
func (p *XAdES122Path) EncapsulatedCertificateValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_CERTIFICATE_VALUES, XAdES122Element_ENCAPSULATED_X509_CERTIFICATE)
}

// AttrAuthoritiesCertValuesPath implements XAdESPath. Ports getAttrAuthoritiesCertValuesPath().
func (p *XAdES122Path) AttrAuthoritiesCertValuesPath() common.XPathQuery {
	return nil
}

// EncapsulatedAttrAuthoritiesCertValuesPath implements XAdESPath. Ports getEncapsulatedAttrAuthoritiesCertValuesPath().
func (p *XAdES122Path) EncapsulatedAttrAuthoritiesCertValuesPath() common.XPathQuery {
	return nil
}

// EncapsulatedTimeStampValidationDataCertValuesPath implements XAdESPath. Ports getEncapsulatedTimeStampValidationDataCertValuesPath().
func (p *XAdES122Path) EncapsulatedTimeStampValidationDataCertValuesPath() common.XPathQuery {
	return nil
}

// TimeStampValidationDataRevocationValuesPath implements XAdESPath. Ports getTimeStampValidationDataRevocationValuesPath().
func (p *XAdES122Path) TimeStampValidationDataRevocationValuesPath() common.XPathQuery {
	return nil
}

// AnyValidationDataPath implements XAdESPath. Ports getAnyValidationDataPath().
func (p *XAdES122Path) AnyValidationDataPath() common.XPathQuery {
	return nil
}

// EncapsulatedAnyValidationDataCertValuesPath implements XAdESPath. Ports getEncapsulatedAnyValidationDataCertValuesPath().
func (p *XAdES122Path) EncapsulatedAnyValidationDataCertValuesPath() common.XPathQuery {
	return nil
}

// AnyValidationDataRevocationValuesPath implements XAdESPath. Ports getAnyValidationDataRevocationValuesPath().
func (p *XAdES122Path) AnyValidationDataRevocationValuesPath() common.XPathQuery {
	return nil
}

// SignatureTimestampPath implements XAdESPath. Ports getSignatureTimestampPath().
func (p *XAdES122Path) SignatureTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIGNATURE_TIMESTAMP)
}

// SigAndRefsTimestampPath implements XAdESPath. Ports getSigAndRefsTimestampPath().
func (p *XAdES122Path) SigAndRefsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES122Element_QUALIFYING_PROPERTIES, XAdES122Element_UNSIGNED_PROPERTIES, XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES122Element_SIG_AND_REFS_TIMESTAMP)
}

// SigAndRefsTimestampV2Path implements XAdESPath. Ports getSigAndRefsTimestampV2Path().
func (p *XAdES122Path) SigAndRefsTimestampV2Path() common.XPathQuery {
	return nil
}

// RefsOnlyTimestampPath implements XAdESPath. Ports getRefsOnlyTimestampPath().
func (p *XAdES122Path) RefsOnlyTimestampPath() common.XPathQuery {
	return nil
}

// RefsOnlyTimestampV2Path implements XAdESPath. Ports getRefsOnlyTimestampV2Path().
func (p *XAdES122Path) RefsOnlyTimestampV2Path() common.XPathQuery {
	return nil
}

// ArchiveTimestampPath implements XAdESPath. Ports getArchiveTimestampPath().
func (p *XAdES122Path) ArchiveTimestampPath() common.XPathQuery {
	return nil
}

// TimestampValidationDataPath implements XAdESPath. Ports getTimestampValidationDataPath().
func (p *XAdES122Path) TimestampValidationDataPath() common.XPathQuery {
	return nil
}

// SignaturePolicyStorePath implements XAdESPath. Ports getSignaturePolicyStorePath().
func (p *XAdES122Path) SignaturePolicyStorePath() common.XPathQuery {
	return nil
}

// SealingEvidenceRecordsPath implements XAdESPath. Ports getSealingEvidenceRecordsPath().
func (p *XAdES122Path) SealingEvidenceRecordsPath() common.XPathQuery {
	return nil
}

// CurrentCRLValuesChildren implements XAdESPath. Ports getCurrentCRLValuesChildren().
func (p *XAdES122Path) CurrentCRLValuesChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CRL_VALUES, XAdES122Element_ENCAPSULATED_CRL_VALUE)
}

// CurrentCRLRefsChildren implements XAdESPath. Ports getCurrentCRLRefsChildren().
func (p *XAdES122Path) CurrentCRLRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CRL_REFS, XAdES122Element_CRL_REF)
}

// CurrentCRLRefCRLIdentifier implements XAdESPath. Ports getCurrentCRLRefCRLIdentifier().
func (p *XAdES122Path) CurrentCRLRefCRLIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CRL_IDENTIFIER)
}

// CurrentCRLRefCRLIdentifierIssuer implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssuer().
func (p *XAdES122Path) CurrentCRLRefCRLIdentifierIssuer() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CRL_IDENTIFIER, XAdES122Element_ISSUER)
}

// CurrentCRLRefCRLIdentifierIssueTime implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssueTime().
func (p *XAdES122Path) CurrentCRLRefCRLIdentifierIssueTime() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CRL_IDENTIFIER, XAdES122Element_ISSUE_TIME)
}

// CurrentCRLRefCRLIdentifierNumber implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierNumber().
func (p *XAdES122Path) CurrentCRLRefCRLIdentifierNumber() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CRL_IDENTIFIER, XAdES122Element_NUMBER)
}

// CurrentOCSPValuesChildren implements XAdESPath. Ports getCurrentOCSPValuesChildren().
func (p *XAdES122Path) CurrentOCSPValuesChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_OCSP_VALUES, XAdES122Element_ENCAPSULATED_OCSP_VALUE)
}

// CurrentOCSPRefsChildren implements XAdESPath. Ports getCurrentOCSPRefsChildren().
func (p *XAdES122Path) CurrentOCSPRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_OCSP_REFS, XAdES122Element_OCSP_REF)
}

// CurrentOCSPRefResponderID implements XAdESPath. Ports getCurrentOCSPRefResponderID().
func (p *XAdES122Path) CurrentOCSPRefResponderID() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_OCSP_IDENTIFIER, XAdES122Element_RESPONDER_ID)
}

// CurrentOCSPRefResponderIDByName implements XAdESPath. Ports getCurrentOCSPRefResponderIDByName().
func (p *XAdES122Path) CurrentOCSPRefResponderIDByName() common.XPathQuery {
	return nil
}

// CurrentOCSPRefResponderIDByKey implements XAdESPath. Ports getCurrentOCSPRefResponderIDByKey().
func (p *XAdES122Path) CurrentOCSPRefResponderIDByKey() common.XPathQuery {
	return nil
}

// CurrentOCSPRefProducedAt implements XAdESPath. Ports getCurrentOCSPRefProducedAt().
func (p *XAdES122Path) CurrentOCSPRefProducedAt() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_OCSP_IDENTIFIER, XAdES122Element_PRODUCED_AT)
}

// CurrentDigestAlgAndValue implements XAdESPath. Ports getCurrentDigestAlgAndValue().
func (p *XAdES122Path) CurrentDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_DIGEST_ALG_AND_VALUE)
}

// CurrentCertRefsCertChildren implements XAdESPath. Ports getCurrentCertRefsCertChildren().
func (p *XAdES122Path) CurrentCertRefsCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CERT_REFS, XAdES122Element_CERT)
}

// CurrentCertRefs141CertChildren implements XAdESPath. Ports getCurrentCertRefs141CertChildren().
func (p *XAdES122Path) CurrentCertRefs141CertChildren() common.XPathQuery {
	return nil
}

// CurrentCertChildren implements XAdESPath. Ports getCurrentCertChildren().
func (p *XAdES122Path) CurrentCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CERT)
}

// CurrentCertDigest implements XAdESPath. Ports getCurrentCertDigest().
func (p *XAdES122Path) CurrentCertDigest() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CERT_DIGEST)
}

// CurrentEncapsulatedTimestamp implements XAdESPath. Ports getCurrentEncapsulatedTimestamp().
func (p *XAdES122Path) CurrentEncapsulatedTimestamp() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_ENCAPSULATED_TIMESTAMP)
}

// CurrentEncapsulatedCertificate implements XAdESPath. Ports getCurrentEncapsulatedCertificate().
func (p *XAdES122Path) CurrentEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_ENCAPSULATED_X509_CERTIFICATE)
}

// CurrentCertificateValuesEncapsulatedCertificate implements XAdESPath. Ports getCurrentCertificateValuesEncapsulatedCertificate().
func (p *XAdES122Path) CurrentCertificateValuesEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CERTIFICATE_VALUES, XAdES122Element_ENCAPSULATED_X509_CERTIFICATE)
}

// CurrentRevocationValuesEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedOCSPValue().
func (p *XAdES122Path) CurrentRevocationValuesEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_REVOCATION_VALUES, XAdES122Element_OCSP_VALUES, XAdES122Element_ENCAPSULATED_OCSP_VALUE)
}

// CurrentEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentEncapsulatedOCSPValue().
func (p *XAdES122Path) CurrentEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_OCSP_VALUES, XAdES122Element_ENCAPSULATED_OCSP_VALUE)
}

// CurrentRevocationValuesEncapsulatedCRLValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedCRLValue().
func (p *XAdES122Path) CurrentRevocationValuesEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_REVOCATION_VALUES, XAdES122Element_CRL_VALUES, XAdES122Element_ENCAPSULATED_CRL_VALUE)
}

// CurrentEncapsulatedCRLValue implements XAdESPath. Ports getCurrentEncapsulatedCRLValue().
func (p *XAdES122Path) CurrentEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_CRL_VALUES, XAdES122Element_ENCAPSULATED_CRL_VALUE)
}

// CurrentIssuerSerialIssuerNamePath implements XAdESPath. Ports getCurrentIssuerSerialIssuerNamePath().
func (p *XAdES122Path) CurrentIssuerSerialIssuerNamePath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_ISSUER_SERIAL, common.XMLDSigElement_X509_ISSUER_NAME)
}

// CurrentIssuerSerialSerialNumberPath implements XAdESPath. Ports getCurrentIssuerSerialSerialNumberPath().
func (p *XAdES122Path) CurrentIssuerSerialSerialNumberPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_ISSUER_SERIAL, common.XMLDSigElement_X509_SERIAL_NUMBER)
}

// CurrentIssuerSerialV2Path implements XAdESPath. Ports getCurrentIssuerSerialV2Path().
func (p *XAdES122Path) CurrentIssuerSerialV2Path() common.XPathQuery {
	return nil
}

// CurrentCommitmentIdentifierPath implements XAdESPath. Ports getCurrentCommitmentIdentifierPath().
func (p *XAdES122Path) CurrentCommitmentIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_COMMITMENT_TYPE_ID, XAdES122Element_IDENTIFIER)
}

// CurrentCommitmentDescriptionPath implements XAdESPath. Ports getCurrentCommitmentDescriptionPath().
func (p *XAdES122Path) CurrentCommitmentDescriptionPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_COMMITMENT_TYPE_ID, XAdES122Element_DESCRIPTION)
}

// CurrentCommitmentDocumentationReferencesPath implements XAdESPath. Ports getCurrentCommitmentDocumentationReferencesPath().
func (p *XAdES122Path) CurrentCommitmentDocumentationReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_COMMITMENT_TYPE_ID, XAdES122Element_DOCUMENTATION_REFERENCES)
}

// CurrentDocumentationReference implements XAdESPath. Ports getCurrentDocumentationReference().
func (p *XAdES122Path) CurrentDocumentationReference() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_DOCUMENTATION_REFERENCE)
}

// CurrentDescription implements XAdESPath. Ports getCurrentDescription().
func (p *XAdES122Path) CurrentDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_DESCRIPTION)
}

// CurrentObjectIdentifier implements XAdESPath. Ports getCurrentObjectIdentifier().
func (p *XAdES122Path) CurrentObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_OBJECT_IDENTIFIER)
}

// CurrentCommitmentObjectReferencesPath implements XAdESPath. Ports getCurrentCommitmentObjectReferencesPath().
func (p *XAdES122Path) CurrentCommitmentObjectReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_OBJECT_REFERENCE)
}

// CurrentCommitmentAllSignedDataObjectsPath implements XAdESPath. Ports getCurrentCommitmentAllSignedDataObjectsPath().
func (p *XAdES122Path) CurrentCommitmentAllSignedDataObjectsPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_ALL_SIGNED_DATA_OBJECTS)
}

// CurrentMimeType implements XAdESPath. Ports getCurrentMimeType().
func (p *XAdES122Path) CurrentMimeType() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_MIME_TYPE)
}

// CurrentEncoding implements XAdESPath. Ports getCurrentEncoding().
func (p *XAdES122Path) CurrentEncoding() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_ENCODING)
}

// CurrentSignaturePolicyId implements XAdESPath. Ports getCurrentSignaturePolicyId().
func (p *XAdES122Path) CurrentSignaturePolicyId() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_SIGNATURE_POLICY_ID, XAdES122Element_SIG_POLICY_ID, XAdES122Element_IDENTIFIER)
}

// CurrentSignaturePolicyDigestAlgAndValue implements XAdESPath. Ports getCurrentSignaturePolicyDigestAlgAndValue().
func (p *XAdES122Path) CurrentSignaturePolicyDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_SIGNATURE_POLICY_ID, XAdES122Element_SIG_POLICY_HASH)
}

// CurrentSignaturePolicySPURI implements XAdESPath. Ports getCurrentSignaturePolicySPURI().
func (p *XAdES122Path) CurrentSignaturePolicySPURI() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_SIGNATURE_POLICY_ID, XAdES122Element_SIG_POLICY_QUALIFIERS, XAdES122Element_SIG_POLICY_QUALIFIER, XAdES122Element_SP_URI)
}

// CurrentSignaturePolicySPUserNotice implements XAdESPath. Ports getCurrentSignaturePolicySPUserNotice().
func (p *XAdES122Path) CurrentSignaturePolicySPUserNotice() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_SIGNATURE_POLICY_ID, XAdES122Element_SIG_POLICY_QUALIFIERS, XAdES122Element_SIG_POLICY_QUALIFIER, XAdES122Element_SP_USER_NOTICE)
}

// CurrentSPUserNoticeNoticeRefOrganization implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefOrganization().
func (p *XAdES122Path) CurrentSPUserNoticeNoticeRefOrganization() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_NOTICE_REF, XAdES122Element_ORGANIZATION)
}

// CurrentSPUserNoticeNoticeRefNoticeNumbers implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefNoticeNumbers().
func (p *XAdES122Path) CurrentSPUserNoticeNoticeRefNoticeNumbers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_NOTICE_REF, XAdES122Element_NOTICE_NUMBERS)
}

// CurrentSPUserNoticeExplicitText implements XAdESPath. Ports getCurrentSPUserNoticeExplicitText().
func (p *XAdES122Path) CurrentSPUserNoticeExplicitText() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_EXPLICIT_TEXT)
}

// CurrentSignaturePolicySPDocSpecification implements XAdESPath. Ports getCurrentSignaturePolicySPDocSpecification().
func (p *XAdES122Path) CurrentSignaturePolicySPDocSpecification() common.XPathQuery {
	return nil
}

// CurrentSignaturePolicySPDocSpecificationIdentifier implements XAdESPath. Ports getCurrentSignaturePolicySPDocSpecificationIdentifier().
func (p *XAdES122Path) CurrentSignaturePolicySPDocSpecificationIdentifier() common.XPathQuery {
	return nil
}

// CurrentSignaturePolicyDescription implements XAdESPath. Ports getCurrentSignaturePolicyDescription().
func (p *XAdES122Path) CurrentSignaturePolicyDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_SIGNATURE_POLICY_ID, XAdES122Element_SIG_POLICY_ID, XAdES122Element_DESCRIPTION)
}

// CurrentSignaturePolicyDocumentationReferences implements XAdESPath. Ports getCurrentSignaturePolicyDocumentationReferences().
func (p *XAdES122Path) CurrentSignaturePolicyDocumentationReferences() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_SIGNATURE_POLICY_ID, XAdES122Element_SIG_POLICY_ID, XAdES122Element_DOCUMENTATION_REFERENCES)
}

// CurrentSignaturePolicyImplied implements XAdESPath. Ports getCurrentSignaturePolicyImplied().
func (p *XAdES122Path) CurrentSignaturePolicyImplied() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_SIGNATURE_POLICY_IMPLIED)
}

// CurrentSignaturePolicyTransforms implements XAdESPath. Ports getCurrentSignaturePolicyTransforms().
func (p *XAdES122Path) CurrentSignaturePolicyTransforms() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_SIGNATURE_POLICY_ID, common.XMLDSigElement_TRANSFORMS)
}

// CurrentSignaturePolicyQualifiers implements XAdESPath. Ports getCurrentSignaturePolicyQualifiers().
func (p *XAdES122Path) CurrentSignaturePolicyQualifiers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_SIGNATURE_POLICY_ID, XAdES122Element_SIG_POLICY_QUALIFIERS)
}

// CurrentInclude implements XAdESPath. Ports getCurrentInclude().
func (p *XAdES122Path) CurrentInclude() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_INCLUDE)
}

// CurrentQualifyingPropertiesPath implements XAdESPath. Ports getCurrentQualifyingPropertiesPath().
func (p *XAdES122Path) CurrentQualifyingPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_QUALIFYING_PROPERTIES)
}

// CurrentSPDocSpecification implements XAdESPath. Ports getCurrentSPDocSpecification().
func (p *XAdES122Path) CurrentSPDocSpecification() common.XPathQuery {
	return nil
}

// CurrentIdentifier implements XAdESPath. Ports getCurrentIdentifier().
func (p *XAdES122Path) CurrentIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES122Element_IDENTIFIER)
}

// CurrentSPDocSpecificationIdentifier implements XAdESPath. Ports getCurrentSPDocSpecificationIdentifier().
func (p *XAdES122Path) CurrentSPDocSpecificationIdentifier() common.XPathQuery {
	return nil
}

// CurrentSPDocSpecificationDescription implements XAdESPath. Ports getCurrentSPDocSpecificationDescription().
func (p *XAdES122Path) CurrentSPDocSpecificationDescription() common.XPathQuery {
	return nil
}

// CurrentDocumentationReferenceElements implements XAdESPath. Ports getCurrentDocumentationReferenceElements().
func (p *XAdES122Path) CurrentDocumentationReferenceElements() common.XPathQuery {
	return nil
}

// CurrentSPDocSpecificationDocumentationReferenceElements implements XAdESPath. Ports getCurrentSPDocSpecificationDocumentationReferenceElements().
func (p *XAdES122Path) CurrentSPDocSpecificationDocumentationReferenceElements() common.XPathQuery {
	return nil
}

// CurrentSignaturePolicyDocument implements XAdESPath. Ports getCurrentSignaturePolicyDocument().
func (p *XAdES122Path) CurrentSignaturePolicyDocument() common.XPathQuery {
	return nil
}

// CurrentSigPolDocLocalURI implements XAdESPath. Ports getCurrentSigPolDocLocalURI().
func (p *XAdES122Path) CurrentSigPolDocLocalURI() common.XPathQuery {
	return nil
}

var _ XAdESPath = (*XAdES122Path)(nil)
