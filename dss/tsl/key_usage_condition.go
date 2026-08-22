// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/dto/condition/KeyUsageCondition.java (DSS 6.5.RC1).
package tsl

import (
	"strconv"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
)

// KeyUsageCondition is a condition based on the certificate key usage.
//
// java.io.Serializable has no Go counterpart and is dropped.
type KeyUsageCondition struct {
	// bit is the key usage bit to be checked.
	bit enumerations.KeyUsageBit

	// value is the required value of the key usage bit.
	value bool
}

// NewKeyUsageCondition constructs a new KeyUsageCondition from a KeyUsageBit and the required
// value of that bit. Port of KeyUsageCondition(KeyUsageBit, boolean).
//
// Panics with the Java message when bit is empty (Objects.requireNonNull(bit, "key usage");
// the empty KeyUsageBit is the Go stand-in for a null enum reference).
func NewKeyUsageCondition(bit enumerations.KeyUsageBit, value bool) *KeyUsageCondition {
	if bit == "" {
		panic("key usage")
	}
	return &KeyUsageCondition{bit: bit, value: value}
}

// NewKeyUsageConditionFromName constructs a new KeyUsageCondition from a key usage enum name
// and the required value of that bit. Port of KeyUsageCondition(String, boolean), which
// delegates to KeyUsageBit.valueOf(usage); Java's thrown IllegalArgumentException becomes a
// returned error, per PORTING.md.
func NewKeyUsageConditionFromName(usage string, value bool) (*KeyUsageCondition, error) {
	bit, err := enumerations.KeyUsageBitValueOf(usage)
	if err != nil {
		return nil, err
	}
	return NewKeyUsageCondition(bit, value), nil
}

// Bit returns the key usage to be checked. Port of the final getBit().
func (c *KeyUsageCondition) Bit() enumerations.KeyUsageBit {
	return c.bit
}

// Value returns the required bit value of the key usage to be checked. Port of the final
// getValue().
func (c *KeyUsageCondition) Value() bool {
	return c.value
}

// Check returns true if the condition is evaluated to true for the given certificate. Port of
// check(CertificateToken).
func (c *KeyUsageCondition) Check(certificateToken *model.CertificateToken) bool {
	keyUsage := certificateToken.CheckKeyUsage(c.bit)
	return keyUsage == c.value
}

// ToString returns a human readable condition using the given indentation. Port of
// toString(String indent).
//
// Go has no null string; the Java `if (indent == null) indent = ""` branch is unreachable here.
// bit.name() is the constant's Go value (PORTING.md's enum rule), and Java's
// StringBuilder#append(boolean) renders "true"/"false".
func (c *KeyUsageCondition) ToString(indent string) string {
	var builder strings.Builder
	builder.WriteString(indent)
	builder.WriteString("KeyUsageCondition: ")
	builder.WriteString(string(c.bit))
	builder.WriteByte('=')
	builder.WriteString(strconv.FormatBool(c.value))
	builder.WriteByte('\n')
	return builder.String()
}

// String ports toString(), i.e. toString("").
func (c *KeyUsageCondition) String() string {
	return c.ToString("")
}

// Equals ports equals(Object): same concrete type, same bit and same required value.
func (c *KeyUsageCondition) Equals(object any) bool {
	that, ok := object.(*KeyUsageCondition)
	if !ok || that == nil {
		return false
	}
	if c == that {
		return true
	}
	return c.value == that.value && c.bit == that.bit
}

var _ tslmodel.Condition = (*KeyUsageCondition)(nil)
