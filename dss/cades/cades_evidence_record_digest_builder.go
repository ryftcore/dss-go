// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/evidencerecord/CAdESEvidenceRecordDigestBuilder.java (DSS 6.5.RC1).
package cades

import (
	"fmt"
	"sort"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// EvidenceRecordDigestBuilder computes message-imprint of a CMS signature to be protected
// by an evidence-record.
type EvidenceRecordDigestBuilder struct {
	validation.AbstractSignatureEvidenceRecordDigestBuilder

	// detachedDocument is the original document in case of a detached signature.
	detachedDocument model.DSSDocument

	// derEncoded sets whether the signature shall be DER-encoded for a hash computation (as per
	// ETSI TS 119 122-3 v1.1.1).
	derEncoded bool
}

// NewEvidenceRecordDigestBuilder instantiates EvidenceRecordDigestBuilder with a
// SHA-256 digest algorithm. Port of CAdESEvidenceRecordDigestBuilder(DSSDocument).
func NewEvidenceRecordDigestBuilder(signatureDocument model.DSSDocument) *EvidenceRecordDigestBuilder {
	return &EvidenceRecordDigestBuilder{
		AbstractSignatureEvidenceRecordDigestBuilder: *validation.NewAbstractSignatureEvidenceRecordDigestBuilder(signatureDocument),
	}
}

// NewEvidenceRecordDigestBuilderWithAlgorithm instantiates EvidenceRecordDigestBuilder
// with a custom digest algorithm. Port of CAdESEvidenceRecordDigestBuilder(DSSDocument,
// DigestAlgorithm).
func NewEvidenceRecordDigestBuilderWithAlgorithm(signatureDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) *EvidenceRecordDigestBuilder {
	return &EvidenceRecordDigestBuilder{
		AbstractSignatureEvidenceRecordDigestBuilder: *validation.NewAbstractSignatureEvidenceRecordDigestBuilderWithAlgorithm(signatureDocument, digestAlgorithm),
	}
}

// newEvidenceRecordDigestBuilderFromSignature instantiates EvidenceRecordDigestBuilder
// from a signature for the given evidenceRecordAttribute. Used on validation of an existing
// evidence record. Port of the protected
// EvidenceRecordDigestBuilder(AdvancedSignature, SignatureAttribute, DigestAlgorithm)
// constructor.
func newEvidenceRecordDigestBuilderFromSignature(signature validation.AdvancedSignature,
	evidenceRecordAttribute validation.SignatureAttribute, digestAlgorithm enumerations.DigestAlgorithm) *EvidenceRecordDigestBuilder {
	return &EvidenceRecordDigestBuilder{
		AbstractSignatureEvidenceRecordDigestBuilder: *validation.NewAbstractSignatureEvidenceRecordDigestBuilderFromSignature(
			signature, evidenceRecordAttribute, digestAlgorithm),
	}
}

// SetDetachedContent sets an original document in case of a detached signature. When set,
// please use BuildExternalEvidenceRecordDigest to compute hash for both the signature and the
// original document. Port of #setDetachedContent, chainable.
func (b *EvidenceRecordDigestBuilder) SetDetachedContent(detachedDocument model.DSSDocument) *EvidenceRecordDigestBuilder {
	b.detachedDocument = detachedDocument
	return b
}

// SetDEREncoded sets whether a signature shall be DER-encoded prior to the hash computation.
// Port of #setDEREncoded, chainable.
func (b *EvidenceRecordDigestBuilder) SetDEREncoded(derEncoded bool) *EvidenceRecordDigestBuilder {
	b.derEncoded = derEncoded
	return b
}

// SetParallelEvidenceRecord overrides AbstractSignatureEvidenceRecordDigestBuilder, keeping the
// concrete return type. Port of the covariant-return #setParallelEvidenceRecord(boolean)
// override.
func (b *EvidenceRecordDigestBuilder) SetParallelEvidenceRecord(parallelEvidenceRecord bool) *EvidenceRecordDigestBuilder {
	b.AbstractSignatureEvidenceRecordDigestBuilder.SetParallelEvidenceRecord(parallelEvidenceRecord)
	return b
}

// Build generates the hash value for the signature enveloping the evidence-record. Port of
// #build().
func (b *EvidenceRecordDigestBuilder) Build() (model.Digest, error) {
	cmsObj, err := b.getCMS()
	if err != nil {
		return model.Digest{}, err
	}
	return b.getDigest(cmsObj)
}

// BuildExternalEvidenceRecordDigest builds a group of digests to be covered in case of an
// external-evidence-record incorporation. Note: the original detached document shall be
// provided within SetDetachedContent. Port of #buildExternalEvidenceRecordDigest().
//
// Returns the signature digest on the first position, and the digest of the detached document
// on the second.
func (b *EvidenceRecordDigestBuilder) BuildExternalEvidenceRecordDigest() ([]model.Digest, error) {
	cmsObj, err := b.getCMS()
	if err != nil {
		return nil, err
	}
	signatureDigest, err := b.getDigest(cmsObj)
	if err != nil {
		return nil, err
	}
	originalDocumentDigest, err := b.getDigestOfDocument(b.detachedDocument)
	if err != nil {
		return nil, err
	}
	return []model.Digest{signatureDigest, originalDocumentDigest}, nil
}

// getCMS ports the protected #getCMS().
//
// Panics when neither a signature nor a signature document was provided (Java's
// IllegalStateException).
func (b *EvidenceRecordDigestBuilder) getCMS() (*cms.CMS, error) {
	if b.Signature() != nil {
		cadesSignature, ok := b.Signature().(*Signature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", b.Signature())
		}
		return cadesSignature.CMS(), nil
	} else if b.SignatureDocument() != nil {
		return cms.UtilsParseToCMS(b.SignatureDocument())
	}
	panic("Either a signature or a signature document shall be provided!")
}

// getDigest ports the protected #getDigest(CMS): the digest of cmsObj to be protected by an
// evidence record.
func (b *EvidenceRecordDigestBuilder) getDigest(cmsObj *cms.CMS) (model.Digest, error) {
	messageImprint, err := b.getCMSContentInfoMessageImprint(cmsObj)
	if err != nil {
		return model.Digest{}, err
	}
	digest, err := spi.DSSUtilsDigest(b.DigestAlgorithm(), messageImprint)
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(b.DigestAlgorithm(), digest), nil
}

// getDigestOfDocument ports the protected getDigest(DSSDocument) inherited from
// AbstractSignatureEvidenceRecordDigestBuilder (unexported there, so reproduced here rather than
// widening that method's visibility for a single non-CMS-typed caller).
func (b *EvidenceRecordDigestBuilder) getDigestOfDocument(document model.DSSDocument) (model.Digest, error) {
	return document.Digest(b.DigestAlgorithm())
}

// getCMSContentInfoMessageImprint returns a CMSSignedData's message-imprint to be protected by
// an evidence record. Port of the protected #getCMSContentInfoMessageImprint(CMS).
func (b *EvidenceRecordDigestBuilder) getCMSContentInfoMessageImprint(cmsObj *cms.CMS) ([]byte, error) {
	if b.IsParallelEvidenceRecord() || b.Signature() != nil {
		var err error
		cmsObj, err = b.getCMSSignedDataBeforeLastEvidenceRecord(cmsObj)
		if err != nil {
			return nil, err
		}
	}
	// ETSI TS 119 122-3 requires the CMS ContentInfo to be DER encoded, while RFC 4998 uses the
	// original CMS encoding. We need to provide a way to compute digest using the specified
	// encoding.
	return b.getEncoded(cmsObj), nil
}

// getCMSSignedDataBeforeLastEvidenceRecord creates a CMS that predates protection by the latest
// evidence-record. Port of the protected #getCMSSignedDataBeforeLastEvidenceRecord(CMS).
func (b *EvidenceRecordDigestBuilder) getCMSSignedDataBeforeLastEvidenceRecord(cmsObj *cms.CMS) (*cms.CMS, error) {
	signerWithERFound := false
	var newSignerInformationList []*cmscore.SignerInfo

	cadesSignature, signatureIsSet := b.Signature().(*Signature)

	for _, signerInformation := range cmsObj.SignerInfos() {
		current := signerInformation
		if !signatureIsSet || cadesSignature.SignerInformation() == signerInformation {
			unsignedAttributes := UnsignedAttributesBuild(signerInformation)

			var targetEvidenceRecordAttribute *Attribute
			if b.IsParallelEvidenceRecord() {
				targetEvidenceRecordAttribute = cadesEvidenceRecordDigestBuilderLatestEvidenceRecordAttribute(unsignedAttributes)
			} else if b.EvidenceRecordAttribute() != nil {
				attribute, ok := b.EvidenceRecordAttribute().(*Attribute)
				if !ok {
					return nil, fmt.Errorf("unexpected evidence record attribute type %T", b.EvidenceRecordAttribute())
				}
				targetEvidenceRecordAttribute = attribute
			}

			if targetEvidenceRecordAttribute != nil {
				if signerWithERFound {
					// At most one of the SignerInfo instances within the SignedData instance
					// shall contain evidence-records attributes. If the SignerInfo instance
					// contains more than one evidence-records attribute, only the ER(s) in the
					// latest added evidence-records attribute shall be updated.
					return nil, exception.NewIllegalInputException(
						"The CMSSignedData contains multiple evidence record attributes! Unable to compute hash.")
				}

				unsignedAttributesTable := cadesEvidenceRecordDigestBuilderRemoveAttributesAtAndAfter(unsignedAttributes, targetEvidenceRecordAttribute)
				var err error
				current, err = cms.UtilsReplaceUnsignedAttributes(signerInformation, unsignedAttributesTable)
				if err != nil {
					return nil, err
				}

				if b.IsParallelEvidenceRecord() {
					signerWithERFound = true
				}
			}
		}
		newSignerInformationList = append(newSignerInformationList, current)
	}
	return cms.UtilsReplaceSigners(cmsObj, newSignerInformationList)
}

// cadesEvidenceRecordDigestBuilderLatestEvidenceRecordAttribute ports the private
// getLatestEvidenceRecordAttribute(UnsignedAttributes).
func cadesEvidenceRecordDigestBuilderLatestEvidenceRecordAttribute(unsignedAttributes *UnsignedAttributes) *Attribute {
	if unsignedAttributes != nil {
		attributes := unsignedAttributes.Attributes() // returns sorted
		if len(attributes) > 0 {
			// Once an evidence-records attribute is included within a SignedData instance, the
			// only changes that might be applied to the SignedData instance are the renewal of
			// the ER within the evidence-records attribute, the adding of a new ER within a new
			// AttributeValue of the latest evidence record attribute or the adding of another
			// evidence-records attribute. No other changes shall be applied to the SignedData
			// instance.
			lastAttribute := attributes[len(attributes)-1]
			if lastAttribute.IsEvidenceRecord() {
				return lastAttribute
			}
		}
	}
	return nil
}

// cadesEvidenceRecordDigestBuilderRemoveAttributesAtAndAfter ports the private
// removeAttributesAtAndAfter(UnsignedAttributes, Attribute).
func cadesEvidenceRecordDigestBuilderRemoveAttributesAtAndAfter(unsignedAttributes *UnsignedAttributes, unsignedAttribute *Attribute) cmscore.Attributes {
	var attributesList []*Attribute
	for _, attribute := range unsignedAttributes.Attributes() {
		if unsignedAttribute.Equals(attribute) {
			break // break if the target attribute is reached
		}
		attributesList = append(attributesList, attribute)
	}

	if len(attributesList) == 0 {
		return nil
	}

	// ensure the original order
	sort.Slice(attributesList, func(i, j int) bool {
		return AttributeCompare(attributesList[i], attributesList[j]) < 0
	})

	attrs := make(cmscore.Attributes, 0, len(attributesList))
	for _, a := range attributesList {
		values := make([][]byte, 0, len(a.AttrValues()))
		for _, v := range a.AttrValues() {
			values = append(values, v.Encoded())
		}
		attrs = append(attrs, cmscore.NewAttribute(a.ASN1Oid(), values...))
	}
	return attrs
}

// getEncoded ports the protected #getEncoded(CMS): the CMS binaries, DER-encoded when
// derEncoded is set.
func (b *EvidenceRecordDigestBuilder) getEncoded(cmsObj *cms.CMS) []byte {
	if b.derEncoded {
		return cmsObj.DEREncoded()
	}
	return cmsObj.Encoded()
}
