// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/certificate/qwac/QWACCertificateProcessExecutor.java
// (DSS 6.5.RC1).

package executor

import (
	"github.com/utain/esig/dss/diagnostic"
)

// QWACCertificateProcessExecutor executes a QWAC certificate validation. Port
// of QWACCertificateProcessExecutor.
type QWACCertificateProcessExecutor struct {
	DefaultCertificateProcessExecutor
}

// compile-time interface assertion.
var _ CertificateProcessExecutor = (*QWACCertificateProcessExecutor)(nil)

// NewQWACCertificateProcessExecutor is the default constructor. Port of
// QWACCertificateProcessExecutor().
func NewQWACCertificateProcessExecutor() *QWACCertificateProcessExecutor {
	e := &QWACCertificateProcessExecutor{
		DefaultCertificateProcessExecutor: *NewDefaultCertificateProcessExecutor(),
	}
	e.InitDefaultCertificateProcessExecutor(e)
	return e
}

// DetailedReportBuilderFor is the port of the overridden
// getDetailedReportBuilder(DiagnosticData); it hands back the embedded
// *DetailedReportForCertificateBuilder of the QWAC builder, which has
// registered itself so Build() still dispatches onto the QWAC overrides.
func (e *QWACCertificateProcessExecutor) DetailedReportBuilderFor(
	diagnosticData *diagnostic.DiagnosticData) *DetailedReportForCertificateBuilder {
	return &NewDetailedReportForQWACBuilder(e.I18nProvider(), diagnosticData, e.Policy,
		e.CurrentTimeValue, e.CertificateId).DetailedReportForCertificateBuilder
}
