// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/ListCertificateSource.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi, so the
// type keeps its Java name unqualified.
package spi

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ListCertificateSource operates on several CertificateSources with the composite design
// pattern.
type ListCertificateSource struct {
	// sources is the list of embedded certificate sources.
	sources []CertificateSource
}

// NewListCertificateSource instantiates an empty ListCertificateSource.
// Port of the default constructor.
func NewListCertificateSource() *ListCertificateSource {
	return &ListCertificateSource{}
}

// NewListCertificateSourceFromOne instantiates a ListCertificateSource with one CertificateSource.
// Port of the ListCertificateSource(CertificateSource) constructor.
func NewListCertificateSourceFromOne(certificateSource CertificateSource) *ListCertificateSource {
	l := NewListCertificateSource()
	l.Add(certificateSource)
	return l
}

// NewListCertificateSourceFromSources instantiates a ListCertificateSource with a slice of
// CertificateSources. Port of the ListCertificateSource(CertificateSource...) and
// ListCertificateSource(List<CertificateSource>) constructors, which Go's variadic parameter
// collapses into one.
func NewListCertificateSourceFromSources(certificateSources ...CertificateSource) *ListCertificateSource {
	l := NewListCertificateSource()
	l.AddAllSources(certificateSources...)
	return l
}

// Add adds a certificate source to the list. Port of add(CertificateSource); returns whether
// the certificateSource has been added successfully.
func (l *ListCertificateSource) Add(certificateSource CertificateSource) bool {
	if certificateSource == nil {
		return false
	}
	for _, existing := range l.sources {
		if existing == certificateSource {
			return false
		}
	}
	l.sources = append(l.sources, certificateSource)
	return true
}

// AddAll adds the embedded certificate sources of listCertificateSource to the list of
// certificate sources. Port of addAll(ListCertificateSource).
func (l *ListCertificateSource) AddAll(listCertificateSource *ListCertificateSource) {
	if listCertificateSource != nil {
		l.AddAllSources(listCertificateSource.Sources()...)
	}
}

// AddAllSources adds a slice of certificate sources. Port of addAll(List<CertificateSource>)
// and addAll(CertificateSource...), which Go's variadic parameter collapses into one.
func (l *ListCertificateSource) AddAllSources(certificateSources ...CertificateSource) {
	for _, certificateSource := range certificateSources {
		l.Add(certificateSource)
	}
}

// Sources returns an unmodifiable list of embedded CertificateSources. Port of getSources().
func (l *ListCertificateSource) Sources() []CertificateSource {
	result := make([]CertificateSource, len(l.sources))
	copy(result, l.sources)
	return result
}

// IsEmpty checks if the embedded sources is empty. Port of isEmpty().
func (l *ListCertificateSource) IsEmpty() bool {
	return len(l.sources) == 0
}

// AreAllCertSourcesTrusted checks if the ListCertificateSource contains only trusted
// CertificateSources. Port of areAllCertSourcesTrusted().
func (l *ListCertificateSource) AreAllCertSourcesTrusted() bool {
	for _, certificateSource := range l.sources {
		if !certificateSource.CertificateSourceType().IsTrusted() {
			return false
		}
	}
	return true
}

// ContainsTrustedCertSources verifies if the current list of certificate sources contains a
// trusted certificate source. Port of containsTrustedCertSources().
func (l *ListCertificateSource) ContainsTrustedCertSources() bool {
	for _, certificateSource := range l.sources {
		if certificateSource.CertificateSourceType().IsTrusted() {
			return true
		}
	}
	return false
}

// AddCertificate always panics: a ListCertificateSource cannot be added to directly.
// Port of addCertificate(CertificateToken), which throws UnsupportedOperationException.
func (l *ListCertificateSource) AddCertificate(certificate *model.CertificateToken) *model.CertificateToken {
	panic("Cannot add a new certificate to a ListCertificateSource!")
}

// CertificateSourceType always panics: use CertificateSourceTypeOf instead.
// Port of getCertificateSourceType(), which throws UnsupportedOperationException.
func (l *ListCertificateSource) CertificateSourceType() enumerations.CertificateSourceType {
	panic("getCertificateSourceType() method is not supported in ListCertificateSource! " +
		"Use CertificateSourceTypeOf(CertificateToken certificate) method instead.")
}

// Certificates returns the deduplicated certificates found in all embedded sources.
// Port of getCertificates().
//
// Java's HashMap-backed dedup would iterate in an order that is arbitrary but stable within a
// JVM run; a bare Go map is randomized on every run instead, so the dedup set is kept
// insertion-ordered (slice + index map, PORTING.md's Collections rule) rather than a bare map,
// while the dedup semantics (last write for a given DSSIDAsString() wins) are unchanged.
func (l *ListCertificateSource) Certificates() []*model.CertificateToken {
	seen := utils.NewOrderedMap[string, *model.CertificateToken]()
	for _, certificateSource := range l.sources {
		for _, token := range certificateSource.Certificates() {
			seen.Set(token.DSSIDAsString(), token)
		}
	}
	return seen.Values()
}

// IsTrusted checks in all sources if the given certificate is trusted.
// Port of isTrusted(CertificateToken).
func (l *ListCertificateSource) IsTrusted(certificateToken *model.CertificateToken) bool {
	for _, source := range l.sources {
		if source.IsTrusted(certificateToken) {
			return true
		}
	}
	return false
}

// IsTrustedAtTime checks in all sources if the given certificate is trusted at controlTime.
// Port of isTrustedAtTime(CertificateToken, Date).
func (l *ListCertificateSource) IsTrustedAtTime(certificateToken *model.CertificateToken, controlTime time.Time) bool {
	for _, source := range l.sources {
		if source.IsTrustedAtTime(certificateToken, controlTime) {
			return true
		}
	}
	return false
}

// IsKnown checks in all sources if the given certificate is known.
// Port of isKnown(CertificateToken).
func (l *ListCertificateSource) IsKnown(certificateToken *model.CertificateToken) bool {
	for _, source := range l.sources {
		if source.IsKnown(certificateToken) {
			return true
		}
	}
	return false
}

// IsAllSelfSigned checks in all sources if all embedded certificates are self-signed.
// Port of isAllSelfSigned().
func (l *ListCertificateSource) IsAllSelfSigned() bool {
	for _, certificateSource := range l.sources {
		if !certificateSource.IsAllSelfSigned() {
			return false
		}
	}
	return true
}

// IsCertificateSourceEqual reports whether the two ListCertificateSources contain the same
// certificate tokens. Port of isCertificateSourceEqual(CertificateSource).
func (l *ListCertificateSource) IsCertificateSourceEqual(certificateSource CertificateSource) bool {
	return certificateTokenSetsEqual(l.Certificates(), certificateSource.Certificates())
}

// IsCertificateSourceEquivalent reports whether the two ListCertificateSources contain the
// same entity keys. Port of isCertificateSourceEquivalent(CertificateSource).
func (l *ListCertificateSource) IsCertificateSourceEquivalent(certificateSource CertificateSource) bool {
	first := l.Entities()
	second := certificateSource.Entities()
	if len(first) != len(second) {
		return false
	}
	for _, a := range first {
		found := false
		for _, b := range second {
			if entitiesEqual(a, b) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// certificateTokenSetsEqual reports whether the two slices contain the same certificate
// tokens, ignoring order and duplicates, keyed the way this package keys certificate sets
// (DSSIDAsString()).
func certificateTokenSetsEqual(first, second []*model.CertificateToken) bool {
	firstSet := make(map[string]struct{}, len(first))
	for _, token := range first {
		firstSet[token.DSSIDAsString()] = struct{}{}
	}
	secondSet := make(map[string]struct{}, len(second))
	for _, token := range second {
		secondSet[token.DSSIDAsString()] = struct{}{}
	}
	if len(firstSet) != len(secondSet) {
		return false
	}
	for key := range firstSet {
		if _, found := secondSet[key]; !found {
			return false
		}
	}
	return true
}

// entitiesEqual reports whether the two CertificateSourceEntity values are equal, delegating
// to the concrete type's Equals when it is the package's own *equivalentCertificatesEntity.
func entitiesEqual(a, b CertificateSourceEntity) bool {
	aEntity, aOK := a.(*equivalentCertificatesEntity)
	bEntity, bOK := b.(*equivalentCertificatesEntity)
	if aOK && bOK {
		return aEntity.Equals(bEntity)
	}
	return a == b
}

// CertificateSourceTypeOf returns the different CertificateSourceTypes where the certificate
// is found. Port of getCertificateSourceType(CertificateToken).
func (l *ListCertificateSource) CertificateSourceTypeOf(certificateToken *model.CertificateToken) map[enumerations.CertificateSourceType]struct{} {
	result := make(map[enumerations.CertificateSourceType]struct{})
	for _, source := range l.sources {
		if source.IsKnown(certificateToken) {
			result[source.CertificateSourceType()] = struct{}{}
		}
	}
	return result
}

// ByEntityKey returns the found CertificateTokens from all CertificateSources for the given
// EntityIdentifier. Port of getByEntityKey(EntityIdentifier).
func (l *ListCertificateSource) ByEntityKey(entityKey *model.EntityIdentifier) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, source := range l.sources {
		for key, token := range source.ByEntityKey(entityKey) {
			result[key] = token
		}
	}
	return result
}

// ByPublicKey returns the found CertificateTokens from all CertificateSources for the given
// PublicKey. Port of getByPublicKey(PublicKey).
func (l *ListCertificateSource) ByPublicKey(publicKey *model.PublicKey) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, source := range l.sources {
		for key, token := range source.ByPublicKey(publicKey) {
			result[key] = token
		}
	}
	return result
}

// BySki returns the found CertificateTokens from all CertificateSources for the given subject
// key identifier (SHA-1 of the public key). Port of getBySki(byte[]).
func (l *ListCertificateSource) BySki(ski []byte) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, source := range l.sources {
		for key, token := range source.BySki(ski) {
			result[key] = token
		}
	}
	return result
}

// FindTokensFromCertRef returns the found CertificateTokens from all CertificateSources for
// the given CertificateRef. Port of findTokensFromCertRef(CertificateRef).
func (l *ListCertificateSource) FindTokensFromCertRef(certificateRef *CertificateRef) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, source := range l.sources {
		for key, token := range source.FindTokensFromCertRef(certificateRef) {
			result[key] = token
		}
	}
	return result
}

// Entities returns the deduplicated certificate source entities found in all sources.
// Port of getEntities().
func (l *ListCertificateSource) Entities() []CertificateSourceEntity {
	var result []CertificateSourceEntity
	for _, certificateSource := range l.sources {
		for _, entity := range certificateSource.Entities() {
			found := false
			for _, existing := range result {
				if entitiesEqual(existing, entity) {
					found = true
					break
				}
			}
			if !found {
				result = append(result, entity)
			}
		}
	}
	return result
}

// BySubject returns the found CertificateTokens from all CertificateSources for the given
// X500PrincipalHelper. Port of getBySubject(X500PrincipalHelper).
func (l *ListCertificateSource) BySubject(subject *model.X500PrincipalHelper) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, source := range l.sources {
		for key, token := range source.BySubject(subject) {
			result[key] = token
		}
	}
	return result
}

// BySignerIdentifier returns the found CertificateTokens from all CertificateSources for the
// given SignerIdentifier. Port of getBySignerIdentifier(SignerIdentifier).
func (l *ListCertificateSource) BySignerIdentifier(signerIdentifier *SignerIdentifier) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, source := range l.sources {
		for key, token := range source.BySignerIdentifier(signerIdentifier) {
			result[key] = token
		}
	}
	return result
}

// ByCertificateDigest returns the found CertificateTokens from all CertificateSources for the
// given Digest. Port of getByCertificateDigest(Digest).
func (l *ListCertificateSource) ByCertificateDigest(digest model.Digest) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, source := range l.sources {
		for key, token := range source.ByCertificateDigest(digest) {
			result[key] = token
		}
	}
	return result
}

// NumberOfSources returns the number of set CertificateSources. Port of getNumberOfSources().
func (l *ListCertificateSource) NumberOfSources() int {
	return len(l.sources)
}

// NumberOfCertificates returns the number of found CertificateTokens in all sources.
// Port of getNumberOfCertificates().
func (l *ListCertificateSource) NumberOfCertificates() int {
	return len(l.Certificates())
}

// NumberOfEntities returns the number of found CertificateSourceEntities in all sources.
// Port of getNumberOfEntities().
func (l *ListCertificateSource) NumberOfEntities() int {
	return len(l.Entities())
}

// compile-time assertion: a ListCertificateSource is a CertificateSource.
var _ CertificateSource = (*ListCertificateSource)(nil)
