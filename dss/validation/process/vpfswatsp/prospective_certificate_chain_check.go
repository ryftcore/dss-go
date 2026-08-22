// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/pcv/checks/ProspectiveCertificateChainCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note. The homonymous class in
// bbb.xcv.checks is a different one and stays in the Go package xcv.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ProspectiveCertificateChainCheck checks if the certificate's chain is trusted.
type ProspectiveCertificateChainCheck struct {
	*process.ChainItemBase[*jaxb.XmlPCV]

	// token is the token's chain to check.
	token diagnostic.TokenProxy
}

// NewProspectiveCertificateChainCheck is the default constructor. Port of
// ProspectiveCertificateChainCheck(Provider, XmlPCV, TokenProxy, LevelRule).
func NewProspectiveCertificateChainCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlPCV],
	token diagnostic.TokenProxy, constraint policy.LevelRule) *ProspectiveCertificateChainCheck {
	c := &ProspectiveCertificateChainCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ProspectiveCertificateChainCheck) Process() bool {
	return c.token.IsTrustedChain()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ProspectiveCertificateChainCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCCCBB
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ProspectiveCertificateChainCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCCCBBANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ProspectiveCertificateChainCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *ProspectiveCertificateChainCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoCertificateChainFound
}
