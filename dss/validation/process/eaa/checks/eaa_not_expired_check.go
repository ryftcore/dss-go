// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAANotExpiredCheck.java (DSS 6.5.RC1).
package checks

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAANotExpiredCheck verifies whether the validation time is within the EAA
// technical validity period range.
type EAANotExpiredCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper

	// validationTime is the EAA validation time.
	validationTime time.Time
}

// NewEAANotExpiredCheck is the default constructor.
func NewEAANotExpiredCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, validationTime time.Time, constraint policy.LevelRule) *EAANotExpiredCheck {
	c := &EAANotExpiredCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:            eaaWrapper,
		validationTime: validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAANotExpiredCheck) Process() bool {
	return c.notBefore() && c.notAtOrAfter()
}

// notBefore ports the private notBefore().
//
// The "nbf" (not before) claim identifies the time before which the JWT MUST
// NOT be accepted for processing. The processing of the "nbf" claim requires
// that the current date/time MUST be after or equal to the not-before
// date/time listed in the "nbf" claim.
func (c *EAANotExpiredCheck) notBefore() bool {
	return c.eaa.EAANotBefore() != nil && !c.validationTime.Before(*c.eaa.EAANotBefore())
}

// notAtOrAfter ports the private notAtOrAfter().
//
// The "exp" (expiration time) claim identifies the expiration time on or
// after which the JWT MUST NOT be accepted for processing. The processing of
// the "exp" claim requires that the current date/time MUST be before the
// expiration date/time listed in the "exp" claim.
func (c *EAANotExpiredCheck) notAtOrAfter() bool {
	return c.eaa.EAAExpiration() != nil && c.validationTime.Before(*c.eaa.EAAExpiration())
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *EAANotExpiredCheck) BuildAdditionalInfo() *string {
	if !c.notBefore() || !c.notAtOrAfter() {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_EAA_VT_ITVR_VALIDITY,
			process.GetFormattedDate(&c.validationTime), process.GetFormattedDate(c.eaa.EAANotBefore()),
			process.GetFormattedDate(c.eaa.EAAExpiration()))
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAANotExpiredCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_VT_ITVR
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAANotExpiredCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_VT_ITVR_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAANotExpiredCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAANotExpiredCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_OUT_OF_BOUNDS_NO_POE
}
