// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/vts/ValidationTimeSliding.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
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
)

// ValidationTimeSliding performs the Validation Time Sliding process.
type ValidationTimeSliding struct {
	*process.ChainBase[*jaxb.XmlVTS]

	// token is the token to process.
	token diagnostic.TokenProxy

	// trustedCertificate is the certificate representing a trust anchor.
	trustedCertificate *diagnostic.CertificateWrapper

	// currentTime is the validation time.
	currentTime time.Time

	// bbbs is the map of all BBBs.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// context is the validation context.
	context enumerations.Context

	// poe is the POE container.
	poe *POEExtraction

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// controlTime is the validation time; nil is Java's null.
	controlTime *time.Time
}

// NewValidationTimeSliding is the default constructor. Port of
// ValidationTimeSliding(I18nProvider, TokenProxy, CertificateWrapper, Date, POEExtraction, Map, Context, ValidationPolicy).
func NewValidationTimeSliding(i18nProvider *i18n.I18nProvider, token diagnostic.TokenProxy,
	trustedCertificate *diagnostic.CertificateWrapper, currentTime time.Time, poe *POEExtraction,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, context enumerations.Context,
	validationPolicy policy.ValidationPolicy) *ValidationTimeSliding {
	xmlVTS := &jaxb.XmlVTS{}
	c := &ValidationTimeSliding{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlVTS,
			&xmlVTS.XmlConstraintsConclusionContent, &xmlVTS.XmlConstraintsConclusionAttrs)),
		token:              token,
		trustedCertificate: trustedCertificate,
		currentTime:        currentTime,
		bbbs:               bbbs,

		context: context,

		poe:    poe,
		policy: validationPolicy,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *ValidationTimeSliding) Title() i18n.MessageTag {
	return i18n.MessageTag_VALIDATION_TIME_SLIDING
}

// InitChain initializes the chain. Port of initChain().
func (c *ValidationTimeSliding) InitChain() {

	tokenBBB := c.bbbs[c.token.Id()]

	var item process.ChainItem[*jaxb.XmlVTS]

	/*
	 * 5.6.2.2.4 Processing
	 *
	 * 1) The building block shall initialize control-time to either:
	 *
	 * a) the trust anchor sunset date when this input is provided, and this date is
	 *    before current date/time; or
	 * b) the current date/time in all other cases.
	 *
	 * NOTE 1: Control-time is an internal variable that is used within the
	 * algorithms and not part of the core results of the validation
	 * process.
	 *
	 * NOTE 2: Initializing control time with current date/time assumes that
	 * the trust anchor is still trusted at the current date/time. The algorithm
	 * can capture the very exotic case where the trust anchor is broken (or becomes
	 * untrusted for any other reason) at a known date by initializing control time
	 * to this date/time.
	 */
	if c.trustedCertificate != nil && c.trustedCertificate.TrustSunsetDate() != nil {
		c.controlTime = c.trustedCertificate.TrustSunsetDate()
		item = c.sunsetDateCheck(c.trustedCertificate)
		c.FirstItem = item

	} else {
		currentTime := c.currentTime
		c.controlTime = &currentTime
	}

	certificateChain := c.token.CertificateChain()
	if utils.IsCollectionNotEmpty(certificateChain) {

		certificateChain = c.reduceChainUntilFirstTrustAnchor(certificateChain)

		/*
		 * 2) For each certificate in the chain starting from the first
		 * certificate (the certificate issued by the trust anchor):
		 */
		certificateChain = utils.ReverseList(certificateChain) // trust anchor -> ... -> signing-certificate

		for _, certificate := range certificateChain {
			if c.isTrustAnchor(certificate) {
				// skip for trust anchor
				continue
			}

			/*
			 * a) The building block shall select revocation status information from
			 * the certificate validation data provided satisfying the following:
			 *
			 * - the revocation status information is consistent with the rules conditioning its use
			 *   to check the revocation status of the considered certificate. In the case of a CRL,
			 *   it shall satisfy the checks specified in IETF RFC 5280 [1], clause 6.3.3 (b) to (l);
			 *   with the exception of the verification if the control-time is within the validity period
			 *   of the certificate of the issuer of the CRL; and
			 *
			 * - the issuance date of the revocation status information is before control time; and
			 *
			 * - the set of POEs contains a proof of existence of the certificate and
			 *   the revocation status information at (or before) control time.
			 *
			 * If at least one revocation status information is selected,
			 * the building block shall go to the next step.
			 * If there is no such information, the building block shall return
			 * the indication INDETERMINATE with the sub indication NO_POE.
			 */

			var latestCompliantRevocation *diagnostic.CertificateRevocationWrapper

			subContext := c.subContext(certificate)
			revocationDataRequiredCheck := c.revocationDataRequired(certificate, subContext)
			revocationDataRequired := revocationDataRequiredCheck.Process()
			if revocationDataRequired {
				revocationIssuerSunsetDateConstraint := c.policy.CertificateSunsetDateConstraint(
					enumerations.Context_REVOCATION, enumerations.SubContext_SIGNING_CERT)
				var certificateRevocationData []*diagnostic.CertificateRevocationWrapper
				if enumerations.SubContext_SIGNING_CERT == subContext {
					certificateRevocationData = process.GetAcceptableRevocationDataForPSVIfExistOrReturnAll(
						c.token, certificate, c.currentTime, c.bbbs, c.poe, revocationIssuerSunsetDateConstraint)
				} else {
					certificateRevocationData = certificate.CertificateRevocationData()
				}

				certificateRevocationSelector := NewValidationTimeSlidingCertificateRevocationSelector(
					c.I18nProvider, certificate, certificateRevocationData, *c.controlTime, c.bbbs,
					xmlCRSId(tokenBBB), c.poe, c.policy)
				xmlCRS := certificateRevocationSelector.Execute()
				c.Result.Value.CRS = append(c.Result.Value.CRS, xmlCRS)

				satisfyingRevocationDataExists := c.satisfyingRevocationDataExists(xmlCRS, certificate, *c.controlTime)
				if item == nil {
					item = satisfyingRevocationDataExists
					c.FirstItem = item
				} else {
					item = item.SetNextItem(satisfyingRevocationDataExists)
				}

				latestCompliantRevocation = certificateRevocationSelector.LatestAcceptableCertificateRevocation()

			} else {
				if item == nil {
					item = revocationDataRequiredCheck
					c.FirstItem = item
				} else {
					item = item.SetNextItem(revocationDataRequiredCheck)
				}
			}

			if latestCompliantRevocation == nil {
				// skip revocation checks
			} else if latestCompliantRevocation.IsRevoked() {
				/*
				 * b) If the certificate is marked as revoked in any of the revocation status information
				 * found in the previous step, the building block shall perform the following steps:
				 *
				 * - select the revocation status information that has been issued the latest;
				 *
				 * - set control time to the revocation time whenever the validation policy requires
				 *   to use the shell model; or, when the validation policy requires to use the chain model and
				 *   the revocation reason is key compromise or unknown.
				 *
				 * - go to step d).
				 */
				validationModel := c.policy.ValidationModel()
				revocationReason := latestCompliantRevocation.Reason()
				// NOTE : HYBRID model is treated as CHAIN for Signing Cert and as SHELL for CAs
				if enumerations.ValidationModel_SHELL == validationModel ||
					(enumerations.ValidationModel_HYBRID == validationModel &&
						enumerations.SubContext_CA_CERTIFICATE == subContext) ||
					enumerations.RevocationReason_KEY_COMPROMISE == revocationReason ||
					enumerations.RevocationReason_UNSPECIFIED == revocationReason {
					c.controlTime = latestCompliantRevocation.RevocationDate()
				}
			} else {
				/*
				 * c) If the certificate is not marked as revoked in all of the revocation data found in step a),
				 * the building block shall select the revocation data that has been issued the latest,
				 * run the Revocation Freshness Checker with that revocation data, the certificate for which
				 * the revocation status is being checked and the control time. If it returns FAILED,
				 * the building block shall set control time to the time that is the earliest between time
				 * A and time B, where time A is the current value of control time and time B is
				 * the issuance time of the revocation status information contained within the revocation data.
				 * Otherwise, the building block shall not change the value of control time.
				 */
				rfc := xcv.NewRevocationFreshnessChecker(c.I18nProvider, &latestCompliantRevocation.RevocationWrapper,
					*c.controlTime, c.context, subContext, c.policy)
				execute := rfc.Execute()
				if execute.Conclusion != nil &&
					enumerations.Indication_FAILED == execute.Conclusion.Indication.Indication() {
					thisUpdate := latestCompliantRevocation.ThisUpdate()
					if thisUpdate.Before(*c.controlTime) {
						c.controlTime = thisUpdate
					}
				}
			}

			/*
			 * d) The building block shall apply the cryptographic constraints to the certificate and
			 * the revocation status information against the control time. If the certificate
			 * (or the revocation status information) does not match these constraints, the building block shall
			 * set control time to the latest time up to which the listed algorithms were all considered reliable.
			 */
			var cryptoNotAfterDate *time.Time

			certificateAOV := c.certificateAlgorithmObsolescenceResult(certificate, *c.controlTime, subContext)
			if !c.IsValidConclusion(certificateAOV.Conclusion) {
				cryptographicValidation := process.GetFailCryptographicValidation(certificateAOV)
				cryptoNotAfterDate = notAfterOf(cryptographicValidation)
			}

			if latestCompliantRevocation != nil {
				revocationAOV := c.revocationDataAlgorithmObsolescenceResult(
					&latestCompliantRevocation.RevocationWrapper, *c.controlTime)
				if !c.IsValidConclusion(revocationAOV.Conclusion) {
					cryptographicValidation := process.GetFailCryptographicValidation(revocationAOV)
					revCryptoNotAfter := notAfterOf(cryptographicValidation)
					if cryptoNotAfterDate == nil ||
						(revCryptoNotAfter != nil && revCryptoNotAfter.Before(*cryptoNotAfterDate)) {
						cryptoNotAfterDate = revCryptoNotAfter
					}
				}
			}

			if cryptoNotAfterDate != nil && cryptoNotAfterDate.Before(*c.controlTime) {
				c.controlTime = cryptoNotAfterDate
			}

			/*
			 * e) The building block shall continue with the next certificate in the chain or,
			 * if no further certificate exists, the building block shall return the status
			 * indication PASSED and the calculated control time.
			 */

		}
	}

	controlTimeConclusiveCheck := c.controlTimeConclusive(c.controlTime)
	if item == nil {
		item = controlTimeConclusiveCheck
		c.FirstItem = item
	} else {
		item = item.SetNextItem(controlTimeConclusiveCheck) //nolint:staticcheck // mirrors upstream ValidationTimeSliding#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
	}
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// addAdditionalInfo(). Java dereferences trustedCertificate unguarded, so a
// null one raises a NullPointerException; the nil pointer panics here in its
// place.
func (c *ValidationTimeSliding) AddAdditionalInfo() {
	if c.controlTime != nil {
		c.Result.Value.ControlTime = jaxb.NewXSDateTime(*c.controlTime)
	} else {
		c.Result.Value.ControlTime = nil
	}
	trustAnchor := c.trustedCertificate.Id()
	c.Result.Value.TrustAnchor = &trustAnchor
}

// subContext ports the private getSubContext(CertificateWrapper).
func (c *ValidationTimeSliding) subContext(certificate *diagnostic.CertificateWrapper) enumerations.SubContext {
	if c.token.SigningCertificate().Id() == certificate.Id() {
		return enumerations.SubContext_SIGNING_CERT
	}
	return enumerations.SubContext_CA_CERTIFICATE
}

// reduceChainUntilFirstTrustAnchor ports the private
// reduceChainUntilFirstTrustAnchor(List).
func (c *ValidationTimeSliding) reduceChainUntilFirstTrustAnchor(
	originalCertificateChain []*diagnostic.CertificateWrapper) []*diagnostic.CertificateWrapper {
	result := make([]*diagnostic.CertificateWrapper, 0)
	for _, cert := range originalCertificateChain {
		result = append(result, cert)
		if c.isTrustAnchor(cert) {
			break
		}
	}
	return result
}

// isTrustAnchor ports the private isTrustAnchor(CertificateWrapper), which is
// certificate.equals(trustedCertificate). A null trustedCertificate answers
// false there (AbstractTokenProxy#equals(null)), which the explicit nil guard
// reproduces - a typed nil pointer is not a nil Go interface.
func (c *ValidationTimeSliding) isTrustAnchor(certificate *diagnostic.CertificateWrapper) bool {
	if c.trustedCertificate == nil {
		return false
	}
	return certificate.Equals(c.trustedCertificate)
}

// sunsetDateCheck ports the private sunsetDateCheck(CertificateWrapper).
func (c *ValidationTimeSliding) sunsetDateCheck(
	trustedCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlVTS] {
	return NewSunsetDateCheck(c.I18nProvider, c.Result, trustedCertificate, c.FailLevelRule())
}

// revocationDataRequired ports the private
// revocationDataRequired(CertificateWrapper, SubContext).
func (c *ValidationTimeSliding) revocationDataRequired(certificate *diagnostic.CertificateWrapper,
	subContext enumerations.SubContext) *xcv.RevocationDataRequiredCheck[*jaxb.XmlVTS] {
	constraint := c.policy.RevocationDataSkipConstraint(c.context, subContext)
	sunsetDateConstraint := c.policy.CertificateSunsetDateConstraint(c.context, subContext)
	return xcv.NewRevocationDataRequiredCheck(c.I18nProvider, c.Result, certificate, c.currentTime,
		sunsetDateConstraint, constraint)
}

// satisfyingRevocationDataExists ports the private
// satisfyingRevocationDataExists(XmlCRS, CertificateWrapper, Date).
func (c *ValidationTimeSliding) satisfyingRevocationDataExists(crsResult *jaxb.XmlCRS,
	certificateWrapper *diagnostic.CertificateWrapper, controlTime time.Time) process.ChainItem[*jaxb.XmlVTS] {
	return NewSatisfyingRevocationDataExistsCheck(c.I18nProvider, c.Result, crsResult, certificateWrapper,
		controlTime, c.FailLevelRule())
}

// controlTimeConclusive ports the private controlTimeConclusive(Date).
func (c *ValidationTimeSliding) controlTimeConclusive(controlTime *time.Time) process.ChainItem[*jaxb.XmlVTS] {
	return NewControlTimeCheck(c.I18nProvider, c.Result, controlTime, c.FailLevelRule())
}

// certificateAlgorithmObsolescenceResult ports the private
// getCertificateAlgorithmObsolescenceResult(CertificateWrapper, Date, SubContext).
func (c *ValidationTimeSliding) certificateAlgorithmObsolescenceResult(
	certificateWrapper *diagnostic.CertificateWrapper, controlTime time.Time,
	subContext enumerations.SubContext) *jaxb.XmlAOV {
	a := aov.NewCertificateAlgorithmObsolescenceValidation(
		c.I18nProvider, certificateWrapper, c.context, subContext, controlTime, c.policy)
	return a.Execute()
}

// revocationDataAlgorithmObsolescenceResult ports the private
// getRevocationDataAlgorithmObsolescenceResult(RevocationWrapper, Date).
func (c *ValidationTimeSliding) revocationDataAlgorithmObsolescenceResult(
	revocationWrapper *diagnostic.RevocationWrapper, controlTime time.Time) *jaxb.XmlAOV {
	a := aov.NewRevocationDataAlgorithmObsolescenceValidation(
		c.I18nProvider, revocationWrapper, controlTime, c.policy)
	return a.Execute()
}

// notAfterOf reads XmlCryptographicValidation#getNotAfter() off a possibly null
// validation: the generated member is a *XSDateTime, whose nil is Java's null.
func notAfterOf(cryptographicValidation *jaxb.XmlCryptographicValidation) *time.Time {
	if cryptographicValidation == nil || cryptographicValidation.NotAfter == nil {
		return nil
	}
	notAfter := cryptographicValidation.NotAfter.Time()
	return &notAfter
}

// xmlCRSId reads XmlBasicBuildingBlocks#getId(): the generated member is a
// plain Go string, so a Java null and an absent id render identically.
func xmlCRSId(bbb *jaxb.XmlBasicBuildingBlocks) string {
	return bbb.Id
}
