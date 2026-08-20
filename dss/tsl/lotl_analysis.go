// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/runnable/LOTLAnalysis.java (DSS 6.5.RC1).
package tsl

import (
	"sync"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/client/http"
	"github.com/utain/esig/dss/validation/job"
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

// GetParsingTask ports the protected getParsingTask(DSSDocument) override.
func (a *LOTLAnalysis) GetParsingTask(document model.DSSDocument) job.ParsingTask {
	return NewLOTLParsingTask(document, a.source)
}

// GetCurrentCertificateSource satisfies job.AbstractRunnableAnalysisOverrides with the base
// (non-overridden in Java, for LOTLAnalysis itself - LOTLWithPivotsAnalysis overrides it below)
// behaviour.
func (a *LOTLAnalysis) GetCurrentCertificateSource() spi.CertificateSource {
	return a.GetCurrentCertificateSourceDefault()
}
