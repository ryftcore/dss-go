// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/evidencerecord/checks/EvidenceRecordSignedFilesCoveredCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
//
// Package placement deviation: Java's vpfswatsp.evidencerecord
// is a package of its own, distinct from vpfswatsp. Everything else in the
// vpfswatsp tree folds into one Go package, but TimestampsValidationBlock
// extends vpftsp.TimestampsValidationBlock while vpftsp imports vpfswatsp
// (POEExtraction) - an import cycle Go forbids. The five evidence-record classes
// therefore keep Java's own vpfswatsp/evidencerecord package boundary; nothing
// in vpfswatsp, vpftsp or vpftspwatsp imports them (only the validation
// executor's DetailedReportBuilder does), so the edge only ever points upward.
package evidencerecord

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignedFilesCoveredCheck verifies whether all files originally
// signed by a signature are covered by the evidence record.
type SignedFilesCoveredCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationProcessEvidenceRecord]

	// evidenceRecord is the evidence record to be validated.
	evidenceRecord *diagnostic.EvidenceRecordWrapper
}

// NewEvidenceRecordSignedFilesCoveredCheck is the default constructor. Port of
// SignedFilesCoveredCheck(Provider, XmlValidationProcessEvidenceRecord, EvidenceRecordWrapper, LevelRule).
func NewEvidenceRecordSignedFilesCoveredCheck(i18nProvider *i18n.Provider,
	result *process.Result[*jaxb.XmlValidationProcessEvidenceRecord],
	evidenceRecord *diagnostic.EvidenceRecordWrapper,
	constraint policy.LevelRule) *SignedFilesCoveredCheck {
	c := &SignedFilesCoveredCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		evidenceRecord: evidenceRecord,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process(): the two nested Stream#allMatch
// / Stream#anyMatch predicates become plain loops.
func (c *SignedFilesCoveredCheck) Process() bool {
	if enumerations.EvidenceRecordOriginSignature == c.evidenceRecord.Origin() {
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
func (c *SignedFilesCoveredCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBCVERHasSDoc
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignedFilesCoveredCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBCVERHasSDocANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignedFilesCoveredCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *SignedFilesCoveredCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}
