// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/AbstractBasicValidationProcess.java (DSS 6.5.RC1).
//
// The abstract class implementing the "5.3 Validation process for Basic
// Signatures" process. Java subclasses: BasicSignatureValidationProcess
// (unmodified initChain), vpftsp.TimestampBasicValidationProcess (unmodified
// initChain), vpfltvd.RevocationBasicValidationProcess (unmodified initChain),
// and eaa.EAAValidationProcess, which fully replaces initChain() rather than
// extending it. Only getContentTimestamps() and getTimestampValidation() are
// genuinely virtually dispatched from this class' own InitChain (only
// BasicSignatureValidationProcess overrides them), so only those two are
// routed through AbstractBasicValidationProcessOverrides; every other
// protected helper below is a plain method - nothing in this call graph ever
// dispatches into a subclass override of them (EAAValidationProcess's own
// signatureAcceptanceValidation override is called only from
// EAAValidationProcess's own initChain, on its own receiver, which Go's method
// shadowing resolves correctly without any interface indirection).
package vpfbs

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractBasicValidationProcessOverrides declares the overridable protected
// methods of AbstractBasicValidationProcess that its own InitChain calls back
// into, on top of the Chain ones. A concrete process registers itself through
// InitAbstractBasicValidationProcess; every method it does not define is
// supplied by the embedded AbstractBasicValidationProcess through ordinary Go
// method promotion.
type AbstractBasicValidationProcessOverrides interface {
	process.ChainOverrides

	// ContentTimestamps returns a list of content timestamps. Port of
	// getContentTimestamps(), whose default is an empty list.
	ContentTimestamps() []*diagnostic.TimestampWrapper

	// TimestampValidation returns the corresponding validation result for a
	// timestamp with the given Id. Port of getTimestampValidation(String),
	// whose default is null.
	TimestampValidation(timestampId string) *jaxb.XmlValidationProcessBasicTimestamp
}

// AbstractBasicValidationProcess is the abstract class implementing the
// "5.3 Validation process for Basic Signatures" process.
type AbstractBasicValidationProcess[T any] struct {
	*process.ChainBase[T]

	// DiagnosticData is the diagnostic data. Exported because Java declares the
	// field protected.
	DiagnosticData *diagnostic.DiagnosticData

	// Token is the token to be validated. Exported because Java declares the
	// field protected.
	Token diagnostic.TokenProxy

	// BBBs is the map of BasicBuildingBlocks. Exported because Java declares
	// the field protected.
	BBBs map[string]*jaxb.XmlBasicBuildingBlocks

	// overrides points back at the concrete process; see
	// InitAbstractBasicValidationProcess.
	overrides AbstractBasicValidationProcessOverrides
}

// NewAbstractBasicValidationProcess is the common constructor. Port of
// AbstractBasicValidationProcess(I18nProvider, T, DiagnosticData, TokenProxy, Map).
func NewAbstractBasicValidationProcess[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	diagnosticData *diagnostic.DiagnosticData, token diagnostic.TokenProxy,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) *AbstractBasicValidationProcess[T] {
	return &AbstractBasicValidationProcess[T]{
		ChainBase:      process.NewChainBase(i18nProvider, result),
		DiagnosticData: diagnosticData,
		Token:          token,
		BBBs:           bbbs,
	}
}

// InitAbstractBasicValidationProcess registers the concrete process with its
// base so that the base can dispatch to the overridden methods, Chain's
// included. It must be called exactly once, by the concrete process'
// constructor, before Execute.
func (c *AbstractBasicValidationProcess[T]) InitAbstractBasicValidationProcess(overrides AbstractBasicValidationProcessOverrides) {
	c.overrides = overrides
	c.InitChainBase(overrides)
}

// abvpOverrides returns the registered overrides, panicking when the concrete
// process forgot to call InitAbstractBasicValidationProcess.
func (c *AbstractBasicValidationProcess[T]) abvpOverrides() AbstractBasicValidationProcessOverrides {
	if c.overrides == nil {
		panic("AbstractBasicValidationProcess was not initialised: the concrete process must call InitAbstractBasicValidationProcess in its constructor")
	}
	return c.overrides
}

// InitChain initializes the chain. Port of the (non-abstract) initChain().
func (c *AbstractBasicValidationProcess[T]) InitChain() {
	overrides := c.abvpOverrides()

	// 5.3.4 Processing (ETSI TS 119 102-1 V1.2.1)
	tokenBBBs := c.BBBs[c.Token.Id()]
	if tokenBBBs == nil {
		// Java's throw new IllegalStateException(...) becomes a panic: InitChain
		// has no error return (process.ChainOverrides), matching the precedent in
		// eaa_type_check.go.
		panic(fmt.Sprintf("Missing Basic Building Blocks result for token with Id '%s'", c.Token.Id()))
	}

	var item process.ChainItem[T]

	/*
	 * 1) The Basic Signature validation process shall perform the format checking
	 * as per clause 5.2.2. If the process returns PASSED, the Basic Signature
	 * validation process shall continue with the next step. Otherwise, the Basic
	 * Signature validation process shall return the indication FAILED with the
	 * sub-indication FORMAT_FAILURE.
	 */
	xmlFC := tokenBBBs.FC
	if xmlFC != nil {
		item = c.formatChecking(xmlFC)
		c.FirstItem = item
	}

	/*
	 * 2) The Basic Signature validation process shall perform the identification
	 * of the signing certificate (as per clause 5.2.3) with the signature and
	 * the signing certificate, if provided as a parameter. If the identification of
	 * the signing certificate process returns the indication INDETERMINATE with
	 * the sub-indication NO_SIGNING_CERTIFICATE_FOUND, the Basic Signature validation
	 * process shall return the indication INDETERMINATE with the sub-indication
	 * NO_SIGNING_CERTIFICATE_FOUND, otherwise it shall go to the next step.
	 */
	xmlISC := tokenBBBs.ISC
	// required for all tokens
	if c.FirstItem == nil {
		item = c.identificationOfSigningCertificate(xmlISC)
		c.FirstItem = item
	} else {
		item = item.SetNextItem(c.identificationOfSigningCertificate(xmlISC))
	}

	/*
	 * 3) The Basic Signature validation process shall perform the Validation Context Initialization
	 * as per clause 5.2.4. If the process returns INDETERMINATE with some sub-indication,
	 * the Basic Signature validation process shall return the indication INDETERMINATE
	 * together with that sub-indication, otherwise it shall go to the next step.
	 */
	xmlVCI := tokenBBBs.VCI
	if xmlVCI != nil {
		item = item.SetNextItem(c.validationContextInitialization(xmlVCI))
	}

	/*
	 * 4) The Basic Signature validation process shall perform the X.509 Certificate Validation
	 * as per clause 5.2.6 with the following inputs:
	 *    a) The signing certificate obtained in step 2). And
	 *    b) X.509 validation constraints, certificate validation-data and
	 *       cryptographic constraints obtained in step 3) or provided as input.
	 */
	contentTimestamps := overrides.ContentTimestamps()

	x509ValidationStatus := &jaxb.XmlConclusion{}
	xmlXCV := tokenBBBs.XCV
	if xmlXCV != nil {

		item = item.SetNextItem(c.x509CertificateValidation(xmlXCV))

		/*
		 * If the X.509 Certificate Validation process returns the indication PASSED,
		 * the Basic Signature validation process shall set X509_validation-status to PASSED
		 * and it shall go to step 5).
		 */
		if c.IsValid(&xmlXCV.XmlConstraintsConclusionContent) {
			x509ValidationStatus.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)

		} else {
			x509ValidationStatus.Indication = xmlXCV.Conclusion.Indication
			x509ValidationStatus.SubIndication = xmlXCV.Conclusion.SubIndication
			/*
			 * If the X.509 Certificate Validation process returns the indication
			 * INDETERMINATE with the sub-indication REVOKED_NO_POE and if
			 * the signature contains a content-time-stamp attribute, the Basic Signature
			 * validation process shall perform the validation process for AdES time-stamps
			 * as defined in clause 5.4. If this process returns the indication PASSED and
			 * the generation time of the time-stamp token is after the revocation time,
			 * the Basic Signature validation process shall set X509_validation-status to FAILED
			 * with the sub-indication REVOKED. In all other cases, the Basic Signature validation
			 * process shall set X509_validation-status to INDETERMINATE with the sub-indication
			 * REVOKED_NO_POE. The process shall continue with step 5)
			 */
			item = item.SetNextItem(c.signingCertificateNotRevoked(xmlXCV))

			if enumerations.IndicationIndeterminate == xcvIndication(xmlXCV) &&
				enumerations.SubIndicationRevokedNoPOE == xcvSubIndication(xmlXCV) &&
				utils.IsCollectionNotEmpty(contentTimestamps) {
				revocationTime := c.getRevocationTimeForSigningCertificate()

				item = item.SetNextItem(c.contentTimestampsPresent(contentTimestamps))

				for _, timestampWrapper := range contentTimestamps {

					timestampValidation := overrides.TimestampValidation(timestampWrapper.Id())
					if timestampValidation != nil {

						item = item.SetNextItem(c.timestampBasicValidation(timestampWrapper, timestampValidation))

						if c.IsValid(&timestampValidation.XmlConstraintsConclusionContent) {

							item = item.SetNextItem(c.timestampNotAfterRevocationTime(timestampWrapper, revocationTime))

							if revocationTime != nil && timestampWrapper.ProductionTime() != nil &&
								timestampWrapper.ProductionTime().After(*revocationTime) {
								x509ValidationStatus.Indication = jaxb.IndicationValue(enumerations.IndicationFailed)
								setXmlSubIndication(x509ValidationStatus, enumerations.SubIndicationRevoked)
								break
							}
						}
					}

				}

			}

			/*
			 * If the X.509 Certificate Validation process returns the indication INDETERMINATE
			 * with the sub-indication OUT_OF_BOUNDS_NO_POE or OUT_OF_BOUNDS_NOT_REVOKED, and if
			 * the signature contains a content-time-stamp attribute, the Basic Signature
			 * validation process shall perform the validation process for AdES time-stamps as defined
			 * in clause 5.4. If it returns the indication PASSED and the generation time of
			 * the time-stamp token is after the expiration date of the signing certificate,
			 * the Basic Signature validation process shall set X509_validation-status to FAILED with
			 * the sub-indication EXPIRED. Otherwise, the Basic Signature validation process shall set
			 * X509_validation-status to INDETERMINATE with the sub-indication OUT_OF_BOUNDS_NO_POE or
			 * OUT_OF_BOUNDS_NOT_REVOKED, respectively. The process shall continue with step 5).
			 */
			item = item.SetNextItem(c.validationTimeAtValidityRange(xmlXCV))

			if enumerations.IndicationIndeterminate == xcvIndication(xmlXCV) &&
				(enumerations.SubIndicationOutOfBoundsNoPOE == xcvSubIndication(xmlXCV) ||
					enumerations.SubIndicationOutOfBoundsNotRevoked == xcvSubIndication(xmlXCV)) &&
				utils.IsCollectionNotEmpty(contentTimestamps) {
				var certificateNotAfter *time.Time
				if signingCertificate := c.Token.SigningCertificate(); signingCertificate != nil {
					certificateNotAfter = signingCertificate.NotAfter()
				}

				item = item.SetNextItem(c.contentTimestampsPresent(contentTimestamps))

				for _, timestampWrapper := range contentTimestamps {

					timestampValidation := overrides.TimestampValidation(timestampWrapper.Id())
					if timestampValidation != nil {

						item = item.SetNextItem(c.timestampBasicValidation(timestampWrapper, timestampValidation))

						if c.IsValid(&timestampValidation.XmlConstraintsConclusionContent) {

							item = item.SetNextItem(c.timestampNotAfterSigningCertificateNotAfterTime(timestampWrapper, certificateNotAfter))

							if certificateNotAfter != nil && timestampWrapper.ProductionTime() != nil &&
								timestampWrapper.ProductionTime().After(*certificateNotAfter) {
								x509ValidationStatus.Indication = jaxb.IndicationValue(enumerations.IndicationFailed)
								setXmlSubIndication(x509ValidationStatus, enumerations.SubIndicationExpired)
								break
							}
						}
					}

				}
			}

			/*
			 * If the X.509 Certificate Validation process returns the indication INDETERMINATE
			 * with the sub-indication NO_CERTIFICATE_CHAIN_FOUND and if the signature algorithm
			 * requires the full certificate chain for determining the public key, the Basic Signature
			 * validation process shall return the indication INDETERMINATE with the sub-indication
			 * NO_CERTIFICATE_CHAIN_FOUND.
			 *
			 * In all other cases, the Basic Signature validation process shall set X509_validation-status
			 * to the indication and sub-indication returned by the X.509 Certificate Validation process
			 * and continue with step 5).
			 */

		}
	}

	/*
	 * 5) The Basic Signature validation process shall perform the Cryptographic Verification
	 * process as per clause 5.2.7 with the following inputs:
	 *    a) The signed data object.
	 *    b) The signing certificate obtained in step 2).
	 *    c) The certificate chain returned in the previous step, if it was returned in step 4). And
	 *    d) The SD or SDR, if given in the input.
	 */
	xmlCV := tokenBBBs.CV
	if xmlCV != nil {

		item = item.SetNextItem(c.cryptographicVerification(xmlCV))

		/*
		 * If the Cryptographic Verification process returns PASSED:
		 */
		if c.IsValid(&xmlCV.XmlConstraintsConclusionContent) {
			/*
			 * a) If the X509_validation-status set in the previous step contains the indication PASSED,
			 * the Basic Signature validation process shall go to the next step;
			 */
			// continue

			/*
			 * b) If the X509_validation-status set in the previous step contains the indication
			 * INDETERMINATE or FAILED with any subindication, the Basic Signature validation process
			 * shall return the indication and subindication contained in X509_validation-status,
			 * with any associated information about the reason.
			 */
			if enumerations.IndicationIndeterminate == x509ValidationStatus.Indication.Indication() ||
				enumerations.IndicationFailed == x509ValidationStatus.Indication.Indication() {

				item = item.SetNextItem(c.basicValidationProcess(x509ValidationStatus))

			}

			/*
			 * Otherwise, the Basic Signature validation process shall return the returned indication,
			 * sub-indication and associated information provided by the Cryptographic Verification process.
			 */
			// returned before
		}

	}

	/*
	 * 6) The Basic Signature validation process shall perform the Signature Acceptance Validation
	 * process as per clause 5.2.8 with the following inputs:
	 *    a) the Signed Data Object(s);
	 *    b) the certificate chain obtained in step 4);
	 *    c) the Cryptographic Constraints; and
	 *    d) the Signature Elements Constraints.
	 */
	xmlSAV := tokenBBBs.SAV
	if xmlSAV != nil {

		item = item.SetNextItem(c.signatureAcceptanceValidation(xmlSAV))

		xmlAOV := tokenBBBs.AOV
		/*
		 * If the signature acceptance validation process returns PASSED, the Basic Signature validation
		 * process shall go to the next step.
		 */
		if enumerations.IndicationPassed == xmlSAV.Conclusion.Indication.Indication() {
			// continue

			/*
			 * If the signature acceptance validation process returns the indication INDETERMINATE
			 * with the sub-indication CRYPTO_CONSTRAINTS_FAILURE_NO_POE and the material concerned by
			 * this failure is the signature value and if the signature contains a content-time-stamp attribute,
			 * the Basic Signature validation process shall perform the validation process for AdES time-stamps
			 * as defined in clause 5.4. If it returns the indication PASSED and the algorithm(s) concerned
			 * were no longer considered reliable at the generation time of the time-stamp token,
			 * the Basic Signature validation process shall return the indication INDETERMINATE with
			 * the sub-indication CRYPTO_CONSTRAINTS_FAILURE. In all other cases, the Basic Signature
			 * validation process shall return the indication INDETERMINATE with the sub-indication
			 * CRYPTO_CONSTRAINTS_FAILURE_NO_POE.
			 */
		} else if enumerations.IndicationIndeterminate == xmlSAV.Conclusion.Indication.Indication() &&
			savSubIndication(xmlSAV) == enumerations.SubIndicationCryptoConstraintsFailureNoPOE &&
			c.isSignatureValueConcernedByFailure(xmlAOV) && utils.IsCollectionNotEmpty(contentTimestamps) {

			item = item.SetNextItem(c.contentTimestampsPresent(contentTimestamps))

			for _, timestampWrapper := range contentTimestamps {

				timestampValidation := overrides.TimestampValidation(timestampWrapper.Id())
				if timestampValidation != nil {

					item = item.SetNextItem(c.timestampBasicValidation(timestampWrapper, timestampValidation))

					if c.IsValid(&timestampValidation.XmlConstraintsConclusionContent) {

						item = item.SetNextItem(c.timestampNotAfterCryptographicAlgorithmsExpiration(
							timestampWrapper, xmlAOV.SignatureCryptographicValidation))

					}
				}

			}
		}

		if !c.IsValid(&xmlSAV.XmlConstraintsConclusionContent) {
			item = item.SetNextItem(c.basicValidationProcess(xmlSAV.Conclusion)) //nolint:staticcheck // mirrors upstream AbstractBasicValidationProcess#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
		}

	}

	/*
	 * 7) The Basic Signature validation process shall return the success indication PASSED
	 * together with the certificate chain obtained in step 4). In addition, the Basic Signature
	 * validation process should return additional information extracted from the signature and/or
	 * used by the intermediate steps. In particular, the SVA should provide to the DA all information
	 * related to signed and unsigned attributes, including those which were not processed during
	 * the validation process.
	 */
}

// xcvIndication reads xmlXCV.Conclusion.Indication.
func xcvIndication(xmlXCV *jaxb.XmlXCV) enumerations.Indication {
	return xmlXCV.Conclusion.Indication.Indication()
}

// xcvSubIndication reads xmlXCV.Conclusion.SubIndication; nil is Java's null.
func xcvSubIndication(xmlXCV *jaxb.XmlXCV) enumerations.SubIndication {
	if xmlXCV.Conclusion.SubIndication == nil {
		return ""
	}
	return xmlXCV.Conclusion.SubIndication.SubIndication()
}

// savSubIndication reads xmlSAV.Conclusion.SubIndication; nil is Java's null.
func savSubIndication(xmlSAV *jaxb.XmlSAV) enumerations.SubIndication {
	if xmlSAV.Conclusion.SubIndication == nil {
		return ""
	}
	return xmlSAV.Conclusion.SubIndication.SubIndication()
}

// setXmlSubIndication ports XmlConclusion#setSubIndication for the locally
// built x509ValidationStatus conclusion.
func setXmlSubIndication(conclusion *jaxb.XmlConclusion, subIndication enumerations.SubIndication) {
	value := jaxb.SubIndicationValue(subIndication)
	conclusion.SubIndication = &value
}

// formatChecking executes "5.2.2 Format Checking" building block for the given
// token. Port of formatChecking(XmlFC).
func (c *AbstractBasicValidationProcess[T]) formatChecking(xmlFC *jaxb.XmlFC) process.ChainItem[T] {
	return NewFormatCheckingResultCheck(c.I18nProvider, c.Result, xmlFC, c.Token, c.FailLevelRule())
}

// identificationOfSigningCertificate executes "5.2.3 Identification of the
// signing certificate" building block for the given token. Port of
// identificationOfSigningCertificate(XmlISC).
func (c *AbstractBasicValidationProcess[T]) identificationOfSigningCertificate(xmlISC *jaxb.XmlISC) process.ChainItem[T] {
	return NewIdentificationOfSigningCertificateResultCheck(c.I18nProvider, c.Result, xmlISC, c.Token, c.FailLevelRule())
}

// validationContextInitialization executes "5.2.4 Validation context
// initialization" building block for the given token. Port of
// validationContextInitialization(XmlVCI).
func (c *AbstractBasicValidationProcess[T]) validationContextInitialization(xmlVCI *jaxb.XmlVCI) process.ChainItem[T] {
	return NewValidationContextInitializationResultCheck(c.I18nProvider, c.Result, xmlVCI, c.Token, c.FailLevelRule())
}

// x509CertificateValidation executes "5.2.6 X.509 certificate validation"
// building block for the given token. Port of x509CertificateValidation(XmlXCV).
func (c *AbstractBasicValidationProcess[T]) x509CertificateValidation(xmlXCV *jaxb.XmlXCV) process.ChainItem[T] {
	return NewX509CertificateValidationResultCheck(c.I18nProvider, c.Result, xmlXCV, c.Token, c.WarnLevelRule())
}

// signingCertificateNotRevoked ports the private signingCertificateNotRevoked(XmlXCV).
func (c *AbstractBasicValidationProcess[T]) signingCertificateNotRevoked(xmlXCV *jaxb.XmlXCV) process.ChainItem[T] {
	return NewSigningCertificateNotRevokedCheck(c.I18nProvider, c.Result, xmlXCV, c.Token, c.WarnLevelRule())
}

// validationTimeAtValidityRange ports the private validationTimeAtValidityRange(XmlXCV).
func (c *AbstractBasicValidationProcess[T]) validationTimeAtValidityRange(xmlXCV *jaxb.XmlXCV) process.ChainItem[T] {
	return NewValidationTimeAtCertificateValidityRangeCheck(c.I18nProvider, c.Result, xmlXCV, c.Token, c.WarnLevelRule())
}

// contentTimestampsPresent ports the private contentTimestampsPresent(List).
func (c *AbstractBasicValidationProcess[T]) contentTimestampsPresent(contentTimestamps []*diagnostic.TimestampWrapper) process.ChainItem[T] {
	return NewContentTimestampsCheck(c.I18nProvider, c.Result, contentTimestamps, c.WarnLevelRule())
}

// timestampBasicValidation ports the private
// timestampBasicValidation(TimestampWrapper, XmlValidationProcessBasicTimestamp).
func (c *AbstractBasicValidationProcess[T]) timestampBasicValidation(timestamp *diagnostic.TimestampWrapper,
	timestampValidation *jaxb.XmlValidationProcessBasicTimestamp) process.ChainItem[T] {
	return NewBasicTimestampValidationWithIdCheck(c.I18nProvider, c.Result, timestamp, timestampValidation, c.WarnLevelRule())
}

// timestampNotAfterRevocationTime ports the private
// timestampNotAfterRevocationTime(TimestampWrapper, Date).
func (c *AbstractBasicValidationProcess[T]) timestampNotAfterRevocationTime(timestamp *diagnostic.TimestampWrapper,
	revocationTime *time.Time) process.ChainItem[T] {
	return NewTimestampGenerationTimeNotAfterRevocationTimeCheck(c.I18nProvider, c.Result, timestamp, revocationTime, c.WarnLevelRule())
}

// timestampNotAfterSigningCertificateNotAfterTime ports the private
// timestampNotAfterSigningCertificateNotAfterTime(TimestampWrapper, Date).
func (c *AbstractBasicValidationProcess[T]) timestampNotAfterSigningCertificateNotAfterTime(timestamp *diagnostic.TimestampWrapper,
	certificateNotAfter *time.Time) process.ChainItem[T] {
	return NewTimestampGenerationTimeNotAfterCertificateExpirationCheck(c.I18nProvider, c.Result, timestamp, certificateNotAfter, c.WarnLevelRule())
}

// cryptographicVerification executes "5.2.7 Cryptographic verification"
// building block for the given token. Port of cryptographicVerification(XmlCV).
func (c *AbstractBasicValidationProcess[T]) cryptographicVerification(xmlCV *jaxb.XmlCV) process.ChainItem[T] {
	return NewCryptographicVerificationResultCheck(c.I18nProvider, c.Result, xmlCV, c.Token, c.FailLevelRule())
}

// signatureAcceptanceValidation executes "5.2.8 Signature Acceptance
// Validation (SAV)" building block for the given token. Port of
// signatureAcceptanceValidation(XmlSAV).
func (c *AbstractBasicValidationProcess[T]) signatureAcceptanceValidation(xmlSAV *jaxb.XmlSAV) process.ChainItem[T] {
	return NewSignatureAcceptanceValidationResultCheck(c.I18nProvider, c.Result, xmlSAV, c.Token, c.WarnLevelRule())
}

// timestampNotAfterCryptographicAlgorithmsExpiration ports the private
// timestampNotAfterCryptographicAlgorithmsExpiration(TimestampWrapper, XmlCryptographicValidation).
func (c *AbstractBasicValidationProcess[T]) timestampNotAfterCryptographicAlgorithmsExpiration(
	timestamp *diagnostic.TimestampWrapper, cryptographicValidation *jaxb.XmlCryptographicValidation) process.ChainItem[T] {
	return NewTimestampGenerationTimeNotAfterCryptographicConstraintsExpirationCheck(c.I18nProvider, c.Result,
		timestamp, cryptographicValidation, c.FailLevelRule())
}

// basicValidationProcess executes a final validation check of the
// "5.3 Validation process for Basic Signatures" block. Port of
// basicValidationProcess(XmlConclusion).
func (c *AbstractBasicValidationProcess[T]) basicValidationProcess(xmlConclusion *jaxb.XmlConclusion) process.ChainItem[T] {
	return NewBasicValidationProcessCheck(c.I18nProvider, c.Result, xmlConclusion, c.Token, c.FailLevelRule())
}

// ContentTimestamps returns a list of content timestamps. Port of
// getContentTimestamps(), whose default is an empty list.
func (c *AbstractBasicValidationProcess[T]) ContentTimestamps() []*diagnostic.TimestampWrapper {
	return nil
}

// TimestampValidation gets the corresponding validation result for a
// timestamp with the given Id. Port of getTimestampValidation(String), whose
// default is null.
func (c *AbstractBasicValidationProcess[T]) TimestampValidation(timestampId string) *jaxb.XmlValidationProcessBasicTimestamp {
	return nil
}

// getRevocationTimeForSigningCertificate ports the private
// getRevocationTimeForSigningCertificate().
func (c *AbstractBasicValidationProcess[T]) getRevocationTimeForSigningCertificate() *time.Time {
	signingCertificate := c.Token.SigningCertificate()
	if signingCertificate != nil && utils.IsCollectionNotEmpty(signingCertificate.CertificateRevocationData()) {
		latest := c.DiagnosticData.LatestRevocationDataForCertificate(signingCertificate)
		if latest != nil {
			return latest.RevocationDate()
		}
	}
	return nil
}

// isSignatureValueConcernedByFailure ports the private
// isSignatureValueConcernedByFailure(XmlAOV).
func (c *AbstractBasicValidationProcess[T]) isSignatureValueConcernedByFailure(xmlAOV *jaxb.XmlAOV) bool {
	if xmlAOV == nil {
		return false
	}
	cryptographicValidation := xmlAOV.SignatureCryptographicValidation
	return cryptographicValidation != nil && !c.IsValidConclusion(cryptographicValidation.Conclusion)
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of the overridden
// collectMessages(XmlConclusion, XmlConstraint).
func (c *AbstractBasicValidationProcess[T]) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	if constraint.BlockType != nil && jaxb.XmlBlockTypeCNTTSTBBB == *constraint.BlockType {
		if constraint.Error != nil {
			conclusion.Errors = append(conclusion.Errors, constraint.Error)
		}
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}

// CollectAdditionalMessages fills additional messages into the conclusion.
// Port of the overridden collectAdditionalMessages(XmlConclusion).
func (c *AbstractBasicValidationProcess[T]) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	tokenBBBs := c.BBBs[c.Token.Id()]
	if tokenBBBs != nil {
		conclusion.Errors = nil
		conclusion.Errors = append(conclusion.Errors, tokenBBBs.Conclusion.Errors...)
		conclusion.Warnings = nil
		conclusion.Warnings = append(conclusion.Warnings, tokenBBBs.Conclusion.Warnings...)
		conclusion.Infos = nil
		conclusion.Infos = append(conclusion.Infos, tokenBBBs.Conclusion.Infos...)

		for _, constraint := range c.Result.Constraint() {
			c.abvpOverrides().CollectMessages(conclusion, constraint)
		}
	}
}
