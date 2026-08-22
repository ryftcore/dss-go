// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/TokenCertificateSource.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi, so the
// type keeps its Java name unqualified.
//
// The three Java LinkedHashMaps (SignerIdentifier -> List<CertificateOrigin>,
// CertificateToken -> List<CertificateOrigin>, CertificateRef -> List<CertificateRefOrigin>)
// each key on a type without a comparable Go representation (SignerIdentifier and
// CertificateRef embed []byte fields; CertificateToken is compared via DSSIDAsString(), not
// pointer identity). They are ported as insertion-ordered slices of pairs, searched with the
// key type's own Equals, which preserves both the value-equality lookup semantics and the
// deterministic iteration order LinkedHashMap provides.
package spi

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// tokenCertificateSourceIdentifierEntry pairs a SignerIdentifier with the CertificateOrigins
// it has been added with.
type tokenCertificateSourceIdentifierEntry struct {
	identifier *SignerIdentifier
	origins    []enumerations.CertificateOrigin
}

// tokenCertificateSourceCertificateEntry pairs a CertificateToken with the CertificateOrigins
// it has been added with.
type tokenCertificateSourceCertificateEntry struct {
	token   *model.CertificateToken
	origins []enumerations.CertificateOrigin
}

// tokenCertificateSourceRefEntry pairs a CertificateRef with the CertificateRefOrigins it has
// been added with.
type tokenCertificateSourceRefEntry struct {
	ref     *CertificateRef
	origins []enumerations.CertificateRefOrigin
}

// TokenCertificateSource represents a source of certificates embedded in a token (signature,
// timestamp, ocsp response).
type TokenCertificateSource struct {
	CommonCertificateSource

	// certificateIdentifierOrigins is a map between SignerIdentifiers and CertificateOrigins.
	certificateIdentifierOrigins []tokenCertificateSourceIdentifierEntry

	// certificateOrigins is a map between CertificateTokens and CertificateOrigins.
	certificateOrigins []tokenCertificateSourceCertificateEntry

	// certificateRefOrigins is a map between CertificateRefs and CertificateRefOrigins.
	certificateRefOrigins []tokenCertificateSourceRefEntry
}

// NewTokenCertificateSource builds the base state of a token certificate source, seeding the
// embedded CommonCertificateSource the way Java's implicit super() chain does.
// Port of the protected default constructor.
func NewTokenCertificateSource() TokenCertificateSource {
	return TokenCertificateSource{
		CommonCertificateSource: NewCommonCertificateSource(),
	}
}

// InitTokenCertificateSource resets a TokenCertificateSource to its default (empty) state.
// Port of the protected default constructor; Go has no automatic superclass constructor
// chaining, so embedders should call this from their own constructor. It exists mainly for
// symmetry with the InitToken pattern used elsewhere in this codebase; the zero value of
// TokenCertificateSource is already usable, calling it is not required.
func (s *TokenCertificateSource) InitTokenCertificateSource() {
	s.certificateIdentifierOrigins = nil
	s.certificateOrigins = nil
	s.certificateRefOrigins = nil
}

// AddCertificateIdentifier adds a SignerIdentifier with its origin.
// Port of the protected addCertificateIdentifier(SignerIdentifier, CertificateOrigin).
//
// Panics with the Java messages on nil arguments (Objects.requireNonNull).
func (s *TokenCertificateSource) AddCertificateIdentifier(signerIdentifier *SignerIdentifier, origin enumerations.CertificateOrigin) {
	if signerIdentifier == nil {
		panic("The certificate identifier cannot be null")
	}
	if origin == "" {
		panic("The origin cannot be null")
	}
	for i := range s.certificateIdentifierOrigins {
		entry := &s.certificateIdentifierOrigins[i]
		if entry.identifier.Equals(signerIdentifier) {
			entry.origins = append(entry.origins, origin)
			return
		}
	}
	s.certificateIdentifierOrigins = append(s.certificateIdentifierOrigins, tokenCertificateSourceIdentifierEntry{
		identifier: signerIdentifier,
		origins:    []enumerations.CertificateOrigin{origin},
	})
}

// AddCertificateWithOrigin adds a CertificateToken with its CertificateOrigin.
// Port of the protected addCertificate(CertificateToken, CertificateOrigin) overload.
//
// Panics with the Java messages on nil arguments (Objects.requireNonNull).
func (s *TokenCertificateSource) AddCertificateWithOrigin(certificate *model.CertificateToken, origin enumerations.CertificateOrigin) {
	if certificate == nil {
		panic("The certificate cannot be null")
	}
	if origin == "" {
		panic("The origin cannot be null")
	}
	found := false
	for i := range s.certificateOrigins {
		entry := &s.certificateOrigins[i]
		if entry.token.DSSIDAsString() == certificate.DSSIDAsString() {
			entry.origins = append(entry.origins, origin)
			found = true
			break
		}
	}
	if !found {
		s.certificateOrigins = append(s.certificateOrigins, tokenCertificateSourceCertificateEntry{
			token:   certificate,
			origins: []enumerations.CertificateOrigin{origin},
		})
	}
	s.AddCertificate(certificate)
}

// AddCertificateRef adds a CertificateRef with its CertificateRefOrigin.
// Port of the protected addCertificateRef(CertificateRef, CertificateRefOrigin).
//
// Panics with the Java messages on nil arguments (Objects.requireNonNull).
func (s *TokenCertificateSource) AddCertificateRef(certificateRef *CertificateRef, origin enumerations.CertificateRefOrigin) {
	if certificateRef == nil {
		panic("The certificateRef cannot be null")
	}
	if origin == "" {
		panic("The origin cannot be null")
	}
	for i := range s.certificateRefOrigins {
		entry := &s.certificateRefOrigins[i]
		if entry.ref.Equals(certificateRef) {
			entry.origins = append(entry.origins, origin)
			return
		}
	}
	s.certificateRefOrigins = append(s.certificateRefOrigins, tokenCertificateSourceRefEntry{
		ref:     certificateRef,
		origins: []enumerations.CertificateRefOrigin{origin},
	})
}

// ReferencesForCertificateToken returns the list of CertificateRefs found for the given
// certificateToken. Port of getReferencesForCertificateToken(CertificateToken).
func (s *TokenCertificateSource) ReferencesForCertificateToken(certificateToken *model.CertificateToken) []*CertificateRef {
	var result []*CertificateRef
	for _, entry := range s.certificateRefOrigins {
		if s.doesCertificateReferenceMatch(certificateToken, entry.ref) {
			result = append(result, entry.ref)
		}
	}
	return result
}

// FindTokensFromRefs returns the set of CertificateTokens for the provided CertificateRefs.
// Port of findTokensFromRefs(Collection<CertificateRef>).
func (s *TokenCertificateSource) FindTokensFromRefs(certificateRefs []*CertificateRef) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, certificateRef := range certificateRefs {
		for key, token := range s.FindTokensFromCertRef(certificateRef) {
			result[key] = token
		}
	}
	return result
}

// AllCertificateIdentifiers returns the set of all SignerIdentifiers.
// For CAdES/PAdES/Timestamp. Port of getAllCertificateIdentifiers().
func (s *TokenCertificateSource) AllCertificateIdentifiers() []*SignerIdentifier {
	result := make([]*SignerIdentifier, 0, len(s.certificateIdentifierOrigins))
	for _, entry := range s.certificateIdentifierOrigins {
		result = append(result, entry.identifier)
	}
	return result
}

// CurrentCertificateIdentifier returns the current SignerIdentifier.
// For CAdES/PAdES/Timestamp. Port of getCurrentCertificateIdentifier().
//
// Panics if more than one current CertificateIdentifier is found (Java's IllegalStateException).
func (s *TokenCertificateSource) CurrentCertificateIdentifier() *SignerIdentifier {
	var current *SignerIdentifier
	for _, signerIdentifier := range s.AllCertificateIdentifiers() {
		if signerIdentifier.IsCurrent() {
			if current != nil {
				panic("More than one current CertificateIdentifier")
			}
			current = signerIdentifier
		}
	}
	return current
}

// AllCertificateRefs returns the set of all certificate references.
// Port of getAllCertificateRefs().
func (s *TokenCertificateSource) AllCertificateRefs() []*CertificateRef {
	result := make([]*CertificateRef, 0, len(s.certificateRefOrigins))
	for _, entry := range s.certificateRefOrigins {
		result = append(result, entry.ref)
	}
	return result
}

// OrphanCertificateRefs returns the list of orphan certificate refs.
// Port of getOrphanCertificateRefs().
func (s *TokenCertificateSource) OrphanCertificateRefs() []*CertificateRef {
	var result []*CertificateRef
	for _, entry := range s.certificateRefOrigins {
		if s.IsOrphan(entry.ref) {
			result = append(result, entry.ref)
		}
	}
	return result
}

// IsOrphan verifies whether the CertificateRef is orphan (does not have related certificates
// embedded to the certificate source). Port of the protected isOrphan(CertificateRef).
func (s *TokenCertificateSource) IsOrphan(certificateRef *CertificateRef) bool {
	for _, entry := range s.certificateOrigins {
		if s.doesCertificateReferenceMatch(entry.token, certificateRef) {
			return false
		}
	}
	return true
}

// CertificateTokenBySignerIdentifier gets a CertificateToken by the given SignerIdentifier.
// Port of the protected getCertificateToken(SignerIdentifier).
//
// A SignerIdentifier.IsRelatedToCertificate error (an unreadable subjectKeyIdentifier
// extension) is treated as "not related", matching Java's letting the unchecked DSSException
// propagate only when it would in fact be thrown for every candidate; here it simply excludes
// the unreadable candidate from the search rather than aborting it.
func (s *TokenCertificateSource) CertificateTokenBySignerIdentifier(signerIdentifier *SignerIdentifier) *model.CertificateToken {
	for _, entry := range s.certificateOrigins {
		if related, err := signerIdentifier.IsRelatedToCertificate(entry.token); err == nil && related {
			return entry.token
		}
	}
	return nil
}

// CertificateTokensByOrigin gets a list of CertificateTokens by the given CertificateOrigin.
// Port of the protected getCertificateTokensByOrigin(CertificateOrigin).
func (s *TokenCertificateSource) CertificateTokensByOrigin(origin enumerations.CertificateOrigin) []*model.CertificateToken {
	var result []*model.CertificateToken
	for _, entry := range s.certificateOrigins {
		if tokenCertificateSourceContainsOrigin(entry.origins, origin) {
			result = append(result, entry.token)
		}
	}
	return result
}

// CertificateRefsByOrigin gets a list of CertificateRefs by the given CertificateRefOrigin.
// Port of the protected getCertificateRefsByOrigin(CertificateRefOrigin).
func (s *TokenCertificateSource) CertificateRefsByOrigin(origin enumerations.CertificateRefOrigin) []*CertificateRef {
	var result []*CertificateRef
	for _, entry := range s.certificateRefOrigins {
		if tokenCertificateSourceContainsRefOrigin(entry.origins, origin) {
			result = append(result, entry.ref)
		}
	}
	return result
}

// CertificateRefOrigins extracts origins for a given certificateRef.
// Port of getCertificateRefOrigins(CertificateRef).
func (s *TokenCertificateSource) CertificateRefOrigins(certificateRef *CertificateRef) []enumerations.CertificateRefOrigin {
	for _, entry := range s.certificateRefOrigins {
		if entry.ref.Equals(certificateRef) {
			return entry.origins
		}
	}
	return nil
}

// tokenCertificateSourceContainsOrigin reports whether origins contains origin.
func tokenCertificateSourceContainsOrigin(origins []enumerations.CertificateOrigin, origin enumerations.CertificateOrigin) bool {
	for _, o := range origins {
		if o == origin {
			return true
		}
	}
	return false
}

// tokenCertificateSourceContainsRefOrigin reports whether origins contains origin.
func tokenCertificateSourceContainsRefOrigin(origins []enumerations.CertificateRefOrigin, origin enumerations.CertificateRefOrigin) bool {
	for _, o := range origins {
		if o == origin {
			return true
		}
	}
	return false
}
