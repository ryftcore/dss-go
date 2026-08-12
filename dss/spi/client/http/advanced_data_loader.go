// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/AdvancedDataLoader.java (DSS 6.5.RC1).
package http

import modelhttp "github.com/utain/esig/dss/model/http"

// AdvancedDataLoader is used to perform a remote request (HTTP, HTTPS, etc.)
// and retrieve a modelhttp.ResponseEnvelope object, containing contextual
// information delivered from the response (response body, headers, session
// data, etc.).
type AdvancedDataLoader interface {
	DataLoader

	// RequestGet executes a GET request and returns a ResponseEnvelope
	// object. This method includes the response message body, context and
	// metadata within the response object.
	RequestGet(url string) *modelhttp.ResponseEnvelope

	// RequestGetDetails executes a GET request and returns a
	// ResponseEnvelope object. Allows configuration of whether the response
	// context (HTTP headers, TLS/SSL certificates, etc.) is to be included
	// within the response object. The response body is always included.
	RequestGetDetails(url string, includeResponseDetails bool) *modelhttp.ResponseEnvelope

	// RequestGetFull executes a GET request and returns a ResponseEnvelope
	// object. Allows configuration of whether the response context and/or
	// the response message body are to be included within the response
	// object. Data which is not included is not read during processing,
	// for time-memory efficiency.
	RequestGetFull(url string, includeResponseDetails, includeResponseBody bool) *modelhttp.ResponseEnvelope

	// RequestPost executes a POST request and returns a ResponseEnvelope
	// object. This method includes the response message body, context and
	// metadata within the response object.
	RequestPost(url string, content []byte) *modelhttp.ResponseEnvelope

	// RequestPostDetails executes a POST request and returns a
	// ResponseEnvelope object. Allows configuration of whether the response
	// context is to be included within the response object. The response
	// body is always included.
	RequestPostDetails(url string, content []byte, includeResponseDetails bool) *modelhttp.ResponseEnvelope

	// RequestPostFull executes a POST request and returns a ResponseEnvelope
	// object. Allows configuration of whether the response context and/or
	// the response message body are to be included within the response
	// object.
	RequestPostFull(url string, content []byte, includeResponseDetails, includeResponseBody bool) *modelhttp.ResponseEnvelope
}
