// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/CRLFirstRevocationDataLoadingStrategy.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (flagged per S2B_BRIEF.md, not invented lightly):
//
//  1. Java declares `RevocationToken getRevocationToken(...)` using the raw type, i.e. a
//     wildcard RevocationToken<?> — the concrete CRL-vs-OCSP kind is erased. Go generics
//     cannot express an existential/wildcard instantiation of spi.RevocationToken[R], so this
//     file introduces AnyRevocationToken, a non-generic interface capturing every
//     RevocationToken[R] member whose signature does not depend on R (everything except
//     Equals). *spi.CRLToken and *spi.OCSPToken already satisfy it with no changes since Go
//     interface satisfaction is structural. Sibling chunks needing the same wildcard
//     capability should reuse this type rather than redeclare an equivalent one.
//
//  2. RevocationDataLoadingStrategy (Java abstract class, owned by sibling chunk VAL-D) has
//     the same shape as model.TokenBase/spi.RevocationTokenBase: concrete protected state plus
//     one abstract method. Per the phase 2a handoff fact "every concrete Token constructor
//     calls the Init pattern... TimestampToken must do the same", this file assumes VAL-D
//     followed the same Init/overrides idiom:
//     - type RevocationDataLoadingStrategyOverrides interface {
//     RevocationToken(certificateToken, issuerToken *model.CertificateToken) AnyRevocationToken
//     }
//     - type RevocationDataLoadingStrategy struct{ overrides RevocationDataLoadingStrategyOverrides; crlSource spi.RevocationSource[revocation.CRL]; ocspSource spi.RevocationSource[revocation.OCSP]; revocationDataVerifier *RevocationDataVerifier; fallbackEnabled bool }
//     - (*RevocationDataLoadingStrategy).InitRevocationDataLoadingStrategy(overrides) registers the concrete strategy
//     - (*RevocationDataLoadingStrategy).RevocationToken(...) forwards to overrides.RevocationToken(...), so a factory
//     can return *RevocationDataLoadingStrategy uniformly (see
//     ocsp_first_revocation_data_loading_strategy_factory.go)
//     - (*RevocationDataLoadingStrategy).checkCRL/.checkOCSP/.isAcceptableToken/.getControlTime carry the concrete
//     Java protected method bodies; setCrlSource/setOcspSource/setRevocationDataVerifier/setFallbackEnabled are
//     unexported (package-private in Java, same-package accessible in Go).
package validation

import (
	"math/big"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
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
