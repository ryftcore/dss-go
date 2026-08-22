// The build side of SignedData, replacing org.bouncycastle.cms.CMSSignedDataGenerator and
// org.bouncycastle.cms.SignerInfoGenerator.
//
// Only SignedData is built: CAdES signing has to produce one, whereas producing a
// TimeStampToken is a TSA's job. Output is DER throughout, since a signature is defined over
// DER; a caller that needs the BER encoding of a parsed document asks asn1ber for it.
package cmscore

import (
	"errors"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// SignerInfoBuilder assembles a SignerInfo. The version is derived from the SignerIdentifier
// per RFC 5652 clause 5.3 rather than being supplied.
type SignerInfoBuilder struct {
	// SID identifies the signing certificate. Required.
	SID *SignerIdentifier
	// DigestAlgorithm is the digest algorithm identifier. Required.
	DigestAlgorithm *asn1ber.AlgorithmIdentifier
	// SignedAttributes holds the attributes covered by the signature. A non-nil, possibly
	// empty value emits the signedAttrs field; nil omits it.
	SignedAttributes Attributes
	// SignatureAlgorithm is the signature algorithm identifier. Required.
	SignatureAlgorithm *asn1ber.AlgorithmIdentifier
	// Signature holds the signature octets. Required.
	Signature []byte
	// UnsignedAttributes holds the attributes outside the signature. A non-nil value emits
	// the unsignedAttrs field; nil omits it.
	UnsignedAttributes Attributes
}

// SignedAttributesDER returns the value the signature has to be computed over: the signed
// attributes as a DER SET OF, per RFC 5652 clause 5.4. It is nil when the builder carries no
// signed attributes, in which case the signature covers the content itself.
func (b *SignerInfoBuilder) SignedAttributesDER() []byte {
	if b.SignedAttributes == nil {
		return nil
	}
	return b.SignedAttributes.DERSetEncoded()
}

// Build assembles the SignerInfo.
func (b *SignerInfoBuilder) Build() (*SignerInfo, error) {
	switch {
	case b.SID == nil:
		return nil, errors.New("cmscore: the SignerInfo builder needs a SignerIdentifier")
	case b.DigestAlgorithm == nil:
		return nil, errors.New("cmscore: the SignerInfo builder needs a digest algorithm")
	case b.SignatureAlgorithm == nil:
		return nil, errors.New("cmscore: the SignerInfo builder needs a signature algorithm")
	case len(b.Signature) == 0:
		return nil, errors.New("cmscore: the SignerInfo builder needs a signature")
	}
	return &SignerInfo{
		Version:            ComputeSignerInfoVersion(b.SID),
		SID:                b.SID,
		DigestAlgorithm:    b.DigestAlgorithm,
		SignedAttributes:   b.SignedAttributes,
		SignatureAlgorithm: b.SignatureAlgorithm,
		Signature:          b.Signature,
		UnsignedAttributes: b.UnsignedAttributes,
		hasSignedAttrs:     b.SignedAttributes != nil,
		hasUnsignedAttrs:   b.UnsignedAttributes != nil,
	}, nil
}

// SignedDataBuilder assembles a SignedData. The version is computed per RFC 5652 clause 5.1
// rather than being supplied.
type SignedDataBuilder struct {
	// DigestAlgorithms holds the members of the digestAlgorithms SET. When empty, the
	// digest algorithms of the SignerInfos are used.
	DigestAlgorithms []*asn1ber.AlgorithmIdentifier
	// EncapContentInfo is the signed content and its type. Required.
	EncapContentInfo *EncapsulatedContentInfo
	// Certificates holds the members of the certificates field; an empty slice omits it.
	Certificates []CertificateChoice
	// CRLs holds the members of the crls field; an empty slice omits it.
	CRLs []RevocationInfoChoice
	// SignerInfos holds the members of the signerInfos SET.
	SignerInfos []*SignerInfo
}

// Build assembles the SignedData.
func (b *SignedDataBuilder) Build() (*SignedData, error) {
	if b.EncapContentInfo == nil {
		return nil, errors.New("cmscore: the SignedData builder needs an EncapsulatedContentInfo")
	}
	signedData := &SignedData{
		DigestAlgorithms: b.digestAlgorithms(),
		EncapContentInfo: b.EncapContentInfo,
		SignerInfos:      b.SignerInfos,
	}
	if len(b.Certificates) != 0 {
		signedData.Certificates = &CertificateSet{Choices: b.Certificates}
	}
	if len(b.CRLs) != 0 {
		signedData.CRLs = &RevocationInfoChoices{Choices: b.CRLs}
	}
	signedData.Version = ComputeSignedDataVersion(b.EncapContentInfo.EContentType,
		signedData.Certificates, signedData.CRLs, signedData.SignerInfos)
	return signedData, nil
}

// BuildCMS assembles the SignedData and wraps it in its id-signedData ContentInfo.
func (b *SignedDataBuilder) BuildCMS() (*CMS, error) {
	signedData, err := b.Build()
	if err != nil {
		return nil, err
	}
	return NewCMS(signedData), nil
}

// digestAlgorithms returns the digestAlgorithms members, defaulting to the distinct digest
// algorithms of the SignerInfos as CMSSignedDataGenerator does.
func (b *SignedDataBuilder) digestAlgorithms() []*asn1ber.AlgorithmIdentifier {
	if len(b.DigestAlgorithms) != 0 {
		return b.DigestAlgorithms
	}
	var algorithms []*asn1ber.AlgorithmIdentifier
	for _, signerInfo := range b.SignerInfos {
		known := false
		for _, algorithm := range algorithms {
			if algorithm.Equals(signerInfo.DigestAlgorithm) {
				known = true
				break
			}
		}
		if !known {
			algorithms = append(algorithms, signerInfo.DigestAlgorithm)
		}
	}
	return algorithms
}
