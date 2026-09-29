// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/runnable/AbstractAnalysis.java (DSS 6.5.RC1).
package job

import (
	"fmt"

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
// Java's download()/parsing()/validation() wrap the whole body (task construction, get(), and
// every subsequent cache call) in a single try/catch(Exception) that records the exception with
// the corresponding *Error cache method. The Go port reproduces that span: an error returned by
// the task's Get() (see DownloadTask, ParsingTask, ValidationTask) and a panic raised anywhere
// inside it - Go's unchecked exception, e.g. the nil dereference standing in for a
// NullPointerException on a malformed document - are both routed to the *Error cache method (see
// catchException). Without the panic half, such a document would abort the analysis without any
// error being recorded (and, inside a goroutine that has no recover of its own, crash the
// process).
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

// catchException runs fn and returns the error it returns or, when it panics, an error built
// from the panic value. It stands in for Java's catch (Exception e) in download()/parsing()/
// validation(): Java's unchecked exceptions (NullPointerException, IllegalStateException, ...)
// are Go panics, and Java's checked ones are the error fn returns. A panic value that is already
// an error (a *model.DSSError, a runtime.Error) is passed through unchanged; anything else
// (the ported Objects.requireNonNull messages are plain strings) becomes a *model.DSSError with
// that message.
func catchException(fn func() error) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if recoveredErr, ok := recovered.(error); ok {
				err = recoveredErr
			} else {
				err = model.NewDSSError(fmt.Sprint(recovered))
			}
		}
	}()
	return fn()
}

// download downloads the document by url. Port of download(String).
func (a *AbstractAnalysis) download(url string) model.DSSDocument {
	overrides := a.abstractAnalysisOverrides()
	var document model.DSSDocument
	if err := catchException(func() error {
		downloadTask := overrides.GetDownloadTask(a.dssFileLoader, url)
		downloadResult, err := downloadTask.Get()
		if err != nil {
			return err
		}
		if !a.cacheAccess.IsUpToDate(downloadResult) {
			a.cacheAccess.UpdateDownloadResult(downloadResult)
			a.expireCache()
		}
		document = downloadResult.DSSDocument()
		return nil
	}); err != nil {
		a.cacheAccess.DownloadError(err)
		return nil
	}
	return document
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
		overrides := a.abstractAnalysisOverrides()
		if err := catchException(func() error {
			parsingTask := overrides.GetParsingTask(document)
			parsingResult, err := parsingTask.Get()
			if err != nil {
				return err
			}
			a.cacheAccess.UpdateParsingResult(parsingResult)
			return nil
		}); err != nil {
			a.cacheAccess.ParsingError(err)
		}
	}
}

// validation validates the document. Port of validation(DSSDocument, CertificateSource).
func (a *AbstractAnalysis) validation(document model.DSSDocument, certificateSource spi.CertificateSource) {
	// True if EMPTY / EXPIRED by document
	if a.cacheAccess.IsValidationRefreshNeeded() {
		overrides := a.abstractAnalysisOverrides()
		if err := catchException(func() error {
			validationTask := overrides.GetValidationTask(document, certificateSource)
			validationResult, err := validationTask.Get()
			if err != nil {
				return err
			}
			a.cacheAccess.UpdateValidationResult(validationResult)
			return nil
		}); err != nil {
			a.cacheAccess.ValidationError(err)
		}
	}
}
