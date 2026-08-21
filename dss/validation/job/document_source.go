// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/source/DocumentSource.java (DSS 6.5.RC1).
package job

import "github.com/ryftcore/dss-go/dss/spi"

// DocumentSource represents a Trusted List source.
type DocumentSource struct {
	// url is the document access URL.
	url string

	// certificateSource is the signing certificate source for the current document.
	certificateSource spi.CertificateSource

	// cacheKey is the cached CacheKey value (computed from the url field). cacheKeySet
	// records whether it has been computed yet, standing in for Java's lazy-init null check.
	cacheKey    CacheKey
	cacheKeySet bool
}

// NewDocumentSource creates a DocumentSource with empty values. Port of the default
// constructor.
func NewDocumentSource() *DocumentSource {
	return &DocumentSource{}
}

// Url returns the document URL. Port of getUrl().
func (s *DocumentSource) Url() string {
	return s.url
}

// SetUrl sets the document access URL. Panics if url is empty, mirroring Java's
// Objects.requireNonNull("URL cannot be null."). Port of setUrl(String).
func (s *DocumentSource) SetUrl(url string) {
	if url == "" {
		panic("URL cannot be null.")
	}
	s.url = url
}

// CertificateSource returns the certificate source to be used for document validation. Port
// of getCertificateSource().
func (s *DocumentSource) CertificateSource() spi.CertificateSource {
	return s.certificateSource
}

// SetCertificateSource sets the certificate source to be used for document validation.
// Panics if certificateSource is nil, mirroring Java's Objects.requireNonNull(). Port of
// setCertificateSource(CertificateSource).
func (s *DocumentSource) SetCertificateSource(certificateSource spi.CertificateSource) {
	if certificateSource == nil {
		panic("certificateSource must not be nil")
	}
	s.certificateSource = certificateSource
}

// CacheKey returns the document cache key, computing and caching it from the url on first
// call. Port of getCacheKey().
func (s *DocumentSource) CacheKey() CacheKey {
	if !s.cacheKeySet {
		s.cacheKey = NewCacheKey(s.url)
		s.cacheKeySet = true
	}
	return s.cacheKey
}
