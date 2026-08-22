// Ported from dss-document/src/main/java/eu/europa/esig/dss/evidencerecord/AbstractEvidenceRecordIncorporationParameters.java (DSS 6.5.RC1).
package document

import "github.com/ryftcore/dss-go/dss/model"

// AbstractEvidenceRecordIncorporationParameters contains parameters used on Evidence Record
// incorporation within an existing signature.
type AbstractEvidenceRecordIncorporationParameters struct {
	// signatureId is the identifier of a signature to include the evidence record into.
	signatureId string

	// detachedContents are the detached documents signed by a signature.
	detachedContents []model.DSSDocument

	// parallelEvidenceRecord defines whether the new evidence-record shall be added to the last
	// available evidence-record attribute, when present. Otherwise, the hash will be computed
	// based on the whole document content (default behavior).
	parallelEvidenceRecord bool
}

// NewAbstractEvidenceRecordIncorporationParameters is the default constructor. Port of the
// protected no-arg constructor.
func NewAbstractEvidenceRecordIncorporationParameters() AbstractEvidenceRecordIncorporationParameters {
	return AbstractEvidenceRecordIncorporationParameters{}
}

// SignatureId ports #getSignatureId.
func (p *AbstractEvidenceRecordIncorporationParameters) SignatureId() string {
	return p.signatureId
}

// SetSignatureId ports #setSignatureId.
func (p *AbstractEvidenceRecordIncorporationParameters) SetSignatureId(signatureId string) {
	p.signatureId = signatureId
}

// DetachedContents ports #getDetachedContents.
func (p *AbstractEvidenceRecordIncorporationParameters) DetachedContents() []model.DSSDocument {
	return p.detachedContents
}

// SetDetachedContents ports #setDetachedContents.
func (p *AbstractEvidenceRecordIncorporationParameters) SetDetachedContents(detachedContents []model.DSSDocument) {
	p.detachedContents = detachedContents
}

// IsParallelEvidenceRecord ports #isParallelEvidenceRecord.
func (p *AbstractEvidenceRecordIncorporationParameters) IsParallelEvidenceRecord() bool {
	return p.parallelEvidenceRecord
}

// SetParallelEvidenceRecord ports #setParallelEvidenceRecord.
func (p *AbstractEvidenceRecordIncorporationParameters) SetParallelEvidenceRecord(parallelEvidenceRecord bool) {
	p.parallelEvidenceRecord = parallelEvidenceRecord
}

// compile-time interface assertion.
var _ model.SerializableEvidenceRecordIncorporationParameters = (*AbstractEvidenceRecordIncorporationParameters)(nil)
