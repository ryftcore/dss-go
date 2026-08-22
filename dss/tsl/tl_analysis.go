// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/runnable/TLAnalysis.java (DSS 6.5.RC1).
package tsl

import (
	"sync"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// TLAnalysis runs the job for a TL analysis.
type TLAnalysis struct {
	AbstractRunnableTLAnalysis

	// source is the TLSource being analyzed.
	source *TLSource
}

var _ job.Runnable = (*TLAnalysis)(nil)
var _ job.AbstractRunnableAnalysisOverrides = (*TLAnalysis)(nil)

// NewTLAnalysis is the default constructor.
func NewTLAnalysis(source *TLSource, cacheAccess job.CacheAccessByKey, dssFileLoader http.DSSFileLoader, latch *sync.WaitGroup) *TLAnalysis {
	a := &TLAnalysis{
		AbstractRunnableTLAnalysis: InitAbstractRunnableTLAnalysis(source, cacheAccess, dssFileLoader, latch),
		source:                     source,
	}
	a.InitAbstractRunnableAnalysis(a)
	return a
}

// GetParsingTask ports the protected getParsingTask(DSSDocument) override. Wrapped in
// parsingTaskAdapter (abstract_runnable_tl_analysis.go) - see that type's header.
func (a *TLAnalysis) GetParsingTask(document model.DSSDocument) job.ParsingTask {
	return parsingTaskAdapter[*TLParsingResult]{get: NewTLParsingTask(document, a.source).Get}
}

// GetCurrentCertificateSource satisfies job.AbstractRunnableAnalysisOverrides with the base
// (non-overridden in Java) behaviour, mirroring the implicit inheritance of
// AbstractRunnableAnalysis.getCurrentCertificateSource().
func (a *TLAnalysis) GetCurrentCertificateSource() spi.CertificateSource {
	return a.GetCurrentCertificateSourceDefault()
}
