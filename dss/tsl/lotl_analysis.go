// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/runnable/LOTLAnalysis.java (DSS 6.5.RC1).
package tsl

import (
	"sync"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// LOTLAnalysis runs the job for a LOTL analysis.
type LOTLAnalysis struct {
	AbstractRunnableTLAnalysis

	// source is the LOTLSource being analyzed.
	source *LOTLSource
}

var _ job.Runnable = (*LOTLAnalysis)(nil)
var _ job.AbstractRunnableAnalysisOverrides = (*LOTLAnalysis)(nil)

// NewLOTLAnalysis is the default constructor.
func NewLOTLAnalysis(source *LOTLSource, cacheAccess job.CacheAccessByKey, dssFileLoader http.DSSFileLoader, latch *sync.WaitGroup) *LOTLAnalysis {
	a := &LOTLAnalysis{
		AbstractRunnableTLAnalysis: InitAbstractRunnableTLAnalysis(&source.TLSource, cacheAccess, dssFileLoader, latch),
		source:                     source,
	}
	a.InitAbstractRunnableAnalysis(a)
	return a
}

// GetParsingTask ports the protected getParsingTask(DSSDocument) override. Wrapped in
// parsingTaskAdapter (abstract_runnable_tl_analysis.go) - see that type's header.
func (a *LOTLAnalysis) GetParsingTask(document model.DSSDocument) job.ParsingTask {
	return parsingTaskAdapter[*LOTLParsingResult]{get: NewLOTLParsingTask(document, a.source).Get}
}

// GetCurrentCertificateSource satisfies job.AbstractRunnableAnalysisOverrides with the base
// (non-overridden in Java, for LOTLAnalysis itself - LOTLWithPivotsAnalysis overrides it below)
// behaviour.
func (a *LOTLAnalysis) GetCurrentCertificateSource() spi.CertificateSource {
	return a.GetCurrentCertificateSourceDefault()
}
