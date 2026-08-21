// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/ExtendedKeyUsages.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// ExtendedKeyUsages is RFC 5280 4.2.1.12. Extended Key Usage.
//
// This extension indicates one or more purposes for which the certified public key may be
// used, in addition to or in place of the basic purposes indicated in the key usage
// extension. In general, this extension will appear only in end entity certificates.
type ExtendedKeyUsages struct {
	CertificateExtension

	// oids lists extended key usage OIDs.
	oids []string
}

// NewExtendedKeyUsages builds an ExtendedKeyUsages extension.
func NewExtendedKeyUsages() *ExtendedKeyUsages {
	return &ExtendedKeyUsages{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnum_EXTENDED_KEY_USAGE.OID()),
	}
}

// Oids returns the extended key usage OIDs.
func (e *ExtendedKeyUsages) Oids() []string {
	return e.oids
}

// SetOids sets the extended key usage OIDs.
func (e *ExtendedKeyUsages) SetOids(oids []string) {
	e.oids = oids
}
