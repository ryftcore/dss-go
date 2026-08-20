// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/signature/DetailedReportBuilder.java
// (DSS 6.5.RC1).
//
// Deviations, all forced and byte-visible only where noted:
//
//   - Java collects the used indications in HashSet<Indication>/
//     HashSet<SubIndication>. Enum hashCode() is the identity hash, so the
//     iteration order - and therefore the order of the <Semantic> elements
//     addSemantics() appends - is not even stable across JVM runs upstream.
//     This port emits them in enumeration declaration order, which is
//     deterministic. Semantics are off by default (includeSemantics=false).
//
//   - The basic building blocks are appended in the insertion order recorded
//     by AbstractDetailedReportBuilder.BBBOrder, standing in for Java's
//     LinkedHashMap.values(); see that file's header.
//
//   - getFinalIndication() throws DSSReportException for an unsupported
//     indication. build() has no error return in Java either, so the port
//     panics with the same *reports.DSSReportException value, matching the
//     convention used across the ported process tree.
//
//   - DetailedReportForEAAPresentationBuilder overrides the protected
//     executeValidation(...), which build() self-calls. Go has no virtual
//     dispatch, so the call is routed through DetailedReportBuilderOverrides,
//     registered by every constructor via InitDetailedReportBuilder - the same
//     arrangement as process/chain_item.go's InitChainItem.

package executor

import (
	"fmt"
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process/qualification"
	"github.com/utain/esig/dss/validation/process/vpfbs"
	"github.com/utain/esig/dss/validation/process/vpfltvdsig"
	"github.com/utain/esig/dss/validation/process/vpfswatsp"
	"github.com/utain/esig/dss/validation/process/vpfswatsp/evidencerecord"
	"github.com/utain/esig/dss/validation/process/vpftsp"
	"github.com/utain/esig/dss/validation/reports"
)

// DetailedReportBuilderOverrides captures the member Java's
// DetailedReportForEAAPresentationBuilder overrides and that build()
// self-calls.
type DetailedReportBuilderOverrides interface {
	// ExecuteValidation performs validation for the given tokens. Port of the
	// protected executeValidation(XmlDetailedReport, Map, POEExtraction).
	ExecuteValidation(detailedReport *jaxb.XmlDetailedReport,
		bbbs map[string]*jaxb.XmlBasicBuildingBlocks, poe *vpfswatsp.POEExtraction)
}

// DetailedReportBuilder builds a DetailedReport for a signature validation.
// Port of DetailedReportBuilder.
type DetailedReportBuilder struct {
	AbstractDetailedReportBuilder

	// overrides points back at the concrete builder; see
	// InitDetailedReportBuilder.
	overrides DetailedReportBuilderOverrides

	// validationLevel is the target highest validation level.
	validationLevel enumerations.ValidationLevel

	// includeSemantics defines if the semantics information shall be
	// included.
	includeSemantics bool

	// allIndications is the set of all used Indications (used for
	// semantics). See the file header for the ordering deviation.
	allIndications map[enumerations.Indication]struct{}

	// allSubIndications is the set of all used SubIndications (used for
	// semantics).
	allSubIndications map[enumerations.SubIndication]struct{}
}

// NewDetailedReportBuilder is the default constructor. Port of
// DetailedReportBuilder(I18nProvider, Date, ValidationPolicy,
// ValidationLevel, DiagnosticData, boolean).
func NewDetailedReportBuilder(i18nProvider *i18n.I18nProvider, currentTime time.Time,
	validationPolicy policy.ValidationPolicy, validationLevel enumerations.ValidationLevel,
	diagnosticData *diagnostic.DiagnosticData, includeSemantics bool) *DetailedReportBuilder {
	b := &DetailedReportBuilder{
		AbstractDetailedReportBuilder: NewAbstractDetailedReportBuilder(i18nProvider, currentTime, validationPolicy, diagnosticData),
		validationLevel:               validationLevel,
		includeSemantics:              includeSemantics,
		allIndications:                make(map[enumerations.Indication]struct{}),
		allSubIndications:             make(map[enumerations.SubIndication]struct{}),
	}
	b.InitDetailedReportBuilder(b)
	return b
}

// InitDetailedReportBuilder registers the concrete builder so that Build()
// dispatches executeValidation() onto it. Every concrete constructor must call
// this before use.
func (b *DetailedReportBuilder) InitDetailedReportBuilder(overrides DetailedReportBuilderOverrides) {
	b.overrides = overrides
}

// detailedReportBuilderOverrides returns the registered overrides, panicking
// when the concrete builder failed to call InitDetailedReportBuilder.
func (b *DetailedReportBuilder) detailedReportBuilderOverrides() DetailedReportBuilderOverrides {
	if b.overrides == nil {
		panic("executor: DetailedReportBuilder used without InitDetailedReportBuilder")
	}
	return b.overrides
}

// Build builds the XmlDetailedReport. Port of build().
func (b *DetailedReportBuilder) Build() *jaxb.XmlDetailedReport {
	detailedReport := b.Init()

	bbbs := b.executeAllBasicBuildingBlocks()
	detailedReport.BasicBuildingBlocks = append(detailedReport.BasicBuildingBlocks, b.BasicBuildingBlocksInOrder(bbbs)...)

	// Init POE
	poe := vpfswatsp.NewPOEExtraction()
	poe.Init(b.DiagnosticData, b.CurrentTime)

	b.detailedReportBuilderOverrides().ExecuteValidation(detailedReport, bbbs, poe)

	if b.includeSemantics {
		b.collectReportIndications(detailedReport)
		b.addSemantics(detailedReport)
	}

	return detailedReport
}

// ExecuteValidation performs validation for the given tokens. Port of the
// protected executeValidation(XmlDetailedReport, Map, POEExtraction).
func (b *DetailedReportBuilder) ExecuteValidation(detailedReport *jaxb.XmlDetailedReport,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, poe *vpfswatsp.POEExtraction) {
	tlAnalysis := detailedReport.TLAnalysis

	attachedTimestamps := make(map[string]struct{})
	attachedEvidenceRecords := make(map[string]struct{})

	timestampValidations := make(map[string]*jaxb.XmlTimestamp)
	evidenceRecordValidations := make(map[string]*jaxb.XmlEvidenceRecord)

	if b.validationLevel == enumerations.ValidationLevel_ARCHIVAL_DATA {
		evidenceRecordsValidationBlock := b.executeEvidenceRecordsValidations(bbbs, tlAnalysis, poe)
		for id, value := range evidenceRecordsValidationBlock.EvidenceRecordValidations() {
			evidenceRecordValidations[id] = value
		}

		erTimestampValidations := evidenceRecordsValidationBlock.TimestampValidations()
		for id, value := range erTimestampValidations {
			timestampValidations[id] = value
			attachedTimestamps[id] = struct{}{}
		}
	}

	if b.validationLevel != enumerations.ValidationLevel_BASIC_SIGNATURES {
		nonEvidenceRecordTimestamps := b.DiagnosticData.NonEvidenceRecordTimestamps()
		for id, value := range b.executeTimestampsValidation(
			nonEvidenceRecordTimestamps, bbbs, evidenceRecordValidations, tlAnalysis, poe, attachedEvidenceRecords) {
			timestampValidations[id] = value
		}
	}

	for _, signature := range b.DiagnosticData.Signatures() {
		signatureAnalysis := &jaxb.XmlSignature{}
		id := signature.Id()
		signatureAnalysis.Id = &id
		if signature.IsCounterSignature() {
			counterSignature := true
			signatureAnalysis.CounterSignature = &counterSignature
		}

		if b.validationLevel != enumerations.ValidationLevel_BASIC_SIGNATURES {
			for _, timestampId := range signature.TimestampIdsList() {
				attachedTimestamps[timestampId] = struct{}{}
			}
			for _, sigTimestamp := range signature.TimestampList() {
				signatureAnalysis.Timestamp = append(signatureAnalysis.Timestamp, timestampValidations[sigTimestamp.Id()])
			}
		}

		// Java types the local as XmlConstraintsConclusionWithProofOfExistence;
		// the Go model has no such supertype, so the port keeps the embedded
		// content struct the qualification block needs, together with the
		// proof-of-existence carried by the concrete result.
		validationBasic := b.executeBasicValidation(signatureAnalysis, signature, signatureAnalysis.Timestamp, bbbs)
		validation := &jaxb.XmlConstraintsConclusionWithProofOfExistence{
			XmlConstraintsConclusionWithProofOfExistenceContent: validationBasic.XmlConstraintsConclusionWithProofOfExistenceContent,
			XmlConstraintsConclusionAttrs:                       validationBasic.XmlConstraintsConclusionAttrs,
		}

		if b.validationLevel == enumerations.ValidationLevel_LONG_TERM_DATA {
			validationLongTerm := b.executeLongTermValidation(signatureAnalysis, signature, bbbs)
			validation = &jaxb.XmlConstraintsConclusionWithProofOfExistence{
				XmlConstraintsConclusionWithProofOfExistenceContent: validationLongTerm.XmlConstraintsConclusionWithProofOfExistenceContent,
				XmlConstraintsConclusionAttrs:                       validationLongTerm.XmlConstraintsConclusionAttrs,
			}

		} else if b.validationLevel == enumerations.ValidationLevel_ARCHIVAL_DATA {
			for _, evidenceRecordId := range signature.EvidenceRecordIdsList() {
				attachedEvidenceRecords[evidenceRecordId] = struct{}{}
			}
			for _, sigEvidenceRecord := range signature.EvidenceRecords() {
				signatureAnalysis.EvidenceRecord = append(signatureAnalysis.EvidenceRecord,
					evidenceRecordValidations[sigEvidenceRecord.Id()])
			}

			b.executeLongTermValidation(signatureAnalysis, signature, bbbs)
			validationArchival := b.executeArchiveValidation(signatureAnalysis, signature, bbbs, poe)
			validation = &jaxb.XmlConstraintsConclusionWithProofOfExistence{
				XmlConstraintsConclusionWithProofOfExistenceContent: validationArchival.XmlConstraintsConclusionWithProofOfExistenceContent,
				XmlConstraintsConclusionAttrs:                       validationArchival.XmlConstraintsConclusionAttrs,
			}
		}

		if b.Policy.EIDASConstraintPresent() {

			// Signature qualification
			signingCertificate := signature.SigningCertificate()
			if signingCertificate != nil {
				qualificationBlock := qualification.NewSignatureQualificationBlock(
					b.I18nProvider, validation, signingCertificate, tlAnalysis)
				signatureAnalysis.ValidationSignatureQualification = qualificationBlock.Execute()
			}

		}

		signatureAnalysis.Conclusion = b.finalConclusion(validation.Conclusion)

		detailedReport.SignatureOrTimestampOrEvidenceRecord =
			append(detailedReport.SignatureOrTimestampOrEvidenceRecord, signatureAnalysis)
	}

	if b.validationLevel == enumerations.ValidationLevel_ARCHIVAL_DATA {
		for _, evidenceRecord := range b.DiagnosticData.EvidenceRecords() {
			if _, attached := attachedEvidenceRecords[evidenceRecord.Id()]; attached {
				continue
			}
			detailedReport.SignatureOrTimestampOrEvidenceRecord =
				append(detailedReport.SignatureOrTimestampOrEvidenceRecord, evidenceRecordValidations[evidenceRecord.Id()])
		}
	}

	if b.validationLevel != enumerations.ValidationLevel_BASIC_SIGNATURES {
		for _, timestamp := range b.DiagnosticData.TimestampList() {
			if _, attached := attachedTimestamps[timestamp.Id()]; attached {
				continue
			}
			detailedReport.SignatureOrTimestampOrEvidenceRecord =
				append(detailedReport.SignatureOrTimestampOrEvidenceRecord, timestampValidations[timestamp.Id()])
		}
	}
}

// executeBasicValidation is the port of the private
// executeBasicValidation(XmlSignature, SignatureWrapper, List, Map).
func (b *DetailedReportBuilder) executeBasicValidation(signatureAnalysis *jaxb.XmlSignature,
	signature *diagnostic.SignatureWrapper, xmlTimestamps []*jaxb.XmlTimestamp,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) *jaxb.XmlValidationProcessBasicSignature {
	vpfbsProcess := vpfbs.NewBasicSignatureValidationProcess(
		b.I18nProvider, b.DiagnosticData, signature, xmlTimestamps, bbbs)
	bs := vpfbsProcess.Execute()
	signatureAnalysis.ValidationProcessBasicSignature = bs
	return bs
}

// executeTimestampsValidation is the port of the private
// executeTimestampsValidation(List, Map, Map, List, POEExtraction, Set).
func (b *DetailedReportBuilder) executeTimestampsValidation(
	timestamps []*diagnostic.TimestampWrapper, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	evidenceRecordValidations map[string]*jaxb.XmlEvidenceRecord, tlAnalysis []*jaxb.XmlTLAnalysis,
	poe *vpfswatsp.POEExtraction, attachedEvidenceRecords map[string]struct{}) map[string]*jaxb.XmlTimestamp {
	allTimestampValidationBlock := vpftsp.NewTimestampsValidationBlock(
		b.I18nProvider, timestamps, b.DiagnosticData, b.Policy, b.CurrentTime, bbbs, evidenceRecordValidations,
		tlAnalysis, b.validationLevel, poe)
	for _, timestampWrapper := range timestamps {
		for _, evidenceRecord := range timestampWrapper.EvidenceRecords() {
			attachedEvidenceRecords[evidenceRecord.Id()] = struct{}{}
		}
	}
	return allTimestampValidationBlock.Execute()
}

// executeLongTermValidation is the port of the private
// executeLongTermValidation(XmlSignature, SignatureWrapper, Map).
func (b *DetailedReportBuilder) executeLongTermValidation(signatureAnalysis *jaxb.XmlSignature,
	signature *diagnostic.SignatureWrapper,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) *jaxb.XmlValidationProcessLongTermData {
	vpfltvd := vpfltvdsig.NewValidationProcessForSignaturesWithLongTermValidationData(
		b.I18nProvider, signatureAnalysis, b.DiagnosticData, signature, bbbs, b.Policy, b.CurrentTime)
	vpfltvdResult := vpfltvd.Execute()
	signatureAnalysis.ValidationProcessLongTermData = vpfltvdResult
	return vpfltvdResult
}

// executeArchiveValidation is the port of the private
// executeArchiveValidation(XmlSignature, SignatureWrapper, Map, POEExtraction).
func (b *DetailedReportBuilder) executeArchiveValidation(signatureAnalysis *jaxb.XmlSignature,
	signature *diagnostic.SignatureWrapper, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	poe *vpfswatsp.POEExtraction) *jaxb.XmlValidationProcessArchivalData {
	vpfswad := vpfswatsp.NewValidationProcessForSignaturesWithArchivalData(
		b.I18nProvider, signatureAnalysis, signature, b.DiagnosticData, bbbs, b.Policy, b.CurrentTime, poe)
	vpfswadResult := vpfswad.Execute()
	signatureAnalysis.ValidationProcessArchivalData = vpfswadResult
	return vpfswadResult
}

// executeEvidenceRecordsValidations is the port of the private
// executeEvidenceRecordsValidations(Map, List, POEExtraction).
func (b *DetailedReportBuilder) executeEvidenceRecordsValidations(
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tlAnalysis []*jaxb.XmlTLAnalysis,
	poe *vpfswatsp.POEExtraction) *evidencerecord.EvidenceRecordsValidationBlock {
	evidenceRecordsValidationBlock := evidencerecord.NewEvidenceRecordsValidationBlock(
		b.I18nProvider, b.DiagnosticData, b.Policy, b.CurrentTime, bbbs, tlAnalysis, b.validationLevel, poe)
	evidenceRecordsValidationBlock.Execute()
	return evidenceRecordsValidationBlock
}

// executeAllBasicBuildingBlocks is the port of the private
// executeAllBasicBuildingBlocks().
func (b *DetailedReportBuilder) executeAllBasicBuildingBlocks() map[string]*jaxb.XmlBasicBuildingBlocks {
	bbbs := make(map[string]*jaxb.XmlBasicBuildingBlocks)
	switch b.validationLevel {
	case enumerations.ValidationLevel_ARCHIVAL_DATA:
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllRevocationData()), enumerations.Context_REVOCATION, bbbs)
		Process(&b.AbstractDetailedReportBuilder, b.DiagnosticData.TimestampList(), enumerations.Context_TIMESTAMP, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllSignatures()), enumerations.Context_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllCounterSignatures()), enumerations.Context_COUNTER_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllKeyBindingSignatures()), enumerations.Context_KEY_BINDING_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllEAARevocationTokens()), enumerations.Context_EAA_REVOCATION, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllEAA()), enumerations.Context_EAA, bbbs)
	case enumerations.ValidationLevel_LONG_TERM_DATA:
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllRevocationData()), enumerations.Context_REVOCATION, bbbs)
		Process(&b.AbstractDetailedReportBuilder, b.DiagnosticData.NonEvidenceRecordTimestamps(), enumerations.Context_TIMESTAMP, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllSignatures()), enumerations.Context_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllCounterSignatures()), enumerations.Context_COUNTER_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllKeyBindingSignatures()), enumerations.Context_KEY_BINDING_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllEAARevocationTokens()), enumerations.Context_EAA_REVOCATION, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllEAA()), enumerations.Context_EAA, bbbs)
	case enumerations.ValidationLevel_TIMESTAMPS:
		Process(&b.AbstractDetailedReportBuilder, b.DiagnosticData.NonEvidenceRecordTimestamps(), enumerations.Context_TIMESTAMP, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllSignatures()), enumerations.Context_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllCounterSignatures()), enumerations.Context_COUNTER_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllKeyBindingSignatures()), enumerations.Context_KEY_BINDING_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllEAARevocationTokens()), enumerations.Context_EAA_REVOCATION, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllEAA()), enumerations.Context_EAA, bbbs)
	case enumerations.ValidationLevel_BASIC_SIGNATURES:
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllSignatures()), enumerations.Context_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllCounterSignatures()), enumerations.Context_COUNTER_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllKeyBindingSignatures()), enumerations.Context_KEY_BINDING_SIGNATURE, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllEAARevocationTokens()), enumerations.Context_EAA_REVOCATION, bbbs)
		Process(&b.AbstractDetailedReportBuilder, JavaHashSetOrder(b.DiagnosticData.AllEAA()), enumerations.Context_EAA, bbbs)
	default:
		panic(fmt.Sprintf("Unsupported validation level %s", b.validationLevel))
	}
	return bbbs
}

// finalConclusion is the port of the private
// getFinalConclusion(XmlConstraintsConclusion).
func (b *DetailedReportBuilder) finalConclusion(conclusion *jaxb.XmlConclusion) *jaxb.XmlConclusion {
	xmlConclusion := &jaxb.XmlConclusion{}
	indication := b.finalIndication(conclusion.Indication.Indication())
	xmlConclusion.Indication = jaxb.IndicationValue(indication)
	xmlConclusion.SubIndication = conclusion.SubIndication
	return xmlConclusion
}

// finalIndication is the port of the private getFinalIndication(Indication).
func (b *DetailedReportBuilder) finalIndication(highestIndication enumerations.Indication) enumerations.Indication {
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

// collectReportIndications is the port of the private
// collectIndications(XmlDetailedReport).
func (b *DetailedReportBuilder) collectReportIndications(detailedReport *jaxb.XmlDetailedReport) {
	for _, xmlObject := range detailedReport.SignatureOrTimestampOrEvidenceRecord {
		switch typed := xmlObject.(type) {
		case *jaxb.XmlSignature:
			b.collectSignatureIndications(typed)
		case *jaxb.XmlTimestamp:
			b.collectTimestampIndications(typed)
		case *jaxb.XmlEvidenceRecord:
			b.collectEvidenceRecordIndications(typed)
		}
	}
	for _, bbb := range detailedReport.BasicBuildingBlocks {
		b.collectBBBIndications(bbb)
	}
}

// collectSignatureIndications is the port of the private
// collectIndications(XmlSignature).
func (b *DetailedReportBuilder) collectSignatureIndications(xmlSignature *jaxb.XmlSignature) {
	b.collectConclusionIndications(xmlSignature.Conclusion)
	if xmlSignature.ValidationProcessBasicSignature != nil {
		b.collectConclusionIndications(xmlSignature.ValidationProcessBasicSignature.Conclusion)
	}
	if xmlSignature.ValidationProcessLongTermData != nil {
		b.collectConclusionIndications(xmlSignature.ValidationProcessLongTermData.Conclusion)
	}
	if xmlSignature.ValidationProcessArchivalData != nil {
		b.collectConclusionIndications(xmlSignature.ValidationProcessArchivalData.Conclusion)
	}
	for _, xmlTimestamp := range xmlSignature.Timestamp {
		b.collectTimestampIndications(xmlTimestamp)
	}
}

// collectTimestampIndications is the port of the private
// collectIndications(XmlTimestamp).
func (b *DetailedReportBuilder) collectTimestampIndications(xmlTimestamp *jaxb.XmlTimestamp) {
	if xmlTimestamp.ValidationProcessBasicTimestamp != nil {
		b.collectConclusionIndications(xmlTimestamp.ValidationProcessBasicTimestamp.Conclusion)
	}
}

// collectEvidenceRecordIndications is the port of the private
// collectIndications(XmlEvidenceRecord).
func (b *DetailedReportBuilder) collectEvidenceRecordIndications(xmlEvidenceRecord *jaxb.XmlEvidenceRecord) {
	if xmlEvidenceRecord.ValidationProcessEvidenceRecord != nil {
		b.collectConclusionIndications(xmlEvidenceRecord.ValidationProcessEvidenceRecord.Conclusion)
	}
}

// collectBBBIndications is the port of the private
// collectIndications(XmlBasicBuildingBlocks).
func (b *DetailedReportBuilder) collectBBBIndications(bbb *jaxb.XmlBasicBuildingBlocks) {
	if bbb.FC != nil {
		b.collectConclusionIndications(bbb.FC.Conclusion)
	}
	if bbb.ISC != nil {
		b.collectConclusionIndications(bbb.ISC.Conclusion)
	}
	if bbb.VCI != nil {
		b.collectConclusionIndications(bbb.VCI.Conclusion)
	}
	if bbb.XCV != nil {
		b.collectConclusionIndications(bbb.XCV.Conclusion)
		for _, subXCV := range bbb.XCV.SubXCV {
			b.collectConclusionIndications(subXCV.Conclusion)
			if subXCV.RFC != nil {
				b.collectConclusionIndications(subXCV.RFC.Conclusion)
			}
			if subXCV.CRS != nil {
				b.collectConclusionIndications(subXCV.CRS.Conclusion)
				for _, rac := range subXCV.CRS.RAC {
					b.collectConclusionIndications(rac.Conclusion)
				}
			}
		}
	}
	if bbb.CV != nil {
		b.collectConclusionIndications(bbb.CV.Conclusion)
	}
	if bbb.SAV != nil {
		b.collectConclusionIndications(bbb.SAV.Conclusion)
	}
	if bbb.PSV != nil {
		b.collectConclusionIndications(bbb.PSV.Conclusion)
	}
	if bbb.PCV != nil {
		b.collectConclusionIndications(bbb.PCV.Conclusion)
	}
	if bbb.VTS != nil {
		b.collectConclusionIndications(bbb.VTS.Conclusion)
	}
}

// collectConclusionIndications is the port of the private
// collectIndications(XmlConclusion).
func (b *DetailedReportBuilder) collectConclusionIndications(xmlConclusion *jaxb.XmlConclusion) {
	if xmlConclusion != nil {
		indication := xmlConclusion.Indication.Indication()
		if indication != "" {
			b.allIndications[indication] = struct{}{}

			subIndication := xmlConclusion.SubIndication
			if subIndication != nil {
				b.allSubIndications[subIndication.SubIndication()] = struct{}{}
			}
		}
	}
}

// addSemantics is the port of the private addSemantics(XmlDetailedReport).
// The enumeration declaration order replaces Java's HashSet iteration order;
// see the file header.
func (b *DetailedReportBuilder) addSemantics(detailedReport *jaxb.XmlDetailedReport) {

	for _, indication := range enumerations.IndicationValues() {
		if _, used := b.allIndications[indication]; !used {
			continue
		}
		semantic := &jaxb.XmlSemantic{}
		semantic.Key = string(indication)
		tag, _ := i18n.MessageTagGetSemantic(string(indication))
		semantic.Value = b.I18nProvider.GetMessage(tag)
		detailedReport.Semantic = append(detailedReport.Semantic, semantic)
	}

	for _, subIndication := range enumerations.SubIndicationValues() {
		if _, used := b.allSubIndications[subIndication]; !used {
			continue
		}
		semantic := &jaxb.XmlSemantic{}
		semantic.Key = string(subIndication)
		tag, _ := i18n.MessageTagGetSemantic(string(subIndication))
		semantic.Value = b.I18nProvider.GetMessage(tag)
		detailedReport.Semantic = append(detailedReport.Semantic, semantic)
	}

}
