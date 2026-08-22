// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/PastSignatureValidation.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note, and
// past_signature_validation_certificate_revocation_selector.go for the assumed
// Go shape of the cross-chunk vpfltvd classes - here
// BestSignatureTimeNotBeforeCertificateIssuanceCheck, a plain generic ChainItem
// whose Go constructor is
//
//	vpfltvd.NewBestSignatureTimeNotBeforeCertificateIssuanceCheck[T](
//	    *i18n.I18nProvider, *process.Result[T], time.Time, *diagnostic.CertificateWrapper, policy.LevelRule)
package vpfswatsp

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
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/xcv"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfltvd"
)

// PastSignatureValidation performs the "5.6.2.4 Past signature validation
// building block".
type PastSignatureValidation struct {
	*process.ChainBase[*jaxb.XmlPSV]

	// token is the token to check.
	token diagnostic.TokenProxy

	// bbbs is the map of all BBBs.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// currentConclusion is the current conclusion.
	currentConclusion *jaxb.XmlConclusion

	// poe is the POE container.
	poe *POEExtraction

	// currentTime is the validation time.
	currentTime time.Time

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// context is the validation context.
	context enumerations.Context
}

// NewPastSignatureValidation is the default constructor. Port of
// PastSignatureValidation(I18nProvider, TokenProxy, Map, XmlConclusion, POEExtraction, Date, ValidationPolicy, Context).
func NewPastSignatureValidation(i18nProvider *i18n.I18nProvider, token diagnostic.TokenProxy,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, currentConclusion *jaxb.XmlConclusion, poe *POEExtraction,
	currentTime time.Time, validationPolicy policy.ValidationPolicy,
	context enumerations.Context) *PastSignatureValidation {
	xmlPSV := &jaxb.XmlPSV{}
	c := &PastSignatureValidation{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlPSV,
			&xmlPSV.XmlConstraintsConclusionContent, &xmlPSV.XmlConstraintsConclusionAttrs)),
		token:             token,
		bbbs:              bbbs,
		currentConclusion: currentConclusion,
		poe:               poe,
		currentTime:       currentTime,
		policy:            validationPolicy,
		context:           context,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *PastSignatureValidation) Title() i18n.MessageTag {
	return i18n.MessageTag_PAST_SIGNATURE_VALIDATION
}

// InitChain initializes the chain. Port of initChain().
func (c *PastSignatureValidation) InitChain() {

	tokenBBB := c.bbbs[c.token.Id()]

	var item process.ChainItem[*jaxb.XmlPSV]

	/*
	 * 1) The building block shall verify that there is at least one revocation data instance
	 * that is known to contain revocation status information about the signing certificate
	 * for which the set of POEs contains a POE for the signing certificate issuer's certificate
	 * after the issuance date and before the expiration date of the signing certificate issuer's certificate:
	 *
	 * a. If there is such a revocation data, the building block shall remove from the Certificate
	 *    Validation Data all revocation data known to contain revocation status information about
	 *    the signing certificate for which there is no such POE and set sig_cert_revocation_poe-status to PASSED.
	 *
	 * b. Otherwise the building block shall set sig_cert_revocation_poe-status to INDETERMINATE with
	 *    the sub-indication REVOCATION_OUT_OF_BOUNDS_NO_POE.
	 */

	signingCertificate := c.token.SigningCertificate()

	sigCertRevocationPoeStatus := &jaxb.XmlConclusion{}
	var signingCertificateRevocations []*diagnostic.CertificateRevocationWrapper

	if c.isRevocationDataRequired(signingCertificate, enumerations.SubContextSigningCert) {
		certificateRevocationSelector := NewPastSignatureValidationCertificateRevocationSelector(
			c.I18nProvider, signingCertificate, c.currentTime, c.bbbs, c.token.Id(), c.poe, c.policy)

		xmlCRS := certificateRevocationSelector.Execute()
		tokenBBB.PSVCRS = xmlCRS

		item = c.checkCertificateRevocationSelectorResult(xmlCRS)
		c.FirstItem = item

		signingCertificateRevocations = certificateRevocationSelector.AcceptableCertificateRevocations()
		if utils.IsCollectionNotEmpty(signingCertificateRevocations) {
			sigCertRevocationPoeStatus.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
		} else {
			sigCertRevocationPoeStatus.Indication = jaxb.IndicationValue(enumerations.IndicationIndeterminate)
			subIndication := jaxb.SubIndicationValue(enumerations.SubIndicationRevocationOutOfBoundsNoPOE)
			sigCertRevocationPoeStatus.SubIndication = &subIndication
			// keep all revocation data if none of the valid instances found
			signingCertificateRevocations = signingCertificate.CertificateRevocationData()
		}

	} else {
		// revocation check is not required
		sigCertRevocationPoeStatus.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
	}

	/*
	 * 2) The building block shall perform the past certificate validation process specified
	 * in clause 5.6.2.1  with the following inputs: the signature, the target certificate,
	 * the X.509 validation parameters, certificate validation data, X.509 validation constraints,
	 * cryptographic constraints and the set of POEs. If it returns PASSED/validation time,
	 * the building block shall go to the next step. Otherwise, the building block shall return
	 * the current time status and sub indication with an explanation of the failure.
	 */
	pcv := NewPastCertificateValidation(c.I18nProvider, c.token, c.bbbs, c.poe, c.currentTime, c.policy, c.context)
	pcvResult := pcv.Execute()
	tokenBBB.PCV = pcvResult

	pastCertificateValidationAcceptableCheck := c.pastCertificateValidationAcceptableCheck(pcvResult)
	if item == nil {
		item = pastCertificateValidationAcceptableCheck
		c.FirstItem = item
	} else {
		item = item.SetNextItem(pastCertificateValidationAcceptableCheck)
	}

	// XmlPCV#getControlTime(): the generated member is a *XSDateTime, whose nil
	// is Java's null.
	var controlTime *time.Time
	if pcvResult.ControlTime != nil {
		t := pcvResult.ControlTime.Time()
		controlTime = &t
		c.Result.Value.ControlTime = jaxb.NewXSDateTime(t)
	} else {
		c.Result.Value.ControlTime = nil
	}

	/*
	 * 3) If there is a POE of the signature value at (or before) the validation time returned in the previous step:
	 */
	poeExistsCheck := c.poeExist(controlTime)

	// NOTE: upstream builds a second, separate instance for the chain; the one
	// above is only asked for its process() result.
	item = item.SetNextItem(c.poeExist(controlTime))

	bestSignatureTime := c.poe.GetLowestPOETime(c.token.Id())

	poeExists := poeExistsCheck.Process()

	currentSubIndication := subIndicationOfConclusion(c.currentConclusion)
	currentIndication := c.currentConclusion.Indication.Indication()

	/*
	 * - If current time indication/sub indication is INDETERMINATE/NO_CERTIFICATE_CHAIN_FOUND_NO_POE:
	 */
	if poeExists && enumerations.IndicationIndeterminate == currentIndication &&
		enumerations.SubIndicationNoCertificateChainFoundNoPOE == currentSubIndication {
		/*
		 * a) If best-signature-time is before the issuance date of the signing certificate (notBefore field), the
		 *    building block shall return the indication FAILED with the sub-indication NOT_YET_VALID.
		 * b) If best-signature-time is after the expiration date of the signing certificate, the building block shall
		 *    return the indication INDETERMINATE with the sub-indication OUT_OF_BOUNDS_NO_POE.
		 * c) Else the building block shall go to step 7).
		 */

		item = item.SetNextItem(c.bestSignatureTimeNotBeforeCertificateIssuance(bestSignatureTime, signingCertificate))

		item = item.SetNextItem(c.bestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpiration(
			bestSignatureTime, signingCertificate, enumerations.SubIndicationOutOfBoundsNoPOE))

	} else if poeExists && enumerations.IndicationIndeterminate == currentIndication &&
		(enumerations.SubIndicationRevokedNoPOE == currentSubIndication ||
			enumerations.SubIndicationRevocationOutOfBoundsNoPOE == currentSubIndication ||
			(enumerations.SubIndicationTryLater == currentSubIndication && c.isCertificateSuspended())) {
		/*
		 * - If current time indication/sub-indication is INDETERMINATE/REVOKED_NO_POE,
		 *   INDETERMINATE/REVOCATION_OUT_OF_BOUNDS_NO_POE or INDETERMINATE/TRY_LATER
		 *   because the certificate has been found to be suspended, then:
		 *
		 * a) If best-signature-time is before the issuance date of the signing certificate,
		 *    the process shall return the indication FAILED with the sub-indication NOT_YET_VALID
		 * b) If best-signature-time is within the validity period of the signing certificate,
		 *    the building block shall go to step 7).
		 * c) Otherwise the building block shall set the current time indication/sub-indication to
		 *    INDETERMINATE/OUT_OF_BOUNDS_NOT_REVOKED and continue the process.
		 */

		item = item.SetNextItem(c.bestSignatureTimeNotBeforeCertificateIssuance(bestSignatureTime, signingCertificate))

		item = item.SetNextItem(c.bestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpiration(
			bestSignatureTime, signingCertificate, enumerations.SubIndicationOutOfBoundsNotRevoked))

	} else if poeExists && enumerations.IndicationIndeterminate == currentIndication &&
		enumerations.SubIndicationRevokedCANoPOE == currentSubIndication {
		/*
		 * - If current time indication/sub-indication is INDETERMINATE/REVOKED_CA_NO_POE then:
		 *
		 * a) If there is a POE for the revocation data containing the revocation status information
		 *    of the signer certificate at (or before) the revocation time of the CA certificate, then:
		 *    i.  If best signature time (lowest time at which there exists a POE for the signature value
		 *        in the set of POEs) is within the validity period of the signing certificate,
		 *        the building block shall go to step 7).
		 *    ii. Otherwise the building block shall set the current time indication/sub-indication to
		 *        OUT_OF_BOUNDS_NOT_REVOKED and continue the process.
		 * b) Otherwise, the building block shall return with the indication INDETERMINATE and the
		 *    sub-indication REVOKED_CA_NO_POE.
		 */

		caCertificate := signingCertificate.SigningCertificate()
		var latestCARevocationData *diagnostic.CertificateRevocationWrapper
		if caCertificate != nil {
			latestCARevocationData = process.GetLatestAcceptableRevocationData(c.token, caCertificate,
				caCertificate.CertificateRevocationData(), c.currentTime, c.bbbs, c.poe)
		}
		if latestCARevocationData != nil {
			item = item.SetNextItem(c.poeExistNotAfterCARevocationTimeCheck(
				signingCertificateRevocations, latestCARevocationData.RevocationDate()))
		}

		// NOTE: executed as a part of the check below (see "continue the process" reference)
		item = item.SetNextItem(c.bestSignatureTimeNotBeforeCertificateIssuance(bestSignatureTime, signingCertificate))

		item = item.SetNextItem(c.bestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpiration(
			bestSignatureTime, signingCertificate, enumerations.SubIndicationOutOfBoundsNotRevoked))

	} else if poeExists && enumerations.IndicationIndeterminate == currentIndication &&
		(enumerations.SubIndicationOutOfBoundsNoPOE == currentSubIndication ||
			enumerations.SubIndicationOutOfBoundsNotRevoked == currentSubIndication) {
		/*
		 * - If current time indication/sub-indication is INDETERMINATE/OUT_OF_BOUNDS_NO_POE or OUT_OF_BOUNDS_NOT_REVOKED:
		 *
		 * a) If best-signature-time (lowest time at which there exists a POE for the signature value in the set of POEs)
		 *    is before the issuance date of the signing certificate (notBefore field), the building block shall
		 *    return the indication FAILED with the sub-indication NOT_YET_VALID.
		 *
		 * b) If best-signature-time (lowest time at which there exists a POE for the signature value in the set of POEs)
		 *    is after the issuance date and before the expiration date of the signing certificate,
		 *    the building block shall go to step 7.
		 */

		item = item.SetNextItem(c.bestSignatureTimeNotBeforeCertificateIssuance(bestSignatureTime, signingCertificate))

		item = item.SetNextItem(c.bestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpiration(
			bestSignatureTime, signingCertificate, currentSubIndication))

	} else if enumerations.IndicationIndeterminate == currentIndication &&
		enumerations.SubIndicationCryptoConstraintsFailureNoPOE == currentSubIndication {
		/*
		 * 4) If current time indication/ sub-indication is INDETERMINATE/CRYPTO_CONSTRAINTS_FAILURE_NO_POE and for
		 * each algorithm (or key size) in the list concerned by the failure, there is a POE for the material that
		 * uses this algorithm (or key size) at a time before the time up to which the algorithm in question was
		 * considered secure, the building block shall go to step 7).
		 */

		item = item.SetNextItem(c.algorithmsObsolescenceValidation())

		item = c.revocationDataAlgorithmsObsolescenceValidation(item, signingCertificateRevocations)

	} else if enumerations.IndicationIndeterminate == currentIndication &&
		enumerations.SubIndicationTryLater == currentSubIndication && !c.isCertificateSuspended() {
		/*
		 * 5) If current time indication/sub indication is INDETERMINATE/TRY_LATER because
		 * the revocation information of the target certificate was not fresh enough:
		 *
		 * a) The building block shall determine from the set of POEs the earliest time at which
		 *    the existence of the signature can be proven.
		 * b) The building block shall run the Revocation Freshness Checker (clause 5.2.5) with
		 *    the corresponding revocation status information, the target certificate and the time
		 *    determined in step a) above.
		 * c) If the checker returns PASSED, the building block shall go to step 7). Otherwise,
		 *    the building block shall return the indication INDETERMINATE, the sub indication
		 *    TRY_LATER and, if returned from the Revocation Freshness Checker, the suggestion
		 *    for when to try the validation again.
		 */
		item = c.revocationIsFresh(item, bestSignatureTime)
	} else {
		/*
		 * 6) In all other cases, the building block shall return the current time indication/sub-indication
		 * together with an explanation of the failure.
		 */
		item = item.SetNextItem(c.currentTimeIndicationCheck())
	}

	/*
	 * 7) The building block shall return the indication and sub-indication contained
	 * in sig_cert_revocation_poe-status.
	 */
	item = item.SetNextItem(c.pastRevocationDataValidationConclusive(sigCertRevocationPoeStatus)) //nolint:staticcheck // mirrors upstream PastSignatureValidation#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

}

// subIndicationOfConclusion reads XmlConclusion#getSubIndication(): the
// generated member is a pointer, whose nil is Java's null.
func subIndicationOfConclusion(conclusion *jaxb.XmlConclusion) enumerations.SubIndication {
	if conclusion.SubIndication == nil {
		return ""
	}
	return conclusion.SubIndication.SubIndication()
}

// isRevocationDataRequired ports the private
// isRevocationDataRequired(CertificateWrapper, SubContext).
func (c *PastSignatureValidation) isRevocationDataRequired(certificate *diagnostic.CertificateWrapper,
	subContext enumerations.SubContext) bool {
	constraint := c.policy.RevocationDataSkipConstraint(c.context, subContext)
	sunsetDateConstraint := c.policy.CertificateSunsetDateConstraint(c.context, subContext)
	return xcv.NewRevocationDataRequiredCheck(c.I18nProvider, c.Result, certificate,
		c.lowestPoeTime(certificate), sunsetDateConstraint, constraint).Process()
}

// checkCertificateRevocationSelectorResult ports the private
// checkCertificateRevocationSelectorResult(XmlCRS).
func (c *PastSignatureValidation) checkCertificateRevocationSelectorResult(
	crsResult *jaxb.XmlCRS) process.ChainItem[*jaxb.XmlPSV] {
	return NewPastSignatureValidationCertificateRevocationSelectorResultCheck(c.I18nProvider, c.Result, crsResult,
		c.WarnLevelRule())
}

// currentTimeIndicationCheck ports the private currentTimeIndicationCheck().
func (c *PastSignatureValidation) currentTimeIndicationCheck() process.ChainItem[*jaxb.XmlPSV] {
	return NewCurrentTimeIndicationCheck(c.I18nProvider, c.Result, c.currentConclusion.Indication.Indication(),
		subIndicationOfConclusion(c.currentConclusion), c.currentConclusion.Errors, c.FailLevelRule())
}

// pastCertificateValidationAcceptableCheck ports the private
// pastCertificateValidationAcceptableCheck(XmlPCV).
func (c *PastSignatureValidation) pastCertificateValidationAcceptableCheck(
	pcvResult *jaxb.XmlPCV) process.ChainItem[*jaxb.XmlPSV] {
	return NewPastCertificateValidationAcceptableCheck(c.I18nProvider, c.Result, pcvResult, c.token.Id(),
		c.currentConclusion.Indication.Indication(), subIndicationOfConclusion(c.currentConclusion),
		c.FailLevelRule())
}

// poeExist ports the private poeExist(Date).
func (c *PastSignatureValidation) poeExist(controlTime *time.Time) *POEExistsCheck {
	return NewPOEExistsCheck(c.I18nProvider, c.Result, c.token, controlTime, c.poe, c.WarnLevelRule())
}

// poeExistNotAfterCARevocationTimeCheck ports the private
// poeExistNotAfterCARevocationTimeCheck(Collection, Date).
func (c *PastSignatureValidation) poeExistNotAfterCARevocationTimeCheck(
	certificateRevocations []*diagnostic.CertificateRevocationWrapper,
	caRevocationTime *time.Time) process.ChainItem[*jaxb.XmlPSV] {
	return NewPOENotAfterCARevocationTimeCheck(c.I18nProvider, c.Result, certificateRevocations, caRevocationTime,
		c.poe, c.FailLevelRule())
}

// pastRevocationDataValidationConclusive ports the private
// pastRevocationDataValidationConclusive(XmlConclusion): the unsupported-Level
// error becomes a panic, the caller being initChain, which cannot propagate one.
func (c *PastSignatureValidation) pastRevocationDataValidationConclusive(
	currentConclusion *jaxb.XmlConclusion) process.ChainItem[*jaxb.XmlPSV] {
	constraint, err := process.GetConstraintOrMaxLevel(
		c.policy.RevocationIssuerNotExpiredConstraint(c.context, enumerations.SubContextSigningCert),
		enumerations.LevelFail)
	if err != nil {
		panic(err)
	}
	return NewPastRevocationDataValidationConclusiveCheck(c.I18nProvider, c.Result, currentConclusion, constraint)
}

// bestSignatureTimeNotBeforeCertificateIssuance ports the private
// bestSignatureTimeNotBeforeCertificateIssuance(Date, CertificateWrapper).
func (c *PastSignatureValidation) bestSignatureTimeNotBeforeCertificateIssuance(bestSignatureTime time.Time,
	signingCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlPSV] {
	// The vpfltvd check keeps Java's nullable Date; the best-signature-time
	// reached here is always the non-null result of POEExtraction#getLowestPOETime.
	return vpfltvd.NewBestSignatureTimeNotBeforeCertificateIssuanceCheck(c.I18nProvider, c.Result,
		&bestSignatureTime, signingCertificate, c.FailLevelRule())
}

// bestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpiration ports
// the private bestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpiration(Date, CertificateWrapper, SubIndication).
func (c *PastSignatureValidation) bestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpiration(
	bestSignatureTime time.Time, signingCertificate *diagnostic.CertificateWrapper,
	currentTimeSubIndication enumerations.SubIndication) process.ChainItem[*jaxb.XmlPSV] {
	return NewBestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck(c.I18nProvider, c.Result,
		bestSignatureTime, signingCertificate, currentTimeSubIndication, c.FailLevelRule())
}

// algorithmsObsolescenceValidation ports the private
// algorithmsObsolescenceValidation(): the unsupported-Context error becomes a
// panic, the caller being initChain, which cannot propagate one.
func (c *PastSignatureValidation) algorithmsObsolescenceValidation() process.ChainItem[*jaxb.XmlPSV] {
	lowestPoeTime := c.lowestPoeTime(c.token)

	algorithmObsolescenceValidation := aov.NewSignatureAlgorithmObsolescenceValidation(
		c.I18nProvider, c.token, c.context, lowestPoeTime, c.policy)
	aovResult := algorithmObsolescenceValidation.Execute()

	position, err := process.GetCryptoPosition(c.context)
	if err != nil {
		panic(err)
	}

	return aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, aovResult, lowestPoeTime, position,
		c.token.Id())
}

// revocationDataAlgorithmsObsolescenceValidation ports the private
// revocationDataAlgorithmsObsolescenceValidation(ChainItem, List) - the
// two-argument entry point of the recursion.
//
// NOTE: it returns the LAST chain item for further appending; it never assigns
// the chain's first item.
func (c *PastSignatureValidation) revocationDataAlgorithmsObsolescenceValidation(item process.ChainItem[*jaxb.XmlPSV],
	signingCertificateRevocations []*diagnostic.CertificateRevocationWrapper) process.ChainItem[*jaxb.XmlPSV] {
	checkedTokens := make([]string, 0)
	return c.revocationDataAlgorithmsObsolescenceValidationRecursive(item, c.token.CertificateChain(),
		signingCertificateRevocations, c.context, &checkedTokens)
}

// revocationDataAlgorithmsObsolescenceValidationRecursive ports the private
// revocationDataAlgorithmsObsolescenceValidation(ChainItem, List, List, Context, List).
//
// Java passes the same ArrayList of checked token ids down the recursion and
// mutates it; the Go form passes a pointer to the slice for the same effect.
func (c *PastSignatureValidation) revocationDataAlgorithmsObsolescenceValidationRecursive(
	item process.ChainItem[*jaxb.XmlPSV], certificateChain []*diagnostic.CertificateWrapper,
	signingCertificateRevocations []*diagnostic.CertificateRevocationWrapper, context enumerations.Context,
	checkedTokens *[]string) process.ChainItem[*jaxb.XmlPSV] {
	for _, certificate := range certificateChain {
		subContext := enumerations.SubContextCACertificate
		if c.token.SigningCertificate().Id() == certificate.Id() {
			subContext = enumerations.SubContextSigningCert
		}
		certificatePoeTime := c.lowestPoeTime(certificate)
		if c.isTrustAnchor(certificate, certificatePoeTime, context, subContext) {
			break
		}
		if containsToken(*checkedTokens, certificate.Id()) {
			continue
		}
		*checkedTokens = append(*checkedTokens, certificate.Id())

		var revocationData []*diagnostic.CertificateRevocationWrapper
		if enumerations.SubContextSigningCert == subContext {
			revocationData = signingCertificateRevocations
		} else {
			revocationData = certificate.CertificateRevocationData()
		}

		latestAcceptableRevocation := process.GetLatestAcceptableRevocationData(c.token, certificate, revocationData,
			c.currentTime, c.bbbs, c.poe)
		if latestAcceptableRevocation != nil && !containsToken(*checkedTokens, latestAcceptableRevocation.Id()) {
			*checkedTokens = append(*checkedTokens, latestAcceptableRevocation.Id())

			revocationPoeTime := c.lowestPoeTime(certificate)

			algorithmObsolescenceValidation := aov.NewRevocationDataAlgorithmObsolescenceValidation(
				c.I18nProvider, &latestAcceptableRevocation.RevocationWrapper, revocationPoeTime, c.policy)
			aovResult := algorithmObsolescenceValidation.Execute()

			position, err := process.GetCryptoPosition(enumerations.ContextRevocation)
			if err != nil {
				panic(err)
			}

			item = item.SetNextItem(aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, aovResult,
				revocationPoeTime, position, c.token.Id()))

			item = c.revocationDataAlgorithmsObsolescenceValidationRecursive(item,
				latestAcceptableRevocation.CertificateChain(), signingCertificateRevocations,
				enumerations.ContextRevocation, checkedTokens)

		}

	}
	return item
}

// containsToken ports java.util.List#contains(Object) for the checked-token id
// list.
func containsToken(checkedTokens []string, tokenId string) bool {
	for _, checked := range checkedTokens {
		if checked == tokenId {
			return true
		}
	}
	return false
}

// isTrustAnchor ports the private
// isTrustAnchor(CertificateWrapper, Date, Context, SubContext).
func (c *PastSignatureValidation) isTrustAnchor(certificateWrapper *diagnostic.CertificateWrapper,
	controlTime time.Time, context enumerations.Context, subContext enumerations.SubContext) bool {
	constraint := c.policy.CertificateSunsetDateConstraint(context, subContext)
	return process.IsTrustAnchor(certificateWrapper, controlTime, constraint)
}

// revocationIsFresh ports the private revocationIsFresh(ChainItem, Date).
//
// NOTE: it returns the LAST chain item for further appending; it never assigns
// the chain's first item.
func (c *PastSignatureValidation) revocationIsFresh(item process.ChainItem[*jaxb.XmlPSV],
	bestSignatureTime time.Time) process.ChainItem[*jaxb.XmlPSV] {
	for _, certificate := range c.token.CertificateChain() {
		subContext := c.subContext(certificate)
		if c.isRevocationDataRequired(certificate, subContext) {
			certificateRevocationSelector := NewPastSignatureValidationCertificateRevocationSelector(
				c.I18nProvider, certificate, c.currentTime, c.bbbs, c.token.Id(), c.poe, c.policy)

			certificateRevocationSelector.Execute()
			latestCertificateRevocation := certificateRevocationSelector.LatestAcceptableCertificateRevocation()

			// Java hands the CertificateRevocationWrapper straight to a
			// RevocationWrapper parameter; the Go wrapper embeds its base by
			// value, so the base is addressed out of it, and a null latest
			// revocation stays a nil *RevocationWrapper.
			var revocationData *diagnostic.RevocationWrapper
			if latestCertificateRevocation != nil {
				revocationData = &latestCertificateRevocation.RevocationWrapper
			}

			rfc := xcv.NewRevocationFreshnessChecker(c.I18nProvider, revocationData, bestSignatureTime, c.context,
				subContext, c.policy)
			xmlRFC := rfc.Execute()
			item = item.SetNextItem(c.checkRevocationFreshnessCheckerResult(xmlRFC))
		}
	}
	return item
}

// subContext ports the private getSubContext(CertificateWrapper).
func (c *PastSignatureValidation) subContext(
	certificateWrapper *diagnostic.CertificateWrapper) enumerations.SubContext {
	if c.token.SigningCertificate().Id() == certificateWrapper.Id() {
		return enumerations.SubContextSigningCert
	}
	return enumerations.SubContextCACertificate
}

// checkRevocationFreshnessCheckerResult ports the private
// checkRevocationFreshnessCheckerResult(XmlRFC).
func (c *PastSignatureValidation) checkRevocationFreshnessCheckerResult(
	rfcResult *jaxb.XmlRFC) process.ChainItem[*jaxb.XmlPSV] {
	return newPSVRevocationFreshnessCheckerResultCheck(c.I18nProvider, c.Result, rfcResult, c.FailLevelRule())
}

// psvRevocationFreshnessCheckerResultCheck is the Go form of the anonymous
// RevocationFreshnessCheckerResultCheck<XmlPSV> subclass declared inside
// checkRevocationFreshnessCheckerResult(XmlRFC).
type psvRevocationFreshnessCheckerResultCheck struct {
	*xcv.RevocationFreshnessCheckerResultCheck[*jaxb.XmlPSV]
}

// newPSVRevocationFreshnessCheckerResultCheck instantiates the anonymous
// subclass.
func newPSVRevocationFreshnessCheckerResultCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlPSV], rfcResult *jaxb.XmlRFC,
	constraint policy.LevelRule) *psvRevocationFreshnessCheckerResultCheck {
	c := &psvRevocationFreshnessCheckerResultCheck{
		RevocationFreshnessCheckerResultCheck: xcv.NewRevocationFreshnessCheckerResultCheck(
			i18nProvider, result, rfcResult, constraint),
	}
	// Re-register with the outer type so the two overrides below dispatch.
	c.InitChainItem(c)
	return c
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// the anonymous subclass's getFailedIndicationForConclusion().
func (c *psvRevocationFreshnessCheckerResultCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the anonymous subclass's getFailedSubIndicationForConclusion().
func (c *psvRevocationFreshnessCheckerResultCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationTryLater
}

// lowestPoeTime ports the private getLowestPoeTime(TokenProxy).
func (c *PastSignatureValidation) lowestPoeTime(token diagnostic.TokenProxy) time.Time {
	return c.poe.GetLowestPOETime(token.Id())
}

// isCertificateSuspended ports the private isCertificateSuspended().
func (c *PastSignatureValidation) isCertificateSuspended() bool {
	for _, certificate := range c.token.CertificateChain() {
		subContext := enumerations.SubContextCACertificate
		if c.token.SigningCertificate().Id() == certificate.Id() {
			subContext = enumerations.SubContextSigningCert
		}
		if c.isTrustAnchor(certificate, c.currentTime, c.context, subContext) {
			break
		}
		revocationData := certificate.CertificateRevocationData()
		latestRevocationData := process.GetLatestAcceptableRevocationData(c.token, certificate, revocationData,
			c.currentTime, c.bbbs, c.poe)
		if latestRevocationData != nil && latestRevocationData.IsRevoked() &&
			enumerations.RevocationReasonCertificateHold == latestRevocationData.Reason() {
			return true
		}
	}
	return false
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of the overridden
// collectMessages(XmlConclusion, XmlConstraint): the generated BlockType member
// is a *XmlBlockType, whose nil is Java's null.
func (c *PastSignatureValidation) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	blockType := jaxb.XmlBlockType("")
	if constraint.BlockType != nil {
		blockType = *constraint.BlockType
	}
	if jaxb.XmlBlockTypePCV == blockType {
		// skip PCV POE message extraction
	} else if jaxb.XmlBlockTypePSVCRS == blockType {
		// skip acceptable revocation message extraction
	} else {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}
