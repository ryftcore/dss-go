// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESSignatureIntegrityValidator.java (DSS 6.5.RC1).
//
// DEVIATION: BouncyCastle's SignerInformation carries an implicit reference back to the
// CMSSignedData's encapsulated (or detached) content, which SignerInformation#verify uses
// directly when the SignerInfo carries no signedAttrs - the signature then covers the content
// itself rather than the DER SET OF the signed attributes. cmscore.SignerInfo, unlike BC's
// type, holds no such back-reference, so NewSignatureIntegrityValidator takes the bytes to
// verify against explicitly (signedContent): internal/cmscore.Attributes#DERSetEncoded() - "the
// message digest is computed on the DER encoding of the SignedAttrs value, with the tag of SET
// OF" (RFC 5652 clause 5.4) - when signed attributes are present (every CAdES baseline profile
// requires them), or the signed content itself otherwise. The caller (Signature, a sibling
// chunk not in this manifest) is expected to compute it accordingly.
package cades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// CAdESSignatureIntegrityValidator validates integrity of a CAdES signature. Port of the class
// SignatureIntegrityValidator, extending spi.SignatureIntegrityValidator.
type SignatureIntegrityValidator struct {
	spi.SignatureIntegrityValidatorBase

	// signerInformation is the corresponding SignerInformation.
	signerInformation *cmscore.SignerInfo

	// signedContent is the bytes the signature is expected to cover; see the file header.
	signedContent []byte

	// contentDigestMismatch mirrors BouncyCastle's SignerInformation#verify: when signed
	// attributes are present, BC recomputes the digest of the content associated with the
	// SignerInformation (via the digest calculator the caller wired in - here, the caller's own
	// message-digest reference validation) and compares it against the message-digest signed
	// attribute BEFORE checking the raw signature bytes against the public key, throwing
	// CMSSignerDigestMismatchException on a mismatch. That failure does not depend on which
	// candidate's public key is being tried, so the caller (Signature.CheckSignatureIntegrity)
	// computes it once, up front, and every Verify call below fails uniformly when set - exactly
	// as every candidate would fail the same BC digest check in Java. Confirmed by
	// pades/testdata/upstream/validation/pdf-byterange-overlap.pdf in the PAdES cross-validation
	// harness, whose second signature has a /ByteRange pointing at content that does not match
	// its own CMS: SignatureIntact must come out false there, not just ReferenceDataIntact.
	contentDigestMismatch bool
}

// NewCAdESSignatureIntegrityValidator is the port of the constructor
// SignatureIntegrityValidator(SignerInformation); see the file header on signedContent and
// the contentDigestMismatch field doc on the extra parameter.
func NewSignatureIntegrityValidator(signerInformation *cmscore.SignerInfo, signedContent []byte, contentDigestMismatch bool) *SignatureIntegrityValidator {
	v := &SignatureIntegrityValidator{
		signerInformation:     signerInformation,
		signedContent:         signedContent,
		contentDigestMismatch: contentDigestMismatch,
	}
	v.InitSignatureIntegrityValidator(v)
	return v
}

// Verify is the port of the protected verify(PublicKey) override.
//
// The DSSException Java wraps a CMSSignerDigestMismatchException or any other verification
// failure in is returned here as an error, its message built the same way.
func (v *SignatureIntegrityValidator) Verify(publicKey *model.PublicKey) (bool, error) {
	if v.contentDigestMismatch {
		return false, model.NewDSSError("Unable to validate CMS Signature : message-digest attribute does not match the digest of the provided content")
	}
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
