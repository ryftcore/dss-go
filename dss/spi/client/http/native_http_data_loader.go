// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/NativeHTTPDataLoader.java (DSS 6.5.RC1).
package http

import (
	"fmt"

	modelhttp "github.com/ryftcore/dss-go/dss/model/http"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// nativeHTTPDataLoaderHTTPMethod are the available HTTP methods. Ports the
// package-private nested enum HttpMethod.
type nativeHTTPDataLoaderHTTPMethod int

const (
	nativeHTTPDataLoaderHTTPMethodGET nativeHTTPDataLoaderHTTPMethod = iota
	nativeHTTPDataLoaderHTTPMethodPOST
)

// NativeHTTPDataLoader is an implementation of a native Go DataLoader using
// net/http.
type NativeHTTPDataLoader struct {
	// maxInputSize is the max InputStream size.
	maxInputSize int

	// connectTimeout is the timeout on the connection establishment with a
	// remote resource, in milliseconds.
	connectTimeout int

	// readTimeout is the timeout on the response reading from a remote
	// resource, in milliseconds.
	readTimeout int
}

// NewNativeHTTPDataLoader is the default constructor instantiating the
// object with zero values.
func NewNativeHTTPDataLoader() *NativeHTTPDataLoader {
	return &NativeHTTPDataLoader{}
}

// SetContentType panics: content type change is not supported by this
// implementation.
func (l *NativeHTTPDataLoader) SetContentType(contentType string) {
	panic("Content type change is not supported by this implementation!")
}

// MaxInputSize gets the maximum InputStream size.
func (l *NativeHTTPDataLoader) MaxInputSize() int {
	return l.maxInputSize
}

// SetMaxInputSize sets the maximum InputStream size.
func (l *NativeHTTPDataLoader) SetMaxInputSize(maxInputSize int) {
	l.maxInputSize = maxInputSize
}

// ConnectTimeout gets the timeout value on connection establishment with a
// remote resource.
func (l *NativeHTTPDataLoader) ConnectTimeout() int {
	return l.connectTimeout
}

// SetConnectTimeout sets the timeout to establish a connection with a remote
// resource (in milliseconds). Zero (0) value is used for no timeout.
// Default: 0 (no timeout).
func (l *NativeHTTPDataLoader) SetConnectTimeout(connectTimeout int) {
	l.connectTimeout = connectTimeout
}

// ReadTimeout gets the timeout value on response reading from a remote
// resource.
func (l *NativeHTTPDataLoader) ReadTimeout() int {
	return l.readTimeout
}

// SetReadTimeout sets the timeout to read a response from a remote resource
// (in milliseconds). Zero (0) value is used for no timeout. Default: 0 (no
// timeout).
func (l *NativeHTTPDataLoader) SetReadTimeout(readTimeout int) {
	l.readTimeout = readTimeout
}

// nativeHTTPDataLoaderCreateCall creates a call to be executed by
// NativeHTTPDataLoader. Ports the protected
// createNativeHTTPDataLoaderCall method.
func (l *NativeHTTPDataLoader) nativeHTTPDataLoaderCreateCall(url string, method nativeHTTPDataLoaderHTTPMethod, content []byte,
	refresh, includeResponseDetails, includeResponseBody bool) *NativeHTTPDataLoaderCall {
	call := NewNativeHTTPDataLoaderCallWithContent(url, content)
	call.SetUseCaches(!refresh)
	call.SetMaxInputSize(l.maxInputSize)
	call.SetConnectTimeout(l.connectTimeout)
	call.SetReadTimeout(l.readTimeout)
	call.SetIncludeResponseDetails(includeResponseDetails)
	call.SetIncludeResponseBody(includeResponseBody)
	return call
}

// nativeHTTPDataLoaderRequest executes the request. Ports the protected
// request method; a Java RuntimeException from task.call() (other than
// DSSExternalResourceException, which is re-thrown as-is) is wrapped into a
// DSSExternalResourceException — Go's Call() already returns
// *exception.DSSExternalResourceException on failure, so this simply
// forwards the error.
func (l *NativeHTTPDataLoader) nativeHTTPDataLoaderRequest(url string, method nativeHTTPDataLoaderHTTPMethod, content []byte,
	refresh, includeResponseDetails, includeResponseBody bool) (*modelhttp.ResponseEnvelope, error) {
	call := l.nativeHTTPDataLoaderCreateCall(url, method, content, refresh, includeResponseDetails, includeResponseBody)
	return call.Call()
}

// GetFromURLs iterates urlStrings, stopping at the first URL that returns
// data. Returns a *exception.DSSExternalResourceException error when no URL
// yields data. Per-URL failures (Get panicking) are recovered and logged as
// dropped (Java logs them at WARN via slf4j; the log call itself is dropped
// as non-load-bearing per PORTING.md, but the catch-and-continue behavior is
// load-bearing and is preserved via recover).
func (l *NativeHTTPDataLoader) GetFromURLs(urlStrings []string) (*DataAndURL, error) {
	for _, urlString := range urlStrings {
		if data := l.tryGet(urlString); data != nil {
			return NewDataAndURL(urlString, data), nil
		}
	}
	return nil, exception.NewDSSExternalResourceException(
		fmt.Sprintf("No data have been obtained from urls : %v", urlStrings))
}

// tryGet calls Get, recovering from the panic Get raises on failure and
// returning nil in that case. Ports the try/catch(Exception) guarding
// get(urlString) inside get(List<String>).
func (l *NativeHTTPDataLoader) tryGet(url string) (data []byte) {
	defer func() {
		if recover() != nil {
			data = nil
		}
	}()
	return l.Get(url)
}

// Get executes a GET request. Panics with the underlying
// *exception.DSSExternalResourceException on failure — the DataLoader
// interface's Get has no error return, mirroring Java's unchecked
// DSSExternalResourceException propagating out of get(String).
func (l *NativeHTTPDataLoader) Get(url string) []byte {
	data, err := l.GetWithRefresh(url, false)
	if err != nil {
		panic(err)
	}
	return data
}

// GetWithRefresh executes a GET request to the provided URL, with a forced
// cache refresh when defined. Ports the public byte[] get(String, boolean).
func (l *NativeHTTPDataLoader) GetWithRefresh(url string, refresh bool) ([]byte, error) {
	resp, err := l.nativeHTTPDataLoaderRequest(url, nativeHTTPDataLoaderHTTPMethodGET, nil, refresh, false, true)
	if err != nil {
		return nil, err
	}
	return resp.ResponseBody(), nil
}

// Post executes a POST request. Panics on failure; see Get.
func (l *NativeHTTPDataLoader) Post(url string, content []byte) []byte {
	resp, err := l.nativeHTTPDataLoaderRequest(url, nativeHTTPDataLoaderHTTPMethodPOST, content, false, false, true)
	if err != nil {
		panic(err)
	}
	return resp.ResponseBody()
}

// RequestGet executes a GET request and returns a ResponseEnvelope object.
// Panics on failure; see Get.
func (l *NativeHTTPDataLoader) RequestGet(url string) *modelhttp.ResponseEnvelope {
	resp, err := l.nativeHTTPDataLoaderRequest(url, nativeHTTPDataLoaderHTTPMethodGET, nil, false, true, true)
	if err != nil {
		panic(err)
	}
	return resp
}

// RequestGetDetails executes a GET request and returns a ResponseEnvelope
// object. Panics on failure; see Get.
func (l *NativeHTTPDataLoader) RequestGetDetails(url string, includeResponseDetails bool) *modelhttp.ResponseEnvelope {
	resp, err := l.nativeHTTPDataLoaderRequest(url, nativeHTTPDataLoaderHTTPMethodGET, nil, false, includeResponseDetails, true)
	if err != nil {
		panic(err)
	}
	return resp
}

// RequestGetFull executes a GET request and returns a ResponseEnvelope
// object. Panics on failure; see Get.
func (l *NativeHTTPDataLoader) RequestGetFull(url string, includeResponseDetails, includeResponseBody bool) *modelhttp.ResponseEnvelope {
	resp, err := l.nativeHTTPDataLoaderRequest(url, nativeHTTPDataLoaderHTTPMethodGET, nil, false, includeResponseDetails, includeResponseBody)
	if err != nil {
		panic(err)
	}
	return resp
}

// RequestPost executes a POST request and returns a ResponseEnvelope object.
// Panics on failure; see Get.
func (l *NativeHTTPDataLoader) RequestPost(url string, content []byte) *modelhttp.ResponseEnvelope {
	resp, err := l.nativeHTTPDataLoaderRequest(url, nativeHTTPDataLoaderHTTPMethodPOST, content, false, true, true)
	if err != nil {
		panic(err)
	}
	return resp
}

// RequestPostDetails executes a POST request and returns a ResponseEnvelope
// object. Panics on failure; see Get.
func (l *NativeHTTPDataLoader) RequestPostDetails(url string, content []byte, includeResponseDetails bool) *modelhttp.ResponseEnvelope {
	resp, err := l.nativeHTTPDataLoaderRequest(url, nativeHTTPDataLoaderHTTPMethodPOST, content, false, includeResponseDetails, true)
	if err != nil {
		panic(err)
	}
	return resp
}

// RequestPostFull executes a POST request and returns a ResponseEnvelope
// object. Panics on failure; see Get.
func (l *NativeHTTPDataLoader) RequestPostFull(url string, content []byte, includeResponseDetails, includeResponseBody bool) *modelhttp.ResponseEnvelope {
	resp, err := l.nativeHTTPDataLoaderRequest(url, nativeHTTPDataLoaderHTTPMethodPOST, content, false, includeResponseDetails, includeResponseBody)
	if err != nil {
		panic(err)
	}
	return resp
}

var (
	_ DataLoader         = (*NativeHTTPDataLoader)(nil)
	_ AdvancedDataLoader = (*NativeHTTPDataLoader)(nil)
)
