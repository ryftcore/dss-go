// Ported from dss-enumerations/.../LoTEServiceTypeIdentifierEnum.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// LoTEServiceTypeIdentifierEnum implements LoTEServiceTypeIdentifier.
type LoTEServiceTypeIdentifierEnum string

const (
	// LoTEServiceTypeIdentifierEnumPIDIssuance identifies the PID Issuance service type.
	LoTEServiceTypeIdentifierEnumPIDIssuance LoTEServiceTypeIdentifierEnum = "PID_ISSUANCE"
	// LoTEServiceTypeIdentifierEnumPIDRevocation identifies the PID Revocation service type.
	LoTEServiceTypeIdentifierEnumPIDRevocation LoTEServiceTypeIdentifierEnum = "PID_REVOCATION"
	// LoTEServiceTypeIdentifierEnumWalletIssuance identifies the Wallet Solution
	// Issuance service type.
	LoTEServiceTypeIdentifierEnumWalletIssuance LoTEServiceTypeIdentifierEnum = "WALLET_ISSUANCE"
	// LoTEServiceTypeIdentifierEnumWalletRevocation identifies the Wallet Solution
	// Revocation service type.
	LoTEServiceTypeIdentifierEnumWalletRevocation LoTEServiceTypeIdentifierEnum = "WALLET_REVOCATION"
	// LoTEServiceTypeIdentifierEnumWRPACIssuance identifies the WRPAC Issuance service type.
	LoTEServiceTypeIdentifierEnumWRPACIssuance LoTEServiceTypeIdentifierEnum = "WRPAC_ISSUANCE"
	// LoTEServiceTypeIdentifierEnumWRPACRevocation identifies the WRPAC Revocation service type.
	LoTEServiceTypeIdentifierEnumWRPACRevocation LoTEServiceTypeIdentifierEnum = "WRPAC_REVOCATION"
	// LoTEServiceTypeIdentifierEnumWRPRCIssuance identifies the WRPRC Issuance service type.
	LoTEServiceTypeIdentifierEnumWRPRCIssuance LoTEServiceTypeIdentifierEnum = "WRPRC_ISSUANCE"
	// LoTEServiceTypeIdentifierEnumWRPRCRevocation identifies the WRPRC Revocation service type.
	LoTEServiceTypeIdentifierEnumWRPRCRevocation LoTEServiceTypeIdentifierEnum = "WRPRC_REVOCATION"
	// LoTEServiceTypeIdentifierEnumPubEAAIssuance identifies the Pub-EAA Issuance service type.
	LoTEServiceTypeIdentifierEnumPubEAAIssuance LoTEServiceTypeIdentifierEnum = "PUB_EAA_ISSUANCE"
	// LoTEServiceTypeIdentifierEnumPubEAARevocation identifies the Pub-EAA Revocation service type.
	LoTEServiceTypeIdentifierEnumPubEAARevocation LoTEServiceTypeIdentifierEnum = "PUB_EAA_REVOCATION"
	// LoTEServiceTypeIdentifierEnumRegister identifies the Register service type.
	LoTEServiceTypeIdentifierEnumRegister LoTEServiceTypeIdentifierEnum = "REGISTER"
)

type loteServiceTypeIdentifierEnumFields struct {
	uri   string
	label string
}

// loteServiceTypeIdentifierEnumData holds the (stiUri, label) pair for each constant.
var loteServiceTypeIdentifierEnumData = map[LoTEServiceTypeIdentifierEnum]loteServiceTypeIdentifierEnumFields{
	LoTEServiceTypeIdentifierEnumPIDIssuance:      {"http://uri.etsi.org/19602/SvcType/PID/Issuance", "PID Issuance"},
	LoTEServiceTypeIdentifierEnumPIDRevocation:    {"http://uri.etsi.org/19602/SvcType/PID/Revocation", "PID Revocation"},
	LoTEServiceTypeIdentifierEnumWalletIssuance:   {"http://uri.etsi.org/19602/SvcType/WalletSolution/Issuance", "Wallet Solution Issuance"},
	LoTEServiceTypeIdentifierEnumWalletRevocation: {"http://uri.etsi.org/19602/SvcType/WalletSolution/Revocation", "Wallet Solution Revocation"},
	LoTEServiceTypeIdentifierEnumWRPACIssuance:    {"http://uri.etsi.org/19602/SvcType/WRPAC/Issuance", "WRPAC Issuance"},
	LoTEServiceTypeIdentifierEnumWRPACRevocation:  {"http://uri.etsi.org/19602/SvcType/WRPAC/Revocation", "WRPAC Revocation"},
	LoTEServiceTypeIdentifierEnumWRPRCIssuance:    {"http://uri.etsi.org/19602/SvcType/WRPRC/Issuance", "WRPRC Issuance"},
	LoTEServiceTypeIdentifierEnumWRPRCRevocation:  {"http://uri.etsi.org/19602/SvcType/WRPRC/Revocation", "WRPRC Revocation"},
	LoTEServiceTypeIdentifierEnumPubEAAIssuance:   {"http://uri.etsi.org/19602/SvcType/PubEAA/Issuance", "Pub-EAA Issuance"},
	LoTEServiceTypeIdentifierEnumPubEAARevocation: {"http://uri.etsi.org/19602/SvcType/PubEAA/Revocation", "Pub-EAA Revocation"},
	LoTEServiceTypeIdentifierEnumRegister:         {"http://uri.etsi.org/19602/SvcType/Register", "Register"},
}

// LoTEServiceTypeIdentifierEnumValues returns all constants in declaration order.
func LoTEServiceTypeIdentifierEnumValues() []LoTEServiceTypeIdentifierEnum {
	return []LoTEServiceTypeIdentifierEnum{
		LoTEServiceTypeIdentifierEnumPIDIssuance,
		LoTEServiceTypeIdentifierEnumPIDRevocation,
		LoTEServiceTypeIdentifierEnumWalletIssuance,
		LoTEServiceTypeIdentifierEnumWalletRevocation,
		LoTEServiceTypeIdentifierEnumWRPACIssuance,
		LoTEServiceTypeIdentifierEnumWRPACRevocation,
		LoTEServiceTypeIdentifierEnumWRPRCIssuance,
		LoTEServiceTypeIdentifierEnumWRPRCRevocation,
		LoTEServiceTypeIdentifierEnumPubEAAIssuance,
		LoTEServiceTypeIdentifierEnumPubEAARevocation,
		LoTEServiceTypeIdentifierEnumRegister,
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
