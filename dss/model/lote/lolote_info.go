// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/LoLoTEInfo.java (DSS 6.5.RC1).
package lote

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/job"
)

// LoLoTEInfo computes a summary for a TS 119 602 List of Lists of Trusted Entities processing
// result.
//
// Implements the assumed job.DocumentListInfo[LoLoTEInfo, LoTEInfo] interface
// (getChildrenInfos, ported as ChildrenInfos()) on top of the fields/behaviour LoTEInfo already
// provides for job.DocumentInfo. See the JUDGMENT CALL note on LoTEInfo for the
// AbstractDocumentInfo flattening this depends on.
type LoLoTEInfo struct {
	LoTEInfo

	// childrenInfos is the list of summary for Lists found inside the current LoTE.
	childrenInfos []*LoTEInfo
	// identifier caches the value returned by DSSID, shadowing the embedded LoTEInfo's own
	// cache so LoLoTEInfo's BuildIdentifier (not LoTEInfo's) is what gets cached and returned.
	identifier model.Identifier
}

// NewLoLoTEInfo is the default constructor.
func NewLoLoTEInfo(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo LoTEParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string) *LoLoTEInfo {
	return &LoLoTEInfo{
		LoTEInfo: *NewLoTEInfo(downloadCacheInfo, parsingCacheInfo, validationCacheInfo, url),
	}
}

// ChildrenInfos gets a list of processing information for other referenced LoTEs.
func (l *LoLoTEInfo) ChildrenInfos() []*LoTEInfo {
	return l.childrenInfos
}

// SetChildrenInfos sets a list of LoTEInfo summary for LoTE found in the LoLoTE.
func (l *LoLoTEInfo) SetChildrenInfos(childrenInfos []*LoTEInfo) {
	l.childrenInfos = childrenInfos
}

// BuildIdentifier builds the identifier of the current LoLoTE. Overrides (shadows) the
// embedded LoTEInfo.BuildIdentifier for direct calls on *LoLoTEInfo.
func (l *LoLoTEInfo) BuildIdentifier() model.Identifier {
	return NewLoLoTEIdentifier(&l.LoTEInfo)
}

// DSSID returns the Identifier of the object, computing and caching it on first access.
// Overrides (shadows) the embedded LoTEInfo.DSSID so the cached value is built from
// LoLoTEInfo's own BuildIdentifier.
func (l *LoLoTEInfo) DSSID() model.Identifier {
	if l.identifier == nil {
		l.identifier = l.BuildIdentifier()
	}
	return l.identifier
}

// DSSIDAsString returns the String representation of the identifier.
func (l *LoLoTEInfo) DSSIDAsString() string {
	return l.DSSID().AsXmlID()
}
