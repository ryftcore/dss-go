// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftsp/TimestampsValidationBlock.java (DSS 6.5.RC1).
//
// Evidence-record support: Java's protected constructor defaults
// evidenceRecordValidations to Collections.emptyMap() with a "TODO: implement
// support" comment; the Go port keeps the same gap (see
// newTimestampsValidationBlock below).
package vpftsp

import (
	"sort"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process/qualification"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfswatsp"
	"github.com/ryftcore/dss-go/dss/validation/process/vpftspwatsp"
)

// TimestampsValidationBlockOverrides declares the overridable protected methods
// of TimestampsValidationBlock that the base implementation calls back into. A
// concrete block registers itself through InitTimestampsValidationBlock; every
// method it does not define is supplied by the embedded TimestampsValidationBlock
// through ordinary Go method promotion. vpfswatsp/evidencerecord's
// EvidenceRecordTimestampsValidationBlock overrides both.
type TimestampsValidationBlockOverrides interface {
	// Timestamps returns a list of time-stamp tokens to be validated. Port of
	// the protected getTimestamps().
	Timestamps() []*diagnostic.TimestampWrapper
	// Poe returns the POE object for the timestamp validation. Port of the
	// protected getPoe(TimestampWrapper).
	Poe(timestamp *diagnostic.TimestampWrapper) *vpfswatsp.POEExtraction
}

// TimestampsValidationBlock is used to perform validation of all available
// timestamps, as well as to extract POE information for valid entries.
type TimestampsValidationBlock struct {
	// i18nProvider is the i18n provider.
	i18nProvider *i18n.I18nProvider

	// TimestampList is the list of time-stamps to be validated. Exported
	// because Java declares the field protected; it cannot keep Java's name
	// "timestamps" capitalised, since the overridable getTimestamps() takes
	// that identifier as the method Timestamps().
	TimestampList []*diagnostic.TimestampWrapper

	// overrides points back at the concrete block; see
	// InitTimestampsValidationBlock.
	overrides TimestampsValidationBlockOverrides

	// diagnosticData is the DiagnosticData to use.
	diagnosticData *diagnostic.DiagnosticData

	// Policy is the validation policy. Exported because Java declares the
	// field protected.
	Policy policy.ValidationPolicy

	// CurrentTime is the validation time. Exported because Java declares the
	// field protected.
	CurrentTime time.Time

	// bbbs is the map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// evidenceRecordValidations is the map of processed evidence records.
	evidenceRecordValidations map[string]*jaxb.XmlEvidenceRecord

	// tlAnalysis is the list of Trusted List validations.
	tlAnalysis []*jaxb.XmlTLAnalysis

	// validationLevel is the target highest validation level.
	validationLevel enumerations.ValidationLevel

	// poe contains the list of all POEs.
	poe *vpfswatsp.POEExtraction
}

// NewTimestampsValidationBlock is the default constructor. Port of
// TimestampsValidationBlock(I18nProvider, List, DiagnosticData, ValidationPolicy, Date, Map, Map, List, ValidationLevel, POEExtraction).
func NewTimestampsValidationBlock(i18nProvider *i18n.I18nProvider, timestamps []*diagnostic.TimestampWrapper,
	diagnosticData *diagnostic.DiagnosticData, validationPolicy policy.ValidationPolicy, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, evidenceRecordValidations map[string]*jaxb.XmlEvidenceRecord,
	tlAnalysis []*jaxb.XmlTLAnalysis, validationLevel enumerations.ValidationLevel,
	poe *vpfswatsp.POEExtraction) *TimestampsValidationBlock {
	b := &TimestampsValidationBlock{}
	b.InitTimestampsValidationBlockState(i18nProvider, timestamps, diagnosticData, validationPolicy, currentTime,
		bbbs, evidenceRecordValidations, tlAnalysis, validationLevel, poe)
	b.InitTimestampsValidationBlock(b)
	return b
}

// InitTimestampsValidationBlockState wires the shared state, the way the Java
// constructor body does. A subclass calls it before
// InitTimestampsValidationBlock, in place of the Java super(...) call.
func (b *TimestampsValidationBlock) InitTimestampsValidationBlockState(i18nProvider *i18n.I18nProvider,
	timestamps []*diagnostic.TimestampWrapper, diagnosticData *diagnostic.DiagnosticData,
	validationPolicy policy.ValidationPolicy, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, evidenceRecordValidations map[string]*jaxb.XmlEvidenceRecord,
	tlAnalysis []*jaxb.XmlTLAnalysis, validationLevel enumerations.ValidationLevel,
	poe *vpfswatsp.POEExtraction) {
	b.i18nProvider = i18nProvider
	b.TimestampList = timestamps
	b.diagnosticData = diagnosticData
	b.Policy = validationPolicy
	b.CurrentTime = currentTime
	b.bbbs = bbbs
	b.evidenceRecordValidations = evidenceRecordValidations
	b.tlAnalysis = tlAnalysis
	b.validationLevel = validationLevel
	b.poe = poe
}

// InitTimestampsValidationBlockStateWithoutPOE wires the shared state of the
// protected Java constructor, which builds its own POEExtraction and leaves the
// evidence-record validations empty.
func (b *TimestampsValidationBlock) InitTimestampsValidationBlockStateWithoutPOE(i18nProvider *i18n.I18nProvider,
	timestamps []*diagnostic.TimestampWrapper, diagnosticData *diagnostic.DiagnosticData,
	validationPolicy policy.ValidationPolicy, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tlAnalysis []*jaxb.XmlTLAnalysis,
	validationLevel enumerations.ValidationLevel) {
	poe := vpfswatsp.NewPOEExtraction()
	poe.Init(diagnosticData, currentTime)
	// evidenceRecordValidations stays nil (Java: Collections.emptyMap();
	// "TODO : implement support").
	b.InitTimestampsValidationBlockState(i18nProvider, timestamps, diagnosticData, validationPolicy, currentTime,
		bbbs, nil, tlAnalysis, validationLevel, poe)
}

// InitTimestampsValidationBlock registers the concrete block with its base so
// that the base can dispatch to the overridden methods. It must be called
// exactly once, by the concrete block's constructor, before Execute.
func (b *TimestampsValidationBlock) InitTimestampsValidationBlock(overrides TimestampsValidationBlockOverrides) {
	b.overrides = overrides
}

// blockOverrides returns the registered overrides, panicking when the concrete
// block forgot to call InitTimestampsValidationBlock.
func (b *TimestampsValidationBlock) blockOverrides() TimestampsValidationBlockOverrides {
	if b.overrides == nil {
		panic("TimestampsValidationBlock was not initialised: the concrete block must call InitTimestampsValidationBlock in its constructor")
	}
	return b.overrides
}

// NewTimestampsValidationBlockWithoutPOE is the constructor without POE. Port
// of the protected
// TimestampsValidationBlock(I18nProvider, List, DiagnosticData, ValidationPolicy, Date, Map, List, ValidationLevel).
func NewTimestampsValidationBlockWithoutPOE(i18nProvider *i18n.I18nProvider, timestamps []*diagnostic.TimestampWrapper,
	diagnosticData *diagnostic.DiagnosticData, validationPolicy policy.ValidationPolicy, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tlAnalysis []*jaxb.XmlTLAnalysis,
	validationLevel enumerations.ValidationLevel) *TimestampsValidationBlock {
	b := &TimestampsValidationBlock{}
	b.InitTimestampsValidationBlockStateWithoutPOE(i18nProvider, timestamps, diagnosticData, validationPolicy,
		currentTime, bbbs, tlAnalysis, validationLevel)
	b.InitTimestampsValidationBlock(b)
	return b
}

// Execute performs validation of timestamps, but also fills the POEExtraction
// object for valid timestamps. Port of execute().
func (b *TimestampsValidationBlock) Execute() map[string]*jaxb.XmlTimestamp {
	result := make(map[string]*jaxb.XmlTimestamp)

	for _, newestTimestamp := range b.blockOverrides().Timestamps() {
		xmlTimestamp := b.buildXmlTimestamp(newestTimestamp)
		result[newestTimestamp.Id()] = xmlTimestamp
	}

	return result
}

// Timestamps returns a list of time-stamp tokens to be validated. Port of
// getTimestamps(): Java sorts by production time, reversed (newest first).
func (b *TimestampsValidationBlock) Timestamps() []*diagnostic.TimestampWrapper {
	timestampList := make([]*diagnostic.TimestampWrapper, len(b.TimestampList))
	copy(timestampList, b.TimestampList)
	sort.SliceStable(timestampList, func(i, j int) bool {
		ti, tj := timestampList[i].ProductionTime(), timestampList[j].ProductionTime()
		if ti == nil || tj == nil {
			return false
		}
		// reversed: newest first
		return ti.After(*tj)
	})
	return timestampList
}

// buildXmlTimestamp ports the private buildXmlTimestamp(TimestampWrapper, Map, List).
func (b *TimestampsValidationBlock) buildXmlTimestamp(timestamp *diagnostic.TimestampWrapper) *jaxb.XmlTimestamp {
	xmlTimestamp := &jaxb.XmlTimestamp{}
	id := timestamp.Id()
	xmlTimestamp.Id = &id

	currentPOE := b.blockOverrides().Poe(timestamp)

	vpftsp := NewTimestampBasicValidationProcess(b.i18nProvider, b.diagnosticData, timestamp, b.bbbs)
	validationProcessBasicTimestamp := vpftsp.Execute()
	xmlTimestamp.ValidationProcessBasicTimestamp = validationProcessBasicTimestamp

	conclusion := validationProcessBasicTimestamp.Conclusion

	// Timestamp qualification
	if b.Policy.EIDASConstraintPresent() {
		timestampQualificationBlock := qualification.NewTimestampQualificationBlock(
			b.i18nProvider, timestamp, b.tlAnalysis, currentPOE)
		xmlTimestamp.ValidationTimestampQualification = timestampQualificationBlock.Execute()
	}

	if enumerations.ValidationLevelArchivalData == b.validationLevel {
		for _, sigEvidenceRecord := range timestamp.EvidenceRecords() {
			xmlTimestamp.EvidenceRecord = append(xmlTimestamp.EvidenceRecord, b.evidenceRecordValidations[sigEvidenceRecord.Id()])
		}

		vpftspwatst := vpftspwatsp.NewValidationProcessForTimestampsWithArchivalData(
			b.i18nProvider, timestamp, validationProcessBasicTimestamp, b.bbbs, b.evidenceRecordValidations,
			b.CurrentTime, b.Policy, currentPOE)
		validationProcessTimestampArchivalData := vpftspwatst.Execute()

		// extract POE for valid time-stamps
		if validationProcessTimestampArchivalData.Conclusion != nil &&
			enumerations.IndicationPassed == validationProcessTimestampArchivalData.Conclusion.Indication.Indication() {
			currentPOE.ExtractPOE(timestamp)
		}

		xmlTimestamp.ValidationProcessArchivalDataTimestamp = validationProcessTimestampArchivalData
		conclusion = validationProcessTimestampArchivalData.Conclusion
	}

	xmlTimestamp.Conclusion = conclusion
	return xmlTimestamp
}

// Poe returns the POE object for the timestamp validation. Port of
// getPoe(TimestampWrapper).
func (b *TimestampsValidationBlock) Poe(timestamp *diagnostic.TimestampWrapper) *vpfswatsp.POEExtraction {
	return b.poe
}
