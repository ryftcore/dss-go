// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/OCSPNoCheck.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// OCSPNoCheck is RFC 6960 "4.2.2.2.1. Revocation Checking of an Authorized Responder".
//
// A CA may specify that an OCSP client can trust a responder for the lifetime of the
// responder's certificate. The CA does so by including the extension id-pkix-ocsp-nocheck.
// This SHOULD be a non-critical extension. The value of the extension SHALL be NULL. CAs
// issuing such a certificate should realize that a compromise of the responder's key is as
// serious as the compromise of a CA key used to sign CRLs, at least for the validity
// period of this certificate. CAs may choose to issue this type of certificate with a very
// short lifetime and renew it frequently.
type OCSPNoCheck struct {
	CertificateExtension

	// ocspNoCheckValue defines if the ocsp-nocheck extension is present and set to true.
	ocspNoCheckValue bool
}

// NewOCSPNoCheck builds an OCSPNoCheck extension.
func NewOCSPNoCheck() *OCSPNoCheck {
	return &OCSPNoCheck{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnumOCSPNoCheck.OID()),
	}
}

// IsOcspNoCheck returns the ocsp-nocheck extension value, when present; FALSE if not
// present.
func (o *OCSPNoCheck) IsOcspNoCheck() bool {
	return o.ocspNoCheckValue
}

// SetOcspNoCheck sets the ocsp-nocheck extension value.
func (o *OCSPNoCheck) SetOcspNoCheck(ocspNoCheckValue bool) {
	o.ocspNoCheckValue = ocspNoCheckValue
}
