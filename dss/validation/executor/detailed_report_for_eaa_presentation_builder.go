// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/eaa/DetailedReportForEAAPresentationBuilder.java
// (DSS 6.5.RC1).
//
// detailedReport.getLoTEAnalysis() is List<XmlLoTEAnalysis> in Java; the Go
// model wraps each entry in XmlLoTEAnalysisEntry to record the xsi:type
// choice, so the list is unwrapped before it is handed to the EAA block (see
// detailed_report_for_certificate_builder.go for the same unwrapping).

package executor

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process/eaa"
	"github.com/utain/esig/dss/validation/process/vpfswatsp"
)

// DetailedReportForEAAPresentationBuilder performs validation of an EAA
// presentation validation. Port of DetailedReportForEAAPresentationBuilder.
type DetailedReportForEAAPresentationBuilder struct {
	DetailedReportBuilder
}

// NewDetailedReportForEAAPresentationBuilder is the default constructor. Port
// of DetailedReportForEAAPresentationBuilder(I18nProvider, Date,
// ValidationPolicy, DiagnosticData, boolean), which hardwires
// ValidationLevel.BASIC_SIGNATURES.
func NewDetailedReportForEAAPresentationBuilder(i18nProvider *i18n.I18nProvider, currentTime time.Time,
	validationPolicy policy.ValidationPolicy, diagnosticData *diagnostic.DiagnosticData,
	includeSemantics bool) *DetailedReportForEAAPresentationBuilder {
	b := &DetailedReportForEAAPresentationBuilder{
		DetailedReportBuilder: *NewDetailedReportBuilder(i18nProvider, currentTime, validationPolicy,
			enumerations.ValidationLevel_BASIC_SIGNATURES, diagnosticData, includeSemantics),
	}
	b.InitDetailedReportBuilder(b)
	return b
}

// ExecuteValidation is the port of the overridden
// executeValidation(XmlDetailedReport, Map, POEExtraction).
func (b *DetailedReportForEAAPresentationBuilder) ExecuteValidation(detailedReport *jaxb.XmlDetailedReport,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, poe *vpfswatsp.POEExtraction) {
	eaas := b.executeEAAValidations(bbbs, detailedReport.TLAnalysis, detailedReport.LoTEAnalysis)
	for _, xmlEAA := range eaas {
		detailedReport.SignatureOrTimestampOrEvidenceRecord =
			append(detailedReport.SignatureOrTimestampOrEvidenceRecord, xmlEAA)
	}
}

// executeEAAValidations is the port of the private executeEAAValidations(Map,
// List, List).
func (b *DetailedReportForEAAPresentationBuilder) executeEAAValidations(
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tlAnalysis []*jaxb.XmlTLAnalysis,
	loteAnalysis []*jaxb.XmlLoTEAnalysisEntry) []*jaxb.XmlEAA {
	unwrappedLoteAnalysis := make([]*jaxb.XmlLoTEAnalysis, 0, len(loteAnalysis))
	for _, entry := range loteAnalysis {
		unwrappedLoteAnalysis = append(unwrappedLoteAnalysis, entry.XmlLoTEAnalysis)
	}
	eaaValidationBlock := eaa.NewEAAValidationBlock(
		b.I18nProvider, b.DiagnosticData, b.Policy, b.CurrentTime, bbbs, tlAnalysis, unwrappedLoteAnalysis)
	return eaaValidationBlock.Execute()
}
