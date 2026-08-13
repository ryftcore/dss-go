// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/security/DSSSignatureSecurityFactory.java (DSS 6.5.RC1).
//
// DEVIATION: Java's Signature.getInstance(String) resolves an opaque JCE algorithm name
// (signatureValue.getAlgorithm().getJCEId(), the sole caller in this manifest -
// AbstractSignatureService#isValidSignatureValue) against the registered JCA security providers
// (see spi.DSSSecurityProviderInitSystemProviders). Go has no such provider registry; this port
// keeps the upstream enumerations.SignatureAlgorithm value itself instead of routing through its
// JCE name string, and resolves it with spi.DSSContentVerifierProviderSecurityFactoryInstance /
// spi.ContentVerifier - the crypto/x509-based verification technique
// DSSContentVerifierProviderSecurityFactory (spi/dss_content_verifier_provider_security_factory.go)
// already established for CMS/timestamp signature verification. Signature below therefore mirrors
// only the subset of java.security.Signature's incremental API (initVerify/update/verify) that
// upstream actually exercises.
package document

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
)

// Signature is the minimal replacement for java.security.Signature used by
// DSSSignatureSecurityFactory: an object bound to a signature algorithm, initialized for
// verification with a public key, fed data via Update, and checked via Verify.
type Signature struct {
	algorithm enumerations.SignatureAlgorithm
	publicKey *model.PublicKey
	data      []byte
}

// InitVerify initializes this signature object for verification, using the public key. Port of
// Signature#initVerify(PublicKey).
func (s *Signature) InitVerify(publicKey *model.PublicKey) {
	s.publicKey = publicKey
	s.data = nil
}

// Update updates the data to be verified. Port of Signature#update(byte[]).
func (s *Signature) Update(data []byte) {
	s.data = append(s.data, data...)
}

// Verify verifies the passed-in signature. Port of Signature#verify(byte[]); a
// GeneralSecurityException from upstream (invalid key, unsupported algorithm, ...) becomes a
// returned error rather than a boolean false, letting AbstractSignatureService#isValidSignatureValue
// decide - as it already does for the failure case - to log and answer false.
func (s *Signature) Verify(signatureBytes []byte) (bool, error) {
	verifier, err := spi.DSSContentVerifierProviderSecurityFactoryInstance.Build(s.publicKey)
	if err != nil {
		return false, err
	}
	if err := verifier.Verify(s.algorithm, s.data, signatureBytes); err != nil {
		return false, err
	}
	return true, nil
}

// dssSignatureSecurityFactoryClassName is Signature.class.getSimpleName().
const dssSignatureSecurityFactoryClassName = "Signature"

// DSSSignatureSecurityFactoryInstance builds a Signature verifier bound to a SignatureAlgorithm.
// Port of DSSSignatureSecurityFactory.INSTANCE (see the file DEVIATION for why this is keyed by
// SignatureAlgorithm rather than a JCE algorithm-name string).
var DSSSignatureSecurityFactoryInstance = &spi.DSSSecurityFactory[enumerations.SignatureAlgorithm, *Signature]{
	FactoryClassName: dssSignatureSecurityFactoryClassName,
	ToString:         func(input enumerations.SignatureAlgorithm) string { return string(input) },
	BuildWithProvider: func(input enumerations.SignatureAlgorithm) (*Signature, error) {
		return &Signature{algorithm: input}, nil
	},
}
