// Ported from org.jose4j.jwx.Headers (jose4j 0.9.6).
package jose

// Headers is a JOSE header: an insertion-ordered member map plus the two cached renderings of
// it. Port of org.jose4j.jwx.Headers.
//
// The caching is not an optimisation, it is the contract. setEncodedHeader stores the received
// base64url text and decodes it for the member map, so EncodedHeader afterwards returns the
// ORIGINAL bytes rather than a re-serialization of the parsed members. That is what lets a JAdES
// signature be validated after a round trip through the parser even if its header was written by
// some other library with a different member order or different escaping - re-encoding it would
// invent a signing input that was never signed. Conversely, every Put invalidates both caches,
// so a header that is being built is always rendered from its current members.
//
// The zero value is not usable; construct with NewHeaders.
type Headers struct {
	members      *Object
	header       string
	headerValid  bool
	encoded      string
	encodedValid bool
}

// NewHeaders returns an empty header. Port of the implicit `new Headers()`, whose backing map is
// a LinkedHashMap.
func NewHeaders() *Headers {
	return &Headers{members: NewObject()}
}

// FullHeaderAsJSONString returns the header serialized as JSON, computing it on first use and
// caching it. Port of getFullHeaderAsJsonString().
func (h *Headers) FullHeaderAsJSONString() string {
	if !h.headerValid {
		h.header = JSON(h.members)
		h.headerValid = true
	}
	return h.header
}

// EncodedHeader returns BASE64URL(UTF8(header)), computing it on first use and caching it. Port
// of getEncodedHeader(). This string is the first part of the JWS signing input.
func (h *Headers) EncodedHeader() string {
	if !h.encodedValid {
		h.encoded = Base64URLEncodeUTF8(h.FullHeaderAsJSONString())
		h.encodedValid = true
	}
	return h.encoded
}

// Put sets a header parameter and invalidates both cached renderings. Port of
// setObjectHeaderValue(String, Object) - and of setStringHeaderValue, which merely delegates.
func (h *Headers) Put(name string, value any) {
	h.members.Put(name, value)
	h.headerValid = false
	h.encodedValid = false
}

// Value returns the header parameter under name, or nil. Port of getObjectHeaderValue(String).
func (h *Headers) Value(name string) any {
	if h == nil {
		return nil
	}
	return h.members.Value(name)
}

// StringValue returns the header parameter under name when it is a string, and "" otherwise.
// Port of getStringHeaderValue(String), which is JsonHelp.getString: an unchecked cast, so a
// non-string value yields the Java null that "" stands in for here.
func (h *Headers) StringValue(name string) string {
	s, _ := h.Value(name).(string)
	return s
}

// Members returns the header's member map. Mutating it bypasses the cache invalidation that Put
// performs, so callers that write must go through Put; the accessor exists for the read paths
// (DSSJsonUtils.extractJOSEHeaderMembersSet and the JAdES header getters).
func (h *Headers) Members() *Object {
	if h == nil {
		return nil
	}
	return h.members
}

// SetFullHeaderAsJSONString replaces the header from its JSON text, keeping that exact text as
// the cached rendering. Port of setFullHeaderAsJsonString(String).
func (h *Headers) SetFullHeaderAsJSONString(header string) error {
	members, err := ParseJSON(header)
	if err != nil {
		return err
	}
	h.encodedValid = false
	h.encoded = ""
	h.header = header
	h.headerValid = true
	h.members = members
	return nil
}

// SetEncodedHeader replaces the header from its base64url text, keeping BOTH that text and its
// decoding as the cached renderings. Port of setEncodedHeader(String).
//
// Keeping the encoded form verbatim is what preserves the signing input across a parse; see the
// type comment.
func (h *Headers) SetEncodedHeader(encodedHeader string) error {
	header := Base64URLDecodeToUTF8String(encodedHeader)
	members, err := ParseJSON(header)
	if err != nil {
		return err
	}
	h.encoded = encodedHeader
	h.encodedValid = true
	h.header = header
	h.headerValid = true
	h.members = members
	return nil
}
