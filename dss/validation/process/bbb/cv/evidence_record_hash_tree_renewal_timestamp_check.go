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
	diagnosticData *diagnostic.Data

	// timestampWrapper is the time-stamp token to check.
	timestampWrapper *diagnostic.TimestampWrapper
}

// NewEvidenceRecordHashTreeRenewalTimestampCheck is the default constructor.
// Port of
// EvidenceRecordHashTreeRenewalTimestampCheck(Provider, XmlCV, Data, TimestampWrapper, LevelRule).
func NewEvidenceRecordHashTreeRenewalTimestampCheck(i18nProvider *i18n.Provider,
	result *process.Result[*jaxb.XmlCV], diagnosticData *diagnostic.Data,
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
	evidenceRecord := c.relatedEvidenceRecord(c.timestampWrapper)
	return c.timestampCoversAllOriginalDocuments(evidenceRecord, c.timestampWrapper)
}

// relatedEvidenceRecord ports the private
// getRelatedEvidenceRecord(TimestampWrapper), whose IllegalStateException is a
// panic here: it reports a diagnostic-data inconsistency, not a validation
// outcome, and process() has no error channel.
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) relatedEvidenceRecord(
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

// coveredDocuments ports the private getCoveredDocuments(List).
func (c *EvidenceRecordHashTreeRenewalTimestampCheck) coveredDocuments(
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
	evidenceRecordCoveredDocuments := c.coveredDocuments(evidenceRecord.DigestMatchers())
	timestampCoveredDocuments := c.coveredDocuments(timestampWrapper.DigestMatchers())
	return multisetCovers(evidenceRecordCoveredDocuments, timestampCoveredDocuments)
}

// multisetCovers reports whether every element of required, counted with its
// multiplicity, is present in available. Java scans the time-stamp's covered
// documents with List#contains and removes each match (List#remove) "to avoid
// checking duplicates" - an O(n*m) multiset-inclusion test; counting the
// available names in a map gives the identical verdict in O(n+m).
func multisetCovers(required, available []string) bool {
	remaining := make(map[string]int, len(available))
	for _, name := range available {
		remaining[name]++
	}
	for _, name := range required {
		if remaining[name] == 0 {
			return false
		}
		remaining[name]-- // consume the object to avoid checking duplicates
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
