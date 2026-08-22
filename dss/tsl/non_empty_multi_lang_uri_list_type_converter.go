// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/converter/NonEmptyMultiLangURIListTypeConverter.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// NonEmptyMultiLangURIListTypeConverter extracts non-empty URI language based values.
type NonEmptyMultiLangURIListTypeConverter struct {
	// predicate is the predicate to be used.
	predicate stringPredicate
}

// NewNonEmptyMultiLangURIListTypeConverter is the default constructor (selects all). Port of
// NonEmptyMultiLangURIListTypeConverter().
func NewNonEmptyMultiLangURIListTypeConverter() *NonEmptyMultiLangURIListTypeConverter {
	return NewNonEmptyMultiLangURIListTypeConverterWithPredicate(alwaysTrueStringPredicate{})
}

// NewNonEmptyMultiLangURIListTypeConverterWithPredicate is the constructor with a filter
// predicate. Port of NonEmptyMultiLangURIListTypeConverter(Predicate).
func NewNonEmptyMultiLangURIListTypeConverterWithPredicate(predicate stringPredicate) *NonEmptyMultiLangURIListTypeConverter {
	return &NonEmptyMultiLangURIListTypeConverter{predicate: predicate}
}

// Apply ports apply(NonEmptyMultiLangURIListType).
func (c *NonEmptyMultiLangURIListTypeConverter) Apply(original *jaxb.NonEmptyMultiLangURIListType) map[string][]string {
	result := make(map[string][]string)
	if original != nil && utils.IsCollectionNotEmpty(original.URI) {
		for _, multiLangURIString := range original.URI {
			lang := multiLangURIString.Lang
			value := multiLangURIString.Value
			if c.predicate.Test(value) {
				result[lang] = append(result[lang], value)
			}
		}
	}
	return result
}
