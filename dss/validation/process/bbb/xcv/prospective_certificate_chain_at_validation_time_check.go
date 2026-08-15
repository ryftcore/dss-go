// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/checks/ProspectiveCertificateChainAtValidationTimeCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// ProspectiveCertificateChainAtValidationTimeCheck verifies whether a
// prospective certificate chain with trust anchors valid at validation time
// has been found.
type ProspectiveCertificateChainAtValidationTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// controlTime is the validation time to check against.
	controlTime time.Time
}

// NewProspectiveCertificateChainAtValidationTimeCheck is the default
// constructor. Port of
// ProspectiveCertificateChainAtValidationTimeCheck(I18nProvider, XmlXCV, CertificateWrapper, Date, LevelRule).
func NewProspectiveCertificateChainAtValidationTimeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlXCV],
	certificate *diagnostic.CertificateWrapper, controlTime time.Time, constraint policy.LevelRule) *ProspectiveCertificateChainAtValidationTimeCheck {
	c := &ProspectiveCertificateChainAtValidationTimeCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
		controlTime:   controlTime,
	}
	c.InitChainItem(c)
	return c
}

// failLevelRule ports the private getFailLevelRule().
func (c *ProspectiveCertificateChainAtValidationTimeCheck) failLevelRule() policy.LevelRule {
	return process.GetLevelRule(enumerations.Level_FAIL)
}

// Process performs the check. Port of process().
func (c *ProspectiveCertificateChainAtValidationTimeCheck) Process() bool {
	// FAIL level constraint is used to fail the check
	if process.IsTrustAnchor(c.certificate, c.controlTime, c.failLevelRule()) {
		return true
	}
	for _, caCertificate := range c.certificate.CertificateChain() {
		if process.IsTrustAnchor(caCertificate, c.controlTime, c.failLevelRule()) {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ProspectiveCertificateChainAtValidationTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_HPCCVVT
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *ProspectiveCertificateChainAtValidationTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_HPCCVVT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ProspectiveCertificateChainAtValidationTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *ProspectiveCertificateChainAtValidationTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_NO_CERTIFICATE_CHAIN_FOUND_NO_POE
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *ProspectiveCertificateChainAtValidationTimeCheck) BuildAdditionalInfo() *string {
	controlTime := c.controlTime
	message := c.I18nProvider.GetMessage(i18n.MessageTag_VALIDATION_TIME, process.GetFormattedDate(&controlTime))
	return &message
}
