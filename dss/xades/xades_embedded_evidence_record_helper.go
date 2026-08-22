// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/evidencerecord/XAdESEmbeddedEvidenceRecordHelper.java (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// EmbeddedEvidenceRecordHelper contains common methods for validation of a XAdES embedded
// evidence record.
type EmbeddedEvidenceRecordHelper struct {
	validation.AbstractEmbeddedEvidenceRecordHelper
}

// NewEmbeddedEvidenceRecordHelperForSignature instantiates a helper for an evidence record
// applied for the whole signature content (not yet embedded). Port of the constructor
// EmbeddedEvidenceRecordHelper(Signature).
func NewEmbeddedEvidenceRecordHelperForSignature(signature *Signature) *EmbeddedEvidenceRecordHelper {
	h := &EmbeddedEvidenceRecordHelper{}
	h.InitAbstractEmbeddedEvidenceRecordHelper(h, signature)
	return h
}

// NewXAdESEmbeddedEvidenceRecordHelper is the default constructor. Port of the constructor
// EmbeddedEvidenceRecordHelper(Signature, Attribute); evidenceRecordAttribute may
// be nil, matching this port's use as the (Signature) overload too (see
// NewEmbeddedEvidenceRecordHelperForSignature, which forwards to the same Init entry point
// with a nil attribute already).
func NewEmbeddedEvidenceRecordHelper(signature *Signature, evidenceRecordAttribute *Attribute) *EmbeddedEvidenceRecordHelper {
	h := &EmbeddedEvidenceRecordHelper{}
	if evidenceRecordAttribute == nil {
		h.InitAbstractEmbeddedEvidenceRecordHelper(h, signature)
	} else {
		h.InitAbstractEmbeddedEvidenceRecordHelperWithAttribute(h, signature, evidenceRecordAttribute)
	}
	return h
}

// DigestBuilder implements AbstractEmbeddedEvidenceRecordHelperOverrides, shadowing
// AbstractEmbeddedEvidenceRecordHelper's own (panicking) default. Port of the protected
// #getDigestBuilder(AdvancedSignature, SignatureAttribute, DigestAlgorithm) override.
func (h *EmbeddedEvidenceRecordHelper) DigestBuilder(signature validation.AdvancedSignature,
	evidenceRecordAttribute validation.SignatureAttribute, digestAlgorithm enumerations.DigestAlgorithm) validation.SignatureEvidenceRecordDigestBuilder {
	digestBuilder := newEvidenceRecordDigestBuilderFromSignature(signature, evidenceRecordAttribute, digestAlgorithm)
	digestBuilder.SetDetachedContent(h.DetachedContents())
	return digestBuilder
}

// SetDEREncoding implements AbstractEmbeddedEvidenceRecordHelperOverrides, shadowing
// AbstractEmbeddedEvidenceRecordHelper's own (panicking) default. Port of the protected
// #setDEREncoding(SignatureEvidenceRecordDigestBuilder, boolean) override.
//
// Panics: the #setEncoding method is not supported for a XAdES signature digest computation
// (Java's UnsupportedOperationException).
func (h *EmbeddedEvidenceRecordHelper) SetDEREncoding(digestBuilder validation.SignatureEvidenceRecordDigestBuilder, derEncoded bool) {
	panic("The #setEncoding method is not supported for a XAdES signature digest computation!")
}

// IsEncodingSelectionSupported reports whether a DER/BER encoding selection is supported by the
// current implementation. Port of #isEncodingSelectionSupported().
func (h *EmbeddedEvidenceRecordHelper) IsEncodingSelectionSupported() bool {
	return false
}

// IsAbsentHashtreeSupported reports whether an absent hashtree computation is supported by the
// current implementation. Port of #isAbsentHashtreeSupported().
func (h *EmbeddedEvidenceRecordHelper) IsAbsentHashtreeSupported() bool {
	return false
}

// compile-time assertion that the helper satisfies the abstract base's contract.
var _ validation.AbstractEmbeddedEvidenceRecordHelperOverrides = (*EmbeddedEvidenceRecordHelper)(nil)
