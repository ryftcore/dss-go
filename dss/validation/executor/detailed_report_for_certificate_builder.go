// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/certificate/DetailedReportForCertificateBuilder.java
// (DSS 6.5.RC1).
//
// Deviations:
//
//   - Java's executeAllBasicBuildingBlocks() builds a plain HashMap and build()
//     appends bbbs.values() to the report. HashMap iteration order is
//     unspecified, so this port uses the insertion order recorded by
//     AbstractDetailedReportBuilder.BBBOrder - deterministic, and identical to
//     Java's for the single-entry map the base builder produces.
//
//   - DetailedReportForQWACBuilder overrides executeAllBasicBuildingBlocks(),
//     executeValidation(...) and getCertificateQualificationBlock(...), all of
//     which the base self-calls; the calls are routed through
//     DetailedReportForCertificateBuilderOverrides, registered by every
//     constructor via InitDetailedReportForCertificateBuilder.
//
//   - detailedReport.getLoTEAnalysis() is List<XmlLoTEAnalysis> in Java; the Go
//     model wraps each entry in XmlLoTEAnalysisEntry to record the xsi:type
//     choice, so the list is unwrapped before it is handed to the approval
//     status block.

package executor

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process/qualification"
)

// DetailedReportForCertificateBuilderOverrides captures the members Java's
// DetailedReportForQWACBuilder overrides and that the base self-calls.
type DetailedReportForCertificateBuilderOverrides interface {
	// ExecuteAllBasicBuildingBlocks executes all basic building blocks
	// required on validation. Port of the protected
	// executeAllBasicBuildingBlocks().
	ExecuteAllBasicBuildingBlocks() map[string]*jaxb.XmlBasicBuildingBlocks

	// ExecuteValidation performs validation for the given tokens. Port of the
	// protected executeValidation(XmlDetailedReport, Map).
	ExecuteValidation(detailedReport *jaxb.XmlDetailedReport, bbbs map[string]*jaxb.XmlBasicBuildingBlocks)

	// CertificateQualificationBlockFor gets the certificate qualification
	// block. Port of the protected
	// getCertificateQualificationBlock(XmlDetailedReport, XmlBasicBuildingBlocks).
	CertificateQualificationBlockFor(detailedReport *jaxb.XmlDetailedReport,
		basicBuildingBlocks *jaxb.XmlBasicBuildingBlocks) *qualification.CertificateQualificationBlock
}

// DetailedReportForCertificateBuilder builds a DetailedReport for a
// certificate validation. Port of DetailedReportForCertificateBuilder.
type DetailedReportForCertificateBuilder struct {
	AbstractDetailedReportBuilder

	// certificateId is the id of the certificate to validate. Port of the
	// private final certificateId field.
	certificateId string

	overrides DetailedReportForCertificateBuilderOverrides
}

// NewDetailedReportForCertificateBuilder is the default constructor. Port of
// DetailedReportForCertificateBuilder(Provider, Data,
// ValidationPolicy, Date, String).
func NewDetailedReportForCertificateBuilder(i18nProvider *i18n.Provider,
	diagnosticData *diagnostic.Data, validationPolicy policy.ValidationPolicy,
	currentTime time.Time, certificateId string) *DetailedReportForCertificateBuilder {
	b := &DetailedReportForCertificateBuilder{
		AbstractDetailedReportBuilder: NewAbstractDetailedReportBuilder(i18nProvider, currentTime, validationPolicy, diagnosticData),
		certificateId:                 certificateId,
	}
	b.InitDetailedReportForCertificateBuilder(b)
	return b
}

// InitDetailedReportForCertificateBuilder registers the concrete builder so
// the base can dispatch its overridden members. Every concrete constructor
// must call this before use.
func (b *DetailedReportForCertificateBuilder) InitDetailedReportForCertificateBuilder(
	overrides DetailedReportForCertificateBuilderOverrides) {
	b.overrides = overrides
}

// certificateBuilderOverrides returns the registered overrides, panicking when
// the concrete builder failed to call InitDetailedReportForCertificateBuilder.
func (b *DetailedReportForCertificateBuilder) certificateBuilderOverrides() DetailedReportForCertificateBuilderOverrides {
	if b.overrides == nil {
		panic("executor: DetailedReportForCertificateBuilder used without InitDetailedReportForCertificateBuilder")
	}
	return b.overrides
}

// CertificateId returns the id of the certificate to validate. It has no Java
// counterpart (the field is private and read by the subclass through the
// inherited members); the Go subclass lives in the same package but the
// accessor keeps the field encapsulated.
func (b *DetailedReportForCertificateBuilder) CertificateId() string {
	return b.certificateId
}

// Build builds the detailed report for the certificate validation. Port of
// build().
func (b *DetailedReportForCertificateBuilder) Build() *jaxb.XmlDetailedReport {
	detailedReport := b.Init()

	bbbs := b.certificateBuilderOverrides().ExecuteAllBasicBuildingBlocks()
	detailedReport.BasicBuildingBlocks = append(detailedReport.BasicBuildingBlocks, b.BasicBuildingBlocksInOrder(bbbs)...)

	b.certificateBuilderOverrides().ExecuteValidation(detailedReport, bbbs)
	return detailedReport
}

// Certificate gets the certificate to be validated. Port of the protected
// getCertificate(); Java's IllegalArgumentException becomes a panic, matching
// the convention used across the ported report builders (build() has no error
// return in Java either).
func (b *DetailedReportForCertificateBuilder) Certificate() *diagnostic.CertificateWrapper {
	certificate := b.DiagnosticData.UsedCertificateById(b.certificateId)
	if certificate == nil {
		panic(fmt.Sprintf("The certificate with the given Id '%s' has not been found in DiagnosticData", b.certificateId))
	}
	return certificate
}

// ExecuteAllBasicBuildingBlocks executes all basic building blocks required on
// validation. Port of the protected executeAllBasicBuildingBlocks().
func (b *DetailedReportForCertificateBuilder) ExecuteAllBasicBuildingBlocks() map[string]*jaxb.XmlBasicBuildingBlocks {
	bbbs := make(map[string]*jaxb.XmlBasicBuildingBlocks)
	b.process([]*diagnostic.CertificateWrapper{b.Certificate()},
		enumerations.ContextCertificate, bbbs)
	return bbbs
}

// ExecuteValidation performs validation for the given tokens. Port of the
// protected executeValidation(XmlDetailedReport, Map).
func (b *DetailedReportForCertificateBuilder) ExecuteValidation(detailedReport *jaxb.XmlDetailedReport,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) {
	b.BuildXmlCertificate(detailedReport, bbbs)
}

// BuildXmlCertificate executes the certificate validation and builds an
// XmlCertificate object. Port of the protected
// buildXmlCertificate(XmlDetailedReport, Map).
func (b *DetailedReportForCertificateBuilder) BuildXmlCertificate(detailedReport *jaxb.XmlDetailedReport,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) *jaxb.XmlCertificate {
	basicBuildingBlocks := bbbs[b.certificateId]

	xmlCertificate := &jaxb.XmlCertificate{}
	certificateId := b.certificateId
	xmlCertificate.Id = &certificateId

	cqb := b.certificateBuilderOverrides().CertificateQualificationBlockFor(detailedReport, basicBuildingBlocks)
	xmlCertificateQualificationProcess := cqb.Execute()
	xmlCertificate.CertificateQualificationProcess = xmlCertificateQualificationProcess

	if b.ValidateCertificateApprovalStatus() {
		cub := b.CertificateApprovalStatusBlockFor(detailedReport, basicBuildingBlocks)
		xmlCertificateApprovalStatusProcess := cub.Execute()
		xmlCertificate.CertificateApprovalStatusProcess = xmlCertificateApprovalStatusProcess
	}

	detailedReport.SignatureOrTimestampOrEvidenceRecord =
		append(detailedReport.SignatureOrTimestampOrEvidenceRecord, xmlCertificate)

	return xmlCertificate
}

// CertificateQualificationBlockFor gets the certificate qualification block.
// Port of the protected
// getCertificateQualificationBlock(XmlDetailedReport, XmlBasicBuildingBlocks).
func (b *DetailedReportForCertificateBuilder) CertificateQualificationBlockFor(detailedReport *jaxb.XmlDetailedReport,
	basicBuildingBlocks *jaxb.XmlBasicBuildingBlocks) *qualification.CertificateQualificationBlock {
	return qualification.NewCertificateQualificationBlock(b.I18nProvider, basicBuildingBlocks.Conclusion,
		b.CurrentTime, b.Certificate(), detailedReport.TLAnalysis)
}

// CertificateApprovalStatusBlockFor gets the certificate approval status
// block. Port of the protected
// getCertificateApprovalStatusBlock(XmlDetailedReport, XmlBasicBuildingBlocks).
func (b *DetailedReportForCertificateBuilder) CertificateApprovalStatusBlockFor(detailedReport *jaxb.XmlDetailedReport,
	basicBuildingBlocks *jaxb.XmlBasicBuildingBlocks) *qualification.CertificateApprovalStatusBlock {
	loteAnalysis := make([]*jaxb.XmlLoTEAnalysis, 0, len(detailedReport.LoTEAnalysis))
	for _, entry := range detailedReport.LoTEAnalysis {
		loteAnalysis = append(loteAnalysis, entry.XmlLoTEAnalysis)
	}
	return qualification.NewCertificateApprovalStatusBlock(b.I18nProvider, basicBuildingBlocks.Conclusion,
		b.CurrentTime, b.Certificate(), loteAnalysis)
}

// ValidateCertificateApprovalStatus checks if the certificate approval status
// is to be validated. Port of the protected
// validateCertificateApprovalStatus().
func (b *DetailedReportForCertificateBuilder) ValidateCertificateApprovalStatus() bool {
	return len(b.DiagnosticData.ListsOfTrustedEntities()) > 0
}
