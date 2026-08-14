// Ported from dss-model/.../model/policy/ValidationPolicy.java (DSS 6.5.RC1).
package policy

import "github.com/utain/esig/dss/enumerations"

// ValidationPolicy encapsulates the constraint file that controls the policy
// used during the validation process.
//
// MultiValuesRule, NumericValueRule, and CertificateApplicabilityRule are
// defined elsewhere in this package (model.policy.crypto flattened in);
// CryptographicSuite is assumed to exist in this package per the porting plan.
type ValidationPolicy interface {

	// PolicyName returns the name of the policy.
	PolicyName() string

	// PolicyDescription returns the policy description.
	PolicyDescription() string

	// SignaturePolicyConstraint indicates if the signature policy should be checked.
	SignaturePolicyConstraint(context enumerations.Context) MultiValuesRule

	// SignaturePolicyIdentifiedConstraint indicates if the signature policy validation should be processed.
	SignaturePolicyIdentifiedConstraint(context enumerations.Context) LevelRule

	// SignaturePolicyStorePresentConstraint indicates if a SignaturePolicyStore unsigned attribute presence shall be checked.
	SignaturePolicyStorePresentConstraint(context enumerations.Context) LevelRule

	// SignaturePolicyPolicyHashValid indicates if the SignaturePolicyIdentifier digest shall match the extracted policy content.
	SignaturePolicyPolicyHashValid(context enumerations.Context) LevelRule

	// StructuralValidationConstraint indicates if the structural validation should be checked.
	StructuralValidationConstraint(context enumerations.Context) LevelRule

	// SigningCertificateRefersCertificateChainConstraint indicates if the Signing Certificate attribute should be checked against the certificate chain.
	SigningCertificateRefersCertificateChainConstraint(context enumerations.Context) LevelRule

	// ReferencesToAllCertificateChainPresentConstraint indicates if the whole certificate chain is covered by the Signing Certificate attribute.
	ReferencesToAllCertificateChainPresentConstraint(context enumerations.Context) LevelRule

	// SigningCertificateDigestAlgorithmConstraint checks the DigestAlgorithm used in signing-certificate-reference creation.
	SigningCertificateDigestAlgorithmConstraint(context enumerations.Context) LevelRule

	// SigningTimeConstraint indicates if the signed property signing-time should be checked.
	SigningTimeConstraint(context enumerations.Context) LevelRule

	// SigningTimeInCertRangeConstraint indicates if signing-time should be checked against the signing-certificate's validity period.
	SigningTimeInCertRangeConstraint(context enumerations.Context) LevelRule

	// SignatureTypeConstraint indicates if the signed property signature type should be checked.
	SignatureTypeConstraint(context enumerations.Context) MultiValuesRule

	// ContentTypeConstraint indicates if the signed property content-type should be checked.
	ContentTypeConstraint(context enumerations.Context) MultiValuesRule

	// ContentHintsConstraint indicates if the signed property content-hints should be checked.
	ContentHintsConstraint(context enumerations.Context) MultiValuesRule

	// ContentIdentifierConstraint indicates if the signed property content-identifier should be checked.
	ContentIdentifierConstraint(context enumerations.Context) MultiValuesRule

	// MessageDigestOrSignedPropertiesConstraint indicates if message-digest (CAdES) or SignedProperties (XAdES) should be checked.
	MessageDigestOrSignedPropertiesConstraint(context enumerations.Context) LevelRule

	// EllipticCurveKeySizeConstraint checks whether a JWA signature has a valid elliptic curve key size.
	EllipticCurveKeySizeConstraint(context enumerations.Context) LevelRule

	// CommitmentTypeIndicationConstraint indicates if the signed property commitment-type-indication should be checked.
	CommitmentTypeIndicationConstraint(context enumerations.Context) MultiValuesRule

	// SignerLocationConstraint indicates if the signed property signer-location should be checked.
	SignerLocationConstraint(context enumerations.Context) LevelRule

	// ContentTimeStampConstraint indicates if the signed property content-time-stamp should be checked.
	ContentTimeStampConstraint(context enumerations.Context) LevelRule

	// ContentTimeStampMessageImprintConstraint indicates if the content-time-stamp message-imprint should be checked.
	ContentTimeStampMessageImprintConstraint(context enumerations.Context) LevelRule

	// ClaimedRoleConstraint indicates if the unsigned property claimed-role should be checked.
	ClaimedRoleConstraint(context enumerations.Context) MultiValuesRule

	// CertifiedRolesConstraint returns the mandated signer role.
	CertifiedRolesConstraint(context enumerations.Context) MultiValuesRule

	// SignatureCryptographicConstraint creates the CryptographicSuite corresponding to the context parameter.
	SignatureCryptographicConstraint(context enumerations.Context) CryptographicSuite

	// CertificateCryptographicConstraint creates the CryptographicSuite corresponding to the context/subContext parameters.
	CertificateCryptographicConstraint(context enumerations.Context, subContext enumerations.SubContext) CryptographicSuite

	// EvidenceRecordCryptographicConstraint returns cryptographic constraints for validation of Evidence Record.
	EvidenceRecordCryptographicConstraint() CryptographicSuite

	// EAACryptographicConstraint returns cryptographic constraints for validation of EAA.
	EAACryptographicConstraint() CryptographicSuite

	// CertificateCAConstraint returns certificate CA constraint.
	CertificateCAConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateMaxPathLengthConstraint returns certificate MaxPathLength constraint.
	CertificateMaxPathLengthConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateKeyUsageConstraint returns certificate key usage constraint.
	CertificateKeyUsageConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateExtendedKeyUsageConstraint returns certificate extended key usage constraint.
	CertificateExtendedKeyUsageConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificatePolicyTreeConstraint returns certificate PolicyTree constraint.
	CertificatePolicyTreeConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateNameConstraintsConstraint returns certificate NameConstraints constraint.
	CertificateNameConstraintsConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateAuthorityKeyIdentifierPresentConstraint returns certificate AuthorityKeyIdentifierPresent constraint.
	CertificateAuthorityKeyIdentifierPresentConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateSubjectKeyIdentifierPresentConstraint returns certificate SubjectKeyIdentifierPresent constraint.
	CertificateSubjectKeyIdentifierPresentConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateNoRevAvailConstraint returns certificate NoRevAvail constraint.
	CertificateNoRevAvailConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateSupportedCriticalExtensionsConstraint returns certificate supported critical extensions constraint.
	CertificateSupportedCriticalExtensionsConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateForbiddenExtensionsConstraint returns certificate forbidden extensions constraint.
	CertificateForbiddenExtensionsConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateNotExpiredConstraint returns certificate's validity range constraint.
	CertificateNotExpiredConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateSunsetDateConstraint returns certificate's sunset date constraint.
	CertificateSunsetDateConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// ProspectiveCertificateChainConstraint requests the presence of the trust anchor in the certificate chain.
	ProspectiveCertificateChainConstraint(context enumerations.Context) LevelRule

	// CertificateSignatureConstraint returns certificate's signature constraint.
	CertificateSignatureConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// UnknownStatusConstraint returns the UnknownStatus constraint.
	UnknownStatusConstraint() LevelRule

	// ThisUpdatePresentConstraint returns the ThisUpdatePresent constraint.
	ThisUpdatePresentConstraint() LevelRule

	// RevocationIssuerKnownConstraint returns the RevocationIssuerKnown constraint.
	RevocationIssuerKnownConstraint() LevelRule

	// RevocationIssuerValidAtProductionTimeConstraint returns the RevocationIssuerValidAtProductionTime constraint.
	RevocationIssuerValidAtProductionTimeConstraint() LevelRule

	// RevocationAfterCertificateIssuanceConstraint returns the RevocationIssuerKnowsCertificate constraint.
	RevocationAfterCertificateIssuanceConstraint() LevelRule

	// RevocationHasInformationAboutCertificateConstraint returns the RevocationIssuerHasInformationAboutCertificate constraint.
	RevocationHasInformationAboutCertificateConstraint() LevelRule

	// OCSPResponseResponderIdMatchConstraint returns the OCSPResponderIdMatch constraint.
	OCSPResponseResponderIdMatchConstraint() LevelRule

	// OCSPResponseCertHashPresentConstraint returns the OCSPCertHashPresent constraint.
	OCSPResponseCertHashPresentConstraint() LevelRule

	// OCSPResponseCertHashMatchConstraint returns the OCSPCertHashMatch constraint.
	OCSPResponseCertHashMatchConstraint() LevelRule

	// SelfIssuedOCSPConstraint returns the SelfIssuedOCSP constraint.
	SelfIssuedOCSPConstraint() LevelRule

	// RevocationDataAvailableConstraint returns revocation data available constraint.
	RevocationDataAvailableConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// AcceptableRevocationDataFoundConstraint returns acceptable revocation data available constraint.
	AcceptableRevocationDataFoundConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CRLNextUpdatePresentConstraint returns CRL's nextUpdate present constraint.
	CRLNextUpdatePresentConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// OCSPNextUpdatePresentConstraint returns OCSP's nextUpdate present constraint.
	OCSPNextUpdatePresentConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// RevocationFreshnessConstraint returns revocation data's freshness constraint.
	RevocationFreshnessConstraint(context enumerations.Context, subContext enumerations.SubContext) DurationRule

	// RevocationFreshnessNextUpdateConstraint returns revocation data's freshness for nextUpdate check constraint.
	RevocationFreshnessNextUpdateConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateNotRevokedConstraint returns certificate's not revoked constraint.
	CertificateNotRevokedConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateNotOnHoldConstraint returns certificate's not onHold constraint.
	CertificateNotOnHoldConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// RevocationIssuerNotExpiredConstraint returns revocation issuer's validity range constraint.
	RevocationIssuerNotExpiredConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateNotSelfSignedConstraint returns certificate's not self-signed constraint.
	CertificateNotSelfSignedConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateSelfSignedConstraint returns certificate's self-signed constraint.
	CertificateSelfSignedConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// TrustServiceTypeIdentifierConstraint returns trusted service type identifier constraint.
	TrustServiceTypeIdentifierConstraint(context enumerations.Context) MultiValuesRule

	// TrustServiceStatusConstraint returns trusted service status constraint.
	TrustServiceStatusConstraint(context enumerations.Context) MultiValuesRule

	// CertificatePolicyIdsConstraint returns the CertificatePolicyIds constraint, if present.
	CertificatePolicyIdsConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificatePolicyQualificationIdsConstraint indicates if the CertificatePolicyIds declare the certificate as qualified.
	CertificatePolicyQualificationIdsConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificatePolicySupportedByQSCDIdsConstraint indicates if the CertificatePolicyIds mandate QSCD support.
	CertificatePolicySupportedByQSCDIdsConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateQCComplianceConstraint indicates if the end user certificate is QC Compliant.
	CertificateQCComplianceConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateQcEuLimitValueCurrencyConstraint indicates the allowed currency for the QCLimitValue statement.
	CertificateQcEuLimitValueCurrencyConstraint(context enumerations.Context, subContext enumerations.SubContext) ValueRule

	// CertificateMinQcEuLimitValueConstraint indicates the minimal allowed QcEuLimitValue transaction limit.
	CertificateMinQcEuLimitValueConstraint(context enumerations.Context, subContext enumerations.SubContext) NumericValueRule

	// CertificateMinQcEuRetentionPeriodConstraint indicates the minimal allowed QC retention period.
	CertificateMinQcEuRetentionPeriodConstraint(context enumerations.Context, subContext enumerations.SubContext) NumericValueRule

	// CertificateQcSSCDConstraint indicates if the end user certificate is mandated to be QSCD-supported.
	CertificateQcSSCDConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateQcEuPDSLocationConstraint indicates the location(s) of PKI Disclosure Statements.
	CertificateQcEuPDSLocationConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateQcTypeConstraint indicates the claimed certificate type(s).
	CertificateQcTypeConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateQcCCLegislationConstraint indicates the country/countries under whose legislation the certificate is issued as qualified.
	CertificateQcCCLegislationConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateIssuedToNaturalPersonConstraint indicates if the end user certificate is issued to a natural person.
	CertificateIssuedToNaturalPersonConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateIssuedToLegalPersonConstraint indicates if the end user certificate is issued to a legal person.
	CertificateIssuedToLegalPersonConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateSemanticsIdentifierConstraint indicates the acceptable QCStatement semantics identifier.
	CertificateSemanticsIdentifierConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificatePS2DQcTypeRolesOfPSPConstraint indicates the acceptable QC PS2D roles.
	CertificatePS2DQcTypeRolesOfPSPConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificatePS2DQcCompetentAuthorityNameConstraint indicates the acceptable QC PS2D names.
	CertificatePS2DQcCompetentAuthorityNameConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificatePS2DQcCompetentAuthorityIdConstraint indicates the acceptable QC PS2D ids.
	CertificatePS2DQcCompetentAuthorityIdConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateQcQSCDLegislationConstraint indicates the country/countries under whose legislation the signature creation device has a qualified status.
	CertificateQcQSCDLegislationConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateQcIdentificationMethodConstraint indicates the verification method used for issuance per eIDAS Article 24.
	CertificateQcIdentificationMethodConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateQcPSBCountryOfLegislationConstraint indicates the verification method for country of legislation of the QcPSB QcStatement.
	CertificateQcPSBCountryOfLegislationConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateQcPSBAuthSourceIdentificationConstraint indicates the verification method for authentic source identification of the QcPSB QcStatement.
	CertificateQcPSBAuthSourceIdentificationConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateQcPSBLegislationIdentificationConstraint indicates the verification method for legislation identification of the QcPSB QcStatement.
	CertificateQcPSBLegislationIdentificationConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// SigningCertificateRecognitionConstraint indicates if signing-certificate has been identified.
	SigningCertificateRecognitionConstraint(context enumerations.Context) LevelRule

	// SigningCertificateAttributePresentConstraint indicates if the signing certificate attribute is present.
	SigningCertificateAttributePresentConstraint(context enumerations.Context) LevelRule

	// UnicitySigningCertificateAttributeConstraint indicates if the signing certificate is not ambiguously determined.
	UnicitySigningCertificateAttributeConstraint(context enumerations.Context) LevelRule

	// SigningCertificateDigestValuePresentConstraint indicates if the signing certificate reference's digest value is present.
	SigningCertificateDigestValuePresentConstraint(context enumerations.Context) LevelRule

	// SigningCertificateDigestValueMatchConstraint indicates if the signing certificate reference's digest value matches.
	SigningCertificateDigestValueMatchConstraint(context enumerations.Context) LevelRule

	// SigningCertificateIssuerSerialMatchConstraint indicates if the signing certificate reference's issuer serial matches.
	SigningCertificateIssuerSerialMatchConstraint(context enumerations.Context) LevelRule

	// KeyIdentifierPresent indicates if the 'kid' header parameter is present.
	KeyIdentifierPresent(context enumerations.Context) LevelRule

	// KeyIdentifierMatch indicates if the value of 'kid' matches the signing-certificate.
	KeyIdentifierMatch(context enumerations.Context) LevelRule

	// X509UrlPresent indicates if the value of 'x5u' header parameter is present.
	X509UrlPresent(context enumerations.Context) LevelRule

	// X509UrlMatch indicates if the signing-certificate can be derived from the 'x5u' header parameter.
	X509UrlMatch(context enumerations.Context) LevelRule

	// ReferenceDataExistenceConstraint indicates if the referenced data is found.
	ReferenceDataExistenceConstraint(context enumerations.Context) LevelRule

	// ReferenceDataIntactConstraint indicates if the referenced data is intact.
	ReferenceDataIntactConstraint(context enumerations.Context) LevelRule

	// ReferenceDataNameMatchConstraint indicates if the referenced document names match the manifest entry references.
	ReferenceDataNameMatchConstraint(context enumerations.Context) LevelRule

	// ManifestEntryObjectExistenceConstraint indicates if the manifested document is found.
	ManifestEntryObjectExistenceConstraint(context enumerations.Context) LevelRule

	// ManifestEntryObjectIntactConstraint indicates if the manifested document is intact.
	ManifestEntryObjectIntactConstraint(context enumerations.Context) LevelRule

	// ManifestEntryObjectGroupConstraint indicates if all manifest entries have been found.
	ManifestEntryObjectGroupConstraint(context enumerations.Context) LevelRule

	// ManifestEntryNameMatchConstraint indicates if names of all matching documents match the manifest entry names.
	ManifestEntryNameMatchConstraint(context enumerations.Context) LevelRule

	// SignatureIntactConstraint indicates if the signature is intact.
	SignatureIntactConstraint(context enumerations.Context) LevelRule

	// SignatureDuplicatedConstraint indicates if the signature is not ambiguous.
	SignatureDuplicatedConstraint(context enumerations.Context) LevelRule

	// SignerInformationStoreConstraint checks if only one SignerInfo is present in a SignerInformationStore (PAdES only).
	SignerInformationStoreConstraint(context enumerations.Context) LevelRule

	// ByteRangeConstraint checks if the ByteRange dictionary is valid (PAdES only).
	ByteRangeConstraint(context enumerations.Context) LevelRule

	// ByteRangeCollisionConstraint checks if ByteRange does not collide with other signature byte ranges (PAdES only).
	ByteRangeCollisionConstraint(context enumerations.Context) LevelRule

	// ByteRangeAllDocumentConstraint checks if ByteRange is valid for all signatures and document timestamps (PAdES only).
	ByteRangeAllDocumentConstraint(context enumerations.Context) LevelRule

	// PdfSignatureDictionaryConstraint checks if signature dictionary is consistent across PDF revisions (PAdES only).
	PdfSignatureDictionaryConstraint(context enumerations.Context) LevelRule

	// PdfPageDifferenceConstraint indicates if a PDF page difference check should be proceeded.
	PdfPageDifferenceConstraint(context enumerations.Context) LevelRule

	// PdfAnnotationOverlapConstraint indicates if a PDF annotation overlapping check should be proceeded.
	PdfAnnotationOverlapConstraint(context enumerations.Context) LevelRule

	// PdfVisualDifferenceConstraint indicates if a PDF visual difference check should be proceeded.
	PdfVisualDifferenceConstraint(context enumerations.Context) LevelRule

	// DocMDPConstraint checks for changes against permission rules identified within a /DocMDP dictionary.
	DocMDPConstraint(context enumerations.Context) LevelRule

	// FieldMDPConstraint checks for changes against permission rules identified within a /FieldMDP dictionary.
	FieldMDPConstraint(context enumerations.Context) LevelRule

	// SigFieldLockConstraint checks for changes against permission rules identified within a /SigFieldLock dictionary.
	SigFieldLockConstraint(context enumerations.Context) LevelRule

	// FormFillChangesConstraint checks whether a PDF document contains form fill or signing modifications after the current signature's revisions.
	FormFillChangesConstraint(context enumerations.Context) LevelRule

	// AnnotationChangesConstraint checks whether a PDF document contains annotation modifications after the current signature's revisions.
	AnnotationChangesConstraint(context enumerations.Context) LevelRule

	// UndefinedChangesConstraint checks whether a PDF document contains undefined object modifications after the current signature's revisions.
	UndefinedChangesConstraint(context enumerations.Context) LevelRule

	// BestSignatureTimeBeforeExpirationDateOfSigningCertificateConstraint checks if the certificate is not expired on best-signature-time.
	BestSignatureTimeBeforeExpirationDateOfSigningCertificateConstraint() LevelRule

	// TimestampCoherenceConstraint checks if the timestamp order is coherent.
	TimestampCoherenceConstraint() LevelRule

	// TimestampDelayConstraint returns the TimestampDelay constraint, if present.
	TimestampDelayConstraint() DurationRule

	// TimestampValidConstraint returns whether the time-stamp is valid.
	TimestampValidConstraint() LevelRule

	// TimestampTSAGeneralNamePresent indicates if the timestamp's TSTInfo.tsa field is present.
	TimestampTSAGeneralNamePresent() LevelRule

	// TimestampTSAGeneralNameContentMatch indicates if TSTInfo.tsa matches the timestamp's issuer distinguishing name.
	TimestampTSAGeneralNameContentMatch() LevelRule

	// TimestampTSAGeneralNameOrderMatch indicates if TSTInfo.tsa value and order match the timestamp's issuer distinguishing name.
	TimestampTSAGeneralNameOrderMatch() LevelRule

	// AtsHashIndexConstraint returns the timestamp AtsHashIndex constraint, if present.
	AtsHashIndexConstraint() LevelRule

	// TimestampContainerSignedAndTimestampedFilesCoveredConstraint returns the timestamp ContainerSignedAndTimestampedFilesCovered constraint, if present.
	TimestampContainerSignedAndTimestampedFilesCoveredConstraint() LevelRule

	// RevocationTimeAgainstBestSignatureTimeConstraint returns the RevocationTimeAgainstBestSignatureTime constraint, if present.
	RevocationTimeAgainstBestSignatureTimeConstraint() LevelRule

	// EvidenceRecordValidConstraint returns whether the evidence record is valid.
	EvidenceRecordValidConstraint() LevelRule

	// EvidenceRecordDataObjectExistenceConstraint returns the DataObjectExistence constraint, if present.
	EvidenceRecordDataObjectExistenceConstraint() LevelRule

	// EvidenceRecordDataObjectIntactConstraint returns the DataObjectIntact constraint, if present.
	EvidenceRecordDataObjectIntactConstraint() LevelRule

	// EvidenceRecordDataObjectFoundConstraint returns the DataObjectFound constraint, if present.
	EvidenceRecordDataObjectFoundConstraint() LevelRule

	// EvidenceRecordDataObjectGroupConstraint returns the DataObjectGroup constraint, if present.
	EvidenceRecordDataObjectGroupConstraint() LevelRule

	// EvidenceRecordSignedFilesCoveredConstraint returns the SignedFilesCovered constraint, if present.
	EvidenceRecordSignedFilesCoveredConstraint() LevelRule

	// EvidenceRecordContainerSignedAndTimestampedFilesCoveredConstraint returns the evidence record ContainerSignedAndTimestampedFilesCovered constraint, if present.
	EvidenceRecordContainerSignedAndTimestampedFilesCoveredConstraint() LevelRule

	// EvidenceRecordHashTreeRenewalConstraint returns the HashTreeRenewal constraint, if present.
	EvidenceRecordHashTreeRenewalConstraint() LevelRule

	// CounterSignatureConstraint returns the CounterSignature constraint, if present.
	CounterSignatureConstraint(context enumerations.Context) LevelRule

	// SignatureTimeStampConstraint indicates if the presence of the signature-time-stamp unsigned property should be checked.
	SignatureTimeStampConstraint(context enumerations.Context) LevelRule

	// ValidationDataTimeStampConstraint indicates if the presence of the validation data timestamp unsigned property should be checked.
	ValidationDataTimeStampConstraint(context enumerations.Context) LevelRule

	// ValidationDataRefsOnlyTimeStampConstraint indicates if the presence of the validation data references only timestamp unsigned property should be checked.
	ValidationDataRefsOnlyTimeStampConstraint(context enumerations.Context) LevelRule

	// ArchiveTimeStampConstraint indicates if the presence of the archive-time-stamp unsigned property should be checked.
	ArchiveTimeStampConstraint(context enumerations.Context) LevelRule

	// DocumentTimeStampConstraint indicates if the presence of the document timestamp unsigned property should be checked.
	DocumentTimeStampConstraint(context enumerations.Context) LevelRule

	// TLevelTimeStampConstraint indicates if the presence of the signature-time-stamp or document timestamp should be checked.
	TLevelTimeStampConstraint(context enumerations.Context) LevelRule

	// LTALevelTimeStampConstraint indicates if the presence of the archive-time-stamp or document timestamp covering the validation data should be checked.
	LTALevelTimeStampConstraint(context enumerations.Context) LevelRule

	// SignatureFormatConstraint returns the SignatureFormat constraint, if present.
	SignatureFormatConstraint(context enumerations.Context) MultiValuesRule

	// CertificateCountryConstraint returns the CertificateCountry constraint, if present.
	CertificateCountryConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateLocalityConstraint returns the CertificateLocality constraint, if present.
	CertificateLocalityConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateStateConstraint returns the CertificateState constraint, if present.
	CertificateStateConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateOrganizationIdentifierConstraint returns the CertificateOrganizationIdentifier constraint, if present.
	CertificateOrganizationIdentifierConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateOrganizationNameConstraint returns the CertificateOrganizationName constraint, if present.
	CertificateOrganizationNameConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateOrganizationUnitConstraint returns the CertificateOrganizationUnit constraint, if present.
	CertificateOrganizationUnitConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateIssuerNameConstraint returns certificate IssuerName constraint.
	CertificateIssuerNameConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateSurnameConstraint returns the CertificateSurname constraint, if present.
	CertificateSurnameConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateGivenNameConstraint returns the CertificateGivenName constraint, if present.
	CertificateGivenNameConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateCommonNameConstraint returns the CertificateCommonName constraint, if present.
	CertificateCommonNameConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificatePseudonymConstraint returns the CertificatePseudonym constraint, if present.
	CertificatePseudonymConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificatePseudoUsageConstraint returns the CertificatePseudoUsage constraint, if present.
	CertificatePseudoUsageConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateTitleConstraint returns the CertificateTitle constraint, if present.
	CertificateTitleConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateEmailConstraint returns the CertificateEmail constraint, if present.
	CertificateEmailConstraint(context enumerations.Context, subContext enumerations.SubContext) MultiValuesRule

	// CertificateSerialNumberConstraint returns the CertificateSerialNumber constraint, if present.
	CertificateSerialNumberConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// CertificateAuthorityInfoAccessPresentConstraint returns the CertificateAuthorityInfoAccessPresent constraint, if present.
	CertificateAuthorityInfoAccessPresentConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// RevocationDataSkipConstraint returns the RevocationDataSkip constraint, if present.
	RevocationDataSkipConstraint(context enumerations.Context, subContext enumerations.SubContext) CertificateApplicabilityRule

	// CertificateRevocationInfoAccessPresentConstraint returns the CertificateRevocationInfoAccessPresent constraint, if present.
	CertificateRevocationInfoAccessPresentConstraint(context enumerations.Context, subContext enumerations.SubContext) LevelRule

	// AcceptedContainerTypesConstraint returns the AcceptedContainerTypes constraint, if present.
	AcceptedContainerTypesConstraint() MultiValuesRule

	// ZipCommentPresentConstraint returns the ZipCommentPresent constraint, if present.
	ZipCommentPresentConstraint() LevelRule

	// AcceptedZipCommentsConstraint returns the AcceptedZipComments constraint, if present.
	AcceptedZipCommentsConstraint() MultiValuesRule

	// MimeTypeFilePresentConstraint returns the MimeTypeFilePresent constraint, if present.
	MimeTypeFilePresentConstraint() LevelRule

	// AcceptedMimeTypeContentsConstraint returns the AcceptedMimeTypeContents constraint, if present.
	AcceptedMimeTypeContentsConstraint() MultiValuesRule

	// ManifestFilePresentConstraint returns the ManifestFilePresent constraint, if present.
	ManifestFilePresentConstraint() LevelRule

	// SignedFilesPresentConstraint returns the SignedFilesPresent constraint, if present.
	SignedFilesPresentConstraint() LevelRule

	// FilenameAdherenceConstraint returns the FilenameAdherence constraint, if present.
	FilenameAdherenceConstraint() LevelRule

	// AllFilesSignedConstraint returns the AllFilesSigned constraint, if present.
	AllFilesSignedConstraint() LevelRule

	// FullScopeConstraint returns the FullScope constraint, if present.
	FullScopeConstraint() LevelRule

	// AcceptablePDFAProfilesConstraint returns the AcceptablePDFAProfiles constraint, if present.
	AcceptablePDFAProfilesConstraint() MultiValuesRule

	// PDFACompliantConstraint returns the PDFACompliant constraint, if present.
	PDFACompliantConstraint() LevelRule

	// EAASignatureUnicityConstraint returns the EAASignatureUnicity constraint, if present.
	EAASignatureUnicityConstraint() LevelRule

	// EAASignatureValidConstraint returns the EAASignatureValid constraint, if present.
	EAASignatureValidConstraint() LevelRule

	// EAADisclosurePresentConstraint returns the DisclosurePresent constraint, if present.
	EAADisclosurePresentConstraint() LevelRule

	// EAADisclosureFoundConstraint returns the DisclosureFound constraint, if present.
	EAADisclosureFoundConstraint() LevelRule

	// EAADisclosureIntactConstraint returns the DisclosureIntact constraint, if present.
	EAADisclosureIntactConstraint() LevelRule

	// EAADisclosureListExhaustiveConstraint returns the DisclosureListExhaustive constraint, if present.
	EAADisclosureListExhaustiveConstraint() LevelRule

	// EAAKeyBindingSignaturePresentConstraint returns the KeyBindingSignaturePresent constraint, if present.
	EAAKeyBindingSignaturePresentConstraint() LevelRule

	// EAAKeyBindingSignatureValidConstraint returns the KeyBindingSignatureValid constraint, if present.
	EAAKeyBindingSignatureValidConstraint() LevelRule

	// EAATypeIntegrityPresentConstraint returns the EAATypeIntegrityPresent constraint, if present.
	EAATypeIntegrityPresentConstraint() LevelRule

	// EAAIdentifierPresentConstraint returns the EAAIdentifierPresent constraint, if present.
	EAAIdentifierPresentConstraint() LevelRule

	// EAAIssuanceDatePresentConstraint returns the EAAIssuanceDatePresent constraint, if present.
	EAAIssuanceDatePresentConstraint() LevelRule

	// EAACategoryConstraint returns the EAACategory constraint, if present.
	EAACategoryConstraint() MultiValuesRule

	// EAASubjectConstraint returns the EAASubject constraint, if present.
	EAASubjectConstraint() MultiValuesRule

	// EAASubjectPseudonymConstraint returns the EAASubjectPseudonym constraint, if present.
	EAASubjectPseudonymConstraint() MultiValuesRule

	// EAAIssuingCountryConstraint returns the EAAIssuingCountry constraint, if present.
	EAAIssuingCountryConstraint() MultiValuesRule

	// EAAIssuingAuthorityConstraint returns the EAAIssuingAuthority constraint, if present.
	EAAIssuingAuthorityConstraint() MultiValuesRule

	// EAAIssuingAuthorityRegistrationIdentifierConstraint returns the EAAIssuingAuthorityRegistrationIdentifier constraint, if present.
	EAAIssuingAuthorityRegistrationIdentifierConstraint() MultiValuesRule

	// EAARevocationPresentConstraint returns the EAARevocationPresent constraint, if present.
	EAARevocationPresentConstraint() LevelRule

	// EAAShortLivedConstraint returns the EAAShortLived constraint, if present.
	EAAShortLivedConstraint() LevelRule

	// EAAOneTimeUseConstraint returns the EAAOneTimeUse constraint, if present.
	EAAOneTimeUseConstraint() LevelRule

	// EAAUsePseudonymConstraint returns the EAAUsePseudonym constraint, if present.
	EAAUsePseudonymConstraint() LevelRule

	// EAAClaimsConstraint returns the EAAClaims constraint, if present.
	EAAClaimsConstraint() MultiValuesRule

	// EAASupportedClaimsConstraint returns the EAASupportedClaims constraint, if present.
	EAASupportedClaimsConstraint() MultiValuesRule

	// EAARevocationAvailableConstraint returns the EAARevocationAvailable constraint, if present.
	EAARevocationAvailableConstraint() LevelRule

	// AcceptableEAARevocationFoundConstraint returns the AcceptableEAARevocationFound constraint, if present.
	AcceptableEAARevocationFoundConstraint() LevelRule

	// EAARevocationNotRevokedConstraint returns the NotRevoked constraint, if present.
	EAARevocationNotRevokedConstraint() LevelRule

	// EAARevocationNotOnHoldConstraint returns the NotOnHold constraint, if present.
	EAARevocationNotOnHoldConstraint() LevelRule

	// EAATypeConstraint returns the EAAType constraint, if present.
	EAATypeConstraint() MultiValuesRule

	// EAANotBeforePresentConstraint returns the EAANotBeforePresent constraint, if present.
	EAANotBeforePresentConstraint() LevelRule

	// EAAExpirationPresentConstraint returns the EAAExpirationPresent constraint, if present.
	EAAExpirationPresentConstraint() LevelRule

	// EAANotExpiredConstraint returns the EAANotExpired constraint, if present.
	EAANotExpiredConstraint() LevelRule

	// EAAAdministrativeIssuanceDatePresentConstraint returns the EAAAdministrativeIssuanceDatePresent constraint, if present.
	EAAAdministrativeIssuanceDatePresentConstraint() LevelRule

	// EAAAdministrativeExpirationDatePresentConstraint returns the EAAAdministrativeExpirationDatePresent constraint, if present.
	EAAAdministrativeExpirationDatePresentConstraint() LevelRule

	// EAAAdministrativePeriodNotExpiredConstraint returns the EAAAdministrativePeriodNotExpired constraint, if present.
	EAAAdministrativePeriodNotExpiredConstraint() LevelRule

	// EAAETSI194721ConformanceConstraint returns the ETSI194721Conformance constraint, if present.
	EAAETSI194721ConformanceConstraint() LevelRule

	// EAARevocationTokenTypeConstraint returns the StatusTokenType constraint, if present.
	EAARevocationTokenTypeConstraint() MultiValuesRule

	// EAARevocationUnknownStatusConstraint returns the UnknownStatus constraint, if present.
	EAARevocationUnknownStatusConstraint() LevelRule

	// EAARevocationIssuanceTimeConstraint returns the IssuanceTime constraint, if present.
	EAARevocationIssuanceTimeConstraint() LevelRule

	// EAARevocationExpirationTimeConstraint returns the ExpirationTime constraint, if present.
	EAARevocationExpirationTimeConstraint() LevelRule

	// EAARevocationNotExpiredConstraint returns the NotExpired constraint, if present.
	EAARevocationNotExpiredConstraint() LevelRule

	// EAARevocationSubjectConstraint returns the EAARevocationSubject constraint, if present.
	EAARevocationSubjectConstraint() MultiValuesRule

	// EAARevocationSubjectMatchConstraint returns the EAARevocationSubjectMatch constraint, if present.
	EAARevocationSubjectMatchConstraint() LevelRule

	// EAARevocationIssuerValidAtIssuanceTimeConstraint returns the EAARevocationIssuerValidAtIssuanceTime constraint, if present.
	EAARevocationIssuerValidAtIssuanceTimeConstraint() LevelRule

	// Article 32.

	// EIDASConstraintPresent returns true if EIDAS constraints are present (qualification check shall be performed).
	EIDASConstraintPresent() bool

	// TLFreshnessConstraint returns the TLFreshness constraint, if present.
	TLFreshnessConstraint() DurationRule

	// TLWellSignedConstraint returns the TLWellSigned constraint, if present.
	TLWellSignedConstraint() LevelRule

	// TLNotExpiredConstraint returns the TLNotExpired constraint, if present.
	TLNotExpiredConstraint() LevelRule

	// TLVersionConstraint returns the TLVersion constraint, if present.
	TLVersionConstraint() MultiValuesRule

	// TLStructureConstraint returns the TLStructure constraint, if present.
	TLStructureConstraint() LevelRule

	// LoTEFreshnessConstraint returns the LoTEFreshness constraint, if present.
	LoTEFreshnessConstraint() DurationRule

	// LoTEWellSignedConstraint returns the LoTEWellSigned constraint, if present.
	LoTEWellSignedConstraint() LevelRule

	// LoTENotExpiredConstraint returns the LoTENotExpired constraint, if present.
	LoTENotExpiredConstraint() LevelRule

	// LoTEVersionConstraint returns the LoTEVersion constraint, if present.
	LoTEVersionConstraint() MultiValuesRule

	// LoTEStructureConstraint returns the LoTEStructure constraint, if present.
	LoTEStructureConstraint() LevelRule

	// ValidationModel returns the used validation model (default SHELL; alternatives CHAIN and HYBRID).
	ValidationModel() enumerations.ValidationModel
}
