// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/TLInfo.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/job"
)

// TLInfo computes a summary for a single Trusted List processing result.
//
// JUDGMENT CALL: Java's AbstractDocumentInfo<P> (package model.job, out of this manifest) is a
// generic abstract base with a self-referential type parameter and a protected
// buildIdentifier() hook invoked polymorphically from its getDSSId(). Go has no virtual
// dispatch through embedding and Go generics do not make self-referential bounds pleasant, so
// this port inlines AbstractDocumentInfo's fields/behaviour directly into TLInfo instead of
// depending on a generic job.AbstractDocumentInfo type. LOTLInfo (same package) embeds TLInfo
// by value and defines its own BuildIdentifier/DSSID pair that shadows these for direct calls
// on a *LOTLInfo; pivot_info.go (already ported) further shadows BuildIdentifier on *PivotInfo
// but does not redefine DSSID, so a caller holding only a *LOTLInfo-shaped or job.DocumentInfo
// value gets LOTLInfo's identifier, not PivotInfo's - the same known deviation pivot_info.go
// documents on its own BuildIdentifier. The integrator should confirm the actual job package
// interfaces (assumed here: DownloadCacheInfo() job.DownloadInfoRecord,
// ParsingCacheInfo() job.ParsingInfoRecord, ValidationCacheInfo() job.ValidationInfoRecord,
// Url() string, Parent() P, DSSID() model.Identifier, DSSIDAsString() string) match what this
// type structurally provides.
type TLInfo struct {
	// downloadCacheInfo is the download result record.
	downloadCacheInfo job.DownloadInfoRecord
	// parsingCacheInfo is the parsing result record.
	parsingCacheInfo TLParsingInfoRecord
	// validationCacheInfo is the validation result record.
	validationCacheInfo job.ValidationInfoRecord
	// url is the address of the source.
	url string
	// parent is the parent LOTL referencing the current Trusted List.
	parent *LOTLInfo
	// otherTSLPointer is the OtherTSLPointer element extracted from the pointing TL/LOTL.
	otherTSLPointer *OtherTSLPointer
	// identifier caches the value returned by DSSID.
	identifier model.Identifier
}

// NewTLInfo is the default constructor.
func NewTLInfo(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo TLParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string) *TLInfo {
	return NewTLInfoWithParent(downloadCacheInfo, parsingCacheInfo, validationCacheInfo, url, nil)
}

// NewTLInfoWithParent is the constructor with a parent LOTLInfo.
func NewTLInfoWithParent(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo TLParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string, parent *LOTLInfo) *TLInfo {
	return NewTLInfoFull(downloadCacheInfo, parsingCacheInfo, validationCacheInfo, url, parent, nil)
}

// NewTLInfoFull is the constructor with a parent LOTLInfo and Mutual Recognition Agreement
// pointer.
func NewTLInfoFull(downloadCacheInfo job.DownloadInfoRecord, parsingCacheInfo TLParsingInfoRecord,
	validationCacheInfo job.ValidationInfoRecord, url string, parent *LOTLInfo,
	otherTSLPointer *OtherTSLPointer) *TLInfo {
	return &TLInfo{
		downloadCacheInfo:   downloadCacheInfo,
		parsingCacheInfo:    parsingCacheInfo,
		validationCacheInfo: validationCacheInfo,
		url:                 url,
		parent:              parent,
		otherTSLPointer:     otherTSLPointer,
	}
}

// DownloadCacheInfo returns the download cache info.
func (t *TLInfo) DownloadCacheInfo() job.DownloadInfoRecord {
	return t.downloadCacheInfo
}

// ParsingCacheInfo returns the parsing cache info. Port of the covariant
// TLInfo#getParsingCacheInfo() override.
//
// INTEGRATION FIX: Java's override narrows the return type to TLParsingInfoRecord
// (covariant return). Go has no covariant interface-method return types, and this method
// must satisfy job.DocumentInfo[P]'s ParsingCacheInfo() job.ParsingInfoRecord exactly for
// *TLInfo/*LOTLInfo to be usable as the D/P type arguments of job.ValidationJob and friends -
// so it answers the base job.ParsingInfoRecord interface here. The underlying value's dynamic
// type is always a TLParsingInfoRecord (e.g. *tsl.TLParsingCacheDTO from package tsl), so a
// caller needing the TL-specific accessors (TSLType, SequenceNumber, NextUpdateDate, ...)
// recovers them with a type assertion: `tlParsingCacheInfo, ok :=
// tlInfo.ParsingCacheInfo().(TLParsingInfoRecord)`.
func (t *TLInfo) ParsingCacheInfo() job.ParsingInfoRecord {
	return t.parsingCacheInfo
}

// TLParsingCacheInfo returns the parsing cache info narrowed to TLParsingInfoRecord, when the
// stored record actually carries the TL-specific accessors (it always does in practice - the
// only two- implementations produced anywhere in this tree are tsl.TLParsingCacheDTO and a nil
// interface). ok is false when parsingCacheInfo is nil or does not implement
// TLParsingInfoRecord. Convenience wrapper over the type assertion documented on
// ParsingCacheInfo, so callers do not need to repeat it inline.
func (t *TLInfo) TLParsingCacheInfo() (TLParsingInfoRecord, bool) {
	tlParsingCacheInfo, ok := t.parsingCacheInfo.(TLParsingInfoRecord)
	return tlParsingCacheInfo, ok
}

// ValidationCacheInfo returns the validation cache info.
func (t *TLInfo) ValidationCacheInfo() job.ValidationInfoRecord {
	return t.validationCacheInfo
}

// Url returns the URL used to download the remote file.
func (t *TLInfo) Url() string {
	return t.url
}

// Parent returns the LOTLInfo referencing the current Trusted List.
func (t *TLInfo) Parent() *LOTLInfo {
	return t.parent
}

// OtherTSLPointer gets the OtherTSLPointer element referencing the current TL from the
// pointing TL/LOTL.
func (t *TLInfo) OtherTSLPointer() *OtherTSLPointer {
	return t.otherTSLPointer
}

// BuildIdentifier builds the identifier of the current TL. Port of the protected
// buildIdentifier(). Exported so LOTLInfo (and further embedders) can shadow it - see the
// JUDGMENT CALL note on TLInfo about the resulting virtual-dispatch limitation.
func (t *TLInfo) BuildIdentifier() model.Identifier {
	return NewTrustedListIdentifier(t)
}

// DSSID returns the Identifier of the object, computing and caching it on first access. Port
// of AbstractDocumentInfo#getDSSId().
func (t *TLInfo) DSSID() model.Identifier {
	if t.identifier == nil {
		t.identifier = t.BuildIdentifier()
	}
	return t.identifier
}

// DSSIDAsString returns the String representation of the identifier.
func (t *TLInfo) DSSIDAsString() string {
	return t.DSSID().AsXmlID()
}
