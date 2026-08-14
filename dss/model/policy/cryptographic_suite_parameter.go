// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteParameter.java (DSS 6.5.RC1).
package policy

import "fmt"

// CryptographicSuiteParameter provides a representation of a "Parameter"
// element extracted from an ETSI TS 119 322 cryptographic suite
// catalogue.
//
// java.io.Serializable is dropped silently (no Go counterpart). Java's
// nullable Integer fields become *int.
type CryptographicSuiteParameter struct {
	// name is the value of the /dssc:Parameter/dssc:name element.
	name string

	// min is the value of the /dssc:Parameter/dssc:Min element.
	min *int

	// max is the value of the /dssc:Parameter/dssc:Max element.
	max *int
}

// NewCryptographicSuiteParameter is the default constructor.
func NewCryptographicSuiteParameter() *CryptographicSuiteParameter {
	return &CryptographicSuiteParameter{}
}

// Name gets the name identifier of parameter.
func (c *CryptographicSuiteParameter) Name() string { return c.name }

// SetName sets a value of the /dssc:Parameter/dssc:name element.
func (c *CryptographicSuiteParameter) SetName(name string) { c.name = name }

// Min gets the min value of parameter.
func (c *CryptographicSuiteParameter) Min() *int { return c.min }

// SetMin sets a value of the /dssc:Parameter/dssc:Min element.
func (c *CryptographicSuiteParameter) SetMin(min *int) { c.min = min }

// Max gets the max value of parameter.
func (c *CryptographicSuiteParameter) Max() *int { return c.max }

// SetMax sets a value of the /dssc:Parameter/dssc:Max element.
func (c *CryptographicSuiteParameter) SetMax(max *int) { c.max = max }

// CryptographicSuiteParameterCopy instantiates a new
// CryptographicSuiteParameter by copying the values of parameter. Ports
// CryptographicSuiteParameter#copy.
func CryptographicSuiteParameterCopy(parameter *CryptographicSuiteParameter) *CryptographicSuiteParameter {
	if parameter == nil {
		return nil
	}
	copyValue := NewCryptographicSuiteParameter()
	copyValue.name = parameter.name
	copyValue.min = parameter.min
	copyValue.max = parameter.max
	return copyValue
}

// Equals ports CryptographicSuiteParameter#equals.
func (c *CryptographicSuiteParameter) Equals(other *CryptographicSuiteParameter) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	return c.name == other.name && cryptographicSuiteParameterIntEqual(c.min, other.min) &&
		cryptographicSuiteParameterIntEqual(c.max, other.max)
}

// cryptographicSuiteParameterIntEqual compares two possibly-nil *int
// values.
func cryptographicSuiteParameterIntEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// String ports CryptographicSuiteParameter#toString.
func (c *CryptographicSuiteParameter) String() string {
	var minStr, maxStr string
	if c.min != nil {
		minStr = fmt.Sprintf("%d", *c.min)
	}
	if c.max != nil {
		maxStr = fmt.Sprintf("%d", *c.max)
	}
	return fmt.Sprintf("CryptographicSuiteParameter [name='%s', min=%s, max=%s]", c.name, minStr, maxStr)
}
