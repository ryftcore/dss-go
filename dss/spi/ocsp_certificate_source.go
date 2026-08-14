// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OCSPCertificateSource.java (DSS 6.5.RC1).
package spi

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// OCSPCertificateSource is the source of the certificates embedded in an OCSP token.
type OCSPCertificateSource struct {
	RevocationCertificateSourceBase

	// basicOCSPResp is the basic OCSP response the certificates are extracted from.
	basicOCSPResp *BasicOCSPResp
	// candidatesForSigningCertificate caches the candidates to the response's signing
	// certificate.
	candidatesForSigningCertificate *CandidatesForSigningCertificate
}

// NewOCSPCertificateSource extracts the certificate tokens and the signing certificate
// reference of the given basic OCSP response.
// Port of the OCSPCertificateSource(BasicOCSPResp) constructor.
//
// Panics with the Java message when the response is missing (Objects.requireNonNull). The
// DSSExceptions the certificate extraction raises upstream are returned as errors.
func NewOCSPCertificateSource(basicOCSPResp *BasicOCSPResp) (*OCSPCertificateSource, error) {
	if basicOCSPResp == nil {
		panic("BasicOCSPResp must be provided!")
	}
	source := &OCSPCertificateSource{
		RevocationCertificateSourceBase: NewRevocationCertificateSourceBase(),
		basicOCSPResp:                   basicOCSPResp,
	}
	if err := source.extractCertificateTokens(); err != nil {
		return nil, err
	}
	if err := source.extractCertificateRefs(); err != nil {
		return nil, err
	}
	return source, nil
}

// extractCertificateTokens ports the private extractCertificateTokens().
//
// Upstream iterates the X509CertificateHolders of the response and converts each with
// DSSASN1Utils.getCertificate; the Go BasicOCSPResp hands back parsed certificates, keeping
// each certificate's original DER, and skips the ones crypto/x509 cannot decode.
func (s *OCSPCertificateSource) extractCertificateTokens() error {
	for _, certificate := range s.basicOCSPResp.Certs() {
		certificateToken, err := model.NewCertificateToken(certificate)
		if err != nil {
			return model.NewDSSErrorMessageCause("Unable to read the certificate of the OCSP response", err)
		}
		s.AddCertificateWithOrigin(certificateToken, enumerations.CertificateOrigin_BASIC_OCSP_RESP)
	}
	return nil
}

// extractCertificateRefs ports the private extractCertificateRefs().
func (s *OCSPCertificateSource) extractCertificateRefs() error {
	responderId, err := DSSRevocationUtilsDSSResponderIDFromRespID(s.basicOCSPResp.ResponderID())
	if err != nil {
		return err
	}
	signingCertificateRef := NewCertificateRef()
	signingCertificateRef.SetResponderId(responderId)
	s.AddCertificateRef(signingCertificateRef, enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)
	return nil
}

// CandidatesForSigningCertificate returns the candidates for the signing certificate of the
// OCSP response, computing them on first use. certificateIssuer is the issuer of a
// certificate covered by the OCSP response.
// Port of getCandidatesForSigningCertificate(CertificateToken).
func (s *OCSPCertificateSource) CandidatesForSigningCertificate(
	certificateIssuer *model.CertificateToken) *CandidatesForSigningCertificate {
	if s.candidatesForSigningCertificate == nil {
		s.candidatesForSigningCertificate = s.extractCandidatesForSigningCertificate(certificateIssuer)
	}
	return s.candidatesForSigningCertificate
}

// extractCandidatesForSigningCertificate ports the private
// extractCandidatesForSigningCertificate(CertificateToken).
func (s *OCSPCertificateSource) extractCandidatesForSigningCertificate(
	certificateIssuer *model.CertificateToken) *CandidatesForSigningCertificate {
	candidates := NewCandidatesForSigningCertificate()

	candidates.Add(NewCertificateValidity(certificateIssuer))
	for _, certificateToken := range s.Certificates() {
		candidates.Add(NewCertificateValidity(certificateToken))
	}

	signingCertificateRefs := s.CertificateRefsByOrigin(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)
	if utils.IsCollectionNotEmpty(signingCertificateRefs) {
		signingCertificateRef := signingCertificateRefs[0]
		for _, certificateValidity := range candidates.CertificateValidityList() {
			certificateValidity.SetResponderIdPresent(signingCertificateRef.ResponderId() != nil)

			certificateToken := certificateValidity.CertificateToken()
			if certificateToken != nil {
				certificateValidity.SetResponderIdMatch(
					s.CertificateMatcher().MatchByResponderId(certificateToken, signingCertificateRef))
			}
		}
	}

	return candidates
}

// CertificateSourceType returns OCSP_RESPONSE. Port of the getCertificateSourceType()
// override.
func (s *OCSPCertificateSource) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceType_OCSP_RESPONSE
}

// compile-time assertion: an OCSPCertificateSource is a revocation certificate source.
var _ RevocationCertificateSource = (*OCSPCertificateSource)(nil)
