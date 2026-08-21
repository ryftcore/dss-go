// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/converter/InternationalNamesTypeConverter.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// stringPredicate is the Go form of java.util.function.Predicate<String>, used by
// InternationalNamesTypeConverter and NonEmptyMultiLangURIListTypeConverter.
type stringPredicate interface {
	Test(t string) bool
}

// alwaysTrueStringPredicate is the Go form of the default constructors' `x -> true` lambda.
type alwaysTrueStringPredicate struct{}

// Test always answers true.
func (alwaysTrueStringPredicate) Test(t string) bool { return true }

// InternationalNamesTypeConverter extracts language based values.
type InternationalNamesTypeConverter struct {
	// predicate is the predicate to be used.
	predicate stringPredicate
}

// NewInternationalNamesTypeConverter is the default constructor (selects all). Port of
// InternationalNamesTypeConverter().
func NewInternationalNamesTypeConverter() *InternationalNamesTypeConverter {
	return NewInternationalNamesTypeConverterWithPredicate(alwaysTrueStringPredicate{})
}

// NewInternationalNamesTypeConverterWithPredicate is the constructor with a filter predicate.
// Port of InternationalNamesTypeConverter(Predicate).
func NewInternationalNamesTypeConverterWithPredicate(predicate stringPredicate) *InternationalNamesTypeConverter {
	return &InternationalNamesTypeConverter{predicate: predicate}
}

// Apply ports apply(InternationalNamesType).
func (c *InternationalNamesTypeConverter) Apply(original *jaxb.InternationalNamesType) map[string][]string {
	result := make(map[string][]string)
	if original != nil && utils.IsCollectionNotEmpty(original.Name) {
		for _, multiLangNormString := range original.Name {
			lang := multiLangNormString.Lang
			value := multiLangNormString.Value
			if c.predicate.Test(value) {
				result[lang] = append(result[lang], value)
			}
		}
	}
	return result
}
