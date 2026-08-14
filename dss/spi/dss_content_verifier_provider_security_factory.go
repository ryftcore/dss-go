// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/security/DSSContentVerifierProviderSecurityFactory.java (DSS 6.5.RC1).
//
// Upstream builds a org.bouncycastle.operator.ContentVerifierProvider bound to a PublicKey:
// callers later ask it for a ContentVerifier per AlgorithmIdentifier, write the signed content
// to that verifier's OutputStream, and compare the digest against a supplied signature - the
// generic BC machinery behind CMS SignerInformation#verify and TimeStampToken validation.
//
// Forward-declared for phase-2b: the CMS/timestamp phase (which owns SignerInformation and
// TimeStampToken) is the actual consumer of this factory. Rather than leave an empty
// placeholder, ContentVerifier.Verify below is implemented now on crypto/x509, mirroring the
// verification technique ocspTokenVerifySignature (dss_asn1_utils.go's OCSP neighbour) already
// uses: build a throwaway *x509.Certificate carrying only the PublicKey and call CheckSignature
// with the crypto/x509.SignatureAlgorithm the SignatureAlgorithm resolves to. This gives the
// CMS phase a working verifier now instead of a type it must still implement from scratch.
//
// DEVIATION: no JCA provider registry to retry against, same as the sibling factories in this
// package (see dss_security_factory.go).
package spi

import (
	"crypto"
	stdx509 "crypto/x509"
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// dssContentVerifierProviderSecurityFactorySignatureAlgorithms maps the algorithms a
// ContentVerifierProvider can be asked to check onto their crypto/x509 counterparts. Algorithms
// crypto/x509 cannot verify (SHA-224, SHA-3, RIPEMD-160, PLAIN-ECDSA and Ed448 flavours) are
// deliberately absent: they yield a verification error rather than a silent success.
var dssContentVerifierProviderSecurityFactorySignatureAlgorithms = map[enumerations.SignatureAlgorithm]stdx509.SignatureAlgorithm{
	enumerations.SignatureAlgorithm_RSA_MD5:                 stdx509.MD5WithRSA,
	enumerations.SignatureAlgorithm_RSA_SHA1:                stdx509.SHA1WithRSA,
	enumerations.SignatureAlgorithm_RSA_SHA256:              stdx509.SHA256WithRSA,
	enumerations.SignatureAlgorithm_RSA_SHA384:              stdx509.SHA384WithRSA,
	enumerations.SignatureAlgorithm_RSA_SHA512:              stdx509.SHA512WithRSA,
	enumerations.SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1: stdx509.SHA256WithRSAPSS,
	enumerations.SignatureAlgorithm_RSA_SSA_PSS_SHA384_MGF1: stdx509.SHA384WithRSAPSS,
	enumerations.SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1: stdx509.SHA512WithRSAPSS,
	enumerations.SignatureAlgorithm_ECDSA_SHA1:              stdx509.ECDSAWithSHA1,
	enumerations.SignatureAlgorithm_ECDSA_SHA256:            stdx509.ECDSAWithSHA256,
	enumerations.SignatureAlgorithm_ECDSA_SHA384:            stdx509.ECDSAWithSHA384,
	enumerations.SignatureAlgorithm_ECDSA_SHA512:            stdx509.ECDSAWithSHA512,
	enumerations.SignatureAlgorithm_DSA_SHA1:                stdx509.DSAWithSHA1,
	enumerations.SignatureAlgorithm_DSA_SHA256:              stdx509.DSAWithSHA256,
	enumerations.SignatureAlgorithm_ED25519:                 stdx509.PureEd25519,
}

// ContentVerifier is the minimal replacement for BouncyCastle's ContentVerifier/
// ContentVerifierProvider pair: an object bound to one PublicKey, able to check whether
// signatureValue is a valid signatureAlgorithm signature over signedContent.
type ContentVerifier struct {
	publicKey crypto.PublicKey
}

// Verify checks signatureValue against signedContent using signatureAlgorithm and the bound
// PublicKey. Port of the ContentVerifier obtained from
// ContentVerifierProvider#get(AlgorithmIdentifier), applied to the signed content's stream.
func (v *ContentVerifier) Verify(signatureAlgorithm enumerations.SignatureAlgorithm, signedContent, signatureValue []byte) error {
	algorithm, supported := dssContentVerifierProviderSecurityFactorySignatureAlgorithms[signatureAlgorithm]
	if !supported {
		return fmt.Errorf("NoSuchAlgorithmException : %s cannot be verified", signatureAlgorithm)
	}
	// CheckSignature verifies "signed" against the receiver's public key, so the signer's key
	// is carried by a throwaway certificate value - there is no other stdlib entry point that
	// verifies an arbitrary signature against a bare PublicKey across every algorithm family.
	signer := &stdx509.Certificate{PublicKey: v.publicKey}
	err := signer.CheckSignature(algorithm, signedContent, signatureValue)
	if err == nil {
		return nil
	}
	// An RSA key whose public exponent is larger than crypto/rsa will work with at all fails
	// CheckSignature before the signature is ever looked at; see rsaLargeExponentVerify for why
	// doing the RSA verification primitive directly is the BouncyCastle-equivalent behaviour
	// and not a weakening of the check. Applied here for the same reason the sibling
	// SignerInformationVerifier.Verify applies it: both stand in for a BouncyCastle verifier
	// that has no such limit.
	if fallbackErr := rsaLargeExponentVerify(
		v.publicKey, signatureAlgorithm, signedContent, signatureValue); fallbackErr == nil {
		return nil
	}
	return err
}

// dssContentVerifierProviderSecurityFactoryClassName is ContentVerifierProvider.class.getSimpleName().
const dssContentVerifierProviderSecurityFactoryClassName = "ContentVerifierProvider"

// DSSContentVerifierProviderSecurityFactoryInstance builds a ContentVerifier for a given
// PublicKey. Port of DSSContentVerifierProviderSecurityFactory.INSTANCE.
var DSSContentVerifierProviderSecurityFactoryInstance = &DSSSecurityFactory[*model.PublicKey, *ContentVerifier]{
	FactoryClassName:  dssContentVerifierProviderSecurityFactoryClassName,
	ToString:          dssContentVerifierProviderSecurityFactoryToString,
	BuildWithProvider: dssContentVerifierProviderSecurityFactoryBuild,
}

func dssContentVerifierProviderSecurityFactoryToString(input *model.PublicKey) string {
	if input == nil {
		return "PublicKey with algorithm 'null'"
	}
	return fmt.Sprintf("PublicKey with algorithm '%s'", input.Algorithm())
}

func dssContentVerifierProviderSecurityFactoryBuild(input *model.PublicKey) (*ContentVerifier, error) {
	// A nil PublicKey is not a caller bug (see the identical fix and its rationale in this
	// package's sibling dss_signer_information_verifier_security_factory.go,
	// dssSignerInformationVerifierSecurityFactoryBuild): upstream's own
	// buildWithProvider(PublicKey, Provider) has no null guard either and simply lets whatever a
	// null key causes surface as a graceful, caught exception downstream. Panicking here would
	// turn that into an unrecoverable crash instead.
	if input == nil {
		return nil, model.NewDSSError("InvalidKeyException : the public key has not been parsed")
	}
	key := input.Key()
	if key == nil {
		return nil, model.NewDSSError("InvalidKeyException : the public key has not been parsed")
	}
	return &ContentVerifier{publicKey: key}, nil
}
