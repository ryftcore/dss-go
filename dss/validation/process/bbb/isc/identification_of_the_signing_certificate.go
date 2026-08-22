// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/isc/IdentificationOfTheSigningCertificate.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.validation.process.bbb.isc and its .checks subpackage
// flatten into this single Go package isc (no name collisions), so the checks
// are referenced unqualified.
//
// Every check constructor takes *process.Result[*jaxb.XmlISC] where Java takes
// the XmlISC itself: see chain.go for why the generated result object has to be
// bound to the JAXB base structs it embeds.
package isc

import (
	"slices"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// IdentificationOfTheSigningCertificate is 5.2.3 Identification of the signing
// certificate. This building block is responsible for identifying the signing
// certificate that will be used to validate the signature.
type IdentificationOfTheSigningCertificate struct {
	*process.ChainBase[*jaxb.XmlISC]

	// token is the token to verify.
	token diagnostic.TokenProxy

	// validationPolicy is the validation policy.
	validationPolicy policy.ValidationPolicy

	// context is the validation context.
	context enumerations.Context
}

// NewIdentificationOfTheSigningCertificate is the default constructor. Port of
// IdentificationOfTheSigningCertificate(Provider, TokenProxy, Context, ValidationPolicy).
func NewIdentificationOfTheSigningCertificate(i18nProvider *i18n.Provider, token diagnostic.TokenProxy,
	context enumerations.Context, validationPolicy policy.ValidationPolicy) *IdentificationOfTheSigningCertificate {
	xmlISC := &jaxb.XmlISC{}
	c := &IdentificationOfTheSigningCertificate{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlISC,
			&xmlISC.XmlConstraintsConclusionContent, &xmlISC.XmlConstraintsConclusionAttrs)),
		token:            token,
		context:          context,
		validationPolicy: validationPolicy,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *IdentificationOfTheSigningCertificate) Title() i18n.MessageTag {
	return i18n.MessageTagIdentificationOfTheSigningCertificate
}

// InitChain initializes the chain. Port of initChain().
func (c *IdentificationOfTheSigningCertificate) InitChain() {

	/*
	 * The common way to unambiguously identify the signing certificate is by using a property/attribute of the
	 * signature containing a reference to it (see clause 4.2.5.2). The certificate can either be found in the
	 * signature or it can be obtained using external sources. The signing certificate can also be provided by the
	 * DA. If no certificate can be retrieved, the building block shall return the indication INDETERMINATE and the
	 * sub-indication NO_SIGNING_CERTIFICATE_FOUND.
	 */
	item := c.signingCertificateRecognition()
	c.FirstItem = item

	isSignature := enumerations.ContextSignature == c.context ||
		enumerations.ContextCounterSignature == c.context ||
		enumerations.ContextKeyBindingSignature == c.context
	isTimestamp := enumerations.ContextTimestamp == c.context

	if isSignature || isTimestamp {
		/*
		 * 1) If the signature format used contains a way to directly identify the reference to the signers'
		 * certificate in the attribute, the building block shall check that the digest of the certificate
		 * referenced matches the result of digesting the signing certificate with the algorithm indicated; if they
		 * match, the building block shall return the signing certificate. Otherwise, the building block shall go to
		 * step 2.
		 */

		if !c.token.IsSigningCertificateReferencePresent() || c.token.SigningCertificate() == nil {
			return
		}

		/*
		 * 2) The building block shall take the first reference and shall check that the digest of the certificate
		 * referenced matches the result of digesting the signing certificate with the algorithm indicated. If they
		 * do not match, the building block shall take the next element and shall repeat this step until a matching
		 * element has been found or all elements have been checked. If they do match, the building block shall
		 * continue with step 3. If the last element is reached without finding any match, the validation of this
		 * property shall be taken as failed and the building block shall return the indication INDETERMINATE with
		 * the sub-indication NO_SIGNING_CERTIFICATE_FOUND.
		 */
		item = item.SetNextItem(c.digestValuePresent())

		item = item.SetNextItem(c.digestValueMatch())

		/*
		 * 3) If the issuer and the serial number are additionally present in that reference, the details of the
		 * issuer's name and the serial number of the IssuerSerial element may be compared with those indicated in
		 * the signing certificate: if they do not match, an additional warning shall be returned with the output.
		 */
		signingCertificateRef := c.token.SigningCertificateReference()
		if signingCertificateRef != nil && signingCertificateRef.IsIssuerSerialPresent() {
			item = item.SetNextItem(c.issuerSerialMatch()) //nolint:staticcheck // mirrors upstream IdentificationOfTheSigningCertificate#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
		}
	}
}

// AddAdditionalInfo adds the certificate chain to the result. Port of the
// overridden addAdditionalInfo().
func (c *IdentificationOfTheSigningCertificate) AddAdditionalInfo() {
	c.ChainBase.AddAdditionalInfo()

	// Java guards on "token.getCertificateChain() != null", but
	// AbstractTokenProxy#getCertificateChain always returns a (possibly empty)
	// list, so the guard is always taken and an empty chain still produces an
	// empty XmlCertificateChain element. Go's zero-length slice is nil, so
	// porting the guard literally would drop that element instead.
	{
		certificateChain := &jaxb.XmlCertificateChain{}
		for _, certificate := range c.token.CertificateChain() {
			chainItem := &jaxb.XmlChainItem{}
			chainItem.Id = certificate.Id()
			sources := certificate.Sources()
			if slices.Contains(sources, enumerations.CertificateSourceTypeTrustedList) {
				chainItem.Source = jaxb.CertificateSourceTypeValue(enumerations.CertificateSourceTypeTrustedList)
			} else if slices.Contains(sources, enumerations.CertificateSourceTypeTrustedStore) {
				chainItem.Source = jaxb.CertificateSourceTypeValue(enumerations.CertificateSourceTypeTrustedStore)
			} else {
				chainItem.Source = jaxb.CertificateSourceTypeValue(sources[0])
			}
			certificateChain.ChainItem = append(certificateChain.ChainItem, chainItem)
		}
		c.Result.Value.CertificateChain = certificateChain
	}
}

// signingCertificateRecognition ports the private signingCertificateRecognition().
func (c *IdentificationOfTheSigningCertificate) signingCertificateRecognition() process.ChainItem[*jaxb.XmlISC] {
	constraint := c.validationPolicy.SigningCertificateRecognitionConstraint(c.context)
	return NewSigningCertificateRecognitionCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// digestValuePresent ports the private digestValuePresent().
func (c *IdentificationOfTheSigningCertificate) digestValuePresent() process.ChainItem[*jaxb.XmlISC] {
	constraint := c.validationPolicy.SigningCertificateDigestValuePresentConstraint(c.context)
	return NewDigestValuePresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// digestValueMatch ports the private digestValueMatch().
func (c *IdentificationOfTheSigningCertificate) digestValueMatch() process.ChainItem[*jaxb.XmlISC] {
	constraint := c.validationPolicy.SigningCertificateDigestValueMatchConstraint(c.context)
	return NewDigestValueMatchCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// issuerSerialMatch ports the private issuerSerialMatch().
func (c *IdentificationOfTheSigningCertificate) issuerSerialMatch() process.ChainItem[*jaxb.XmlISC] {
	constraint := c.validationPolicy.SigningCertificateIssuerSerialMatchConstraint(c.context)
	return NewIssuerSerialMatchCheck(c.I18nProvider, c.Result, c.token, constraint)
}
