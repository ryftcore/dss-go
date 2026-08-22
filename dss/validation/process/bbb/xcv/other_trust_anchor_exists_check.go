// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/OtherTrustAnchorExistsCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// OtherTrustAnchorExistsCheck verifies if other trust anchor exists in the
// certificate chain.
type OtherTrustAnchorExistsCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewOtherTrustAnchorExistsCheck is the default constructor. Port of
// OtherTrustAnchorExistsCheck(Provider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewOtherTrustAnchorExistsCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
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
	return jaxb.XmlBlockTypeSubXCVTA
}

// Process performs the check. Port of process().
func (c *OtherTrustAnchorExistsCheck) Process() bool {
	return c.certificate.IsTrustedChain()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *OtherTrustAnchorExistsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIOTAA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *OtherTrustAnchorExistsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIOTAAANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *OtherTrustAnchorExistsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *OtherTrustAnchorExistsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoCertificateChainFoundNoPOE
}
