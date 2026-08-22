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
	return XAdESNamespaceXAdES132
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
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties)
}

// SignedPropertiesPath implements XAdESPath. Ports getSignedPropertiesPath().
func (p *XAdES132Path) SignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties)
}

// SignedSignaturePropertiesPath implements XAdESPath. Ports getSignedSignaturePropertiesPath().
func (p *XAdES132Path) SignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties)
}

// SigningTimePath implements XAdESPath. Ports getSigningTimePath().
func (p *XAdES132Path) SigningTimePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSigningTime)
}

// SigningCertificatePath implements XAdESPath. Ports getSigningCertificatePath().
func (p *XAdES132Path) SigningCertificatePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSigningCertificate)
}

// SigningCertificateChildren implements XAdESPath. Ports getSigningCertificateChildren().
func (p *XAdES132Path) SigningCertificateChildren() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSigningCertificate, XAdES132ElementCert)
}

// SigningCertificateV2Path implements XAdESPath. Ports getSigningCertificateV2Path().
func (p *XAdES132Path) SigningCertificateV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSigningCertificateV2)
}

// SigningCertificateV2Children implements XAdESPath. Ports getSigningCertificateV2Children().
func (p *XAdES132Path) SigningCertificateV2Children() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSigningCertificateV2, XAdES132ElementCert)
}

// SignatureProductionPlacePath implements XAdESPath. Ports getSignatureProductionPlacePath().
func (p *XAdES132Path) SignatureProductionPlacePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignatureProductionPlace)
}

// SignatureProductionPlaceV2Path implements XAdESPath. Ports getSignatureProductionPlaceV2Path().
func (p *XAdES132Path) SignatureProductionPlaceV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignatureProductionPlaceV2)
}

// SignaturePolicyIdentifierPath implements XAdESPath. Ports getSignaturePolicyIdentifierPath().
func (p *XAdES132Path) SignaturePolicyIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignaturePolicyIdentifier)
}

// SignerRolePath implements XAdESPath. Ports getSignerRolePath().
func (p *XAdES132Path) SignerRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignerRole)
}

// ClaimedRolePath implements XAdESPath. Ports getClaimedRolePath().
func (p *XAdES132Path) ClaimedRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignerRole, XAdES132ElementClaimedRoles, XAdES132ElementClaimedRole)
}

// SignedAssertionPath implements XAdESPath. Ports getSignedAssertionPath().
func (p *XAdES132Path) SignedAssertionPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignerRoleV2, XAdES132ElementSignedAssertions, XAdES132ElementSignedAssertion)
}

// SignerRoleV2Path implements XAdESPath. Ports getSignerRoleV2Path().
func (p *XAdES132Path) SignerRoleV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignerRoleV2)
}

// ClaimedRoleV2Path implements XAdESPath. Ports getClaimedRoleV2Path().
func (p *XAdES132Path) ClaimedRoleV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignerRoleV2, XAdES132ElementClaimedRoles, XAdES132ElementClaimedRole)
}

// CertifiedRolePath implements XAdESPath. Ports getCertifiedRolePath().
func (p *XAdES132Path) CertifiedRolePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignerRole, XAdES132ElementCertifiedRoles, XAdES132ElementCertifiedRole)
}

// CertifiedRoleV2Path implements XAdESPath. Ports getCertifiedRoleV2Path().
func (p *XAdES132Path) CertifiedRoleV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedSignatureProperties, XAdES132ElementSignerRoleV2, XAdES132ElementCertifiedRolesV2, XAdES132ElementCertifiedRole)
}

// SignedDataObjectPropertiesPath implements XAdESPath. Ports getSignedDataObjectPropertiesPath().
func (p *XAdES132Path) SignedDataObjectPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedDataObjectProperties)
}

// AllDataObjectsTimestampPath implements XAdESPath. Ports getAllDataObjectsTimestampPath().
func (p *XAdES132Path) AllDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedDataObjectProperties, XAdES132ElementAllDataObjectsTimestamp)
}

// IndividualDataObjectsTimestampPath implements XAdESPath. Ports getIndividualDataObjectsTimestampPath().
func (p *XAdES132Path) IndividualDataObjectsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedDataObjectProperties, XAdES132ElementIndividualDataObjectsTimestamp)
}

// DataObjectFormat implements XAdESPath. Ports getDataObjectFormat().
func (p *XAdES132Path) DataObjectFormat() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedDataObjectProperties, XAdES132ElementDataObjectFormat)
}

// DataObjectFormatMimeType implements XAdESPath. Ports getDataObjectFormatMimeType().
func (p *XAdES132Path) DataObjectFormatMimeType() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedDataObjectProperties, XAdES132ElementDataObjectFormat, XAdES132ElementMIMEType)
}

// DataObjectFormatObjectIdentifier implements XAdESPath. Ports getDataObjectFormatObjectIdentifier().
func (p *XAdES132Path) DataObjectFormatObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedDataObjectProperties, XAdES132ElementDataObjectFormat, XAdES132ElementObjectIdentifier)
}

// CommitmentTypeIndicationPath implements XAdESPath. Ports getCommitmentTypeIndicationPath().
func (p *XAdES132Path) CommitmentTypeIndicationPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementSignedProperties, XAdES132ElementSignedDataObjectProperties, XAdES132ElementCommitmentTypeIndication)
}

// UnsignedPropertiesPath implements XAdESPath. Ports getUnsignedPropertiesPath().
func (p *XAdES132Path) UnsignedPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties)
}

// UnsignedSignaturePropertiesPath implements XAdESPath. Ports getUnsignedSignaturePropertiesPath().
func (p *XAdES132Path) UnsignedSignaturePropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties)
}

// CounterSignaturePath implements XAdESPath. Ports getCounterSignaturePath().
func (p *XAdES132Path) CounterSignaturePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementCounterSignature)
}

// AttributeRevocationRefsPath implements XAdESPath. Ports getAttributeRevocationRefsPath().
func (p *XAdES132Path) AttributeRevocationRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementAttributeRevocationRefs)
}

// CompleteRevocationRefsPath implements XAdESPath. Ports getCompleteRevocationRefsPath().
func (p *XAdES132Path) CompleteRevocationRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementCompleteRevocationRefs)
}

// CompleteCertificateRefsPath implements XAdESPath. Ports getCompleteCertificateRefsPath().
func (p *XAdES132Path) CompleteCertificateRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementCompleteCertificateRefs)
}

// CompleteCertificateRefsCertPath implements XAdESPath. Ports getCompleteCertificateRefsCertPath().
func (p *XAdES132Path) CompleteCertificateRefsCertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementCompleteCertificateRefs, XAdES132ElementCertRefs, XAdES132ElementCert)
}

// CompleteCertificateRefsV2Path implements XAdESPath. Ports getCompleteCertificateRefsV2Path().
func (p *XAdES132Path) CompleteCertificateRefsV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementCompleteCertificateRefsV2)
}

// CompleteCertificateRefsV2CertPath implements XAdESPath. Ports getCompleteCertificateRefsV2CertPath().
func (p *XAdES132Path) CompleteCertificateRefsV2CertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementCompleteCertificateRefsV2, XAdES141ElementCertRefs, XAdES132ElementCert)
}

// AttributeCertificateRefsPath implements XAdESPath. Ports getAttributeCertificateRefsPath().
func (p *XAdES132Path) AttributeCertificateRefsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementAttributeCertificateRefs)
}

// AttributeCertificateRefsCertPath implements XAdESPath. Ports getAttributeCertificateRefsCertPath().
func (p *XAdES132Path) AttributeCertificateRefsCertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementAttributeCertificateRefs, XAdES132ElementCertRefs, XAdES132ElementCert)
}

// AttributeCertificateRefsV2Path implements XAdESPath. Ports getAttributeCertificateRefsV2Path().
func (p *XAdES132Path) AttributeCertificateRefsV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementAttributeCertificateRefsV2)
}

// AttributeCertificateRefsV2CertPath implements XAdESPath. Ports getAttributeCertificateRefsV2CertPath().
func (p *XAdES132Path) AttributeCertificateRefsV2CertPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementAttributeCertificateRefsV2, XAdES141ElementCertRefs, XAdES132ElementCert)
}

// CertificateValuesPath implements XAdESPath. Ports getCertificateValuesPath().
func (p *XAdES132Path) CertificateValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementCertificateValues)
}

// RevocationValuesPath implements XAdESPath. Ports getRevocationValuesPath().
func (p *XAdES132Path) RevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementRevocationValues)
}

// AttributeRevocationValuesPath implements XAdESPath. Ports getAttributeRevocationValuesPath().
func (p *XAdES132Path) AttributeRevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementAttributeRevocationValues)
}

// EncapsulatedCertificateValuesPath implements XAdESPath. Ports getEncapsulatedCertificateValuesPath().
func (p *XAdES132Path) EncapsulatedCertificateValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementCertificateValues, XAdES132ElementEncapsulatedX509Certificate)
}

// AttrAuthoritiesCertValuesPath implements XAdESPath. Ports getAttrAuthoritiesCertValuesPath().
func (p *XAdES132Path) AttrAuthoritiesCertValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementAttrAuthoritiesCertValues)
}

// EncapsulatedAttrAuthoritiesCertValuesPath implements XAdESPath. Ports getEncapsulatedAttrAuthoritiesCertValuesPath().
func (p *XAdES132Path) EncapsulatedAttrAuthoritiesCertValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementAttrAuthoritiesCertValues, XAdES132ElementEncapsulatedX509Certificate)
}

// EncapsulatedTimeStampValidationDataCertValuesPath implements XAdESPath. Ports getEncapsulatedTimeStampValidationDataCertValuesPath().
func (p *XAdES132Path) EncapsulatedTimeStampValidationDataCertValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementTimestampValidationData, XAdES132ElementCertificateValues, XAdES132ElementEncapsulatedX509Certificate)
}

// TimeStampValidationDataRevocationValuesPath implements XAdESPath. Ports getTimeStampValidationDataRevocationValuesPath().
func (p *XAdES132Path) TimeStampValidationDataRevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementTimestampValidationData, XAdES132ElementRevocationValues)
}

// AnyValidationDataPath implements XAdESPath. Ports getAnyValidationDataPath().
func (p *XAdES132Path) AnyValidationDataPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementAnyValidationData)
}

// EncapsulatedAnyValidationDataCertValuesPath implements XAdESPath. Ports getEncapsulatedAnyValidationDataCertValuesPath().
func (p *XAdES132Path) EncapsulatedAnyValidationDataCertValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementAnyValidationData, XAdES132ElementCertificateValues, XAdES132ElementEncapsulatedX509Certificate)
}

// AnyValidationDataRevocationValuesPath implements XAdESPath. Ports getAnyValidationDataRevocationValuesPath().
func (p *XAdES132Path) AnyValidationDataRevocationValuesPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementAnyValidationData, XAdES132ElementRevocationValues)
}

// SignatureTimestampPath implements XAdESPath. Ports getSignatureTimestampPath().
func (p *XAdES132Path) SignatureTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementSignatureTimestamp)
}

// SigAndRefsTimestampPath implements XAdESPath. Ports getSigAndRefsTimestampPath().
func (p *XAdES132Path) SigAndRefsTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementSigAndRefsTimestamp)
}

// SigAndRefsTimestampV2Path implements XAdESPath. Ports getSigAndRefsTimestampV2Path().
func (p *XAdES132Path) SigAndRefsTimestampV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementSigAndRefsTimestampV2)
}

// RefsOnlyTimestampPath implements XAdESPath. Ports getRefsOnlyTimestampPath().
func (p *XAdES132Path) RefsOnlyTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES132ElementRefsOnlyTimestamp)
}

// RefsOnlyTimestampV2Path implements XAdESPath. Ports getRefsOnlyTimestampV2Path().
func (p *XAdES132Path) RefsOnlyTimestampV2Path() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementRefsOnlyTimestampV2)
}

// ArchiveTimestampPath implements XAdESPath. Ports getArchiveTimestampPath().
func (p *XAdES132Path) ArchiveTimestampPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementArchiveTimestamp)
}

// TimestampValidationDataPath implements XAdESPath. Ports getTimestampValidationDataPath().
func (p *XAdES132Path) TimestampValidationDataPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementTimestampValidationData)
}

// SignaturePolicyStorePath implements XAdESPath. Ports getSignaturePolicyStorePath().
func (p *XAdES132Path) SignaturePolicyStorePath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdES141ElementSignaturePolicyStore)
}

// SealingEvidenceRecordsPath implements XAdESPath. Ports getSealingEvidenceRecordsPath().
func (p *XAdES132Path) SealingEvidenceRecordsPath() common.XPathQuery {
	return common.FromCurrentPosition(common.XMLDSigElementObject, XAdES132ElementQualifyingProperties, XAdES132ElementUnsignedProperties, XAdES132ElementUnsignedSignatureProperties, XAdESEvidencerecordNamespaceElementSealingEvidenceRecords)
}

// CurrentCRLValuesChildren implements XAdESPath. Ports getCurrentCRLValuesChildren().
func (p *XAdES132Path) CurrentCRLValuesChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCRLValues, XAdES132ElementEncapsulatedCRLValue)
}

// CurrentCRLRefsChildren implements XAdESPath. Ports getCurrentCRLRefsChildren().
func (p *XAdES132Path) CurrentCRLRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCRLRefs, XAdES132ElementCRLRef)
}

// CurrentCRLRefCRLIdentifier implements XAdESPath. Ports getCurrentCRLRefCRLIdentifier().
func (p *XAdES132Path) CurrentCRLRefCRLIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCRLIdentifier)
}

// CurrentCRLRefCRLIdentifierIssuer implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssuer().
func (p *XAdES132Path) CurrentCRLRefCRLIdentifierIssuer() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCRLIdentifier, XAdES132ElementIssuer)
}

// CurrentCRLRefCRLIdentifierIssueTime implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierIssueTime().
func (p *XAdES132Path) CurrentCRLRefCRLIdentifierIssueTime() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCRLIdentifier, XAdES132ElementIssueTime)
}

// CurrentCRLRefCRLIdentifierNumber implements XAdESPath. Ports getCurrentCRLRefCRLIdentifierNumber().
func (p *XAdES132Path) CurrentCRLRefCRLIdentifierNumber() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCRLIdentifier, XAdES132ElementNumber)
}

// CurrentOCSPValuesChildren implements XAdESPath. Ports getCurrentOCSPValuesChildren().
func (p *XAdES132Path) CurrentOCSPValuesChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementOCSPValues, XAdES132ElementEncapsulatedOCSPValue)
}

// CurrentOCSPRefsChildren implements XAdESPath. Ports getCurrentOCSPRefsChildren().
func (p *XAdES132Path) CurrentOCSPRefsChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementOCSPRefs, XAdES132ElementOCSPRef)
}

// CurrentOCSPRefResponderID implements XAdESPath. Ports getCurrentOCSPRefResponderID().
func (p *XAdES132Path) CurrentOCSPRefResponderID() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementOCSPIdentifier, XAdES132ElementResponderID)
}

// CurrentOCSPRefResponderIDByName implements XAdESPath. Ports getCurrentOCSPRefResponderIDByName().
func (p *XAdES132Path) CurrentOCSPRefResponderIDByName() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementOCSPIdentifier, XAdES132ElementResponderID, XAdES132ElementByName)
}

// CurrentOCSPRefResponderIDByKey implements XAdESPath. Ports getCurrentOCSPRefResponderIDByKey().
func (p *XAdES132Path) CurrentOCSPRefResponderIDByKey() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementOCSPIdentifier, XAdES132ElementResponderID, XAdES132ElementByKey)
}

// CurrentOCSPRefProducedAt implements XAdESPath. Ports getCurrentOCSPRefProducedAt().
func (p *XAdES132Path) CurrentOCSPRefProducedAt() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementOCSPIdentifier, XAdES132ElementProducedAt)
}

// CurrentDigestAlgAndValue implements XAdESPath. Ports getCurrentDigestAlgAndValue().
func (p *XAdES132Path) CurrentDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementDigestAlgAndValue)
}

// CurrentCertRefsCertChildren implements XAdESPath. Ports getCurrentCertRefsCertChildren().
func (p *XAdES132Path) CurrentCertRefsCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCertRefs, XAdES132ElementCert)
}

// CurrentCertRefs141CertChildren implements XAdESPath. Ports getCurrentCertRefs141CertChildren().
func (p *XAdES132Path) CurrentCertRefs141CertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141ElementCertRefs, XAdES132ElementCert)
}

// CurrentCertChildren implements XAdESPath. Ports getCurrentCertChildren().
func (p *XAdES132Path) CurrentCertChildren() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCert)
}

// CurrentCertDigest implements XAdESPath. Ports getCurrentCertDigest().
func (p *XAdES132Path) CurrentCertDigest() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCertDigest)
}

// CurrentEncapsulatedTimestamp implements XAdESPath. Ports getCurrentEncapsulatedTimestamp().
func (p *XAdES132Path) CurrentEncapsulatedTimestamp() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementEncapsulatedTimestamp)
}

// CurrentEncapsulatedCertificate implements XAdESPath. Ports getCurrentEncapsulatedCertificate().
func (p *XAdES132Path) CurrentEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementEncapsulatedX509Certificate)
}

// CurrentCertificateValuesEncapsulatedCertificate implements XAdESPath. Ports getCurrentCertificateValuesEncapsulatedCertificate().
func (p *XAdES132Path) CurrentCertificateValuesEncapsulatedCertificate() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCertificateValues, XAdES132ElementEncapsulatedX509Certificate)
}

// CurrentRevocationValuesEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedOCSPValue().
func (p *XAdES132Path) CurrentRevocationValuesEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementRevocationValues, XAdES132ElementOCSPValues, XAdES132ElementEncapsulatedOCSPValue)
}

// CurrentEncapsulatedOCSPValue implements XAdESPath. Ports getCurrentEncapsulatedOCSPValue().
func (p *XAdES132Path) CurrentEncapsulatedOCSPValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementOCSPValues, XAdES132ElementEncapsulatedOCSPValue)
}

// CurrentRevocationValuesEncapsulatedCRLValue implements XAdESPath. Ports getCurrentRevocationValuesEncapsulatedCRLValue().
func (p *XAdES132Path) CurrentRevocationValuesEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementRevocationValues, XAdES132ElementCRLValues, XAdES132ElementEncapsulatedCRLValue)
}

// CurrentEncapsulatedCRLValue implements XAdESPath. Ports getCurrentEncapsulatedCRLValue().
func (p *XAdES132Path) CurrentEncapsulatedCRLValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCRLValues, XAdES132ElementEncapsulatedCRLValue)
}

// CurrentIssuerSerialIssuerNamePath implements XAdESPath. Ports getCurrentIssuerSerialIssuerNamePath().
func (p *XAdES132Path) CurrentIssuerSerialIssuerNamePath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementIssuerSerial, common.XMLDSigElementX509IssuerName)
}

// CurrentIssuerSerialSerialNumberPath implements XAdESPath. Ports getCurrentIssuerSerialSerialNumberPath().
func (p *XAdES132Path) CurrentIssuerSerialSerialNumberPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementIssuerSerial, common.XMLDSigElementX509SerialNumber)
}

// CurrentIssuerSerialV2Path implements XAdESPath. Ports getCurrentIssuerSerialV2Path().
func (p *XAdES132Path) CurrentIssuerSerialV2Path() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementIssuerSerialV2)
}

// CurrentCommitmentIdentifierPath implements XAdESPath. Ports getCurrentCommitmentIdentifierPath().
func (p *XAdES132Path) CurrentCommitmentIdentifierPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCommitmentTypeID, XAdES132ElementIdentifier)
}

// CurrentCommitmentDescriptionPath implements XAdESPath. Ports getCurrentCommitmentDescriptionPath().
func (p *XAdES132Path) CurrentCommitmentDescriptionPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCommitmentTypeID, XAdES132ElementDescription)
}

// CurrentCommitmentDocumentationReferencesPath implements XAdESPath. Ports getCurrentCommitmentDocumentationReferencesPath().
func (p *XAdES132Path) CurrentCommitmentDocumentationReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementCommitmentTypeID, XAdES132ElementDocumentationReferences)
}

// CurrentDocumentationReference implements XAdESPath. Ports getCurrentDocumentationReference().
func (p *XAdES132Path) CurrentDocumentationReference() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementDocumentationReference)
}

// CurrentDescription implements XAdESPath. Ports getCurrentDescription().
func (p *XAdES132Path) CurrentDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementDescription)
}

// CurrentObjectIdentifier implements XAdESPath. Ports getCurrentObjectIdentifier().
func (p *XAdES132Path) CurrentObjectIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementObjectIdentifier)
}

// CurrentCommitmentObjectReferencesPath implements XAdESPath. Ports getCurrentCommitmentObjectReferencesPath().
func (p *XAdES132Path) CurrentCommitmentObjectReferencesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementObjectReference)
}

// CurrentCommitmentAllSignedDataObjectsPath implements XAdESPath. Ports getCurrentCommitmentAllSignedDataObjectsPath().
func (p *XAdES132Path) CurrentCommitmentAllSignedDataObjectsPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementAllSignedDataObjects)
}

// CurrentMimeType implements XAdESPath. Ports getCurrentMimeType().
func (p *XAdES132Path) CurrentMimeType() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementMIMEType)
}

// CurrentEncoding implements XAdESPath. Ports getCurrentEncoding().
func (p *XAdES132Path) CurrentEncoding() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementEncoding)
}

// CurrentSignaturePolicyId implements XAdESPath. Ports getCurrentSignaturePolicyId().
func (p *XAdES132Path) CurrentSignaturePolicyId() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, XAdES132ElementSigPolicyID, XAdES132ElementIdentifier)
}

// CurrentSignaturePolicyDigestAlgAndValue implements XAdESPath. Ports getCurrentSignaturePolicyDigestAlgAndValue().
func (p *XAdES132Path) CurrentSignaturePolicyDigestAlgAndValue() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, XAdES132ElementSigPolicyHash)
}

// CurrentSignaturePolicySPURI implements XAdESPath. Ports getCurrentSignaturePolicySPURI().
func (p *XAdES132Path) CurrentSignaturePolicySPURI() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, XAdES132ElementSigPolicyQualifiers, XAdES132ElementSigPolicyQualifier, XAdES132ElementSPURI)
}

// CurrentSignaturePolicySPUserNotice implements XAdESPath. Ports getCurrentSignaturePolicySPUserNotice().
func (p *XAdES132Path) CurrentSignaturePolicySPUserNotice() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, XAdES132ElementSigPolicyQualifiers, XAdES132ElementSigPolicyQualifier, XAdES132ElementSPUserNotice)
}

// CurrentSPUserNoticeNoticeRefOrganization implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefOrganization().
func (p *XAdES132Path) CurrentSPUserNoticeNoticeRefOrganization() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementNoticeRef, XAdES132ElementOrganization)
}

// CurrentSPUserNoticeNoticeRefNoticeNumbers implements XAdESPath. Ports getCurrentSPUserNoticeNoticeRefNoticeNumbers().
func (p *XAdES132Path) CurrentSPUserNoticeNoticeRefNoticeNumbers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementNoticeRef, XAdES132ElementNoticeNumbers)
}

// CurrentSPUserNoticeExplicitText implements XAdESPath. Ports getCurrentSPUserNoticeExplicitText().
func (p *XAdES132Path) CurrentSPUserNoticeExplicitText() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementExplicitText)
}

// CurrentSignaturePolicySPDocSpecification implements XAdESPath. Ports getCurrentSignaturePolicySPDocSpecification().
func (p *XAdES132Path) CurrentSignaturePolicySPDocSpecification() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, XAdES132ElementSigPolicyQualifiers, XAdES132ElementSigPolicyQualifier, XAdES141ElementSPDocSpecification)
}

// CurrentSignaturePolicySPDocSpecificationIdentifier implements XAdESPath. Ports getCurrentSignaturePolicySPDocSpecificationIdentifier().
func (p *XAdES132Path) CurrentSignaturePolicySPDocSpecificationIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, XAdES132ElementSigPolicyQualifiers, XAdES132ElementSigPolicyQualifier, XAdES141ElementSPDocSpecification, XAdES132ElementIdentifier)
}

// CurrentSignaturePolicyDescription implements XAdESPath. Ports getCurrentSignaturePolicyDescription().
func (p *XAdES132Path) CurrentSignaturePolicyDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, XAdES132ElementSigPolicyID, XAdES132ElementDescription)
}

// CurrentSignaturePolicyDocumentationReferences implements XAdESPath. Ports getCurrentSignaturePolicyDocumentationReferences().
func (p *XAdES132Path) CurrentSignaturePolicyDocumentationReferences() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, XAdES132ElementSigPolicyID, XAdES132ElementDocumentationReferences)
}

// CurrentSignaturePolicyImplied implements XAdESPath. Ports getCurrentSignaturePolicyImplied().
func (p *XAdES132Path) CurrentSignaturePolicyImplied() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyImplied)
}

// CurrentSignaturePolicyTransforms implements XAdESPath. Ports getCurrentSignaturePolicyTransforms().
func (p *XAdES132Path) CurrentSignaturePolicyTransforms() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, common.XMLDSigElementTransforms)
}

// CurrentSignaturePolicyQualifiers implements XAdESPath. Ports getCurrentSignaturePolicyQualifiers().
func (p *XAdES132Path) CurrentSignaturePolicyQualifiers() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementSignaturePolicyID, XAdES132ElementSigPolicyQualifiers)
}

// CurrentInclude implements XAdESPath. Ports getCurrentInclude().
func (p *XAdES132Path) CurrentInclude() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementInclude)
}

// CurrentQualifyingPropertiesPath implements XAdESPath. Ports getCurrentQualifyingPropertiesPath().
func (p *XAdES132Path) CurrentQualifyingPropertiesPath() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementQualifyingProperties)
}

// CurrentSPDocSpecification implements XAdESPath. Ports getCurrentSPDocSpecification().
func (p *XAdES132Path) CurrentSPDocSpecification() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141ElementSPDocSpecification)
}

// CurrentIdentifier implements XAdESPath. Ports getCurrentIdentifier().
func (p *XAdES132Path) CurrentIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementIdentifier)
}

// CurrentSPDocSpecificationIdentifier implements XAdESPath. Ports getCurrentSPDocSpecificationIdentifier().
func (p *XAdES132Path) CurrentSPDocSpecificationIdentifier() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141ElementSPDocSpecification, XAdES132ElementIdentifier)
}

// CurrentSPDocSpecificationDescription implements XAdESPath. Ports getCurrentSPDocSpecificationDescription().
func (p *XAdES132Path) CurrentSPDocSpecificationDescription() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141ElementSPDocSpecification, XAdES132ElementDescription)
}

// CurrentDocumentationReferenceElements implements XAdESPath. Ports getCurrentDocumentationReferenceElements().
func (p *XAdES132Path) CurrentDocumentationReferenceElements() common.XPathQuery {
	return common.FromCurrentPosition(XAdES132ElementDocumentationReferences, XAdES132ElementDocumentationReference)
}

// CurrentSPDocSpecificationDocumentationReferenceElements implements XAdESPath. Ports getCurrentSPDocSpecificationDocumentationReferenceElements().
func (p *XAdES132Path) CurrentSPDocSpecificationDocumentationReferenceElements() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141ElementSPDocSpecification, XAdES132ElementDocumentationReferences, XAdES132ElementDocumentationReference)
}

// CurrentSignaturePolicyDocument implements XAdESPath. Ports getCurrentSignaturePolicyDocument().
func (p *XAdES132Path) CurrentSignaturePolicyDocument() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141ElementSignaturePolicyDocument)
}

// CurrentSigPolDocLocalURI implements XAdESPath. Ports getCurrentSigPolDocLocalURI().
func (p *XAdES132Path) CurrentSigPolDocLocalURI() common.XPathQuery {
	return common.FromCurrentPosition(XAdES141ElementSigPolDocLocalURI)
}

var _ XAdESPath = (*XAdES132Path)(nil)
