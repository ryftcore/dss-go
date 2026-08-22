// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/EvidenceRecordValidationCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EvidenceRecordValidationCheck verifies validity of the performed evidence
// record validation process.
type EvidenceRecordValidationCheck[T any] struct {
	*process.ChainItemBase[T]

	// evidenceRecord is the evidence record to check.
	evidenceRecord *diagnostic.EvidenceRecordWrapper

	// erValidationResult is the evidence record validation result.
	erValidationResult *jaxb.XmlValidationProcessEvidenceRecord
}

// NewEvidenceRecordValidationCheck is the default constructor. Port of
// EvidenceRecordValidationCheck(I18nProvider, T, EvidenceRecordWrapper, XmlValidationProcessEvidenceRecord, LevelRule).
func NewEvidenceRecordValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	evidenceRecord *diagnostic.EvidenceRecordWrapper, erValidationResult *jaxb.XmlValidationProcessEvidenceRecord,
	constraint policy.LevelRule) *EvidenceRecordValidationCheck[T] {
	c := &EvidenceRecordValidationCheck[T]{
		ChainItemBase:      process.NewChainItemBaseWithId(i18nProvider, result, constraint, evidenceRecord.Id()),
		evidenceRecord:     evidenceRecord,
		erValidationResult: erValidationResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *EvidenceRecordValidationCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeER
}

// Process performs the check. Port of process().
func (c *EvidenceRecordValidationCheck[T]) Process() bool {
	return c.IsValid(&c.erValidationResult.XmlConstraintsConclusionContent)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(): Java dereferences getProofOfExistence() unguarded, so
// a result without one raises a NullPointerException; the nil pointer panics
// here in its place.
func (c *EvidenceRecordValidationCheck[T]) BuildAdditionalInfo() *string {
	proofOfExistenceTime := c.erValidationResult.ProofOfExistence.Time.Time()
	date := process.GetFormattedDate(&proofOfExistenceTime)
	message := c.I18nProvider.GetMessage(i18n.MessageTag_EVIDENCE_RECORD_VALIDATION, c.evidenceRecord.Id(), date)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EvidenceRecordValidationCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_IRERVPC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EvidenceRecordValidationCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_IRERVPC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EvidenceRecordValidationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.erValidationResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(): the generated SubIndication
// member is a pointer, whose nil is Java's null.
func (c *EvidenceRecordValidationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.erValidationResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.erValidationResult.Conclusion.SubIndication.SubIndication()
}
