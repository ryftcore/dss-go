// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/signature/DefaultSignatureProcessExecutor.java
// (DSS 6.5.RC1).
//
// EAAPresentationProcessExecutor overrides the protected
// getDetailedReportBuilder(DiagnosticData), which buildReports() self-calls,
// so the call is routed through DefaultSignatureProcessExecutorOverrides,
// registered by every constructor via InitDefaultSignatureProcessExecutor.

package executor

import (
	"github.com/ryftcore/dss-go/dss/detailedreport"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/validation/reports"
	validationreportjaxb "github.com/ryftcore/dss-go/dss/validationreport/jaxb"
)

// DefaultSignatureProcessExecutorOverrides captures the member Java's
// EAAPresentationProcessExecutor overrides and that buildReports()
// self-calls.
type DefaultSignatureProcessExecutorOverrides interface {
	// DetailedReportBuilderFor instantiates a builder for the Detailed
	// Report. Port of the protected getDetailedReportBuilder(DiagnosticData).
	DetailedReportBuilderFor(diagnosticData *diagnostic.DiagnosticData) *DetailedReportBuilder
}

// DefaultSignatureProcessExecutor executes a signature validation process and
// produces the SimpleReport, DetailedReport and ETSI Validation report. Port
// of DefaultSignatureProcessExecutor.
type DefaultSignatureProcessExecutor struct {
	AbstractProcessExecutor

	// ValidationLevel is the target highest validation level. Port of the
	// protected validationLevel field (default: ARCHIVAL_DATA).
	ValidationLevel enumerations.ValidationLevel

	// EnableEtsiValidationReport defines if the ETSI Validation Report shall
	// be generated. Port of the protected enableEtsiValidationReport field
	// (default: true).
	EnableEtsiValidationReport bool

	// IncludeSemantics defines if the semantics information shall be
	// included. Port of the protected includeSemantics field (default:
	// false).
	IncludeSemantics bool

	overrides DefaultSignatureProcessExecutorOverrides
}

// compile-time interface assertion.
var _ DocumentProcessExecutor = (*DefaultSignatureProcessExecutor)(nil)

// NewDefaultSignatureProcessExecutor is the default constructor instantiating
// the object with the default configuration. Port of
// DefaultSignatureProcessExecutor() plus the field initializers Java runs
// with it.
func NewDefaultSignatureProcessExecutor() *DefaultSignatureProcessExecutor {
	e := &DefaultSignatureProcessExecutor{
		AbstractProcessExecutor:    NewAbstractProcessExecutor(),
		ValidationLevel:            enumerations.ValidationLevel_ARCHIVAL_DATA,
		EnableEtsiValidationReport: true,
		IncludeSemantics:           false,
	}
	e.InitDefaultSignatureProcessExecutor(e)
	return e
}

// InitDefaultSignatureProcessExecutor registers the concrete executor so that
// BuildReports() dispatches getDetailedReportBuilder() onto it. Every concrete
// constructor must call this before use.
func (e *DefaultSignatureProcessExecutor) InitDefaultSignatureProcessExecutor(
	overrides DefaultSignatureProcessExecutorOverrides) {
	e.overrides = overrides
}

// signatureExecutorOverrides returns the registered overrides, panicking when
// the concrete executor failed to call InitDefaultSignatureProcessExecutor.
func (e *DefaultSignatureProcessExecutor) signatureExecutorOverrides() DefaultSignatureProcessExecutorOverrides {
	if e.overrides == nil {
		panic("executor: DefaultSignatureProcessExecutor used without InitDefaultSignatureProcessExecutor")
	}
	return e.overrides
}

// SetValidationLevel is the port of setValidationLevel(ValidationLevel).
func (e *DefaultSignatureProcessExecutor) SetValidationLevel(validationLevel enumerations.ValidationLevel) {
	e.ValidationLevel = validationLevel
}

// SetEnableEtsiValidationReport is the port of
// setEnableEtsiValidationReport(boolean).
func (e *DefaultSignatureProcessExecutor) SetEnableEtsiValidationReport(enableEtsiValidationReport bool) {
	e.EnableEtsiValidationReport = enableEtsiValidationReport
}

// SetIncludeSemantics is the port of setIncludeSemantics(boolean).
func (e *DefaultSignatureProcessExecutor) SetIncludeSemantics(includeSemantics bool) {
	e.IncludeSemantics = includeSemantics
}

// Execute runs the validation process. Port of execute().
//
// Java's Objects.requireNonNull(validationLevel, "The validation level is
// missing") signals a programming error, so it is ported as a panic - matching
// AbstractProcessExecutor.AssertConfigurationValid.
func (e *DefaultSignatureProcessExecutor) Execute() *reports.Reports {
	e.AssertConfigurationValid()
	if e.ValidationLevel == "" {
		panic("The validation level is missing")
	}
	diagnosticData := e.DiagnosticData()
	return e.BuildReports(diagnosticData)
}

// DiagnosticData gets the DiagnosticData. Port of the protected
// getDiagnosticData().
func (e *DefaultSignatureProcessExecutor) DiagnosticData() *diagnostic.DiagnosticData {
	return diagnostic.NewDiagnosticData(e.JaxbDiagnosticData)
}

// BuildReports builds the reports. Port of the protected
// buildReports(DiagnosticData).
func (e *DefaultSignatureProcessExecutor) BuildReports(diagnosticData *diagnostic.DiagnosticData) *reports.Reports {

	detailedReportBuilder := e.signatureExecutorOverrides().DetailedReportBuilderFor(diagnosticData)
	jaxbDetailedReport := detailedReportBuilder.Build()

	detailedReportWrapper := detailedreport.NewDetailedReport(jaxbDetailedReport)

	simpleReportBuilder := e.SimpleReportBuilderFor(diagnosticData, detailedReportWrapper)
	simpleReport := simpleReportBuilder.Build()

	var validationReport *validationreportjaxb.ValidationReportType
	if e.EnableEtsiValidationReport {
		etsiValidationReportBuilder := e.ETSIValidationReportBuilderFor(diagnosticData, detailedReportWrapper)
		validationReport = etsiValidationReportBuilder.Build()
	}

	return reports.NewReports(e.JaxbDiagnosticData, jaxbDetailedReport, simpleReport, validationReport)
}

// DetailedReportBuilderFor instantiates a builder for the Detailed Report.
// Port of the protected getDetailedReportBuilder(DiagnosticData).
func (e *DefaultSignatureProcessExecutor) DetailedReportBuilderFor(
	diagnosticData *diagnostic.DiagnosticData) *DetailedReportBuilder {
	return NewDetailedReportBuilder(e.I18nProvider(), e.CurrentTimeValue, e.Policy,
		e.ValidationLevel, diagnosticData, e.IncludeSemantics)
}

// SimpleReportBuilderFor instantiates a builder for the Simple Report. Port of
// the protected getSimpleReportBuilder(DiagnosticData, DetailedReport). No
// upstream subclass overrides it, so it is not routed through the overrides
// interface.
func (e *DefaultSignatureProcessExecutor) SimpleReportBuilderFor(diagnosticData *diagnostic.DiagnosticData,
	detailedReport *detailedreport.DetailedReport) *SimpleReportBuilder {
	return NewSimpleReportBuilder(e.I18nProvider(), e.CurrentTimeValue, e.Policy,
		diagnosticData, detailedReport, e.IncludeSemantics)
}

// ETSIValidationReportBuilderFor instantiates a builder for the ETSI
// Validation Report 102-2. Port of the protected
// getETSIValidationReportBuilder(DiagnosticData, DetailedReport). No upstream
// subclass overrides it, so it is not routed through the overrides interface.
func (e *DefaultSignatureProcessExecutor) ETSIValidationReportBuilderFor(diagnosticData *diagnostic.DiagnosticData,
	detailedReport *detailedreport.DetailedReport) *ETSIValidationReportBuilder {
	return NewETSIValidationReportBuilder(e.CurrentTimeValue, diagnosticData, detailedReport)
}
