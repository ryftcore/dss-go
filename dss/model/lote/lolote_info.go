// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/LoLoTEInfo.java (DSS 6.5.RC1).
package lote

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/job"
)

// LoLoTEInfo computes a summary for a TS 119 602 List of Lists of Trusted Entities processing
// result.
//
// Implements the assumed job.DocumentListInfo[LoLoTEInfo, Info] interface
// (getChildrenInfos, ported as ChildrenInfos()) on top of the fields/behaviour Info already
// provides for job.DocumentInfo. See the JUDGMENT CALL note on Info for the
// AbstractDocumentInfo flattening this depends on.
type LoLoTEInfo struct {
	Info

	// childrenInfos is the list of summary for Lists found inside the current LoTE.
	childrenInfos []*Info
	// identifier caches the value returned by DSSID, shadowing the embedded Info's own
	// cache so LoLoTEInfo's BuildIdentifier (not Info's) is what gets cached and returned.
	identifier model.Identifier
}

// NewLoLoTEInfo is the default constructor.
func NewLoLoTEInfo(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo ParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string) *LoLoTEInfo {
	return &LoLoTEInfo{
		Info: *NewLoTEInfo(downloadCacheInfo, parsingCacheInfo, validationCacheInfo, url),
	}
}

// ChildrenInfos gets a list of processing information for other referenced LoTEs.
func (l *LoLoTEInfo) ChildrenInfos() []*Info {
	return l.childrenInfos
}

// SetChildrenInfos sets a list of Info summary for LoTE found in the LoLoTE.
func (l *LoLoTEInfo) SetChildrenInfos(childrenInfos []*Info) {
	l.childrenInfos = childrenInfos
}

// BuildIdentifier builds the identifier of the current LoLoTE. Overrides (shadows) the
// embedded Info.BuildIdentifier for direct calls on *LoLoTEInfo.
func (l *LoLoTEInfo) BuildIdentifier() model.Identifier {
	return NewLoLoTEIdentifier(&l.Info)
}

// DSSID returns the Identifier of the object, computing and caching it on first access.
// Overrides (shadows) the embedded Info.DSSID so the cached value is built from
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
