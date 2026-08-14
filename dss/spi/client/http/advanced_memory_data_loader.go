// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/AdvancedMemoryDataLoader.java (DSS 6.5.RC1).
package http

import (
	"fmt"
	"strings"

	modelhttp "github.com/utain/esig/dss/model/http"
	"github.com/utain/esig/dss/spi/exception"
)

// AdvancedMemoryDataLoader defines a map between URL and document to load the
// response data from an offline source.
type AdvancedMemoryDataLoader struct {
	// dataMap is the map between URLs and the corresponding response data
	// content.
	dataMap map[string]*modelhttp.ResponseEnvelope
}

// NewAdvancedMemoryDataLoader is the default constructor. Copies dataMap,
// mirroring the Java constructor's putAll into a fresh backing map.
func NewAdvancedMemoryDataLoader(dataMap map[string]*modelhttp.ResponseEnvelope) *AdvancedMemoryDataLoader {
	m := make(map[string]*modelhttp.ResponseEnvelope, len(dataMap))
	for k, v := range dataMap {
		m[k] = v
	}
	return &AdvancedMemoryDataLoader{dataMap: m}
}

// Get returns the response body stored under url, via RequestGet.
func (l *AdvancedMemoryDataLoader) Get(url string) []byte {
	return l.RequestGet(url).ResponseBody()
}

// GetFromURLs returns the data and URL for the first url with stored data.
// Returns a *exception.DSSExternalResourceException error when none of the
// URLs have data.
func (l *AdvancedMemoryDataLoader) GetFromURLs(urlStrings []string) (*DataAndURL, error) {
	for _, url := range urlStrings {
		if data := l.Get(url); data != nil {
			return NewDataAndURL(url, data), nil
		}
	}
	return nil, exception.NewDSSExternalResourceException(
		fmt.Sprintf("A content for URLs [%s] does not exist!", strings.Join(urlStrings, ", ")))
}

// Post returns the response body stored under url, via RequestPost.
func (l *AdvancedMemoryDataLoader) Post(url string, content []byte) []byte {
	return l.RequestPost(url, content).ResponseBody()
}

// SetContentType panics: content type change is not supported by this
// implementation.
func (l *AdvancedMemoryDataLoader) SetContentType(contentType string) {
	panic("Content type change is not supported by this implementation!")
}

// RequestGet performs a GET, including headers/TLS certificates and body.
func (l *AdvancedMemoryDataLoader) RequestGet(url string) *modelhttp.ResponseEnvelope {
	return l.RequestGetDetails(url, true)
}

// RequestGetDetails performs a GET, optionally including headers/TLS
// certificates; the body is always included.
func (l *AdvancedMemoryDataLoader) RequestGetDetails(url string, includeResponseDetails bool) *modelhttp.ResponseEnvelope {
	return l.RequestGetFull(url, includeResponseDetails, true)
}

// RequestGetFull performs a GET, optionally including headers/TLS
// certificates and/or the response body.
func (l *AdvancedMemoryDataLoader) RequestGetFull(url string, includeResponseDetails, includeResponseBody bool) *modelhttp.ResponseEnvelope {
	storedValue, ok := l.dataMap[url]
	if !ok {
		return modelhttp.NewResponseEnvelope()
	}

	response := modelhttp.NewResponseEnvelope()
	if includeResponseDetails {
		response.SetHeaders(storedValue.Headers())
		response.SetTLSCertificates(storedValue.TLSCertificates())
	}
	if includeResponseBody {
		response.SetResponseBody(storedValue.ResponseBody())
	}
	return response
}

// RequestPost performs a POST, including headers/TLS certificates and body.
// Delegates to RequestGetFull, mirroring upstream's identical behaviour for
// GET and POST against this offline data source.
func (l *AdvancedMemoryDataLoader) RequestPost(url string, content []byte) *modelhttp.ResponseEnvelope {
	return l.RequestPostDetails(url, content, true)
}

// RequestPostDetails performs a POST, optionally including headers/TLS
// certificates; the body is always included.
func (l *AdvancedMemoryDataLoader) RequestPostDetails(url string, content []byte, includeResponseDetails bool) *modelhttp.ResponseEnvelope {
	return l.RequestPostFull(url, content, includeResponseDetails, true)
}

// RequestPostFull performs a POST, optionally including headers/TLS
// certificates and/or the response body.
func (l *AdvancedMemoryDataLoader) RequestPostFull(url string, content []byte, includeResponseDetails, includeResponseBody bool) *modelhttp.ResponseEnvelope {
	return l.RequestGetFull(url, includeResponseDetails, includeResponseBody)
}

var _ AdvancedDataLoader = (*AdvancedMemoryDataLoader)(nil)
