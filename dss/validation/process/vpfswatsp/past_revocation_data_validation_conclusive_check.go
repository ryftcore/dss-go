// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/checks/PastRevocationDataValidationConclusiveCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// PastRevocationDataValidationConclusiveCheck checks if an acceptable
// revocation data is found.
type PastRevocationDataValidationConclusiveCheck struct {
	*process.ChainItemBase[*jaxb.XmlPSV]

	// conclusion is the validation conclusion.
	conclusion *jaxb.XmlConclusion
}

// NewPastRevocationDataValidationConclusiveCheck is the constructor. Port of
// PastRevocationDataValidationConclusiveCheck(I18nProvider, XmlPSV, XmlConclusion, LevelRule).
func NewPastRevocationDataValidationConclusiveCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlPSV], conclusion *jaxb.XmlConclusion,
	constraint policy.LevelRule) *PastRevocationDataValidationConclusiveCheck {
	c := &PastRevocationDataValidationConclusiveCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		conclusion:    conclusion,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *PastRevocationDataValidationConclusiveCheck) Process() bool {
	return c.IsValidConclusion(c.conclusion)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PastRevocationDataValidationConclusiveCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_DIURDSCHPVR
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *PastRevocationDataValidationConclusiveCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_DIURDSCHPVR_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *PastRevocationDataValidationConclusiveCheck) FailedIndicationForConclusion() enumerations.Indication {
	return c.conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(): the generated SubIndication
// member is a pointer, whose nil is Java's null.
func (c *PastRevocationDataValidationConclusiveCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.conclusion.SubIndication == nil {
		return ""
	}
	return c.conclusion.SubIndication.SubIndication()
}
