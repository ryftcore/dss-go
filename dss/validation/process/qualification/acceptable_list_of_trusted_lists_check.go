// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/AcceptableListOfTrustedListsCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// AcceptableListOfTrustedListsCheck verifies whether the validation of a
// List of Trusted Lists is conclusive.
type AcceptableListOfTrustedListsCheck[T any] struct {
	*AbstractTrustedListCheck[T]
}

// NewAcceptableListOfTrustedListsCheck is the default constructor. Port of
// AcceptableListOfTrustedListsCheck(I18nProvider, T, XmlTLAnalysis, LevelRule).
func NewAcceptableListOfTrustedListsCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	lotlAnalysis *jaxb.XmlTLAnalysis, constraint policy.LevelRule) *AcceptableListOfTrustedListsCheck[T] {
	c := &AcceptableListOfTrustedListsCheck[T]{
		AbstractTrustedListCheck: NewAbstractTrustedListCheck(i18nProvider, result, lotlAnalysis, constraint),
	}
	c.InitChainItem(c)
	return c
}

// MessageTag returns the check's message tag. Port of the overridden getMessageTag().
func (c *AcceptableListOfTrustedListsCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT
}

// ErrorMessageTag returns the check's error message tag. Port of the
// overridden getErrorMessageTag().
func (c *AcceptableListOfTrustedListsCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_LIST_OF_TRUSTED_LISTS_ACCEPT_ANS
}
