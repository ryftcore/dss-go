// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/MemoryDataLoader.java (DSS 6.5.RC1).
package http

import (
	"fmt"
	"strings"

	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// MemoryDataLoader defines a map between URL and document to load the data
// from an offline source.
type MemoryDataLoader struct {
	// dataMap is the map between URLs and the corresponding binary content.
	dataMap map[string][]byte
}

// NewMemoryDataLoader is the default constructor. Copies dataMap, mirroring
// the Java constructor's putAll into a fresh backing map.
func NewMemoryDataLoader(dataMap map[string][]byte) *MemoryDataLoader {
	m := make(map[string][]byte, len(dataMap))
	for k, v := range dataMap {
		m[k] = v
	}
	return &MemoryDataLoader{dataMap: m}
}

// Get returns the data stored under url, or nil.
func (l *MemoryDataLoader) Get(url string) []byte {
	return l.dataMap[url]
}

// GetFromURLs returns the data and URL for the first url with stored data.
// Returns a *exception.DSSExternalResourceException error when none of the
// URLs have data.
func (l *MemoryDataLoader) GetFromURLs(urlStrings []string) (*DataAndURL, error) {
	for _, url := range urlStrings {
		if data := l.Get(url); data != nil {
			return NewDataAndURL(url, data), nil
		}
	}
	return nil, exception.NewDSSExternalResourceException(
		fmt.Sprintf("A content for URLs [%s] does not exist!", strings.Join(urlStrings, ", ")))
}

// Post returns the data stored under url (the content is ignored, matching
// upstream's delegation to Get).
func (l *MemoryDataLoader) Post(url string, content []byte) []byte {
	return l.Get(url)
}

// SetContentType panics: content type change is not supported by this
// implementation.
func (l *MemoryDataLoader) SetContentType(contentType string) {
	panic("Content type change is not supported by this implementation!")
}

var _ DataLoader = (*MemoryDataLoader)(nil)
