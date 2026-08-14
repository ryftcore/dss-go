// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/jwt/JWTClaimNames.java
// (DSS 6.5.RC1).
//
// The Java package eu.europa.esig.dss.jades.jwt is folded into this one Go package, per the
// phase-6 package layout; the JWTClaimNames prefix keeps the class-scoped names unambiguous in
// Go's single package namespace.
package jades

// The registered claim names of RFC 7519 "JSON Web Token (JWT)" section 4.1. Values are verbatim.
const (
	// JWTClaimNamesIss is the "iss" (Issuer) claim: it identifies the principal that issued the
	// JWT. The processing of this claim is generally application specific. The value is a
	// case-sensitive string containing a StringOrURI value. Use of this claim is OPTIONAL.
	// Port of ISS.
	JWTClaimNamesIss = "iss"

	// JWTClaimNamesSub is the "sub" (Subject) claim: it identifies the principal that is the
	// subject of the JWT. The subject value MUST either be scoped to be locally unique in the
	// context of the issuer or be globally unique. The value is a case-sensitive string
	// containing a StringOrURI value. Use of this claim is OPTIONAL. Port of SUB.
	JWTClaimNamesSub = "sub"

	// JWTClaimNamesAud is the "aud" (Audience) claim: it identifies the recipients that the JWT
	// is intended for. Each principal intended to process the JWT MUST identify itself with a
	// value in the audience claim; if it does not, the JWT MUST be rejected. In the general case
	// the value is an array of case-sensitive strings, each containing a StringOrURI value, and
	// in the single-audience case it MAY be one such string. Use of this claim is OPTIONAL.
	// Port of AUD.
	JWTClaimNamesAud = "aud"

	// JWTClaimNamesExp is the "exp" (Expiration Time) claim: it identifies the expiration time
	// on or after which the JWT MUST NOT be accepted for processing. Its value MUST be a number
	// containing a NumericDate value. Use of this claim is OPTIONAL. Port of EXP.
	JWTClaimNamesExp = "exp"

	// JWTClaimNamesNbf is the "nbf" (Not Before) claim: it identifies the time before which the
	// JWT MUST NOT be accepted for processing. Its value MUST be a number containing a
	// NumericDate value. Use of this claim is OPTIONAL. Port of NBF.
	JWTClaimNamesNbf = "nbf"

	// JWTClaimNamesIat is the "iat" (Issued At) claim: it identifies the time at which the JWT
	// was issued, and can be used to determine its age. Its value MUST be a number containing a
	// NumericDate value. Use of this claim is OPTIONAL. Port of IAT.
	JWTClaimNamesIat = "iat"

	// JWTClaimNamesJti is the "jti" (JWT ID) claim: it provides a unique identifier for the JWT,
	// assigned so that the same value will not be accidentally assigned to a different data
	// object. It can be used to prevent the JWT from being replayed. The value is a
	// case-sensitive string. Use of this claim is OPTIONAL. Port of JTI.
	JWTClaimNamesJti = "jti"
)
