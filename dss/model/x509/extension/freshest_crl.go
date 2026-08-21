// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/FreshestCRL.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// FreshestCRL is RFC 5280 4.2.1.15. Freshest CRL (a.k.a. Delta CRL Distribution Point).
//
// The freshest CRL extension identifies how delta CRL information is obtained. The
// extension MUST be marked as non-critical by conforming CAs. Further discussion of CRL
// management is contained in Section 5.
type FreshestCRL struct {
	CertificateExtension

	// crlUrls lists Freshest CRL distribution points.
	crlUrls []string
}

// NewFreshestCRL builds a FreshestCRL extension.
func NewFreshestCRL() *FreshestCRL {
	return &FreshestCRL{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnum_FRESHEST_CRL.OID()),
	}
}

// CrlUrls returns a list of Freshest CRL distribution point URLs.
func (f *FreshestCRL) CrlUrls() []string {
	return f.crlUrls
}

// SetCrlUrls sets a list of Freshest CRL distribution point URLs.
func (f *FreshestCRL) SetCrlUrls(crlUrls []string) {
	f.crlUrls = crlUrls
}
