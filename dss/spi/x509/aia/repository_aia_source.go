// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/aia/RepositoryAIASource.java (DSS 6.5.RC1).
//
// Java's abstract methods (getExistingAIAKeys, findCertificates, insertCertificate,
// removeCertificates) have no Go inheritance equivalent. Following the self-registration
// pattern established in Phase 1b (model.TokenBase.InitToken / TokenOverrides), a concrete
// repository embeds RepositoryAIASource and must call InitRepositoryAIASource(self) from its
// constructor so RepositoryAIASource's methods can dispatch to the concrete implementation;
// forgetting to do so panics, matching Token's contract.
package aia

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// RepositoryAIASourceOverrides declares the abstract methods a concrete repository AIA
// source, embedding RepositoryAIASource, must implement.
type RepositoryAIASourceOverrides interface {
	// ExistingAIAKeys returns a list of all existing AIA keys present in the DB.
	ExistingAIAKeys() []string

	// FindCertificates returns a set of certificates from a DB with the given key.
	FindCertificates(key string) []*model.CertificateToken

	// InsertCertificate allows inserting a certificate into the DB.
	InsertCertificate(aiaKey string, certificateToken *model.CertificateToken)

	// RemoveCertificates removes the certificates from the DB with the given aiaKey.
	RemoveCertificates(aiaKey string)
}

// RepositoryAIASource is the abstract repository AIA source.
type RepositoryAIASource struct {
	// ProxiedSource is used to access certificate tokens that are not present in the
	// repository.
	ProxiedSource AIASource

	// overrides points back at the concrete repository; see InitRepositoryAIASource.
	overrides RepositoryAIASourceOverrides
}

// InitRepositoryAIASource registers the concrete repository with its base so that the base can
// dispatch to the abstract methods. Must be called by every concrete repository's constructor.
func (r *RepositoryAIASource) InitRepositoryAIASource(overrides RepositoryAIASourceOverrides) {
	r.overrides = overrides
}

// repositoryAIASourceOverrides returns the registered overrides, panicking when the concrete
// repository forgot to call InitRepositoryAIASource.
func (r *RepositoryAIASource) repositoryAIASourceOverrides() RepositoryAIASourceOverrides {
	if r.overrides == nil {
		panic("RepositoryAIASource was not initialised: the concrete repository must call InitRepositoryAIASource in its constructor")
	}
	return r.overrides
}

// SetProxySource sets a source to access an AIA in case the requested certificates are not
// present in the repository.
func (r *RepositoryAIASource) SetProxySource(proxiedSource AIASource) {
	r.ProxiedSource = proxiedSource
}

// CertificatesByAIA loads the AIA certificates for certificateToken, without forcing a
// refresh.
func (r *RepositoryAIASource) CertificatesByAIA(certificateToken *model.CertificateToken) []*model.CertificateToken {
	return r.CertificatesByAIAWithRefresh(certificateToken, false)
}

// CertificatesByAIAWithRefresh allows populating the source with new AIA certificates obtained
// from the proxied source, by forcing the refresh. Ports the public
// getCertificatesByAIA(CertificateToken, boolean). Panics if certificateToken is nil (Java
// Objects.requireNonNull).
func (r *RepositoryAIASource) CertificatesByAIAWithRefresh(certificateToken *model.CertificateToken, forceRefresh bool) []*model.CertificateToken {
	if certificateToken == nil {
		panic("CertificateToken shall be provided!")
	}
	urls := spi.CertificateExtensionsUtilsCAIssuersAccessUrls(certificateToken)
	if utils.IsCollectionEmpty(urls) {
		// "There is no AIA extension for certificate download." LOG.info dropped, not
		// load-bearing per PORTING.md.
		return nil
	}

	aiaKeys := r.InitCertificateAIAKeys(urls)
	if !forceRefresh {
		aiaCertificates := r.extractAIAFromCacheSource(aiaKeys)
		if utils.IsCollectionNotEmpty(aiaCertificates) {
			// "Certificate tokens with AIA '{}' have been loaded from the cache" LOG.info
			// dropped, not load-bearing per PORTING.md.
			return aiaCertificates
		}
	}
	// forceRefresh == true: "Cache is skipped..." LOG.info dropped, not load-bearing per
	// PORTING.md.

	return r.extractAndInsertCertificatesFromProxiedSource(certificateToken, aiaKeys)
}

// extractAndInsertCertificatesFromProxiedSource extracts a set of CertificateTokens from
// ProxiedSource and inserts/updates values in the cache source if required.
func (r *RepositoryAIASource) extractAndInsertCertificatesFromProxiedSource(certificateToken *model.CertificateToken, aiaKeys []string) []*model.CertificateToken {
	if r.ProxiedSource == nil {
		// "Proxied AIASource is not provided!" LOG.warn dropped, not load-bearing per
		// PORTING.md.
		return nil
	}

	overrides := r.repositoryAIASourceOverrides()

	existingAIAKeys := overrides.ExistingAIAKeys()
	for _, aiaKey := range aiaKeys {
		if repositoryAIASourceContainsString(existingAIAKeys, aiaKey) {
			// "AIA Certificates with key '{}' have been removed from DB" LOG.info dropped,
			// not load-bearing per PORTING.md.
			overrides.RemoveCertificates(aiaKey)
		}
	}

	var result []*model.CertificateToken
	seen := make(map[string]struct{})

	certificatesTokenByAIA := r.ProxiedSource.CertificatesByAIA(certificateToken)
	if utils.IsCollectionNotEmpty(certificatesTokenByAIA) {
		for _, certificate := range certificatesTokenByAIA {
			sourceURL := r.GetCertificateTokenAIAUrl(certificate)
			if sourceURL == "" {
				// "Not able to find AIA CA issuers URL for certificate..." LOG.warn dropped,
				// not load-bearing per PORTING.md.
				return certificatesTokenByAIA
			}
			aiaKey := r.GetAIAKey(sourceURL)
			overrides.InsertCertificate(aiaKey, certificate)
			id := certificate.DSSIDAsString()
			if _, found := seen[id]; !found {
				seen[id] = struct{}{}
				result = append(result, certificate)
			}
		}
		// "CA issuers for a certificate with Id '{}' are added into the cache" LOG.info
		// dropped, not load-bearing per PORTING.md.
	}

	return result
}

// repositoryAIASourceContainsString reports whether values contains value.
func repositoryAIASourceContainsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// GetCertificateTokenAIAUrl returns a caIssuers access URL, preferring the certificate's own
// recorded SourceURL. Returns "" (Java's null) if none can be found.
func (r *RepositoryAIASource) GetCertificateTokenAIAUrl(certificateToken *model.CertificateToken) string {
	sourceURL := certificateToken.SourceURL()
	if sourceURL == "" {
		aiaUrls := spi.CertificateExtensionsUtilsCAIssuersAccessUrls(certificateToken)
		if len(aiaUrls) == 0 {
			// "No AIA distribution points have been found for this certificate Token..."
			// LOG.warn dropped, not load-bearing per PORTING.md.
		} else {
			// For both the single- and multiple-URL cases upstream picks the first URL,
			// logging at DEBUG in the multiple-URL case (dropped, not load-bearing).
			sourceURL = aiaUrls[0]
		}
	}
	return sourceURL
}

// InitCertificateAIAKeys initializes a list of AIA certificate token keys from the given URLs.
func (r *RepositoryAIASource) InitCertificateAIAKeys(aiaUrls []string) []string {
	keys := make([]string, 0, len(aiaUrls))
	for _, url := range aiaUrls {
		keys = append(keys, r.GetAIAKey(url))
	}
	return keys
}

// GetAIAKey creates a key corresponding to the given aiaURL. Panics wrapping the underlying
// error on failure, mirroring Java's unchecked DSSException propagating out of a method with
// no throws clause.
func (r *RepositoryAIASource) GetAIAKey(aiaURL string) string {
	key, err := spi.DSSUtilsSHA1Digest(aiaURL)
	if err != nil {
		panic(err)
	}
	return key
}

// GetUniqueCertificateAiaID generates a unique identifier for the certificateToken and aiaURL
// pair. Panics wrapping the underlying error on failure; see GetAIAKey.
func (r *RepositoryAIASource) GetUniqueCertificateAiaID(certificateToken *model.CertificateToken, aiaURL string) string {
	key, err := spi.DSSUtilsSHA1Digest(certificateToken.DSSIDAsString() + aiaURL)
	if err != nil {
		panic(err)
	}
	return key
}

// extractAIAFromCacheSource ports the private extractAIAFromCacheSource(List<String>),
// de-duplicating by DSSIDAsString() while preserving order (LinkedHashSet semantics).
func (r *RepositoryAIASource) extractAIAFromCacheSource(aiaKeys []string) []*model.CertificateToken {
	overrides := r.repositoryAIASourceOverrides()
	seen := make(map[string]struct{})
	var certificateTokens []*model.CertificateToken
	for _, key := range aiaKeys {
		for _, certificateToken := range overrides.FindCertificates(key) {
			id := certificateToken.DSSIDAsString()
			if _, found := seen[id]; !found {
				seen[id] = struct{}{}
				certificateTokens = append(certificateTokens, certificateToken)
			}
		}
	}
	return certificateTokens
}

var _ AIASource = (*RepositoryAIASource)(nil)
