// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/ValidityAssuredShortTerm.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// ValidityAssuredShortTerm is ETSI EN 319 412-1 "5.2.2 Validity Assured - Short Term".
//
// This extension indicates that the validity of the certificate is assured because the
// certificate is a "short-term certificate". That is, the time as indicated in the
// certificate attribute from notBefore through notAfter, inclusive, is shorter than the
// maximum time to process a revocation request as specified by the certificate practice
// statement or certificate policy.
type ValidityAssuredShortTerm struct {
	CertificateExtension

	// valAssuredSTCerts defines the value of the ext-etsi-valassured-ST-certs extension.
	valAssuredSTCerts bool
}

// NewValidityAssuredShortTerm builds a ValidityAssuredShortTerm extension.
func NewValidityAssuredShortTerm() *ValidityAssuredShortTerm {
	return &ValidityAssuredShortTerm{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnum_VALIDITY_ASSURED_SHORT_TERM.OID()),
	}
}

// IsValAssuredSTCerts returns the ext-etsi-valassured-ST-certs extension value: TRUE if
// the ext-etsi-valassured-ST-certs extension is present, FALSE otherwise.
func (v *ValidityAssuredShortTerm) IsValAssuredSTCerts() bool {
	return v.valAssuredSTCerts
}

// SetValAssuredSTCerts sets the ext-etsi-valassured-ST-certs extension value.
func (v *ValidityAssuredShortTerm) SetValAssuredSTCerts(valAssuredSTCerts bool) {
	v.valAssuredSTCerts = valAssuredSTCerts
}
