// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/evidencerecord/AbstractSignatureEvidenceRecordDigestBuilder.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.validation.evidencerecord.AbstractSignatureEvidenceRecordDigestBuilder
// lands in this same Go package per S2B_BRIEF.md's package layout table.
//
// FORWARD DEPENDENCY: SignatureEvidenceRecordDigestBuilder (Java
// spi.validation.evidencerecord.SignatureEvidenceRecordDigestBuilder, the interface this class
// implements) is not in this manifest and is left unreferenced: Go has no "implements" clause to
// satisfy explicitly, and this abstract base's own methods do not need to call back into it.
// SignatureAttribute (Java spi.validation.SignatureAttribute) is a forward dependency handled
// opaquely, following evidence_record.go's precedent.
package validation

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// AbstractSignatureEvidenceRecordDigestBuilder is the abstract implementation of the
// SignatureEvidenceRecordDigestGenerator.
type AbstractSignatureEvidenceRecordDigestBuilder struct {
	// signatureDocument is the signature document to compute a hash value for.
	signatureDocument model.DSSDocument

	// digestAlgorithm is the digest algorithm to be used on hash computation. Default:
	// DigestAlgorithm.SHA256.
	digestAlgorithm enumerations.DigestAlgorithm

	// signature is the signature incorporating the evidence record.
	signature AdvancedSignature

	// evidenceRecordAttribute is the attribute containing an evidence record to compute digest
	// for.
	evidenceRecordAttribute SignatureAttribute

	// parallelEvidenceRecord defines whether the new evidence-record shall be added to the last
	// available evidence-record attribute, when present. Otherwise, the hash is computed based
	// on the whole document content (default behavior).
	parallelEvidenceRecord bool
}

// NewAbstractSignatureEvidenceRecordDigestBuilder instantiates the builder with a SHA-256
// digest algorithm. Port of the AbstractSignatureEvidenceRecordDigestBuilder(DSSDocument)
// constructor; Java's Objects.requireNonNull becomes a panic.
func NewAbstractSignatureEvidenceRecordDigestBuilder(signatureDocument model.DSSDocument) *AbstractSignatureEvidenceRecordDigestBuilder {
	return NewAbstractSignatureEvidenceRecordDigestBuilderWithAlgorithm(signatureDocument, enumerations.DigestAlgorithm_SHA256)
}

// NewAbstractSignatureEvidenceRecordDigestBuilderWithAlgorithm instantiates the builder with a
// custom digest algorithm. Port of the AbstractSignatureEvidenceRecordDigestBuilder(DSSDocument,
// DigestAlgorithm) constructor; Java's Objects.requireNonNull calls become panics.
func NewAbstractSignatureEvidenceRecordDigestBuilderWithAlgorithm(signatureDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) *AbstractSignatureEvidenceRecordDigestBuilder {
	if signatureDocument == nil {
		panic("Signature document cannot be null!")
	}
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	return &AbstractSignatureEvidenceRecordDigestBuilder{
		signatureDocument: signatureDocument,
		digestAlgorithm:   digestAlgorithm,
	}
}

// NewAbstractSignatureEvidenceRecordDigestBuilderFromSignature instantiates the builder from a
// signature for the given evidenceRecordAttribute. Port of the
// AbstractSignatureEvidenceRecordDigestBuilder(AdvancedSignature, SignatureAttribute,
// DigestAlgorithm) constructor; Java's Objects.requireNonNull calls become panics.
func NewAbstractSignatureEvidenceRecordDigestBuilderFromSignature(signature AdvancedSignature, evidenceRecordAttribute SignatureAttribute,
	digestAlgorithm enumerations.DigestAlgorithm) *AbstractSignatureEvidenceRecordDigestBuilder {
	if signature == nil {
		panic("Signature cannot be null!")
	}
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	return &AbstractSignatureEvidenceRecordDigestBuilder{
		signature:               signature,
		evidenceRecordAttribute: evidenceRecordAttribute,
		digestAlgorithm:         digestAlgorithm,
	}
}

// SignatureDocument returns the signature document to compute a hash value for, nil when the
// builder was instantiated from a signature instead. Exported accessor for the Java `protected`
// field of the same name, needed by concrete-format subclasses living in other Go packages
// (cades, ...) - flagged out-of-manifest addition to this otherwise-frozen file, see S3_BRIEF.md
// "flag needs in notes; the integrator arbitrates".
func (b *AbstractSignatureEvidenceRecordDigestBuilder) SignatureDocument() model.DSSDocument {
	return b.signatureDocument
}

// DigestAlgorithm returns the digest algorithm to be used on hash computation. Exported accessor
// for the Java `protected` field of the same name; see SignatureDocument's doc.
func (b *AbstractSignatureEvidenceRecordDigestBuilder) DigestAlgorithm() enumerations.DigestAlgorithm {
	return b.digestAlgorithm
}

// Signature returns the signature incorporating the evidence record, nil when the builder was
// instantiated from a signature document instead. Exported accessor for the Java `protected`
// field of the same name; see SignatureDocument's doc.
func (b *AbstractSignatureEvidenceRecordDigestBuilder) Signature() AdvancedSignature {
	return b.signature
}

// EvidenceRecordAttribute returns the attribute containing an evidence record to compute digest
// for. Exported accessor for the Java `protected` field of the same name; see SignatureDocument's
// doc.
func (b *AbstractSignatureEvidenceRecordDigestBuilder) EvidenceRecordAttribute() SignatureAttribute {
	return b.evidenceRecordAttribute
}

// IsParallelEvidenceRecord returns whether the new evidence-record shall be added to the last
// available evidence-record attribute, when present. Exported accessor for the Java `protected`
// field of the same name; see SignatureDocument's doc.
func (b *AbstractSignatureEvidenceRecordDigestBuilder) IsParallelEvidenceRecord() bool {
	return b.parallelEvidenceRecord
}

// SetParallelEvidenceRecord sets whether the message-imprint for an evidence record shall be
// computed as for a parallel evidence-record (i.e. to be incorporated within the latest
// evidence-record attribute, when available). Otherwise, computes the message-imprint based on
// the whole signature's content, including coverage of other existing evidence-record.
// Default: FALSE (computes digest based on the whole signature's content). Port of
// setParallelEvidenceRecord(boolean); Java returns `this` for chaining, reproduced here.
func (b *AbstractSignatureEvidenceRecordDigestBuilder) SetParallelEvidenceRecord(parallelEvidenceRecord bool) *AbstractSignatureEvidenceRecordDigestBuilder {
	b.parallelEvidenceRecord = parallelEvidenceRecord
	return b
}

// getDigest returns the digest of the given document. Port of the protected
// getDigest(DSSDocument).
//
// DEVIATION: Java builds `new Digest(digestAlgorithm, document.getDigestValue(digestAlgorithm))`
// from two separate document calls; model.DSSDocument.Digest(DigestAlgorithm) already returns
// the equivalent model.Digest struct directly (algorithm + value), so this simply forwards to
// it instead of reassembling the same struct from a raw byte value. Java's method is unchecked;
// the Go port surfaces document.Digest's error instead of hiding it, since every other data-
// dependent throw in this phase becomes a returned error per PORTING.md.
func (b *AbstractSignatureEvidenceRecordDigestBuilder) getDigest(document model.DSSDocument) (model.Digest, error) {
	return document.Digest(b.digestAlgorithm)
}
