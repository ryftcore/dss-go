// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/runnable/AbstractRunnableAnalysis.java (DSS 6.5.RC1).
package job

import (
	"sync"

	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/client/http"
)

// AbstractRunnableAnalysisOverrides extends AbstractAnalysisOverrides with the operation
// Java's AbstractRunnableAnalysis class declares as overridable (though not abstract) -
// getCurrentCertificateSource() - standing in for the virtual dispatch doAnalyze needs to
// reach a concrete override. A concrete analysis registers itself with
// AbstractRunnableAnalysis.InitAbstractRunnableAnalysis.
type AbstractRunnableAnalysisOverrides interface {
	AbstractAnalysisOverrides

	// GetCurrentCertificateSource returns the certificate source to be used to validate the
	// document. Port of the (overridable, non-abstract) protected
	// getCurrentCertificateSource(); a concrete analysis wanting the base behaviour calls
	// AbstractRunnableAnalysis.GetCurrentCertificateSourceDefault explicitly, mirroring
	// Java's super.getCurrentCertificateSource().
	GetCurrentCertificateSource() spi.CertificateSource
}

// abstractRunnableAnalysisLogErrorPerformAnalysis mirrors LOG_ERROR_PERFORM_ANALYSIS; kept
// only as documentation of the dropped log message (slf4j logging is not ported).
const abstractRunnableAnalysisLogErrorPerformAnalysis = "Error performing analysis."

// AbstractRunnableAnalysis is the runnable facade that processes the document validation job
// (download - parse - validate).
type AbstractRunnableAnalysis struct {
	AbstractAnalysis

	// overrides points back at the concrete analysis; see InitAbstractRunnableAnalysis.
	overrides AbstractRunnableAnalysisOverrides

	// latch is the tasks counter (Java's CountDownLatch), ported as a *sync.WaitGroup whose
	// Done is called exactly once by Run's finally-equivalent, mirroring latch.countDown().
	latch *sync.WaitGroup
}

// NewAbstractRunnableAnalysis creates an AbstractRunnableAnalysis. Port of the protected
// constructor.
func NewAbstractRunnableAnalysis(source *DocumentSource, cacheAccess CacheAccessByKey,
	dssFileLoader http.DSSFileLoader, latch *sync.WaitGroup) AbstractRunnableAnalysis {
	return AbstractRunnableAnalysis{
		AbstractAnalysis: NewAbstractAnalysis(source, cacheAccess, dssFileLoader),
		latch:            latch,
	}
}

// InitAbstractRunnableAnalysis registers the concrete analysis with its base (and with the
// embedded AbstractAnalysis base) so that self-calls dispatch to the concrete
// implementation. It must be called exactly once, by the concrete analysis's constructor,
// before any other method.
func (a *AbstractRunnableAnalysis) InitAbstractRunnableAnalysis(overrides AbstractRunnableAnalysisOverrides) {
	a.overrides = overrides
	a.AbstractAnalysis.InitAbstractAnalysis(overrides)
}

func (a *AbstractRunnableAnalysis) abstractRunnableAnalysisOverrides() AbstractRunnableAnalysisOverrides {
	if a.overrides == nil {
		panic("AbstractRunnableAnalysis was not initialised: the concrete analysis must call InitAbstractRunnableAnalysis in its constructor")
	}
	return a.overrides
}

// DoAnalyze performs the analysis: download, then (if a document was obtained) parse and
// validate. Port of doAnalyze().
func (a *AbstractRunnableAnalysis) DoAnalyze() {
	document := a.download(a.Source().Url())
	if document != nil {
		a.parsing(document)
		a.validation(document, a.abstractRunnableAnalysisOverrides().GetCurrentCertificateSource())
	}
}

// GetCurrentCertificateSourceDefault is the base (non-abstract) implementation of the
// GetCurrentCertificateSource hook. Port of the protected getCurrentCertificateSource()
// default body; an overriding concrete analysis calls this explicitly to get the base
// behaviour, mirroring Java's super.getCurrentCertificateSource().
func (a *AbstractRunnableAnalysis) GetCurrentCertificateSourceDefault() spi.CertificateSource {
	return a.Source().CertificateSource()
}

// Run performs the analysis, recovering from any panic (mirroring Java's catch(Throwable)),
// and always signals the latch afterward (mirroring the finally block's
// latch.countDown()).
//
// DEVIATION: slf4j warn logging of the caught Throwable is dropped, per this module's
// logging convention; see abstractRunnableAnalysisLogErrorPerformAnalysis.
func (a *AbstractRunnableAnalysis) Run() {
	defer a.latch.Done()
	defer func() {
		_ = recover()
	}()
	a.DoAnalyze()
}
