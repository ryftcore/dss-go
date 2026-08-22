// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/AcceptableTrustedListPresenceCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AcceptableTrustedListPresenceCheck verifies whether acceptable Trusted
// Lists have been found. Java's Set<String> validTLUrls -> map[string]struct{}
// per PORTING.md's collections mapping.
type AcceptableTrustedListPresenceCheck[T any] struct {
	*process.ChainItemBase[T]

	// validTLUrls is the set of URLs of acceptable Trusted Lists.
	validTLUrls map[string]struct{}
}

// NewAcceptableTrustedListPresenceCheck is the default constructor. Port of
// AcceptableTrustedListPresenceCheck(I18nProvider, T, Set, LevelRule).
func NewAcceptableTrustedListPresenceCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	validTLUrls map[string]struct{}, constraint policy.LevelRule) *AcceptableTrustedListPresenceCheck[T] {
	c := &AcceptableTrustedListPresenceCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		validTLUrls:   validTLUrls,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AcceptableTrustedListPresenceCheck[T]) Process() bool {
	return utils.IsMapNotEmpty(c.validTLUrls)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AcceptableTrustedListPresenceCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_VALID_TRUSTED_LIST_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *AcceptableTrustedListPresenceCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_VALID_TRUSTED_LIST_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AcceptableTrustedListPresenceCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *AcceptableTrustedListPresenceCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
