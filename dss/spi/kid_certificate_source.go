// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/KidCertificateSource.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi, so the
// type keeps its Java name unqualified.
package spi

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// KidCertificateSource is the certificate source containing a map of certificates by KIDs.
// The type is used for JAdES and CB-AdES processing.
type KidCertificateSource struct {
	CommonCertificateSource

	// mapByKid maps kids and related certificate tokens.
	mapByKid map[string]*model.CertificateToken
}

// NewKidCertificateSource instantiates an object with an empty map of 'kid' identifiers and
// certificate tokens relation. Port of the default constructor.
func NewKidCertificateSource() *KidCertificateSource {
	return &KidCertificateSource{
		CommonCertificateSource: NewCommonCertificateSource(),
		mapByKid:                make(map[string]*model.CertificateToken),
	}
}

// AddCertificate adds a certificate, generating its 'kid' following the JAdES specification.
// Port of addCertificate(CertificateToken), overriding CommonCertificateSource.
//
// Panics on a nil certificate (Objects.requireNonNull).
func (k *KidCertificateSource) AddCertificate(certificateToAdd *model.CertificateToken) *model.CertificateToken {
	if certificateToAdd == nil {
		panic("The certificate must be filled")
	}
	return k.AddCertificateWithKid(k.generateKidBase64String(certificateToAdd), certificateToAdd)
}

// generateKidBase64String generates the 'kid' value as in IETF RFC 5035.
// Port of the protected generateKidBase64String(CertificateToken).
func (k *KidCertificateSource) generateKidBase64String(signingCertificate *model.CertificateToken) string {
	return utils.ToBase64(DSSUtilsGenerateKid(signingCertificate))
}

// AddCertificateWithKid adds a certificate for a given String 'kid' (JWS).
// Port of addCertificate(String, CertificateToken).
//
// Panics on an empty kid (Objects.requireNonNull).
func (k *KidCertificateSource) AddCertificateWithKid(kid string, certificate *model.CertificateToken) *model.CertificateToken {
	if kid == "" {
		panic("kid cannot be null!")
	}
	addedCertificate := k.CommonCertificateSource.AddCertificate(certificate)
	k.mapByKid[kid] = addedCertificate
	return addedCertificate
}

// AddCertificateWithBinaryKid adds a certificate for a given byte array 'kid' (COSE).
// NOTE: when used, the value of the byte array kid is base64-encoded for preservation within
// the map. Port of addCertificate(byte[], CertificateToken).
//
// Panics on a nil kid (Objects.requireNonNull).
func (k *KidCertificateSource) AddCertificateWithBinaryKid(kid []byte, certificate *model.CertificateToken) *model.CertificateToken {
	if kid == nil {
		panic("kid cannot be null!")
	}
	return k.AddCertificateWithKid(utils.ToBase64(kid), certificate)
}

// CertificateByKid gets a CertificateToken by the given base64-encoded KID value.
// Port of getCertificateByKid(String).
func (k *KidCertificateSource) CertificateByKid(kid string) *model.CertificateToken {
	return k.mapByKid[kid]
}

// CertificateByBinaryKid gets a CertificateToken by the given binary KID value.
// Port of getCertificateByKid(byte[]).
func (k *KidCertificateSource) CertificateByBinaryKid(kid []byte) *model.CertificateToken {
	return k.mapByKid[utils.ToBase64(kid)]
}

// FindTokensFromCertRef also matches the certificate reference's 'kid', in addition to the
// base CommonCertificateSource lookup. Port of findTokensFromCertRef(CertificateRef),
// overriding CommonCertificateSource.
func (k *KidCertificateSource) FindTokensFromCertRef(certificateRef *CertificateRef) map[string]*model.CertificateToken {
	certificates := k.CommonCertificateSource.FindTokensFromCertRef(certificateRef)
	if certificateRef.Kid() != "" {
		if kidCertificate, found := k.mapByKid[certificateRef.Kid()]; found {
			if certificates == nil {
				certificates = make(map[string]*model.CertificateToken)
			}
			certificates[kidCertificate.DSSIDAsString()] = kidCertificate
		}
	}
	return certificates
}

// Reset removes all certificates from the source, including the 'kid' relation.
// Port of the protected reset(), overriding CommonCertificateSource.
func (k *KidCertificateSource) Reset() {
	k.CommonCertificateSource.reset()
	k.mapByKid = make(map[string]*model.CertificateToken)
}

// compile-time assertion: a KidCertificateSource is a CertificateSource.
var _ CertificateSource = (*KidCertificateSource)(nil)
