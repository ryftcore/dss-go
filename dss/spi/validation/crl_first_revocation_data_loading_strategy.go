// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/CRLFirstRevocationDataLoadingStrategy.java (DSS 6.5.RC1).
//
// Java declares `RevocationToken getRevocationToken(...)` using the raw type, i.e. a wildcard
// RevocationToken<?> — the concrete CRL-vs-OCSP kind is erased. Go generics cannot express an
// existential/wildcard instantiation of spi.RevocationToken[R], so this file introduces
// AnyRevocationToken, a non-generic interface capturing every RevocationToken[R] member whose
// signature does not depend on R (everything except Equals). *spi.CRLToken and *spi.OCSPToken
// already satisfy it with no changes since Go interface satisfaction is structural.
package validation

import (
	"math/big"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// AnyRevocationToken is the Go analogue of Java's raw/wildcard RevocationToken<?>: every
// member of spi.RevocationToken[R] whose signature does not mention R. Both *spi.CRLToken and
// *spi.OCSPToken satisfy it.
type AnyRevocationToken interface {
	model.Token

	RevocationType() enumerations.RevocationType
	RelatedCertificate() *model.CertificateToken
	SetRelatedCertificate(certificate *model.CertificateToken)
	RelatedCertificateID() string
	IssuerCertificateToken() *model.CertificateToken
	SourceURL() string
	SetSourceURL(sourceURL string)
	Status() enumerations.CertificateStatus
	SetStatus(status enumerations.CertificateStatus)
	ProductionDate() time.Time
	SetProductionDate(productionDate time.Time)
	ThisUpdate() time.Time
	SetThisUpdate(thisUpdate time.Time)
	NextUpdate() time.Time
	SetNextUpdate(nextUpdate time.Time)
	RevocationDate() time.Time
	SetRevocationDate(revocationDate time.Time)
	CRLNumber() *big.Int
	SetCRLNumber(crlNumber *big.Int)
	ExpiredCertsOnCRL() time.Time
	SetExpiredCertsOnCRL(expiredCertsOnCRL time.Time)
	ArchiveCutOff() time.Time
	SetArchiveCutOff(archiveCutOff time.Time)
	CertHashPresent() bool
	SetCertHashPresent(certHashPresent bool)
	CertHashMatch() bool
	SetCertHashMatch(certHashMatch bool)
	Reason() enumerations.RevocationReason
	SetReason(reason enumerations.RevocationReason)
	CertificateSource() spi.RevocationCertificateSource
	Certificates() []*model.CertificateToken
	SetExternalOrigin(origin enumerations.RevocationOrigin)
	ExternalOrigin() enumerations.RevocationOrigin
	IsInternal() bool
}

// CRLFirstRevocationDataLoadingStrategy fetches firstly the CRL response, and if not
// available, tries OCSP, returning the first succeeded result.
type CRLFirstRevocationDataLoadingStrategy struct {
	RevocationDataLoadingStrategy
}

// NewCRLFirstRevocationDataLoadingStrategy is the default constructor.
func NewCRLFirstRevocationDataLoadingStrategy() *CRLFirstRevocationDataLoadingStrategy {
	s := &CRLFirstRevocationDataLoadingStrategy{}
	s.InitRevocationDataLoadingStrategy(s)
	return s
}

// RevocationToken retrieves the revocation token for the given certificate, trying CRL first
// and falling back to OCSP. Port of getRevocationToken(CertificateToken, CertificateToken).
func (s *CRLFirstRevocationDataLoadingStrategy) RevocationToken(certificateToken, issuerToken *model.CertificateToken) AnyRevocationToken {
	crlToken := s.checkCRL(certificateToken, issuerToken)
	if crlToken != nil && s.isAcceptableToken(crlToken) {
		return crlToken
	}
	ocspToken := s.checkOCSP(certificateToken, issuerToken)
	if ocspToken != nil && s.isAcceptableToken(ocspToken) {
		return ocspToken
	}
	if s.fallbackEnabled {
		// return first successful result
		if crlToken != nil {
			return crlToken
		}
		if ocspToken != nil {
			return ocspToken
		}
	}
	return nil
}
