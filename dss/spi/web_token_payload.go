// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/WebTokenPayload.java (DSS 6.5.RC1).
package spi

import "time"

// WebTokenPayload represents a payload of a web token (e.g. RFC 7519 token or RFC 8392 CWT).
//
// Java's Date-typed claims (ExpirationTime, NotBefore, IssuedAt) become time.Time; an absent
// claim is the zero time.Time, matching the convention used elsewhere in this port (e.g.
// model.CertificateToken.IsValidOn).
type WebTokenPayload interface {
	// Issuer gets the value of the Issuer claim identifying the principal that issued the
	// token. Port of getIssuer().
	Issuer() string

	// Subject gets the value of the Subject claim identifying the principal that is the
	// subject of the token. Port of getSubject().
	Subject() string

	// Audience gets the value of the Audience claim identifying the recipients that the token
	// is intended for. Port of getAudience().
	Audience() string

	// ExpirationTime gets the value of the Expiration Time claim identifying the expiration
	// time on or after which the token MUST NOT be accepted for processing. Port of
	// getExpirationTime().
	ExpirationTime() time.Time

	// NotBefore gets the value of the Not Before claim identifying the time before which the
	// token MUST NOT be accepted for processing. Port of getNotBefore().
	NotBefore() time.Time

	// IssuedAt gets the value of the Issued At claim identifying the time before which the
	// token MUST NOT be accepted for processing. Port of getIssuedAt().
	IssuedAt() time.Time

	// TokenID gets the value of the token ID claim representing a unique identifier for the
	// token. Port of getTokenId().
	TokenID() string
}
