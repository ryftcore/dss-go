// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/RoleOfPSP.java (DSS 6.5.RC1).
package extension

import "github.com/utain/esig/dss/enumerations"

// RoleOfPSP is an Object Identifier for roles of payment service providers.
type RoleOfPSP struct {
	// pspOid is the role OID.
	pspOid enumerations.RoleOfPspOid

	// pspName is the PSP name.
	pspName string
}

// NewRoleOfPSP instantiates the object with null values. Ports the default constructor.
func NewRoleOfPSP() *RoleOfPSP {
	return &RoleOfPSP{}
}

// PspOid gets the role OID.
func (r *RoleOfPSP) PspOid() enumerations.RoleOfPspOid {
	return r.pspOid
}

// SetPspOid sets the role OID.
func (r *RoleOfPSP) SetPspOid(pspOid enumerations.RoleOfPspOid) {
	r.pspOid = pspOid
}

// PspName gets the PSP name.
func (r *RoleOfPSP) PspName() string {
	return r.pspName
}

// SetPspName sets the PSP name.
func (r *RoleOfPSP) SetPspName(pspName string) {
	r.pspName = pspName
}
