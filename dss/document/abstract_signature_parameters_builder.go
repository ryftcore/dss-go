// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/AbstractSignatureParametersBuilder.java (DSS 6.5.RC1).
package document

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// AbstractSignatureParametersBuilderTarget is the subset of AbstractSignatureParameters[TP]'s
// exported API AbstractSignatureParametersBuilder needs to populate a built instance, factored
// out as an interface (independent of the TP type parameter) so the builder itself can be
// generic over SP alone, matching the Java @SuppressWarnings("rawtypes") raw-typed
// AbstractSignatureParametersBuilder<SP extends AbstractSignatureParameters>. Every concrete
// signature parameters type across every format (CAdES, XAdES, ... ported in later phases)
// satisfies it by embedding AbstractSignatureParameters[TP].
type AbstractSignatureParametersBuilderTarget interface {
	model.SerializableSignatureParameters

	SetSigningCertificate(signingCertificate *model.CertificateToken)
	SetCertificateChain(certificateChain []*model.CertificateToken)
	SetBLevelParams(bLevelParams *model.BLevelParameters)
	SetEncryptionAlgorithm(encryptionAlgorithm enumerations.EncryptionAlgorithm)
	SetDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm)
}

// AbstractSignatureParametersBuilder builds a signature parameters instance, generic over the SP
// implementation of AbstractSignatureParametersBuilderTarget to be created.
type AbstractSignatureParametersBuilder[SP AbstractSignatureParametersBuilderTarget] struct {
	// signingCertificate is a signing certificate to be used for a signature creation.
	signingCertificate *model.CertificateToken

	// certificateChain is a certificate chain of the signing certificate.
	certificateChain []*model.CertificateToken

	// encryptionAlgorithm is the encryption algorithm used for a signature creation by the
	// current signing-certificate.
	encryptionAlgorithm enumerations.EncryptionAlgorithm

	// digestAlgorithm is the digest algorithm used to hash signed data on signing.
	digestAlgorithm enumerations.DigestAlgorithm

	// bLevelParams holds the BLevelParameters.
	bLevelParams *model.BLevelParameters

	// InitParameters initializes and returns empty signature parameters. Go has no method
	// overriding across an embedded base, so Java's abstract initParameters() becomes this
	// exported function field - the same template-method-as-struct-field idiom used by
	// spi.DSSSecurityFactory (spi/dss_security_factory.go) and AbstractResourcesHandler
	// (document/abstract_resources_handler.go) - to be set by the concrete format builder
	// (CAdESSignatureParametersBuilder etc., ported in later phases) before #Build is called.
	InitParameters func() SP
}

// NewAbstractSignatureParametersBuilder is the default constructor. Port of
// AbstractSignatureParametersBuilder(CertificateToken).
func NewAbstractSignatureParametersBuilder[SP AbstractSignatureParametersBuilderTarget](signingCertificate *model.CertificateToken) *AbstractSignatureParametersBuilder[SP] {
	return NewAbstractSignatureParametersBuilderWithChain[SP](signingCertificate, nil)
}

// NewAbstractSignatureParametersBuilderWithChain is a constructor with a certificateChain. Port
// of AbstractSignatureParametersBuilder(CertificateToken, List).
func NewAbstractSignatureParametersBuilderWithChain[SP AbstractSignatureParametersBuilderTarget](signingCertificate *model.CertificateToken, certificateChain []*model.CertificateToken) *AbstractSignatureParametersBuilder[SP] {
	return &AbstractSignatureParametersBuilder[SP]{
		signingCertificate: signingCertificate,
		certificateChain:   certificateChain,
		bLevelParams:       model.NewBLevelParameters(),
	}
}

// SetEncryptionAlgorithm sets an encryption algorithm used by the signing-certificate's key
// pair. Port of #setEncryptionAlgorithm.
func (b *AbstractSignatureParametersBuilder[SP]) SetEncryptionAlgorithm(encryptionAlgorithm enumerations.EncryptionAlgorithm) *AbstractSignatureParametersBuilder[SP] {
	b.encryptionAlgorithm = encryptionAlgorithm
	return b
}

// SetDigestAlgorithm sets a digest algorithm to be used to hash the signed data. Port of
// #setDigestAlgorithm.
func (b *AbstractSignatureParametersBuilder[SP]) SetDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *AbstractSignatureParametersBuilder[SP] {
	b.digestAlgorithm = digestAlgorithm
	return b
}

// BLevel returns BLevelParameters. Port of #bLevel.
func (b *AbstractSignatureParametersBuilder[SP]) BLevel() *model.BLevelParameters {
	return b.bLevelParams
}

// SetBLevelParams sets a BLevelParameters (e.g. a SigningDate). Port of #setBLevelParams.
func (b *AbstractSignatureParametersBuilder[SP]) SetBLevelParams(bLevelParams *model.BLevelParameters) *AbstractSignatureParametersBuilder[SP] {
	b.bLevelParams = bLevelParams
	return b
}

// Build ports #build.
func (b *AbstractSignatureParametersBuilder[SP]) Build() SP {
	signatureParameters := b.InitParameters()
	signatureParameters.SetSigningCertificate(b.signingCertificate)
	signatureParameters.SetCertificateChain(b.certificateChain)
	signatureParameters.SetBLevelParams(b.bLevelParams)
	if b.encryptionAlgorithm != "" {
		signatureParameters.SetEncryptionAlgorithm(b.encryptionAlgorithm)
	}
	if b.digestAlgorithm != "" {
		signatureParameters.SetDigestAlgorithm(b.digestAlgorithm)
	}
	return signatureParameters
}

// compile-time interface assertion.
var _ model.SignatureParametersBuilder[AbstractSignatureParametersBuilderTarget] = (*AbstractSignatureParametersBuilder[AbstractSignatureParametersBuilderTarget])(nil)
