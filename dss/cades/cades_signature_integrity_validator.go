// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESSignatureIntegrityValidator.java (DSS 6.5.RC1).
//
// DEVIATION: BouncyCastle's SignerInformation carries an implicit reference back to the
// CMSSignedData's encapsulated (or detached) content, which SignerInformation#verify uses
// directly when the SignerInfo carries no signedAttrs - the signature then covers the content
// itself rather than the DER SET OF the signed attributes. cmscore.SignerInfo, unlike BC's
// type, holds no such back-reference, so NewCAdESSignatureIntegrityValidator takes the bytes to
// verify against explicitly (signedContent): internal/cmscore.Attributes#DERSetEncoded() - "the
// message digest is computed on the DER encoding of the SignedAttrs value, with the tag of SET
// OF" (RFC 5652 clause 5.4) - when signed attributes are present (every CAdES baseline profile
// requires them), or the signed content itself otherwise. The caller (CAdESSignature, a sibling
// chunk not in this manifest) is expected to compute it accordingly.
package cades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
)

// CAdESSignatureIntegrityValidator validates integrity of a CAdES signature. Port of the class
// CAdESSignatureIntegrityValidator, extending spi.SignatureIntegrityValidator.
type CAdESSignatureIntegrityValidator struct {
	spi.SignatureIntegrityValidatorBase

	// signerInformation is the corresponding SignerInformation.
	signerInformation *cmscore.SignerInfo

	// signedContent is the bytes the signature is expected to cover; see the file header.
	signedContent []byte
}

// NewCAdESSignatureIntegrityValidator is the port of the constructor
// CAdESSignatureIntegrityValidator(SignerInformation); see the file header on signedContent.
func NewCAdESSignatureIntegrityValidator(signerInformation *cmscore.SignerInfo, signedContent []byte) *CAdESSignatureIntegrityValidator {
	v := &CAdESSignatureIntegrityValidator{
		signerInformation: signerInformation,
		signedContent:     signedContent,
	}
	v.InitSignatureIntegrityValidator(v)
	return v
}

// Verify is the port of the protected verify(PublicKey) override.
//
// The DSSException Java wraps a CMSSignerDigestMismatchException or any other verification
// failure in is returned here as an error, its message built the same way.
func (v *CAdESSignatureIntegrityValidator) Verify(publicKey *model.PublicKey) (bool, error) {
	signerInformationVerifier, err := spi.DSSSignerInformationVerifierSecurityFactoryPublicTokenInstance.Build(publicKey)
	if err != nil {
		return false, model.NewDSSErrorMessageCause("Unable to validate CMS Signature : "+err.Error(), err)
	}
	signatureAlgorithm := v.signerInformation.SignatureAlgorithm
	digestAlgorithm := v.signerInformation.DigestAlgorithm
	sigAlg, err := cadesSignatureIntegrityValidatorSignatureAlgorithm(signatureAlgorithm, digestAlgorithm)
	if err != nil {
		return false, model.NewDSSErrorMessageCause("Unable to validate CMS Signature : "+err.Error(), err)
	}
	if err := signerInformationVerifier.Verify(sigAlg, v.signedContent, v.signerInformation.Signature); err != nil {
		return false, model.NewDSSErrorMessageCause("Unable to validate CMS Signature : "+err.Error(), err)
	}
	return true, nil
}

// cadesSignatureIntegrityValidatorSignatureAlgorithm resolves the combined SignatureAlgorithm
// (encryption + digest) a CMS SignerInfo's signatureAlgorithm/digestAlgorithm fields designate,
// mirroring what BouncyCastle's CMSSignedHelper does internally before verify() runs. The
// signatureAlgorithm field commonly names only the encryption algorithm (e.g. plain
// rsaEncryption), so it is resolved on its own first and combined with the digest algorithm;
// only if that OID already designates a combined algorithm (e.g. an RSASSA-PSS OID with
// parameters) is it used directly.
func cadesSignatureIntegrityValidatorSignatureAlgorithm(signatureAlgorithm, digestAlgorithm *asn1ber.AlgorithmIdentifier) (enumerations.SignatureAlgorithm, error) {
	oid := signatureAlgorithm.Algorithm.String()
	if combined, err := enumerations.SignatureAlgorithmForOIDAndParams(oid, signatureAlgorithm.Parameters); err == nil {
		return combined, nil
	}
	encryptionAlgorithm, err := enumerations.EncryptionAlgorithmForOID(oid)
	if err != nil {
		return "", err
	}
	digestAlg, err := enumerations.DigestAlgorithmForOID(digestAlgorithm.Algorithm.String())
	if err != nil {
		return "", err
	}
	return enumerations.SignatureAlgorithmGetAlgorithm(encryptionAlgorithm, digestAlg), nil
}
