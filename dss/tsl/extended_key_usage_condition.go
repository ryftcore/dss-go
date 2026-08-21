// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/dto/condition/ExtendedKeyUsageCondition.java (DSS 6.5.RC1).
package tsl

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ExtendedKeyUsageCondition implements the ExtendedKeyUsage criterion.
//
// Presence: This field is optional.
//
// Description: It provides a non empty list of key purposes values to match with the
// correspondent KeyPurposes present in the ExtendedKeyUsage certificate Extension. The assertion
// is verified if the ExtendedKeyUsage Extension is present in the certificate and all key
// purposes provided are present in the certificate ExtendedKeyUsage Extension.
//
// Format: A non-empty sequence of KeyPurposes, whose semantic shall be as defined in X.509 for
// the ExtendedKeyUsage Extension. For the formal definition see ExtendedKeyUsage element in the
// schema referenced by clause C.2 (point 3).
//
// java.io.Serializable has no Go counterpart and is dropped.
type ExtendedKeyUsageCondition struct {
	// extendedKeyUsageOids is the list of extended key usages to check.
	extendedKeyUsageOids []string
}

// NewExtendedKeyUsageCondition is the default constructor. Port of
// ExtendedKeyUsageCondition(List<String>).
func NewExtendedKeyUsageCondition(oids []string) *ExtendedKeyUsageCondition {
	return &ExtendedKeyUsageCondition{extendedKeyUsageOids: oids}
}

// KeyPurposeIds returns the list of key purpose IDs to be checked against the certificate's
// extended key usage extension: possibly empty, never nil. Port of the final
// getKeyPurposeIds().
func (c *ExtendedKeyUsageCondition) KeyPurposeIds() []string {
	if c.extendedKeyUsageOids == nil {
		return []string{}
	}
	return c.extendedKeyUsageOids
}

// Check returns true if the condition is evaluated to true for the given certificate. Port of
// check(CertificateToken).
func (c *ExtendedKeyUsageCondition) Check(certificateToken *model.CertificateToken) bool {
	if utils.IsCollectionNotEmpty(c.extendedKeyUsageOids) {
		extendedKeyUsage := spi.CertificateExtensionsUtilsExtendedKeyUsage(certificateToken)
		if extendedKeyUsage == nil || utils.IsCollectionEmpty(extendedKeyUsage.Oids()) {
			return false
		}
		for _, oid := range c.extendedKeyUsageOids {
			if !extendedKeyUsageConditionContains(extendedKeyUsage.Oids(), oid) {
				return false
			}
		}
	}
	return true
}

// ToString returns a human readable condition using the given indentation. Port of
// toString(String indent).
//
// Go has no null string; the Java `if (indent == null) indent = ""` branch is unreachable here.
func (c *ExtendedKeyUsageCondition) ToString(indent string) string {
	var builder strings.Builder
	builder.WriteString(indent)
	builder.WriteString("ExtendedKeyUsageCondition: ")
	builder.WriteString(extendedKeyUsageConditionJavaListString(c.extendedKeyUsageOids))
	builder.WriteByte('\n')
	return builder.String()
}

// String ports toString(), i.e. toString("").
func (c *ExtendedKeyUsageCondition) String() string {
	return c.ToString("")
}

// Equals ports equals(Object): same concrete type and equal OID lists.
func (c *ExtendedKeyUsageCondition) Equals(object any) bool {
	that, ok := object.(*ExtendedKeyUsageCondition)
	if !ok || that == nil {
		return false
	}
	if c == that {
		return true
	}
	return extendedKeyUsageConditionListEquals(c.extendedKeyUsageOids, that.extendedKeyUsageOids)
}

// extendedKeyUsageConditionContains ports List#contains(Object).
func extendedKeyUsageConditionContains(list []string, value string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

// extendedKeyUsageConditionJavaListString renders a List<String> the way Java's
// StringBuilder#append(Object) does: "null" for a null list, otherwise "[a, b]".
func extendedKeyUsageConditionJavaListString(list []string) string {
	if list == nil {
		return "null"
	}
	return "[" + strings.Join(list, ", ") + "]"
}

// extendedKeyUsageConditionListEquals ports Objects.equals(List, List): a null list is equal
// only to another null list, never to an empty one.
func extendedKeyUsageConditionListEquals(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

var _ tslmodel.Condition = (*ExtendedKeyUsageCondition)(nil)
