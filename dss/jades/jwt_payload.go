// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/jwt/JWTPayload.java
// (DSS 6.5.RC1).
//
// The Java package eu.europa.esig.dss.jades.jwt is folded into this one Go package, per the
// phase-6 package layout.
package jades

import (
	"time"

	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/spi"
)

// JWTPayload represents a payload of the RFC 7519 "JSON Web Token (JWT)". Port of the class
// JWTPayload, which implements eu.europa.esig.dss.spi.WebTokenPayload.
//
// The interface's Date-typed claims become time.Time, an absent claim being the zero time.Time -
// the convention spi.WebTokenPayload already documents.
type JWTPayload struct {
	// payload is the map representing the JSON payload. Port of the private final field of the
	// same name.
	payload *jose.Object
}

// Compile-time assertion that the port really satisfies the interface it declares.
var _ spi.WebTokenPayload = (*JWTPayload)(nil)

// NewJWTPayload is the default constructor. Port of JWTPayload(Map).
func NewJWTPayload(payload *jose.Object) *JWTPayload {
	return &JWTPayload{payload: payload}
}

// Issuer gets the value of the 'iss' (Issuer) claim identifying the principal that issued the
// JWT. Port of getIssuer().
func (p *JWTPayload) Issuer() string {
	return p.GetAsString(JWTClaimNamesIss)
}

// Subject gets the value of the 'sub' (Subject) claim identifying the principal that is the
// subject of the JWT. Port of getSubject().
func (p *JWTPayload) Subject() string {
	return p.GetAsString(JWTClaimNamesSub)
}

// Audience gets the value of the 'aud' (Audience) claim identifying the recipients that the JWT
// is intended for. Port of getAudience().
//
// RFC 7519 allows 'aud' to be either a string or an array of strings; upstream reads it as a
// string and answers the empty string for the array form, which this reproduces rather than
// quietly improving on.
func (p *JWTPayload) Audience() string {
	return p.GetAsString(JWTClaimNamesAud)
}

// ExpirationTime gets the value of the 'exp' (Expiration Time) claim identifying the expiration
// time on or after which the JWT MUST NOT be accepted for processing. Port of
// getExpirationTime().
func (p *JWTPayload) ExpirationTime() time.Time {
	return p.GetAsDate(JWTClaimNamesExp)
}

// NotBefore gets the value of the 'nbf' (Not Before) claim identifying the time before which the
// JWT MUST NOT be accepted for processing. Port of getNotBefore().
func (p *JWTPayload) NotBefore() time.Time {
	return p.GetAsDate(JWTClaimNamesNbf)
}

// IssuedAt gets the value of the 'iat' (Issued At) claim identifying the time at which the JWT
// was issued. Port of getIssuedAt().
func (p *JWTPayload) IssuedAt() time.Time {
	return p.GetAsDate(JWTClaimNamesIat)
}

// TokenID gets the value of the 'jti' (JWT ID) claim representing a unique identifier for the
// JWT. Port of getTokenId().
func (p *JWTPayload) TokenID() string {
	return p.GetAsString(JWTClaimNamesJti)
}

// GetAsString gets the claim under name as a String. If the claim is absent or is not a string,
// it returns the empty string. Port of getAsString(String).
func (p *JWTPayload) GetAsString(name string) string {
	return DSSJsonUtilsGetAsString(p.payload, name)
}

// GetAsNumber gets the claim under name as a Number. If the claim is absent or is not a number,
// it returns nil. Port of getAsNumber(String).
func (p *JWTPayload) GetAsNumber(name string) *jose.Number {
	return DSSJsonUtilsGetAsNumber(p.payload, name)
}

// GetAsDate gets the claim under name as a Date. If the claim is absent or is not a NumericDate,
// it returns the zero time.Time. Port of getAsDate(String).
func (p *JWTPayload) GetAsDate(name string) time.Time {
	return DSSJsonUtilsGetAsNumericDate(p.payload, name)
}

// GetAsMap gets the claim under name as a JSON object. If the claim is absent or is not an
// object, it returns an empty one. Port of getAsMap(String).
func (p *JWTPayload) GetAsMap(name string) *jose.Object {
	return DSSJsonUtilsGetAsMap(p.payload, name)
}
