// Ported from dss-enumerations/.../EAAStatus.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// EAAStatus represents a status of an EAA.
type EAAStatus string

const (
	// EAAStatus_VALID: the status of the Referenced Token is valid, correct or legal.
	EAAStatus_VALID EAAStatus = "VALID"
	// EAAStatus_INVALID: the status of the Referenced Token is revoked, annulled,
	// taken back, recalled or cancelled.
	EAAStatus_INVALID EAAStatus = "INVALID"
	// EAAStatus_SUSPENDED: the status of the Referenced Token is temporarily invalid,
	// hanging, debarred from privilege. This status is usually temporary.
	EAAStatus_SUSPENDED EAAStatus = "SUSPENDED"
	// EAAStatus_APPLICATION_SPECIFIC: the status of the Referenced Token is
	// application specific.
	EAAStatus_APPLICATION_SPECIFIC EAAStatus = "APPLICATION_SPECIFIC"
	// EAAStatus_UNKNOWN: the EAA status is not known.
	EAAStatus_UNKNOWN EAAStatus = "UNKNOWN"
)

// eaaStatusBitValue maps each EAAStatus with a defined bit value to its bit
// representation. EAAStatus_UNKNOWN has no bit value (Java's null), matching the
// empty constructor overload.
var eaaStatusBitValue = map[EAAStatus]int{
	EAAStatus_VALID:                0x00,
	EAAStatus_INVALID:              0x01,
	EAAStatus_SUSPENDED:            0x02,
	EAAStatus_APPLICATION_SPECIFIC: 0x03,
}

// EAAStatusValues returns all EAAStatus constants in declaration order.
func EAAStatusValues() []EAAStatus {
	return []EAAStatus{
		EAAStatus_VALID,
		EAAStatus_INVALID,
		EAAStatus_SUSPENDED,
		EAAStatus_APPLICATION_SPECIFIC,
		EAAStatus_UNKNOWN,
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
	return e == EAAStatus_VALID
}

// BitValue returns the bit representation of the Status Type in a byte hex
// representation, and false if the status has no defined bit value (EAAStatus_UNKNOWN).
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
