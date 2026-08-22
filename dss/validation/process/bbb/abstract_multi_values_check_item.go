// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/AbstractMultiValuesCheckItem.java (DSS 6.5.RC1).
package bbb

import (
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractMultiValuesCheckItem is the abstract class to check if the given value
// is one of the allowed values by ValidationPolicy. A concrete check embeds it
// instead of process.ChainItemBase and registers itself with InitChainItem the
// same way.
type AbstractMultiValuesCheckItem[T any] struct {
	*process.ChainItemBase[T]

	// constraint is the constraint value.
	constraint policy.MultiValuesRule
}

// NewAbstractMultiValuesCheckItem is the default constructor. Port of
// AbstractMultiValuesCheckItem(I18nProvider, T, MultiValuesRule).
func NewAbstractMultiValuesCheckItem[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	constraint policy.MultiValuesRule) *AbstractMultiValuesCheckItem[T] {
	return &AbstractMultiValuesCheckItem[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		constraint:    constraint,
	}
}

// ProcessValueCheck checks the value, returning TRUE if the value is allowed by
// the constraint. Port of processValueCheck(String).
func (c *AbstractMultiValuesCheckItem[T]) ProcessValueCheck(value string) bool {
	return process.ProcessValueCheck(value, c.constraint.Values())
}

// ProcessValuesCheck checks the values, returning TRUE if the values are allowed
// by the constraint. Port of processValuesCheck(List).
func (c *AbstractMultiValuesCheckItem[T]) ProcessValuesCheck(values []string) bool {
	return process.ProcessValuesCheck(values, c.constraint.Values())
}

// ProcessAllValuesCheck checks the values, returning TRUE if all the values are
// allowed by the constraint. Port of processAllValuesCheck(List).
func (c *AbstractMultiValuesCheckItem[T]) ProcessAllValuesCheck(values []string) bool {
	return process.ProcessAllValuesCheck(values, c.constraint.Values())
}

// ProcessValuesForEachExpectedCheck checks whether values contain all the
// expected values specified in the policy constraint. Port of
// processValuesForEachExpectedCheck(List).
func (c *AbstractMultiValuesCheckItem[T]) ProcessValuesForEachExpectedCheck(values []string) bool {
	return process.ProcessValuesForEachExpectedCheck(values, c.constraint.Values())
}

// Values gets a list of expected values as specified within the policy
// constraint. Port of getValues().
func (c *AbstractMultiValuesCheckItem[T]) Values() []string {
	return c.constraint.Values()
}
