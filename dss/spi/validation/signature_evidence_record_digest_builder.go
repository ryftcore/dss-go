// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/evidencerecord/SignatureEvidenceRecordDigestBuilder.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.validation.evidencerecord.SignatureEvidenceRecordDigestBuilder lands
// in this same Go package per S2B_BRIEF.md's package layout table.
package validation

import "github.com/utain/esig/dss/model"

// SignatureEvidenceRecordDigestBuilder generates a digest for an evidence record to be embedded
// within a given signature.
type SignatureEvidenceRecordDigestBuilder interface {
	// Build generates the hash value for the signature enveloping the evidence-record.
	// Note: this method is not supported for ASiC containers.
	//
	// Java's build() is unchecked but may throw at runtime (see
	// AbstractEmbeddedEvidenceRecordHelper.buildDigest's broad catch); per PORTING.md's
	// data-dependent-throw convention, that becomes a returned error here instead. Port of
	// build().
	Build() (model.Digest, error)
}
