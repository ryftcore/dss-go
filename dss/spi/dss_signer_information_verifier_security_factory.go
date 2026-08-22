// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/security/DSSSignerInformationVerifierSecurityFactory.java (DSS 6.5.RC1).
//
// Upstream builds a org.bouncycastle.cms.SignerInformationVerifier bound to a CertificateToken
// or PublicKey, later passed to CMS SignerInformation#verify(SignerInformationVerifier).
//
// Structural stand-in: the CMS package (which owns SignerInformation) is the actual
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
	"crypto/rsa"
	stdx509 "crypto/x509"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// dssSignerInformationVerifierSecurityFactorySignatureAlgorithms maps the algorithms a
// SignerInformationVerifier can be asked to check onto their crypto/x509 counterparts. See
// dss_content_verifier_provider_security_factory.go for why the table is duplicated rather than
// shared (PORTING.md: "no shared helpers across files").
var dssSignerInformationVerifierSecurityFactorySignatureAlgorithms = map[enumerations.SignatureAlgorithm]stdx509.SignatureAlgorithm{
	enumerations.SignatureAlgorithmRSAMD5:              stdx509.MD5WithRSA,
	enumerations.SignatureAlgorithmRSASHA1:             stdx509.SHA1WithRSA,
	enumerations.SignatureAlgorithmRSASHA256:           stdx509.SHA256WithRSA,
	enumerations.SignatureAlgorithmRSASHA384:           stdx509.SHA384WithRSA,
	enumerations.SignatureAlgorithmRSASHA512:           stdx509.SHA512WithRSA,
	enumerations.SignatureAlgorithmRSASSAPSSSHA256MGF1: stdx509.SHA256WithRSAPSS,
	enumerations.SignatureAlgorithmRSASSAPSSSHA384MGF1: stdx509.SHA384WithRSAPSS,
	enumerations.SignatureAlgorithmRSASSAPSSSHA512MGF1: stdx509.SHA512WithRSAPSS,
	enumerations.SignatureAlgorithmECDSASHA1:           stdx509.ECDSAWithSHA1,
	enumerations.SignatureAlgorithmECDSASHA256:         stdx509.ECDSAWithSHA256,
	enumerations.SignatureAlgorithmECDSASHA384:         stdx509.ECDSAWithSHA384,
	enumerations.SignatureAlgorithmECDSASHA512:         stdx509.ECDSAWithSHA512,
	enumerations.SignatureAlgorithmDSASHA1:             stdx509.DSAWithSHA1,
	enumerations.SignatureAlgorithmDSASHA256:           stdx509.DSAWithSHA256,
	enumerations.SignatureAlgorithmED25519:             stdx509.PureEd25519,
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
//
// When the canonical check fails, RSA PKCS#1 v1.5 signatures get a second chance against the
// legacy "DigestInfo without the AlgorithmIdentifier NULL parameter" encoding; see
// dssSignerInformationVerifierVerifyRSAWithoutDigestInfoNullParameter for why that mirrors
// upstream rather than weakening the check.
func (v *SignerInformationVerifier) Verify(signatureAlgorithm enumerations.SignatureAlgorithm, signedContent, signatureValue []byte) error {
	algorithm, supported := dssSignerInformationVerifierSecurityFactorySignatureAlgorithms[signatureAlgorithm]
	if !supported {
		return fmt.Errorf("NoSuchAlgorithmException : %s cannot be verified", signatureAlgorithm)
	}
	signer := &stdx509.Certificate{PublicKey: v.publicKey}
	err := signer.CheckSignature(algorithm, signedContent, signatureValue)
	if err == nil {
		return nil
	}
	if fallbackErr := dssSignerInformationVerifierVerifyRSAWithoutDigestInfoNullParameter(
		v.publicKey, signatureAlgorithm, signedContent, signatureValue); fallbackErr == nil {
		return nil
	}
	// Second fallback: an RSA key whose public exponent is larger than crypto/rsa will work
	// with at all. CheckSignature above failed without ever looking at the signature in that
	// case; see rsaLargeExponentVerify for why doing the RSA verification primitive directly is
	// the BouncyCastle-equivalent behaviour and not a weakening of the check.
	if fallbackErr := rsaLargeExponentVerify(
		v.publicKey, signatureAlgorithm, signedContent, signatureValue); fallbackErr == nil {
		return nil
	}
	// The canonical failure is the one worth reporting; the fallback is only ever a second
	// chance, never a different diagnosis.
	return err
}

// dssSignerInformationVerifierSecurityFactoryPKCS1v15RSAAlgorithms is the exact set of RSA
// PKCS#1 v1.5 algorithms the legacy-DigestInfo fallback below applies to. RSASSA-PSS is
// deliberately absent (its encoding carries no DigestInfo at all), and so is every non-RSA
// family.
var dssSignerInformationVerifierSecurityFactoryPKCS1v15RSAAlgorithms = map[enumerations.SignatureAlgorithm]bool{
	enumerations.SignatureAlgorithmRSAMD5:    true,
	enumerations.SignatureAlgorithmRSASHA1:   true,
	enumerations.SignatureAlgorithmRSASHA256: true,
	enumerations.SignatureAlgorithmRSASHA384: true,
	enumerations.SignatureAlgorithmRSASHA512: true,
}

// dssSignerInformationVerifierVerifyRSAWithoutDigestInfoNullParameter verifies an RSA PKCS#1
// v1.5 signature whose DigestInfo omits the digest AlgorithmIdentifier's NULL parameter:
//
//	DigestInfo ::= SEQUENCE { digestAlgorithm SEQUENCE { OID },        -- no NULL
//	                          digest          OCTET STRING }
//
// instead of the canonical SEQUENCE { OID, NULL }. Some real-world signers produce it (two of
// the CAdES fixtures in cades/testdata/upstream do), and it is exactly the encoding
// BouncyCastle's org.bouncycastle.crypto.signers.RSADigestSigner#verifySignature accepts as its
// documented "NULL left out" fallback - so upstream DSS, which verifies CMS SignerInfos through
// that class, treats those signatures as intact. Go's crypto/rsa.VerifyPKCS1v15 hardcodes the
// NULL-bearing prefix and rejects them, which without the fallback below would report a valid
// signature as broken. Reproducing BouncyCastle's fallback keeps validation results identical
// to upstream's.
//
// This does NOT relax the padding check: the full EM block is still rebuilt and compared in
// whole by crypto/rsa (called with crypto.Hash(0), which means "the caller supplies the complete
// DigestInfo"), so no trailing bytes and no malformed padding are ever accepted - the only thing
// that changes is which of the two well-known DigestInfo spellings is expected.
func dssSignerInformationVerifierVerifyRSAWithoutDigestInfoNullParameter(publicKey crypto.PublicKey,
	signatureAlgorithm enumerations.SignatureAlgorithm, signedContent, signatureValue []byte) error {
	if !dssSignerInformationVerifierSecurityFactoryPKCS1v15RSAAlgorithms[signatureAlgorithm] {
		return fmt.Errorf("%s is not an RSA PKCS#1 v1.5 algorithm", signatureAlgorithm)
	}
	rsaPublicKey, isRSA := publicKey.(*rsa.PublicKey)
	if !isRSA {
		return fmt.Errorf("the public key is %T, not an RSA one", publicKey)
	}
	digestAlgorithm := signatureAlgorithm.DigestAlgorithm()
	digest, err := DSSUtilsDigest(digestAlgorithm, signedContent)
	if err != nil {
		return err
	}
	objectIdentifier, err := asn1ber.OIDFromString(digestAlgorithm.OID())
	if err != nil {
		return err
	}
	digestInfo := asn1ber.WriteSequence(append(
		asn1ber.WriteSequence(asn1ber.EncodeOID(objectIdentifier)),
		asn1ber.WriteTLV(asn1ber.TagOctetString, digest)...))
	// crypto.Hash(0) makes crypto/rsa compare the signature against the DigestInfo bytes given
	// here verbatim, rather than prepending its own NULL-bearing prefix.
	return rsa.VerifyPKCS1v15(rsaPublicKey, 0, digestInfo, signatureValue)
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
	// A nil publicKey is a legitimately reachable value here, not a caller bug: it flows straight
	// from CertificateValidity.PublicKey() (spi/certificate_validity.go), which returns nil
	// whenever a signing-certificate candidate is known only by its SignerIdentifier (no embedded
	// certificate resolved for it - e.g. a PDF revision whose CMS SignerInfo names a certificate
	// this port never located). Upstream's PUBLIC_TOKEN_INSTANCE.buildWithProvider(null, ...)
	// passes a null PublicKey straight into BouncyCastle's JcaSimpleSignerInfoVerifierBuilder
	// without a null check either, relying entirely on CAdESSignatureIntegrityValidator.verify's
	// own "catch (Exception e)" to turn whatever BC throws into a graceful DSSException - so a
	// panic here (as this used to do) turns a candidate correctly rejected as "not the signing
	// certificate" into a process-ending crash the moment a PDF/CMS carries more than one signer
	// candidate and only one of them resolves to a real certificate (confirmed by
	// pades/testdata/upstream/validation/PAdES-LT.pdf, a 3-signature PDF, in the PAdES
	// cross-validation harness). Matching Java's behavior means returning the same graceful
	// error every other failure path here returns, not panicking.
	if publicKey == nil {
		return nil, model.NewDSSError("InvalidKeyException : the public key has not been parsed")
	}
	key := publicKey.Key()
	if key == nil {
		return nil, model.NewDSSError("InvalidKeyException : the public key has not been parsed")
	}
	return &SignerInformationVerifier{publicKey: key}, nil
}
