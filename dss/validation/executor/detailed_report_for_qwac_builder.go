// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/certificate/qwac/DetailedReportForQWACBuilder.java
// (DSS 6.5.RC1).
//
// Deviations:
//
//   - Java's XmlConstraintsConclusionWithProofOfExistence supertype has no Go
//     counterpart; the port rebuilds that view from the two embedded structs
//     the concrete result carries, exactly as detailed_report_builder.go does.
//     The Conclusion pointer is shared, which matters: Java mutates the
//     conclusion returned by the basic validation IN PLACE and stores the same
//     object on the XmlSignature, so both <ValidationProcessBasicSignature>
//     and <Signature> end up carrying the TOTAL_* indication.
//
//   - getCertificateQualificationBlock() returns the embedded
//     *CertificateQualificationBlock of the QWAC block; the QWAC constructor
//     has registered itself with InitCertificateQualificationBlock, so
//     Execute() still dispatches onto the QWAC overrides.
//
//   - getSignatureFinalIndication() throws DSSReportException for an
//     unsupported indication; build() has no error return in Java either, so
//     the port panics with the same *reports.DSSReportException value.

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
	"github.com/ryftcore/dss-go/dss/validation/process/vpfbs"
	"github.com/ryftcore/dss-go/dss/validation/reports"
)

// DetailedReportForQWACBuilder builds a Detailed Report for a QWAC
// certificate validation. Port of DetailedReportForQWACBuilder.
type DetailedReportForQWACBuilder struct {
	DetailedReportForCertificateBuilder
}

// NewDetailedReportForQWACBuilder is the default constructor. Port of
// DetailedReportForQWACBuilder(I18nProvider, DiagnosticData, ValidationPolicy,
// Date, String).
func NewDetailedReportForQWACBuilder(i18nProvider *i18n.I18nProvider,
	diagnosticData *diagnostic.DiagnosticData, validationPolicy policy.ValidationPolicy,
	currentTime time.Time, certificateId string) *DetailedReportForQWACBuilder {
	b := &DetailedReportForQWACBuilder{
		DetailedReportForCertificateBuilder: *NewDetailedReportForCertificateBuilder(
			i18nProvider, diagnosticData, validationPolicy, currentTime, certificateId),
	}
	b.InitDetailedReportForCertificateBuilder(b)
	return b
}

// ExecuteAllBasicBuildingBlocks is the port of the overridden
// executeAllBasicBuildingBlocks().
func (b *DetailedReportForQWACBuilder) ExecuteAllBasicBuildingBlocks() map[string]*jaxb.XmlBasicBuildingBlocks {
	bbbs := b.DetailedReportForCertificateBuilder.ExecuteAllBasicBuildingBlocks()
	b.process(javaHashSetOrder(b.DiagnosticData.AllSignatures()), enumerations.Context_SIGNATURE, bbbs)
	return bbbs
}

// ExecuteValidation is the port of the overridden
// executeValidation(XmlDetailedReport, Map).
func (b *DetailedReportForQWACBuilder) ExecuteValidation(detailedReport *jaxb.XmlDetailedReport,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) {
	xmlSignature := b.BuildXmlTLSCertificateBindingSignature(detailedReport, bbbs)
	b.BuildXmlCertificateWithBindingSignature(detailedReport, bbbs, xmlSignature)
}

// BuildXmlTLSCertificateBindingSignature validates and builds the report
// information for the validation of a TLS Certificate Binding signature, when
// present. Port of the protected
// buildXmlTLSCertificateBindingSignature(XmlDetailedReport, Map).
func (b *DetailedReportForQWACBuilder) BuildXmlTLSCertificateBindingSignature(detailedReport *jaxb.XmlDetailedReport,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) *jaxb.XmlSignature {
	bindingSignature := b.DiagnosticData.TLSCertificateBindingSignature()
	if bindingSignature == nil {
		return nil
	}

	xmlSignature := &jaxb.XmlSignature{}
	id := bindingSignature.Id()
	xmlSignature.Id = &id

	basic := b.executeBasicValidation(xmlSignature, bindingSignature, bbbs)
	validation := &jaxb.XmlConstraintsConclusionWithProofOfExistence{
		XmlConstraintsConclusionWithProofOfExistenceContent: basic.XmlConstraintsConclusionWithProofOfExistenceContent,
		XmlConstraintsConclusionAttrs:                       basic.XmlConstraintsConclusionAttrs,
	}

	if b.Policy.EIDASConstraintPresent() {

		// Signature qualification
		qualificationBlock := qualification.NewTLSBindingSignatureQualificationBlock(
			b.I18nProvider, bbbs, validation, bindingSignature, detailedReport.TLAnalysis, b.DiagnosticData.WebsiteUrl())
		xmlSignature.ValidationSignatureQualification = qualificationBlock.Execute()

	}

	conclusion := validation.Conclusion
	conclusion.Indication = jaxb.IndicationValue(b.signatureFinalIndication(conclusion.Indication.Indication()))
	xmlSignature.Conclusion = conclusion

	detailedReport.SignatureOrTimestampOrEvidenceRecord =
		append(detailedReport.SignatureOrTimestampOrEvidenceRecord, xmlSignature)

	return xmlSignature
}

// executeBasicValidation is the port of the private
// executeBasicValidation(XmlSignature, SignatureWrapper, Map).
func (b *DetailedReportForQWACBuilder) executeBasicValidation(signatureAnalysis *jaxb.XmlSignature,
	signature *diagnostic.SignatureWrapper,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) *jaxb.XmlValidationProcessBasicSignature {
	// Java passes Collections.emptyList() for the timestamps.
	vpfbsProcess := vpfbs.NewBasicSignatureValidationProcess(
		b.I18nProvider, b.DiagnosticData, signature, nil, bbbs)
	bs := vpfbsProcess.Execute()
	signatureAnalysis.ValidationProcessBasicSignature = bs
	return bs
}

// BuildXmlCertificateWithBindingSignature builds the XmlCertificate. Port of
// the protected buildXmlCertificate(XmlDetailedReport, Map, XmlSignature) -
// an OVERLOAD of the two-argument inherited method, not an override, so it
// keeps a distinct Go name.
func (b *DetailedReportForQWACBuilder) BuildXmlCertificateWithBindingSignature(detailedReport *jaxb.XmlDetailedReport,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, bindingSignature *jaxb.XmlSignature) *jaxb.XmlCertificate {
	xmlCertificate := b.DetailedReportForCertificateBuilder.BuildXmlCertificate(detailedReport, bbbs)

	qwacForTLSCertificateValidationBlock := qualification.NewQWACForTLSCertificateValidationBlock(
		b.I18nProvider, b.DiagnosticData, b.Certificate(), bbbs, xmlCertificate.CertificateQualificationProcess,
		b.qwacProfile(bindingSignature), b.DiagnosticData.WebsiteUrl())
	xmlQWACProcess := qwacForTLSCertificateValidationBlock.Execute()
	xmlCertificate.QWACProcess = xmlQWACProcess

	return xmlCertificate
}

// CertificateQualificationBlockFor is the port of the overridden
// getCertificateQualificationBlock(XmlDetailedReport, XmlBasicBuildingBlocks).
func (b *DetailedReportForQWACBuilder) CertificateQualificationBlockFor(detailedReport *jaxb.XmlDetailedReport,
	basicBuildingBlocks *jaxb.XmlBasicBuildingBlocks) *qualification.CertificateQualificationBlock {
	return qualification.NewCertificateQualificationForQWACBlock(b.I18nProvider, basicBuildingBlocks.Conclusion,
		b.CurrentTime, b.Certificate(), detailedReport.TLAnalysis).CertificateQualificationBlock
}

// qwacProfile is the port of the private getQWACProfile(XmlSignature).
func (b *DetailedReportForQWACBuilder) qwacProfile(bindingSignature *jaxb.XmlSignature) enumerations.QWACProfile {
	if bindingSignature != nil && bindingSignature.ValidationSignatureQualification != nil {
		qwacProcess := bindingSignature.ValidationSignatureQualification.QWACProcess
		if qwacProcess != nil {
			return qwacProcess.QWACType.QWACProfile()
		}
	}
	return enumerations.QWACProfile_NOT_QWAC
}

// signatureFinalIndication is the port of the private
// getSignatureFinalIndication(Indication).
func (b *DetailedReportForQWACBuilder) signatureFinalIndication(highestIndication enumerations.Indication) enumerations.Indication {
	switch highestIndication {
	case enumerations.Indication_PASSED:
		return enumerations.Indication_TOTAL_PASSED
	case enumerations.Indication_INDETERMINATE:
		return enumerations.Indication_INDETERMINATE
	case enumerations.Indication_FAILED:
		return enumerations.Indication_TOTAL_FAILED
	default:
		panic(reports.NewDSSReportExceptionMessage(
			fmt.Sprintf("The Indication '%s' is not supported!", highestIndication)))
	}
}
