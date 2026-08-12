// Ported from dss-enumerations/.../LoTEServiceTypeIdentifierEnum.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// LoTEServiceTypeIdentifierEnum implements LoTEServiceTypeIdentifier.
type LoTEServiceTypeIdentifierEnum string

const (
	// LoTEServiceTypeIdentifierEnum_PID_ISSUANCE identifies the PID Issuance service type.
	LoTEServiceTypeIdentifierEnum_PID_ISSUANCE LoTEServiceTypeIdentifierEnum = "PID_ISSUANCE"
	// LoTEServiceTypeIdentifierEnum_PID_REVOCATION identifies the PID Revocation service type.
	LoTEServiceTypeIdentifierEnum_PID_REVOCATION LoTEServiceTypeIdentifierEnum = "PID_REVOCATION"
	// LoTEServiceTypeIdentifierEnum_WALLET_ISSUANCE identifies the Wallet Solution
	// Issuance service type.
	LoTEServiceTypeIdentifierEnum_WALLET_ISSUANCE LoTEServiceTypeIdentifierEnum = "WALLET_ISSUANCE"
	// LoTEServiceTypeIdentifierEnum_WALLET_REVOCATION identifies the Wallet Solution
	// Revocation service type.
	LoTEServiceTypeIdentifierEnum_WALLET_REVOCATION LoTEServiceTypeIdentifierEnum = "WALLET_REVOCATION"
	// LoTEServiceTypeIdentifierEnum_WRPAC_ISSUANCE identifies the WRPAC Issuance service type.
	LoTEServiceTypeIdentifierEnum_WRPAC_ISSUANCE LoTEServiceTypeIdentifierEnum = "WRPAC_ISSUANCE"
	// LoTEServiceTypeIdentifierEnum_WRPAC_REVOCATION identifies the WRPAC Revocation service type.
	LoTEServiceTypeIdentifierEnum_WRPAC_REVOCATION LoTEServiceTypeIdentifierEnum = "WRPAC_REVOCATION"
	// LoTEServiceTypeIdentifierEnum_WRPRC_ISSUANCE identifies the WRPRC Issuance service type.
	LoTEServiceTypeIdentifierEnum_WRPRC_ISSUANCE LoTEServiceTypeIdentifierEnum = "WRPRC_ISSUANCE"
	// LoTEServiceTypeIdentifierEnum_WRPRC_REVOCATION identifies the WRPRC Revocation service type.
	LoTEServiceTypeIdentifierEnum_WRPRC_REVOCATION LoTEServiceTypeIdentifierEnum = "WRPRC_REVOCATION"
	// LoTEServiceTypeIdentifierEnum_PUB_EAA_ISSUANCE identifies the Pub-EAA Issuance service type.
	LoTEServiceTypeIdentifierEnum_PUB_EAA_ISSUANCE LoTEServiceTypeIdentifierEnum = "PUB_EAA_ISSUANCE"
	// LoTEServiceTypeIdentifierEnum_PUB_EAA_REVOCATION identifies the Pub-EAA Revocation service type.
	LoTEServiceTypeIdentifierEnum_PUB_EAA_REVOCATION LoTEServiceTypeIdentifierEnum = "PUB_EAA_REVOCATION"
	// LoTEServiceTypeIdentifierEnum_REGISTER identifies the Register service type.
	LoTEServiceTypeIdentifierEnum_REGISTER LoTEServiceTypeIdentifierEnum = "REGISTER"
)

type loteServiceTypeIdentifierEnumFields struct {
	uri   string
	label string
}

// loteServiceTypeIdentifierEnumData holds the (stiUri, label) pair for each constant.
var loteServiceTypeIdentifierEnumData = map[LoTEServiceTypeIdentifierEnum]loteServiceTypeIdentifierEnumFields{
	LoTEServiceTypeIdentifierEnum_PID_ISSUANCE:       {"http://uri.etsi.org/19602/SvcType/PID/Issuance", "PID Issuance"},
	LoTEServiceTypeIdentifierEnum_PID_REVOCATION:     {"http://uri.etsi.org/19602/SvcType/PID/Revocation", "PID Revocation"},
	LoTEServiceTypeIdentifierEnum_WALLET_ISSUANCE:    {"http://uri.etsi.org/19602/SvcType/WalletSolution/Issuance", "Wallet Solution Issuance"},
	LoTEServiceTypeIdentifierEnum_WALLET_REVOCATION:  {"http://uri.etsi.org/19602/SvcType/WalletSolution/Revocation", "Wallet Solution Revocation"},
	LoTEServiceTypeIdentifierEnum_WRPAC_ISSUANCE:     {"http://uri.etsi.org/19602/SvcType/WRPAC/Issuance", "WRPAC Issuance"},
	LoTEServiceTypeIdentifierEnum_WRPAC_REVOCATION:   {"http://uri.etsi.org/19602/SvcType/WRPAC/Revocation", "WRPAC Revocation"},
	LoTEServiceTypeIdentifierEnum_WRPRC_ISSUANCE:     {"http://uri.etsi.org/19602/SvcType/WRPRC/Issuance", "WRPRC Issuance"},
	LoTEServiceTypeIdentifierEnum_WRPRC_REVOCATION:   {"http://uri.etsi.org/19602/SvcType/WRPRC/Revocation", "WRPRC Revocation"},
	LoTEServiceTypeIdentifierEnum_PUB_EAA_ISSUANCE:   {"http://uri.etsi.org/19602/SvcType/PubEAA/Issuance", "Pub-EAA Issuance"},
	LoTEServiceTypeIdentifierEnum_PUB_EAA_REVOCATION: {"http://uri.etsi.org/19602/SvcType/PubEAA/Revocation", "Pub-EAA Revocation"},
	LoTEServiceTypeIdentifierEnum_REGISTER:           {"http://uri.etsi.org/19602/SvcType/Register", "Register"},
}

// LoTEServiceTypeIdentifierEnumValues returns all constants in declaration order.
func LoTEServiceTypeIdentifierEnumValues() []LoTEServiceTypeIdentifierEnum {
	return []LoTEServiceTypeIdentifierEnum{
		LoTEServiceTypeIdentifierEnum_PID_ISSUANCE,
		LoTEServiceTypeIdentifierEnum_PID_REVOCATION,
		LoTEServiceTypeIdentifierEnum_WALLET_ISSUANCE,
		LoTEServiceTypeIdentifierEnum_WALLET_REVOCATION,
		LoTEServiceTypeIdentifierEnum_WRPAC_ISSUANCE,
		LoTEServiceTypeIdentifierEnum_WRPAC_REVOCATION,
		LoTEServiceTypeIdentifierEnum_WRPRC_ISSUANCE,
		LoTEServiceTypeIdentifierEnum_WRPRC_REVOCATION,
		LoTEServiceTypeIdentifierEnum_PUB_EAA_ISSUANCE,
		LoTEServiceTypeIdentifierEnum_PUB_EAA_REVOCATION,
		LoTEServiceTypeIdentifierEnum_REGISTER,
	}
}

// LoTEServiceTypeIdentifierEnumValueOf returns the constant matching the given Java enum name.
func LoTEServiceTypeIdentifierEnumValueOf(name string) (LoTEServiceTypeIdentifierEnum, error) {
	for _, v := range LoTEServiceTypeIdentifierEnumValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant LoTEServiceTypeIdentifierEnum.%s", name)
}

// URI returns the Service Type Identifier URI. Implements LoTEServiceTypeIdentifier.
func (l LoTEServiceTypeIdentifierEnum) URI() string {
	return loteServiceTypeIdentifierEnumData[l].uri
}

// Label returns the user-friendly label defining the certificate approval status
// type. Implements LoTEServiceTypeIdentifier.
func (l LoTEServiceTypeIdentifierEnum) Label() string {
	return loteServiceTypeIdentifierEnumData[l].label
}
