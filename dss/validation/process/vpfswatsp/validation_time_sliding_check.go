// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/pcv/checks/ValidationTimeSlidingCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ValidationTimeSlidingCheck checks if the Validation Time Sliding result is
// valid.
type ValidationTimeSlidingCheck struct {
	*process.ChainItemBase[*jaxb.XmlPCV]

	// vts is the Validation Time Sliding result.
	vts *jaxb.XmlVTS

	// trustedCertificate is the certificate used as a trust anchor during the
	// VTS process.
	trustedCertificate *diagnostic.CertificateWrapper
}

// NewValidationTimeSlidingCheck is the default constructor. Port of
// ValidationTimeSlidingCheck(I18nProvider, XmlPCV, XmlVTS, String, CertificateWrapper, LevelRule).
func NewValidationTimeSlidingCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlPCV],
	vts *jaxb.XmlVTS, tokenId string, trustedCertificate *diagnostic.CertificateWrapper,
	constraint policy.LevelRule) *ValidationTimeSlidingCheck {
	c := &ValidationTimeSlidingCheck{
		ChainItemBase:      process.NewChainItemBaseWithId(i18nProvider, result, constraint, tokenId),
		vts:                vts,
		trustedCertificate: trustedCertificate,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *ValidationTimeSlidingCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeVTS
}

// Process performs the check. Port of process().
func (c *ValidationTimeSlidingCheck) Process() bool {
	return c.IsValid(&c.vts.XmlConstraintsConclusionContent)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ValidationTimeSlidingCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PCV_IVTSC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ValidationTimeSlidingCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PCV_IVTSC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ValidationTimeSlidingCheck) FailedIndicationForConclusion() enumerations.Indication {
	return c.vts.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(): the generated SubIndication
// member is a pointer, whose nil is Java's null.
func (c *ValidationTimeSlidingCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.vts.Conclusion.SubIndication == nil {
		return ""
	}
	return c.vts.Conclusion.SubIndication.SubIndication()
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(), whose null result leaves the element absent.
func (c *ValidationTimeSlidingCheck) BuildAdditionalInfo() *string {
	// XmlVTS#getControlTime(): the generated member is a *XSDateTime, whose nil
	// is Java's null.
	if c.vts.ControlTime != nil {
		controlTime := c.vts.ControlTime.Time()
		var message string
		if c.trustedCertificate != nil {
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CONTROL_TIME_WITH_TRUST_ANCHOR,
				c.trustedCertificate.Id(), process.GetFormattedDate(&controlTime))
		} else {
			message = c.I18nProvider.GetMessage(i18n.MessageTag_CONTROL_TIME_ALONE,
				process.GetFormattedDate(&controlTime))
		}
		return &message
	}
	return nil
}
