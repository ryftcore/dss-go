// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/UninterruptedChainItem.java (DSS 6.5.RC1).
package process

import (
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
)

// UninterruptedChainItemBase allows to continue the chain validation process in
// case of a check failure. It is the Go form of the abstract class
// UninterruptedChainItem<T>, which adds nothing but the overridden
// continueProcessOnFail(); a concrete check embeds it instead of ChainItemBase
// and registers itself with InitChainItem exactly the same way.
type UninterruptedChainItemBase[T any] struct {
	*ChainItemBase[T]
}

// NewUninterruptedChainItemBase is the default constructor. Port of
// UninterruptedChainItem(I18nProvider, T, LevelRule).
func NewUninterruptedChainItemBase[T any](i18nProvider *i18n.I18nProvider, result *Result[T],
	constraint policy.LevelRule) *UninterruptedChainItemBase[T] {
	return &UninterruptedChainItemBase[T]{
		ChainItemBase: NewChainItemBase(i18nProvider, result, constraint),
	}
}

// NewUninterruptedChainItemBaseWithId is the constructor with custom Id. Port of
// UninterruptedChainItem(I18nProvider, T, LevelRule, String).
func NewUninterruptedChainItemBaseWithId[T any](i18nProvider *i18n.I18nProvider, result *Result[T],
	constraint policy.LevelRule, id string) *UninterruptedChainItemBase[T] {
	return &UninterruptedChainItemBase[T]{
		ChainItemBase: NewChainItemBaseWithId(i18nProvider, result, constraint, id),
	}
}

// ContinueProcessOnFail returns TRUE: the validation process is continued on a
// check failure. Port of the overridden continueProcessOnFail().
func (c *UninterruptedChainItemBase[T]) ContinueProcessOnFail() bool {
	return true
}
