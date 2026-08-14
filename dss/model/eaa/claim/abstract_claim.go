// Ported from dss-model/.../claim/AbstractClaim.java (DSS 6.5.RC1).
package claim

import (
	"reflect"
	"time"
)

// AbstractClaim is the common base for the disclosable claim
// implementations, designed for embedding by concrete Claim
// implementations. It supplies default (mostly zero-value) implementations
// for every Claim accessor except IsNullOrEmpty and ValueAsString, which
// remain abstract in Java (AbstractClaim itself does not implement the
// full Claim interface, matching Java's abstract class not being
// instantiable).
type AbstractClaim struct {
	// name is the name of the claim.
	name string

	// selectivelyDisclosable is whether the claim is selectively
	// disclosable.
	selectivelyDisclosable bool

	// namespace is the namespace of the claim's origin (used in mdoc).
	namespace string

	// parent is the parent claim, containing the current claim in its
	// body.
	parent Claim
}

// NewAbstractClaim ports the default (no-arg) constructor.
func NewAbstractClaim() AbstractClaim {
	return AbstractClaim{}
}

// NewAbstractClaimWithName ports the constructor with claim name provided.
func NewAbstractClaimWithName(name string) AbstractClaim {
	return AbstractClaim{name: name}
}

// NewAbstractClaimWithDisclosable ports the constructor with claim name
// and selectively disclosable status provided.
func NewAbstractClaimWithDisclosable(name string, selectivelyDisclosable bool) AbstractClaim {
	return NewAbstractClaimWithParent(name, selectivelyDisclosable, nil)
}

// NewAbstractClaimWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewAbstractClaimWithParent(name string, selectivelyDisclosable bool, parent Claim) AbstractClaim {
	return NewAbstractClaimFull(name, "", selectivelyDisclosable, parent)
}

// NewAbstractClaimFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewAbstractClaimFull(name, namespace string, selectivelyDisclosable bool, parent Claim) AbstractClaim {
	return AbstractClaim{
		name:                   name,
		namespace:              namespace,
		selectivelyDisclosable: selectivelyDisclosable,
		parent:                 parent,
	}
}

// Name gets the claim name.
func (a *AbstractClaim) Name() string { return a.name }

// IsSelectivelyDisclosable gets whether the claim was made selectively
// disclosable.
func (a *AbstractClaim) IsSelectivelyDisclosable() bool { return a.selectivelyDisclosable }

// Namespace gets the origin namespace of the claim (used in mdoc).
func (a *AbstractClaim) Namespace() string { return a.namespace }

// Parent gets the parent claim, when applicable.
func (a *AbstractClaim) Parent() Claim { return a.parent }

// StringValue is the default (nil) implementation. Ports
// AbstractClaim#getStringValue.
func (a *AbstractClaim) StringValue() string { return "" }

// NumberValue is the default (nil) implementation. Ports
// AbstractClaim#getNumberValue.
func (a *AbstractClaim) NumberValue() any { return nil }

// MapValue is the default (nil) implementation. Ports
// AbstractClaim#getMapValue.
func (a *AbstractClaim) MapValue() map[string]Claim { return nil }

// DateValue is the default (nil) implementation. Ports
// AbstractClaim#getDateValue.
func (a *AbstractClaim) DateValue() *time.Time { return nil }

// BooleanValue is the default (nil) implementation. Ports
// AbstractClaim#getBooleanValue.
func (a *AbstractClaim) BooleanValue() *bool { return nil }

// BinaryValue is the default (nil) implementation. Ports
// AbstractClaim#getBinaryValue.
func (a *AbstractClaim) BinaryValue() []byte { return nil }

// ListValue is the default (nil) implementation. Ports
// AbstractClaim#getListValue.
func (a *AbstractClaim) ListValue() []Claim { return nil }

// IsStringValueType is the default (false) implementation.
func (a *AbstractClaim) IsStringValueType() bool { return false }

// IsBinaryValueType is the default (false) implementation.
func (a *AbstractClaim) IsBinaryValueType() bool { return false }

// IsBooleanValueType is the default (false) implementation.
func (a *AbstractClaim) IsBooleanValueType() bool { return false }

// IsNumberValueType is the default (false) implementation.
func (a *AbstractClaim) IsNumberValueType() bool { return false }

// IsDateValueType is the default (false) implementation.
func (a *AbstractClaim) IsDateValueType() bool { return false }

// IsArrayValueType is the default (false) implementation.
func (a *AbstractClaim) IsArrayValueType() bool { return false }

// IsMapValueType is the default (false) implementation.
func (a *AbstractClaim) IsMapValueType() bool { return false }

// IsSubresourceIntegrityType is the default (false) implementation.
func (a *AbstractClaim) IsSubresourceIntegrityType() bool { return false }

// IsNullValueType is the default (false) implementation.
func (a *AbstractClaim) IsNullValueType() bool { return false }

// AbstractClaimString renders a claim the way AbstractClaim#toString does:
// the concrete claim's simple class name, then " {", the quoted claim name,
// " (disclosure)" when the claim is selectively disclosable, ": ", the
// claim's value as a string and "}".
//
// Java inherits toString() from the abstract base, where getClass() and
// getValueAsString() dispatch to the concrete subclass. Go has neither
// inheritance nor virtual dispatch from an embedded struct, so every
// concrete claim's String() method forwards here, passing itself.
func AbstractClaimString(c Claim) string {
	disclosure := ""
	if c.IsSelectivelyDisclosable() {
		disclosure = " (disclosure)"
	}
	return abstractClaimSimpleName(c) + " {" + "'" + c.Name() + "'" + disclosure + ": " + c.ValueAsString() + "}"
}

// abstractClaimSimpleName returns the Go type name of the concrete claim,
// which is kept identical to the Java class name getClass().getSimpleName()
// reports.
func abstractClaimSimpleName(c Claim) string {
	t := reflect.TypeOf(c)
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t == nil {
		return ""
	}
	return t.Name()
}

// Equals ports AbstractClaim#equals: NOTE this compares name,
// selectivelyDisclosable and parent only (namespace is intentionally
// excluded, matching upstream). Parent equality falls back to
// reflect.DeepEqual since Parent is a Claim interface value that may hold
// any concrete claim type.
func (a *AbstractClaim) Equals(other *AbstractClaim) bool {
	if a == other {
		return true
	}
	if other == nil {
		return false
	}
	return a.selectivelyDisclosable == other.selectivelyDisclosable &&
		a.name == other.name &&
		reflect.DeepEqual(a.parent, other.parent)
}
