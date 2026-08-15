// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/OtherTrustAnchorExistsCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// OtherTrustAnchorExistsCheck verifies if other trust anchor exists in the
// certificate chain.
type OtherTrustAnchorExistsCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewOtherTrustAnchorExistsCheck is the default constructor. Port of
// OtherTrustAnchorExistsCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewOtherTrustAnchorExistsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *OtherTrustAnchorExistsCheck {
	c := &OtherTrustAnchorExistsCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *OtherTrustAnchorExistsCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockType_SUB_XCV_TA
}

// Process performs the check. Port of process().
func (c *OtherTrustAnchorExistsCheck) Process() bool {
	return c.certificate.IsTrustedChain()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *OtherTrustAnchorExistsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IOTAA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *OtherTrustAnchorExistsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IOTAA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *OtherTrustAnchorExistsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *OtherTrustAnchorExistsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_NO_CERTIFICATE_CHAIN_FOUND_NO_POE
}
