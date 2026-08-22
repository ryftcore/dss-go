// Ported from specs-jades/src/main/java/eu/europa/esig/jades/AbstractJAdESUtils.java (DSS 6.5.RC1).
//
// Java's AbstractJAdESUtils extends specs-jws' eu.europa.esig.jws.AbstractJWSUtils, itself
// extending dss-json-common's eu.europa.esig.json.JSONSchemaAbstractUtils - neither of which
// belongs to a module this port ports separately (they live outside dss-jades/specs-jades).
// Since specs is dependency-closed, the three-level Java hierarchy collapses into this single Go
// type: the two base classes' behaviour (getSchemaDefinitions's JSON-draft-07 + RFC7515/7517
// entries, the validateAgainstSchema/getSchema/getValidator caching machinery) is folded in here
// rather than split into a separate package.
package specs

import (
	"io/fs"
)

// jadesSchemaRoot is the schema/ subtree of the embedded resources, so that URIs used as
// definitions-map keys match the on-disk relative paths one-to-one (e.g. "rfcs/rfc7519.json").
var jadesSchemaRoot = func() fs.FS {
	sub, err := fs.Sub(schemaFS, "schema")
	if err != nil {
		panic(err)
	}
	return sub
}()

const (
	// jadesRFC7519SchemaURI is the RFC 7519 schema name URI as registered under the "rfcs/"
	// subdirectory. Port of RFC7519_SCHEMA_URI (combined with RFC_SUBDIRECTORY).
	jadesRFC7519SchemaURI = "rfcs/rfc7519.json"

	// jadesRFC7797SchemaURI is the RFC 7797 schema name URI, likewise under "rfcs/". Port of
	// RFC7797_SCHEMA_URI (combined with RFC_SUBDIRECTORY).
	jadesRFC7797SchemaURI = "rfcs/rfc7797.json"

	// jadesSchemaDefinitionsURI is the JAdES schema name URI. Port of
	// JAdES_SCHEMA_DEFINITIONS_URI.
	jadesSchemaDefinitionsURI = "19182-jsonSchema.json"

	// jwsRFC7515SchemaURI and jwsRFC7517SchemaURI are the unprefixed RFC 7515/7517 schema URIs,
	// as registered by AbstractJWSUtils#getRFCDefinitions before the JAdES "rfcs/" prefix is
	// applied.
	jwsRFC7515SchemaURI = "rfc7515.json"
	jwsRFC7517SchemaURI = "rfc7517.json"
)

// abstractJAdESUtils is the abstract class for JAdES signature validation against JSON schemas.
// Concrete validators (JAdESProtectedHeaderUtils, JAdESUnprotectedHeaderUtils) embed it by value
// and call SchemaDefinitions() to obtain the shared base entries before adding their own.
type abstractJAdESUtils struct{}

// rfcDefinitions gets a map of RFC definitions, all placed under the "rfcs/" subdirectory (so
// that a JAdES schema's own reference to, say, "19182-jsonSchema.json" cannot collide with a
// same-named RFC definitions schema). Port of #getRFCDefinitions.
func (abstractJAdESUtils) rfcDefinitions() map[string]string {
	definitions := map[string]string{
		"rfcs/" + jwsRFC7515SchemaURI: jwsRFC7515SchemaURI,
		"rfcs/" + jwsRFC7517SchemaURI: jwsRFC7517SchemaURI,
	}
	definitions[jadesRFC7519SchemaURI] = jadesRFC7519SchemaURI
	definitions[jadesRFC7797SchemaURI] = jadesRFC7797SchemaURI
	return definitions
}

// schemaDefinitions gets a list of schema definitions: a map between schema URIs and their
// on-disk (embedded) relative path, folding in JSONSchemaAbstractUtils#getJSONSchemaDefinitions'
// draft-07 declaration (a no-op here - see the file header on the JSON-draft-07 meta-schema
// deviation, schema_engine.go never dereferences it) and AbstractJWSUtils#getRFCDefinitions via
// the "rfcs/"-prefixed rfcDefinitions above. Port of #getSchemaDefinitions.
func (u abstractJAdESUtils) schemaDefinitions() map[string]string {
	definitions := map[string]string{}
	for k, v := range u.rfcDefinitions() {
		definitions[k] = v
	}
	definitions[jadesSchemaDefinitionsURI] = jadesSchemaDefinitionsURI
	return definitions
}

// validateAgainstSchema validates a JSON document against the schema registered under schemaURI
// within definitions, returning a list of error messages - empty when the JSON is valid. Port of
// JSONSchemaAbstractUtils#validateAgainstSchema(String), folded in here per the file header.
func validateAgainstSchema(schemaURI string, definitions map[string]string, jsonString string) []string {
	return validateJSONAgainstSchema(jadesSchemaRoot, definitions, schemaURI, jsonString)
}
