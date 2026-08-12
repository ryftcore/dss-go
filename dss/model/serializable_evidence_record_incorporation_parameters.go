// Ported from dss-model/.../SerializableEvidenceRecordIncorporationParameters.java (DSS 6.5.RC1).
package model

// SerializableEvidenceRecordIncorporationParameters contains common
// methods for evidence records incorporation within existing signatures.
type SerializableEvidenceRecordIncorporationParameters interface {
	// SignatureId gets an identifier of the signature to include the
	// evidence record into. Ports #getSignatureId.
	SignatureId() string

	// SetSignatureId sets an identifier of the signature to include the
	// evidence record into. When a document with a single signature is
	// provided, the value can be set to empty. Otherwise, the signature
	// with the given identifier shall be found in order to perform the
	// operation.
	SetSignatureId(signatureId string)

	// DetachedContents gets detached documents signed by a signature.
	// Ports #getDetachedContents.
	DetachedContents() []DSSDocument

	// SetDetachedContents sets detached documents signed by a signature.
	SetDetachedContents(detachedContents []DSSDocument)

	// IsParallelEvidenceRecord gets whether the evidence record should be
	// incorporated within an existing (latest) evidence-record unsigned
	// property, when available. Otherwise, a new evidence record
	// attribute is to be created for incorporation of the evidence
	// record. Ports #isParallelEvidenceRecord.
	IsParallelEvidenceRecord() bool

	// SetParallelEvidenceRecord sets whether the evidence record should
	// be incorporated within an existing (latest) evidence-record
	// unsigned property, when available. Otherwise, a new evidence
	// record attribute is to be created for incorporation of the
	// evidence record.
	//
	// Default: false (a new evidence record unsigned property is to be
	// created).
	SetParallelEvidenceRecord(parallelEvidenceRecord bool)
}
