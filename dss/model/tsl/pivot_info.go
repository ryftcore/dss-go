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
	// identifier caches the value returned by DSSID, shadowing the embedded LOTLInfo's own
	// cache so PivotInfo's BuildIdentifier (not LOTLInfo's) is what gets cached and returned.
	identifier model.Identifier
}

// NewPivotInfo is the default constructor.
//
// certificates maps CertificateToken to CertificatePivotStatus (map between certificates and
// their statuses in the current pivot); lotlLocation is the associated LOTL location.
//
// JUDGMENT CALL: Java's PivotInfo overrides the protected virtual method buildIdentifier(),
// invoked polymorphically from AbstractDocumentInfo#getDSSId() so the identifier is built with
// the most-derived class's logic. Go embedding gives no such virtual dispatch: LOTLInfo.DSSID
// calls LOTLInfo.BuildIdentifier, never PivotInfo's. This port therefore defines
// BuildIdentifier/DSSID/DSSIDAsString/IsPivot as ordinary methods that shadow the embedded
// LOTLInfo ones, so any call made on a *PivotInfo yields the PivotIdentifier ("P-" prefix)
// exactly as Java's PivotInfo#getDSSId() does. The one remaining limitation is a call made
// through the embedded value - &pivotInfo.LOTLInfo, or a *LOTLInfo the caller derived from it -
// which yields LOTLInfo's identifier; hold the *PivotInfo itself (as
// LOTLInfo.PivotInfos() returns) to get Java's behaviour. This is the same limitation
// tl_info.go documents.
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

// DSSID returns the Identifier of the object, computing and caching it on first access.
// Overrides (shadows) the embedded LOTLInfo.DSSID so the cached value is built from
// PivotInfo's own BuildIdentifier: a PivotIdentifier, as Java's PivotInfo#getDSSId() returns
// through its buildIdentifier() override, not the LOTLIdentifier the promoted LOTLInfo.DSSID
// would build.
func (p *PivotInfo) DSSID() model.Identifier {
	if p.identifier == nil {
		p.identifier = p.BuildIdentifier()
	}
	return p.identifier
}

// DSSIDAsString returns the String representation of the identifier.
func (p *PivotInfo) DSSIDAsString() string {
	return p.DSSID().AsXmlID()
}
