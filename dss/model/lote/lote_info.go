// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/LoTEInfo.java (DSS 6.5.RC1).
package lote

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/job"
)

// Info computes a summary for a single TS 119 602 List of Trusted Entities processing
// result.
//
// JUDGMENT CALL: mirrors the AbstractDocumentInfo flattening already documented on
// tsl.TLInfo/tsl.LOTLInfo (same reasoning applies here): Java's
// model.job.AbstractDocumentInfo<P> (out of this manifest) is a generic abstract base with a
// self-referential type parameter and a protected buildIdentifier() hook invoked
// polymorphically from its getDSSId(). Go has no virtual dispatch through embedding and Go
// generics do not make self-referential bounds pleasant, so this port inlines
// AbstractDocumentInfo's fields/behaviour directly into Info instead of depending on a
// generic job.AbstractDocumentInfo type. LoLoTEInfo (same package) embeds Info by value
// and defines its own BuildIdentifier/DSSID pair that shadows these for direct calls on a
// *LoLoTEInfo. The integrator should confirm the actual job package interfaces (assumed here:
// job.DownloadInfoRecord, job.ParsingInfoRecord, job.ValidationInfoRecord) match what this type
// structurally provides.
type Info struct {
	// downloadCacheInfo is the download result record.
	downloadCacheInfo job.DownloadInfoRecord
	// parsingCacheInfo is the parsing result record.
	parsingCacheInfo ParsingInfoRecord
	// validationCacheInfo is the validation result record.
	validationCacheInfo job.ValidationInfoRecord
	// url is the address used to extract the entry.
	url string
	// parent is the parent LoLoTEInfo referencing the current List.
	parent *LoLoTEInfo
	// otherListPointer is the OtherListPointer element extracted from the pointing LoLoTE/LoTE.
	otherListPointer *OtherListPointer
	// identifier caches the value returned by DSSID.
	identifier model.Identifier
}

// NewInfo is the default constructor.
func NewInfo(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo ParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string) *Info {
	return NewInfoWithParent(downloadCacheInfo, parsingCacheInfo, validationCacheInfo, url, nil)
}

// NewInfoWithParent is the constructor with a parent LoLoTEInfo.
func NewInfoWithParent(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo ParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string, parent *LoLoTEInfo) *Info {
	return NewInfoFull(downloadCacheInfo, parsingCacheInfo, validationCacheInfo, url, parent, nil)
}

// NewInfoFull is the constructor with a parent LoLoTEInfo and OtherListPointer.
func NewInfoFull(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo ParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string, parent *LoLoTEInfo,
	otherListPointer *OtherListPointer) *Info {
	return &Info{
		downloadCacheInfo:   downloadCacheInfo,
		parsingCacheInfo:    parsingCacheInfo,
		validationCacheInfo: validationCacheInfo,
		url:                 url,
		parent:              parent,
		otherListPointer:    otherListPointer,
	}
}

// DownloadCacheInfo returns the download cache info.
func (l *Info) DownloadCacheInfo() job.DownloadInfoRecord {
	return l.downloadCacheInfo
}

// ParsingCacheInfo returns the parsing cache info. Port of the covariant
// LoTEInfo#getParsingCacheInfo() override, narrowed to LoTEParsingInfoRecord.
func (l *Info) ParsingCacheInfo() ParsingInfoRecord {
	return l.parsingCacheInfo
}

// ValidationCacheInfo returns the validation cache info.
func (l *Info) ValidationCacheInfo() job.ValidationInfoRecord {
	return l.validationCacheInfo
}

// Url returns the URL used to download the remote file.
func (l *Info) Url() string {
	return l.url
}

// Parent returns the LoLoTEInfo referencing the current List.
func (l *Info) Parent() *LoLoTEInfo {
	return l.parent
}

// ListPointer gets the pointer to the current LoTE.
func (l *Info) ListPointer() *OtherListPointer {
	return l.otherListPointer
}

// BuildIdentifier builds the identifier of the current LoTE. Port of the protected
// buildIdentifier(). Exported so LoLoTEInfo can shadow it - see the JUDGMENT CALL note above
// about the resulting virtual-dispatch limitation.
func (l *Info) BuildIdentifier() model.Identifier {
	return NewIdentifier(l)
}

// DSSID returns the Identifier of the object, computing and caching it on first access. Port
// of LoTEInfo#getDSSId() (originally declared on AbstractDocumentInfo, see the JUDGMENT CALL
// note above).
func (l *Info) DSSID() model.Identifier {
	if l.identifier == nil {
		l.identifier = l.BuildIdentifier()
	}
	return l.identifier
}

// DSSIDAsString returns the String representation of the identifier.
func (l *Info) DSSIDAsString() string {
	return l.DSSID().AsXmlID()
}
