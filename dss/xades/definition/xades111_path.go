// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades111/XAdES111Path.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// XAdES111Path holds the XAdES 111 paths.
type XAdES111Path struct{}

// NewXAdES111Path creates a XAdES111Path.
func NewXAdES111Path() *XAdES111Path {
	return &XAdES111Path{}
}

// Namespace implements XAdESPath. Ports getNamespace().
func (p *XAdES111Path) Namespace() *common.DSSNamespace {
	return XAdESNamespace_XADES_111
}

// SignedPropertiesUri implements XAdESPath. Ports getSignedPropertiesUri().
func (p *XAdES111Path) SignedPropertiesUri() string {
	return "http://uri.etsi.org/01903/v1.1.1#SignedProperties"
}

// CounterSignatureUri implements XAdESPath. Ports getCounterSignatureUri().
func (p *XAdES111Path) CounterSignatureUri() string {
	return "http://uri.etsi.org/01903#CountersignedSignature"
}

// QualifyingPropertiesPath implements XAdESPath. Ports getQualifyingPropertiesPath().
func (p *XAdES111Path) QualifyingPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES)
}

// SignedPropertiesPath implements XAdESPath. Ports getSignedPropertiesPath().
func (p *XAdES111Path) SignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES)
}

// SignedSignaturePropertiesPath implements XAdESPath. Ports getSignedSignaturePropertiesPath().
func (p *XAdES111Path) SignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_SIGNATURE_PROPERTIES)
}

// SigningTimePath implements XAdESPath. Ports getSigningTimePath().
func (p *XAdES111Path) SigningTimePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIGNING_TIME)
}

// SigningCertificatePath implements XAdESPath. Ports getSigningCertificatePath().
func (p *XAdES111Path) SigningCertificatePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIGNING_CERTIFICATE)
}

// SigningCertificateChildren implements XAdESPath. Ports getSigningCertificateChildren().
func (p *XAdES111Path) SigningCertificateChildren() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIGNING_CERTIFICATE, XAdES111Element_CERT)
}

// SigningCertificateV2Path implements XAdESPath. Ports getSigningCertificateV2Path().
func (p *XAdES111Path) SigningCertificateV2Path() common.XPathQuery {
	return nil
}

// SigningCertificateV2Children implements XAdESPath. Ports getSigningCertificateV2Children().
func (p *XAdES111Path) SigningCertificateV2Children() common.XPathQuery {
	return nil
}

// SignatureProductionPlacePath implements XAdESPath. Ports getSignatureProductionPlacePath().
func (p *XAdES111Path) SignatureProductionPlacePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIGNATURE_PRODUCTION_PLACE)
}

// SignatureProductionPlaceV2Path implements XAdESPath. Ports getSignatureProductionPlaceV2Path().
func (p *XAdES111Path) SignatureProductionPlaceV2Path() common.XPathQuery {
	return nil
}

// SignaturePolicyIdentifierPath implements XAdESPath. Ports getSignaturePolicyIdentifierPath().
func (p *XAdES111Path) SignaturePolicyIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIGNATURE_POLICY_IDENTIFIER)
}

// SignerRolePath implements XAdESPath. Ports getSignerRolePath().
func (p *XAdES111Path) SignerRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIGNER_ROLE)
}

// ClaimedRolePath implements XAdESPath. Ports getClaimedRolePath().
func (p *XAdES111Path) ClaimedRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIGNER_ROLE, XAdES111Element_CLAIMED_ROLES, XAdES111Element_CLAIMED_ROLE)
}

// SignedAssertionPath implements XAdESPath. Ports getSignedAssertionPath().
func (p *XAdES111Path) SignedAssertionPath() common.XPathQuery {
	return nil
}

// SignerRoleV2Path implements XAdESPath. Ports getSignerRoleV2Path().
func (p *XAdES111Path) SignerRoleV2Path() common.XPathQuery {
	return nil
}

// ClaimedRoleV2Path implements XAdESPath. Ports getClaimedRoleV2Path().
func (p *XAdES111Path) ClaimedRoleV2Path() common.XPathQuery {
	return nil
}

// CertifiedRolePath implements XAdESPath. Ports getCertifiedRolePath().
func (p *XAdES111Path) CertifiedRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIGNER_ROLE, XAdES111Element_CERTIFIED_ROLES, XAdES111Element_CERTIFIED_ROLE)
}

// CertifiedRoleV2Path implements XAdESPath. Ports getCertifiedRoleV2Path().
func (p *XAdES111Path) CertifiedRoleV2Path() common.XPathQuery {
	return nil
}

// SignedDataObjectPropertiesPath implements XAdESPath. Ports getSignedDataObjectPropertiesPath().
func (p *XAdES111Path) SignedDataObjectPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES)
}

// AllDataObjectsTimestampPath implements XAdESPath. Ports getAllDataObjectsTimestampPath().
func (p *XAdES111Path) AllDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP)
}

// IndividualDataObjectsTimestampPath implements XAdESPath. Ports getIndividualDataObjectsTimestampPath().
func (p *XAdES111Path) IndividualDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES111Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP)
}

// DataObjectFormat implements XAdESPath. Ports getDataObjectFormat().
func (p *XAdES111Path) DataObjectFormat() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES111Element_DATA_OBJECT_FORMAT)
}

// DataObjectFormatMimeType implements XAdESPath. Ports getDataObjectFormatMimeType().
func (p *XAdES111Path) DataObjectFormatMimeType() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES111Element_DATA_OBJECT_FORMAT, XAdES111Element_MIME_TYPE)
}

// DataObjectFormatObjectIdentifier implements XAdESPath. Ports getDataObjectFormatObjectIdentifier().
func (p *XAdES111Path) DataObjectFormatObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES111Element_DATA_OBJECT_FORMAT, XAdES111Element_OBJECT_IDENTIFIER)
}

// CommitmentTypeIndicationPath implements XAdESPath. Ports getCommitmentTypeIndicationPath().
func (p *XAdES111Path) CommitmentTypeIndicationPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_SIGNED_PROPERTIES, XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES, XAdES111Element_COMMITMENT_TYPE_INDICATION)
}

// UnsignedPropertiesPath implements XAdESPath. Ports getUnsignedPropertiesPath().
func (p *XAdES111Path) UnsignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES)
}

// UnsignedSignaturePropertiesPath implements XAdESPath. Ports getUnsignedSignaturePropertiesPath().
func (p *XAdES111Path) UnsignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES)
}

// CounterSignaturePath implements XAdESPath. Ports getCounterSignaturePath().
func (p *XAdES111Path) CounterSignaturePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES111Element_COUNTER_SIGNATURE)
}

// AttributeRevocationRefsPath implements XAdESPath. Ports getAttributeRevocationRefsPath().
func (p *XAdES111Path) AttributeRevocationRefsPath() common.XPathQuery {
	return nil
}

// CompleteRevocationRefsPath implements XAdESPath. Ports getCompleteRevocationRefsPath().
func (p *XAdES111Path) CompleteRevocationRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES111Element_COMPLETE_REVOCATION_REFS)
}

// CompleteCertificateRefsPath implements XAdESPath. Ports getCompleteCertificateRefsPath().
func (p *XAdES111Path) CompleteCertificateRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES111Element_COMPLETE_CERTIFICATE_REFS)
}

// CompleteCertificateRefsCertPath implements XAdESPath. Ports getCompleteCertificateRefsCertPath().
func (p *XAdES111Path) CompleteCertificateRefsCertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES111Element_COMPLETE_CERTIFICATE_REFS, XAdES111Element_CERT_REFS, XAdES111Element_CERT)
}

// CompleteCertificateRefsV2Path implements XAdESPath. Ports getCompleteCertificateRefsV2Path().
func (p *XAdES111Path) CompleteCertificateRefsV2Path() common.XPathQuery {
	return nil
}

// CompleteCertificateRefsV2CertPath implements XAdESPath. Ports getCompleteCertificateRefsV2CertPath().
func (p *XAdES111Path) CompleteCertificateRefsV2CertPath() common.XPathQuery {
	return nil
}

// AttributeCertificateRefsPath implements XAdESPath. Ports getAttributeCertificateRefsPath().
func (p *XAdES111Path) AttributeCertificateRefsPath() common.XPathQuery {
	return nil
}

// AttributeCertificateRefsCertPath implements XAdESPath. Ports getAttributeCertificateRefsCertPath().
func (p *XAdES111Path) AttributeCertificateRefsCertPath() common.XPathQuery {
	return nil
}

// AttributeCertificateRefsV2Path implements XAdESPath. Ports getAttributeCertificateRefsV2Path().
func (p *XAdES111Path) AttributeCertificateRefsV2Path() common.XPathQuery {
	return nil
}

// AttributeCertificateRefsV2CertPath implements XAdESPath. Ports getAttributeCertificateRefsV2CertPath().
func (p *XAdES111Path) AttributeCertificateRefsV2CertPath() common.XPathQuery {
	return nil
}

// CertificateValuesPath implements XAdESPath. Ports getCertificateValuesPath().
func (p *XAdES111Path) CertificateValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES111Element_CERTIFICATE_VALUES)
}

// RevocationValuesPath implements XAdESPath. Ports getRevocationValuesPath().
func (p *XAdES111Path) RevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES111Element_REVOCATION_VALUES)
}

// AttributeRevocationValuesPath implements XAdESPath. Ports getAttributeRevocationValuesPath().
func (p *XAdES111Path) AttributeRevocationValuesPath() common.XPathQuery {
	return nil
}

// EncapsulatedCertificateValuesPath implements XAdESPath. Ports getEncapsulatedCertificateValuesPath().
func (p *XAdES111Path) EncapsulatedCertificateValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES111Element_CERTIFICATE_VALUES, XAdES111Element_ENCAPSULATED_X509_CERTIFICATE)
}

// AttrAuthoritiesCertValuesPath implements XAdESPath. Ports getAttrAuthoritiesCertValuesPath().
func (p *XAdES111Path) AttrAuthoritiesCertValuesPath() common.XPathQuery {
	return nil
}

// EncapsulatedAttrAuthoritiesCertValuesPath implements XAdESPath. Ports getEncapsulatedAttrAuthoritiesCertValuesPath().
func (p *XAdES111Path) EncapsulatedAttrAuthoritiesCertValuesPath() common.XPathQuery {
	return nil
}

// EncapsulatedTimeStampValidationDataCertValuesPath implements XAdESPath. Ports getEncapsulatedTimeStampValidationDataCertValuesPath().
func (p *XAdES111Path) EncapsulatedTimeStampValidationDataCertValuesPath() common.XPathQuery {
	return nil
}

// TimeStampValidationDataRevocationValuesPath implements XAdESPath. Ports getTimeStampValidationDataRevocationValuesPath().
func (p *XAdES111Path) TimeStampValidationDataRevocationValuesPath() common.XPathQuery {
	return nil
}

// AnyValidationDataPath implements XAdESPath. Ports getAnyValidationDataPath().
func (p *XAdES111Path) AnyValidationDataPath() common.XPathQuery {
	return nil
}

// EncapsulatedAnyValidationDataCertValuesPath implements XAdESPath. Ports getEncapsulatedAnyValidationDataCertValuesPath().
func (p *XAdES111Path) EncapsulatedAnyValidationDataCertValuesPath() common.XPathQuery {
	return nil
}

// AnyValidationDataRevocationValuesPath implements XAdESPath. Ports getAnyValidationDataRevocationValuesPath().
func (p *XAdES111Path) AnyValidationDataRevocationValuesPath() common.XPathQuery {
	return nil
}

// SignatureTimestampPath implements XAdESPath. Ports getSignatureTimestampPath().
func (p *XAdES111Path) SignatureTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIGNATURE_TIMESTAMP)
}

// SigAndRefsTimestampPath implements XAdESPath. Ports getSigAndRefsTimestampPath().
func (p *XAdES111Path) SigAndRefsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElement_OBJECT, XAdES111Element_QUALIFYING_PROPERTIES, XAdES111Element_UNSIGNED_PROPERTIES, XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES, XAdES111Element_SIG_AND_REFS_TIMESTAMP)
}

// SigAndRefsTimestampV2Path implements XAdESPath. Ports getSigAndRefsTimestampV2Path().
func (p *XAdES111Path) SigAndRefsTimestampV2Path() common.XPathQuery {
	return nil
}

// RefsOnlyTimestampPath implements XAdESPath. Ports getRefsOnlyTimestampPath().
func (p *XAdES111Path) RefsOnlyTimestampPath() common.XPathQuery {
	return nil
}

// RefsOnlyTimestampV2Path implements XAdESPath. Ports getRefsOnlyTimestampV2Path().
func (p *XAdES111Path) RefsOnlyTimestampV2Path() common.XPathQuery {
	return nil
}

// ArchiveTimestampPath implements XAdESPath. Ports getArchiveTimestampPath().
func (p *XAdES111Path) ArchiveTimestampPath() common.XPathQuery {
	return nil
}

// TimestampValidationDataPath implements XAdESPath. Ports getTimestampValidationDataPath().
func (p *XAdES111Path) TimestampValidationDataPath() common.XPathQuery {
	return nil
}

// SignaturePolicyStorePath implements XAdESPath. Ports getSignaturePolicyStorePath().
func (p *XAdES111Path) SignaturePolicyStorePath() common.XPathQuery {
	return nil
}

// SealingEvidenceRecordsPath implements XAdESPath. Ports getSealingEvidenceRecordsPath().
func (p *XAdES111Path) SealingEvidenceRecordsPath() common.XPathQuery {
	return nil
}

// CurrentCRLValuesChildren implements XAdESPath. Ports getCurrentCRLValuesChildren().
func (p *XAdES111Path) CurrentCRLValuesChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CRL_VALUES, XAdES111Element_ENCAPSULATED_CRL_VALUE)
}

// CurrentCRLRefsChildren implements XAdESPath. Ports getCurrentCRLRefsChildren().
func (p *XAdES111Path) CurrentCRLRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CRL_REFS, XAdES111Element_CRL_REF)
}

// CurrentCRLRefCRLIdentifier implements XAdESPath. Ports getCurrentCRLRefCRLIdentifier().
func (p *XAdES111Path) CurrentCRLRefCRLIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CRL_IDENTIFIER)
}

// CurrentCRLRefCRLIdentifierIssuer implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssuer().
func (p *XAdES111Path) CurrentCRLRefCRLIdentifierIssuer() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CRL_IDENTIFIER, XAdES111Element_ISSUER)
}

// CurrentCRLRefCRLIdentifierIssueTime implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssueTime().
func (p *XAdES111Path) CurrentCRLRefCRLIdentifierIssueTime() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CRL_IDENTIFIER, XAdES111Element_ISSUE_TIME)
}

// CurrentCRLRefCRLIdentifierNumber implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierNumber().
func (p *XAdES111Path) CurrentCRLRefCRLIdentifierNumber() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CRL_IDENTIFIER, XAdES111Element_NUMBER)
}

// CurrentOCSPValuesChildren implements XAdESPath. Ports getCurrentOCSPValuesChildren().
func (p *XAdES111Path) CurrentOCSPValuesChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_OCSP_VALUES, XAdES111Element_ENCAPSULATED_OCSP_VALUE)
}

// CurrentOCSPRefsChildren implements XAdESPath. Ports getCurrentOCSPRefsChildren().
func (p *XAdES111Path) CurrentOCSPRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_OCSP_REFS, XAdES111Element_OCSP_REF)
}

// CurrentOCSPRefResponderID implements XAdESPath. Ports getCurrentOCSPRefResponderID().
func (p *XAdES111Path) CurrentOCSPRefResponderID() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_OCSP_IDENTIFIER, XAdES111Element_RESPONDER_ID)
}

// CurrentOCSPRefResponderIDByName implements XAdESPath. Ports getCurrentOCSPRefResponderIDByName().
func (p *XAdES111Path) CurrentOCSPRefResponderIDByName() common.XPathQuery {
	return nil
}

// CurrentOCSPRefResponderIDByKey implements XAdESPath. Ports getCurrentOCSPRefResponderIDByKey().
func (p *XAdES111Path) CurrentOCSPRefResponderIDByKey() common.XPathQuery {
	return nil
}

// CurrentOCSPRefProducedAt implements XAdESPath. Ports getCurrentOCSPRefProducedAt().
func (p *XAdES111Path) CurrentOCSPRefProducedAt() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_OCSP_IDENTIFIER, XAdES111Element_PRODUCED_AT)
}

// CurrentDigestAlgAndValue implements XAdESPath. Ports getCurrentDigestAlgAndValue().
func (p *XAdES111Path) CurrentDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_DIGEST_ALG_AND_VALUE)
}

// CurrentCertRefsCertChildren implements XAdESPath. Ports getCurrentCertRefsCertChildren().
func (p *XAdES111Path) CurrentCertRefsCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CERT_REFS, XAdES111Element_CERT)
}

// CurrentCertRefs141CertChildren implements XAdESPath. Ports getCurrentCertRefs141CertChildren().
func (p *XAdES111Path) CurrentCertRefs141CertChildren() common.XPathQuery {
	return nil
}

// CurrentCertChildren implements XAdESPath. Ports getCurrentCertChildren().
func (p *XAdES111Path) CurrentCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CERT)
}

// CurrentCertDigest implements XAdESPath. Ports getCurrentCertDigest().
func (p *XAdES111Path) CurrentCertDigest() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CERT_DIGEST)
}

// CurrentEncapsulatedTimestamp implements XAdESPath. Ports getCurrentEncapsulatedTimestamp().
func (p *XAdES111Path) CurrentEncapsulatedTimestamp() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_ENCAPSULATED_TIMESTAMP)
}

// CurrentEncapsulatedCertificate implements XAdESPath. Ports getCurrentEncapsulatedCertificate().
func (p *XAdES111Path) CurrentEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_ENCAPSULATED_X509_CERTIFICATE)
}

// CurrentCertificateValuesEncapsulatedCertificate implements XAdESPath. Ports getCurrentCertificateValuesEncapsulatedCertificate().
func (p *XAdES111Path) CurrentCertificateValuesEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CERTIFICATE_VALUES, XAdES111Element_ENCAPSULATED_X509_CERTIFICATE)
}

// CurrentRevocationValuesEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedOCSPValue().
func (p *XAdES111Path) CurrentRevocationValuesEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_REVOCATION_VALUES, XAdES111Element_OCSP_VALUES, XAdES111Element_ENCAPSULATED_OCSP_VALUE)
}

// CurrentEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentEncapsulatedOCSPValue().
func (p *XAdES111Path) CurrentEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_OCSP_VALUES, XAdES111Element_ENCAPSULATED_OCSP_VALUE)
}

// CurrentRevocationValuesEncapsulatedCRLValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedCRLValue().
func (p *XAdES111Path) CurrentRevocationValuesEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_REVOCATION_VALUES, XAdES111Element_CRL_VALUES, XAdES111Element_ENCAPSULATED_CRL_VALUE)
}

// CurrentEncapsulatedCRLValue implements XAdESPath. Ports getCurrentEncapsulatedCRLValue().
func (p *XAdES111Path) CurrentEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_CRL_VALUES, XAdES111Element_ENCAPSULATED_CRL_VALUE)
}

// CurrentIssuerSerialIssuerNamePath implements XAdESPath. Ports getCurrentIssuerSerialIssuerNamePath().
func (p *XAdES111Path) CurrentIssuerSerialIssuerNamePath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_ISSUER_SERIAL, common.XMLDSigElement_X509_ISSUER_NAME)
}

// CurrentIssuerSerialSerialNumberPath implements XAdESPath. Ports getCurrentIssuerSerialSerialNumberPath().
func (p *XAdES111Path) CurrentIssuerSerialSerialNumberPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_ISSUER_SERIAL, common.XMLDSigElement_X509_SERIAL_NUMBER)
}

// CurrentIssuerSerialV2Path implements XAdESPath. Ports getCurrentIssuerSerialV2Path().
func (p *XAdES111Path) CurrentIssuerSerialV2Path() common.XPathQuery {
	return nil
}

// CurrentCommitmentIdentifierPath implements XAdESPath. Ports getCurrentCommitmentIdentifierPath().
func (p *XAdES111Path) CurrentCommitmentIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_COMMITMENT_TYPE_ID, XAdES111Element_IDENTIFIER)
}

// CurrentCommitmentDescriptionPath implements XAdESPath. Ports getCurrentCommitmentDescriptionPath().
func (p *XAdES111Path) CurrentCommitmentDescriptionPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_COMMITMENT_TYPE_ID, XAdES111Element_DESCRIPTION)
}

// CurrentCommitmentDocumentationReferencesPath implements XAdESPath. Ports getCurrentCommitmentDocumentationReferencesPath().
func (p *XAdES111Path) CurrentCommitmentDocumentationReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_COMMITMENT_TYPE_ID, XAdES111Element_DOCUMENTATION_REFERENCES)
}

// CurrentDocumentationReference implements XAdESPath. Ports getCurrentDocumentationReference().
func (p *XAdES111Path) CurrentDocumentationReference() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_DOCUMENTATION_REFERENCE)
}

// CurrentDescription implements XAdESPath. Ports getCurrentDescription().
func (p *XAdES111Path) CurrentDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_DESCRIPTION)
}

// CurrentObjectIdentifier implements XAdESPath. Ports getCurrentObjectIdentifier().
func (p *XAdES111Path) CurrentObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_OBJECT_IDENTIFIER)
}

// CurrentCommitmentObjectReferencesPath implements XAdESPath. Ports getCurrentCommitmentObjectReferencesPath().
func (p *XAdES111Path) CurrentCommitmentObjectReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_OBJECT_REFERENCE)
}

// CurrentCommitmentAllSignedDataObjectsPath implements XAdESPath. Ports getCurrentCommitmentAllSignedDataObjectsPath().
func (p *XAdES111Path) CurrentCommitmentAllSignedDataObjectsPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_ALL_SIGNED_DATA_OBJECTS)
}

// CurrentMimeType implements XAdESPath. Ports getCurrentMimeType().
func (p *XAdES111Path) CurrentMimeType() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_MIME_TYPE)
}

// CurrentEncoding implements XAdESPath. Ports getCurrentEncoding().
func (p *XAdES111Path) CurrentEncoding() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_ENCODING)
}

// CurrentSignaturePolicyId implements XAdESPath. Ports getCurrentSignaturePolicyId().
func (p *XAdES111Path) CurrentSignaturePolicyId() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_SIGNATURE_POLICY_ID, XAdES111Element_SIG_POLICY_ID, XAdES111Element_IDENTIFIER)
}

// CurrentSignaturePolicyDigestAlgAndValue implements XAdESPath. Ports getCurrentSignaturePolicyDigestAlgAndValue().
func (p *XAdES111Path) CurrentSignaturePolicyDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_SIGNATURE_POLICY_ID, XAdES111Element_SIG_POLICY_HASH)
}

// CurrentSignaturePolicySPURI implements XAdESPath. Ports getCurrentSignaturePolicySPURI().
func (p *XAdES111Path) CurrentSignaturePolicySPURI() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_SIGNATURE_POLICY_ID, XAdES111Element_SIG_POLICY_QUALIFIERS, XAdES111Element_SIG_POLICY_QUALIFIER, XAdES111Element_SP_URI)
}

// CurrentSignaturePolicySPUserNotice implements XAdESPath. Ports getCurrentSignaturePolicySPUserNotice().
func (p *XAdES111Path) CurrentSignaturePolicySPUserNotice() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_SIGNATURE_POLICY_ID, XAdES111Element_SIG_POLICY_QUALIFIERS, XAdES111Element_SIG_POLICY_QUALIFIER, XAdES111Element_SP_USER_NOTICE)
}

// CurrentSPUserNoticeNoticeRefOrganization implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefOrganization().
func (p *XAdES111Path) CurrentSPUserNoticeNoticeRefOrganization() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_NOTICE_REF, XAdES111Element_ORGANIZATION)
}

// CurrentSPUserNoticeNoticeRefNoticeNumbers implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefNoticeNumbers().
func (p *XAdES111Path) CurrentSPUserNoticeNoticeRefNoticeNumbers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_NOTICE_REF, XAdES111Element_NOTICE_NUMBERS)
}

// CurrentSPUserNoticeExplicitText implements XAdESPath. Ports getCurrentSPUserNoticeExplicitText().
func (p *XAdES111Path) CurrentSPUserNoticeExplicitText() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_EXPLICIT_TEXT)
}

// CurrentSignaturePolicySPDocSpecification implements XAdESPath. Ports getCurrentSignaturePolicySPDocSpecification().
func (p *XAdES111Path) CurrentSignaturePolicySPDocSpecification() common.XPathQuery {
	return nil
}

// CurrentSignaturePolicySPDocSpecificationIdentifier implements XAdESPath. Ports getCurrentSignaturePolicySPDocSpecificationIdentifier().
func (p *XAdES111Path) CurrentSignaturePolicySPDocSpecificationIdentifier() common.XPathQuery {
	return nil
}

// CurrentSignaturePolicyDescription implements XAdESPath. Ports getCurrentSignaturePolicyDescription().
func (p *XAdES111Path) CurrentSignaturePolicyDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_SIGNATURE_POLICY_ID, XAdES111Element_SIG_POLICY_ID, XAdES111Element_DESCRIPTION)
}

// CurrentSignaturePolicyDocumentationReferences implements XAdESPath. Ports getCurrentSignaturePolicyDocumentationReferences().
func (p *XAdES111Path) CurrentSignaturePolicyDocumentationReferences() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_SIGNATURE_POLICY_ID, XAdES111Element_SIG_POLICY_ID, XAdES111Element_DOCUMENTATION_REFERENCES)
}

// CurrentSignaturePolicyImplied implements XAdESPath. Ports getCurrentSignaturePolicyImplied().
func (p *XAdES111Path) CurrentSignaturePolicyImplied() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_SIGNATURE_POLICY_IMPLIED)
}

// CurrentSignaturePolicyTransforms implements XAdESPath. Ports getCurrentSignaturePolicyTransforms().
func (p *XAdES111Path) CurrentSignaturePolicyTransforms() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_SIGNATURE_POLICY_ID, common.XMLDSigElement_TRANSFORMS)
}

// CurrentSignaturePolicyQualifiers implements XAdESPath. Ports getCurrentSignaturePolicyQualifiers().
func (p *XAdES111Path) CurrentSignaturePolicyQualifiers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_SIGNATURE_POLICY_ID, XAdES111Element_SIG_POLICY_QUALIFIERS)
}

// CurrentInclude implements XAdESPath. Ports getCurrentInclude().
func (p *XAdES111Path) CurrentInclude() common.XPathQuery {
	return nil
}

// CurrentQualifyingPropertiesPath implements XAdESPath. Ports getCurrentQualifyingPropertiesPath().
func (p *XAdES111Path) CurrentQualifyingPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_QUALIFYING_PROPERTIES)
}

// CurrentSPDocSpecification implements XAdESPath. Ports getCurrentSPDocSpecification().
func (p *XAdES111Path) CurrentSPDocSpecification() common.XPathQuery {
	return nil
}

// CurrentIdentifier implements XAdESPath. Ports getCurrentIdentifier().
func (p *XAdES111Path) CurrentIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111Element_IDENTIFIER)
}

// CurrentSPDocSpecificationIdentifier implements XAdESPath. Ports getCurrentSPDocSpecificationIdentifier().
func (p *XAdES111Path) CurrentSPDocSpecificationIdentifier() common.XPathQuery {
	return nil
}

// CurrentSPDocSpecificationDescription implements XAdESPath. Ports getCurrentSPDocSpecificationDescription().
func (p *XAdES111Path) CurrentSPDocSpecificationDescription() common.XPathQuery {
	return nil
}

// CurrentDocumentationReferenceElements implements XAdESPath. Ports getCurrentDocumentationReferenceElements().
func (p *XAdES111Path) CurrentDocumentationReferenceElements() common.XPathQuery {
	return nil
}

// CurrentSPDocSpecificationDocumentationReferenceElements implements XAdESPath. Ports getCurrentSPDocSpecificationDocumentationReferenceElements().
func (p *XAdES111Path) CurrentSPDocSpecificationDocumentationReferenceElements() common.XPathQuery {
	return nil
}

// CurrentSignaturePolicyDocument implements XAdESPath. Ports getCurrentSignaturePolicyDocument().
func (p *XAdES111Path) CurrentSignaturePolicyDocument() common.XPathQuery {
	return nil
}

// CurrentSigPolDocLocalURI implements XAdESPath. Ports getCurrentSigPolDocLocalURI().
func (p *XAdES111Path) CurrentSigPolDocLocalURI() common.XPathQuery {
	return nil
}

var _ XAdESPath = (*XAdES111Path)(nil)
