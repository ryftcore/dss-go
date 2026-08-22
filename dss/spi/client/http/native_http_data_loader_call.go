// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/NativeHTTPDataLoaderCall.java (DSS 6.5.RC1).
//
// Deviation: Java's URLConnection separates a connect timeout from a
// read timeout on the same connection. net/http does not expose that split
// cleanly; this port uses connectTimeout as the dial timeout
// (http.Transport.DialContext / net.Dialer.Timeout) and readTimeout as the
// overall per-request timeout (http.Client.Timeout, which bounds connection,
// any redirects, and reading the response body). When both are zero the
// request has no deadline, matching upstream's "0 = no timeout" default.
package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	modelhttp "github.com/ryftcore/dss-go/dss/model/http"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// nativeHTTPErrorMessage is the default error message format.
const nativeHTTPErrorMessage = "An error occurred while reading from url '%s' : %s"

// NativeHTTPDataLoaderCall is the call of the native Go net/http DataLoader.
// It returns a *modelhttp.ResponseEnvelope object. Set includeResponseDetails
// when extended HTTP content data (headers, TLS session parameters, etc.) is
// needed; if not set, only the response message body is returned.
type NativeHTTPDataLoaderCall struct {
	// url is the request URL.
	url string

	// content is the body of the request, when performing a POST request.
	content []byte

	// useCaches defines whether the cache is used. Not meaningful for
	// net/http (no shared URL cache exists like java.net.URLConnection's);
	// kept for API parity and currently unused, matching a stdlib client
	// with no caching layer.
	useCaches bool

	// maxInputSize is the max input size.
	maxInputSize int

	// connectTimeout is the timeout on opening a connection with the
	// remote resource, in milliseconds.
	connectTimeout int

	// readTimeout is the timeout on reading a response from the remote
	// resource, in milliseconds.
	readTimeout int

	// includeResponseBody defines whether the response message body is to
	// be included in the final response object.
	includeResponseBody bool

	// includeResponseDetails defines whether the extended HTTP content
	// (headers, TLS/SSL certificates, etc.) should be included in the
	// result.
	includeResponseDetails bool
}

// NewNativeHTTPDataLoaderCall is the constructor for a GET call
// instantiation.
func NewNativeHTTPDataLoaderCall(url string) *NativeHTTPDataLoaderCall {
	return NewNativeHTTPDataLoaderCallWithContent(url, nil)
}

// NewNativeHTTPDataLoaderCallWithContent is the constructor for a POST call
// instantiation.
func NewNativeHTTPDataLoaderCallWithContent(url string, content []byte) *NativeHTTPDataLoaderCall {
	return &NativeHTTPDataLoaderCall{url: url, content: content, includeResponseBody: true}
}

// SetUseCaches sets whether the caches shall be used.
func (c *NativeHTTPDataLoaderCall) SetUseCaches(useCaches bool) {
	c.useCaches = useCaches
}

// SetMaxInputSize sets the limit size of the HTTP response body. If the HTTP
// response exceeds the given value, an error will be returned. For 0 or
// negative value, no limit will be exposed.
func (c *NativeHTTPDataLoaderCall) SetMaxInputSize(maxInputSize int) {
	c.maxInputSize = maxInputSize
}

// SetConnectTimeout sets a connection timeout, in milliseconds.
func (c *NativeHTTPDataLoaderCall) SetConnectTimeout(connectTimeout int) {
	c.connectTimeout = connectTimeout
}

// SetReadTimeout sets a read timeout, in milliseconds.
func (c *NativeHTTPDataLoaderCall) SetReadTimeout(readTimeout int) {
	c.readTimeout = readTimeout
}

// SetIncludeResponseBody sets whether the response message body is to be
// included in the final response object. Default: true.
func (c *NativeHTTPDataLoaderCall) SetIncludeResponseBody(includeResponseBody bool) {
	c.includeResponseBody = includeResponseBody
}

// SetIncludeResponseDetails sets whether additional HTTP response content is
// to be included in the final object (e.g. TLS/SSL certificates, headers,
// etc.). Default: false.
func (c *NativeHTTPDataLoaderCall) SetIncludeResponseDetails(includeResponseDetails bool) {
	c.includeResponseDetails = includeResponseDetails
}

// Call executes the HTTP request. Ports Callable<ResponseEnvelope>#call.
func (c *NativeHTTPDataLoaderCall) Call() (*modelhttp.ResponseEnvelope, error) {
	client := c.createClient()

	var bodyReader io.Reader
	method := http.MethodGet
	if c.content != nil {
		method = http.MethodPost
		bodyReader = bytes.NewReader(c.content)
	}

	req, err := http.NewRequest(method, c.url, bodyReader)
	if err != nil {
		return nil, exception.NewDSSExternalResourceExceptionMessageCause(
			fmt.Sprintf(nativeHTTPErrorMessage, c.url, err.Error()), err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, exception.NewDSSExternalResourceExceptionMessageCause(
			fmt.Sprintf(nativeHTTPErrorMessage, c.url, err.Error()), err)
	}
	defer utils.CloseQuietly(resp.Body)

	result := modelhttp.NewResponseEnvelope()

	if c.includeResponseBody {
		var reader io.Reader = resp.Body
		if c.maxInputSize > 0 {
			reader = NewMaxSizeInputStream(resp.Body, c.maxInputSize, c.url)
		}
		body, readErr := utils.ToByteArray(reader)
		if readErr != nil {
			return nil, exception.NewDSSExternalResourceExceptionMessageCause(
				fmt.Sprintf(nativeHTTPErrorMessage, c.url, readErr.Error()), readErr)
		}
		result.SetResponseBody(body)
	}

	if c.includeResponseDetails {
		result.SetHeaders(map[string][]string(resp.Header))

		// If it's HTTPS, retrieve the TLS session's peer certificates.
		if resp.TLS != nil {
			result.SetTLSCertificates(resp.TLS.PeerCertificates)
		}
	}

	return result, nil
}

// createClient builds the http.Client used for the request. Ports
// createConnection: the connect timeout bounds dialing, the read timeout
// bounds the overall request (see the file-level deviation note).
func (c *NativeHTTPDataLoaderCall) createClient() *http.Client {
	dialer := &net.Dialer{}
	if c.connectTimeout > 0 {
		dialer.Timeout = time.Duration(c.connectTimeout) * time.Millisecond
	}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		},
	}
	client := &http.Client{Transport: transport}
	if c.readTimeout > 0 {
		client.Timeout = time.Duration(c.readTimeout) * time.Millisecond
	}
	return client
}
