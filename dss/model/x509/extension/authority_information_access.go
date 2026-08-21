// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/AuthorityInformationAccess.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// AuthorityInformationAccess is RFC 5280 4.2.2.1. Authority Information Access.
//
// The authority information access extension indicates how to access information and
// services for the issuer of the certificate in which the extension appears. Information
// and services may include on-line validation services and CA policy data. (The location
// of CRLs is not specified in this extension; that information is provided by the
// cRLDistributionPoints extension.) This extension may be included in end entity or CA
// certificates. Conforming CAs MUST mark this extension as non-critical.
type AuthorityInformationAccess struct {
	CertificateExtension

	// caIssuers lists certificates that were issued to the CA that issued this certificate.
	caIssuers []string

	// ocsp defines the location of the OCSP responder.
	ocsp []string
}

// NewAuthorityInformationAccess builds an AuthorityInformationAccess extension.
func NewAuthorityInformationAccess() *AuthorityInformationAccess {
	return &AuthorityInformationAccess{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnum_AUTHORITY_INFORMATION_ACCESS.OID()),
	}
}

// CaIssuers returns a list of CA issuers URLs.
func (a *AuthorityInformationAccess) CaIssuers() []string {
	return a.caIssuers
}

// SetCaIssuers sets a list of CA issuers URLs.
func (a *AuthorityInformationAccess) SetCaIssuers(caIssuers []string) {
	a.caIssuers = caIssuers
}

// Ocsp returns a list of OCSP access URLs.
func (a *AuthorityInformationAccess) Ocsp() []string {
	return a.ocsp
}

// SetOcsp sets a list of OCSP access URLs.
func (a *AuthorityInformationAccess) SetOcsp(ocsp []string) {
	a.ocsp = ocsp
}
