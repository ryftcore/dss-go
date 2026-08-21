// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/OtherTSLPointerPredicate.java (DSS 6.5.RC1).
package tsl

import "github.com/ryftcore/dss-go/dss/trustedlist/jaxb"

// OtherTSLPointerPredicate is a predicate allowing to filter TSL pointers. Port of
// java.util.function.Predicate<OtherTSLPointerType> (the marker sub-interface carries no members
// of its own).
type OtherTSLPointerPredicate interface {
	// Test tests the predicate on the given OtherTSLPointerType. Port of test(OtherTSLPointerType).
	Test(o *jaxb.OtherTSLPointerType) bool
}

// otherTSLPointerPredicateAnd is the Go form of Predicate#and(Predicate), used by
// TLPredicateFactory to compose OtherTSLPointerPredicates.
type otherTSLPointerPredicateAnd struct {
	first, second OtherTSLPointerPredicate
}

// Test evaluates both predicates, short-circuiting like Java's default and().
func (a *otherTSLPointerPredicateAnd) Test(o *jaxb.OtherTSLPointerType) bool {
	return a.first.Test(o) && a.second.Test(o)
}

// otherTSLPointerPredicateAndOf builds the Go equivalent of `first.and(second)`.
func otherTSLPointerPredicateAndOf(first, second OtherTSLPointerPredicate) OtherTSLPointerPredicate {
	return &otherTSLPointerPredicateAnd{first: first, second: second}
}
