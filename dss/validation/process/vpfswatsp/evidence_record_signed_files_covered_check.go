// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/evidencerecord/checks/EvidenceRecordSignedFilesCoveredCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// EvidenceRecordSignedFilesCoveredCheck verifies whether all files originally
// signed by a signature are covered by the evidence record.
type EvidenceRecordSignedFilesCoveredCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationProcessEvidenceRecord]

	// evidenceRecord is the evidence record to be validated.
	evidenceRecord *diagnostic.EvidenceRecordWrapper
}

// NewEvidenceRecordSignedFilesCoveredCheck is the default constructor. Port of
// EvidenceRecordSignedFilesCoveredCheck(I18nProvider, XmlValidationProcessEvidenceRecord, EvidenceRecordWrapper, LevelRule).
func NewEvidenceRecordSignedFilesCoveredCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessEvidenceRecord],
	evidenceRecord *diagnostic.EvidenceRecordWrapper,
	constraint policy.LevelRule) *EvidenceRecordSignedFilesCoveredCheck {
	c := &EvidenceRecordSignedFilesCoveredCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		evidenceRecord: evidenceRecord,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process(): the two nested Stream#allMatch
// / Stream#anyMatch predicates become plain loops.
func (c *EvidenceRecordSignedFilesCoveredCheck) Process() bool {
	if enumerations.EvidenceRecordOrigin_SIGNATURE == c.evidenceRecord.Origin() {
		// embedded signature covers all original documents
		return true
	}

	coveredSignatures := c.evidenceRecord.CoveredSignatures()
	evidenceRecordDigestMatchers := c.evidenceRecord.DigestMatchers()
	if utils.IsCollectionNotEmpty(coveredSignatures) {
		for _, signature := range coveredSignatures {
			digestMatchers := signature.DigestMatchers()
			for _, s := range digestMatchers {
				if !signedFileCovered(s, evidenceRecordDigestMatchers) {
					return false
				}
			}
		}
	}
	return true
}

// signedFileCovered ports the inner predicate of process():
// "s.getDocumentName() == null || evidenceRecordDigestMatchers.stream()
// .anyMatch(e -> s.getDocumentName().equals(e.getDocumentName()))". The
// generated DocumentName member is a *string, whose nil is Java's null; the
// equals() call therefore only matches a non-nil name against a non-nil one.
func signedFileCovered(signatureDigestMatcher *diagnosticjaxb.XmlDigestMatcher,
	evidenceRecordDigestMatchers []*diagnosticjaxb.XmlDigestMatcher) bool {
	if signatureDigestMatcher.DocumentName == nil {
		return true
	}
	for _, e := range evidenceRecordDigestMatchers {
		if e.DocumentName != nil && *signatureDigestMatcher.DocumentName == *e.DocumentName {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EvidenceRecordSignedFilesCoveredCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_ER_HASSDOC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EvidenceRecordSignedFilesCoveredCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_ER_HASSDOC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EvidenceRecordSignedFilesCoveredCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EvidenceRecordSignedFilesCoveredCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
