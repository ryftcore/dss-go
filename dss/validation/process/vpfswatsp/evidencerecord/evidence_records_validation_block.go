// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/evidencerecord/EvidenceRecordsValidationBlock.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note, and
// evidence_record_timestamps_validation_block.go for the cross-chunk vpftsp
// dependency and the import cycle it exposes.
//
// HASH-ORDER (closed in phase 8f, see Execute()). Java fills
// XmlEvidenceRecord#getTimestamps() from currentTimestampValidations.values(),
// where currentTimestampValidations is the java.util.HashMap that
// TimestampsValidationBlock#execute() returns. That iteration order is the
// String-key hash order of the time-stamp ids - unspecified, yet it decides the
// order of the <Timestamp> elements of the detailed report, which is
// order-sensitive output. The Go port reproduces that order exactly, by
// walking the same insertion order (EvidenceRecordTimestampsValidationBlock#
// getTimestamps(), production time descending) through
// utils.JavaHashMapStringKeyOrder.
//
// The other two maps (timestampValidations, evidenceRecordValidations) are read
// by DetailedReportBuilder through get(id) and keySet() only, so their iteration
// order never reaches an output and a plain Go map is faithful.//
// PACKAGE-BOUNDARY DEVIATION (LTVA, phase 8e): Java's vpfswatsp.evidencerecord
// is a package of its own, distinct from vpfswatsp; the phase 8e layout folds
// the whole vpfswatsp tree into one Go package, but EvidenceRecordTimestampsValidationBlock
// extends vpftsp.TimestampsValidationBlock while vpftsp imports vpfswatsp
// (POEExtraction) - an import cycle Go forbids. The five evidence-record classes
// therefore keep Java's own vpfswatsp/evidencerecord package boundary; nothing
// in vpfswatsp, vpftsp or vpftspwatsp imports them (only the validation
// executor's DetailedReportBuilder does), so the edge only ever points upward.
package evidencerecord

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process/vpfswatsp"
)

// EvidenceRecordsValidationBlock performs validation of all evidence records
// provided to the validator.
type EvidenceRecordsValidationBlock struct {
	// i18nProvider is the i18n provider.
	i18nProvider *i18n.I18nProvider

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.DiagnosticData

	// Policy is the validation policy. Exported because Java declares the
	// field protected.
	Policy policy.ValidationPolicy

	// CurrentTime is the validation time. Exported because Java declares the
	// field protected.
	CurrentTime time.Time

	// bbbs is the map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// tlAnalysis is the list of Trusted List validations.
	tlAnalysis []*jaxb.XmlTLAnalysis

	// validationLevel is the target highest validation level.
	validationLevel enumerations.ValidationLevel

	// poe contains list of all POEs.
	poe *vpfswatsp.POEExtraction

	// timestampValidations is the map of all performed time-stamp validations.
	timestampValidations map[string]*jaxb.XmlTimestamp

	// evidenceRecordValidations is the map of all performed evidence record
	// validations.
	evidenceRecordValidations map[string]*jaxb.XmlEvidenceRecord
}

// NewEvidenceRecordsValidationBlock is the default constructor. Port of
// EvidenceRecordsValidationBlock(I18nProvider, DiagnosticData, ValidationPolicy, Date, Map, List, ValidationLevel, POEExtraction).
func NewEvidenceRecordsValidationBlock(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	validationPolicy policy.ValidationPolicy, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tlAnalysis []*jaxb.XmlTLAnalysis,
	validationLevel enumerations.ValidationLevel, poe *vpfswatsp.POEExtraction) *EvidenceRecordsValidationBlock {
	return &EvidenceRecordsValidationBlock{
		i18nProvider:              i18nProvider,
		diagnosticData:            diagnosticData,
		Policy:                    validationPolicy,
		CurrentTime:               currentTime,
		bbbs:                      bbbs,
		tlAnalysis:                tlAnalysis,
		validationLevel:           validationLevel,
		poe:                       poe,
		timestampValidations:      make(map[string]*jaxb.XmlTimestamp),
		evidenceRecordValidations: make(map[string]*jaxb.XmlEvidenceRecord),
	}
}

// Execute performs validation of evidence records. Port of execute().
func (b *EvidenceRecordsValidationBlock) Execute() {
	for _, evidenceRecord := range b.diagnosticData.EvidenceRecords() {
		evidenceRecordAnalysis := &jaxb.XmlEvidenceRecord{}
		id := evidenceRecord.Id()
		evidenceRecordAnalysis.Id = &id

		allTimestampValidationBlock := NewEvidenceRecordTimestampsValidationBlock(
			b.i18nProvider, evidenceRecord, b.diagnosticData, b.Policy, b.CurrentTime, b.bbbs, b.tlAnalysis,
			b.validationLevel)
		currentTimestampValidations := allTimestampValidationBlock.Execute()
		for timestampId, xmlTimestamp := range currentTimestampValidations {
			b.timestampValidations[timestampId] = xmlTimestamp
		}

		// Java appends currentTimestampValidations.values(), i.e. the values of
		// the HashMap<String, XmlTimestamp> TimestampsValidationBlock#execute()
		// filled with put() while walking getTimestamps() (production time,
		// descending). That iteration order reaches the <Timestamp> element
		// sequence of the marshalled detailed report, so it is byte-compared
		// output and has to be REPRODUCED, not substituted: walk the same
		// insertion order and reorder it the way a java.util.HashMap keyed by
		// those ids iterates. Found by the phase-8f full-corpus report
		// byte-parity run on er-validation/er-valid.xml and
		// er-validation/sig-with-er-valid.xml.
		var timestampInsertionOrder []string
		for _, timestamp := range allTimestampValidationBlock.Timestamps() {
			timestampInsertionOrder = append(timestampInsertionOrder, timestamp.Id())
		}
		for _, timestampId := range utils.JavaHashMapStringKeyOrder(timestampInsertionOrder) {
			if xmlTimestamp, ok := currentTimestampValidations[timestampId]; ok {
				evidenceRecordAnalysis.Timestamp = append(evidenceRecordAnalysis.Timestamp, xmlTimestamp)
			}
		}

		// Java hands EvidenceRecordValidationProcess the map's values() view,
		// a collection distinct from the analysis' list; the process only reads
		// it by looking a time-stamp id up linearly, so the same elements in
		// the deterministic order above are equivalent.
		xmlTimestamps := make([]*jaxb.XmlTimestamp, len(evidenceRecordAnalysis.Timestamp))
		copy(xmlTimestamps, evidenceRecordAnalysis.Timestamp)

		ervp := NewEvidenceRecordValidationProcess(b.i18nProvider, b.diagnosticData, evidenceRecord, xmlTimestamps,
			b.bbbs, b.Policy, b.CurrentTime)
		validationProcessEvidenceRecord := ervp.Execute()
		evidenceRecordAnalysis.ValidationProcessEvidenceRecord = validationProcessEvidenceRecord

		conclusion := validationProcessEvidenceRecord.Conclusion
		evidenceRecordAnalysis.Conclusion = conclusion

		if conclusion != nil && enumerations.Indication_PASSED == conclusion.Indication.Indication() {
			b.poe.ExtractEvidenceRecordPOE(evidenceRecord)
		}

		b.evidenceRecordValidations[evidenceRecord.Id()] = evidenceRecordAnalysis
	}
}

// TimestampValidations returns a map of performed time-stamp validations: a map
// of time-stamp identifiers and their corresponding validations. Port of
// getTimestampValidations().
func (b *EvidenceRecordsValidationBlock) TimestampValidations() map[string]*jaxb.XmlTimestamp {
	return b.timestampValidations
}

// EvidenceRecordValidations returns a map of performed evidence record
// validations: a map of evidence record identifiers and their corresponding
// validations. Port of getEvidenceRecordValidations().
func (b *EvidenceRecordsValidationBlock) EvidenceRecordValidations() map[string]*jaxb.XmlEvidenceRecord {
	return b.evidenceRecordValidations
}
