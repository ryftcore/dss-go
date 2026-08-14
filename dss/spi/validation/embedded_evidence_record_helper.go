// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/evidencerecord/EmbeddedEvidenceRecordHelper.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.validation.evidencerecord.EmbeddedEvidenceRecordHelper lands in this
// same Go package per S2B_BRIEF.md's package layout table.
//
// FORWARD DEPENDENCY: this file is already relied upon opaquely by evidence_record.go (a sibling
// chunk file landed earlier), which documents it as one of two forward dependencies of
// EvidenceRecord.
//
// Java's Integer getOrderOfAttribute()/getOrderWithinAttribute() are nullable box types; ported
// as *int, matching this codebase's nullable-numeric convention (see model/eaa/claim.Claim's
// BooleanValue()/DateValue() doc comment for the same pointer-for-nullability rationale).
//
// getMasterSignatureDigest has two Java overloads (DigestAlgorithm) and (DigestAlgorithm,
// boolean); Go has no overloading, so the second lands as MasterSignatureDigestWithEncoding,
// matching the naming already used for AbstractSignatureEvidenceRecordDigestBuilder's sibling
// constructor overloads in this package (e.g. NewAbstractSignatureEvidenceRecordDigestBuilderWithAlgorithm).
package validation

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// EmbeddedEvidenceRecordHelper contains utility methods required for processing and validation
// of an embedded evidence record.
type EmbeddedEvidenceRecordHelper interface {
	// MasterSignature gets a master signature, enveloping the current evidence record. Port of
	// getMasterSignature().
	MasterSignature() AdvancedSignature

	// EvidenceRecordAttribute gets the unsigned attribute property embedding the evidence
	// record. NOTE: can be nil in case of a not yet embedded evidence record. Port of
	// getEvidenceRecordAttribute().
	EvidenceRecordAttribute() SignatureAttribute

	// OrderOfAttribute gets position of the evidence record carrying attribute within the
	// signature. NOTE: can be nil in case of a not yet embedded evidence record. Port of
	// getOrderOfAttribute().
	OrderOfAttribute() *int

	// OrderWithinAttribute gets position of the evidence record within its carrying attribute.
	// NOTE: can be nil in case of a not yet embedded evidence record. Port of
	// getOrderWithinAttribute().
	OrderWithinAttribute() *int

	// DetachedContents gets a list of detached documents. Port of getDetachedContents().
	DetachedContents() []model.DSSDocument

	// MasterSignatureDigest builds digest for the embedded evidence record for the given
	// DigestAlgorithm. This method uses an existing coding of a signature for hash generation.
	// Port of getMasterSignatureDigest(DigestAlgorithm).
	MasterSignatureDigest(digestAlgorithm enumerations.DigestAlgorithm) model.Digest

	// MasterSignatureDigestWithEncoding builds digest for the embedded evidence record for the
	// given DigestAlgorithm using a specified encoding. The method can be called only for a
	// CAdES signature implementation.
	//
	// NOTE: please use IsEncodingSelectionSupported to check whether the encoding choice is
	// supported by the current implementation. Use MasterSignatureDigest otherwise. Port of the
	// getMasterSignatureDigest(DigestAlgorithm, boolean) overload.
	MasterSignatureDigestWithEncoding(digestAlgorithm enumerations.DigestAlgorithm, derEncoded bool) model.Digest

	// IsEncodingSelectionSupported gets whether the selection of a target encoding is supported
	// by the current implementation. This method is used to resolve the interoperability issues
	// between ETSI TS 119 122-3 and RFC 4998 embedded ERS, requiring hash computation in
	// different ways. Port of isEncodingSelectionSupported().
	IsEncodingSelectionSupported() bool

	// IsAbsentHashtreeSupported gets whether the embedded evidence records without the reduced
	// hashtree are supported by the current signature implementation. This method resolves the
	// difference on processing between CAdES (TS 119 122-3) and XAdES (TS 119 132-3) embedded
	// evidence records, with the CAdES allowing omitted reduced hashtree, while XAdES requiring
	// such. Port of isAbsentHashtreeSupported().
	IsAbsentHashtreeSupported() bool
}
