// Ported from dss-enumerations/.../ASiCContainerType.java (DSS 6.5.RC1).

package enumerations

import (
	"fmt"
	"strings"
)

// ASiCContainerType defines possible types for an ASiC container.
type ASiCContainerType string

const (
	// ASiCContainerType_ASiC_S is the Associated Signature Container Simple.
	ASiCContainerType_ASiC_S ASiCContainerType = "ASiC_S"
	// ASiCContainerType_ASiC_E is the Associated Signature Container Extended.
	ASiCContainerType_ASiC_E ASiCContainerType = "ASiC_E"
)

// ASiCContainerTypeValues returns all ASiCContainerType constants in declaration order.
func ASiCContainerTypeValues() []ASiCContainerType {
	return []ASiCContainerType{ASiCContainerType_ASiC_S, ASiCContainerType_ASiC_E}
}

// ASiCContainerTypeValueByName returns the ASiCContainerType based on the name (String),
// accepting either '-' or '_' as separator (Java's valueByName).
func ASiCContainerTypeValueByName(name string) (ASiCContainerType, error) {
	return ASiCContainerTypeValueOf(strings.ReplaceAll(name, "-", "_"))
}

// ASiCContainerTypeValueOf returns the ASiCContainerType matching the given Java enum name.
func ASiCContainerTypeValueOf(name string) (ASiCContainerType, error) {
	for _, v := range ASiCContainerTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant ASiCContainerType.%s", name)
}

// String returns the Java toString() representation, with '_' replaced by '-'.
func (a ASiCContainerType) String() string {
	return strings.ReplaceAll(string(a), "_", "-")
}
