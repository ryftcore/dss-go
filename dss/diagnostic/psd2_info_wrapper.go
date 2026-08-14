// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/PSD2InfoWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// PSD2InfoWrapper provides a user-friendly interface for dealing with jaxb.XmlPSD2QcInfo.
type PSD2InfoWrapper struct {
	// psd2QcInfo is the wrapped XmlPSD2QcInfo object.
	psd2QcInfo *jaxb.XmlPSD2QcInfo
}

// NewPSD2InfoWrapper is the default constructor.
func NewPSD2InfoWrapper(psd2QcInfo *jaxb.XmlPSD2QcInfo) *PSD2InfoWrapper {
	return &PSD2InfoWrapper{psd2QcInfo: psd2QcInfo}
}

// GetRoleOfPSPNames returns names of roles of PSP. Port of getRoleOfPSPNames().
func (w *PSD2InfoWrapper) GetRoleOfPSPNames() []string {
	var result []string
	for _, roleOfPSP := range w.psd2QcInfo.RolesOfPSP {
		result = append(result, roleOfPSP.Name)
	}
	return result
}

// GetRoleOfPSPOids returns OIDs of roles of PSP. Port of getRoleOfPSPOids().
func (w *PSD2InfoWrapper) GetRoleOfPSPOids() []enumerations.RoleOfPspOid {
	var result []enumerations.RoleOfPspOid
	for _, roleOfPSP := range w.psd2QcInfo.RolesOfPSP {
		pspOid := roleOfPSP.Oid
		if pspOid != nil {
			result = append(result, enumerations.RoleOfPspOidFromOid(pspOid.Value))
		}
	}
	return result
}

// GetNcaId returns the Competent Authority Id. Port of getNcaId().
func (w *PSD2InfoWrapper) GetNcaId() string {
	return w.psd2QcInfo.NcaId
}

// GetNcaName returns the Competent Authority name. Port of getNcaName().
func (w *PSD2InfoWrapper) GetNcaName() string {
	return w.psd2QcInfo.NcaName
}
