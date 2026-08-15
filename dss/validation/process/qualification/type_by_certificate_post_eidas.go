// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/type/TypeByCertificatePostEIDAS.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
)

// typeByCertificatePostEIDAS gets certificate approval status type after
// eIDAS.
type typeByCertificatePostEIDAS struct {
	// signingCertificate is the certificate to get usage type for.
	signingCertificate *diagnostic.CertificateWrapper
}

// newTypeByCertificatePostEIDAS is the default constructor. Port of
// TypeByCertificatePostEIDAS(CertificateWrapper).
func newTypeByCertificatePostEIDAS(signingCertificate *diagnostic.CertificateWrapper) *typeByCertificatePostEIDAS {
	return &typeByCertificatePostEIDAS{signingCertificate: signingCertificate}
}

// Type is the port of the overridden getType().
func (t *typeByCertificatePostEIDAS) Type() enumerations.CertificateType {
	esign := IsQCTypeEsign(t.signingCertificate)
	eseal := IsQCTypeEseal(t.signingCertificate)
	web := IsQCTypeWeb(t.signingCertificate)

	noneType := !(esign || eseal || web)

	// multiple qcTypes are possible (mistake) but MUST be overruled by the trusted list
	count := 0
	for _, b := range []bool{esign, eseal, web} {
		if b {
			count++
		}
	}
	onlyOne := count == 1

	if (noneType && t.signingCertificate.IsQcCompliance()) || (esign && onlyOne) {
		return enumerations.CertificateType_ESIGN
	} else if eseal && onlyOne {
		return enumerations.CertificateType_ESEAL
	} else if web && onlyOne {
		return enumerations.CertificateType_WSA
	}
	return enumerations.CertificateType_UNKNOWN
}
