// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades111/XAdES111Path.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES111Path holds the XAdES 111 paths.
type XAdES111Path struct{}

// NewXAdES111Path creates a XAdES111Path.
func NewXAdES111Path() *XAdES111Path {
	return &XAdES111Path{}
}

// Namespace implements XAdESPath. Ports getNamespace().
func (p *XAdES111Path) Namespace() *common.DSSNamespace {
	return XAdESNamespaceXAdES111
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
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties)
}

// SignedPropertiesPath implements XAdESPath. Ports getSignedPropertiesPath().
func (p *XAdES111Path) SignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties)
}

// SignedSignaturePropertiesPath implements XAdESPath. Ports getSignedSignaturePropertiesPath().
func (p *XAdES111Path) SignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedSignatureProperties)
}

// SigningTimePath implements XAdESPath. Ports getSigningTimePath().
func (p *XAdES111Path) SigningTimePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedSignatureProperties, XAdES111ElementSigningTime)
}

// SigningCertificatePath implements XAdESPath. Ports getSigningCertificatePath().
func (p *XAdES111Path) SigningCertificatePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedSignatureProperties, XAdES111ElementSigningCertificate)
}

// SigningCertificateChildren implements XAdESPath. Ports getSigningCertificateChildren().
func (p *XAdES111Path) SigningCertificateChildren() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedSignatureProperties, XAdES111ElementSigningCertificate, XAdES111ElementCert)
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
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedSignatureProperties, XAdES111ElementSignatureProductionPlace)
}

// SignatureProductionPlaceV2Path implements XAdESPath. Ports getSignatureProductionPlaceV2Path().
func (p *XAdES111Path) SignatureProductionPlaceV2Path() common.XPathQuery {
	return nil
}

// SignaturePolicyIdentifierPath implements XAdESPath. Ports getSignaturePolicyIdentifierPath().
func (p *XAdES111Path) SignaturePolicyIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedSignatureProperties, XAdES111ElementSignaturePolicyIdentifier)
}

// SignerRolePath implements XAdESPath. Ports getSignerRolePath().
func (p *XAdES111Path) SignerRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedSignatureProperties, XAdES111ElementSignerRole)
}

// ClaimedRolePath implements XAdESPath. Ports getClaimedRolePath().
func (p *XAdES111Path) ClaimedRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedSignatureProperties, XAdES111ElementSignerRole, XAdES111ElementClaimedRoles, XAdES111ElementClaimedRole)
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
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedSignatureProperties, XAdES111ElementSignerRole, XAdES111ElementCertifiedRoles, XAdES111ElementCertifiedRole)
}

// CertifiedRoleV2Path implements XAdESPath. Ports getCertifiedRoleV2Path().
func (p *XAdES111Path) CertifiedRoleV2Path() common.XPathQuery {
	return nil
}

// SignedDataObjectPropertiesPath implements XAdESPath. Ports getSignedDataObjectPropertiesPath().
func (p *XAdES111Path) SignedDataObjectPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedDataObjectProperties)
}

// AllDataObjectsTimestampPath implements XAdESPath. Ports getAllDataObjectsTimestampPath().
func (p *XAdES111Path) AllDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedDataObjectProperties, XAdES111ElementAllDataObjectsTimestamp)
}

// IndividualDataObjectsTimestampPath implements XAdESPath. Ports getIndividualDataObjectsTimestampPath().
func (p *XAdES111Path) IndividualDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedDataObjectProperties, XAdES111ElementIndividualDataObjectsTimestamp)
}

// DataObjectFormat implements XAdESPath. Ports getDataObjectFormat().
func (p *XAdES111Path) DataObjectFormat() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedDataObjectProperties, XAdES111ElementDataObjectFormat)
}

// DataObjectFormatMimeType implements XAdESPath. Ports getDataObjectFormatMimeType().
func (p *XAdES111Path) DataObjectFormatMimeType() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedDataObjectProperties, XAdES111ElementDataObjectFormat, XAdES111ElementMIMEType)
}

// DataObjectFormatObjectIdentifier implements XAdESPath. Ports getDataObjectFormatObjectIdentifier().
func (p *XAdES111Path) DataObjectFormatObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedDataObjectProperties, XAdES111ElementDataObjectFormat, XAdES111ElementObjectIdentifier)
}

// CommitmentTypeIndicationPath implements XAdESPath. Ports getCommitmentTypeIndicationPath().
func (p *XAdES111Path) CommitmentTypeIndicationPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementSignedProperties, XAdES111ElementSignedDataObjectProperties, XAdES111ElementCommitmentTypeIndication)
}

// UnsignedPropertiesPath implements XAdESPath. Ports getUnsignedPropertiesPath().
func (p *XAdES111Path) UnsignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties)
}

// UnsignedSignaturePropertiesPath implements XAdESPath. Ports getUnsignedSignaturePropertiesPath().
func (p *XAdES111Path) UnsignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties)
}

// CounterSignaturePath implements XAdESPath. Ports getCounterSignaturePath().
func (p *XAdES111Path) CounterSignaturePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties, XAdES111ElementCounterSignature)
}

// AttributeRevocationRefsPath implements XAdESPath. Ports getAttributeRevocationRefsPath().
func (p *XAdES111Path) AttributeRevocationRefsPath() common.XPathQuery {
	return nil
}

// CompleteRevocationRefsPath implements XAdESPath. Ports getCompleteRevocationRefsPath().
func (p *XAdES111Path) CompleteRevocationRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties, XAdES111ElementCompleteRevocationRefs)
}

// CompleteCertificateRefsPath implements XAdESPath. Ports getCompleteCertificateRefsPath().
func (p *XAdES111Path) CompleteCertificateRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties, XAdES111ElementCompleteCertificateRefs)
}

// CompleteCertificateRefsCertPath implements XAdESPath. Ports getCompleteCertificateRefsCertPath().
func (p *XAdES111Path) CompleteCertificateRefsCertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties, XAdES111ElementCompleteCertificateRefs, XAdES111ElementCertRefs, XAdES111ElementCert)
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
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties, XAdES111ElementCertificateValues)
}

// RevocationValuesPath implements XAdESPath. Ports getRevocationValuesPath().
func (p *XAdES111Path) RevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties, XAdES111ElementRevocationValues)
}

// AttributeRevocationValuesPath implements XAdESPath. Ports getAttributeRevocationValuesPath().
func (p *XAdES111Path) AttributeRevocationValuesPath() common.XPathQuery {
	return nil
}

// EncapsulatedCertificateValuesPath implements XAdESPath. Ports getEncapsulatedCertificateValuesPath().
func (p *XAdES111Path) EncapsulatedCertificateValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties, XAdES111ElementCertificateValues, XAdES111ElementEncapsulatedX509Certificate)
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
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties, XAdES111ElementSignatureTimestamp)
}

// SigAndRefsTimestampPath implements XAdESPath. Ports getSigAndRefsTimestampPath().
func (p *XAdES111Path) SigAndRefsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES111ElementQualifyingProperties, XAdES111ElementUnsignedProperties, XAdES111ElementUnsignedSignatureProperties, XAdES111ElementSigAndRefsTimestamp)
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
	return common.FromCurrentPosition(XAdES111ElementCRLValues, XAdES111ElementEncapsulatedCRLValue)
}

// CurrentCRLRefsChildren implements XAdESPath. Ports getCurrentCRLRefsChildren().
func (p *XAdES111Path) CurrentCRLRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCRLRefs, XAdES111ElementCRLRef)
}

// CurrentCRLRefCRLIdentifier implements XAdESPath. Ports getCurrentCRLRefCRLIdentifier().
func (p *XAdES111Path) CurrentCRLRefCRLIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCRLIdentifier)
}

// CurrentCRLRefCRLIdentifierIssuer implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssuer().
func (p *XAdES111Path) CurrentCRLRefCRLIdentifierIssuer() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCRLIdentifier, XAdES111ElementIssuer)
}

// CurrentCRLRefCRLIdentifierIssueTime implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssueTime().
func (p *XAdES111Path) CurrentCRLRefCRLIdentifierIssueTime() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCRLIdentifier, XAdES111ElementIssueTime)
}

// CurrentCRLRefCRLIdentifierNumber implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierNumber().
func (p *XAdES111Path) CurrentCRLRefCRLIdentifierNumber() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCRLIdentifier, XAdES111ElementNumber)
}

// CurrentOCSPValuesChildren implements XAdESPath. Ports getCurrentOCSPValuesChildren().
func (p *XAdES111Path) CurrentOCSPValuesChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementOCSPValues, XAdES111ElementEncapsulatedOCSPValue)
}

// CurrentOCSPRefsChildren implements XAdESPath. Ports getCurrentOCSPRefsChildren().
func (p *XAdES111Path) CurrentOCSPRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementOCSPRefs, XAdES111ElementOCSPRef)
}

// CurrentOCSPRefResponderID implements XAdESPath. Ports getCurrentOCSPRefResponderID().
func (p *XAdES111Path) CurrentOCSPRefResponderID() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementOCSPIdentifier, XAdES111ElementResponderID)
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
	return common.FromCurrentPosition(XAdES111ElementOCSPIdentifier, XAdES111ElementProducedAt)
}

// CurrentDigestAlgAndValue implements XAdESPath. Ports getCurrentDigestAlgAndValue().
func (p *XAdES111Path) CurrentDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementDigestAlgAndValue)
}

// CurrentCertRefsCertChildren implements XAdESPath. Ports getCurrentCertRefsCertChildren().
func (p *XAdES111Path) CurrentCertRefsCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCertRefs, XAdES111ElementCert)
}

// CurrentCertRefs141CertChildren implements XAdESPath. Ports getCurrentCertRefs141CertChildren().
func (p *XAdES111Path) CurrentCertRefs141CertChildren() common.XPathQuery {
	return nil
}

// CurrentCertChildren implements XAdESPath. Ports getCurrentCertChildren().
func (p *XAdES111Path) CurrentCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCert)
}

// CurrentCertDigest implements XAdESPath. Ports getCurrentCertDigest().
func (p *XAdES111Path) CurrentCertDigest() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCertDigest)
}

// CurrentEncapsulatedTimestamp implements XAdESPath. Ports getCurrentEncapsulatedTimestamp().
func (p *XAdES111Path) CurrentEncapsulatedTimestamp() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementEncapsulatedTimestamp)
}

// CurrentEncapsulatedCertificate implements XAdESPath. Ports getCurrentEncapsulatedCertificate().
func (p *XAdES111Path) CurrentEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementEncapsulatedX509Certificate)
}

// CurrentCertificateValuesEncapsulatedCertificate implements XAdESPath. Ports getCurrentCertificateValuesEncapsulatedCertificate().
func (p *XAdES111Path) CurrentCertificateValuesEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCertificateValues, XAdES111ElementEncapsulatedX509Certificate)
}

// CurrentRevocationValuesEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedOCSPValue().
func (p *XAdES111Path) CurrentRevocationValuesEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementRevocationValues, XAdES111ElementOCSPValues, XAdES111ElementEncapsulatedOCSPValue)
}

// CurrentEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentEncapsulatedOCSPValue().
func (p *XAdES111Path) CurrentEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementOCSPValues, XAdES111ElementEncapsulatedOCSPValue)
}

// CurrentRevocationValuesEncapsulatedCRLValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedCRLValue().
func (p *XAdES111Path) CurrentRevocationValuesEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementRevocationValues, XAdES111ElementCRLValues, XAdES111ElementEncapsulatedCRLValue)
}

// CurrentEncapsulatedCRLValue implements XAdESPath. Ports getCurrentEncapsulatedCRLValue().
func (p *XAdES111Path) CurrentEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCRLValues, XAdES111ElementEncapsulatedCRLValue)
}

// CurrentIssuerSerialIssuerNamePath implements XAdESPath. Ports getCurrentIssuerSerialIssuerNamePath().
func (p *XAdES111Path) CurrentIssuerSerialIssuerNamePath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementIssuerSerial, common.XMLDSigElementX509IssuerName)
}

// CurrentIssuerSerialSerialNumberPath implements XAdESPath. Ports getCurrentIssuerSerialSerialNumberPath().
func (p *XAdES111Path) CurrentIssuerSerialSerialNumberPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementIssuerSerial, common.XMLDSigElementX509SerialNumber)
}

// CurrentIssuerSerialV2Path implements XAdESPath. Ports getCurrentIssuerSerialV2Path().
func (p *XAdES111Path) CurrentIssuerSerialV2Path() common.XPathQuery {
	return nil
}

// CurrentCommitmentIdentifierPath implements XAdESPath. Ports getCurrentCommitmentIdentifierPath().
func (p *XAdES111Path) CurrentCommitmentIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCommitmentTypeID, XAdES111ElementIdentifier)
}

// CurrentCommitmentDescriptionPath implements XAdESPath. Ports getCurrentCommitmentDescriptionPath().
func (p *XAdES111Path) CurrentCommitmentDescriptionPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCommitmentTypeID, XAdES111ElementDescription)
}

// CurrentCommitmentDocumentationReferencesPath implements XAdESPath. Ports getCurrentCommitmentDocumentationReferencesPath().
func (p *XAdES111Path) CurrentCommitmentDocumentationReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementCommitmentTypeID, XAdES111ElementDocumentationReferences)
}

// CurrentDocumentationReference implements XAdESPath. Ports getCurrentDocumentationReference().
func (p *XAdES111Path) CurrentDocumentationReference() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementDocumentationReference)
}

// CurrentDescription implements XAdESPath. Ports getCurrentDescription().
func (p *XAdES111Path) CurrentDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementDescription)
}

// CurrentObjectIdentifier implements XAdESPath. Ports getCurrentObjectIdentifier().
func (p *XAdES111Path) CurrentObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementObjectIdentifier)
}

// CurrentCommitmentObjectReferencesPath implements XAdESPath. Ports getCurrentCommitmentObjectReferencesPath().
func (p *XAdES111Path) CurrentCommitmentObjectReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementObjectReference)
}

// CurrentCommitmentAllSignedDataObjectsPath implements XAdESPath. Ports getCurrentCommitmentAllSignedDataObjectsPath().
func (p *XAdES111Path) CurrentCommitmentAllSignedDataObjectsPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementAllSignedDataObjects)
}

// CurrentMimeType implements XAdESPath. Ports getCurrentMimeType().
func (p *XAdES111Path) CurrentMimeType() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementMIMEType)
}

// CurrentEncoding implements XAdESPath. Ports getCurrentEncoding().
func (p *XAdES111Path) CurrentEncoding() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementEncoding)
}

// CurrentSignaturePolicyId implements XAdESPath. Ports getCurrentSignaturePolicyId().
func (p *XAdES111Path) CurrentSignaturePolicyId() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementSignaturePolicyID, XAdES111ElementSigPolicyID, XAdES111ElementIdentifier)
}

// CurrentSignaturePolicyDigestAlgAndValue implements XAdESPath. Ports getCurrentSignaturePolicyDigestAlgAndValue().
func (p *XAdES111Path) CurrentSignaturePolicyDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementSignaturePolicyID, XAdES111ElementSigPolicyHash)
}

// CurrentSignaturePolicySPURI implements XAdESPath. Ports getCurrentSignaturePolicySPURI().
func (p *XAdES111Path) CurrentSignaturePolicySPURI() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementSignaturePolicyID, XAdES111ElementSigPolicyQualifiers, XAdES111ElementSigPolicyQualifier, XAdES111ElementSPURI)
}

// CurrentSignaturePolicySPUserNotice implements XAdESPath. Ports getCurrentSignaturePolicySPUserNotice().
func (p *XAdES111Path) CurrentSignaturePolicySPUserNotice() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementSignaturePolicyID, XAdES111ElementSigPolicyQualifiers, XAdES111ElementSigPolicyQualifier, XAdES111ElementSPUserNotice)
}

// CurrentSPUserNoticeNoticeRefOrganization implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefOrganization().
func (p *XAdES111Path) CurrentSPUserNoticeNoticeRefOrganization() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementNoticeRef, XAdES111ElementOrganization)
}

// CurrentSPUserNoticeNoticeRefNoticeNumbers implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefNoticeNumbers().
func (p *XAdES111Path) CurrentSPUserNoticeNoticeRefNoticeNumbers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementNoticeRef, XAdES111ElementNoticeNumbers)
}

// CurrentSPUserNoticeExplicitText implements XAdESPath. Ports getCurrentSPUserNoticeExplicitText().
func (p *XAdES111Path) CurrentSPUserNoticeExplicitText() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementExplicitText)
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
	return common.FromCurrentPosition(XAdES111ElementSignaturePolicyID, XAdES111ElementSigPolicyID, XAdES111ElementDescription)
}

// CurrentSignaturePolicyDocumentationReferences implements XAdESPath. Ports getCurrentSignaturePolicyDocumentationReferences().
func (p *XAdES111Path) CurrentSignaturePolicyDocumentationReferences() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementSignaturePolicyID, XAdES111ElementSigPolicyID, XAdES111ElementDocumentationReferences)
}

// CurrentSignaturePolicyImplied implements XAdESPath. Ports getCurrentSignaturePolicyImplied().
func (p *XAdES111Path) CurrentSignaturePolicyImplied() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementSignaturePolicyImplied)
}

// CurrentSignaturePolicyTransforms implements XAdESPath. Ports getCurrentSignaturePolicyTransforms().
func (p *XAdES111Path) CurrentSignaturePolicyTransforms() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementSignaturePolicyID, common.XMLDSigElementTransforms)
}

// CurrentSignaturePolicyQualifiers implements XAdESPath. Ports getCurrentSignaturePolicyQualifiers().
func (p *XAdES111Path) CurrentSignaturePolicyQualifiers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementSignaturePolicyID, XAdES111ElementSigPolicyQualifiers)
}

// CurrentInclude implements XAdESPath. Ports getCurrentInclude().
func (p *XAdES111Path) CurrentInclude() common.XPathQuery {
	return nil
}

// CurrentQualifyingPropertiesPath implements XAdESPath. Ports getCurrentQualifyingPropertiesPath().
func (p *XAdES111Path) CurrentQualifyingPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementQualifyingProperties)
}

// CurrentSPDocSpecification implements XAdESPath. Ports getCurrentSPDocSpecification().
func (p *XAdES111Path) CurrentSPDocSpecification() common.XPathQuery {
	return nil
}

// CurrentIdentifier implements XAdESPath. Ports getCurrentIdentifier().
func (p *XAdES111Path) CurrentIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES111ElementIdentifier)
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
