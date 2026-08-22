// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/evidencerecord/CAdESEmbeddedEvidenceRecordHelper.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// EmbeddedEvidenceRecordHelper contains common methods for validation of a CAdES embedded
// evidence record.
type EmbeddedEvidenceRecordHelper struct {
	validation.AbstractEmbeddedEvidenceRecordHelper
}

// NewEmbeddedEvidenceRecordHelperForSignature instantiates a helper for an evidence record
// applied for the whole signature content (not yet embedded). Port of the constructor
// EmbeddedEvidenceRecordHelper(Signature).
func NewEmbeddedEvidenceRecordHelperForSignature(sig *Signature) *EmbeddedEvidenceRecordHelper {
	h := &EmbeddedEvidenceRecordHelper{}
	h.InitAbstractEmbeddedEvidenceRecordHelper(h, sig)
	return h
}

// NewCAdESEmbeddedEvidenceRecordHelper is the default constructor. Port of the constructor
// EmbeddedEvidenceRecordHelper(Signature, Attribute); evidenceRecordAttribute may
// be nil, matching this port's use as the (Signature) overload too (see
// NewEmbeddedEvidenceRecordHelperForSignature, which forwards to the same Init entry point
// with a nil attribute already).
func NewEmbeddedEvidenceRecordHelper(sig *Signature, evidenceRecordAttribute *Attribute) *EmbeddedEvidenceRecordHelper {
	h := &EmbeddedEvidenceRecordHelper{}
	h.InitAbstractEmbeddedEvidenceRecordHelperWithAttribute(h, sig, evidenceRecordAttribute)
	return h
}

// SetDetachedContents overrides AbstractEmbeddedEvidenceRecordHelper#setDetachedContents.
//
// Panics when detachedContents does not contain exactly one document (Java's
// IllegalArgumentException("One and only one detached document is allowed for an embedded
// evidence record in CAdES!")).
func (h *EmbeddedEvidenceRecordHelper) SetDetachedContents(detachedContents []model.DSSDocument) {
	if len(detachedContents) != 1 {
		panic("One and only one detached document is allowed for an embedded evidence record in CAdES!")
	}
	h.AbstractEmbeddedEvidenceRecordHelper.SetDetachedContents(detachedContents)
}

// DigestBuilder implements AbstractEmbeddedEvidenceRecordHelperOverrides, shadowing
// AbstractEmbeddedEvidenceRecordHelper's own (panicking) default. Port of the protected
// #getDigestBuilder(AdvancedSignature, SignatureAttribute, DigestAlgorithm) override.
func (h *EmbeddedEvidenceRecordHelper) DigestBuilder(sig validation.AdvancedSignature,
	evidenceRecordAttribute validation.SignatureAttribute, digestAlgorithm enumerations.DigestAlgorithm) validation.SignatureEvidenceRecordDigestBuilder {
	digestBuilder := newEvidenceRecordDigestBuilderFromSignature(sig, evidenceRecordAttribute, digestAlgorithm)
	if isDetached, err := h.isDetached(sig); err == nil && isDetached {
		digestBuilder.SetDetachedContent(h.detachedDocument())
	}
	return digestBuilder
}

// isDetached ports the private isDetached(AdvancedSignature).
//
// Panics for a signature that is not a *CAdESSignature (Java's
// IllegalStateException("Only instance of Signature is supported by
// EmbeddedEvidenceRecordHelper")).
func (h *EmbeddedEvidenceRecordHelper) isDetached(sig validation.AdvancedSignature) (bool, error) {
	cadesSignature, ok := sig.(*Signature)
	if !ok {
		panic("Only instance of CAdESSignature is supported by CAdESEmbeddedEvidenceRecordHelper")
	}
	return cadesSignature.CMS().IsDetachedSignature(), nil
}

// detachedDocument gets the detached document covered by a detached CAdES. Port of the
// protected #getDetachedDocument().
func (h *EmbeddedEvidenceRecordHelper) detachedDocument() model.DSSDocument {
	detachedContents := h.DetachedContents()
	if len(detachedContents) == 1 {
		return detachedContents[0]
	}
	return nil
}

// SetDEREncoding implements AbstractEmbeddedEvidenceRecordHelperOverrides, shadowing
// AbstractEmbeddedEvidenceRecordHelper's own (panicking) default. Port of the protected
// #setDEREncoding(SignatureEvidenceRecordDigestBuilder, boolean) override.
//
// Panics when digestBuilder is not a *CAdESEvidenceRecordDigestBuilder (Java's
// IllegalArgumentException("The digestBuilder shall be an instance of
// EvidenceRecordDigestBuilder!")).
func (h *EmbeddedEvidenceRecordHelper) SetDEREncoding(digestBuilder validation.SignatureEvidenceRecordDigestBuilder, derEncoded bool) {
	cadesDigestBuilder, ok := digestBuilder.(*EvidenceRecordDigestBuilder)
	if !ok {
		panic("The digestBuilder shall be an instance of CAdESEvidenceRecordDigestBuilder!")
	}
	cadesDigestBuilder.SetDEREncoded(derEncoded)
}

// IsEncodingSelectionSupported reports whether a DER/BER encoding selection is supported by the
// current implementation. Port of #isEncodingSelectionSupported().
func (h *EmbeddedEvidenceRecordHelper) IsEncodingSelectionSupported() bool {
	return true
}

// IsAbsentHashtreeSupported reports whether an absent hashtree computation is supported by the
// current implementation. Port of #isAbsentHashtreeSupported().
func (h *EmbeddedEvidenceRecordHelper) IsAbsentHashtreeSupported() bool {
	return true
}
