// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/vts/checks/SunsetDateCheck.java (DSS 6.5.RC1).
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

// SunsetDateCheck checks if the sunset date is defined for the current trust
// anchor.
type SunsetDateCheck struct {
	*process.ChainItemBase[*jaxb.XmlVTS]

	// trustedCertificate is the trust anchor to check the sunset date.
	trustedCertificate *diagnostic.CertificateWrapper
}

// NewSunsetDateCheck is the default constructor. Port of
// SunsetDateCheck(I18nProvider, XmlVTS, CertificateWrapper, LevelRule).
func NewSunsetDateCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlVTS],
	trustedCertificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *SunsetDateCheck {
	c := &SunsetDateCheck{
		ChainItemBase:      process.NewChainItemBaseWithId(i18nProvider, result, constraint, trustedCertificate.Id()),
		trustedCertificate: trustedCertificate,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *SunsetDateCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeSubXCVTA
}

// Process performs the check. Port of process().
func (c *SunsetDateCheck) Process() bool {
	return c.trustedCertificate != nil && c.trustedCertificate.TrustSunsetDate() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SunsetDateCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_ISDDTA
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SunsetDateCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_ISDDTA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion(), which returns null.
func (c *SunsetDateCheck) FailedIndicationForConclusion() enumerations.Indication {
	return ""
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), which returns null.
func (c *SunsetDateCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(), whose null result leaves the element absent.
func (c *SunsetDateCheck) BuildAdditionalInfo() *string {
	if c.trustedCertificate != nil && c.trustedCertificate.TrustSunsetDate() != nil {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_CERTIFICATE_SUNSET_DATE_TRUST_ANCHOR,
			c.trustedCertificate.Id(), process.GetFormattedDate(c.trustedCertificate.TrustSunsetDate()))
		return &message
	}
	return nil
}
