// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/evidencerecord/checks/EvidenceRecordSignedAndTimestampedFilesCoveredCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
//
// Package placement deviation: Java's vpfswatsp.evidencerecord
// is a package of its own, distinct from vpfswatsp. Everything else in the
// vpfswatsp tree folds into one Go package, but EvidenceRecordTimestampsValidationBlock
// extends vpftsp.TimestampsValidationBlock while vpftsp imports vpfswatsp
// (POEExtraction) - an import cycle Go forbids. The five evidence-record classes
// therefore keep Java's own vpfswatsp/evidencerecord package boundary; nothing
// in vpfswatsp, vpftsp or vpftspwatsp imports them (only the validation
// executor's DetailedReportBuilder does), so the edge only ever points upward.
package evidencerecord

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/fc"
)

// EvidenceRecordSignedAndTimestampedFilesCoveredCheck verifies whether all
// signed and/or time-asserted file objects are subsequently covered by the
// evidence record.
type EvidenceRecordSignedAndTimestampedFilesCoveredCheck struct {
	*fc.AbstractSignedAndTimestampedFilesCoveredCheck[*jaxb.XmlValidationProcessEvidenceRecord]

	// evidenceRecordWrapper is the evidence record to be validated.
	evidenceRecordWrapper *diagnostic.EvidenceRecordWrapper
}

// NewEvidenceRecordSignedAndTimestampedFilesCoveredCheck is the default
// constructor. Port of
// EvidenceRecordSignedAndTimestampedFilesCoveredCheck(I18nProvider, XmlValidationProcessEvidenceRecord, DiagnosticData, EvidenceRecordWrapper, LevelRule).
func NewEvidenceRecordSignedAndTimestampedFilesCoveredCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessEvidenceRecord], diagnosticData *diagnostic.DiagnosticData,
	evidenceRecordWrapper *diagnostic.EvidenceRecordWrapper,
	constraint policy.LevelRule) *EvidenceRecordSignedAndTimestampedFilesCoveredCheck {
	c := &EvidenceRecordSignedAndTimestampedFilesCoveredCheck{
		AbstractSignedAndTimestampedFilesCoveredCheck: &fc.AbstractSignedAndTimestampedFilesCoveredCheck[*jaxb.XmlValidationProcessEvidenceRecord]{},
		evidenceRecordWrapper:                         evidenceRecordWrapper,
	}
	c.InitAbstractSignedAndTimestampedFilesCoveredCheck(i18nProvider, result, diagnosticData,
		evidenceRecordWrapper.Filename(), constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of the overridden process(), whose
// super.process() call is the embedded base's.
func (c *EvidenceRecordSignedAndTimestampedFilesCoveredCheck) Process() bool {
	if !c.AbstractSignedAndTimestampedFilesCoveredCheck.Process() {
		return false
	}

	coveredDocumentEntries := c.coveredDocumentEntries()
	return c.CheckManifestFilesCovered(coveredDocumentEntries)
}

// coveredDocumentEntries ports the private getCoveredDocumentEntries(): the
// generated DocumentName and Filename members are pointers/strings, whose nil
// and "" are Java's null (SignatureWrapper#getFilename() returns the plain
// string the wrapper holds, empty when absent).
func (c *EvidenceRecordSignedAndTimestampedFilesCoveredCheck) coveredDocumentEntries() []string {
	result := make([]string, 0)
	for _, digestMatcher := range c.evidenceRecordWrapper.DigestMatchers() {
		if digestMatcher.DocumentName != nil {
			result = append(result, *digestMatcher.DocumentName)
		}
	}
	if c.evidenceRecordWrapper.IsEmbedded() && c.evidenceRecordWrapper.Parent() != nil &&
		c.evidenceRecordWrapper.Parent().Filename() != "" {
		result = append(result, c.evidenceRecordWrapper.Parent().Filename())
	}
	return result
}
