// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/RepositoryRevocationSource.java (DSS 6.5.RC1).
package spi

import (
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// RepositoryRevocationSourceOverrides captures what RepositoryRevocationSourceBase needs to
// reach through virtual dispatch: the operations Java declares abstract on
// RepositoryRevocationSource<R> (initRevocationTokenKeys, findRevocations, insertRevocation,
// updateRevocation, removeRevocation, getRevocationAccessUrls, getRevocationTokenKey).
type RepositoryRevocationSourceOverrides[R revocation.Revocation] interface {
	// InitRevocationTokenKeys initializes a list of revocation token keys from the given
	// CertificateToken. Port of the abstract initRevocationTokenKeys(CertificateToken).
	InitRevocationTokenKeys(certificateToken *model.CertificateToken) []string
	// FindRevocations finds a list of RevocationTokens in the cache for the given
	// certificateToken with the corresponding key. Port of the abstract
	// findRevocations(String, CertificateToken, CertificateToken).
	FindRevocations(key string, certificateToken, issuerCertToken *model.CertificateToken) []RevocationToken[R]
	// InsertRevocation inserts a new RevocationToken into the cache. Port of the abstract
	// insertRevocation(String, RevocationToken).
	InsertRevocation(revocationKey string, token RevocationToken[R])
	// UpdateRevocation updates the RevocationToken into cache. Port of the abstract
	// updateRevocation(String, RevocationToken).
	UpdateRevocation(revocationKey string, token RevocationToken[R])
	// RemoveRevocation removes the RevocationToken from cache with the given key. Port of the
	// abstract removeRevocation(String).
	RemoveRevocation(revocationKey string)
	// RevocationAccessURLs returns a revocation access URLs of the given revocation type for
	// the provided CertificateToken. Port of the abstract
	// getRevocationAccessUrls(CertificateToken).
	RevocationAccessURLs(certificateToken *model.CertificateToken) []string
	// RevocationTokenKey gets a unique revocation token identifier used to store the
	// revocation token for certificateToken within a repository. Port of the abstract
	// getRevocationTokenKey(CertificateToken, String).
	RevocationTokenKey(certificateToken *model.CertificateToken, urlString string) string
}

// RepositoryRevocationSourceBase allows storing and retrieving of revocation data to/from a
// repository (e.g. database). Concrete sources (FileRevocationSourceBase, JdbcRevocationSourceBase)
// embed it and register themselves with InitRepositoryRevocationSource.
type RepositoryRevocationSourceBase[R revocation.Revocation] struct {
	// overrides points back at the concrete source; see InitRepositoryRevocationSource.
	overrides RepositoryRevocationSourceOverrides[R]

	// proxiedSource is the data source used to access a revocation token that is not present
	// in the repository.
	proxiedSource RevocationSource[R]

	// defaultNextUpdateDelay is the default cache delay (milliseconds) in case of a nil
	// nextUpdate in the revocation data; nil when unset.
	defaultNextUpdateDelay *int64
	// maxNextUpdateDelay is the maximum cache delay (milliseconds) for the revocation data;
	// nil when unset.
	maxNextUpdateDelay *int64
	// removeExpired, if true, removes revocation tokens from the repository whose nextUpdate
	// is before the current date.
	removeExpired bool
}

// NewRepositoryRevocationSourceBase instantiates the base state of a repository revocation
// source with the Java default values (removeExpired: TRUE). Port of the protected default
// constructor; the concrete source must still call InitRepositoryRevocationSource.
func NewRepositoryRevocationSourceBase[R revocation.Revocation]() RepositoryRevocationSourceBase[R] {
	return RepositoryRevocationSourceBase[R]{removeExpired: true}
}

// InitRepositoryRevocationSource registers the concrete source with its base so that the base
// can dispatch to the abstract Java operations. It must be called by the outermost concrete
// source's constructor before the source is used.
func (s *RepositoryRevocationSourceBase[R]) InitRepositoryRevocationSource(overrides RepositoryRevocationSourceOverrides[R]) {
	s.overrides = overrides
}

// repositoryRevocationSourceBaseOverrides returns the registered overrides, panicking when the
// concrete source forgot to call InitRepositoryRevocationSource.
func (s *RepositoryRevocationSourceBase[R]) repositoryRevocationSourceBaseOverrides() RepositoryRevocationSourceOverrides[R] {
	if s.overrides == nil {
		panic("RepositoryRevocationSource was not initialised: the concrete source must call InitRepositoryRevocationSource in its constructor")
	}
	return s.overrides
}

// SetDefaultNextUpdateDelay sets the default next update delay for the cached files, in
// seconds. If more time has passed from the revocation token's thisUpdate and next update time
// is not specified, then a fresh copy is downloaded and cached, otherwise a cached copy is
// used: if revocation.nextUpdate == nil, then nextUpdate = revocation.thisUpdate +
// defaultNextUpdateDelay. Port of setDefaultNextUpdateDelay(Long).
func (s *RepositoryRevocationSourceBase[R]) SetDefaultNextUpdateDelay(defaultNextUpdateDelaySeconds *int64) {
	s.defaultNextUpdateDelay = repositoryRevocationSourceToMillis(defaultNextUpdateDelaySeconds)
}

// SetMaxNextUpdateDelay sets the maximum allowed nextUpdate delay for cached files, in seconds.
// Allows forcing a refresh in case of long periods between revocation publication (e.g. 6
// months for an ARL): if revocation.nextUpdate > revocation.thisUpdate + maxNextUpdateDelay,
// then nextUpdate = revocation.thisUpdate + maxNextUpdateDelay. Port of
// setMaxNextUpdateDelay(Long).
func (s *RepositoryRevocationSourceBase[R]) SetMaxNextUpdateDelay(maxNextUpdateDelaySeconds *int64) {
	s.maxNextUpdateDelay = repositoryRevocationSourceToMillis(maxNextUpdateDelaySeconds)
}

// repositoryRevocationSourceToMillis converts a *int64 of seconds to a *int64 of milliseconds,
// preserving a nil input.
func repositoryRevocationSourceToMillis(seconds *int64) *int64 {
	if seconds == nil {
		return nil
	}
	millis := *seconds * 1000
	return &millis
}

// SetProxySource sets the proxied revocation source to be called if the data is not available
// in the cache. Port of setProxySource(RevocationSource).
func (s *RepositoryRevocationSourceBase[R]) SetProxySource(proxiedSource RevocationSource[R]) {
	s.proxiedSource = proxiedSource
}

// SetRemoveExpired sets whether the expired revocation data shall be removed from the cache.
// Default: TRUE. Port of setRemoveExpired(boolean).
func (s *RepositoryRevocationSourceBase[R]) SetRemoveExpired(removeExpired bool) {
	s.removeExpired = removeExpired
}

// RevocationToken retrieves a revocation token for the given CertificateToken, using the cache.
// Port of the getRevocationToken(CertificateToken, CertificateToken) override.
func (s *RepositoryRevocationSourceBase[R]) RevocationToken(certificateToken, issuerCertificateToken *model.CertificateToken) RevocationToken[R] {
	return s.RevocationTokenForceRefresh(certificateToken, issuerCertificateToken, false)
}

// RevocationTokenForceRefresh retrieves a revocation token for the given CertificateToken;
// forceRefresh, if true, explicitly skips the cache. Port of getRevocationToken(CertificateToken,
// CertificateToken, boolean).
func (s *RepositoryRevocationSourceBase[R]) RevocationTokenForceRefresh(certificateToken, issuerCertificateToken *model.CertificateToken, forceRefresh bool) RevocationToken[R] {
	revocationTokens := s.RevocationTokensForceRefresh(certificateToken, issuerCertificateToken, forceRefresh)
	if utils.IsCollectionNotEmpty(revocationTokens) {
		if len(revocationTokens) == 1 {
			return revocationTokens[0]
		}
		// Upstream logs "More than one revocation token has been found for certificate with
		// Id '{}'. Return the latest revocation data."
		return s.latestRevocationData(revocationTokens)
	}
	return nil
}

// RevocationTokens retrieves a list of revocation tokens for the given CertificateToken, using
// the cache. Port of the getRevocationTokens(CertificateToken, CertificateToken) override.
//
// Java declares no throws clause here, and the cache/proxied-source path never raises a
// DSSException upstream either; the error return only exists to satisfy the
// MultipleRevocationSource[R] interface and is always nil.
func (s *RepositoryRevocationSourceBase[R]) RevocationTokens(certificateToken, issuerCertificateToken *model.CertificateToken) ([]RevocationToken[R], error) {
	return s.RevocationTokensForceRefresh(certificateToken, issuerCertificateToken, false), nil
}

// RevocationTokensForceRefresh retrieves a list of revocation tokens for the given
// CertificateToken; forceRefresh, if true, explicitly skips the cache. Port of
// getRevocationTokens(CertificateToken, CertificateToken, boolean).
func (s *RepositoryRevocationSourceBase[R]) RevocationTokensForceRefresh(certificateToken, issuerCertificateToken *model.CertificateToken, forceRefresh bool) []RevocationToken[R] {
	if certificateToken == nil || issuerCertificateToken == nil {
		// Upstream logs "Certificate token or issuer's certificate token is null. Cannot get a
		// revocation token!"
		return nil
	}

	overrides := s.repositoryRevocationSourceBaseOverrides()
	keys := overrides.InitRevocationTokenKeys(certificateToken)
	if forceRefresh {
		// Upstream logs "Cache is skipped to retrieve the revocation token for certificate
		// with Id '{}'"
	} else {
		cachedRevocationTokensMap := s.extractRevocationFromCacheSource(certificateToken, issuerCertificateToken, keys)
		keys = nil
		for key := range cachedRevocationTokensMap {
			keys = append(keys, key)
		}
		if utils.IsMapNotEmpty(cachedRevocationTokensMap) {
			var flattened []RevocationToken[R]
			for _, tokens := range cachedRevocationTokensMap {
				flattened = append(flattened, tokens...)
			}
			return flattened
		}
	}

	revocationToken := s.extractAndInsertRevocationTokenFromProxiedSource(certificateToken, issuerCertificateToken, keys)
	if revocationToken != nil {
		return []RevocationToken[R]{revocationToken}
	}
	return nil
}

// extractRevocationFromCacheSource returns a map of correspondence between requested revocation
// keys and extracted, still fresh, revocation data tokens. The map contains entries only for
// keys with available and still fresh revocation data. Port of the private
// extractRevocationFromCacheSource(CertificateToken, CertificateToken, Collection<String>).
func (s *RepositoryRevocationSourceBase[R]) extractRevocationFromCacheSource(
	certificateToken, issuerCertificateToken *model.CertificateToken, keys []string) map[string][]RevocationToken[R] {
	overrides := s.repositoryRevocationSourceBaseOverrides()
	result := make(map[string][]RevocationToken[R])
	for _, key := range keys {
		revocationTokens := overrides.FindRevocations(key, certificateToken, issuerCertificateToken)
		if utils.IsCollectionNotEmpty(revocationTokens) {
			var fresh []RevocationToken[R]
			for _, r := range revocationTokens {
				if s.isNotExpired(r, issuerCertificateToken) {
					fresh = append(fresh, r)
				}
			}
			if utils.IsCollectionNotEmpty(fresh) {
				result[key] = fresh
			} else {
				// Upstream logs "Revocation token is expired in the cache for certificate with
				// Id '{}'"
				if s.removeExpired {
					overrides.RemoveRevocation(key)
				}
			}
		}
	}
	// Upstream logs "Revocation token for certificate with Id '{}' has been loaded from the
	// cache" when result is not empty.
	return result
}

// latestRevocationData returns the RevocationToken with the latest thisUpdate among
// revocationTokens. Port of the private getLatestRevocationData(Collection<RevocationToken>).
func (s *RepositoryRevocationSourceBase[R]) latestRevocationData(revocationTokens []RevocationToken[R]) RevocationToken[R] {
	var latestRevocationData RevocationToken[R]
	for _, revocationToken := range revocationTokens {
		if latestRevocationData == nil ||
			(!revocationToken.ThisUpdate().IsZero() && latestRevocationData.ThisUpdate().Before(revocationToken.ThisUpdate())) {
			latestRevocationData = revocationToken
		}
	}
	return latestRevocationData
}

// extractAndInsertRevocationTokenFromProxiedSource extracts a RevocationToken from the defined
// proxiedSource and inserts/updates it in the cache source if required. Port of the private
// extractAndInsertRevocationTokenFromProxiedSource(CertificateToken, CertificateToken,
// Collection<String>).
func (s *RepositoryRevocationSourceBase[R]) extractAndInsertRevocationTokenFromProxiedSource(
	certificateToken, issuerCertificateToken *model.CertificateToken, keys []string) RevocationToken[R] {
	if s.proxiedSource == nil {
		// Upstream logs "Proxied revocation source is not initialized for the called
		// RevocationSource!"
		return nil
	}

	revocationToken := s.proxiedSource.RevocationToken(certificateToken, issuerCertificateToken)
	if revocationToken != nil {
		if revocationToken.IsValid() {
			overrides := s.repositoryRevocationSourceBaseOverrides()
			sourceURL := s.revocationSourceURL(certificateToken, revocationToken)
			if sourceURL == "" {
				// Upstream logs "Not able to find revocation source URL for certificate '{}'.
				// Revocation will not be added to the cache"
				return revocationToken
			}
			revocationTokenKey := overrides.RevocationTokenKey(certificateToken, sourceURL)
			if !repositoryRevocationSourceContainsString(keys, revocationTokenKey) {
				overrides.InsertRevocation(revocationTokenKey, revocationToken)
				// Upstream logs "Revocation token for certificate '{}' is added into the cache"
			} else {
				overrides.UpdateRevocation(revocationTokenKey, revocationToken)
				// Upstream logs "Revocation token for certificate '{}' is updated in the cache"
			}
		}
		// Upstream logs "The extracted revocation token with Id '{}' is invalid! Reason: {}"
		// when invalid.
	}
	return revocationToken
}

// revocationSourceURL returns a revocation URL for the given revocationToken. Port of the
// protected getRevocationSourceUrl(CertificateToken, RevocationToken).
func (s *RepositoryRevocationSourceBase[R]) revocationSourceURL(certificateToken *model.CertificateToken, revocationToken RevocationToken[R]) string {
	sourceURL := revocationToken.SourceURL()
	if sourceURL == "" {
		urls := s.repositoryRevocationSourceBaseOverrides().RevocationAccessURLs(certificateToken)
		if len(urls) == 0 {
			// Upstream logs "No revocation distribution points have been found for this
			// certificate Token with ID {}"
		} else {
			sourceURL = urls[0]
			// Upstream logs "There are multiple revocation distribution points..." when
			// len(urls) > 1.
		}
	}
	return sourceURL
}

// isNotExpired checks if the nextUpdate date is currently valid with respect to the
// nextUpdateDelay and maxNextUpdateDelay parameters. Port of the protected
// isNotExpired(RevocationToken, CertificateToken).
func (s *RepositoryRevocationSourceBase[R]) isNotExpired(revocationToken RevocationToken[R], certificateTokenIssuer *model.CertificateToken) bool {
	validationDate := time.Now()

	nextUpdate := revocationToken.NextUpdate()
	if nextUpdate.IsZero() {
		revocationIssuer := revocationToken.IssuerCertificateToken()
		if revocationIssuer == nil {
			revocationIssuer = certificateTokenIssuer
		}
		if !revocationIssuer.IsValidOn(validationDate) {
			return false
		}
	}

	thisUpdate := revocationToken.ThisUpdate()
	if nextUpdate.IsZero() && s.defaultNextUpdateDelay != nil && !thisUpdate.IsZero() {
		nextUpdate = thisUpdate.Add(time.Duration(*s.defaultNextUpdateDelay) * time.Millisecond)
	}
	if !nextUpdate.IsZero() {
		if s.maxNextUpdateDelay != nil && !thisUpdate.IsZero() {
			maxNextUpdate := thisUpdate.Add(time.Duration(*s.maxNextUpdateDelay) * time.Millisecond)
			if nextUpdate.After(maxNextUpdate) {
				nextUpdate = maxNextUpdate
			}
		}
		return nextUpdate.After(validationDate)
	}

	return false
}

// repositoryRevocationSourceContainsString reports whether values contains value.
func repositoryRevocationSourceContainsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// compile-time assertion: a RepositoryRevocationSourceBase is a RevocationSource and a
// MultipleRevocationSource.
var (
	_ RevocationSource[revocation.CRL]         = (*RepositoryRevocationSourceBase[revocation.CRL])(nil)
	_ MultipleRevocationSource[revocation.CRL] = (*RepositoryRevocationSourceBase[revocation.CRL])(nil)
)
