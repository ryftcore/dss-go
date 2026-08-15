// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/AcceptableLoLoTECheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// AcceptableLoLoTECheck verifies if the validation of a LoLoTE was
// successful.
type AcceptableLoLoTECheck[T any] struct {
	*AcceptableLoTECheck[T]
}

// NewAcceptableLoLoTECheck is the default constructor. Port of
// AcceptableLoLoTECheck(I18nProvider, T, XmlLoTEAnalysis, LevelRule).
func NewAcceptableLoLoTECheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	loloteAnalysis *jaxb.XmlLoTEAnalysis, constraint policy.LevelRule) *AcceptableLoLoTECheck[T] {
	c := &AcceptableLoLoTECheck[T]{
		AcceptableLoTECheck: NewAcceptableLoTECheck(i18nProvider, result, loloteAnalysis, constraint),
	}
	c.InitChainItem(c)
	return c
}

// MessageTag returns the check's message tag. Port of the overridden getMessageTag().
func (c *AcceptableLoLoTECheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_LOLOTE_ACCEPT
}

// ErrorMessageTag returns the check's error message tag. Port of the
// overridden getErrorMessageTag().
func (c *AcceptableLoLoTECheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_LOLOTE_ACCEPT_ANS
}
