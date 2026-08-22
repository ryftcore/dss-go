// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/pcv/PastCertificateValidation.java (DSS 6.5.RC1).
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
)

// PastCertificateValidation validates certificate in a past.
type PastCertificateValidation struct {
	*process.ChainBase[*jaxb.XmlPCV]

	// token is the token to be validated.
	token diagnostic.TokenProxy

	// bbbs is the map of all BBBs.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// poe is the POE container.
	poe *POEExtraction

	// currentTime is the validation time.
	currentTime time.Time

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// context is the validation context.
	context enumerations.Context

	// controlTime is the control time; nil is Java's null.
	controlTime *time.Time
}

// NewPastCertificateValidation is the default constructor. Port of
// PastCertificateValidation(Provider, TokenProxy, Map, POEExtraction, Date, ValidationPolicy, Context).
func NewPastCertificateValidation(i18nProvider *i18n.Provider, token diagnostic.TokenProxy,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, poe *POEExtraction, currentTime time.Time,
	validationPolicy policy.ValidationPolicy, context enumerations.Context) *PastCertificateValidation {
	xmlPCV := &jaxb.XmlPCV{}
	c := &PastCertificateValidation{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlPCV,
			&xmlPCV.XmlConstraintsConclusionContent, &xmlPCV.XmlConstraintsConclusionAttrs)),

		token:       token,
		bbbs:        bbbs,
		poe:         poe,
		currentTime: currentTime,

		policy:  validationPolicy,
		context: context,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *PastCertificateValidation) Title() i18n.MessageTag {
	return i18n.MessageTagPastCertificateValidation
}

// InitChain initializes the chain. Port of initChain().
func (c *PastCertificateValidation) InitChain() {

	/*
	 * 1) The building block shall build a new prospective certificate chain that
	 * has not yet been evaluated:
	 * a) If no new chain can be built, the building block shall return the current
	 *    status and the last chain built or, if no chain was built, the indication
	 *    INDETERMINATE with the sub-indication NO_CERTIFICATE_CHAIN_FOUND.
	 * b) Otherwise, the building block shall go to the next ste
	 */
	item := c.prospectiveCertificateChain()
	c.FirstItem = item

	/*
	 * 2) The building block shall run the Certification Path Validation of
	 * IETF RFC 5280 [1], clause 6.1, with the following inputs: the prospective
	 * certificate chain built in the previous step, the trust anchor used in the
	 * previous step, the X.509 parameters provided in the inputs and either:
	 *
	 * i) when the validation policy requires to use the shell model, a date from
	 * the intersection of the validity intervals of all the certificates in the
	 * prospective certificate chain; or
	 *
	 * ii) when the validation policy requires to use the chain model, a date from
	 * the validity of the signer's certificate.
	 *
	 * The validation shall not include revocation checking nor verifying that current
	 * time is before a trust anchor sunset date when the X.509 validation constraints
	 * define such a sunset date:
	 *
	 * a) If the certificate path validation returns PASSED, the building block
	 * shall go to the next step.
	 *
	 * b) If the certificate path validation returns a failure indication, the building
	 * block shall set the current status to
	 * INDETERMINATE/CERTIFICATE_CHAIN_GENERAL_FAILURE and shall go to step 1.
	 */

	// Certificates are validated before (see (Sub)X509CertificateValidation)

	// TODO : process different validation models (?)

	/*
	 * 3) The building block shall perform the validation time sliding process as per
	 * clause 5.6.2.2 with the following inputs: the prospective chain, the set of POEs,
	 * the set of certificate validation data, the sunset date of the trust anchor from
	 * which the current chain has been built when the X.509 validation constraint specify
	 * such a date, and the cryptographic constraints. If it outputs a success indication,
	 * the building block shall go to the next step. Otherwise, the building block shall
	 * set the current status to the returned indication and sub-indication and shall go
	 * back to step 1).
	 */

	var trustedCertificate *diagnostic.CertificateWrapper

	vtsList := make([]*jaxb.XmlVTS, 0)
	for _, certificateWrapper := range c.token.CertificateChain() {
		if certificateWrapper.IsTrusted() {
			trustedCertificate = certificateWrapper
			vts := c.vtsResult(trustedCertificate)
			vtsList = append(vtsList, vts)

			item = item.SetNextItem(c.validationTimeSliding(vts, trustedCertificate))

			if trustedCertificate.IsSelfSigned() ||
				trustedCertificate.TrustSunsetDate() == nil || !trustedCertificate.IsTrustedChain() {
				// no sunset date change -> no reason to restart VTS
				break
			}
		}
	}

	vts := c.bestValidationTimeSliding(vtsList)
	if vts != nil {
		bbb := c.bbbs[c.token.Id()]
		bbb.VTS = vts
		if c.IsValid(&vts.XmlConstraintsConclusionContent) {
			// XmlVTS#getControlTime(): the generated member is a *XSDateTime,
			// whose nil is Java's null.
			if vts.ControlTime != nil {
				controlTime := vts.ControlTime.Time()
				c.controlTime = &controlTime
			} else {
				c.controlTime = nil
			}
		}
	}

	item = item.SetNextItem(c.successfulValidationTimeSlidingFound(vts))

	/*
	 * 4) The building block shall apply the X.509 validation constraints to the chain.
	 * If the chain does not match these constraints, the building block shall set the
	 * current status to INDETERMINATE/CHAIN_CONSTRAINTS_FAILURE and shall go to step 1).
	 */
	if c.controlTime != nil {

		item = item.SetNextItem(c.cryptographicCheck(*c.controlTime)) //nolint:staticcheck // mirrors upstream PastCertificateValidation#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

	}

	/*
	 * 5) The building block shall return the current status. If the
	 * current status is PASSED, the building block shall also return the
	 * certificate chain as well as the calculated validation time returned
	 * in step 3.
	 */
}

// prospectiveCertificateChain ports the private prospectiveCertificateChain().
func (c *PastCertificateValidation) prospectiveCertificateChain() process.ChainItem[*jaxb.XmlPCV] {
	constraint := c.policy.ProspectiveCertificateChainConstraint(c.context)
	return NewProspectiveCertificateChainCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// vtsResult ports the private getVTSResult(CertificateWrapper).
func (c *PastCertificateValidation) vtsResult(trustedCertificate *diagnostic.CertificateWrapper) *jaxb.XmlVTS {
	validationTimeSliding := NewValidationTimeSliding(c.I18nProvider, c.token, trustedCertificate, c.currentTime,
		c.poe, c.bbbs, c.context, c.policy)
	return validationTimeSliding.Execute()
}

// validationTimeSliding ports the private
// validationTimeSliding(XmlVTS, CertificateWrapper).
func (c *PastCertificateValidation) validationTimeSliding(vts *jaxb.XmlVTS,
	trustedCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlPCV] {
	return NewValidationTimeSlidingCheck(c.I18nProvider, c.Result, vts, c.token.Id(), trustedCertificate,
		c.WarnLevelRule())
}

// successfulValidationTimeSlidingFound ports the private
// successfulValidationTimeSlidingFound(XmlVTS).
func (c *PastCertificateValidation) successfulValidationTimeSlidingFound(vts *jaxb.XmlVTS) process.ChainItem[*jaxb.XmlPCV] {
	return NewSuccessfulValidationTimeSlidingFoundCheck(c.I18nProvider, c.Result, vts, c.FailLevelRule())
}

// cryptographicCheck ports the private cryptographicCheck(XmlPCV, Date): the
// Java parameter is this chain's own result object, which the Go form reaches
// through c.Result. The unsupported-Context error becomes a panic, the caller
// being initChain, which cannot propagate one.
func (c *PastCertificateValidation) cryptographicCheck(validationTime time.Time) process.ChainItem[*jaxb.XmlPCV] {
	algorithmObsolescenceValidation := aov.NewTokenCertificateChainAlgorithmObsolescenceValidation(
		c.I18nProvider, c.token, c.context, validationTime, c.policy)
	aovResult := algorithmObsolescenceValidation.Execute()

	position, err := process.GetCertificateChainCryptoPosition(c.context)
	if err != nil {
		panic(err)
	}
	return aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, aovResult, validationTime,
		position, c.token.Id())
}

// bestValidationTimeSliding returns a successful VTS result with the latest
// control time (enough to proof validity of the signature), when applicable.
// Port of the private getBestValidationTimeSliding(List).
//
// Java dereferences getControlTime() unguarded in the comparisons, so a VTS
// result carrying none raises a NullPointerException; the nil *XSDateTime
// panics here in its place.
func (c *PastCertificateValidation) bestValidationTimeSliding(vtsList []*jaxb.XmlVTS) *jaxb.XmlVTS {
	if utils.IsCollectionEmpty(vtsList) {
		return nil
	}
	var bestVTS *jaxb.XmlVTS
	for _, xmlVTS := range vtsList {
		if bestVTS == nil ||
			(!c.IsValid(&bestVTS.XmlConstraintsConclusionContent) &&
				(c.IsValid(&xmlVTS.XmlConstraintsConclusionContent) ||
					bestVTS.ControlTime.Time().Before(xmlVTS.ControlTime.Time())) ||
				(c.IsValid(&xmlVTS.XmlConstraintsConclusionContent) &&
					bestVTS.ControlTime.Time().Before(xmlVTS.ControlTime.Time()))) {
			bestVTS = xmlVTS
		}
	}
	return bestVTS
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// addAdditionalInfo(): the control time can be null.
func (c *PastCertificateValidation) AddAdditionalInfo() {
	if c.controlTime != nil {
		c.Result.Value.ControlTime = jaxb.NewXSDateTime(*c.controlTime)
	} else {
		c.Result.Value.ControlTime = nil
	}
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of the overridden
// collectMessages(XmlConclusion, XmlConstraint): the generated BlockType member
// is a *XmlBlockType, whose nil is Java's null.
func (c *PastCertificateValidation) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	if constraint.BlockType != nil && jaxb.XmlBlockTypeVTS == *constraint.BlockType {
		// skip validation for VTS
	} else {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}
