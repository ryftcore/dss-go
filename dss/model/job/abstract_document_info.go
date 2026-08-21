// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/AbstractDocumentInfo.java (DSS 6.5.RC1).
package job

import "github.com/ryftcore/dss-go/dss/model"

// AbstractDocumentInfoOverrides declares the operation Java's abstract AbstractDocumentInfo
// class leaves abstract, standing in for the virtual dispatch the base needs to reach the
// concrete document info. A concrete document info registers itself with
// AbstractDocumentInfoBase.InitAbstractDocumentInfo.
type AbstractDocumentInfoOverrides interface {
	// BuildIdentifier builds the identifier. Port of the abstract protected
	// buildIdentifier().
	BuildIdentifier() model.Identifier
}

// AbstractDocumentInfoBase carries the state and behaviour of the abstract Java class
// AbstractDocumentInfo<P>: an abstract representation of a document validation result,
// containing information about the download, parsing and validation statuses. Concrete
// document infos embed it and register themselves with InitAbstractDocumentInfo.
//
// P is the parent DocumentInfo type, mirroring Java's self-bound "P extends DocumentInfo<P>".
type AbstractDocumentInfoBase[P any] struct {
	// overrides points back at the concrete document info; see InitAbstractDocumentInfo.
	overrides AbstractDocumentInfoOverrides

	// url is the address of the source.
	url string
	// parent is the parent LOTL/TL referencing the current Trusted List. The zero value of P
	// stands for Java's null when P is an interface type, matching the single-argument
	// constructor.
	parent P

	// downloadCacheInfo is the download result record.
	downloadCacheInfo DownloadInfoRecord
	// parsingCacheInfo is the parsing result record.
	parsingCacheInfo ParsingInfoRecord
	// validationCacheInfo is the validation result record.
	validationCacheInfo ValidationInfoRecord

	// identifier caches the built Identifier.
	identifier model.Identifier
}

// NewAbstractDocumentInfoBase is the default constructor without a parent TLInfo. Port of
// AbstractDocumentInfo(DownloadInfoRecord, ParsingInfoRecord, ValidationInfoRecord, String).
//
// Panics with the Java message when url is empty (Java Objects.requireNonNull("URL String
// shall be provided!")); Go strings cannot be null, so the empty string is the port of the
// missing-value case.
func NewAbstractDocumentInfoBase[P any](downloadCacheInfo DownloadInfoRecord, parsingCacheInfo ParsingInfoRecord,
	validationCacheInfo ValidationInfoRecord, url string) AbstractDocumentInfoBase[P] {
	var zero P
	return NewAbstractDocumentInfoBaseWithParent(downloadCacheInfo, parsingCacheInfo, validationCacheInfo, url, zero)
}

// NewAbstractDocumentInfoBaseWithParent is the default constructor with a parent TLInfo. Port
// of AbstractDocumentInfo(DownloadInfoRecord, ParsingInfoRecord, ValidationInfoRecord,
// String, P).
//
// Panics with the Java message when url is empty (Java Objects.requireNonNull("URL String
// shall be provided!")); Go strings cannot be null, so the empty string is the port of the
// missing-value case.
func NewAbstractDocumentInfoBaseWithParent[P any](downloadCacheInfo DownloadInfoRecord, parsingCacheInfo ParsingInfoRecord,
	validationCacheInfo ValidationInfoRecord, url string, parent P) AbstractDocumentInfoBase[P] {
	if url == "" {
		panic("URL String shall be provided!")
	}
	return AbstractDocumentInfoBase[P]{
		downloadCacheInfo:   downloadCacheInfo,
		parsingCacheInfo:    parsingCacheInfo,
		validationCacheInfo: validationCacheInfo,
		url:                 url,
		parent:              parent,
	}
}

// InitAbstractDocumentInfo registers the concrete document info with its base so that the
// base can dispatch to BuildIdentifier. It must be called exactly once, by the concrete
// document info's constructor, before any other method.
func (a *AbstractDocumentInfoBase[P]) InitAbstractDocumentInfo(overrides AbstractDocumentInfoOverrides) {
	a.overrides = overrides
}

// abstractDocumentInfoOverrides returns the registered overrides, panicking when the
// concrete document info forgot to call InitAbstractDocumentInfo.
func (a *AbstractDocumentInfoBase[P]) abstractDocumentInfoOverrides() AbstractDocumentInfoOverrides {
	if a.overrides == nil {
		panic("AbstractDocumentInfoBase was not initialised: the concrete document info must call InitAbstractDocumentInfo in its constructor")
	}
	return a.overrides
}

// DownloadCacheInfo returns the download result record. Port of getDownloadCacheInfo().
func (a *AbstractDocumentInfoBase[P]) DownloadCacheInfo() DownloadInfoRecord {
	return a.downloadCacheInfo
}

// ParsingCacheInfo returns the parsing result record. Port of getParsingCacheInfo().
func (a *AbstractDocumentInfoBase[P]) ParsingCacheInfo() ParsingInfoRecord {
	return a.parsingCacheInfo
}

// ValidationCacheInfo returns the validation result record. Port of
// getValidationCacheInfo().
func (a *AbstractDocumentInfoBase[P]) ValidationCacheInfo() ValidationInfoRecord {
	return a.validationCacheInfo
}

// URL returns the address used to extract the entry. Port of getUrl().
func (a *AbstractDocumentInfoBase[P]) URL() string {
	return a.url
}

// Parent returns the DocumentInfo referencing the current Trusted List. Port of getParent().
func (a *AbstractDocumentInfoBase[P]) Parent() P {
	return a.parent
}

// DSSID returns the Identifier, building it lazily via the registered overrides. Port of
// getDSSId().
func (a *AbstractDocumentInfoBase[P]) DSSID() model.Identifier {
	if a.identifier == nil {
		a.identifier = a.abstractDocumentInfoOverrides().BuildIdentifier()
	}
	return a.identifier
}

// DSSIDAsString returns the String representation of the identifier. Port of
// getDSSIdAsString().
func (a *AbstractDocumentInfoBase[P]) DSSIDAsString() string {
	return a.DSSID().AsXmlID()
}
