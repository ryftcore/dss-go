// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/AbstractValueCheckItem.java (DSS 6.5.RC1).
package bbb

import (
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// allValue accepts all values.
const allValue = "*"

// AbstractValueCheckItem checks if the value is allowed. A concrete check embeds
// it instead of process.ChainItemBase and registers itself with InitChainItem
// the same way.
type AbstractValueCheckItem[T any] struct {
	*process.ChainItemBase[T]

	// constraint is the value constraint.
	constraint policy.ValueRule
}

// NewAbstractValueCheckItem is the default constructor. Port of
// AbstractValueCheckItem(I18nProvider, T, ValueRule).
func NewAbstractValueCheckItem[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	constraint policy.ValueRule) *AbstractValueCheckItem[T] {
	return &AbstractValueCheckItem[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		constraint:    constraint,
	}
}

// ProcessValueCheck processes the value check, returning TRUE if the value
// matches the expected one. Port of processValueCheck(String).
func (c *AbstractValueCheckItem[T]) ProcessValueCheck(value string) bool {
	if utils.IsStringEmpty(value) {
		return false
	}
	expected := c.constraint.Value()
	if allValue == expected {
		return true
	} else {
		return utils.AreStringsEqual(expected, value)
	}
}
