// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/eaa/EAAPresentationProcessExecutor.java
// (DSS 6.5.RC1).

package executor

import (
	"github.com/utain/esig/dss/diagnostic"
)

// EAAPresentationProcessExecutor performs validation and reports building for
// an EAA Presentation. Port of EAAPresentationProcessExecutor.
type EAAPresentationProcessExecutor struct {
	DefaultSignatureProcessExecutor
}

// compile-time interface assertion.
var _ DocumentProcessExecutor = (*EAAPresentationProcessExecutor)(nil)

// NewEAAPresentationProcessExecutor is the default constructor. Port of
// EAAPresentationProcessExecutor().
func NewEAAPresentationProcessExecutor() *EAAPresentationProcessExecutor {
	e := &EAAPresentationProcessExecutor{
		DefaultSignatureProcessExecutor: *NewDefaultSignatureProcessExecutor(),
	}
	e.InitDefaultSignatureProcessExecutor(e)
	return e
}

// DetailedReportBuilderFor is the port of the overridden
// getDetailedReportBuilder(DiagnosticData); it hands back the embedded
// *DetailedReportBuilder of the EAA builder, which has registered itself so
// Build() still dispatches onto the EAA overrides.
func (e *EAAPresentationProcessExecutor) DetailedReportBuilderFor(
	diagnosticData *diagnostic.DiagnosticData) *DetailedReportBuilder {
	return &NewDetailedReportForEAAPresentationBuilder(e.I18nProvider(), e.CurrentTimeValue, e.Policy,
		diagnosticData, e.IncludeSemantics).DetailedReportBuilder
}
