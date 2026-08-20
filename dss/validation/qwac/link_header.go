// Ported from the nested class
// dss-validation/src/main/java/eu/europa/esig/dss/validation/qwac/LinkHeaderParser.java$LinkHeader
// (DSS 6.5.RC1).
package qwac

// LinkHeader represents a parsed value of the "Link" HTTP response header.
type LinkHeader struct {
	// url is the "Link" URL.
	url string

	// attributes is a map of the "Link" header attributes. A value of nil (as opposed to an
	// absent key) marks an RFC 8288 flag attribute given without a "=value" part - matching
	// Java's HashMap allowing a null value while distinguishing it from a missing key.
	attributes map[string]*string
}

// URL gets the "Link" header value URL. Port of getUrl().
func (h *LinkHeader) URL() string {
	return h.url
}

// SetURL sets the "Link" header value URL. Port of setUrl(String).
func (h *LinkHeader) SetURL(url string) {
	h.url = url
}

// Attributes gets a map of "Link" header attributes. Port of getAttributes().
func (h *LinkHeader) Attributes() map[string]*string {
	return h.attributes
}

// SetAttributes sets a map of "Link" header attributes. Port of setAttributes(Map).
func (h *LinkHeader) SetAttributes(attributes map[string]*string) {
	h.attributes = attributes
}
