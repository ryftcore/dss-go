// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/evidencerecord/AbstractEmbeddedEvidenceRecordHelper.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.validation.evidencerecord.AbstractEmbeddedEvidenceRecordHelper lands
// in this same Go package per S2B_BRIEF.md's package layout table.
//
// Java's protected abstract getDigestBuilder(...)/setDEREncoding(...) methods become the
// Overrides interface + Init pattern established in phase 1b/2a and already used throughout
// this package (see DefaultAdvancedSignatureOverrides in default_advanced_signature.go for the
// precedent this file follows).
//
// slf4j LOG.warn calls in buildDigest are dropped per phase 2a handoff fact "slf4j dropped
// unless load-bearing" - the swallow-error-and-return-empty-Digest behavior itself is preserved
// exactly, only the logging side effect is dropped.
package validation

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// AbstractEmbeddedEvidenceRecordHelperOverrides captures the two protected abstract methods of
// Java's AbstractEmbeddedEvidenceRecordHelper that concrete subclasses must supply.
type AbstractEmbeddedEvidenceRecordHelperOverrides interface {
	// DigestBuilder gets the implementation of the signature digest builder for the given
	// evidence record. Port of the protected abstract getDigestBuilder(AdvancedSignature,
	// SignatureAttribute, DigestAlgorithm).
	DigestBuilder(signature AdvancedSignature, evidenceRecordAttribute SignatureAttribute, digestAlgorithm enumerations.DigestAlgorithm) SignatureEvidenceRecordDigestBuilder

	// SetDEREncoding sets the encoding to be used on the hash computation to the
	// SignatureEvidenceRecordDigestBuilder, whether applicable. Port of the protected abstract
	// setDEREncoding(SignatureEvidenceRecordDigestBuilder, boolean).
	SetDEREncoding(digestBuilder SignatureEvidenceRecordDigestBuilder, derEncoded bool)
}

// AbstractEmbeddedEvidenceRecordHelper is the abstract implementation of
// InternalEvidenceRecordHelper containing common implementation methods.
type AbstractEmbeddedEvidenceRecordHelper struct {
	// overrides points back at the concrete helper; see InitAbstractEmbeddedEvidenceRecordHelper.
	overrides AbstractEmbeddedEvidenceRecordHelperOverrides

	// signature is the master signature.
	signature AdvancedSignature

	// evidenceRecordAttribute is the unsigned signature attribute embedding the evidence
	// record.
	evidenceRecordAttribute SignatureAttribute

	// orderOfAttribute is the position of the attribute within the signature.
	orderOfAttribute *int

	// orderWithinAttribute is the position of the current evidence record within the evidence
	// record attribute.
	orderWithinAttribute *int

	// detachedContents is the list of detached documents provided to the validation.
	detachedContents []model.DSSDocument
}

// InitAbstractEmbeddedEvidenceRecordHelper initializes an evidence record applied for the whole
// signature content (not yet embedded). Port of the
// AbstractEmbeddedEvidenceRecordHelper(AdvancedSignature) constructor.
func (h *AbstractEmbeddedEvidenceRecordHelper) InitAbstractEmbeddedEvidenceRecordHelper(overrides AbstractEmbeddedEvidenceRecordHelperOverrides, signature AdvancedSignature) {
	h.InitAbstractEmbeddedEvidenceRecordHelperWithAttribute(overrides, signature, nil)
}

// InitAbstractEmbeddedEvidenceRecordHelperWithAttribute is the port of the default
// AbstractEmbeddedEvidenceRecordHelper(AdvancedSignature, SignatureAttribute) constructor.
func (h *AbstractEmbeddedEvidenceRecordHelper) InitAbstractEmbeddedEvidenceRecordHelperWithAttribute(overrides AbstractEmbeddedEvidenceRecordHelperOverrides, signature AdvancedSignature, evidenceRecordAttribute SignatureAttribute) {
	h.overrides = overrides
	h.signature = signature
	h.evidenceRecordAttribute = evidenceRecordAttribute
}

// MasterSignature gets a master signature, enveloping the current evidence record. Port of
// getMasterSignature().
func (h *AbstractEmbeddedEvidenceRecordHelper) MasterSignature() AdvancedSignature {
	return h.signature
}

// EvidenceRecordAttribute gets the unsigned attribute property embedding the evidence record.
// Port of getEvidenceRecordAttribute().
func (h *AbstractEmbeddedEvidenceRecordHelper) EvidenceRecordAttribute() SignatureAttribute {
	return h.evidenceRecordAttribute
}

// OrderOfAttribute gets position of the evidence record carrying attribute within the
// signature. Port of getOrderOfAttribute().
func (h *AbstractEmbeddedEvidenceRecordHelper) OrderOfAttribute() *int {
	return h.orderOfAttribute
}

// SetOrderOfAttribute sets position of the evidence record carrying attribute within the
// signature. Port of setOrderOfAttribute(Integer).
func (h *AbstractEmbeddedEvidenceRecordHelper) SetOrderOfAttribute(orderOfAttribute *int) {
	h.orderOfAttribute = orderOfAttribute
}

// OrderWithinAttribute gets position of the evidence record within its carrying attribute.
// Port of getOrderWithinAttribute().
func (h *AbstractEmbeddedEvidenceRecordHelper) OrderWithinAttribute() *int {
	return h.orderWithinAttribute
}

// SetOrderWithinAttribute sets position of the evidence record within its carrying attribute.
// Port of setOrderWithinAttribute(Integer).
func (h *AbstractEmbeddedEvidenceRecordHelper) SetOrderWithinAttribute(orderWithinAttribute *int) {
	h.orderWithinAttribute = orderWithinAttribute
}

// DetachedContents gets a list of detached documents. Port of getDetachedContents().
func (h *AbstractEmbeddedEvidenceRecordHelper) DetachedContents() []model.DSSDocument {
	return h.detachedContents
}

// SetDetachedContents sets a list of documents used for validation of a detached signature.
// Port of setDetachedContents(List).
func (h *AbstractEmbeddedEvidenceRecordHelper) SetDetachedContents(detachedContents []model.DSSDocument) {
	h.detachedContents = detachedContents
}

// MasterSignatureDigest builds digest for the embedded evidence record for the given
// DigestAlgorithm. Port of getMasterSignatureDigest(DigestAlgorithm).
func (h *AbstractEmbeddedEvidenceRecordHelper) MasterSignatureDigest(digestAlgorithm enumerations.DigestAlgorithm) model.Digest {
	digestBuilder := h.overrides.DigestBuilder(h.signature, h.evidenceRecordAttribute, digestAlgorithm)
	return h.buildDigest(digestBuilder)
}

// MasterSignatureDigestWithEncoding builds digest for the embedded evidence record for the
// given DigestAlgorithm using a specified encoding. Port of the getMasterSignatureDigest(
// DigestAlgorithm, boolean) overload.
func (h *AbstractEmbeddedEvidenceRecordHelper) MasterSignatureDigestWithEncoding(digestAlgorithm enumerations.DigestAlgorithm, derEncoded bool) model.Digest {
	digestBuilder := h.overrides.DigestBuilder(h.signature, h.evidenceRecordAttribute, digestAlgorithm)
	h.overrides.SetDEREncoding(digestBuilder, derEncoded)
	return h.buildDigest(digestBuilder)
}

// buildDigest is the port of the private buildDigest(SignatureEvidenceRecordDigestBuilder):
// errors are swallowed and an empty Digest is returned, matching Java's broad catch (Exception
// e) branch.
func (h *AbstractEmbeddedEvidenceRecordHelper) buildDigest(digestBuilder SignatureEvidenceRecordDigestBuilder) model.Digest {
	digest, err := digestBuilder.Build()
	if err != nil {
		return model.Digest{} // return empty digest
	}
	return digest
}
