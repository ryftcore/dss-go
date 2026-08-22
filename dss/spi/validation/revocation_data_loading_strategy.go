// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/RevocationDataLoadingStrategy.java (DSS 6.5.RC1).
//
// This class allows retrieving of Revocation data from CRL or OCSP sources, based on the
// defined strategy.
//
// NOTE: The implemented object does not require setting of OCSP/CRL/RevocationDataVerifier
// sources on instantiation from the user. All the values are automatically configured and set
// in eu.europa.esig.dss.validation.SignatureValidationContext based on the parameters defined
// in the provided eu.europa.esig.dss.validation.CertificateVerifier.
//
// Java's abstract class + "protected abstract RevocationToken getRevocationToken(...)" become
// the Init/overrides idiom used throughout this port (see model.TokenBase): a
// RevocationDataLoadingStrategyOverrides interface carries the one abstract method, concrete
// strategies (CRLFirstRevocationDataLoadingStrategy, OCSPFirstRevocationDataLoadingStrategy)
// embed RevocationDataLoadingStrategy and call InitRevocationDataLoadingStrategy(self) in their
// constructor; see ocsp_first_revocation_data_loading_strategy_factory.go for how a factory
// hands out the result.
//
// Java's slf4j logging has no Go equivalent and is not ported.
package validation

import (
	"errors"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
)

// RevocationDataLoadingStrategyOverrides carries the operation Java declares abstract on
// RevocationDataLoadingStrategy. Every concrete strategy implements it and registers itself
// via InitRevocationDataLoadingStrategy so RevocationDataLoadingStrategy's own methods can
// call back into the concrete behaviour.
type RevocationDataLoadingStrategyOverrides interface {
	// RevocationToken retrieves a RevocationToken for the given certificateToken. Port of the
	// abstract getRevocationToken(CertificateToken, CertificateToken).
	RevocationToken(certificateToken, issuerCertificateToken *model.CertificateToken) AnyRevocationToken
}

// RevocationDataLoadingStrategy allows retrieving of Revocation data from CRL or OCSP
// sources, based on the defined strategy.
type RevocationDataLoadingStrategy struct {
	// overrides is the concrete strategy registered via InitRevocationDataLoadingStrategy.
	overrides RevocationDataLoadingStrategyOverrides

	// crlSource is the CRL revocation source.
	crlSource spi.RevocationSource[revocation.CRL]

	// ocspSource is the OCSP revocation source.
	ocspSource spi.RevocationSource[revocation.OCSP]

	// revocationDataVerifier is used to verify the validity of obtained revocation data.
	revocationDataVerifier *RevocationDataVerifier

	// fallbackEnabled, when enabled, returns first obtained revocation token, if both OCSP and
	// CRL requests failed.
	fallbackEnabled bool
}

// InitRevocationDataLoadingStrategy registers the concrete strategy embedding this struct, so
// RevocationToken (if forwarded by the embedder) and the check* helpers below can invoke its
// overridden behaviour. Every concrete strategy's constructor must call this, mirroring
// model.TokenBase.InitToken.
func (s *RevocationDataLoadingStrategy) InitRevocationDataLoadingStrategy(overrides RevocationDataLoadingStrategyOverrides) {
	s.overrides = overrides
}

// RevocationToken retrieves a RevocationToken for the given certificateToken by forwarding to
// the registered concrete strategy. Port of the abstract getRevocationToken(CertificateToken,
// CertificateToken), exposed here so a RevocationDataLoadingStrategyFactory can hand out
// *RevocationDataLoadingStrategy uniformly (see
// ocsp_first_revocation_data_loading_strategy_factory.go).
func (s *RevocationDataLoadingStrategy) RevocationToken(certificateToken, issuerCertificateToken *model.CertificateToken) AnyRevocationToken {
	return s.overrides.RevocationToken(certificateToken, issuerCertificateToken)
}

// setCrlSource sets the CRLSource. Package-private in Java (void setCrlSource(...), no access
// modifier).
func (s *RevocationDataLoadingStrategy) setCrlSource(crlSource spi.RevocationSource[revocation.CRL]) {
	s.crlSource = crlSource
}

// setOcspSource sets the OCSPSource. Package-private in Java.
func (s *RevocationDataLoadingStrategy) setOcspSource(ocspSource spi.RevocationSource[revocation.OCSP]) {
	s.ocspSource = ocspSource
}

// setRevocationDataVerifier sets the RevocationDataVerifier. Package-private in Java.
func (s *RevocationDataLoadingStrategy) setRevocationDataVerifier(revocationDataVerifier *RevocationDataVerifier) {
	s.revocationDataVerifier = revocationDataVerifier
}

// setFallbackEnabled sets whether the fallback shall be enabled. When set to true, returns
// the first obtained token, even when it is not acceptable by the verifier.
//
// Default: false - no fallback. If tokens fail the validation, nil is returned.
//
// Package-private in Java.
func (s *RevocationDataLoadingStrategy) setFallbackEnabled(fallbackEnabled bool) {
	s.fallbackEnabled = fallbackEnabled
}

// getControlTime returns a control time to verify validity of revocation data against.
// Port of getControlTime(); protected in Java, current time is returned by default.
func (s *RevocationDataLoadingStrategy) getControlTime() time.Time {
	return time.Now()
}

// checkCRL retrieves and verifies the obtained CRL token.
//
// NOTE: returns only if a valid entry has been obtained!
//
// Port of checkCRL(CertificateToken, CertificateToken); protected in Java. Java's
// try/catch(DSSException) around the source call becomes a recover() here, since
// spi.RevocationSource.RevocationToken turns a failed extraction into a panic (per
// spi/revocation_source.go's own header comment) rather than returning an error.
func (s *RevocationDataLoadingStrategy) checkCRL(certificateToken, issuerToken *model.CertificateToken) (revocationToken spi.RevocationToken[revocation.CRL]) {
	if s.crlSource == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			var dssErr *model.DSSError
			if err, ok := r.(error); ok && errors.As(err, &dssErr) {
				revocationToken = nil
				return
			}
			panic(r)
		}
	}()
	return s.crlSource.RevocationToken(certificateToken, issuerToken)
}

// checkOCSP retrieves and verifies the obtained OCSP token.
//
// NOTE: returns only if a valid entry has been obtained!
//
// Port of checkOCSP(CertificateToken, CertificateToken); protected in Java. See checkCRL for
// the DSSException-to-recover() translation.
func (s *RevocationDataLoadingStrategy) checkOCSP(certificateToken, issuerToken *model.CertificateToken) (revocationToken spi.RevocationToken[revocation.OCSP]) {
	if s.ocspSource == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			var dssErr *model.DSSError
			if err, ok := r.(error); ok && errors.As(err, &dssErr) {
				revocationToken = nil
				return
			}
			panic(r)
		}
	}()
	return s.ocspSource.RevocationToken(certificateToken, issuerToken)
}

// isAcceptableToken verifies whether the obtained revocation token is acceptable. Port of
// isAcceptableToken(RevocationToken<?>); protected in Java.
func (s *RevocationDataLoadingStrategy) isAcceptableToken(revocationToken AnyRevocationToken) bool {
	if s.revocationDataVerifier == nil {
		return true
	}
	return s.revocationDataVerifier.IsAcceptable(revocationToken, s.getControlTime())
}
