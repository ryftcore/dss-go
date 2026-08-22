// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/SignatureCertificateSource.java (DSS 6.5.RC1).
package spi

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// SignatureCertificateSourceOverrides is the contract a concrete signature certificate source
// (CAdES/XAdES/JAdES/timestamp, ported in later phases) must implement. Port of the protected
// abstract extractCandidatesForSigningCertificate(CertificateSource) method; dispatched the way
// model.TokenBase dispatches to model.TokenOverrides via InitToken.
type SignatureCertificateSourceOverrides interface {
	// ExtractCandidatesForSigningCertificate extracts candidates to be a signing certificate
	// from the source. signingCertificateSource is optional (nil is a valid value, matching
	// Java's nullable parameter).
	ExtractCandidatesForSigningCertificate(signingCertificateSource CertificateSource) *CandidatesForSigningCertificate
}

// SignatureCertificateSource is a basic skeleton able to retrieve the certificate needed to
// validate a signature from a list; the concrete subclass supplies the list of wrapped
// certificates. Embeds TokenCertificateSource. Port of the abstract class
// SignatureCertificateSource.
type SignatureCertificateSource struct {
	TokenCertificateSource

	// overrides points back at the concrete signature certificate source; see
	// InitSignatureCertificateSource. Left nil (as with a bare "default constructor initializing
	// object with null signing certificate candidates list") panics on first use, matching the
	// TokenBase.InitToken convention used elsewhere in this port.
	overrides SignatureCertificateSourceOverrides

	// candidatesForSigningCertificate is the reference to the object containing all candidates
	// to the signing certificate; caches the result of CandidatesForSigningCertificate.
	candidatesForSigningCertificate *CandidatesForSigningCertificate
}

// InitSignatureCertificateSource registers the concrete signature certificate source with its
// base so that the base can dispatch ExtractCandidatesForSigningCertificate. Every concrete
// subclass constructor must call this once.
func (s *SignatureCertificateSource) InitSignatureCertificateSource(overrides SignatureCertificateSourceOverrides) {
	s.overrides = overrides
}

// SignedDataCertificates retrieves the list of all certificates present in a signed element
// (i.e. the CMS Signed data (CAdES)). Port of getSignedDataCertificates().
func (s *SignatureCertificateSource) SignedDataCertificates() []*model.CertificateToken {
	return s.CertificateTokensByOrigin(enumerations.CertificateOriginSignedData)
}

// KeyInfoCertificates retrieves the list of all certificates present in the KeyInfo element
// (XAdES) (can be unsigned). Port of getKeyInfoCertificates().
func (s *SignatureCertificateSource) KeyInfoCertificates() []*model.CertificateToken {
	return s.CertificateTokensByOrigin(enumerations.CertificateOriginKeyInfo)
}

// CertificateValues retrieves the list of all certificates from CertificateValues
// (XAdES/CAdES). Port of getCertificateValues().
func (s *SignatureCertificateSource) CertificateValues() []*model.CertificateToken {
	return s.CertificateTokensByOrigin(enumerations.CertificateOriginCertificateValues)
}

// AttrAuthoritiesCertValues retrieves the list of all certificates from the
// AttrAuthoritiesCertValues (XAdES). Port of getAttrAuthoritiesCertValues().
func (s *SignatureCertificateSource) AttrAuthoritiesCertValues() []*model.CertificateToken {
	return s.CertificateTokensByOrigin(enumerations.CertificateOriginAttrAuthoritiesCertValues)
}

// TimeStampValidationDataCertValues retrieves the list of all certificates from the
// TimeStampValidationData. Port of getTimeStampValidationDataCertValues().
func (s *SignatureCertificateSource) TimeStampValidationDataCertValues() []*model.CertificateToken {
	return s.CertificateTokensByOrigin(enumerations.CertificateOriginTimestampValidationData)
}

// AnyValidationDataCertValues retrieves the list of all certificates from the
// AnyValidationData element. Port of getAnyValidationDataCertValues().
func (s *SignatureCertificateSource) AnyValidationDataCertValues() []*model.CertificateToken {
	return s.CertificateTokensByOrigin(enumerations.CertificateOriginAnyValidationData)
}

// DSSDictionaryCertValues retrieves the list of all certificates from the DSS dictionary
// (PAdES). Port of getDSSDictionaryCertValues().
func (s *SignatureCertificateSource) DSSDictionaryCertValues() []*model.CertificateToken {
	return s.CertificateTokensByOrigin(enumerations.CertificateOriginDSSDictionary)
}

// VRIDictionaryCertValues retrieves the list of all certificates from the VRI dictionary
// (PAdES). Port of getVRIDictionaryCertValues().
func (s *SignatureCertificateSource) VRIDictionaryCertValues() []*model.CertificateToken {
	return s.CertificateTokensByOrigin(enumerations.CertificateOriginVRIDictionary)
}

// UnprotectedHeaderCertificates retrieves the list of all certificates present in the
// unprotected header parameters (JWS, COSE). Port of getUnprotectedHeaderCertificates().
func (s *SignatureCertificateSource) UnprotectedHeaderCertificates() []*model.CertificateToken {
	return s.CertificateTokensByOrigin(enumerations.CertificateOriginUnprotectedHeader)
}

// SigningCertificateRefs retrieves the list of CertificateRefs for the signing certificate
// (V1/V2). Port of getSigningCertificateRefs().
func (s *SignatureCertificateSource) SigningCertificateRefs() []*CertificateRef {
	return s.CertificateRefsByOrigin(enumerations.CertificateRefOriginSigningCertificate)
}

// CompleteCertificateRefs retrieves the list of CertificateRefs included in the attribute
// complete-certificate-references (CAdES) or the
// CompleteCertificateRefs/CompleteCertificateRefsV2 (XAdES). Port of
// getCompleteCertificateRefs().
func (s *SignatureCertificateSource) CompleteCertificateRefs() []*CertificateRef {
	return s.CertificateRefsByOrigin(enumerations.CertificateRefOriginCompleteCertificateRefs)
}

// AttributeCertificateRefs retrieves the list of CertificateRefs included in the attribute
// attribute-certificate-references (CAdES) or the
// AttributeCertificateRefs/AttributeCertificateRefsV2 (XAdES). Port of
// getAttributeCertificateRefs().
func (s *SignatureCertificateSource) AttributeCertificateRefs() []*CertificateRef {
	return s.CertificateRefsByOrigin(enumerations.CertificateRefOriginAttributeCertificateRefs)
}

// SigningCertificates retrieves the Set of CertificateTokens for the signing certificate
// (V1/V2), keyed by DSSIDAsString() per this package's Set<CertificateToken> convention.
// Port of getSigningCertificates().
func (s *SignatureCertificateSource) SigningCertificates() map[string]*model.CertificateToken {
	return s.FindTokensFromRefs(s.SigningCertificateRefs())
}

// CompleteCertificates retrieves the Set of CertificateTokens according to references included
// in the attribute complete-certificate-references (CAdES) or the
// CompleteCertificateRefs/CompleteCertificateRefsV2 (XAdES), keyed by DSSIDAsString() per this
// package's Set<CertificateToken> convention. Port of getCompleteCertificates().
func (s *SignatureCertificateSource) CompleteCertificates() map[string]*model.CertificateToken {
	return s.FindTokensFromRefs(s.CompleteCertificateRefs())
}

// AttributeCertificates retrieves the Set of CertificateTokens according to references included
// in the attribute attribute-certificate-references (CAdES) or the
// AttributeCertificateRefs/AttributeCertificateRefsV2 (XAdES), keyed by DSSIDAsString() per
// this package's Set<CertificateToken> convention. Port of getAttributeCertificates().
func (s *SignatureCertificateSource) AttributeCertificates() map[string]*model.CertificateToken {
	return s.FindTokensFromRefs(s.AttributeCertificateRefs())
}

// CandidatesForSigningCertificate gets an object containing the signing certificate or
// information indicating why it is impossible to extract it from the signature. If the signing
// certificate is identified then it is cached and subsequent calls return this cached value.
// Never returns nil. Port of getCandidatesForSigningCertificate(CertificateSource).
func (s *SignatureCertificateSource) CandidatesForSigningCertificate(signingCertificateSource CertificateSource) *CandidatesForSigningCertificate {
	if s.candidatesForSigningCertificate == nil {
		if s.overrides == nil {
			panic("SignatureCertificateSource was not initialised: the concrete signature certificate source must call InitSignatureCertificateSource in its constructor")
		}
		s.candidatesForSigningCertificate = s.overrides.ExtractCandidatesForSigningCertificate(signingCertificateSource)
	}
	return s.candidatesForSigningCertificate
}

// InitCandidatesList is used to init the candidates list from a provided signing certificate
// source. Port of the protected initCandidatesList(CertificateSource).
//
// Ported verbatim, including an upstream quirk: when certificate tokens are found via the
// proof-of-possession source's certificate refs, the branch re-adds the CertificateValidity
// entries from the source's own certificate list (`certificates`) rather than from the tokens
// just resolved (`certificateTokens`) - see the Java source. This is reproduced as-is per
// PORTING.md's fidelity requirement rather than silently "fixed".
func (s *SignatureCertificateSource) InitCandidatesList(signingCertificateSource CertificateSource) *CandidatesForSigningCertificate {
	if popCertificateSource, ok := signingCertificateSource.(ProofOfPossessionCertificateSource); ok {
		candidates := NewCandidatesForSigningCertificate()
		certificates := popCertificateSource.Certificates()
		if len(certificates) > 0 {
			for _, certificateToken := range certificates {
				candidates.Add(NewCertificateValidity(certificateToken))
			}
		}
		certificateRefs := popCertificateSource.AllCertificateRefs()
		if len(certificateRefs) > 0 {
			certificateTokens := s.FindTokensFromRefs(certificateRefs)
			if len(certificateTokens) > 0 {
				// NOTE: upstream re-iterates `certificates`, not `certificateTokens`; see the
				// doc comment above.
				for _, certificateToken := range certificates {
					candidates.Add(NewCertificateValidity(certificateToken))
				}
			} else {
				for _, certificateRef := range certificateRefs {
					if certificateRef.PublicKey() != nil {
						candidates.Add(NewCertificateValidityFromPublicKey(certificateRef.PublicKey()))
					}
				}
			}
		}
		return candidates

	} else if listCertificateSource, ok := signingCertificateSource.(*ListCertificateSource); ok {
		for _, certificateSource := range listCertificateSource.Sources() {
			candidates := s.InitCandidatesList(certificateSource)
			if !candidates.IsEmpty() {
				return candidates
			}
		}
	}
	return NewCandidatesForSigningCertificate()
}

// CertificateSourceType returns the certificate source type associated with this
// implementation. Port of the getCertificateSourceType() override.
func (s *SignatureCertificateSource) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceTypeSignature
}
