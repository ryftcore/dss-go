// Ported from specs-jades/src/main/java/eu/europa/esig/jades/JAdESUtils.java (DSS 6.5.RC1).
//
// Contains utils for JAdES signature validation against JSON schemas. Upstream's getSchemaURI
// and getSchemaDefinitions both delegate entirely to specs-jws' eu.europa.esig.jws.JWSUtils
// (rather than to AbstractJAdESUtils' own JAdES-specific schemaDefinitions()): this validates a
// JWS document's generic envelope shape (compact/JSON/flattened serialization), independent of
// any JAdES-specific header constraint - those are checked separately through
// JAdESProtectedHeaderUtils / JAdESUnprotectedHeaderUtils. Since specs-jws is not a module this
// port ports separately, and specs is dependency-closed, JWSUtils' definitions are inlined here
// rather than routed through a separate package - the same collapsing documented in
// abstract_jades_utils.go.
package specs

import "sync"

// jwsSchemaURI is the main JWS signature wrapper schema URI. Port of JWSUtils#JWS_SCHEMA_URI.
const jwsSchemaURI = "rfc7515-jws.json"

// JAdESUtils contains utils for JAdES signature validation against JSON schemas.
type JAdESUtils struct{}

var (
	jadesUtilsSingleton     *JAdESUtils
	jadesUtilsSingletonOnce sync.Once
)

// JAdESUtilsInstance returns the singleton instance of JAdESUtils. Port of the private
// constructor + #getInstance lazy singleton.
func JAdESUtilsInstance() *JAdESUtils {
	jadesUtilsSingletonOnce.Do(func() {
		jadesUtilsSingleton = &JAdESUtils{}
	})
	return jadesUtilsSingleton
}

// SchemaURI ports the overridden #getSchemaURI, delegating to JWSUtils#getSchemaURI.
func (u *JAdESUtils) SchemaURI() string {
	return jwsSchemaURI
}

// SchemaDefinitions ports the overridden #getSchemaDefinitions, delegating to
// JWSUtils#getSchemaDefinitions: the JSON-draft-07 declaration (a no-op here, see
// abstract_jades_utils.go) plus the unprefixed RFC 7515/7517 definitions plus this schema's own
// entry.
func (u *JAdESUtils) SchemaDefinitions() map[string]string {
	return map[string]string{
		jwsRFC7515SchemaURI: jwsRFC7515SchemaURI,
		jwsRFC7517SchemaURI: jwsRFC7517SchemaURI,
		jwsSchemaURI:        jwsSchemaURI,
	}
}

// ValidateAgainstSchema validates a JSON against the generic JWS envelope schema (RFC 7515),
// returning a list of error messages - empty when the JSON is valid. Port of
// JSONSchemaAbstractUtils#validateAgainstSchema(String) as inherited here.
func (u *JAdESUtils) ValidateAgainstSchema(jsonString string) []string {
	return validateAgainstSchema(u.SchemaURI(), u.SchemaDefinitions(), jsonString)
}
