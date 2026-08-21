// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/crl/OfflineCRLSource.java (DSS 6.5.RC1).
package spi

import (
	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// OfflineCRLSourceBase is the skeleton able to retrieve the needed CRL data from a contained
// list; the concrete source supplies the wrapped CRLs. It is the port of the Java abstract
// class OfflineCRLSource, whose state and concrete behaviour it carries.
//
// A concrete source embeds it and must register itself with InitOfflineRevocationSource in
// its constructor - directly when it overrides RevocationTokens, with the base value
// otherwise - so that the RevocationToken dispatch of OfflineRevocationSourceBase reaches
// the right implementation.
type OfflineCRLSourceBase struct {
	OfflineRevocationSourceBase[revocation.CRL]

	// cachedValidCRLValidities is the cached list of processed CRLValidity values.
	cachedValidCRLValidities []*crlparser.CRLValidity
}

// NewOfflineCRLSourceBase builds the base state of an offline CRL source, whose references
// are matched with a CRLTokenRefMatcher. Port of the protected OfflineCRLSource()
// constructor.
func NewOfflineCRLSourceBase() OfflineCRLSourceBase {
	return OfflineCRLSourceBase{
		OfflineRevocationSourceBase: NewOfflineRevocationSourceBase[revocation.CRL](NewCRLTokenRefMatcher()),
	}
}

// RevocationTokens builds a CRL token per contained CRL that is valid for the certificate's
// issuer. Port of the getRevocationTokens(CertificateToken, CertificateToken) override.
//
// Panics with the Java messages when either certificate is missing (Objects.requireNonNull).
// The DSSException the CRLToken constructor raises for a CRL signed by another issuer is
// returned as an error, as upstream lets it propagate.
func (s *OfflineCRLSourceBase) RevocationTokens(certificateToken *model.CertificateToken,
	issuerToken *model.CertificateToken) ([]RevocationToken[revocation.CRL], error) {
	if certificateToken == nil {
		panic("The certificate to be verified cannot be null")
	}
	if issuerToken == nil {
		panic("The issuer of the certificate to be verified cannot be null")
	}

	result := make([]RevocationToken[revocation.CRL], 0)

	validCRLValiditiesForIssuer := s.fromCachedCRLValidities(issuerToken)

	if len(validCRLValiditiesForIssuer) == 0 {
		collectedBinaries := s.AllRevocationBinaries()
		// Upstream logs "--> OfflineCRLSource queried for {} contains: {} element(s)."

		for _, binary := range collectedBinaries {
			crlBinary, ok := binary.(*crlparser.CRLBinary)
			if !ok {
				panic("ClassCastException : the revocation identifier is not a CRLBinary")
			}
			crlValidity, err := crlparser.CRLUtilsBuildCRLValidity(crlBinary, issuerToken)
			if err != nil {
				// Upstream logs "Unable to retrieve the CRLValidity for CRL with Id '{}' : {}".
				continue
			}
			if crlValidity.IsValid() {
				s.cachedValidCRLValidities = append(s.cachedValidCRLValidities, crlValidity)
				validCRLValiditiesForIssuer = append(validCRLValiditiesForIssuer, crlValidity)
			}
		}
	}

	for _, crlValidity := range validCRLValiditiesForIssuer {
		crlToken, err := NewCRLToken(certificateToken, crlValidity)
		if err != nil {
			return nil, err
		}
		s.AddRevocationWithBinary(crlToken, crlValidity.CrlBinary())
		result = append(result, crlToken)
	}

	// Upstream logs "--> OfflineCRLSource found result(s) : {}".
	return result, nil
}

// fromCachedCRLValidities ports the private getFromCachedCRLValidities(CertificateToken).
func (s *OfflineCRLSourceBase) fromCachedCRLValidities(issuerToken *model.CertificateToken) []*crlparser.CRLValidity {
	result := make([]*crlparser.CRLValidity, 0)
	for _, validity := range s.cachedValidCRLValidities {
		if issuerToken.Equals(validity.IssuerToken()) {
			result = append(result, validity)
		}
	}
	return result
}
