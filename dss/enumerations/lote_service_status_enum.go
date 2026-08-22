// Ported from dss-enumerations/.../LoTEServiceStatusEnum.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// LoTEServiceStatusEnum implements LoTEServiceStatus.
type LoTEServiceStatusEnum string

const (
	// LoTEServiceStatusEnumPubEAAProviderNotified identifies a notified
	// Pub-EAA provider service.
	LoTEServiceStatusEnumPubEAAProviderNotified LoTEServiceStatusEnum = "PUB_EAA_PROVIDER_NOTIFIED"
	// LoTEServiceStatusEnumPubEAAProviderWithdrawn identifies a withdrawn
	// Pub-EAA provider service.
	LoTEServiceStatusEnumPubEAAProviderWithdrawn LoTEServiceStatusEnum = "PUB_EAA_PROVIDER_WITHDRAWN"
)

type loteServiceStatusEnumFields struct {
	uri   string
	label string
}

// loteServiceStatusEnumData holds the (statusUri, label) pair for each constant.
var loteServiceStatusEnumData = map[LoTEServiceStatusEnum]loteServiceStatusEnumFields{
	LoTEServiceStatusEnumPubEAAProviderNotified:  {"http://uri.etsi.org/19602/PubEAAProvidersList/SvcStatus/notified", "Notified Pub-EAA provider service"},
	LoTEServiceStatusEnumPubEAAProviderWithdrawn: {"http://uri.etsi.org/19602/PubEAAProvidersList/SvcStatus/withdrawn", "Withdrawn Pub-EAA provider service"},
}

// LoTEServiceStatusEnumValues returns all constants in declaration order.
func LoTEServiceStatusEnumValues() []LoTEServiceStatusEnum {
	return []LoTEServiceStatusEnum{
		LoTEServiceStatusEnumPubEAAProviderNotified,
		LoTEServiceStatusEnumPubEAAProviderWithdrawn,
	}
}

// LoTEServiceStatusEnumValueOf returns the constant matching the given Java enum name.
func LoTEServiceStatusEnumValueOf(name string) (LoTEServiceStatusEnum, error) {
	for _, v := range LoTEServiceStatusEnumValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant LoTEServiceStatusEnum.%s", name)
}

// URI returns the Service Status URI. Implements LoTEServiceStatus.
func (l LoTEServiceStatusEnum) URI() string {
	return loteServiceStatusEnumData[l].uri
}

// Label returns the user-friendly label. Implements LoTEServiceStatus.
func (l LoTEServiceStatusEnum) Label() string {
	return loteServiceStatusEnumData[l].label
}
