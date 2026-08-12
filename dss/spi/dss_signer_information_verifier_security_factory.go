// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/security/DSSSignerInformationVerifierSecurityFactory.java (DSS 6.5.RC1).
//
// Upstream builds a org.bouncycastle.cms.SignerInformationVerifier bound to a CertificateToken
// or PublicKey, later passed to CMS SignerInformation#verify(SignerInformationVerifier).
//
// Forward-declared for phase-2b: the CMS phase (which owns SignerInformation) is the actual
// consumer. As with DSSContentVerifierProviderSecurityFactory (this package's sibling file),
// SignerInformationVerifier.Verify is implemented now rather than left a stub, using the same
// crypto/x509.Certificate#CheckSignature technique - see that file's header for the rationale
// and the shared signature-algorithm table's provenance (ocspTokenVerifySignature).
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
	"github.com/utain/esig/dss/utils"
)

// dssSignerInformationVerifierSecurityFactorySignatureAlgorithms maps the algorithms a
// SignerInformationVerifier can be asked to check onto their crypto/x509 counterparts. See
// dss_content_verifier_provider_security_factory.go for why the table is duplicated rather than
// shared (PORTING.md: "no shared helpers across files").
var dssSignerInformationVerifierSecurityFactorySignatureAlgorithms = map[enumerations.SignatureAlgorithm]stdx509.SignatureAlgorithm{
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

// SignerInformationVerifier is the minimal replacement for BouncyCastle's
// org.bouncycastle.cms.SignerInformationVerifier: an object bound to one PublicKey, able to
// check whether signatureValue is a valid signatureAlgorithm signature over signedContent
// (a CMS SignerInfo's signed attributes, in the CMS phase's intended use).
type SignerInformationVerifier struct {
	publicKey crypto.PublicKey
}

// Verify checks signatureValue against signedContent using signatureAlgorithm and the bound
// PublicKey.
func (v *SignerInformationVerifier) Verify(signatureAlgorithm enumerations.SignatureAlgorithm, signedContent, signatureValue []byte) error {
	algorithm, supported := dssSignerInformationVerifierSecurityFactorySignatureAlgorithms[signatureAlgorithm]
	if !supported {
		return fmt.Errorf("NoSuchAlgorithmException : %s cannot be verified", signatureAlgorithm)
	}
	signer := &stdx509.Certificate{PublicKey: v.publicKey}
	return signer.CheckSignature(algorithm, signedContent, signatureValue)
}

// dssSignerInformationVerifierSecurityFactoryClassName is SignerInformationVerifier.class.getSimpleName().
const dssSignerInformationVerifierSecurityFactoryClassName = "SignerInformationVerifier"

// DSSSignerInformationVerifierSecurityFactoryCertificateTokenInstance builds a
// SignerInformationVerifier for the given CertificateToken. Port of
// DSSSignerInformationVerifierSecurityFactory.CERTIFICATE_TOKEN_INSTANCE.
var DSSSignerInformationVerifierSecurityFactoryCertificateTokenInstance = &DSSSecurityFactory[*model.CertificateToken, *SignerInformationVerifier]{
	FactoryClassName: dssSignerInformationVerifierSecurityFactoryClassName,
	ToString:         dssSignerInformationVerifierSecurityFactoryCertificateTokenToString,
	BuildWithProvider: func(input *model.CertificateToken) (*SignerInformationVerifier, error) {
		if input == nil {
			panic("Input cannot be null")
		}
		return dssSignerInformationVerifierSecurityFactoryBuild(input.PublicKey())
	},
}

func dssSignerInformationVerifierSecurityFactoryCertificateTokenToString(input *model.CertificateToken) string {
	if input == nil {
		return ""
	}
	return utils.ToBase64(input.Encoded())
}

// DSSSignerInformationVerifierSecurityFactoryPublicTokenInstance builds a
// SignerInformationVerifier for the given PublicKey. Port of
// DSSSignerInformationVerifierSecurityFactory.PUBLIC_TOKEN_INSTANCE.
var DSSSignerInformationVerifierSecurityFactoryPublicTokenInstance = &DSSSecurityFactory[*model.PublicKey, *SignerInformationVerifier]{
	FactoryClassName: dssSignerInformationVerifierSecurityFactoryClassName,
	ToString: func(input *model.PublicKey) string {
		if input == nil {
			return "PublicKey with algorithm 'null'"
		}
		return fmt.Sprintf("PublicKey with algorithm '%s'", input.Algorithm())
	},
	BuildWithProvider: dssSignerInformationVerifierSecurityFactoryBuild,
}

func dssSignerInformationVerifierSecurityFactoryBuild(publicKey *model.PublicKey) (*SignerInformationVerifier, error) {
	if publicKey == nil {
		panic("Input cannot be null")
	}
	key := publicKey.Key()
	if key == nil {
		return nil, model.NewDSSError("InvalidKeyException : the public key has not been parsed")
	}
	return &SignerInformationVerifier{publicKey: key}, nil
}
