// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/Condition.java (DSS 6.5.RC1).
package tsl

import "github.com/utain/esig/dss/model"

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
