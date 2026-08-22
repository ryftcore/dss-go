// Ported from dss-enumerations/.../LoTETypeEnum.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// LoTETypeEnum contains a list of registered LoTE types within ETSI TS 119 602.
// Implements ListType.
type LoTETypeEnum string

const (
	// LoTETypeEnumEUPIDProvidersList is the EU PID Providers List.
	LoTETypeEnumEUPIDProvidersList LoTETypeEnum = "EUPIDProvidersList"
	// LoTETypeEnumEUWalletProvidersList is the EU Wallet Providers List.
	LoTETypeEnumEUWalletProvidersList LoTETypeEnum = "EUWalletProvidersList"
	// LoTETypeEnumEUWRPACProvidersList is the EU WRPAC Providers List.
	LoTETypeEnumEUWRPACProvidersList LoTETypeEnum = "EUWRPACProvidersList"
	// LoTETypeEnumEUWRPRCProvidersList is the EU WRPRC Providers List.
	LoTETypeEnumEUWRPRCProvidersList LoTETypeEnum = "EUWRPRCProvidersList"
	// LoTETypeEnumEUPubEAAProvidersList is the EU Pub-EAA Providers List.
	LoTETypeEnumEUPubEAAProvidersList LoTETypeEnum = "EUPubEAAProvidersList"
	// LoTETypeEnumEURegistrarsAndRegistersList is the EU Registrars and
	// Registers List.
	LoTETypeEnumEURegistrarsAndRegistersList LoTETypeEnum = "EURegistrarsAndRegistersList"
)

type loteTypeEnumFields struct {
	uri   string
	label string
}

// loteTypeEnumData holds the (uri, label) pair for each constant.
var loteTypeEnumData = map[LoTETypeEnum]loteTypeEnumFields{
	LoTETypeEnumEUPIDProvidersList:           {"http://uri.etsi.org/19602/LoTEType/EUPIDProvidersList", "EU List of providers of person identity data"},
	LoTETypeEnumEUWalletProvidersList:        {"http://uri.etsi.org/19602/LoTEType/EUWalletProvidersList", "EU List of wallet providers"},
	LoTETypeEnumEUWRPACProvidersList:         {"https://uri.etsi.org/19602/LoTEType/EUWRPACProvidersList", "EU List of providers of wallet relying party access certificates"},
	LoTETypeEnumEUWRPRCProvidersList:         {"http://uri.etsi.org/19602/LoTEType/EUWRPRCProvidersList", "EU List of providers of wallet relying party registration certificates"},
	LoTETypeEnumEUPubEAAProvidersList:        {"http://uri.etsi.org/19602/LoTEType/EUPubEAAProvidersList", "EU List of public sector bodies issuing electronic attestation of attributes"},
	LoTETypeEnumEURegistrarsAndRegistersList: {"http://uri.etsi.org/19602/LoTEType/EURegistrarsAndRegistersList", "EU List of registrars and registers"},
}

// LoTETypeEnumValues returns all constants in declaration order.
func LoTETypeEnumValues() []LoTETypeEnum {
	return []LoTETypeEnum{
		LoTETypeEnumEUPIDProvidersList,
		LoTETypeEnumEUWalletProvidersList,
		LoTETypeEnumEUWRPACProvidersList,
		LoTETypeEnumEUWRPRCProvidersList,
		LoTETypeEnumEUPubEAAProvidersList,
		LoTETypeEnumEURegistrarsAndRegistersList,
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
