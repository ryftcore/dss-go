// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/IgnoreDataLoader.java (DSS 6.5.RC1).
package http

// IgnoreDataLoader allows avoiding downloading resources: every operation is
// a no-op returning nil.
type IgnoreDataLoader struct{}

// NewIgnoreDataLoader is the default constructor.
func NewIgnoreDataLoader() *IgnoreDataLoader {
	return &IgnoreDataLoader{}
}

// Get is a no-op; slf4j debug logging of the ignored URL is dropped
// (non-load-bearing per PORTING.md).
func (l *IgnoreDataLoader) Get(url string) []byte {
	return nil
}

// GetFromURLs is a no-op.
func (l *IgnoreDataLoader) GetFromURLs(urlStrings []string) (*DataAndURL, error) {
	return nil, nil
}

// Post is a no-op.
func (l *IgnoreDataLoader) Post(url string, content []byte) []byte {
	return nil
}

// SetContentType panics: content type change is not supported by this
// implementation. Ports the Java UnsupportedOperationException.
func (l *IgnoreDataLoader) SetContentType(contentType string) {
	panic("Content type change is not supported by this implementation!")
}

var _ DataLoader = (*IgnoreDataLoader)(nil)
