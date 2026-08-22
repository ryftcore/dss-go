// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/tsl/TrustedListV6SignatureParametersBuilder.java (DSS 6.5.RC1).
//
// Java's covariant-return setter overrides (each re-declaring the inherited setter to return
// TrustedListV6SignatureParametersBuilder instead of the parent type) become explicit wrapper
// methods here: Go has no covariant return through embedding, so each wrapper calls the
// embedded AbstractTrustedListSignatureParametersBuilder's setter and returns the receiver.
package tsl

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TrustedListV6SignatureParametersBuilder creates Signature parameters for a Trusted List V6
// creation.
//
// NOTE: the same instance of SignatureParameters shall be used on calls
// SignatureService.GetDataToSign(...) and SignatureService.SignDocument(...).
type TrustedListV6SignatureParametersBuilder struct {
	AbstractTrustedListSignatureParametersBuilder
}

// NewTrustedListV6SignatureParametersBuilder is the constructor to build Signature Parameters
// for a Trusted List V6 signing with respect to ETSI TS 119 612.
// NOTE: This class creates a new XAdES signature, according to ETSI EN 319 132-1. Port of
// TrustedListV6SignatureParametersBuilder(CertificateToken, DSSDocument).
func NewTrustedListV6SignatureParametersBuilder(signingCertificate *model.CertificateToken, tlXmlDocument model.DSSDocument) *TrustedListV6SignatureParametersBuilder {
	b := &TrustedListV6SignatureParametersBuilder{}
	b.InitAbstractTrustedListSignatureParametersBuilder(b, signingCertificate, tlXmlDocument)
	return b
}

// SetReferenceId overrides AbstractTrustedListSignatureParametersBuilder, keeping the concrete
// return type. Port of the covariant-return #setReferenceId(String) override.
func (b *TrustedListV6SignatureParametersBuilder) SetReferenceId(referenceId string) *TrustedListV6SignatureParametersBuilder {
	b.AbstractTrustedListSignatureParametersBuilder.SetReferenceId(referenceId)
	return b
}

// SetReferenceDigestAlgorithm overrides AbstractTrustedListSignatureParametersBuilder, keeping
// the concrete return type. Port of the covariant-return
// #setReferenceDigestAlgorithm(DigestAlgorithm) override.
func (b *TrustedListV6SignatureParametersBuilder) SetReferenceDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *TrustedListV6SignatureParametersBuilder {
	b.AbstractTrustedListSignatureParametersBuilder.SetReferenceDigestAlgorithm(digestAlgorithm)
	return b
}

// SetDigestAlgorithm overrides AbstractTrustedListSignatureParametersBuilder, keeping the
// concrete return type. Port of the covariant-return #setDigestAlgorithm(DigestAlgorithm)
// override.
func (b *TrustedListV6SignatureParametersBuilder) SetDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *TrustedListV6SignatureParametersBuilder {
	b.AbstractSignatureParametersBuilder.SetDigestAlgorithm(digestAlgorithm)
	return b
}

// SetEncryptionAlgorithm overrides AbstractTrustedListSignatureParametersBuilder, keeping the
// concrete return type. Port of the covariant-return #setEncryptionAlgorithm(EncryptionAlgorithm)
// override.
func (b *TrustedListV6SignatureParametersBuilder) SetEncryptionAlgorithm(encryptionAlgorithm enumerations.EncryptionAlgorithm) *TrustedListV6SignatureParametersBuilder {
	b.AbstractSignatureParametersBuilder.SetEncryptionAlgorithm(encryptionAlgorithm)
	return b
}

// SetBLevelParams overrides AbstractTrustedListSignatureParametersBuilder, keeping the concrete
// return type. Port of the covariant-return #setBLevelParams(BLevelParameters) override.
func (b *TrustedListV6SignatureParametersBuilder) SetBLevelParams(bLevelParams *model.BLevelParameters) *TrustedListV6SignatureParametersBuilder {
	b.AbstractSignatureParametersBuilder.SetBLevelParams(bLevelParams)
	return b
}

// IsEn319132 implements AbstractTrustedListSignatureParametersBuilderOverrides. Port of the
// overridden protected isEn319132().
func (b *TrustedListV6SignatureParametersBuilder) IsEn319132() bool {
	return true
}

// TargetTLVersion implements AbstractTrustedListSignatureParametersBuilderOverrides. Port of
// the overridden protected getTargetTLVersion().
func (b *TrustedListV6SignatureParametersBuilder) TargetTLVersion() int {
	return 6
}

// compile-time assertion that the builder satisfies the abstract base's contract.
var _ AbstractTrustedListSignatureParametersBuilderOverrides = (*TrustedListV6SignatureParametersBuilder)(nil)
