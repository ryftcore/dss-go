// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CertificateRefIdentifier.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: this file references CertificateRef (a spi.x509 type flattened into
// this package, ported in a sibling chunk of phase 2a). Its assumed shape, inferred from the
// Java getters actually called below, is a struct with accessor methods CertDigest()
// model.Digest, CertificateIdentifier() *SignerIdentifier, ResponderId() *ResponderId,
// Kid() string, X509Url() string and PublicKey() *model.PublicKey.
package spi

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// certificateRefIdentifierSkiDigestAlgorithm is the digest algorithm used to compute SKI, by
// RFC 6960. Port of the private static SKI_DIGEST_ALGO constant.
const certificateRefIdentifierSkiDigestAlgorithm = enumerations.DigestAlgorithmSHA1

// CertificateRefIdentifier is an identifier for a certificate token reference.
type CertificateRefIdentifier struct {
	model.IdentifierBase
}

// NewCertificateRefIdentifier builds the identifier of the given reference.
// Port of the protected CertificateRefIdentifier(CertificateRef) constructor.
//
// Java's DSSException("One of [certDigest, publicKeyDigest, issuerInfo, kid, x509Uri,
// publicKey] must be defined for a CertificateRef!") is data-dependent on the caller-supplied
// certificateRef, so it is returned as an error rather than raised as a panic.
func NewCertificateRefIdentifier(certificateRef *CertificateRef) (*CertificateRefIdentifier, error) {
	digest, err := certificateRefIdentifierDigest(certificateRef)
	if err != nil {
		return nil, err
	}
	return &CertificateRefIdentifier{
		model.NewIdentifierBaseFromDigest("CertificateRefIdentifier", "C-", digest),
	}, nil
}

// certificateRefIdentifierDigest ports the private static getDigest(CertificateRef).
func certificateRefIdentifierDigest(certificateRef *CertificateRef) (model.Digest, error) {
	if certDigest := certificateRef.CertDigest(); certDigest.Value() != nil {
		return certDigest, nil
	}
	if signerIdentifier := certificateRef.CertificateIdentifier(); signerIdentifier != nil {
		if ski := signerIdentifier.Ski(); ski != nil {
			return model.NewDigest(certificateRefIdentifierSkiDigestAlgorithm, ski), nil
		}
		if issuerSerialEncoded := signerIdentifier.IssuerSerialEncoded(); issuerSerialEncoded != nil {
			value, err := DSSUtilsDigest(model.IdentifierDigestAlgorithm, issuerSerialEncoded)
			if err != nil {
				return model.Digest{}, err
			}
			return model.NewDigest(model.IdentifierDigestAlgorithm, value), nil
		}
		if issuerName := signerIdentifier.IssuerName(); issuerName != nil {
			value, err := DSSUtilsDigest(model.IdentifierDigestAlgorithm, issuerName.Encoded())
			if err != nil {
				return model.Digest{}, err
			}
			return model.NewDigest(model.IdentifierDigestAlgorithm, value), nil
		}
	}
	if responderId := certificateRef.ResponderId(); responderId != nil {
		if ski := responderId.Ski(); ski != nil {
			return model.NewDigest(certificateRefIdentifierSkiDigestAlgorithm, ski), nil
		}
		if principal := responderId.X500Principal(); principal != nil {
			value, err := DSSUtilsDigest(model.IdentifierDigestAlgorithm, principal.Encoded())
			if err != nil {
				return model.Digest{}, err
			}
			return model.NewDigest(model.IdentifierDigestAlgorithm, value), nil
		}
	}
	if kid := certificateRef.Kid(); kid != "" {
		value, err := DSSUtilsDigest(model.IdentifierDigestAlgorithm, []byte(kid))
		if err != nil {
			return model.Digest{}, err
		}
		return model.NewDigest(model.IdentifierDigestAlgorithm, value), nil
	}
	if x509Url := certificateRef.X509Url(); x509Url != "" {
		value, err := DSSUtilsDigest(model.IdentifierDigestAlgorithm, []byte(x509Url))
		if err != nil {
			return model.Digest{}, err
		}
		return model.NewDigest(model.IdentifierDigestAlgorithm, value), nil
	}
	if publicKey := certificateRef.PublicKey(); publicKey != nil {
		value, err := DSSUtilsDigest(model.IdentifierDigestAlgorithm, publicKey.Encoded())
		if err != nil {
			return model.Digest{}, err
		}
		return model.NewDigest(model.IdentifierDigestAlgorithm, value), nil
	}
	return model.Digest{}, model.NewDSSError(
		"One of [certDigest, publicKeyDigest, issuerInfo, kid, x509Uri, publicKey] must be defined for a CertificateRef!")
}

// compile-time assertion: a CertificateRefIdentifier is an identifier.
var _ model.Identifier = (*CertificateRefIdentifier)(nil)
