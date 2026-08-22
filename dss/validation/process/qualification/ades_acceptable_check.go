// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/AdESAcceptableCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AdESAcceptableCheck checks whether AdES signature validation as per EN
// 319 102-1 succeeded.
type AdESAcceptableCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationSignatureQualification]

	// etsi319102Conclusion is the final conclusion of EN 319 102-1 AdES
	// signature validation.
	etsi319102Conclusion *jaxb.XmlConclusion

	// error is the internal cached error message.
	error i18n.MessageTag
}

// NewAdESAcceptableCheck is the default constructor. Port of
// AdESAcceptableCheck(I18nProvider, XmlValidationSignatureQualification, XmlConclusion, LevelRule).
func NewAdESAcceptableCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationSignatureQualification],
	etsi319102Conclusion *jaxb.XmlConclusion, constraint policy.LevelRule) *AdESAcceptableCheck {
	c := &AdESAcceptableCheck{
		ChainItemBase:        process.NewChainItemBase(i18nProvider, result, constraint),
		etsi319102Conclusion: etsi319102Conclusion,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AdESAcceptableCheck) Process() bool {
	valid := c.IsValidConclusion(c.etsi319102Conclusion)
	if !valid {
		if c.IsIndeterminateConclusion(c.etsi319102Conclusion) {
			c.error = i18n.MessageTagQualIsAdESInd
		} else if c.IsInvalidConclusion(c.etsi319102Conclusion) {
			c.error = i18n.MessageTagQualIsAdESINV
		}
		return false
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AdESAcceptableCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualIsAdES
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *AdESAcceptableCheck) ErrorMessageTag() i18n.MessageTag {
	return c.error
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AdESAcceptableCheck) FailedIndicationForConclusion() enumerations.Indication {
	return c.etsi319102Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AdESAcceptableCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.etsi319102Conclusion.SubIndication == nil {
		return ""
	}
	return c.etsi319102Conclusion.SubIndication.SubIndication()
}
