// Ported from dss-enumerations/.../RoleOfPspOid.java (DSS 6.5.RC1).
//
// RoleOfPspOid ::= OBJECT IDENTIFIER -- Object Identifier arc for roles of
// payment service providers -- defined in the present document etsi-psd2-roles
// OBJECT IDENTIFIER ::={ itu-t(0) identified-organization(4) etsi(0)
// psd2(19495) id-roles(1) }
package enumerations

import "fmt"

// RoleOfPspOid represents a PSD2 role of payment service provider.
// Implements OidDescription.
type RoleOfPspOid string

const (
	// RoleOfPspOidPSPAs: Account Servicing Payment Service Provider role.
	RoleOfPspOidPSPAs RoleOfPspOid = "PSP_AS"
	// RoleOfPspOidPSPPI: Payment Initiation Service Provider role.
	RoleOfPspOidPSPPI RoleOfPspOid = "PSP_PI"
	// RoleOfPspOidPSPAI: Account Information Service Provider role.
	RoleOfPspOidPSPAI RoleOfPspOid = "PSP_AI"
	// RoleOfPspOidPSPIC: Payment Service Provider issuing card-based
	// payment instruments role.
	RoleOfPspOidPSPIC RoleOfPspOid = "PSP_IC"
)

// roleOfPspOidFields holds the (description, oid) pair for each constant.
type roleOfPspOidFields struct {
	description string
	oid         string
}

// roleOfPspOidData holds the full field tuple for each constant, copied
// verbatim from the Java enum constructors.
var roleOfPspOidData = map[RoleOfPspOid]roleOfPspOidFields{
	RoleOfPspOidPSPAs: {"psp-as", "0.4.0.19495.1.1"},
	RoleOfPspOidPSPPI: {"psp-pi", "0.4.0.19495.1.2"},
	RoleOfPspOidPSPAI: {"psp-ai", "0.4.0.19495.1.3"},
	RoleOfPspOidPSPIC: {"psp-ic", "0.4.0.19495.1.4"},
}

// RoleOfPspOidValues returns all constants in declaration order.
func RoleOfPspOidValues() []RoleOfPspOid {
	return []RoleOfPspOid{
		RoleOfPspOidPSPAs,
		RoleOfPspOidPSPPI,
		RoleOfPspOidPSPAI,
		RoleOfPspOidPSPIC,
	}
}

// RoleOfPspOidValueOf returns the RoleOfPspOid matching the given Java enum name.
func RoleOfPspOidValueOf(name string) (RoleOfPspOid, error) {
	for _, v := range RoleOfPspOidValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant RoleOfPspOid.%s", name)
}

// OID returns the OID. Implements OidDescription.
func (r RoleOfPspOid) OID() string {
	return roleOfPspOidData[r].oid
}

// Description returns the name of the PSP role. Implements OidDescription.
func (r RoleOfPspOid) Description() string {
	return roleOfPspOidData[r].description
}

// RoleOfPspOidFromOid returns a RoleOfPspOid by the given OID, or "" if not found.
func RoleOfPspOidFromOid(oid string) RoleOfPspOid {
	for _, v := range RoleOfPspOidValues() {
		if v.OID() == oid {
			return v
		}
	}
	return ""
}

// compile-time interface assertion.
var _ OidDescription = RoleOfPspOid("")
