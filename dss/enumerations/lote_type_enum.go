// Ported from dss-enumerations/.../LoTETypeEnum.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// LoTETypeEnum contains a list of registered LoTE types within ETSI TS 119 602.
// Implements ListType.
type LoTETypeEnum string

const (
	// LoTETypeEnum_EUPIDProvidersList is the EU PID Providers List.
	LoTETypeEnum_EUPIDProvidersList LoTETypeEnum = "EUPIDProvidersList"
	// LoTETypeEnum_EUWalletProvidersList is the EU Wallet Providers List.
	LoTETypeEnum_EUWalletProvidersList LoTETypeEnum = "EUWalletProvidersList"
	// LoTETypeEnum_EUWRPACProvidersList is the EU WRPAC Providers List.
	LoTETypeEnum_EUWRPACProvidersList LoTETypeEnum = "EUWRPACProvidersList"
	// LoTETypeEnum_EUWRPRCProvidersList is the EU WRPRC Providers List.
	LoTETypeEnum_EUWRPRCProvidersList LoTETypeEnum = "EUWRPRCProvidersList"
	// LoTETypeEnum_EUPubEAAProvidersList is the EU Pub-EAA Providers List.
	LoTETypeEnum_EUPubEAAProvidersList LoTETypeEnum = "EUPubEAAProvidersList"
	// LoTETypeEnum_EURegistrarsAndRegistersList is the EU Registrars and
	// Registers List.
	LoTETypeEnum_EURegistrarsAndRegistersList LoTETypeEnum = "EURegistrarsAndRegistersList"
)

type loteTypeEnumFields struct {
	uri   string
	label string
}

// loteTypeEnumData holds the (uri, label) pair for each constant.
var loteTypeEnumData = map[LoTETypeEnum]loteTypeEnumFields{
	LoTETypeEnum_EUPIDProvidersList:           {"http://uri.etsi.org/19602/LoTEType/EUPIDProvidersList", "EU List of providers of person identity data"},
	LoTETypeEnum_EUWalletProvidersList:        {"http://uri.etsi.org/19602/LoTEType/EUWalletProvidersList", "EU List of wallet providers"},
	LoTETypeEnum_EUWRPACProvidersList:         {"https://uri.etsi.org/19602/LoTEType/EUWRPACProvidersList", "EU List of providers of wallet relying party access certificates"},
	LoTETypeEnum_EUWRPRCProvidersList:         {"http://uri.etsi.org/19602/LoTEType/EUWRPRCProvidersList", "EU List of providers of wallet relying party registration certificates"},
	LoTETypeEnum_EUPubEAAProvidersList:        {"http://uri.etsi.org/19602/LoTEType/EUPubEAAProvidersList", "EU List of public sector bodies issuing electronic attestation of attributes"},
	LoTETypeEnum_EURegistrarsAndRegistersList: {"http://uri.etsi.org/19602/LoTEType/EURegistrarsAndRegistersList", "EU List of registrars and registers"},
}

// LoTETypeEnumValues returns all constants in declaration order.
func LoTETypeEnumValues() []LoTETypeEnum {
	return []LoTETypeEnum{
		LoTETypeEnum_EUPIDProvidersList,
		LoTETypeEnum_EUWalletProvidersList,
		LoTETypeEnum_EUWRPACProvidersList,
		LoTETypeEnum_EUWRPRCProvidersList,
		LoTETypeEnum_EUPubEAAProvidersList,
		LoTETypeEnum_EURegistrarsAndRegistersList,
	}
}

// LoTETypeEnumValueOf returns the constant matching the given Java enum name.
func LoTETypeEnumValueOf(name string) (LoTETypeEnum, error) {
	for _, v := range LoTETypeEnumValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant LoTETypeEnum.%s", name)
}

// URI returns the URI associated with the LoTE type. Implements ListType.
func (l LoTETypeEnum) URI() string {
	return loteTypeEnumData[l].uri
}

// Label returns the name of the LoTE type. Implements ListType.
func (l LoTETypeEnum) Label() string {
	return loteTypeEnumData[l].label
}
