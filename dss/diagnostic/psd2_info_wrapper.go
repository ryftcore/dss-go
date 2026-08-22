// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/PSD2InfoWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
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

// RoleOfPSPNames returns names of roles of PSP. Port of getRoleOfPSPNames().
func (w *PSD2InfoWrapper) RoleOfPSPNames() []string {
	var result []string
	for _, roleOfPSP := range w.psd2QcInfo.RolesOfPSP.All() {
		if roleOfPSP.Name != nil {
			result = append(result, *roleOfPSP.Name)
		}
	}
	return result
}

// RoleOfPSPOids returns OIDs of roles of PSP. Port of getRoleOfPSPOids().
func (w *PSD2InfoWrapper) RoleOfPSPOids() []enumerations.RoleOfPspOid {
	var result []enumerations.RoleOfPspOid
	for _, roleOfPSP := range w.psd2QcInfo.RolesOfPSP.All() {
		pspOid := roleOfPSP.Oid
		if pspOid != nil {
			result = append(result, enumerations.RoleOfPspOidFromOid(pspOid.Value))
		}
	}
	return result
}

// NcaId returns the Competent Authority Id. Port of getNcaId().
func (w *PSD2InfoWrapper) NcaId() string {
	if w.psd2QcInfo.NcaId != nil {
		return *w.psd2QcInfo.NcaId
	}
	return ""
}

// NcaName returns the Competent Authority name. Port of getNcaName().
func (w *PSD2InfoWrapper) NcaName() string {
	if w.psd2QcInfo.NcaName != nil {
		return *w.psd2QcInfo.NcaName
	}
	return ""
}
