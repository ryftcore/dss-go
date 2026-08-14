// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/PSD2QcType.java (DSS 6.5.RC1).
package extension

// PSD2QcType represents a PSD-2-QC type.
type PSD2QcType struct {
	// rolesOfPSP is a list of RoleOfPSPs.
	rolesOfPSP []*RoleOfPSP

	// ncaName is the NCA name.
	ncaName string

	// ncaId is the NCA Id.
	ncaId string
}

// NewPSD2QcType instantiates the object with null values. Ports the default constructor.
func NewPSD2QcType() *PSD2QcType {
	return &PSD2QcType{}
}

// RolesOfPSP gets a list of RoleOfPSPs.
func (p *PSD2QcType) RolesOfPSP() []*RoleOfPSP {
	return p.rolesOfPSP
}

// SetRolesOfPSP sets a list of RoleOfPSPs.
func (p *PSD2QcType) SetRolesOfPSP(rolesOfPSP []*RoleOfPSP) {
	p.rolesOfPSP = rolesOfPSP
}

// NcaName gets the NCA name.
func (p *PSD2QcType) NcaName() string {
	return p.ncaName
}

// SetNcaName sets the NCA name.
func (p *PSD2QcType) SetNcaName(ncaName string) {
	p.ncaName = ncaName
}

// NcaId gets the NCA Id.
func (p *PSD2QcType) NcaId() string {
	return p.ncaId
}

// SetNcaId sets the NCA Id.
func (p *PSD2QcType) SetNcaId(ncaId string) {
	p.ncaId = ncaId
}
