// Ported from dss-model/.../claim/Claim.java (DSS 6.5.RC1).
package claim

import "time"

// Claim defines a claim that may be made selectively disclosable.
//
// java.io.Serializable is dropped silently (no Go counterpart).
//
// Nullability notes on value accessors (Java's Number/Boolean/Date/String
// are nullable wrapper types):
//   - StringValue: "" stands for Java null (matches this codebase's
//     string-nullability convention elsewhere in dss/model).
//   - NumberValue: returns nil for Java null, otherwise the underlying Go
//     numeric value in whatever type it was constructed with (Java's
//     Number is itself a polymorphic wrapper over Integer/Long/Double/
//     BigDecimal/...; `any` is the closest faithful mapping).
//   - BooleanValue/DateValue: nil-able via pointer, since a false/zero
//     value is itself meaningful and must be distinguished from "absent".
type Claim interface {
	// Name gets the claim name. Ports Claim#getName.
	Name() string

	// IsSelectivelyDisclosable gets whether the claim was made selectively
	// disclosable and its value has been obtained from a provided
	// disclosure. Ports Claim#isSelectivelyDisclosable.
	IsSelectivelyDisclosable() bool

	// Namespace gets the origin namespace of the claim (NOTE: used in
	// mdoc). Ports Claim#getNamespace.
	Namespace() string

	// Parent gets the parent claim, when applicable (e.g. for claims
	// nested within a map or an array). Ports Claim#getParent.
	Parent() Claim

	// ListValue gets the value as a list. If the value is nil or not of a
	// list type, returns nil. Ports Claim#getListValue.
	ListValue() []Claim

	// BinaryValue gets the value as binaries. If the value is nil or not
	// of binaries type, returns nil. Ports Claim#getBinaryValue.
	BinaryValue() []byte

	// BooleanValue gets the value as a boolean. If the value is nil or not
	// of a boolean type, returns nil. Ports Claim#getBooleanValue.
	BooleanValue() *bool

	// DateValue gets the value as a date. If the value is nil or not of a
	// date type, returns nil. Ports Claim#getDateValue.
	DateValue() *time.Time

	// MapValue gets the value as a map. If the value is nil or not of a
	// map type, returns nil. Ports Claim#getMapValue.
	MapValue() map[string]Claim

	// NumberValue gets the value as a number. If the value is nil or not
	// of a number type, returns nil. Ports Claim#getNumberValue.
	NumberValue() any

	// StringValue gets the value as a string. If the value is nil or not
	// of a string type, returns "". Ports Claim#getStringValue.
	StringValue() string

	// IsStringValueType gets whether the claim value is of String type.
	// Ports Claim#isStringValueType.
	IsStringValueType() bool

	// IsBinaryValueType gets whether the claim value is of Binary type.
	// Ports Claim#isBinaryValueType.
	IsBinaryValueType() bool

	// IsBooleanValueType gets whether the claim value is of Boolean type.
	// Ports Claim#isBooleanValueType.
	IsBooleanValueType() bool

	// IsNumberValueType gets whether the claim value is of Number type.
	// Ports Claim#isNumberValueType.
	IsNumberValueType() bool

	// IsDateValueType gets whether the claim value is of Date type. Ports
	// Claim#isDateValueType.
	IsDateValueType() bool

	// IsArrayValueType gets whether the claim value is of Array type.
	// Ports Claim#isArrayValueType.
	IsArrayValueType() bool

	// IsMapValueType gets whether the claim value is of Map type. Ports
	// Claim#isMapValueType.
	IsMapValueType() bool

	// IsNullValueType gets whether the claim value is of Null type. Ports
	// Claim#isNullValueType.
	IsNullValueType() bool

	// IsSubresourceIntegrityType gets whether the claim provides an
	// integrity validation material for another claim, implementation
	// specific. Ports Claim#isSubresourceIntegrityType.
	IsSubresourceIntegrityType() bool

	// IsNullOrEmpty gets whether the value of the claim is null or empty.
	// Ports Claim#isNullOrEmpty.
	IsNullOrEmpty() bool

	// ValueAsString converts the claim's value to its corresponding
	// string representation. Ports Claim#getValueAsString.
	ValueAsString() string
}
