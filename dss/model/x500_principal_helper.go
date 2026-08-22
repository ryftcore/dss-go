// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/X500PrincipalHelper.java (DSS 6.5.RC1).
package model

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// X500PrincipalHelper extracts the String representations of an X500Principal
// distinguishing name.
type X500PrincipalHelper struct {
	// principal is the X500Principal to be processed.
	principal *X500Principal
}

// NewX500PrincipalHelper wraps the given principal.
//
// Panics with the Java message when the principal is missing (Objects.requireNonNull).
func NewX500PrincipalHelper(principal *X500Principal) *X500PrincipalHelper {
	if principal == nil {
		panic("X500Principal cannot be null!")
	}
	return &X500PrincipalHelper{principal: principal}
}

// Principal returns the wrapped X500Principal. Port of getPrincipal().
func (h *X500PrincipalHelper) Principal() *X500Principal {
	return h.principal
}

// Canonical returns the canonical name. Port of getCanonical(), i.e.
// principal.getName(X500Principal.CANONICAL).
func (h *X500PrincipalHelper) Canonical() string {
	return h.principal.Canonical()
}

// RFC2253 returns the RFC 2253 standard name. Port of getRFC2253(), i.e.
// principal.getName(X500Principal.RFC2253).
func (h *X500PrincipalHelper) RFC2253() string {
	return h.principal.RFC2253Name()
}

// PrettyPrintRFC2253 returns the pretty-printed RFC 2253 standard name, i.e. the RFC 2253
// form with the X520Attributes OID descriptions substituted for the attribute type
// keywords. Port of getPrettyPrintRFC2253().
//
// Java throws IllegalArgumentException when an OID in the name maps to an improperly
// specified keyword; that becomes the returned error. Every X520Attributes description is a
// valid keyword, so the error cannot trigger for the built-in map.
func (h *X500PrincipalHelper) PrettyPrintRFC2253() (string, error) {
	return h.principal.RFC2253NameWithOIDMap(enumerations.GetOidDescriptions())
}

// Encoded returns the encoded X500Principal binaries. Port of getEncoded().
func (h *X500PrincipalHelper) Encoded() []byte {
	return h.principal.Encoded()
}

// Equals reports whether both helpers wrap a principal with the same DER encoding.
// Port of equals(Object).
//
// NOTE: this deliberately differs from X500Principal.Equals, which compares canonical names;
// the helper compares the raw encodings, exactly as upstream does.
func (h *X500PrincipalHelper) Equals(other *X500PrincipalHelper) bool {
	if h == other {
		return true
	}
	if other == nil {
		return false
	}
	return bytes.Equal(h.principal.Encoded(), other.principal.Encoded())
}
