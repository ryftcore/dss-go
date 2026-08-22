// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/ValidationProcessForSignaturesWithLongTermValidationData.java (DSS 6.5.RC1).
//
// 5.5 Validation process for Signatures with Time and Signatures with
// Long-Term Validation Data.
//
// PACKAGE-BOUNDARY DEVIATION (LTVB, phase 8e): Java's ValidationProcessForSignaturesWithLongTermValidationData
// lives in eu.europa.esig.dss.validation.process.vpfltvd, the same package as
// its checks and as LongTermValidationCertificateRevocationSelector /
// RevocationBasicValidationProcess. It is filed in this separate package
// instead, because it is the only vpfltvd-family class that also needs
// bbb/sav.TLevelTimeStampCheck and LTALevelTimeStampCheck (tLevelTimeStamp()
// / ltaLevelTimeStamp() below) - and bbb/sav (frozen, phase 8c) already
// imports package vpfltvd unconditionally for
// vpfltvd.NewTimestampMessageImprintWithIdCheck (see
// bbb/sav/signature_acceptance_validation.go's own header). Since Go forbids
// package import cycles outright and bbb/sav cannot be edited, package
// vpfltvd can never import bbb/sav; this file - and only this file - is
// filed under vpfltvdsig so it can import both vpfltvd (for its own sibling
// checks: AcceptableBasicSignatureValidationCheck,
// LongTermValidationCertificateRevocationSelector, the BestSignatureTimeXxx
// family, CertificateKnownToBeNotRevokedXxx,
// RevocationDateAfterBestSignatureTimeCheck, TimestampCoherenceOrderCheck,
// SigningTimeAttributePresentCheck, TimestampDelayCheck,
// TimestampMessageImprintWithIdCheck) and bbb/sav (TLevelTimeStampCheck,
// LTALevelTimeStampCheck). The same relocation technique vpfbs uses for
// BasicTimestampValidationCheck/WithIdCheck (see
// vpfbs/basic_timestamp_validation_check.go) for the mirror-image cycle
// between vpfbs and vpftsp.
//
// certificateRevocationMap is Java's LinkedHashMap<CertificateWrapper,
// CertificateRevocationWrapper>: insertion order (the order the signature's
// certificate chain is walked in InitChain) is iterated directly in
// revocationDateAfterBestSignatureTimeValidation and revocationDataReliableAtTime,
// so it is order-sensitive output and cannot become a plain Go map (random
// iteration order). It is kept as a parallel ordered-keys slice plus a lookup
// map, both keyed by CertificateWrapper pointer identity - the same identity
// comparison bbb/xcv's X509CertificateValidation already relies on for
// "trustAnchor == certificate".
package vpfltvdsig

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/aov"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/sav"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/xcv"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfbs"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfltvd"
)

// ValidationProcessForSignaturesWithLongTermValidationData is 5.5 Validation
// process for Signatures with Time and Signatures with Long-Term Validation
// Data.
type ValidationProcessForSignaturesWithLongTermValidationData struct {
	*process.ChainBase[*jaxb.XmlValidationProcessLongTermData]

	// basicSignatureValidation is the basic signature validation conclusion.
	basicSignatureValidation *jaxb.XmlConstraintsConclusionContent

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.DiagnosticData

	// currentSignature is the signature.
	currentSignature *diagnostic.SignatureWrapper

	// bbbs is the map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// xmlTimestamps is the list of timestamps.
	xmlTimestamps []*jaxb.XmlTimestamp

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// currentDate is the validation time.
	currentDate time.Time

	// certificateRevocationMap defines the map between certificates in the
	// chain and their latest valid revocation data.
	// Keyed by the certificate's ID, not by wrapper pointer: Java's
	// CertificateWrapper (AbstractTokenProxy) overrides equals/hashCode by
	// getId(), and every CertificateWrapper this class hands to the map -
	// currentSignature.getCertificateChain(), getSigningCertificate() - is a
	// FRESH wrapper object in the ported diagnostic package, so a
	// pointer-identity key would never match on lookup.
	certificateRevocationMap map[string]*diagnostic.CertificateRevocationWrapper

	// certificateRevocationOrder is the insertion order of
	// certificateRevocationMap's keys (see the package/file header).
	certificateRevocationOrder []*diagnostic.CertificateWrapper
}

// NewValidationProcessForSignaturesWithLongTermValidationData is the default
// constructor. Port of
// ValidationProcessForSignaturesWithLongTermValidationData(I18nProvider, XmlSignature, DiagnosticData, SignatureWrapper, Map, ValidationPolicy, Date).
func NewValidationProcessForSignaturesWithLongTermValidationData(i18nProvider *i18n.I18nProvider,
	signatureAnalysis *jaxb.XmlSignature, diagnosticData *diagnostic.DiagnosticData, currentSignature *diagnostic.SignatureWrapper,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, validationPolicy policy.ValidationPolicy,
	currentDate time.Time) *ValidationProcessForSignaturesWithLongTermValidationData {
	xmlResult := &jaxb.XmlValidationProcessLongTermData{}
	result := process.NewResult(xmlResult, &xmlResult.XmlConstraintsConclusionWithProofOfExistenceContent.XmlConstraintsConclusionContent,
		&xmlResult.XmlConstraintsConclusionAttrs)
	c := &ValidationProcessForSignaturesWithLongTermValidationData{
		ChainBase:        process.NewChainBase(i18nProvider, result),
		diagnosticData:   diagnosticData,
		currentSignature: currentSignature,
		bbbs:             bbbs,
		policy:           validationPolicy,
		currentDate:      currentDate,
	}
	if signatureAnalysis.ValidationProcessBasicSignature != nil {
		c.basicSignatureValidation = &signatureAnalysis.ValidationProcessBasicSignature.XmlConstraintsConclusionWithProofOfExistenceContent.XmlConstraintsConclusionContent
	}
	c.xmlTimestamps = signatureAnalysis.Timestamp
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *ValidationProcessForSignaturesWithLongTermValidationData) Title() i18n.MessageTag {
	return i18n.MessageTag_VPFLTVD
}

// InitChain initializes the chain. Port of initChain().
func (c *ValidationProcessForSignaturesWithLongTermValidationData) InitChain() {

	currentContext := enumerations.ContextSignature
	if c.currentSignature.IsCounterSignature() {
		currentContext = enumerations.ContextCounterSignature
	} else if c.currentSignature.IsKeyBindingSignature() {
		currentContext = enumerations.ContextKeyBindingSignature
	}

	/*
	 * 5.5.4 1) The process shall initialize the set of signature time-stamp tokens from the signature time-stamp
	 * attributes present in the signature and shall initialize the best-signature-time to the current time.
	 * NOTE 1: Best-signature-time is an internal variable for the algorithm denoting the earliest time when it can
	 * be proven that a signature has existed.
	 */
	bestSignatureTime := c.getCurrentTime()
	c.Result.Value.ProofOfExistence = bestSignatureTime

	/*
	 * 2) Signature validation: the process shall perform the validation process for Basic Signatures as per
	 * clause 5.3 with all the inputs, including the processing of any signed attributes as specified. If the
	 * Signature contains long-term validation data, this data shall be passed to the validation process for Basic
	 * Signatures.
	 *
	 * If this validation returns PASSED, INDETERMINATE/CRYPTO_CONSTRAINTS_FAILURE_NO_POE,
	 * INDETERMINATE/REVOKED_NO_POE, INDETERMINATE/REVOKED_CA_NO_POE,
	 * INDETERMINATE/TRY_LATER, INDETERMINATE/OUT_OF_BOUNDS_NO_POE or
	 * INDETERMINATE/OUT_OF_BOUNDS_NOT_REVOKED, the SVA shall go to the next step. Otherwise, the
	 * process shall return the status and information returned by the validation process for Basic Signatures.
	 */
	item := c.isAcceptableBasicSignatureValidation()
	c.FirstItem = item

	bsConclusion := c.basicSignatureValidation.Conclusion
	if !process.IsAllowedBasicSignatureValidation(bsConclusion) {
		return
	}

	/* Revocation BBBs analysis */
	c.certificateRevocationMap = make(map[string]*diagnostic.CertificateRevocationWrapper)

	for _, certificateWrapper := range c.currentSignature.CertificateChain() {
		subContext := enumerations.SubContextCACertificate
		if c.currentSignature.SigningCertificate() != nil && c.currentSignature.SigningCertificate().Id() == certificateWrapper.Id() {
			subContext = enumerations.SubContextSigningCert
		}
		if c.isTrustAnchor(certificateWrapper, bestSignatureTime.Time.Time(), currentContext, subContext) {
			break
		}

		revocationDataRequired := c.revocationDataRequired(certificateWrapper, currentContext, subContext)

		if !revocationDataRequired.Process() {
			item = item.SetNextItem(revocationDataRequired)
			continue
		}

		item = item.SetNextItem(c.revocationDataPresent(certificateWrapper, currentContext, c.getSubContext(certificateWrapper)))

		if utils.IsCollectionEmpty(certificateWrapper.CertificateRevocationData()) {
			continue
		}

		certificateRevocationSelector := vpfltvd.NewLongTermValidationCertificateRevocationSelector(
			c.I18nProvider, certificateWrapper, c.currentDate, c.diagnosticData, c.bbbs, c.currentSignature.Id(), c.policy)
		xmlCRS := certificateRevocationSelector.Execute()
		c.Result.Value.CRS = append(c.Result.Value.CRS, xmlCRS)

		item = item.SetNextItem(c.checkCertificateRevocationSelectorResult(xmlCRS, currentContext, c.getSubContext(certificateWrapper)))

		latestCertificateRevocation := certificateRevocationSelector.LatestAcceptableCertificateRevocation()

		if latestCertificateRevocation != nil {
			if _, ok := c.certificateRevocationMap[certificateWrapper.Id()]; !ok {
				c.certificateRevocationOrder = append(c.certificateRevocationOrder, certificateWrapper)
			}
			c.certificateRevocationMap[certificateWrapper.Id()] = latestCertificateRevocation
		}
	}

	var filteredTimestamps []*diagnostic.TimestampWrapper

	/*
	 * 3) Signature time-stamp validation:
	 */
	signatureTimestamps := c.currentSignature.TLevelTimestamps()

	if utils.IsCollectionNotEmpty(signatureTimestamps) {

		/*
		 * a) For each time-stamp token in the set of signature time-stamp tokens, the process shall check that the
		 * message imprint has been generated according to the corresponding signature format specification
		 * verification. If the verification fails, the process shall remove the token from the set.
		 */
		for _, timestampWrapper := range signatureTimestamps {

			item = item.SetNextItem(c.timestampMessageImprint(timestampWrapper))

			if timestampWrapper.IsMessageImprintDataFound() && timestampWrapper.IsMessageImprintDataIntact() {

				/*
				 * b) Time-stamp token validation: For each time-stamp token remaining in the set of signature
				 * time-stamp tokens, the process shall perform the time-stamp validation process as per clause 5.4:
				 *
				 * If PASSED is returned and if the returned generation time is before best-signature-time,
				 * the process shall set best-signature-time to this date and shall try the next token.
				 *
				 * In all other cases:
				 *
				 * - If no specific constraints mandating the validity of the attribute are specified in the validation
				 *   constraints, the process shall remove the time-stamp token from the set of signature time-stamp
				 *   tokens and shall try the next token.
				 *
				 * - Otherwise, the process shall return the indication/sub-indication and associated explanations
				 *   returned from the Time-stamp token validation process.
				 */
				timestampValidationProcess := c.getTimestampValidationProcess(timestampWrapper.Id())
				if timestampValidationProcess != nil {
					item = item.SetNextItem(c.timestampBasicSignatureValidation(timestampWrapper, timestampValidationProcess))
				}

				if timestampValidationProcess != nil && c.IsValid(&timestampValidationProcess.XmlConstraintsConclusionContent) {

					filteredTimestamps = append(filteredTimestamps, timestampWrapper)

					productionTime := timestampWrapper.ProductionTime()
					if productionTime != nil && productionTime.Before(bestSignatureTime.Time.Time()) {
						bestSignatureTime = c.getProofOfExistence(timestampWrapper)
					}

				}

			}

		}

	}

	// If no LTA material, perform *-level timestamp validation
	if !process.IsLongTermAvailabilityAndIntegrityMaterialPresent(c.currentSignature) {

		item = item.SetNextItem(c.tLevelTimeStamp(currentContext))

		item = item.SetNextItem(c.ltaLevelTimeStamp(currentContext))

	}

	/*
	 * 4) Comparing times:
	 * a) If step 2) returned the indication INDETERMINATE with the sub-indication REVOKED_NO_POE or
	 * REVOKED_CA_NO_POE:
	 * a. If the returned revocation time is posterior to best-signature-time, then:
	 *     i.   If best-signature-time is before the issuance date of the signing certificate, the process shall
	 *          return the indication FAILED with the sub-indication NOT_YET_VALID.
	 *     ii.  If best-signature-time is before the expiration date of the signing certificate, the process shall
	 *          perform step 4)e).
	 *     iii. Otherwise, the process shall return the indication
	 *          INDETERMINATE/OUT_OF_BOUNDS_NOT_REVOKED.
	 * b. Otherwise, the process shall return the indication INDETERMINATE with the sub-indication
	 *    REVOKED_NO_POE or REVOKED_CA_NO_POE, respectively.
	 */
	bsSubIndication := xmlSubIndication(bsConclusion)
	if enumerations.IndicationIndeterminate == bsConclusion.Indication.Indication() &&
		(enumerations.SubIndicationRevokedNoPOE == bsSubIndication || enumerations.SubIndicationRevokedCANoPOE == bsSubIndication) {

		bestSignatureTimeTime := bestSignatureTime.Time.Time()
		item = c.revocationDateAfterBestSignatureTimeValidation(item, &bestSignatureTimeTime, bsSubIndication)

		if c.currentSignature.SigningCertificate() != nil {

			item = item.SetNextItem(c.bestSignatureTimeNotBeforeCertificateIssuance(&bestSignatureTimeTime))

			item = item.SetNextItem(c.bestSignatureTimeBeforeCertificateExpiration(&bestSignatureTimeTime))

		}

	}

	/*
	 * b) If step 2) returned the indication PASSED or the indication INDETERMINATE with the sub-indication
	 * OUT_OF_BOUNDS_NO_POE: If best-signature-time is before the issuance date of the signing
	 * certificate, the process shall return the indication FAILED with the sub-indication NOT_YET_VALID.
	 * Otherwise:
	 * a. If the returned indication was PASSED, the process shall continue with step 4)e);
	 * b. Else, the process shall return the indication and sub-indication which was returned by step 2).
	 */
	if enumerations.IndicationPassed == bsConclusion.Indication.Indication() ||
		(enumerations.IndicationIndeterminate == bsConclusion.Indication.Indication() && enumerations.SubIndicationOutOfBoundsNoPOE == bsSubIndication) {

		// verify signing certificate presence for the check
		if c.currentSignature.SigningCertificate() != nil {

			bestSignatureTimeTime := bestSignatureTime.Time.Time()
			item = item.SetNextItem(c.bestSignatureTimeNotBeforeCertificateIssuance(&bestSignatureTimeTime))

			if enumerations.IndicationPassed != bsConclusion.Indication.Indication() {

				item = item.SetNextItem(c.certificateKnownToBeNotRevokedFail(bsConclusion, &bestSignatureTimeTime)) //nolint:staticcheck // mirrors upstream ValidationProcessForSignaturesWithLongTermValidationData#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

				return
			}

		}
	}

	/*
	 * c) If step 2) returned INDETERMINATE with the sub-indication CRYPTO_CONSTRAINTS_FAILURE_NO_POE and the
	 * material concerned by this failure is the signature value or a signed attribute: If the algorithm(s)
	 * concerned were still considered reliable at best-signature-time, the process shall continue with step 4-e).
	 * Otherwise, the process shall return the indication INDETERMINATE with the sub-indication
	 * CRYPTO_CONSTRAINTS_FAILURE_NO_POE.
	 */
	if c.isCryptoConstraintFailureNoPoe(bsConclusion) {

		bestSignatureTimeTime := bestSignatureTime.Time.Time()
		item = item.SetNextItem(c.signatureValueAndSignedAttributesAlgorithmsAcceptable(bestSignatureTimeTime, currentContext))

	}

	/*
	 * d) If step 2) returned the indication INDETERMINATE with the sub-indication OUT_OF_BOUNDS_NOT_REVOKED: If
	 * best-signature-time is before the issuance date of the signing certificate, the process shall return the indication
	 * FAILED with the sub-indication NOT_YET_VALID.
	 * Otherwise:
	 * a. If best-signature-time is before the expiration date of the signing certificate, the process shall
	 * perform step 4)e).
	 * b. Else, the process shall return the indication INDETERMINATE/OUT_OF_BOUNDS_NOT_REVOKED.
	 */
	if enumerations.IndicationIndeterminate == bsConclusion.Indication.Indication() && enumerations.SubIndicationOutOfBoundsNotRevoked == bsSubIndication {

		bestSignatureTimeTime := bestSignatureTime.Time.Time()
		item = item.SetNextItem(c.bestSignatureTimeNotBeforeCertificateIssuance(&bestSignatureTimeTime))

		item = item.SetNextItem(c.bestSignatureTimeBeforeCertificateExpiration(&bestSignatureTimeTime))

		if c.revocationDataRequired(c.currentSignature.SigningCertificate(), currentContext, enumerations.SubContextSigningCert).Process() {
			item = item.SetNextItem(c.certificateKnownToBeNotRevokedWarn(bsConclusion, &bestSignatureTimeTime, currentContext))
		}

	}

	if utils.IsCollectionNotEmpty(filteredTimestamps) {
		/*
		 * e) For each time-stamp token remaining in the set of signature time-stamp tokens, the process shall check
		 * the coherence in the values of the times indicated in the time-stamp tokens. They shall be posterior to
		 * the times indicated in any time-stamp token computed on the signed data (content-time-stamp). The process shall apply the
		 * rules specified in IETF RFC 3161 [3], clause 2.4.2 regarding the order of time-stamp tokens generated by the
		 * same or different TSAs given the accuracy and ordering fields' values of the TSTInfo field,
		 * unless stated differently by the signature validation constraints. If all the checks end successfully,
		 * the process shall go to the next step. Otherwise the process shall return the indication INDETERMINATE with the
		 * sub-indication TIMESTAMP_ORDER_FAILURE.
		 */
		item = item.SetNextItem(c.timestampCoherenceOrder(c.currentSignature.TimestampList()))

		/*
		 * 5) Handling Time-stamp delay: If the signature contains a signature time-stamp token and the validation
		 * constraints specify a time-stamp delay:
		 */

		if len(signatureTimestamps) != 0 && c.policy.TimestampDelayConstraint() != nil {
			/*
			 * a) If no signing-time property/attribute is present, the process shall return the indication
			 * INDETERMINATE with the sub-indication SIG_CONSTRAINTS_FAILURE.
			 */
			item = item.SetNextItem(c.signingTimeAttributePresent(currentContext))
			/*
			 * b) If a signing-time property/attribute is present, the process shall check that the claimed time in the
			 * attribute plus the time-stamp delay is after the best-signature-time. If the check is successful, the
			 * process shall go to the next step. Otherwise, the process shall return the indication INDETERMINATE with the
			 * sub-indication SIG_CONSTRAINTS_FAILURE.
			 */
			bestSignatureTimeTime := bestSignatureTime.Time.Time()
			item = item.SetNextItem(c.timestampDelay(&bestSignatureTimeTime))
		}
	}

	/*
	 * 6) If step 2) returned the indication INDETERMINATE with the sub indication
	 * TRY_LATER because the revocation information was not fresh enough: the building block
	 * shall run the Revocation Freshness Checker (clause 5.2.5) with the revocation status
	 * information returned in step 2), the certificate for which the revocation status
	 * is being checked and best signature time. If the checker returns PASSED, the building block
	 * shall go to the next step. Otherwise, the building block shall return the indication INDETERMINATE,
	 * the sub indication TRY_LATER and, if returned from the Revocation Freshness Checker,
	 * the suggestion for when to try the validation again.
	 *
	 * 7) If step 2) returned the indication INDETERMINATE with the sub indication TRY_LATER
	 * because the certificate has been found to be suspended:
	 *    a. If best-signature-time is before the time of suspension of the certificate:
	 *       the process shall go to the step 8).
	 *    b. Otherwise, the building block shall return the indication INDETERMINATE,
	 *       the sub indication TRY_LATER and a suggestion on when to try the validation gain,
	 *       if returned by the validation process in step 2).
	 */
	if enumerations.IndicationIndeterminate == bsConclusion.Indication.Indication() && enumerations.SubIndicationTryLater == bsSubIndication {
		bestSignatureTimeTime := bestSignatureTime.Time.Time()
		item = c.revocationIsFresh(item, &bestSignatureTimeTime, currentContext)
	}

	/*
	 * 8) The SVA shall perform the Signature Acceptance Validation process as per clause 5.2.8 with the following
	 * inputs:
	 * a) The Signed Data Object(s).
	 * b) best-signature-time as the validation time parameter.
	 * c) The Cryptographic Constraints.
	 */
	bestSignatureTimeTime := bestSignatureTime.Time.Time()
	item = item.SetNextItem(c.signedDataObjectAlgorithmsAcceptable(bestSignatureTimeTime, currentContext))

	/*
	 * 9) If the signature acceptance validation process returns PASSED, the SVA shall go to the next step.
	 * Otherwise, the SVA shall return the indication and sub-indication returned by
	 * the Signature Acceptance Validation Process.
	 */
	item = item.SetNextItem(c.signatureIsAcceptable())

	/*
	 * 10) The SVA shall apply the cryptographic constraints to all the certificates and revocation status information used
	 * in the validation process against the current time. If any of those certificates or revocation status information
	 * do not match these constraints, the SVA shall return the indication INDETERMINATE with the sub-indication
	 * CRYPTO_CONSTRAINTS_FAILURE_NO_POE together with the list of algorithms and key sizes, if
	 * applicable, that are concerned and the time for each of the algorithms up to which the respective algorithm has
	 * been considered secure by the cryptographic constraints. Otherwise, the SVA shall go to the next step.
	 */

	// check validity of Cryptographic Constraints for the Signing Certificate and CA Certificates
	item = c.certificateChainReliableAtTime(item, c.currentSignature, c.currentDate, currentContext)

	// check validity of revocation data
	item = c.revocationDataReliableAtTime(item, c.currentDate) //nolint:staticcheck // mirrors upstream ValidationProcessForSignaturesWithLongTermValidationData#initChain: Java's `item = revocationDataReliableAtTime(item, currentDate);` is the same dead store - the helper links and returns the new tail, which nothing reads.

	/*
	 * 11) Data extraction: the process shall return the success indication PASSED,
	 * the certificate chain obtained in step 2 and best-signature-time.
	 * In addition, the process should return additional information extracted from the signature and/or
	 * used by the intermediate steps.
	 * In particular, the process should return intermediate results such as the validation results
	 * of any signature time-stamp token.
	 */
	c.Result.Value.ProofOfExistence = bestSignatureTime
}

// xmlSubIndication reads conclusion.SubIndication; nil is Java's null.
func xmlSubIndication(conclusion *jaxb.XmlConclusion) enumerations.SubIndication {
	if conclusion.SubIndication == nil {
		return ""
	}
	return conclusion.SubIndication.SubIndication()
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) isAcceptableBasicSignatureValidation() process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	return vpfltvd.NewAcceptableBasicSignatureValidationCheck(c.I18nProvider, c.Result, c.basicSignatureValidation, c.FailLevelRule())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) revocationDataRequired(certificate *diagnostic.CertificateWrapper,
	context enumerations.Context, subContext enumerations.SubContext) *xcv.RevocationDataRequiredCheck[*jaxb.XmlValidationProcessLongTermData] {
	constraint := c.policy.RevocationDataSkipConstraint(context, subContext)
	sunsetDateConstraint := c.policy.CertificateSunsetDateConstraint(context, subContext)
	return xcv.NewRevocationDataRequiredCheck(c.I18nProvider, c.Result, certificate, c.currentDate, sunsetDateConstraint, constraint)
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) revocationDataPresent(certificate *diagnostic.CertificateWrapper,
	context enumerations.Context, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	constraint := c.policy.RevocationDataAvailableConstraint(context, subContext)
	return xcv.NewRevocationDataAvailableCheckWithId(c.I18nProvider, c.Result, certificate, constraint, certificate.Id())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) checkCertificateRevocationSelectorResult(
	crsResult *jaxb.XmlCRS, context enumerations.Context, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	constraint := c.policy.AcceptableRevocationDataFoundConstraint(context, subContext)
	return xcv.NewCertificateRevocationSelectorResultCheck(c.I18nProvider, c.Result, crsResult, constraint)
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) timestampMessageImprint(
	timestampWrapper *diagnostic.TimestampWrapper) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	return vpfltvd.NewTimestampMessageImprintWithIdCheck(c.I18nProvider, c.Result, timestampWrapper, c.WarnLevelRule())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) timestampBasicSignatureValidation(
	timestampWrapper *diagnostic.TimestampWrapper,
	timestampValidationResult *jaxb.XmlValidationProcessBasicTimestamp) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	return vpfbs.NewBasicTimestampValidationWithIdCheck(c.I18nProvider, c.Result, timestampWrapper,
		timestampValidationResult, c.getTimestampBasicValidationConstraintLevel())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) tLevelTimeStamp(
	context enumerations.Context) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	constraint := c.policy.TLevelTimeStampConstraint(context)
	return sav.NewTLevelTimeStampCheck(c.I18nProvider, c.Result, c.currentSignature, c.bbbs, c.xmlTimestamps, constraint)
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) ltaLevelTimeStamp(
	context enumerations.Context) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	constraint := c.policy.LTALevelTimeStampConstraint(context)
	return sav.NewLTALevelTimeStampCheck(c.I18nProvider, c.Result, c.currentSignature, c.bbbs, c.xmlTimestamps, constraint)
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) getTimestampBasicValidationConstraintLevel() policy.LevelRule {
	constraint := c.policy.TimestampValidConstraint()
	// continue if LTA is present
	if constraint == nil || process.IsLongTermAvailabilityAndIntegrityMaterialPresent(c.currentSignature) {
		constraint = c.WarnLevelRule()
	}
	return constraint
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) revocationIsFresh(
	item process.ChainItem[*jaxb.XmlValidationProcessLongTermData], bestSignatureTime *time.Time,
	currentContext enumerations.Context) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	for _, certificate := range c.currentSignature.CertificateChain() {
		revocationData, ok := c.certificateRevocationMap[certificate.Id()]
		if !ok {
			continue
		}

		subContext := c.getSubContext(certificate)

		if enumerations.RevocationReasonCertificateHold == revocationData.Reason() {
			item = item.SetNextItem(c.checkCertificateSuspensionNotBeforeBestSignatureTime(revocationData,
				bestSignatureTime, currentContext, subContext))

		} else {
			rfc := xcv.NewRevocationFreshnessChecker(c.I18nProvider, &revocationData.RevocationWrapper,
				*bestSignatureTime, currentContext, subContext, c.policy)
			xmlRFC := rfc.Execute()
			c.Result.Value.RFC = append(c.Result.Value.RFC, xmlRFC)

			item = item.SetNextItem(c.checkRevocationFreshnessCheckerResult(xmlRFC))

			if !c.IsValid(&xmlRFC.XmlConstraintsConclusionContent) {
				break
			}
		}

	}
	return item
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) checkRevocationFreshnessCheckerResult(
	rfcResult *jaxb.XmlRFC) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	return newRevocationFreshnessCheckerTryLaterResultCheck(c.I18nProvider, c.Result, rfcResult, c.FailLevelRule())
}

// revocationFreshnessCheckerTryLaterResultCheck is the Go form of the
// anonymous RevocationFreshnessCheckerResultCheck subclass returned by
// checkRevocationFreshnessCheckerResult(XmlRFC).
type revocationFreshnessCheckerTryLaterResultCheck struct {
	*xcv.RevocationFreshnessCheckerResultCheck[*jaxb.XmlValidationProcessLongTermData]
}

func newRevocationFreshnessCheckerTryLaterResultCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessLongTermData], rfcResult *jaxb.XmlRFC,
	constraint policy.LevelRule) *revocationFreshnessCheckerTryLaterResultCheck {
	c := &revocationFreshnessCheckerTryLaterResultCheck{
		RevocationFreshnessCheckerResultCheck: xcv.NewRevocationFreshnessCheckerResultCheck(i18nProvider, result, rfcResult, constraint),
	}
	c.InitChainItem(c)
	return c
}

func (c *revocationFreshnessCheckerTryLaterResultCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

func (c *revocationFreshnessCheckerTryLaterResultCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationTryLater
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) checkCertificateSuspensionNotBeforeBestSignatureTime(
	certificateRevocationWrapper *diagnostic.CertificateRevocationWrapper, bestSignatureTime *time.Time,
	context enumerations.Context, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	constraint := c.policy.CertificateNotOnHoldConstraint(context, subContext)
	return vpfltvd.NewBestSignatureTimeBeforeSuspensionTimeCheck(c.I18nProvider, c.Result, certificateRevocationWrapper, bestSignatureTime, constraint)
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) revocationDateAfterBestSignatureTimeValidation(
	item process.ChainItem[*jaxb.XmlValidationProcessLongTermData], bestSignatureTime *time.Time,
	subIndication enumerations.SubIndication) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {

	constraint := c.policy.RevocationTimeAgainstBestSignatureTimeConstraint()

	for _, certificate := range c.certificateRevocationOrder {
		revocationData := c.certificateRevocationMap[certificate.Id()]
		subContext := c.getSubContext(certificate)

		// separate cases to check based on the returned subIndication
		if (enumerations.SubContextSigningCert == subContext && enumerations.SubIndicationRevokedNoPOE == subIndication) ||
			(enumerations.SubContextCACertificate == subContext && enumerations.SubIndicationRevokedCANoPOE == subIndication) {

			item = item.SetNextItem(vpfltvd.NewRevocationDateAfterBestSignatureTimeCheck(c.I18nProvider, c.Result, revocationData,
				bestSignatureTime, constraint, subContext))

		}
	}

	return item
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) bestSignatureTimeNotBeforeCertificateIssuance(
	bestSignatureTime *time.Time) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	signingCertificate := c.currentSignature.SigningCertificate()
	return vpfltvd.NewBestSignatureTimeNotBeforeCertificateIssuanceCheck(c.I18nProvider, c.Result, bestSignatureTime, signingCertificate, c.FailLevelRule())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) certificateKnownToBeNotRevokedFail(
	bsConclusion *jaxb.XmlConclusion, bestSignatureTime *time.Time) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	signingCertificate := c.currentSignature.SigningCertificate()
	var revocationWrapper *diagnostic.CertificateRevocationWrapper
	if signingCertificate != nil {
		revocationWrapper = c.certificateRevocationMap[signingCertificate.Id()]
	}
	isRevocationIssuerTrusted := c.isRevocationIssuerTrusted(revocationWrapper, bestSignatureTime)
	return vpfltvd.NewCertificateKnownToBeNotRevokedEnforceFailCheck(c.I18nProvider, c.Result, signingCertificate, revocationWrapper,
		isRevocationIssuerTrusted, &c.currentDate, bsConclusion, c.FailLevelRule())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) isRevocationIssuerTrusted(
	revocationWrapper *diagnostic.CertificateRevocationWrapper, bestSignatureTime *time.Time) bool {
	sunsetDateConstraint := c.policy.CertificateSunsetDateConstraint(enumerations.ContextRevocation, enumerations.SubContextSigningCert)
	return revocationWrapper != nil && revocationWrapper.SigningCertificate() != nil && bestSignatureTime != nil &&
		process.IsTrustAnchor(revocationWrapper.SigningCertificate(), *bestSignatureTime, sunsetDateConstraint)
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) certificateKnownToBeNotRevokedWarn(
	bsConclusion *jaxb.XmlConclusion, bestSignatureTime *time.Time,
	context enumerations.Context) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	signingCertificate := c.currentSignature.SigningCertificate()
	var revocationWrapper *diagnostic.CertificateRevocationWrapper
	if signingCertificate != nil {
		revocationWrapper = c.certificateRevocationMap[signingCertificate.Id()]
	}
	isRevocationIssuerTrusted := c.isRevocationIssuerTrusted(revocationWrapper, bestSignatureTime)
	constraint, err := process.GetConstraintOrMaxLevel(c.policy.RevocationIssuerNotExpiredConstraint(context, enumerations.SubContextSigningCert), enumerations.LevelWarn)
	if err != nil {
		panic(err)
	}
	return vpfltvd.NewCertificateKnownToBeNotRevokedCheck(c.I18nProvider, c.Result, signingCertificate, revocationWrapper,
		isRevocationIssuerTrusted, &c.currentDate, bsConclusion, constraint)
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) bestSignatureTimeBeforeCertificateExpiration(
	bestSignatureTime *time.Time) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	signingCertificate := c.currentSignature.SigningCertificate()
	return vpfltvd.NewBestSignatureTimeBeforeCertificateExpirationCheck(c.I18nProvider, c.Result, bestSignatureTime, signingCertificate,
		c.policy.BestSignatureTimeBeforeExpirationDateOfSigningCertificateConstraint())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) timestampCoherenceOrder(
	timestamps []*diagnostic.TimestampWrapper) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	return vpfltvd.NewTimestampCoherenceOrderCheck(c.I18nProvider, c.Result, timestamps, c.policy.TimestampCoherenceConstraint())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) signingTimeAttributePresent(
	context enumerations.Context) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	return vpfltvd.NewSigningTimeAttributePresentCheck(c.I18nProvider, c.Result, c.currentSignature, c.policy.SigningTimeConstraint(context))
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) timestampDelay(
	bestSignatureTime *time.Time) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	return vpfltvd.NewTimestampDelayCheck(c.I18nProvider, c.Result, c.currentSignature, bestSignatureTime, c.policy.TimestampDelayConstraint())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) tokenUsedAlgorithmsAreSecureAtTimeWithId(
	currentToken diagnostic.TokenProxy, validationDate time.Time,
	position i18n.MessageTag) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	revocationBBB := c.bbbs[currentToken.Id()]
	xmlAOV := revocationBBB.AOV
	return aov.NewAlgorithmObsolescenceValidationCheckWithId(c.I18nProvider, c.Result, xmlAOV, validationDate, position, currentToken.Id())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) signatureValueAndSignedAttributesAlgorithmsAcceptable(
	bestSignatureTime time.Time, context enumerations.Context) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	algorithmObsolescenceValidation := aov.NewSignatureValueAndSignedAttributesAlgorithmObsolescenceValidation(
		c.I18nProvider, c.currentSignature, context, bestSignatureTime, c.policy)
	aovResult := algorithmObsolescenceValidation.Execute()

	return aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, aovResult, bestSignatureTime,
		i18n.MessageTag_ACCM_POS_SIG_VAL_AND_PRT, c.currentSignature.Id())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) signedDataObjectAlgorithmsAcceptable(
	bestSignatureTime time.Time, context enumerations.Context) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	// NOTE: we execute AlgorithmObsolescenceValidation check, as the only time sensitive part of SAV
	algorithmObsolescenceValidation := aov.NewSignatureSignedDataAlgorithmObsolescenceValidation(
		c.I18nProvider, c.currentSignature, context, bestSignatureTime, c.policy)
	aovResult := algorithmObsolescenceValidation.Execute()

	return aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, aovResult, bestSignatureTime,
		i18n.MessageTag_ACCM_POS_SIGND_OBJ, c.currentSignature.Id())
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) signatureIsAcceptable() process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	signatureBBB := c.bbbs[c.currentSignature.Id()]
	return vpfbs.NewSignatureAcceptanceValidationNoCryptoResultCheck(c.I18nProvider, c.Result, signatureBBB.SAV, c.currentSignature, c.FailLevelRule())
}

// certificateChainReliableAtTime sets up cryptographic check for certificates
// used in the certificate chain of the signature. Port of
// certificateChainReliableAtTime(ChainItem, TokenProxy, Date, Context).
func (c *ValidationProcessForSignaturesWithLongTermValidationData) certificateChainReliableAtTime(
	item process.ChainItem[*jaxb.XmlValidationProcessLongTermData], token diagnostic.TokenProxy, validationTime time.Time,
	context enumerations.Context) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	if token.SigningCertificate() == nil || c.isTrustAnchor(token.SigningCertificate(), validationTime, context, enumerations.SubContextSigningCert) ||
		utils.IsCollectionEmpty(token.CertificateChain()) {
		return item
	}

	algorithmObsolescenceValidation := aov.NewTokenCertificateChainAlgorithmObsolescenceValidation[*diagnostic.SignatureWrapper](
		c.I18nProvider, c.currentSignature, context, validationTime, c.policy)
	aovResult := algorithmObsolescenceValidation.Execute()

	position, err := process.GetCertificateChainCryptoPosition(context)
	if err != nil {
		panic(err)
	}

	item = item.SetNextItem(aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, aovResult, validationTime, position, c.currentSignature.Id()))

	return item
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) revocationDataReliableAtTime(
	item process.ChainItem[*jaxb.XmlValidationProcessLongTermData],
	validationTime time.Time) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	checkedTokenIds := make(map[string]struct{})
	for _, certificate := range c.certificateRevocationOrder {
		revocationData := c.certificateRevocationMap[certificate.Id()]
		item = c.checkRevocationAgainstBestSignatureTime(item, &revocationData.RevocationWrapper, validationTime, checkedTokenIds)
	}
	return item
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) checkRevocationAgainstBestSignatureTime(
	item process.ChainItem[*jaxb.XmlValidationProcessLongTermData], revocationData *diagnostic.RevocationWrapper,
	validationTime time.Time, checkedTokenIds map[string]struct{}) process.ChainItem[*jaxb.XmlValidationProcessLongTermData] {
	revocationBBB := c.bbbs[revocationData.Id()]
	if _, ok := checkedTokenIds[revocationData.Id()]; !ok && revocationBBB != nil {

		position, err := process.GetCryptoPosition(enumerations.ContextRevocation)
		if err != nil {
			panic(err)
		}
		item = item.SetNextItem(c.tokenUsedAlgorithmsAreSecureAtTimeWithId(revocationData, validationTime, position))

		checkedTokenIds[revocationData.Id()] = struct{}{}
	}
	return item
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) isTrustAnchor(certificateWrapper *diagnostic.CertificateWrapper,
	currentTime time.Time, context enumerations.Context, subContext enumerations.SubContext) bool {
	constraint := c.policy.CertificateSunsetDateConstraint(context, subContext)
	return process.IsTrustAnchor(certificateWrapper, currentTime, constraint)
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) getSubContext(certificateWrapper *diagnostic.CertificateWrapper) enumerations.SubContext {
	if c.currentSignature.SigningCertificate() != nil && c.currentSignature.SigningCertificate().Id() == certificateWrapper.Id() {
		return enumerations.SubContextSigningCert
	}
	return enumerations.SubContextCACertificate
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) getCurrentTime() *jaxb.XmlProofOfExistence {
	poe := &jaxb.XmlProofOfExistence{}
	poe.Time = jaxb.XSDateTime(c.currentDate)
	return poe
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) getProofOfExistence(timestampWrapper *diagnostic.TimestampWrapper) *jaxb.XmlProofOfExistence {
	xpoe := &jaxb.XmlProofOfExistence{}
	if productionTime := timestampWrapper.ProductionTime(); productionTime != nil {
		xpoe.Time = jaxb.XSDateTime(*productionTime)
	}
	id := timestampWrapper.Id()
	xpoe.TimestampId = &id
	return xpoe
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) getTimestampValidationProcess(timestampId string) *jaxb.XmlValidationProcessBasicTimestamp {
	for _, xmlTimestamp := range c.xmlTimestamps {
		if xmlTimestamp.Id != nil && timestampId == *xmlTimestamp.Id {
			return xmlTimestamp.ValidationProcessBasicTimestamp
		}
	}
	return nil
}

func (c *ValidationProcessForSignaturesWithLongTermValidationData) isCryptoConstraintFailureNoPoe(conclusion *jaxb.XmlConclusion) bool {
	return enumerations.IndicationIndeterminate == conclusion.Indication.Indication() &&
		enumerations.SubIndicationCryptoConstraintsFailureNoPOE == xmlSubIndication(conclusion)
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of the overridden
// collectMessages(XmlConclusion, XmlConstraint).
func (c *ValidationProcessForSignaturesWithLongTermValidationData) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	if constraint.BlockType != nil && jaxb.XmlBlockTypeTSTBBB == *constraint.BlockType && c.policy.TimestampValidConstraint() == nil {
		// skip validation messages for TSTs
	} else {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}

// CollectAdditionalMessages fills additional messages into the conclusion.
// Port of the overridden collectAdditionalMessages(XmlConclusion).
func (c *ValidationProcessForSignaturesWithLongTermValidationData) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	if !process.IsAllowedBasicSignatureValidation(c.basicSignatureValidation.Conclusion) {
		conclusion.Warnings = append(conclusion.Warnings, c.basicSignatureValidation.Conclusion.Warnings...)
		conclusion.Infos = append(conclusion.Infos, c.basicSignatureValidation.Conclusion.Infos...)
	}
}
