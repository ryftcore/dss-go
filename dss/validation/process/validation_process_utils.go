// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/ValidationProcessUtils.java (DSS 6.5.RC1).
//
// The three methods that take a vpfswatsp POEExtraction
// (getLatestAcceptableRevocationData, getAcceptableRevocationDataForPSVIfExistOrReturnAll
// and its private helper filterRevocationDataForPastSignatureValidation) live in
// validation_process_utils_vpfswatsp.go behind the "phase8e" build tag, since
// eu.europa.esig.dss.validation.process.vpfswatsp is ported in phase 8e (same
// tag-split precedent as the ASiC phase8 files). Everything else is here.
//
// Java's IllegalArgumentException becomes a returned error (PORTING.md); the
// static utility class becomes package-level functions.
package process

import (
	"fmt"
	"regexp"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
)

// dateFormat is the Validation policy date format.
const dateFormat = "2006-01-02 15:04" // Java: "yyyy-MM-dd HH:mm"

// urnOidPrefix is the prefix used for "urn:oid:" definition as per RFC 3061.
const urnOidPrefix = "urn:oid:"

// allValue is the value used to accept all values.
const allValue = "*"

// domainNamePattern is the regular expression getDomainName strips from a URI.
var domainNamePattern = regexp.MustCompile(`(^.*://)|(www\.)|([?=:#/].*)`)

// IsAllowedBasicSignatureValidation checks if the given conclusion is allowed as
// a basic signature validation in order to continue the validation process with
// Long-Term Validation Data. Port of isAllowedBasicSignatureValidation(XmlConclusion).
func IsAllowedBasicSignatureValidation(conclusion *jaxb.XmlConclusion) bool {
	return conclusion != nil && (enumerations.IndicationPassed == conclusion.Indication.Indication() ||
		(enumerations.IndicationIndeterminate == conclusion.Indication.Indication() &&
			(enumerations.SubIndicationCryptoConstraintsFailureNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationRevokedNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationRevokedCANoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationTryLater == subIndicationOf(conclusion) ||
				enumerations.SubIndicationOutOfBoundsNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationOutOfBoundsNotRevoked == subIndicationOf(conclusion))))
}

// IsAllowedBasicRevocationDataValidation checks if the given conclusion is
// allowed as a basic revocation validation in order to continue the validation
// process with Long-Term Validation Data. Port of
// isAllowedBasicRevocationDataValidation(XmlConclusion).
func IsAllowedBasicRevocationDataValidation(conclusion *jaxb.XmlConclusion) bool {
	return conclusion != nil && (enumerations.IndicationPassed == conclusion.Indication.Indication() ||
		(enumerations.IndicationIndeterminate == conclusion.Indication.Indication() &&
			(enumerations.SubIndicationRevokedNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationRevokedCANoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationOutOfBoundsNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationOutOfBoundsNotRevoked == subIndicationOf(conclusion) ||
				enumerations.SubIndicationCryptoConstraintsFailureNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationRevocationOutOfBoundsNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationNoCertificateChainFoundNoPOE == subIndicationOf(conclusion))))
}

// IsAllowedBasicTimestampValidation checks if the given conclusion is allowed as
// a basic timestamp validation in order to continue the validation process with
// Archival Data. Port of isAllowedBasicTimestampValidation(XmlConclusion).
func IsAllowedBasicTimestampValidation(conclusion *jaxb.XmlConclusion) bool {
	return conclusion != nil && (enumerations.IndicationPassed == conclusion.Indication.Indication() ||
		(enumerations.IndicationIndeterminate == conclusion.Indication.Indication() &&
			(enumerations.SubIndicationRevokedNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationRevokedCANoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationOutOfBoundsNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationOutOfBoundsNotRevoked == subIndicationOf(conclusion) ||
				enumerations.SubIndicationCryptoConstraintsFailureNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationRevocationOutOfBoundsNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationNoCertificateChainFoundNoPOE == subIndicationOf(conclusion))))
}

// IsAllowedValidationWithLongTermData checks if the given conclusion is allowed
// as a validation process with a long-term validation data in order to continue
// the validation process with Archival Data. Port of
// isAllowedValidationWithLongTermData(XmlConclusion).
func IsAllowedValidationWithLongTermData(conclusion *jaxb.XmlConclusion) bool {
	return conclusion != nil && (enumerations.IndicationPassed == conclusion.Indication.Indication() ||
		(enumerations.IndicationIndeterminate == conclusion.Indication.Indication() &&
			(enumerations.SubIndicationRevokedNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationRevokedCANoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationOutOfBoundsNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationOutOfBoundsNotRevoked == subIndicationOf(conclusion) ||
				enumerations.SubIndicationCryptoConstraintsFailureNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationRevocationOutOfBoundsNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationNoCertificateChainFoundNoPOE == subIndicationOf(conclusion) ||
				enumerations.SubIndicationSigConstraintsFailure == subIndicationOf(conclusion) ||
				enumerations.SubIndicationTryLater == subIndicationOf(conclusion))))
}

// subIndicationOf reads XmlConclusion#getSubIndication(): the generated member
// is a pointer, whose nil is Java's null.
func subIndicationOf(conclusion *jaxb.XmlConclusion) enumerations.SubIndication {
	if conclusion.SubIndication == nil {
		return ""
	}
	return conclusion.SubIndication.SubIndication()
}

// IsTrustAnchor verifies whether the given certificateWrapper can be considered
// as a trust anchor at the currentTime. Port of
// isTrustAnchor(CertificateWrapper, Date, LevelRule).
func IsTrustAnchor(certificateWrapper *diagnostic.CertificateWrapper, currentTime time.Time,
	certificateSunsetDateConstraint policy.LevelRule) bool {
	return certificateWrapper.IsTrusted() &&
		(certificateWrapper.TrustSunsetDate() == nil || currentTime.Before(*certificateWrapper.TrustSunsetDate()) ||
			!certificateSunsetDateCheckEnforced(certificateSunsetDateConstraint))
}

// certificateSunsetDateCheckEnforced ports the private
// certificateSunsetDateCheckEnforced(LevelRule).
func certificateSunsetDateCheckEnforced(constraint policy.LevelRule) bool {
	return constraint != nil && enumerations.LevelFail == constraint.Level()
}

// IsRevocationDataAcceptable verifies if a revocation data is acceptable for the
// given certificate according to the validation performed within bbb. Port of
// isRevocationDataAcceptable(XmlBasicBuildingBlocks, CertificateWrapper, RevocationWrapper).
func IsRevocationDataAcceptable(bbb *jaxb.XmlBasicBuildingBlocks, certificate *diagnostic.CertificateWrapper,
	revocationData *diagnostic.RevocationWrapper) bool {
	xmlRAC := GetRevocationAcceptanceCheckerResult(bbb, certificate.Id(), revocationData.Id())
	return xmlRAC != nil && xmlRAC.Conclusion != nil &&
		enumerations.IndicationPassed == xmlRAC.Conclusion.Indication.Indication()
}

// IsLongTermAvailabilityAndIntegrityMaterialPresent verifies if the signature
// contains long-term availability and integrity material within its structure.
// Port of isLongTermAvailabilityAndIntegrityMaterialPresent(SignatureWrapper).
func IsLongTermAvailabilityAndIntegrityMaterialPresent(signature *diagnostic.SignatureWrapper) bool {
	return signature.IsThereALevel() || timestampCoveringOtherSignatureTimestampsPresent(signature) ||
		utils.IsCollectionNotEmpty(signature.EvidenceRecords())
}

// timestampCoveringOtherSignatureTimestampsPresent ports the private
// timestampCoveringOtherSignatureTimestampsPresent(SignatureWrapper).
func timestampCoveringOtherSignatureTimestampsPresent(signature *diagnostic.SignatureWrapper) bool {
	for _, timestamp := range signature.TimestampList() {
		timestampedTimestamps := timestamp.TimestampedTimestamps()
		if utils.IsCollectionNotEmpty(timestampedTimestamps) {
			for _, timestampedTimestamp := range timestampedTimestamps {
				if !timestampedTimestamp.Type().IsContentTimestamp() &&
					utils.IsCollectionNotEmpty(timestampedTimestamp.TimestampedSignatures()) {
					return true
				}
			}
		}
	}
	return false
}

// GetRevocationAcceptanceCheckerResult returns a corresponding XmlRAC result for
// the given certificate and revocationData. Port of
// getRevocationAcceptanceCheckerResult(XmlBasicBuildingBlocks, String, String).
func GetRevocationAcceptanceCheckerResult(bbb *jaxb.XmlBasicBuildingBlocks, certificateId string,
	revocationDataId string) *jaxb.XmlRAC {
	if bbb != nil {
		xcv := bbb.XCV
		if xcv != nil {
			subXCV := getXmlSubXCVForId(xcv.SubXCV, certificateId)
			if subXCV != nil {
				crs := subXCV.CRS
				if crs != nil {
					racs := crs.RAC
					rac := getXmlRACForId(racs, revocationDataId)
					if rac != nil {
						return rac
					}
				}
			}
		}
	}
	return nil
}

// getXmlSubXCVForId ports the private getXmlSubXCVForId(List, String).
func getXmlSubXCVForId(subXCVs []*jaxb.XmlSubXCV, tokenId string) *jaxb.XmlSubXCV {
	for _, subXCV := range subXCVs {
		if tokenId == subXCV.Id {
			return subXCV
		}
	}
	return nil
}

// getXmlRACForId ports the private getXmlRACForId(List, String).
func getXmlRACForId(racs []*jaxb.XmlRAC, tokenId string) *jaxb.XmlRAC {
	if utils.IsCollectionNotEmpty(racs) {
		for _, rac := range racs {
			if rac.Id != nil && tokenId == *rac.Id {
				return rac
			}
		}
	}
	return nil
}

// GetFormattedDate returns a formatted String representation of a given Date.
// Port of getFormattedDate(Date).
//
// Java returns null for a null Date; the Go port returns the empty string, so
// that the result stays directly usable as an I18nProvider argument (a *string
// there would render as a pointer, and a nil interface as "<nil>" rather than
// Java's "null").
func GetFormattedDate(date *time.Time) string {
	if date != nil {
		return date.UTC().Format(dateFormat)
	}
	return ""
}

// BuildStringMessage builds a String message from the provided messageTag. Port
// of buildStringMessage(I18nProvider, MessageTag, Object...): Java's null result
// (no message tag defined) is nil here, since the callers propagate it into
// members where absent and empty differ.
func BuildStringMessage(i18nProvider *i18n.I18nProvider, messageTag i18n.MessageTag, args ...interface{}) *string {
	if messageTag != "" {
		message := i18nProvider.GetMessage(messageTag, args...)
		return &message
	}
	return nil
}

// GetCryptoPosition returns the message tag for the given context (signature
// creation,...). Port of getCryptoPosition(Context).
func GetCryptoPosition(context enumerations.Context) (i18n.MessageTag, error) {
	switch context {
	case enumerations.ContextSignature, enumerations.ContextCounterSignature,
		enumerations.ContextKeyBindingSignature:
		return i18n.MessageTag_ACCM_POS_SIG_SIG, nil
	case enumerations.ContextTimestamp:
		return i18n.MessageTag_ACCM_POS_TST_SIG, nil
	case enumerations.ContextRevocation:
		return i18n.MessageTag_ACCM_POS_REVOC_SIG, nil
	case enumerations.ContextCertificate:
		return i18n.MessageTag_ACCM_POS_CERT_CHAIN, nil
	case enumerations.ContextEvidenceRecord:
		return i18n.MessageTag_ACCM_POS_EV_RECORD, nil
	case enumerations.ContextEAA:
		return i18n.MessageTag_ACCM_POS_EAA, nil
	case enumerations.ContextEAARevocation:
		return i18n.MessageTag_ACCM_POS_EAA, nil
	default:
		return "", fmt.Errorf("Unsupported context %s", context)
	}
}

// GetCertificateChainCryptoPosition returns the message tag for the certificate
// chain of the given context. Port of getCertificateChainCryptoPosition(Context).
func GetCertificateChainCryptoPosition(context enumerations.Context) (i18n.MessageTag, error) {
	switch context {
	case enumerations.ContextSignature, enumerations.ContextCounterSignature,
		enumerations.ContextKeyBindingSignature:
		return i18n.MessageTag_ACCM_POS_CERT_CHAIN_SIG, nil
	case enumerations.ContextTimestamp:
		return i18n.MessageTag_ACCM_POS_CERT_CHAIN_TST, nil
	case enumerations.ContextRevocation:
		return i18n.MessageTag_ACCM_POS_CERT_CHAIN_REVOC, nil
	case enumerations.ContextEAARevocation:
		return i18n.MessageTag_ACCM_POS_CERT_CHAIN_EAA_REV, nil
	case enumerations.ContextCertificate:
		return i18n.MessageTag_ACCM_POS_CERT_CHAIN, nil
	default:
		return "", fmt.Errorf("Unsupported context %s", context)
	}
}

// GetDigestMatcherCryptoPosition returns crypto position MessageTag for the
// given XmlDigestMatcher. Port of getDigestMatcherCryptoPosition(XmlDigestMatcher).
//
// A digest matcher with no type - which Java would fail on with a
// NullPointerException in the switch - falls into the unsupported-type error.
func GetDigestMatcherCryptoPosition(digestMatcher *diagnosticjaxb.XmlDigestMatcher) (i18n.MessageTag, error) {
	switch digestMatcherTypeOf(digestMatcher) {
	case enumerations.DigestMatcherTypeObject, enumerations.DigestMatcherTypeReference,
		enumerations.DigestMatcherTypeXPointer:
		return i18n.MessageTag_ACCM_POS_REF, nil
	case enumerations.DigestMatcherTypeManifest:
		return i18n.MessageTag_ACCM_POS_MAN, nil
	case enumerations.DigestMatcherTypeManifestEntry:
		return i18n.MessageTag_ACCM_POS_MAN_ENT, nil
	case enumerations.DigestMatcherTypeSignedProperties:
		return i18n.MessageTag_ACCM_POS_SIGND_PRT, nil
	case enumerations.DigestMatcherTypeKeyInfo:
		return i18n.MessageTag_ACCM_POS_KEY, nil
	case enumerations.DigestMatcherTypeSignatureProperties:
		return i18n.MessageTag_ACCM_POS_SIGNTR_PRT, nil
	case enumerations.DigestMatcherTypeCounterSignature,
		enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return i18n.MessageTag_ACCM_POS_CNTR_SIG, nil
	case enumerations.DigestMatcherTypeMessageDigest:
		return i18n.MessageTag_ACCM_POS_MES_DIG, nil
	case enumerations.DigestMatcherTypeContentDigest:
		return i18n.MessageTag_ACCM_POS_CON_DIG, nil
	case enumerations.DigestMatcherTypeJWSSigningInput:
		return i18n.MessageTag_ACCM_POS_JWS, nil
	case enumerations.DigestMatcherTypeCoseSigStructure:
		return i18n.MessageTag_ACCM_POS_COSE, nil
	case enumerations.DigestMatcherTypeSigDEntry:
		return i18n.MessageTag_ACCM_POS_SIG_D_ENT, nil
	case enumerations.DigestMatcherTypeMessageImprint:
		return i18n.MessageTag_ACCM_POS_MESS_IMP, nil
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveObject:
		return i18n.MessageTag_ACCM_POS_ER_ADO, nil
	case enumerations.DigestMatcherTypeEvidenceRecordOrphanReference:
		return i18n.MessageTag_ACCM_POS_ER_OR, nil
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		return i18n.MessageTag_ACCM_POS_ER_TST, nil
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		return i18n.MessageTag_ACCM_POS_ER_TST_SEQ, nil
	case enumerations.DigestMatcherTypeEvidenceRecordMasterSignature:
		return i18n.MessageTag_ACCM_POS_ER_MST_SIG, nil
	case enumerations.DigestMatcherTypeEAADisclosure:
		return i18n.MessageTag_ACCM_POS_EAA_SD, nil
	case enumerations.DigestMatcherTypeEAANestedDisclosure:
		return i18n.MessageTag_ACCM_POS_EAA_NSD, nil
	case enumerations.DigestMatcherTypeEAAOrphanSelectivelyDisclosableClaim:
		return i18n.MessageTag_ACCM_POS_EAA_OSDC, nil
	case enumerations.DigestMatcherTypeEAAKeyBinding:
		return i18n.MessageTag_ACCM_POS_EAA_KB, nil
	default:
		return "", fmt.Errorf("The provided DigestMatcherType '%s' is not supported!",
			digestMatcherTypeOf(digestMatcher))
	}
}

// GetDigestMatchersCryptoPosition returns crypto position MessageTag for the
// given collection of XmlDigestMatchers. Port of the
// getDigestMatcherCryptoPosition(Collection) overload, renamed since Go has no
// overloading.
func GetDigestMatchersCryptoPosition(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) (i18n.MessageTag, error) {
	if utils.IsCollectionEmpty(digestMatchers) {
		return "", fmt.Errorf("Collection of DigestMatchers cannot be null!")
	} else if utils.CollectionSize(digestMatchers) == 1 {
		return GetDigestMatcherCryptoPosition(digestMatchers[0])
	} else {
		// if more than 1 digest matcher
		digestMatcherType := getDigestMatcherType(digestMatchers)
		switch digestMatcherType {
		case enumerations.DigestMatcherTypeObject, enumerations.DigestMatcherTypeReference,
			enumerations.DigestMatcherTypeXPointer:
			return i18n.MessageTag_ACCM_POS_REF_PL, nil
		case enumerations.DigestMatcherTypeManifest:
			return i18n.MessageTag_ACCM_POS_MAN_PL, nil
		case enumerations.DigestMatcherTypeManifestEntry:
			return i18n.MessageTag_ACCM_POS_MAN_ENT_PL, nil
		case enumerations.DigestMatcherTypeSignedProperties:
			return i18n.MessageTag_ACCM_POS_SIGND_PRT, nil
		case enumerations.DigestMatcherTypeKeyInfo:
			return i18n.MessageTag_ACCM_POS_KEY_PL, nil
		case enumerations.DigestMatcherTypeSignatureProperties:
			return i18n.MessageTag_ACCM_POS_SIGNTR_PRT, nil
		case enumerations.DigestMatcherTypeCounterSignature,
			enumerations.DigestMatcherTypeCounterSignedSignatureValue:
			return i18n.MessageTag_ACCM_POS_CNTR_SIG_PL, nil
		case enumerations.DigestMatcherTypeSigDEntry:
			return i18n.MessageTag_ACCM_POS_SIG_D_ENT_PL, nil
		case enumerations.DigestMatcherTypeEvidenceRecordArchiveObject:
			return i18n.MessageTag_ACCM_POS_ER_ADO_PL, nil
		case enumerations.DigestMatcherTypeEvidenceRecordOrphanReference:
			return i18n.MessageTag_ACCM_POS_ER_OR_PL, nil
		case enumerations.DigestMatcherTypeEAADisclosure:
			return i18n.MessageTag_ACCM_POS_EAA_SD_PL, nil
		case enumerations.DigestMatcherTypeEAANestedDisclosure:
			return i18n.MessageTag_ACCM_POS_EAA_NSD_PL, nil
		case enumerations.DigestMatcherTypeEAAOrphanSelectivelyDisclosableClaim:
			return i18n.MessageTag_ACCM_POS_EAA_OSDC_PL, nil
		case enumerations.DigestMatcherTypeEAAKeyBinding:
			return i18n.MessageTag_ACCM_POS_EAA_KB, nil
		default:
			return "", fmt.Errorf("The provided DigestMatcherType '%s' is not supported for multiple digest matchers!",
				digestMatcherType)
		}
	}
}

// getDigestMatcherType ports the private getDigestMatcherType(Collection).
func getDigestMatcherType(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) enumerations.DigestMatcherType {
	return digestMatcherTypeOf(digestMatchers[0]) // same position shall be provided
}

// digestMatcherTypeOf reads XmlDigestMatcher#getType(): the generated member is
// a pointer, whose nil is Java's null.
func digestMatcherTypeOf(digestMatcher *diagnosticjaxb.XmlDigestMatcher) enumerations.DigestMatcherType {
	if digestMatcher.Type == nil {
		return ""
	}
	return digestMatcher.Type.DigestMatcherType()
}

// GetTimestampTypeMessageTag returns MessageTag associated with the given
// timestamp type. Port of getTimestampTypeMessageTag(TimestampType).
func GetTimestampTypeMessageTag(timestampType enumerations.TimestampType) (i18n.MessageTag, error) {
	if timestampType.IsContentTimestamp() {
		return i18n.MessageTag_TST_TYPE_CONTENT_TST, nil
	} else if timestampType.IsSignatureTimestamp() {
		return i18n.MessageTag_TST_TYPE_SIGNATURE_TST, nil
	} else if timestampType.IsValidationDataTimestamp() {
		return i18n.MessageTag_TST_TYPE_VD_TST, nil
	} else if timestampType.IsDocumentTimestamp() {
		return i18n.MessageTag_TST_TYPE_DOC_TST, nil
	} else if timestampType.IsContainerTimestamp() {
		return i18n.MessageTag_TST_TYPE_CONTAINER_TST, nil
	} else if timestampType.IsArchivalTimestamp() {
		return i18n.MessageTag_TST_TYPE_ARCHIVE_TST, nil
	} else if timestampType.IsEvidenceRecordTimestamp() {
		return i18n.MessageTag_TST_TYPE_ER_TST, nil
	} else {
		return "", fmt.Errorf("The TimestampType '%s' is not supported!", timestampType)
	}
}

// GetContextPosition returns the message tag for the given context. Port of
// getContextPosition(Context).
//
// Deprecated: since DSS 6.5. To be removed.
func GetContextPosition(context enumerations.Context) (i18n.MessageTag, error) {
	switch context {
	case enumerations.ContextSignature, enumerations.ContextCounterSignature,
		enumerations.ContextKeyBindingSignature, enumerations.ContextCertificate:
		return i18n.MessageTag_SIGNATURE, nil
	case enumerations.ContextTimestamp:
		return i18n.MessageTag_TIMESTAMP, nil
	case enumerations.ContextRevocation:
		return i18n.MessageTag_REVOCATION, nil
	default:
		return "", fmt.Errorf("Unsupported context %s", context)
	}
}

// GetSubContextPosition returns the message tag for the given subContext. Port
// of getSubContextPosition(Context, SubContext).
func GetSubContextPosition(context enumerations.Context, subContext enumerations.SubContext) (i18n.MessageTag, error) {
	switch context {
	case enumerations.ContextCertificate:
		return i18n.MessageTag_CERTIFICATE, nil
	case enumerations.ContextSignature, enumerations.ContextCounterSignature,
		enumerations.ContextKeyBindingSignature:
		switch subContext {
		case enumerations.SubContextSigningCert:
			return i18n.MessageTag_SIGNING_CERTIFICATE, nil
		case enumerations.SubContextCACertificate:
			return i18n.MessageTag_CA_CERTIFICATE, nil
		default:
			return "", fmt.Errorf("Unsupported subContext %s", subContext)
		}
	case enumerations.ContextTimestamp:
		switch subContext {
		case enumerations.SubContextSigningCert:
			return i18n.MessageTag_TIMESTAMP_SIG_CERT, nil
		case enumerations.SubContextCACertificate:
			return i18n.MessageTag_TIMESTAMP_CA_CERT, nil
		default:
			return "", fmt.Errorf("Unsupported subContext %s", subContext)
		}
	case enumerations.ContextRevocation:
		switch subContext {
		case enumerations.SubContextSigningCert:
			return i18n.MessageTag_REVOCATION_SIG_CERT, nil
		case enumerations.SubContextCACertificate:
			return i18n.MessageTag_REVOCATION_CA_CERT, nil
		default:
			return "", fmt.Errorf("Unsupported subContext %s", subContext)
		}
	case enumerations.ContextEAARevocation:
		switch subContext {
		case enumerations.SubContextSigningCert:
			return i18n.MessageTag_EAA_REV_SIG_CERT, nil
		case enumerations.SubContextCACertificate:
			return i18n.MessageTag_EAA_REV_CA_CERT, nil
		default:
			return "", fmt.Errorf("Unsupported subContext %s", subContext)
		}
	default:
		return "", fmt.Errorf("Unsupported context %s", context)
	}
}

// GetValidationTimeMessageTag returns a MessageTag corresponding to the given
// ValidationTime type. Port of getValidationTimeMessageTag(ValidationTime).
func GetValidationTimeMessageTag(validationTime enumerations.ValidationTime) (i18n.MessageTag, error) {
	switch validationTime {
	case enumerations.ValidationTimeBESTSignatureTime:
		return i18n.MessageTag_VT_BEST_SIGNATURE_TIME, nil
	case enumerations.ValidationTimeCertificateIssuanceTime:
		return i18n.MessageTag_VT_CERTIFICATE_ISSUANCE_TIME, nil
	case enumerations.ValidationTimeValidationTime:
		return i18n.MessageTag_VT_VALIDATION_TIME, nil
	case enumerations.ValidationTimeTimestampGenerationTime:
		return i18n.MessageTag_VT_TST_GENERATION_TIME, nil
	case enumerations.ValidationTimeTimestampPOETime:
		return i18n.MessageTag_VT_TST_POE_TIME, nil
	default:
		return "", fmt.Errorf("The validation time [%s] is not supported", validationTime)
	}
}

// GetQWACValidationMessageTag returns a MessageTag corresponding to the given
// QWACProfile. Port of getQWACValidationMessageTag(QWACProfile).
func GetQWACValidationMessageTag(qwacProfile enumerations.QWACProfile) (i18n.MessageTag, error) {
	switch qwacProfile {
	case enumerations.QWACProfileQWAC1:
		return i18n.MessageTag_QWAC1_PROFILE, nil
	case enumerations.QWACProfileQWAC2:
		return i18n.MessageTag_QWAC2_PROFILE, nil
	case enumerations.QWACProfileTLSByQWAC2:
		return i18n.MessageTag_TLS_BY_QWAC2_PROFILE, nil
	default:
		return "", fmt.Errorf("The QWAC profile  [%s] is not supported", qwacProfile)
	}
}

// ToUrnOid transforms the given OID to a URN format as per RFC 3061, e.g. "1.2.3"
// to "urn:oid:1.2.3". Port of toUrnOid(String): Java's null argument and null
// result are the empty string here, the OID being carried as a plain Go string.
func ToUrnOid(oid string) string {
	if oid == "" {
		return ""
	}
	return urnOidPrefix + oid
}

// GetDomainName returns a domain name for any given valid URI. Port of
// getDomainName(String).
func GetDomainName(uri string) string {
	if uri == "" {
		return ""
	}
	return domainNamePattern.ReplaceAllString(uri, "")
}

// ProcessValueCheck checks the value against the list of expected values. Port
// of processValueCheck(String, List).
func ProcessValueCheck(value string, expectedValues []string) bool {
	if utils.IsStringNotEmpty(value) && utils.IsCollectionNotEmpty(expectedValues) {
		return containsString(expectedValues, allValue) || containsString(expectedValues, value)
	}
	return false
}

// ProcessValuesCheck checks the values against the expected values. It returns
// TRUE if at least one of the values is allowed by the list of expected values.
// Port of processValuesCheck(List, List).
func ProcessValuesCheck(values []string, expectedValues []string) bool {
	if utils.IsCollectionNotEmpty(values) {
		for _, value := range values {
			if ProcessValueCheck(value, expectedValues) {
				return true
			}
		}
		return false
	} else {
		return utils.IsCollectionEmpty(expectedValues)
	}
}

// ProcessAllValuesCheck checks the values against the expected values. It
// returns TRUE if all the values are allowed by the list of expected values.
// Port of processAllValuesCheck(List, List).
func ProcessAllValuesCheck(values []string, expectedValues []string) bool {
	if utils.IsCollectionNotEmpty(values) {
		for _, value := range values {
			if !ProcessValueCheck(value, expectedValues) {
				return false
			}
		}
		return true
	} else {
		return utils.IsCollectionEmpty(expectedValues)
	}
}

// ProcessValuesForEachExpectedCheck checks whether values contain all the
// expectedValues. Port of processValuesForEachExpectedCheck(List, List).
func ProcessValuesForEachExpectedCheck(values []string, expectedValues []string) bool {
	if utils.IsCollectionNotEmpty(values) {
		for _, expectedValue := range expectedValues {
			if !ProcessValueCheck(expectedValue, values) {
				return false
			}
		}
		return true
	} else {
		return utils.IsCollectionEmpty(expectedValues)
	}
}

// containsString ports java.util.List#contains(Object) for the String lists this
// class matches against.
func containsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// GetFinalCryptographicValidation returns final cryptographic validation from
// the AOV block. This method returns the first algorithm which is going to
// expire in case of failure, or the first applicable algorithm (which is
// SignatureValue's signature algorithm in most of the cases). In case of a valid
// cryptographic validation, returns the first available entry. Port of
// getFinalCryptographicValidation(XmlAOV).
func GetFinalCryptographicValidation(aov *jaxb.XmlAOV) *jaxb.XmlCryptographicValidation {
	if aov == nil || aov.Conclusion == nil {
		return nil
	}
	if enumerations.IndicationPassed == aov.Conclusion.Indication.Indication() {
		return GetPrimaryCryptographicValidation(aov)
	} else {
		return GetFailCryptographicValidation(aov)
	}
}

// GetFailCryptographicValidation returns final cryptographic validation from the
// AOV block. This method returns the first algorithm which is going to expire in
// case of failure, or the first applicable algorithm (which is SignatureValue's
// signature algorithm in most of the cases). Port of
// getFailCryptographicValidation(XmlAOV).
func GetFailCryptographicValidation(aov *jaxb.XmlAOV) *jaxb.XmlCryptographicValidation {
	if aov == nil {
		return nil
	}
	var result *jaxb.XmlCryptographicValidation
	if aov.SignatureCryptographicValidation != nil {
		result = aov.SignatureCryptographicValidation
	}
	if aov.SignedAttributesValidation != nil &&
		(result == nil || (enumerations.IndicationPassed == result.Conclusion.Indication.Indication() && result.NotAfter == nil) ||
			(aov.SignedAttributesValidation.NotAfter != nil && result.NotAfter != nil &&
				result.NotAfter.Time().After(aov.SignedAttributesValidation.NotAfter.Time()))) {
		result = aov.SignedAttributesValidation
	}
	if aov.DigestMatchersValidation != nil &&
		(result == nil || (enumerations.IndicationPassed == result.Conclusion.Indication.Indication() && result.NotAfter == nil) ||
			(aov.DigestMatchersValidation.NotAfter != nil && result.NotAfter != nil &&
				result.NotAfter.Time().After(aov.DigestMatchersValidation.NotAfter.Time()))) {
		result = aov.DigestMatchersValidation
	}
	if aov.CertificateChainCryptographicValidation != nil &&
		utils.IsCollectionNotEmpty(aov.CertificateChainCryptographicValidation.CertificateCryptographicValidation) {
		for _, cryptographicValidation := range aov.CertificateChainCryptographicValidation.CertificateCryptographicValidation {
			if cryptographicValidation != nil &&
				(result == nil || (enumerations.IndicationPassed == result.Conclusion.Indication.Indication() && result.NotAfter == nil) ||
					(cryptographicValidation.NotAfter != nil && result.NotAfter != nil &&
						result.NotAfter.Time().After(cryptographicValidation.NotAfter.Time()))) {
				result = cryptographicValidation
			}
		}
	}
	return result
}

// GetPrimaryCryptographicValidation returns the first available Cryptographic
// Validation entry. Port of getPrimaryCryptographicValidation(XmlAOV).
func GetPrimaryCryptographicValidation(aov *jaxb.XmlAOV) *jaxb.XmlCryptographicValidation {
	if aov == nil {
		return nil
	}
	var result *jaxb.XmlCryptographicValidation
	if aov.SignatureCryptographicValidation != nil {
		result = aov.SignatureCryptographicValidation
	}
	if result == nil && aov.SignedAttributesValidation != nil {
		result = aov.SignedAttributesValidation
	}
	if result == nil && aov.DigestMatchersValidation != nil {
		result = aov.DigestMatchersValidation
	}
	if result == nil && aov.CertificateChainCryptographicValidation != nil &&
		utils.IsCollectionNotEmpty(aov.CertificateChainCryptographicValidation.CertificateCryptographicValidation) {
		for _, cryptographicValidation := range aov.CertificateChainCryptographicValidation.CertificateCryptographicValidation {
			if cryptographicValidation != nil {
				result = cryptographicValidation
				break
			}
		}
	}
	return result
}

// GetConstraintOrMaxLevel returns the current level with a max limit of the
// maxLevel. Port of getConstraintOrMaxLevel(LevelRule, Level), fall-through of
// the Java switch included.
func GetConstraintOrMaxLevel(constraint policy.LevelRule, maxLevel enumerations.Level) (policy.LevelRule, error) {
	if constraint == nil || maxLevel == "" {
		return nil, nil
	}
	var level enumerations.Level
	switch constraint.Level() {
	case enumerations.LevelFail:
		if enumerations.LevelFail == maxLevel {
			level = enumerations.LevelFail
			break
		}
		fallthrough
	case enumerations.LevelWarn:
		if enumerations.LevelWarn == maxLevel {
			level = enumerations.LevelWarn
			break
		}
		fallthrough
	case enumerations.LevelInform:
		if enumerations.LevelInform == maxLevel {
			level = enumerations.LevelInform
			break
		}
		fallthrough
	case enumerations.LevelIgnore:
		if enumerations.LevelIgnore == maxLevel {
			level = enumerations.LevelIgnore
			break
		}
		level = constraint.Level()
	default:
		return nil, fmt.Errorf("The support of Level '%s' is not implemented!", constraint.Level())
	}

	return GetLevelRule(level), nil
}

// GetLevelRule generates an anonymous implementation of the LevelRule with the
// given Level. Port of getLevelRule(Level).
func GetLevelRule(level enumerations.Level) policy.LevelRule {
	if level == "" {
		return nil
	}
	return levelRule(level)
}

// levelRule is the Go form of the Java lambda "() -> level".
type levelRule enumerations.Level

// Level gets the constraint execution level.
func (l levelRule) Level() enumerations.Level {
	return enumerations.Level(l)
}
