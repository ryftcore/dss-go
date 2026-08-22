// Ported from
// dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESCertificateRefExtractionUtils.java
// (DSS 6.5.RC1).
package jades

import (
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/spi"
)

// JAdESCertificateRefExtractionUtilsCreateCertificateRef parses the xRefs component value and
// returns a CertificateRef, or nil if the value could not be parsed. Port of the public static
// createCertificateRef(Map).
func JAdESCertificateRefExtractionUtilsCreateCertificateRef(certificateRefMap *jose.Object) *spi.CertificateRef {
	digest, ok := DSSJsonUtilsDigest(certificateRefMap)
	if !ok {
		return nil
	}

	certificateRef := spi.NewCertificateRef()
	certificateRef.SetCertDigest(digest)

	kid, _ := certificateRefMap.Value(jose.HeaderKeyID).(string)
	issuerSerial := DSSJsonUtilsIssuerSerial(kid)
	if issuerSerial != nil {
		certificateRef.SetCertificateIdentifier(spi.DSSASN1UtilsToSignerIdentifierFromIssuerSerial(issuerSerial))
	} else {
		certificateRef.SetKid(kid)
	}
	return certificateRef
}
