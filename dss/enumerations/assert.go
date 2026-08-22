// Ported from dss-enumerations/.../Assert.java (DSS 6.5.RC1).
package enumerations

// Assert is the Trusted List assert list. See TS 119 612.
type Assert string

const (
	// AssertAll applies if all of the assertion shall be met.
	AssertAll Assert = "ALL"
	// AssertAtLeastOne applies if at least one of the assertion shall
	// be met.
	AssertAtLeastOne Assert = "AT_LEAST_ONE"
	// AssertNone applies if all of the assertion shall be met.
	AssertNone Assert = "NONE"
)

var assertValueTable = map[Assert]string{
	AssertAll:        "all",
	AssertAtLeastOne: "atLeastOne",
	AssertNone:       "none",
}

// AssertValues returns all constants in declaration order.
func AssertValues() []Assert {
	return []Assert{
		AssertAll,
		AssertAtLeastOne,
		AssertNone,
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
