// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/dto/condition/CompositeCondition.java (DSS 6.5.RC1).
package tsl

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
)

// CompositeCondition is the condition resulting of the matchingCriteriaIndicator of other
// Conditions. Upstream's class name recalls its XML origin: its toString() still prints
// "CriteriaListCondition".
//
// java.io.Serializable has no Go counterpart and is dropped.
type CompositeCondition struct {
	// matchingCriteriaIndicator is the matching criteria: atLeastOne, all, none.
	matchingCriteriaIndicator enumerations.Assert

	// children is the list of child conditions.
	children []tslmodel.Condition
}

// NewCompositeCondition is the default constructor for CriteriaListCondition: all conditions
// must match. Port of CompositeCondition().
//
// Java initialises the children field to an empty ArrayList at declaration; Go's nil slice
// behaves identically for append and len, but getChildren() must not hand back nil where Java
// hands back an empty list, so Children() normalises (see below).
func NewCompositeCondition() *CompositeCondition {
	return &CompositeCondition{matchingCriteriaIndicator: enumerations.AssertAll}
}

// NewCompositeConditionWithMatchingCriteriaIndicator is the constructor for
// CriteriaListCondition with an explicit matching criteria indicator (atLeastOne, all, none).
// Port of CompositeCondition(Assert).
func NewCompositeConditionWithMatchingCriteriaIndicator(matchingCriteriaIndicator enumerations.Assert) *CompositeCondition {
	return &CompositeCondition{matchingCriteriaIndicator: matchingCriteriaIndicator}
}

// Children returns the list of child conditions: possibly empty, never nil. Port of the final
// getChildren(), which answers an unmodifiable view of a field Java never leaves null.
func (c *CompositeCondition) Children() []tslmodel.Condition {
	if c.children == nil {
		return []tslmodel.Condition{}
	}
	return c.children
}

// AddChild adds a child condition, allowing embedded conditions to be handled. Port of
// addChild(Condition).
func (c *CompositeCondition) AddChild(condition tslmodel.Condition) {
	c.children = append(c.children, condition)
}

// MatchingCriteriaIndicator returns the matching criteria indicator: atLeastOne, all, none.
// Port of getMatchingCriteriaIndicator().
func (c *CompositeCondition) MatchingCriteriaIndicator() enumerations.Assert {
	return c.matchingCriteriaIndicator
}

// Check executes the composite condition on the given certificate, returning true if the
// condition matches. Port of check(CertificateToken).
//
// Java's default branch throws DSSException("Unsupported MatchingCriteriaIndicator : " + ...)
// for an indicator outside the enum. Condition#check declares no checked exception and the Go
// interface returns a bare bool, so an unsupported indicator answers false here - the
// conservative reading (an unmatched condition never grants trust). Java can only reach that
// branch through a null indicator, which the two constructors above make unrepresentable
// except by explicitly passing the empty Assert.
func (c *CompositeCondition) Check(certificateToken *model.CertificateToken) bool {
	switch c.matchingCriteriaIndicator {
	case enumerations.AssertAll:
		for _, condition := range c.children {
			if !condition.Check(certificateToken) {
				return false
			}
		}
		return true
	case enumerations.AssertAtLeastOne:
		for _, condition := range c.children {
			if condition.Check(certificateToken) {
				return true
			}
		}
		return false
	case enumerations.AssertNone:
		for _, condition := range c.children {
			if condition.Check(certificateToken) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// ToString returns a human readable condition using the given indentation. Port of
// toString(String indent), including the "CriteriaListCondition: " prefix and the tab-deepened
// indentation handed to each child.
//
// Go has no null string; the Java `if (indent == null) indent = ""` branch is unreachable here.
// The `if (children != null)` guard is likewise vacuous over a nil slice, whose range is empty.
func (c *CompositeCondition) ToString(indent string) string {
	var builder strings.Builder
	builder.WriteString(indent)
	builder.WriteString("CriteriaListCondition: ")
	builder.WriteString(string(c.matchingCriteriaIndicator))
	builder.WriteByte('\n')
	indent += "\t"
	for _, condition := range c.children {
		builder.WriteString(condition.ToString(indent))
	}
	return builder.String()
}

// String ports toString(), i.e. toString("").
func (c *CompositeCondition) String() string {
	return c.ToString("")
}

// Equals ports equals(Object): same concrete type, same indicator and equal child lists.
//
// Java compares the children with List#equals, which delegates to each element's equals();
// this port therefore compares each child through its own Equals where the child provides one
// (every condition in this package does), falling back to interface comparison otherwise.
func (c *CompositeCondition) Equals(object any) bool {
	that, ok := object.(*CompositeCondition)
	if !ok || that == nil {
		return false
	}
	if c == that {
		return true
	}
	if c.matchingCriteriaIndicator != that.matchingCriteriaIndicator {
		return false
	}
	if (c.children == nil) != (that.children == nil) || len(c.children) != len(that.children) {
		return false
	}
	for index := range c.children {
		if !compositeConditionEquals(c.children[index], that.children[index]) {
			return false
		}
	}
	return true
}

// compositeConditionEquals compares two conditions the way Java's List#equals compares two
// elements: through the left element's equals(Object).
func compositeConditionEquals(left, right tslmodel.Condition) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	type equatable interface{ Equals(object any) bool }
	if comparable, ok := left.(equatable); ok {
		return comparable.Equals(right)
	}
	return left == right
}

var _ tslmodel.Condition = (*CompositeCondition)(nil)
