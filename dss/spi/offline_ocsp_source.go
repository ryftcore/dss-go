// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OfflineOCSPSource.java (DSS 6.5.RC1).
package spi

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/utils"
)

// OfflineOCSPSourceBase helps implement an OCSPSource over an already loaded list of basic
// OCSP responses. It is the port of the Java abstract class OfflineOCSPSource, whose state
// and concrete behaviour it carries.
//
// A concrete source embeds it and must register itself with InitOfflineRevocationSource in
// its constructor - directly when it overrides RevocationTokens, with the base value
// otherwise - so that the RevocationToken dispatch of OfflineRevocationSourceBase reaches
// the right implementation.
type OfflineOCSPSourceBase struct {
	OfflineRevocationSourceBase[revocation.OCSP]
}

// NewOfflineOCSPSourceBase builds the base state of an offline OCSP source, whose references
// are matched with an OCSPTokenRefMatcher. Port of the protected OfflineOCSPSource()
// constructor.
func NewOfflineOCSPSourceBase() OfflineOCSPSourceBase {
	return OfflineOCSPSourceBase{
		OfflineRevocationSourceBase: NewOfflineRevocationSourceBase[revocation.OCSP](NewOCSPTokenRefMatcher()),
	}
}

// RevocationTokens returns the already built tokens concerning the certificate, and
// otherwise builds one per contained response that carries a single response for it.
// Port of the getRevocationTokens(CertificateToken, CertificateToken) override.
//
// Panics with the Java messages when either certificate is missing (Objects.requireNonNull).
func (s *OfflineOCSPSourceBase) RevocationTokens(certificate *model.CertificateToken,
	issuer *model.CertificateToken) ([]RevocationToken[revocation.OCSP], error) {
	if certificate == nil {
		panic("The certificate to be verified cannot be null")
	}
	if issuer == nil {
		panic("The issuer of the certificate to be verified cannot be null")
	}

	result := make([]RevocationToken[revocation.OCSP], 0)

	for _, revocationToken := range s.AllRevocationTokens() {
		if certificate.DSSIDAsString() == revocationToken.RelatedCertificateID() {
			result = append(result, revocationToken)
		}
	}

	if utils.IsCollectionEmpty(result) {
		collectedBinaries := s.AllRevocationBinaries()
		// Upstream logs "--> OfflineOCSPSource queried for {} contains: {} element(s)."
		for _, binary := range collectedBinaries {
			ocspBinary, ok := binary.(*OCSPResponseBinary)
			if !ok {
				panic("ClassCastException : the revocation identifier is not an OCSPResponseBinary")
			}
			basicOCSPResp := ocspBinary.BasicOCSPResp()
			latestSingleResponse := DSSRevocationUtilsLatestSingleResponse(basicOCSPResp, certificate, issuer)
			if latestSingleResponse != nil {
				ocspToken, err := NewOCSPToken(basicOCSPResp, latestSingleResponse, certificate, issuer)
				if err != nil {
					return nil, err
				}
				s.AddRevocationWithBinary(ocspToken, ocspBinary)
				result = append(result, ocspToken)
			}
		}
	}

	// Upstream logs "--> OfflineOCSPSource found result(s) : {}".
	return result, nil
}
