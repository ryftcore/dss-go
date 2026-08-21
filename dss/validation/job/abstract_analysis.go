// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/runnable/AbstractAnalysis.java (DSS 6.5.RC1).
package job

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
)

// AbstractAnalysisOverrides declares the operations Java's abstract AbstractAnalysis class
// leaves abstract, standing in for the virtual dispatch the base needs to reach the concrete
// analysis. A concrete analysis registers itself with AbstractAnalysis.InitAbstractAnalysis.
type AbstractAnalysisOverrides interface {
	// GetDownloadTask returns the corresponding download task for the source on the given
	// document. Port of the abstract protected getDownloadTask(DSSFileLoader, String).
	GetDownloadTask(dssFileLoader http.DSSFileLoader, url string) DownloadTask

	// GetParsingTask returns the corresponding parsing task for the source on the given
	// document. Port of the abstract protected getParsingTask(DSSDocument).
	GetParsingTask(document model.DSSDocument) ParsingTask

	// GetValidationTask returns the corresponding validation task for the source on the
	// given document using the provided certificate source. Port of the abstract protected
	// getValidationTask(DSSDocument, CertificateSource).
	GetValidationTask(document model.DSSDocument, certificateSource spi.CertificateSource) ValidationTask
}

// AbstractAnalysis processes the document/document list validation job (download - parse -
// validate).
//
// slf4j debug/warn logging is dropped (no observable behavior).
//
// DEVIATION: Java's download()/parsing()/validation() wrap the whole body (task
// construction, get(), and every subsequent cache call) in a single try/catch(Exception).
// The Go port only routes the error returned by the task's Get() (see DownloadTask,
// ParsingTask, ValidationTask) to the corresponding *Error cache method; the surrounding
// cache-access calls are ordinary (non-erroring) Go method calls, not further guarded.
type AbstractAnalysis struct {
	// overrides points back at the concrete analysis; see InitAbstractAnalysis.
	overrides AbstractAnalysisOverrides

	// source is the document/document list source.
	source *DocumentSource

	// cacheAccess is the cache access of the record.
	cacheAccess CacheAccessByKey

	// dssFileLoader is the file loader.
	dssFileLoader http.DSSFileLoader
}

// NewAbstractAnalysis creates an AbstractAnalysis. Port of the protected constructor.
func NewAbstractAnalysis(source *DocumentSource, cacheAccess CacheAccessByKey, dssFileLoader http.DSSFileLoader) AbstractAnalysis {
	return AbstractAnalysis{source: source, cacheAccess: cacheAccess, dssFileLoader: dssFileLoader}
}

// InitAbstractAnalysis registers the concrete analysis with its base so that the base can
// dispatch to GetDownloadTask/GetParsingTask/GetValidationTask. It must be called exactly
// once, by the concrete analysis's constructor, before any other method.
func (a *AbstractAnalysis) InitAbstractAnalysis(overrides AbstractAnalysisOverrides) {
	a.overrides = overrides
}

func (a *AbstractAnalysis) abstractAnalysisOverrides() AbstractAnalysisOverrides {
	if a.overrides == nil {
		panic("AbstractAnalysis was not initialised: the concrete analysis must call InitAbstractAnalysis in its constructor")
	}
	return a.overrides
}

// Source returns the current DocumentSource. Port of getSource().
func (a *AbstractAnalysis) Source() *DocumentSource {
	return a.source
}

// CacheAccessByKey returns the ReadOnlyCacheAccessByKey view of the analysis's cache
// access. Port of the protected final getCacheAccessByKey().
func (a *AbstractAnalysis) CacheAccessByKey() ReadOnlyCacheAccessByKey {
	return a.cacheAccess
}

// download downloads the document by url. Port of download(String).
func (a *AbstractAnalysis) download(url string) model.DSSDocument {
	downloadTask := a.abstractAnalysisOverrides().GetDownloadTask(a.dssFileLoader, url)
	downloadResult, err := downloadTask.Get()
	if err != nil {
		a.cacheAccess.DownloadError(err)
		return nil
	}
	if !a.cacheAccess.IsUpToDate(downloadResult) {
		a.cacheAccess.UpdateDownloadResult(downloadResult)
		a.expireCache()
	}
	return downloadResult.DSSDocument()
}

// expireCache expires the cache in order to trigger the corresponding tasks on refresh.
// Port of expireCache().
func (a *AbstractAnalysis) expireCache() {
	a.cacheAccess.ExpireParsing()
	a.cacheAccess.ExpireValidation()
}

// parsing parses the document. Port of parsing(DSSDocument).
func (a *AbstractAnalysis) parsing(document model.DSSDocument) {
	// True if EMPTY / EXPIRED by a document / document list
	if a.cacheAccess.IsParsingRefreshNeeded() {
		parsingTask := a.abstractAnalysisOverrides().GetParsingTask(document)
		parsingResult, err := parsingTask.Get()
		if err != nil {
			a.cacheAccess.ParsingError(err)
			return
		}
		a.cacheAccess.UpdateParsingResult(parsingResult)
	}
}

// validation validates the document. Port of validation(DSSDocument, CertificateSource).
func (a *AbstractAnalysis) validation(document model.DSSDocument, certificateSource spi.CertificateSource) {
	// True if EMPTY / EXPIRED by document
	if a.cacheAccess.IsValidationRefreshNeeded() {
		validationTask := a.abstractAnalysisOverrides().GetValidationTask(document, certificateSource)
		validationResult, err := validationTask.Get()
		if err != nil {
			a.cacheAccess.ValidationError(err)
			return
		}
		a.cacheAccess.UpdateValidationResult(validationResult)
	}
}
