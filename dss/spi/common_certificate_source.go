// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CommonCertificateSource.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: this file references CertificateRef and CertificateTokenRefMatcher, two
// spi.x509 types flattened into this package and ported in a sibling chunk of phase 2a.
// CertificateTokenRefMatcher's assumed shape, inferred from the Java signature actually called
// below, is a struct with a default constructor NewCertificateTokenRefMatcher() and a method
// Match(*model.CertificateToken, *CertificateRef) bool.
//
// Java's synchronized blocks around the entitiesByEntityKey/entitiesByPublicKey/tokensBySubject
// mutations exist only to make concurrent AddCertificate/removeCertificate calls safe; per
// PORTING.md this is not the lazy-init-caching case that earns a mutex, so the port is not
// goroutine-safe, matching the rest of this package.
//
// Java's Map<EntityIdentifier, ...>/Map<KeyIdentifier, ...>/Map<X500NameIdentifier, ...> rely on
// hashCode()/equals(); since these identifiers wrap a Digest (a []byte-backed, non-comparable
// Go struct), they cannot be used directly as Go map keys. The port keys on AsXmlID() instead -
// the same problem equivalent_certificates_entity.go (chunk X509-B) already solved for
// Set<CertificateToken>.
package spi

import (
	"bytes"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CommonCertificateSource is the common implementation for all CertificateSource. It stores
// added certificates and allows retrieving them with several methods.
type CommonCertificateSource struct {
	// certificateMatcher is used to match CertificateTokens and CertificateRefs.
	certificateMatcher *CertificateTokenRefMatcher

	// entitiesByEntityKey holds entries keyed by AsXmlID() of a hash of the entity key
	// (public key + subject name combination). All entries sharing a key share the same key
	// pair and a subject name.
	//
	// Java's HashMap iteration order is arbitrary but stable within a JVM run; a bare Go map
	// is randomized on every run instead, and Certificates()/Entities() return in this map's
	// iteration order, so it is kept insertion-ordered (slice + index map, PORTING.md's
	// Collections rule) rather than a bare map.
	entitiesByEntityKey *utils.OrderedMap[string, *equivalentCertificatesEntity]

	// entitiesByPublicKey holds entries keyed by AsXmlID() of a hash of a public key. For a
	// same KeyIdentifier, different subject names (and certificates) are possible.
	entitiesByPublicKey map[string]*equivalentCertificatesEntity

	// tokensBySubject holds tokens keyed by AsXmlID() of an X500Name (RDN). For a same
	// SubjectX500Principal, different key pairs (and certificates) are possible. The inner map
	// is keyed by DSSIDAsString(), standing in for Java's Set<CertificateToken>.
	tokensBySubject map[string]map[string]*model.CertificateToken
}

// NewCommonCertificateSource builds the default certificate source. Port of the default
// constructor.
func NewCommonCertificateSource() CommonCertificateSource {
	source := CommonCertificateSource{}
	source.commonCertificateSourceEnsureInitialized()
	return source
}

// commonCertificateSourceEnsureInitialized brings a zero-value CommonCertificateSource into
// the state NewCommonCertificateSource builds.
//
// Java's field initializers run through the implicit super() call of every subclass, so a
// CommonCertificateSource is never observable with null maps. Go has no constructor chaining:
// a subclass that embeds CommonCertificateSource by value and forgets to seed it with
// NewCommonCertificateSource() would start from nil maps, and the first AddCertificate would
// panic with "assignment to entry in nil map". Every method that writes to the state calls
// this first, so the zero value behaves exactly like a freshly constructed one.
func (s *CommonCertificateSource) commonCertificateSourceEnsureInitialized() {
	if s.certificateMatcher == nil {
		s.certificateMatcher = NewCertificateTokenRefMatcher()
	}
	if s.entitiesByEntityKey == nil {
		s.entitiesByEntityKey = utils.NewOrderedMap[string, *equivalentCertificatesEntity]()
	}
	if s.entitiesByPublicKey == nil {
		s.entitiesByPublicKey = make(map[string]*equivalentCertificatesEntity)
	}
	if s.tokensBySubject == nil {
		s.tokensBySubject = make(map[string]map[string]*model.CertificateToken)
	}
}

// AddCertificate adds an external certificate to the source. If the public key is already
// known, the certificate is merged into the CertificateSourceEntity.
// Port of addCertificate(CertificateToken).
//
// Panics with the Java message when certificateToAdd is missing (Objects.requireNonNull).
// DSSASN1UtilsComputeSkiFromCert's error, raised inside newEquivalentCertificatesEntity /
// addEquivalentCertificate for a malformed public key encoding, is a genuinely exceptional
// programmer/data error Java has no room to propagate either (the constructor and
// addEquivalentCertificate are called from contexts with no checked exception), so it panics
// with the Java DSSException message it stands in for.
func (s *CommonCertificateSource) AddCertificate(certificateToAdd *model.CertificateToken) *model.CertificateToken {
	if certificateToAdd == nil {
		panic("The certificate must be filled")
	}
	s.commonCertificateSourceEnsureInitialized()

	entityKey := certificateToAdd.EntityKey()
	entityKeyID := entityKey.AsXmlID()
	if poolEntity, found := s.entitiesByEntityKey.Get(entityKeyID); found {
		if err := poolEntity.addEquivalentCertificate(certificateToAdd); err != nil {
			panic(err.Error())
		}
	} else {
		newEntity, err := newEquivalentCertificatesEntity(certificateToAdd)
		if err != nil {
			panic(err.Error())
		}
		s.entitiesByEntityKey.Set(entityKeyID, newEntity)
	}

	keyIdentifier := model.NewKeyIdentifier(certificateToAdd.PublicKey())
	keyIdentifierID := keyIdentifier.AsXmlID()
	if poolEntity, found := s.entitiesByPublicKey[keyIdentifierID]; found {
		if err := poolEntity.addEquivalentCertificate(certificateToAdd); err != nil {
			panic(err.Error())
		}
	} else {
		newEntity, err := newEquivalentCertificatesEntity(certificateToAdd)
		if err != nil {
			panic(err.Error())
		}
		s.entitiesByPublicKey[keyIdentifierID] = newEntity
	}

	x500NameIdentifierID := model.NewX500NameIdentifier(certificateToAdd.Subject().Principal()).AsXmlID()
	tokens, found := s.tokensBySubject[x500NameIdentifierID]
	if !found {
		tokens = make(map[string]*model.CertificateToken)
		s.tokensBySubject[x500NameIdentifierID] = tokens
	}
	tokens[certificateToAdd.DSSIDAsString()] = certificateToAdd

	return certificateToAdd
}

// removeCertificate removes the corresponding certificate token from the certificate source.
// Port of the protected removeCertificate(CertificateToken).
//
// Panics with the Java message when certificateToRemove is missing (Objects.requireNonNull).
func (s *CommonCertificateSource) removeCertificate(certificateToRemove *model.CertificateToken) {
	if certificateToRemove == nil {
		panic("The certificate must be filled")
	}

	entityKeyID := certificateToRemove.EntityKey().AsXmlID()
	if poolEntity, found := s.entitiesByEntityKey.Get(entityKeyID); found {
		if len(poolEntity.EquivalentCertificates()) == 1 {
			s.entitiesByEntityKey.Delete(entityKeyID)
		} else {
			poolEntity.removeEquivalentCertificate(certificateToRemove)
		}
	}

	keyIdentifierID := model.NewKeyIdentifier(certificateToRemove.PublicKey()).AsXmlID()
	if poolEntity, found := s.entitiesByPublicKey[keyIdentifierID]; found {
		if len(poolEntity.EquivalentCertificates()) == 1 {
			delete(s.entitiesByPublicKey, keyIdentifierID)
		} else {
			poolEntity.removeEquivalentCertificate(certificateToRemove)
		}
	}

	x500NameIdentifierID := model.NewX500NameIdentifier(certificateToRemove.Subject().Principal()).AsXmlID()
	if certificateTokens, found := s.tokensBySubject[x500NameIdentifierID]; found && len(certificateTokens) > 0 {
		if len(certificateTokens) == 1 {
			delete(s.tokensBySubject, x500NameIdentifierID)
		} else {
			delete(certificateTokens, certificateToRemove.DSSIDAsString())
		}
	}
}

// reset removes all certificates from the source. Port of the protected reset().
func (s *CommonCertificateSource) reset() {
	s.entitiesByEntityKey = nil
	s.entitiesByPublicKey = nil
	s.tokensBySubject = nil
	s.commonCertificateSourceEnsureInitialized()
}

// IsKnown checks if a given certificate is known in the current source. Port of isKnown(CertificateToken).
func (s *CommonCertificateSource) IsKnown(token *model.CertificateToken) bool {
	poolEntity, found := s.entitiesByEntityKey.Get(token.EntityKey().AsXmlID())
	if !found {
		return false
	}
	certsByPublicKey := poolEntity.EquivalentCertificates()
	certsBySubject := s.BySubject(token.Subject())
	for id := range certsByPublicKey {
		if _, ok := certsBySubject[id]; ok {
			return true
		}
	}
	return false
}

// Certificates retrieves the unmodifiable list of all certificate tokens from this source.
// Port of getCertificates().
func (s *CommonCertificateSource) Certificates() []*model.CertificateToken {
	var allCertificates []*model.CertificateToken
	for _, entity := range s.entitiesByEntityKey.Values() {
		allCertificates = append(allCertificates, entity.orderedEquivalentCertificates()...)
	}
	return allCertificates
}

// Entities returns a list of certificates grouped by their public keys. Port of getEntities().
func (s *CommonCertificateSource) Entities() []CertificateSourceEntity {
	entityValues := s.entitiesByEntityKey.Values()
	entities := make([]CertificateSourceEntity, 0, len(entityValues))
	for _, entity := range entityValues {
		entities = append(entities, entity)
	}
	return entities
}

// ByPublicKey returns the certificate tokens with the given PublicKey, keyed by DSSIDAsString().
// Port of getByPublicKey(PublicKey).
func (s *CommonCertificateSource) ByPublicKey(publicKey *model.PublicKey) map[string]*model.CertificateToken {
	entity, found := s.entitiesByPublicKey[model.NewKeyIdentifier(publicKey).AsXmlID()]
	if !found {
		return map[string]*model.CertificateToken{}
	}
	return entity.EquivalentCertificates()
}

// ByEntityKey returns the certificate tokens with the given EntityIdentifier, keyed by
// DSSIDAsString(). Port of getByEntityKey(EntityIdentifier).
func (s *CommonCertificateSource) ByEntityKey(entityKey *model.EntityIdentifier) map[string]*model.CertificateToken {
	entity, found := s.entitiesByEntityKey.Get(entityKey.AsXmlID())
	if !found {
		return map[string]*model.CertificateToken{}
	}
	return entity.EquivalentCertificates()
}

// BySki returns the certificate tokens with the given SKI (SubjectKeyIdentifier, SHA-1 of the
// PublicKey), keyed by DSSIDAsString(). Port of getBySki(byte[]).
func (s *CommonCertificateSource) BySki(ski []byte) map[string]*model.CertificateToken {
	for _, entity := range s.entitiesByPublicKey {
		if bytes.Equal(entity.Ski(), ski) {
			return entity.EquivalentCertificates()
		}
	}
	return map[string]*model.CertificateToken{}
}

// BySubject returns the certificate tokens with the same subjectDN, keyed by DSSIDAsString().
// Port of getBySubject(X500PrincipalHelper).
func (s *CommonCertificateSource) BySubject(subject *model.X500PrincipalHelper) map[string]*model.CertificateToken {
	tokensSet, found := s.tokensBySubject[model.NewX500NameIdentifier(subject.Principal()).AsXmlID()]
	if !found {
		return map[string]*model.CertificateToken{}
	}
	return tokensSet
}

// BySignerIdentifier returns the certificate tokens with the given SignerIdentifier, keyed by
// DSSIDAsString(). Port of getBySignerIdentifier(SignerIdentifier).
func (s *CommonCertificateSource) BySignerIdentifier(signerIdentifier *SignerIdentifier) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, entry := range s.entitiesByEntityKey.Values() {
		for _, certificateToken := range entry.EquivalentCertificates() {
			// run over all entries to compare with the SN too
			related, err := signerIdentifier.IsRelatedToCertificate(certificateToken)
			if err == nil && related {
				result[certificateToken.DSSIDAsString()] = certificateToken
			}
		}
	}
	return result
}

// ByCertificateDigest returns the certificate tokens with the given Digest, keyed by
// DSSIDAsString(). Port of getByCertificateDigest(Digest).
func (s *CommonCertificateSource) ByCertificateDigest(digest model.Digest) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, entry := range s.entitiesByEntityKey.Values() {
		for _, certificateToken := range entry.EquivalentCertificates() {
			value, err := certificateToken.Digest(digest.Algorithm())
			if err == nil && bytes.Equal(digest.Value(), value) {
				result[certificateToken.DSSIDAsString()] = certificateToken
			}
		}
	}
	return result
}

// FindTokensFromCertRef returns the certificate tokens for the provided CertificateRef, keyed
// by DSSIDAsString(). Port of findTokensFromCertRef(CertificateRef).
func (s *CommonCertificateSource) FindTokensFromCertRef(certificateRef *CertificateRef) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, entry := range s.entitiesByEntityKey.Values() {
		for _, certificateToken := range entry.EquivalentCertificates() {
			if s.doesCertificateReferenceMatch(certificateToken, certificateRef) {
				result[certificateToken.DSSIDAsString()] = certificateToken
			}
		}
	}
	return result
}

// doesCertificateReferenceMatch verifies whether the CertificateRef matches the
// CertificateToken. Port of the protected doesCertificateReferenceMatch(CertificateToken, CertificateRef).
func (s *CommonCertificateSource) doesCertificateReferenceMatch(certificateToken *model.CertificateToken, certificateRef *CertificateRef) bool {
	s.commonCertificateSourceEnsureInitialized()
	return s.certificateMatcher.Match(certificateToken, certificateRef)
}

// CertificateMatcher returns the CertificateTokenRefMatcher used to match CertificateTokens
// and CertificateRefs. Port of the protected final certificateMatcher field access; exposed
// as an accessor since Go has no protected-field equivalent and subclasses in this package
// (e.g. OCSPCertificateSource) need direct access to it.
func (s *CommonCertificateSource) CertificateMatcher() *CertificateTokenRefMatcher {
	s.commonCertificateSourceEnsureInitialized()
	return s.certificateMatcher
}

// NumberOfCertificates returns the number of stored certificates in this source.
// Port of getNumberOfCertificates().
func (s *CommonCertificateSource) NumberOfCertificates() int {
	return len(s.Certificates())
}

// NumberOfEntities returns the number of stored entities (unique public key) in this source.
// Port of getNumberOfEntities().
func (s *CommonCertificateSource) NumberOfEntities() int {
	return s.entitiesByEntityKey.Len()
}

// CertificateSourceType returns the certificate source type associated with the
// implementation. Port of getCertificateSourceType().
func (s *CommonCertificateSource) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceType_OTHER
}

// IsTrusted checks if a given certificate is trusted. Port of isTrusted(CertificateToken).
func (s *CommonCertificateSource) IsTrusted(certificateToken *model.CertificateToken) bool {
	return false
}

// IsTrustedAtTime checks if a given certificate is trusted at controlTime.
// Port of isTrustedAtTime(CertificateToken, Date).
func (s *CommonCertificateSource) IsTrustedAtTime(certificateToken *model.CertificateToken, controlTime time.Time) bool {
	return s.IsTrusted(certificateToken)
}

// IsAllSelfSigned checks if all certificates are self-signed. Port of isAllSelfSigned().
func (s *CommonCertificateSource) IsAllSelfSigned() bool {
	for _, certificate := range s.Certificates() {
		if !certificate.IsSelfSigned() {
			return false
		}
	}
	return true
}

// IsCertificateSourceEqual checks if the current and the given CertificateSources contain the
// same certificate tokens. Port of isCertificateSourceEqual(CertificateSource).
func (s *CommonCertificateSource) IsCertificateSourceEqual(certificateSource CertificateSource) bool {
	return commonCertificateSourceTokenSetsEqual(s.Certificates(), certificateSource.Certificates())
}

// IsCertificateSourceEquivalent checks if the current and the given CertificateSources
// contain the same entity keys. Port of isCertificateSourceEquivalent(CertificateSource).
func (s *CommonCertificateSource) IsCertificateSourceEquivalent(certificateSource CertificateSource) bool {
	return commonCertificateSourceEntitySetsEqual(s.Entities(), certificateSource.Entities())
}

// commonCertificateSourceTokenSetsEqual reports whether two collections of certificate tokens contain
// the same tokens (by DSSIDAsString()), standing in for
// new HashSet<>(a).equals(new HashSet<>(b)).
func commonCertificateSourceTokenSetsEqual(a, b []*model.CertificateToken) bool {
	if len(a) != len(b) {
		return false
	}
	setA := make(map[string]struct{}, len(a))
	for _, token := range a {
		setA[token.DSSIDAsString()] = struct{}{}
	}
	setB := make(map[string]struct{}, len(b))
	for _, token := range b {
		setB[token.DSSIDAsString()] = struct{}{}
	}
	if len(setA) != len(setB) {
		return false
	}
	for id := range setA {
		if _, found := setB[id]; !found {
			return false
		}
	}
	return true
}

// commonCertificateSourceEntitySetsEqual reports whether two collections of CertificateSourceEntity
// contain the same entities, standing in for new HashSet<>(a).equals(new HashSet<>(b)).
//
// Java relies on EquivalentCertificatesEntity#equals/hashCode (entity identifier equality);
// the port compares pairwise with Equals since CertificateSourceEntity is a marker interface
// with no exported equality method of its own.
func commonCertificateSourceEntitySetsEqual(a, b []CertificateSourceEntity) bool {
	if len(a) != len(b) {
		return false
	}
	matched := make([]bool, len(b))
	for _, entityA := range a {
		concreteA, ok := entityA.(*equivalentCertificatesEntity)
		if !ok {
			return false
		}
		found := false
		for i, entityB := range b {
			if matched[i] {
				continue
			}
			concreteB, ok := entityB.(*equivalentCertificatesEntity)
			if ok && concreteA.Equals(concreteB) {
				matched[i] = true
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

// compile-time assertion: a CommonCertificateSource is a CertificateSource.
var _ CertificateSource = (*CommonCertificateSource)(nil)
