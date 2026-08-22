// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/tsl/TrustedListV5SignatureParametersBuilder.java (DSS 6.5.RC1).
//
// Java's covariant-return setter overrides (each re-declaring the inherited setter to return
// TrustedListV5SignatureParametersBuilder instead of the parent type) become explicit wrapper
// methods here: Go has no covariant return through embedding, so each wrapper calls the
// embedded AbstractTrustedListSignatureParametersBuilder's setter and returns the receiver.
package tsl

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TrustedListV5SignatureParametersBuilder creates Signature parameters for a Trusted List V5
// creation.
//
// NOTE: the same instance of SignatureParameters shall be used on calls
// SignatureService.GetDataToSign(...) and SignatureService.SignDocument(...).
type TrustedListV5SignatureParametersBuilder struct {
	AbstractTrustedListSignatureParametersBuilder
}

// NewTrustedListV5SignatureParametersBuilder is the constructor to build Signature Parameters
// for a Trusted List V5 signing with respect to ETSI TS 119 612. Port of
// TrustedListV5SignatureParametersBuilder(CertificateToken, DSSDocument).
func NewTrustedListV5SignatureParametersBuilder(signingCertificate *model.CertificateToken, tlXmlDocument model.DSSDocument) *TrustedListV5SignatureParametersBuilder {
	b := &TrustedListV5SignatureParametersBuilder{}
	b.InitAbstractTrustedListSignatureParametersBuilder(b, signingCertificate, tlXmlDocument)
	return b
}

// SetReferenceId overrides AbstractTrustedListSignatureParametersBuilder, keeping the concrete
// return type. Port of the covariant-return #setReferenceId(String) override.
func (b *TrustedListV5SignatureParametersBuilder) SetReferenceId(referenceId string) *TrustedListV5SignatureParametersBuilder {
	b.AbstractTrustedListSignatureParametersBuilder.SetReferenceId(referenceId)
	return b
}

// SetReferenceDigestAlgorithm overrides AbstractTrustedListSignatureParametersBuilder, keeping
// the concrete return type. Port of the covariant-return
// #setReferenceDigestAlgorithm(DigestAlgorithm) override.
func (b *TrustedListV5SignatureParametersBuilder) SetReferenceDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *TrustedListV5SignatureParametersBuilder {
	b.AbstractTrustedListSignatureParametersBuilder.SetReferenceDigestAlgorithm(digestAlgorithm)
	return b
}

// SetDigestAlgorithm overrides AbstractTrustedListSignatureParametersBuilder, keeping the
// concrete return type. Port of the covariant-return #setDigestAlgorithm(DigestAlgorithm)
// override.
func (b *TrustedListV5SignatureParametersBuilder) SetDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *TrustedListV5SignatureParametersBuilder {
	b.AbstractSignatureParametersBuilder.SetDigestAlgorithm(digestAlgorithm)
	return b
}

// SetEncryptionAlgorithm overrides AbstractTrustedListSignatureParametersBuilder, keeping the
// concrete return type. Port of the covariant-return #setEncryptionAlgorithm(EncryptionAlgorithm)
// override.
func (b *TrustedListV5SignatureParametersBuilder) SetEncryptionAlgorithm(encryptionAlgorithm enumerations.EncryptionAlgorithm) *TrustedListV5SignatureParametersBuilder {
	b.AbstractSignatureParametersBuilder.SetEncryptionAlgorithm(encryptionAlgorithm)
	return b
}

// SetBLevelParams overrides AbstractTrustedListSignatureParametersBuilder, keeping the concrete
// return type. Port of the covariant-return #setBLevelParams(BLevelParameters) override.
func (b *TrustedListV5SignatureParametersBuilder) SetBLevelParams(bLevelParams *model.BLevelParameters) *TrustedListV5SignatureParametersBuilder {
	b.AbstractSignatureParametersBuilder.SetBLevelParams(bLevelParams)
	return b
}

// IsEn319132 implements AbstractTrustedListSignatureParametersBuilderOverrides. Port of the
// overridden protected isEn319132().
func (b *TrustedListV5SignatureParametersBuilder) IsEn319132() bool {
	return false
}

// TargetTLVersion implements AbstractTrustedListSignatureParametersBuilderOverrides. Port of
// the overridden protected getTargetTLVersion().
func (b *TrustedListV5SignatureParametersBuilder) TargetTLVersion() int {
	return 5
}

// compile-time assertion that the builder satisfies the abstract base's contract.
var _ AbstractTrustedListSignatureParametersBuilderOverrides = (*TrustedListV5SignatureParametersBuilder)(nil)
