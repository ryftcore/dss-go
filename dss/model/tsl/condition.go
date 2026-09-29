// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/Condition.java (DSS 6.5.RC1).
package tsl

import (
	"reflect"

	"github.com/ryftcore/dss-go/dss/model"
)

// Condition represents a condition defined in the trusted list on a certificate.
//
// java.io.Serializable has no Go counterpart and is dropped.
type Condition interface {
	// Check returns true if the condition is evaluated to true for the given certificate.
	Check(certificateToken *model.CertificateToken) bool
	// ToString returns a human readable condition using the given indentation.
	//
	// Named ToString (not String) because it takes an argument and therefore does not
	// implement fmt.Stringer; it ports Java's Condition#toString(String indent) overload.
	ToString(indent string) string
}

// conditionEquals ports Objects.equals(condition, that.condition) as used by
// ConditionForQualifiers#equals and CertificateContentEquivalence#equals: two nil conditions
// are equal, a nil and a non-nil one are not, and otherwise the condition's own structural
// equals decides. Every Condition implementation in the port (dss/tsl/*_condition.go) provides
// Equals(object any) bool, the counterpart of the Java classes' equals(Object) overrides;
// comparing the interface values with == would compare pointer identity instead, so two
// independently built but equal conditions would not be equal as they are upstream. An
// implementation without an Equals method falls back to reflect.DeepEqual, which never panics
// on an uncomparable dynamic type the way == can.
func conditionEquals(a, b Condition) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if equatable, ok := a.(interface{ Equals(object any) bool }); ok {
		return equatable.Equals(b)
	}
	return reflect.DeepEqual(a, b)
}
