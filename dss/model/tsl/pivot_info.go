// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/PivotInfo.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/job"
)

// PivotInfo contains information about a pivot.
type PivotInfo struct {
	LOTLInfo

	// certificateStatusMap maps certificates to their change statuses in the current pivot.
	certificateStatusMap map[*model.CertificateToken]CertificatePivotStatus
	// lotlLocation is the associated XML LOTL location.
	lotlLocation string
}

// NewPivotInfo is the default constructor.
//
// certificates maps CertificateToken to CertificatePivotStatus (map between certificates and
// their statuses in the current pivot); lotlLocation is the associated LOTL location.
//
// JUDGMENT CALL: LOTLInfo/NewLOTLInfo are out of this manifest (owned by whichever chunk ports
// tsl/LOTLInfo.java and tsl/TLInfo.java, same target package). Java's PivotInfo overrides the
// protected virtual method buildIdentifier(), invoked polymorphically from deep inside the
// AbstractDocumentInfo constructor chain so the identifier is built with the most-derived
// class's logic. Go embedding gives no such virtual dispatch: LOTLInfo's own construction
// cannot call back into PivotInfo.BuildIdentifier. This port defines BuildIdentifier/IsPivot as
// ordinary methods that shadow the embedded LOTLInfo ones for direct calls on a *PivotInfo, but
// any code that holds a value only as a LOTLInfo (or relies on identifier construction
// happening inside NewLOTLInfo) will get LOTLInfo's behavior, not PivotInfo's. The integrator
// owning LOTLInfo must decide how identifier construction is actually wired (e.g. a
// buildIdentifier function/interface passed into NewLOTLInfo, or lazy computation on first
// access) and this constructor may need to change accordingly.
func NewPivotInfo(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo TLParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string,
	certificates map[*model.CertificateToken]CertificatePivotStatus, lotlLocation string) *PivotInfo {
	p := &PivotInfo{
		LOTLInfo:             NewLOTLInfo(downloadCacheInfo, parsingCacheInfo, validationCacheInfo, url),
		certificateStatusMap: certificates,
		lotlLocation:         lotlLocation,
	}
	return p
}

// CertificateStatusMap returns a map of certificate tokens with a status regarding the current
// pivot.
func (p *PivotInfo) CertificateStatusMap() map[*model.CertificateToken]CertificatePivotStatus {
	return p.certificateStatusMap
}

// LOTLLocation returns the associated LOTL location url.
func (p *PivotInfo) LOTLLocation() string {
	return p.lotlLocation
}

// IsPivot overrides LOTLInfo.IsPivot: a PivotInfo is always a pivot.
func (p *PivotInfo) IsPivot() bool {
	return true
}

// BuildIdentifier overrides LOTLInfo.BuildIdentifier: a PivotInfo's identifier is a
// PivotIdentifier.
func (p *PivotInfo) BuildIdentifier() model.Identifier {
	return NewPivotIdentifier(p)
}
