package dss

import (
	"time"

	"github.com/ryftcore/dss-go/dss/simplereport"
	"github.com/ryftcore/dss-go/dss/validation/reports"
)

// Reports is what [Validate] returns: the four DSS validation reports, plus a
// few accessors for the questions most callers actually have.
//
// It embeds the upstream [reports.Reports], so the whole DSS report API stays
// reachable - GetSimpleReport, GetDetailedReport, GetDiagnosticData,
// GetEtsiValidationReportJaxb and the JAXB models behind them - without going
// through this type. The methods declared here add nothing the embedded API
// cannot express; they only spare the caller the walk.
type Reports struct {
	*reports.Reports
}

// Verdict is the outcome of validating one signature, gathered from the
// SimpleReport.
type Verdict struct {
	// ID is the signature identifier the reports use throughout. It is stable
	// for a given signature and is the key into every other report.
	ID string

	// Indication is the EN 319 102-1 verdict: TOTAL_PASSED, INDETERMINATE or
	// TOTAL_FAILED.
	Indication Indication

	// SubIndication says why, and is empty for a TOTAL_PASSED signature.
	SubIndication SubIndication

	// SignatureLevel is the level the signature was recognised as, for
	// example "XAdES-BASELINE-LTA" when printed. It reports what the
	// signature IS, not what it should have been.
	SignatureLevel SignatureLevel

	// Qualification is the eIDAS qualification determined from the trusted
	// lists, for example "QESig". It is "NA" when no trusted-list information
	// was supplied - see ValidateOptions.TrustedCertificateSources - which is
	// a statement that the question could not be answered, not that the
	// signature is unqualified.
	Qualification SignatureQualification

	// SignedBy is the signing certificate's subject as the report renders it.
	SignedBy string

	// SigningTime is the signing time claimed in the signed attributes, which
	// nothing but the signer vouches for. Nil when the signature carries none.
	SigningTime *time.Time

	// BestSignatureTime is the earliest time the signature is PROVEN to have
	// existed at, from the time-stamps covering it. It falls back to the
	// validation time when no time-stamp proves anything.
	BestSignatureTime *time.Time

	// Errors, Warnings and Infos are the AdES validation messages behind the
	// indication, already localised.
	Errors   []simplereport.Message
	Warnings []simplereport.Message
	Infos    []simplereport.Message
}

// Valid reports whether the signature reached TOTAL_PASSED.
func (v Verdict) Valid() bool { return v.Indication == IndicationTotalPassed }

// TimestampVerdict is the outcome of validating one time-stamp token.
type TimestampVerdict struct {
	// ID is the time-stamp identifier the reports use.
	ID string

	// Indication is the verdict for the time-stamp token itself.
	Indication Indication

	// SubIndication says why, and is empty for a passed time-stamp.
	SubIndication SubIndication

	// Qualification is the eIDAS qualification of the time-stamp, "NA" unless
	// trusted-list information was supplied.
	Qualification TimestampQualification

	// ProductionTime is the time the TSA asserts, and ProducedBy names the
	// TSA that asserted it.
	ProductionTime *time.Time
	ProducedBy     string
}

// Verdicts returns one [Verdict] per signature found in the document, in the
// order the SimpleReport lists them. A document with no signature yields an
// empty slice, not an error.
func (r *Reports) Verdicts() []Verdict {
	simple := r.GetSimpleReport()
	ids := simple.GetSignatureIdList()
	verdicts := make([]Verdict, 0, len(ids))
	for _, id := range ids {
		verdicts = append(verdicts, Verdict{
			ID:                id,
			Indication:        simple.GetIndication(id),
			SubIndication:     simple.GetSubIndication(id),
			SignatureLevel:    simple.GetSignatureFormat(id),
			Qualification:     simple.GetSignatureQualification(id),
			SignedBy:          simple.GetSignedBy(id),
			SigningTime:       simple.GetSigningTime(id),
			BestSignatureTime: simple.GetBestSignatureTime(id),
			Errors:            simple.GetAdESValidationErrors(id),
			Warnings:          simple.GetAdESValidationWarnings(id),
			Infos:             simple.GetAdESValidationInfo(id),
		})
	}
	return verdicts
}

// TimestampVerdicts returns one [TimestampVerdict] per detached time-stamp
// token the document carries. Time-stamps embedded in a signature are reported
// under that signature in the detailed report, not here.
func (r *Reports) TimestampVerdicts() []TimestampVerdict {
	simple := r.GetSimpleReport()
	ids := simple.GetTimestampIdList()
	verdicts := make([]TimestampVerdict, 0, len(ids))
	for _, id := range ids {
		verdicts = append(verdicts, TimestampVerdict{
			ID:             id,
			Indication:     simple.GetIndication(id),
			SubIndication:  simple.GetSubIndication(id),
			Qualification:  simple.GetTimestampQualification(id),
			ProductionTime: simple.GetProductionTime(id),
			ProducedBy:     simple.GetProducedBy(id),
		})
	}
	return verdicts
}

// Valid reports whether the document carries at least one signature and every
// signature reached TOTAL_PASSED. It is the single-boolean answer; anything
// more nuanced needs [Reports.Verdicts].
func (r *Reports) Valid() bool {
	simple := r.GetSimpleReport()
	count := simple.GetSignaturesCount()
	return count > 0 && simple.GetValidSignaturesCount() == count
}

// SignatureCount returns the number of signatures found in the document.
func (r *Reports) SignatureCount() int {
	return r.GetSimpleReport().GetSignaturesCount()
}

// ValidSignatureCount returns the number of signatures that reached
// TOTAL_PASSED.
func (r *Reports) ValidSignatureCount() int {
	return r.GetSimpleReport().GetValidSignaturesCount()
}

// SimpleReportXML marshals the SimpleReport - the short, human-oriented
// verdict document. Passthrough of GetXmlSimpleReport.
func (r *Reports) SimpleReportXML() (string, error) { return r.GetXmlSimpleReport() }

// DetailedReportXML marshals the DetailedReport - every EN 319 102-1 building
// block and check, with its own conclusion. Passthrough of
// GetXmlDetailedReport.
func (r *Reports) DetailedReportXML() (string, error) { return r.GetXmlDetailedReport() }

// DiagnosticDataXML marshals the diagnostic data - the raw facts the process
// reasoned over: certificates, revocation data, time-stamps, signature
// properties. Passthrough of GetXmlDiagnosticData.
func (r *Reports) DiagnosticDataXML() (string, error) { return r.GetXmlDiagnosticData() }

// ETSIValidationReportXML marshals the ETSI TS 119 102-2 validation report,
// the standardised, machine-readable report format. Passthrough of
// GetXmlValidationReport.
func (r *Reports) ETSIValidationReportXML() (string, error) { return r.GetXmlValidationReport() }
