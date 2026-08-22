// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/RevocationAcceptanceChecker.java (DSS 6.5.RC1).
//
// See x509_certificate_validation.go for the package-flattening note. The
// crs <-> rac mutual recursion (a selector spawns an acceptance checker, which
// may spawn a selector for the revocation issuer's certificate) is a plain
// function call inside the flattened package; the shared validatedTokens set
// remains the loop breaker it is upstream.
package xcv

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/cv"
)

// RevocationAcceptanceChecker checks if the revocation is acceptable and can be
// used.
type RevocationAcceptanceChecker struct {
	*process.ChainBase[*jaxb.XmlRAC]

	// certificate is the certificate in question.
	certificate *diagnostic.CertificateWrapper

	// revocationData is the revocation data.
	revocationData *diagnostic.CertificateRevocationWrapper

	// controlTime is the validation time.
	controlTime time.Time

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// validatedTokens is the internal set of processed tokens (avoids infinite
	// loop). Java's Set<String> is shared by reference with the selector that
	// spawned this checker; the Go map is the same shared object.
	validatedTokens map[string]struct{}
}

// NewRevocationAcceptanceChecker is the default constructor. Port of
// RevocationAcceptanceChecker(I18nProvider, CertificateWrapper,
// CertificateRevocationWrapper, Date, ValidationPolicy, Set).
//
// Java copies the revocation's thisUpdate and productionDate into the result;
// both are nullable Dates there, and the generated Go model carries the two
// members as plain (non-pointer) XSDateTime values, so a Java null maps to the
// zero time. The one reader of the distinction,
// RevocationAcceptanceCheckerResultCheck#buildAdditionalInfo, tests the zero
// time in its place - see that file.
func NewRevocationAcceptanceChecker(i18nProvider *i18n.I18nProvider, certificate *diagnostic.CertificateWrapper,
	revocationData *diagnostic.CertificateRevocationWrapper, controlTime time.Time,
	validationPolicy policy.ValidationPolicy, validatedTokens map[string]struct{}) *RevocationAcceptanceChecker {
	xmlRAC := &jaxb.XmlRAC{}
	c := &RevocationAcceptanceChecker{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlRAC,
			&xmlRAC.XmlConstraintsConclusionContent, &xmlRAC.XmlConstraintsConclusionAttrs)),
		certificate:     certificate,
		revocationData:  revocationData,
		controlTime:     controlTime,
		policy:          validationPolicy,
		validatedTokens: validatedTokens,
	}

	id := revocationData.Id()
	xmlRAC.Id = &id
	if thisUpdate := revocationData.ThisUpdate(); thisUpdate != nil {
		xmlRAC.RevocationThisUpdate = jaxb.XSDateTime(*thisUpdate)
	}
	if productionDate := revocationData.ProductionDate(); productionDate != nil {
		xmlRAC.RevocationProductionDate = jaxb.XSDateTime(*productionDate)
	}
	c.validatedTokens[certificate.Id()] = struct{}{}

	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *RevocationAcceptanceChecker) Title() i18n.MessageTag {
	return i18n.MessageTagRAC
}

// InitChain initializes the chain. Port of initChain().
func (c *RevocationAcceptanceChecker) InitChain() {

	item := c.revocationDataKnown()
	c.FirstItem = item

	item = item.SetNextItem(c.issuerCertificateKnown())

	item = item.SetNextItem(c.prospectiveCertificateChain(c.revocationData.SigningCertificate()))

	item = item.SetNextItem(c.revocationDataIntact())

	item = item.SetNextItem(c.thisUpdate())

	/*
	 * certHash extension can be present in an OCSP Response. If present, a digest match indicates the OCSP
	 * responder knows the certificate as we have it, and so also its revocation state
	 */
	if enumerations.RevocationTypeOCSP == c.revocationData.RevocationType() {

		item = item.SetNextItem(c.issuerValidAtProductionTime())

		item = item.SetNextItem(c.revocationResponderIdMatch())

		item = item.SetNextItem(c.revocationCertHashPresent())

		if c.revocationData.IsCertHashExtensionPresent() {
			item = item.SetNextItem(c.revocationCertHashMatch())
		}

		item = item.SetNextItem(c.selfIssuedOcsp())

	}

	item = item.SetNextItem(c.revocationAfterCertIssuance())

	item = item.SetNextItem(c.revocationHasInformationAboutCertificate())

	for _, revocationCertificate := range c.revocationData.CertificateChain() {
		subContext := enumerations.SubContextCACertificate
		if c.revocationData.SigningCertificate().Id() == revocationCertificate.Id() {
			subContext = enumerations.SubContextSigningCert
		}

		if c.isTrustAnchor(revocationCertificate, subContext) {
			break
		}

		if c.isTokenValidated(revocationCertificate) {
			continue
		}

		item = item.SetNextItem(c.certificateIntact(revocationCertificate))

		if revocationCertificate.IsSelfSigned() {
			item = item.SetNextItem(c.selfSigned(revocationCertificate))
		}

		revocationDataRequired := c.revocationDataRequired(revocationCertificate, subContext)
		if revocationDataRequired.Process() {

			item = item.SetNextItem(c.revocationDataPresentForRevocationChain(revocationCertificate, subContext))

			if utils.IsCollectionNotEmpty(revocationCertificate.CertificateRevocationData()) {

				certificateRevocationSelector := NewCertificateRevocationSelectorWithValidatedTokens(
					c.I18nProvider, revocationCertificate, c.controlTime, c.policy, c.validatedTokens)
				xmlCRS := certificateRevocationSelector.Execute()
				c.Result.Value.CRS = xmlCRS

				item = item.SetNextItem(c.checkCertificateRevocationSelectorResult(xmlCRS, subContext))

			}

		} else {
			item = item.SetNextItem(revocationDataRequired)
		}

	}

}

// revocationDataKnown ports the private revocationDataKnown().
func (c *RevocationAcceptanceChecker) revocationDataKnown() process.ChainItem[*jaxb.XmlRAC] {
	return NewRevocationDataKnownCheck(c.I18nProvider, c.Result, c.revocationData, c.policy.UnknownStatusConstraint())
}

// revocationResponderIdMatch ports the private revocationResponderIdMatch().
func (c *RevocationAcceptanceChecker) revocationResponderIdMatch() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.OCSPResponseResponderIdMatchConstraint()
	return NewRevocationResponderIdMatchCheck(c.I18nProvider, c.Result, &c.revocationData.RevocationWrapper, constraint)
}

// revocationCertHashPresent ports the private revocationCertHashPresent().
func (c *RevocationAcceptanceChecker) revocationCertHashPresent() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.OCSPResponseCertHashPresentConstraint()
	return NewRevocationCertHashPresenceCheck(c.I18nProvider, c.Result, &c.revocationData.RevocationWrapper, constraint)
}

// revocationCertHashMatch ports the private revocationCertHashMatch().
func (c *RevocationAcceptanceChecker) revocationCertHashMatch() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.OCSPResponseCertHashMatchConstraint()
	return NewRevocationCertHashMatchCheck(c.I18nProvider, c.Result, &c.revocationData.RevocationWrapper, constraint)
}

// selfIssuedOcsp ports the private selfIssuedOcsp().
func (c *RevocationAcceptanceChecker) selfIssuedOcsp() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.SelfIssuedOCSPConstraint()
	return NewSelfIssuedOCSPCheck(c.I18nProvider, c.Result, c.certificate, &c.revocationData.RevocationWrapper, constraint)
}

// thisUpdate ports the private thisUpdate().
func (c *RevocationAcceptanceChecker) thisUpdate() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.ThisUpdatePresentConstraint()
	return NewThisUpdatePresenceCheck(c.I18nProvider, c.Result, &c.revocationData.RevocationWrapper, constraint)
}

// issuerCertificateKnown ports the private issuerCertificateKnown().
func (c *RevocationAcceptanceChecker) issuerCertificateKnown() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.RevocationIssuerKnownConstraint()
	return NewRevocationIssuerKnownCheck(c.I18nProvider, c.Result, &c.revocationData.RevocationWrapper, constraint)
}

// issuerValidAtProductionTime ports the private issuerValidAtProductionTime().
func (c *RevocationAcceptanceChecker) issuerValidAtProductionTime() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.RevocationIssuerValidAtProductionTimeConstraint()
	return NewRevocationIssuerValidAtProductionTimeCheck(c.I18nProvider, c.Result, &c.revocationData.RevocationWrapper, constraint)
}

// revocationAfterCertIssuance ports the private revocationAfterCertIssuance().
func (c *RevocationAcceptanceChecker) revocationAfterCertIssuance() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.RevocationAfterCertificateIssuanceConstraint()
	return NewRevocationAfterCertificateIssuanceCheck(c.I18nProvider, c.Result, c.certificate, &c.revocationData.RevocationWrapper, constraint)
}

// revocationHasInformationAboutCertificate ports the private
// revocationHasInformationAboutCertificate().
func (c *RevocationAcceptanceChecker) revocationHasInformationAboutCertificate() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.RevocationHasInformationAboutCertificateConstraint()
	return NewRevocationHasInformationAboutCertificateCheck(c.I18nProvider, c.Result, c.certificate, &c.revocationData.RevocationWrapper, constraint)
}

// revocationDataIntact ports the private revocationDataIntact().
func (c *RevocationAcceptanceChecker) revocationDataIntact() process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.SignatureIntactConstraint(enumerations.ContextRevocation)
	return cv.NewSignatureIntactCheck(c.I18nProvider, c.Result, c.revocationData, enumerations.ContextRevocation, constraint)
}

// prospectiveCertificateChain ports the private
// prospectiveCertificateChain(CertificateWrapper).
func (c *RevocationAcceptanceChecker) prospectiveCertificateChain(
	signingCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.ProspectiveCertificateChainConstraint(enumerations.ContextRevocation)
	return NewProspectiveCertificateChainCheck(c.I18nProvider, c.Result, signingCertificate, enumerations.ContextRevocation, constraint)
}

// isTokenValidated ports the private isTokenValidated(TokenProxy).
func (c *RevocationAcceptanceChecker) isTokenValidated(token diagnostic.TokenProxy) bool {
	_, validated := c.validatedTokens[token.Id()]
	c.validatedTokens[token.Id()] = struct{}{}
	return validated
}

// certificateIntact ports the private certificateIntact(CertificateWrapper).
func (c *RevocationAcceptanceChecker) certificateIntact(
	certificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.SignatureIntactConstraint(enumerations.ContextCertificate)
	return cv.NewSignatureIntactWithIdCheck(c.I18nProvider, c.Result, certificate, enumerations.ContextCertificate, constraint)
}

// selfSigned ports the private selfSigned(CertificateWrapper).
func (c *RevocationAcceptanceChecker) selfSigned(
	certificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlRAC] {
	return NewCertificateSelfSignedCheck(c.I18nProvider, c.Result, certificate, c.WarnLevelRule())
}

// revocationDataPresentForRevocationChain ports the private
// revocationDataPresentForRevocationChain(CertificateWrapper, SubContext).
func (c *RevocationAcceptanceChecker) revocationDataPresentForRevocationChain(
	certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.RevocationDataAvailableConstraint(enumerations.ContextRevocation, subContext)
	return NewRevocationIssuerRevocationDataAvailableCheck(c.I18nProvider, c.Result, certificate, constraint)
}

// checkCertificateRevocationSelectorResult ports the private
// checkCertificateRevocationSelectorResult(XmlCRS, SubContext).
func (c *RevocationAcceptanceChecker) checkCertificateRevocationSelectorResult(crsResult *jaxb.XmlCRS,
	subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlRAC] {
	constraint := c.policy.AcceptableRevocationDataFoundConstraint(enumerations.ContextRevocation, subContext)
	return NewCertificateRevocationSelectorResultCheck(c.I18nProvider, c.Result, crsResult, constraint)
}

// revocationDataRequired ports the private
// revocationDataRequired(CertificateWrapper, SubContext).
func (c *RevocationAcceptanceChecker) revocationDataRequired(certificate *diagnostic.CertificateWrapper,
	subContext enumerations.SubContext) *RevocationDataRequiredCheck[*jaxb.XmlRAC] {
	constraint := c.policy.RevocationDataSkipConstraint(enumerations.ContextRevocation, subContext)
	sunsetDateConstraint := c.policy.CertificateSunsetDateConstraint(enumerations.ContextRevocation, subContext)
	return NewRevocationDataRequiredCheck(c.I18nProvider, c.Result, certificate, c.controlTime, sunsetDateConstraint, constraint)
}

// isTrustAnchor ports the private isTrustAnchor(CertificateWrapper, SubContext).
func (c *RevocationAcceptanceChecker) isTrustAnchor(certificateWrapper *diagnostic.CertificateWrapper,
	subContext enumerations.SubContext) bool {
	sunsetDateConstraint := c.policy.CertificateSunsetDateConstraint(enumerations.ContextRevocation, subContext)
	return process.IsTrustAnchor(certificateWrapper, c.controlTime, sunsetDateConstraint)
}

// CollectAdditionalMessages fills additional messages into the conclusion. Port
// of the overridden collectAdditionalMessages(XmlConclusion).
func (c *RevocationAcceptanceChecker) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	c.ChainBase.CollectAdditionalMessages(conclusion)

	crs := c.Result.Value.CRS
	if crs != nil {
		c.CollectAllMessages(conclusion, crs.Conclusion)
	}
}
