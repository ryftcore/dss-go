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
	// RoleOfPspOid_PSP_AS: Account Servicing Payment Service Provider role.
	RoleOfPspOid_PSP_AS RoleOfPspOid = "PSP_AS"
	// RoleOfPspOid_PSP_PI: Payment Initiation Service Provider role.
	RoleOfPspOid_PSP_PI RoleOfPspOid = "PSP_PI"
	// RoleOfPspOid_PSP_AI: Account Information Service Provider role.
	RoleOfPspOid_PSP_AI RoleOfPspOid = "PSP_AI"
	// RoleOfPspOid_PSP_IC: Payment Service Provider issuing card-based
	// payment instruments role.
	RoleOfPspOid_PSP_IC RoleOfPspOid = "PSP_IC"
)

// roleOfPspOidFields holds the (description, oid) pair for each constant.
type roleOfPspOidFields struct {
	description string
	oid         string
}

// roleOfPspOidData holds the full field tuple for each constant, copied
// verbatim from the Java enum constructors.
var roleOfPspOidData = map[RoleOfPspOid]roleOfPspOidFields{
	RoleOfPspOid_PSP_AS: {"psp-as", "0.4.0.19495.1.1"},
	RoleOfPspOid_PSP_PI: {"psp-pi", "0.4.0.19495.1.2"},
	RoleOfPspOid_PSP_AI: {"psp-ai", "0.4.0.19495.1.3"},
	RoleOfPspOid_PSP_IC: {"psp-ic", "0.4.0.19495.1.4"},
}

// RoleOfPspOidValues returns all constants in declaration order.
func RoleOfPspOidValues() []RoleOfPspOid {
	return []RoleOfPspOid{
		RoleOfPspOid_PSP_AS,
		RoleOfPspOid_PSP_PI,
		RoleOfPspOid_PSP_AI,
		RoleOfPspOid_PSP_IC,
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
