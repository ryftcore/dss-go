// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/AcceptableTrustedListCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AcceptableTrustedListCheck verifies whether the validation of a Trusted
// List is conclusive.
type AcceptableTrustedListCheck[T any] struct {
	*AbstractTrustedListCheck[T]
}

// NewAcceptableTrustedListCheck is the default constructor. Port of
// AcceptableTrustedListCheck(Provider, T, XmlTLAnalysis, LevelRule).
func NewAcceptableTrustedListCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	tlAnalysis *jaxb.XmlTLAnalysis, constraint policy.LevelRule) *AcceptableTrustedListCheck[T] {
	c := &AcceptableTrustedListCheck[T]{
		AbstractTrustedListCheck: NewAbstractTrustedListCheck(i18nProvider, result, tlAnalysis, constraint),
	}
	c.InitChainItem(c)
	return c
}

// MessageTag returns the check's message tag. Port of the overridden getMessageTag().
func (c *AcceptableTrustedListCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualTrustedListAccept
}

// ErrorMessageTag returns the check's error message tag. Port of the
// overridden getErrorMessageTag().
func (c *AcceptableTrustedListCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualTrustedListAcceptANS
}
