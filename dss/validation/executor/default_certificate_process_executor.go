// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/certificate/DefaultCertificateProcessExecutor.java
// (DSS 6.5.RC1).
//
// QWACCertificateProcessExecutor overrides the protected
// getDetailedReportBuilder(DiagnosticData), which execute() self-calls, so the
// call is routed through DefaultCertificateProcessExecutorOverrides,
// registered by every constructor via InitDefaultCertificateProcessExecutor -
// the same arrangement as process/chain_item.go's InitChainItem. The override
// returns the EMBEDDED *DetailedReportForCertificateBuilder of the QWAC
// builder, whose own InitDetailedReportForCertificateBuilder call keeps
// Build() dispatching onto the QWAC overrides.

package executor

import (
	"github.com/utain/esig/dss/detailedreport"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/validation/reports"
)

// DefaultCertificateProcessExecutorOverrides captures the members Java's
// QWACCertificateProcessExecutor overrides and that execute() self-calls.
type DefaultCertificateProcessExecutorOverrides interface {
	// DetailedReportBuilderFor gets the Detailed report builder. Port of the
	// protected getDetailedReportBuilder(DiagnosticData).
	DetailedReportBuilderFor(diagnosticData *diagnostic.DiagnosticData) *DetailedReportForCertificateBuilder
}

// DefaultCertificateProcessExecutor executes a certificate validation. Port of
// DefaultCertificateProcessExecutor.
type DefaultCertificateProcessExecutor struct {
	AbstractProcessExecutor

	// CertificateId is the id of a certificate to validate. Port of the
	// protected certificateId field.
	CertificateId string

	overrides DefaultCertificateProcessExecutorOverrides
}

// compile-time interface assertion.
var _ CertificateProcessExecutor = (*DefaultCertificateProcessExecutor)(nil)

// NewDefaultCertificateProcessExecutor is the default constructor
// instantiating the object with a null certificate id. Port of
// DefaultCertificateProcessExecutor().
func NewDefaultCertificateProcessExecutor() *DefaultCertificateProcessExecutor {
	e := &DefaultCertificateProcessExecutor{
		AbstractProcessExecutor: NewAbstractProcessExecutor(),
	}
	e.InitDefaultCertificateProcessExecutor(e)
	return e
}

// InitDefaultCertificateProcessExecutor registers the concrete executor so
// that Execute() dispatches getDetailedReportBuilder() onto it. Every concrete
// constructor must call this before use.
func (e *DefaultCertificateProcessExecutor) InitDefaultCertificateProcessExecutor(
	overrides DefaultCertificateProcessExecutorOverrides) {
	e.overrides = overrides
}

// certificateExecutorOverrides returns the registered overrides, panicking
// when the concrete executor failed to call
// InitDefaultCertificateProcessExecutor.
func (e *DefaultCertificateProcessExecutor) certificateExecutorOverrides() DefaultCertificateProcessExecutorOverrides {
	if e.overrides == nil {
		panic("executor: DefaultCertificateProcessExecutor used without InitDefaultCertificateProcessExecutor")
	}
	return e.overrides
}

// SetCertificateId is the port of setCertificateId(String).
func (e *DefaultCertificateProcessExecutor) SetCertificateId(certificateId string) {
	e.CertificateId = certificateId
}

// Execute executes the certificate validation. Port of execute().
//
// Java's Objects.requireNonNull(certificateId, "The certificate id is
// missing") signals a programming error (an NPE escapes execute()), so it is
// ported as a panic - matching AbstractProcessExecutor.AssertConfigurationValid.
func (e *DefaultCertificateProcessExecutor) Execute() *reports.CertificateReports {
	e.AssertConfigurationValid()
	if e.CertificateId == "" {
		panic("The certificate id is missing")
	}

	diagnosticData := e.DiagnosticData()

	detailedReportBuilder := e.certificateExecutorOverrides().DetailedReportBuilderFor(diagnosticData)
	xmlDetailedReport := detailedReportBuilder.Build()
	detailedReport := detailedreport.NewDetailedReport(xmlDetailedReport)

	simpleReportBuilder := e.SimpleReportBuilderFor(diagnosticData, detailedReport)
	simpleReport := simpleReportBuilder.Build()

	return reports.NewCertificateReports(e.JaxbDiagnosticData, xmlDetailedReport, simpleReport)
}

// DiagnosticData gets the Diagnostic Data. Port of the protected
// getDiagnosticData().
func (e *DefaultCertificateProcessExecutor) DiagnosticData() *diagnostic.DiagnosticData {
	return diagnostic.NewDiagnosticData(e.JaxbDiagnosticData)
}

// DetailedReportBuilderFor gets the Detailed report builder. Port of the
// protected getDetailedReportBuilder(DiagnosticData).
func (e *DefaultCertificateProcessExecutor) DetailedReportBuilderFor(
	diagnosticData *diagnostic.DiagnosticData) *DetailedReportForCertificateBuilder {
	return NewDetailedReportForCertificateBuilder(e.I18nProvider(), diagnosticData, e.Policy,
		e.CurrentTimeValue, e.CertificateId)
}

// SimpleReportBuilderFor gets the Simple report builder. Port of the protected
// getSimpleReportBuilder(DiagnosticData, DetailedReport). No upstream subclass
// overrides it, so it is not routed through the overrides interface.
func (e *DefaultCertificateProcessExecutor) SimpleReportBuilderFor(diagnosticData *diagnostic.DiagnosticData,
	detailedReport *detailedreport.DetailedReport) *SimpleReportForCertificateBuilder {
	return NewSimpleReportForCertificateBuilder(diagnosticData, detailedReport, e.Policy,
		e.CurrentTimeValue, e.CertificateId)
}
