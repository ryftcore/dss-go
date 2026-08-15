// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftsp/TimestampsValidationBlock.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY (LTVB, phase 8e): this file uses three packages ported by
// sibling porters in the same phase, not yet present when this file was
// written:
//
//   - github.com/utain/esig/dss/validation/process/vpfswatsp for POEExtraction
//     (LTVA). This one IS already present with the assumed shape:
//     NewPOEExtraction(), Init(*diagnostic.DiagnosticData, time.Time),
//     ExtractPOE(*diagnostic.TimestampWrapper).
//   - github.com/utain/esig/dss/validation/process/qualification (shared
//     QCERT/QTRUST/QSIG porters) for TimestampQualificationBlock, assumed to
//     have the shape
//     NewTimestampQualificationBlock(*i18n.I18nProvider, *diagnostic.TimestampWrapper, []*jaxb.XmlTLAnalysis, *vpfswatsp.POEExtraction) *TimestampQualificationBlock
//     with an Execute() *jaxb.XmlValidationTimestampQualification method,
//     mirroring Java's TimestampQualificationBlock(I18nProvider, TimestampWrapper, List<XmlTLAnalysis>, POEExtraction).
//   - github.com/utain/esig/dss/validation/process/vpftspwatsp (LTVA) for
//     ValidationProcessForTimestampsWithArchivalData, assumed to have the shape
//     NewValidationProcessForTimestampsWithArchivalData(*i18n.I18nProvider, *diagnostic.TimestampWrapper, *jaxb.XmlValidationProcessBasicTimestamp, map[string]*jaxb.XmlBasicBuildingBlocks, map[string]*jaxb.XmlEvidenceRecord, time.Time, policy.ValidationPolicy, *vpfswatsp.POEExtraction) *ValidationProcessForTimestampsWithArchivalData
//     with an Execute() *jaxb.XmlValidationProcessArchivalDataTimestamp method,
//     mirroring Java's
//     ValidationProcessForTimestampsWithArchivalData(I18nProvider, TimestampWrapper, XmlValidationProcessBasicTimestamp, Map, Map, Date, ValidationPolicy, POEExtraction).
//
// Evidence-record support: Java's protected constructor defaults
// evidenceRecordValidations to Collections.emptyMap() with a "TODO: implement
// support" comment; the Go port keeps the same gap (see
// newTimestampsValidationBlock below).
package vpftsp

import (
	"sort"
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process/qualification"
	"github.com/utain/esig/dss/validation/process/vpfswatsp"
	"github.com/utain/esig/dss/validation/process/vpftspwatsp"
)

// TimestampsValidationBlock is used to perform validation of all available
// timestamps, as well as to extract POE information for valid entries.
type TimestampsValidationBlock struct {
	// i18nProvider is the i18n provider.
	i18nProvider *i18n.I18nProvider

	// Timestamps is the list of time-stamps to be validated. Exported because
	// Java declares the field protected.
	Timestamps []*diagnostic.TimestampWrapper

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
	return &TimestampsValidationBlock{
		i18nProvider:              i18nProvider,
		Timestamps:                timestamps,
		diagnosticData:            diagnosticData,
		Policy:                    validationPolicy,
		CurrentTime:               currentTime,
		bbbs:                      bbbs,
		evidenceRecordValidations: evidenceRecordValidations,
		tlAnalysis:                tlAnalysis,
		validationLevel:           validationLevel,
		poe:                       poe,
	}
}

// NewTimestampsValidationBlockWithoutPOE is the constructor without POE. Port
// of the protected
// TimestampsValidationBlock(I18nProvider, List, DiagnosticData, ValidationPolicy, Date, Map, List, ValidationLevel).
func NewTimestampsValidationBlockWithoutPOE(i18nProvider *i18n.I18nProvider, timestamps []*diagnostic.TimestampWrapper,
	diagnosticData *diagnostic.DiagnosticData, validationPolicy policy.ValidationPolicy, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tlAnalysis []*jaxb.XmlTLAnalysis,
	validationLevel enumerations.ValidationLevel) *TimestampsValidationBlock {
	poe := vpfswatsp.NewPOEExtraction()
	poe.Init(diagnosticData, currentTime)
	return &TimestampsValidationBlock{
		i18nProvider:   i18nProvider,
		Timestamps:     timestamps,
		diagnosticData: diagnosticData,
		Policy:         validationPolicy,
		CurrentTime:    currentTime,
		bbbs:           bbbs,
		// evidenceRecordValidations stays nil (Java: Collections.emptyMap();
		// "TODO : implement support").
		tlAnalysis:      tlAnalysis,
		validationLevel: validationLevel,
		poe:             poe,
	}
}

// Execute performs validation of timestamps, but also fills the POEExtraction
// object for valid timestamps. Port of execute().
func (b *TimestampsValidationBlock) Execute() map[string]*jaxb.XmlTimestamp {
	result := make(map[string]*jaxb.XmlTimestamp)

	for _, newestTimestamp := range b.getTimestamps() {
		xmlTimestamp := b.buildXmlTimestamp(newestTimestamp)
		result[newestTimestamp.Id()] = xmlTimestamp
	}

	return result
}

// getTimestamps returns a list of time-stamp tokens to be validated. Port of
// getTimestamps(): Java sorts by production time, reversed (newest first).
func (b *TimestampsValidationBlock) getTimestamps() []*diagnostic.TimestampWrapper {
	timestampList := make([]*diagnostic.TimestampWrapper, len(b.Timestamps))
	copy(timestampList, b.Timestamps)
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

	currentPOE := b.getPoe(timestamp)

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

	if enumerations.ValidationLevel_ARCHIVAL_DATA == b.validationLevel {
		for _, sigEvidenceRecord := range timestamp.EvidenceRecords() {
			xmlTimestamp.EvidenceRecord = append(xmlTimestamp.EvidenceRecord, b.evidenceRecordValidations[sigEvidenceRecord.Id()])
		}

		vpftspwatst := vpftspwatsp.NewValidationProcessForTimestampsWithArchivalData(
			b.i18nProvider, timestamp, validationProcessBasicTimestamp, b.bbbs, b.evidenceRecordValidations,
			b.CurrentTime, b.Policy, currentPOE)
		validationProcessTimestampArchivalData := vpftspwatst.Execute()

		// extract POE for valid time-stamps
		if validationProcessTimestampArchivalData.Conclusion != nil &&
			enumerations.Indication_PASSED == validationProcessTimestampArchivalData.Conclusion.Indication.Indication() {
			currentPOE.ExtractPOE(timestamp)
		}

		xmlTimestamp.ValidationProcessArchivalDataTimestamp = validationProcessTimestampArchivalData
		conclusion = validationProcessTimestampArchivalData.Conclusion
	}

	xmlTimestamp.Conclusion = conclusion
	return xmlTimestamp
}

// getPoe returns the POE object for the timestamp validation. Port of
// getPoe(TimestampWrapper).
func (b *TimestampsValidationBlock) getPoe(timestamp *diagnostic.TimestampWrapper) *vpfswatsp.POEExtraction {
	return b.poe
}
