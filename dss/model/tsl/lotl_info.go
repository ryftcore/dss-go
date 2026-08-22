// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/LOTLInfo.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/job"
)

// LOTLInfo computes a summary for a List of Trusted Lists processing result.
//
// Implements the assumed job.DocumentListInfo[LOTLInfo, TLInfo] interface (getChildrenInfos,
// ported as ChildrenInfos()) on top of the fields/behaviour TLInfo already provides for
// job.DocumentInfo. See the JUDGMENT CALL note on TLInfo for the AbstractDocumentInfo
// flattening this and PivotInfo (already ported, see pivot_info.go) depend on.
type LOTLInfo struct {
	TLInfo

	// tlInfos is the list of summaries for TLs found inside the current LOTL.
	tlInfos []*TLInfo
	// pivotInfos is the list of summaries for pivots found inside the current LOTL.
	pivotInfos []*PivotInfo
	// identifier caches the value returned by DSSID, shadowing the embedded TLInfo's own cache
	// so LOTLInfo's BuildIdentifier (not TLInfo's) is what gets cached and returned.
	identifier model.Identifier
}

// NewLOTLInfo is the default constructor.
func NewLOTLInfo(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo TLParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string) LOTLInfo {
	return LOTLInfo{
		TLInfo: *NewTLInfo(downloadCacheInfo, parsingCacheInfo, validationCacheInfo, url),
	}
}

// TLInfos returns a list of TLInfo summaries for TLs found in the LOTL.
func (l *LOTLInfo) TLInfos() []*TLInfo {
	return l.tlInfos
}

// SetTlInfos sets a list of TLInfo summaries for TLs found in the LOTL.
func (l *LOTLInfo) SetTlInfos(tlInfos []*TLInfo) {
	l.tlInfos = tlInfos
}

// PivotInfos returns a list of PivotInfo summaries for pivots found in the LOTL.
func (l *LOTLInfo) PivotInfos() []*PivotInfo {
	return l.pivotInfos
}

// SetPivotInfos sets a list of PivotInfo summaries for pivots found in the LOTL.
func (l *LOTLInfo) SetPivotInfos(pivotInfos []*PivotInfo) {
	l.pivotInfos = pivotInfos
}

// IsPivot checks if the current entry is a pivot info: a LOTLInfo never is.
func (l *LOTLInfo) IsPivot() bool {
	return false
}

// BuildIdentifier builds the identifier of the current LOTL. Overrides (shadows) the embedded
// TLInfo.BuildIdentifier for direct calls on *LOTLInfo.
func (l *LOTLInfo) BuildIdentifier() model.Identifier {
	return NewLOTLIdentifier(l)
}

// DSSID returns the Identifier of the object, computing and caching it on first access.
// Overrides (shadows) the embedded TLInfo.DSSID so the cached value is built from LOTLInfo's
// own BuildIdentifier.
func (l *LOTLInfo) DSSID() model.Identifier {
	if l.identifier == nil {
		l.identifier = l.BuildIdentifier()
	}
	return l.identifier
}

// DSSIDAsString returns the String representation of the identifier.
func (l *LOTLInfo) DSSIDAsString() string {
	return l.DSSID().AsXmlID()
}

// ChildrenInfos returns a list of DocumentInfo summaries for documents referenced from the
// current LOTL.
func (l *LOTLInfo) ChildrenInfos() []*TLInfo {
	return l.TLInfos()
}
