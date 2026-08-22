// Ported from dss-enumerations/.../EAAStatus.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// EAAStatus represents a status of an EAA.
type EAAStatus string

const (
	// EAAStatusValid: the status of the Referenced Token is valid, correct or legal.
	EAAStatusValid EAAStatus = "VALID"
	// EAAStatusInvalid: the status of the Referenced Token is revoked, annulled,
	// taken back, recalled or cancelled.
	EAAStatusInvalid EAAStatus = "INVALID"
	// EAAStatusSuspended: the status of the Referenced Token is temporarily invalid,
	// hanging, debarred from privilege. This status is usually temporary.
	EAAStatusSuspended EAAStatus = "SUSPENDED"
	// EAAStatusApplicationSpecific: the status of the Referenced Token is
	// application specific.
	EAAStatusApplicationSpecific EAAStatus = "APPLICATION_SPECIFIC"
	// EAAStatusUnknown: the EAA status is not known.
	EAAStatusUnknown EAAStatus = "UNKNOWN"
)

// eaaStatusBitValue maps each EAAStatus with a defined bit value to its bit
// representation. EAAStatusUnknown has no bit value (Java's null), matching the
// empty constructor overload.
var eaaStatusBitValue = map[EAAStatus]int{
	EAAStatusValid:               0x00,
	EAAStatusInvalid:             0x01,
	EAAStatusSuspended:           0x02,
	EAAStatusApplicationSpecific: 0x03,
}

// EAAStatusValues returns all EAAStatus constants in declaration order.
func EAAStatusValues() []EAAStatus {
	return []EAAStatus{
		EAAStatusValid,
		EAAStatusInvalid,
		EAAStatusSuspended,
		EAAStatusApplicationSpecific,
		EAAStatusUnknown,
	}
}

// EAAStatusValueOf returns the EAAStatus matching the given Java enum name.
func EAAStatusValueOf(name string) (EAAStatus, error) {
	for _, v := range EAAStatusValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant EAAStatus.%s", name)
}

// IsValid checks if the EAA status is valid.
func (e EAAStatus) IsValid() bool {
	return e == EAAStatusValid
}

// BitValue returns the bit representation of the Status Type in a byte hex
// representation, and false if the status has no defined bit value (EAAStatusUnknown).
func (e EAAStatus) BitValue() (int, bool) {
	v, ok := eaaStatusBitValue[e]
	return v, ok
}

// EAAStatusForBitValue gets a corresponding EAAStatus for the given bitValue.
//
// Upstream appears to fall back to UNKNOWN, but that return is dead code: the loop
// compares an int against the Integer field with `bitValue == eaaStatus.bitValue`, which
// unboxes, and UNKNOWN — declared last — carries a null bitValue. Any value outside
// 0x00-0x03 therefore reaches UNKNOWN and throws NullPointerException before the fallback
// can run. That failure is reproduced here as an error return rather than silently
// resolving to UNKNOWN, which would let an unrecognized status pass as a known one.
func EAAStatusForBitValue(bitValue int) (EAAStatus, error) {
	for _, v := range EAAStatusValues() {
		bv, ok := eaaStatusBitValue[v]
		if !ok {
			// Upstream throws here when unboxing the null bitValue.
			return "", fmt.Errorf("no EAAStatus for bit value %#02x", bitValue)
		}
		if bv == bitValue {
			return v, nil
		}
	}
	return "", fmt.Errorf("no EAAStatus for bit value %#02x", bitValue)
}
