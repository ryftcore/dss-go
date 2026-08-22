// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/status/EAARevocationNotExpiredCheck.java (DSS 6.5.RC1).
package checks

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAARevocationNotExpiredCheck verifies whether the EAA revocation is not yet
// expired.
type EAARevocationNotExpiredCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper

	// validationTime is the validation time.
	validationTime time.Time
}

// NewEAARevocationNotExpiredCheck is the default constructor.
func NewEAARevocationNotExpiredCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper, validationTime time.Time,
	constraint policy.LevelRule) *EAARevocationNotExpiredCheck {
	c := &EAARevocationNotExpiredCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaaStatusToken: eaaStatusToken,
		validationTime: validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
//
// The "exp" (expiration time) claim identifies the expiration time on or
// after which the JWT MUST NOT be accepted for processing.
func (c *EAARevocationNotExpiredCheck) Process() bool {
	return c.eaaStatusToken.ExpirationTime() != nil && c.validationTime.Before(*c.eaaStatusToken.ExpirationTime())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAARevocationNotExpiredCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevNotExp
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationNotExpiredCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAARevNotExpANS
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *EAARevocationNotExpiredCheck) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTagEAARevTime, process.GetFormattedDate(&c.validationTime),
		process.GetFormattedDate(c.eaaStatusToken.IssuedAt()), process.GetFormattedDate(c.eaaStatusToken.ExpirationTime()))
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationNotExpiredCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationNotExpiredCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
