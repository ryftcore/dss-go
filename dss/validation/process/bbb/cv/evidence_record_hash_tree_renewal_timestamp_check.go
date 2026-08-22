// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/EvidenceRecordHashTreeRenewalTimestampCheck.java (DSS 6.5.RC1).
package cv

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EvidenceRecordHashTreeRenewalTimestampCheck verifies whether the HashTree
// renewal time-stamp is conclusive and covers all original archive data objects
// covered by the evidence record.
type EvidenceRecordHashTreeRenewalTimestampCheck struct {
	*process.ChainItemBase[*jaxb.XmlCV]

	// diagnosticData is the Diagnostic Data.
	diagnosticData *diagnostic.DiagnosticData

	// timestampWrapper is the time-stamp token to check.
	timestampWrapper *diagnostic.TimestampWrapper
}

// NewEvidenceRecordHashTreeRenewalTimestampCheck is the default constructor.
// Port of
// EvidenceRecordHashTreeRenewalTimestampCheck(I18nProvider, XmlCV, DiagnosticData, TimestampWrapper, LevelRule).
func NewEvidenceRecordHashTreeRenewalTimestampCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlCV], diagnosticData *diagnostic.DiagnosticData,
	timestampWrapper *diagnostic.TimestampWrapper,
	constraint policy.LevelRule) *EvidenceRecordHashTreeRenewalTimestampCheck {
	c := &EvidenceRecordHashTreeRenewalTimestampCheck{
		ChainItemBase:    process.NewChainItemBase(i18nProvider, result, constraint),
		diagnosticData:   diagnosticData,
		timestampWrapper: timestampWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) Process() bool {
	evidenceRecord := c.getRelatedEvidenceRecord(c.timestampWrapper)
	return c.timestampCoversAllOriginalDocuments(evidenceRecord, c.timestampWrapper)
}

// getRelatedEvidenceRecord ports the private
// getRelatedEvidenceRecord(TimestampWrapper), whose IllegalStateException is a
// panic here: it reports a diagnostic-data inconsistency, not a validation
// outcome, and process() has no error channel.
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) getRelatedEvidenceRecord(
	timestampWrapper *diagnostic.TimestampWrapper) *diagnostic.EvidenceRecordWrapper {
	for _, evidenceRecordWrapper := range c.diagnosticData.EvidenceRecords() {
		// List#contains uses AbstractTokenProxy#equals, i.e. same wrapper type
		// and same Id - not object identity, which would never match since the
		// wrappers are rebuilt on every getTimestampList() call.
		for _, timestamp := range evidenceRecordWrapper.TimestampList() {
			if timestamp.Equals(timestampWrapper) {
				return evidenceRecordWrapper
			}
		}
	}
	panic(fmt.Sprintf("Not found a corresponding evidence record for a time-stamp with Id '%s'",
		timestampWrapper.Id()))
}

// getCoveredDocuments ports the private getCoveredDocuments(List).
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) getCoveredDocuments(
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher) []string {
	var documentNames []string
	for _, d := range digestMatchers {
		if enumerations.DigestMatcherTypeEvidenceRecordArchiveObject == digestMatcherType(d) && d.DataFound {
			documentNames = append(documentNames, digestMatcherDocumentName(d))
		}
	}
	return documentNames
}

// timestampCoversAllOriginalDocuments ports the private
// timestampCoversAllOriginalDocuments(EvidenceRecordWrapper, TimestampWrapper).
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) timestampCoversAllOriginalDocuments(
	evidenceRecord *diagnostic.EvidenceRecordWrapper, timestampWrapper *diagnostic.TimestampWrapper) bool {
	evidenceRecordCoveredDocuments := c.getCoveredDocuments(evidenceRecord.DigestMatchers())
	timestampCoveredDocuments := c.getCoveredDocuments(timestampWrapper.DigestMatchers())
	for _, originalDataObject := range evidenceRecordCoveredDocuments {
		index := -1
		for i, covered := range timestampCoveredDocuments {
			if covered == originalDataObject {
				index = i
				break
			}
		}
		if index < 0 {
			return false
		}
		// remove object to avoid checking duplicates
		timestampCoveredDocuments = append(timestampCoveredDocuments[:index], timestampCoveredDocuments[index+1:]...)
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBCVERTSTRN
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) ErrorMessageTag() i18n.MessageTag {
	if c.containsOtherDigests() {
		return i18n.MessageTagBBBCVERTSTRNANS2
	} else {
		return i18n.MessageTagBBBCVERTSTRNANS1
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) FailedIndicationForConclusion() enumerations.Indication {
	if c.containsOtherDigests() {
		return enumerations.IndicationFailed
	} else {
		return enumerations.IndicationIndeterminate
	}
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.containsOtherDigests() {
		return enumerations.SubIndicationHashFailure
	} else {
		return enumerations.SubIndicationSignedDataNotFound
	}
}

// containsOtherDigests ports the private containsOtherDigests().
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) containsOtherDigests() bool {
	for _, d := range c.timestampWrapper.DigestMatchers() {
		if enumerations.DigestMatcherTypeEvidenceRecordOrphanReference == digestMatcherType(d) {
			return true
		}
	}
	return false
}
