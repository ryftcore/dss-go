// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/SigningTimeCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SigningTimeCheck checks if a claimed signing time is present.
type SigningTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewSigningTimeCheck is the default constructor. Port of
// SigningTimeCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewSigningTimeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *SigningTimeCheck {
	c := &SigningTimeCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SigningTimeCheck) Process() bool {
	return c.signature.ClaimedSigningTime() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SigningTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISQPSTP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SigningTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISQPSTP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SigningTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SigningTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
