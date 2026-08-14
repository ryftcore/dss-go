// Ported from dss-enumerations/.../Assert.java (DSS 6.5.RC1).
package enumerations

// Assert is the Trusted List assert list. See TS 119 612.
type Assert string

const (
	// Assert_ALL applies if all of the assertion shall be met.
	Assert_ALL Assert = "ALL"
	// Assert_AT_LEAST_ONE applies if at least one of the assertion shall
	// be met.
	Assert_AT_LEAST_ONE Assert = "AT_LEAST_ONE"
	// Assert_NONE applies if all of the assertion shall be met.
	Assert_NONE Assert = "NONE"
)

var assertValueTable = map[Assert]string{
	Assert_ALL:          "all",
	Assert_AT_LEAST_ONE: "atLeastOne",
	Assert_NONE:         "none",
}

// AssertValues returns all constants in declaration order.
func AssertValues() []Assert {
	return []Assert{
		Assert_ALL,
		Assert_AT_LEAST_ONE,
		Assert_NONE,
	}
}

// Value returns the value of the enumeration, corresponding to the one
// defined in the XSD.
func (a Assert) Value() string {
	return assertValueTable[a]
}

// AssertValueOf returns the constant matching the given Java enum name.
func AssertValueOf(name string) (Assert, error) {
	for _, v := range AssertValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &assertInvalidValueError{name}
}

type assertInvalidValueError struct {
	name string
}

func (e *assertInvalidValueError) Error() string {
	return "no enum constant Assert." + e.name
}
