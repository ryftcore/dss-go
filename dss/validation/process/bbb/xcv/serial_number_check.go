// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/SerialNumberCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// SerialNumberCheck checks if the certificate's serial number is present.
type SerialNumberCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewSerialNumberCheck is the default constructor. Port of
// SerialNumberCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewSerialNumberCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *SerialNumberCheck {
	c := &SerialNumberCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SerialNumberCheck) Process() bool {
	return utils.IsStringNotBlank(c.certificate.SerialNumber())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SerialNumberCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_UNIQUE_CERT
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *SerialNumberCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_UNIQUE_CERT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SerialNumberCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *SerialNumberCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}
