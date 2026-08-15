// Ported from specs-jades/src/main/java/eu/europa/esig/jades/JAdESUnprotectedHeaderUtils.java (DSS 6.5.RC1).
//
// Contains utils for JAdES Unprotected header validation against JSON schemas.
package specs

import "sync"

// jadesUnprotectedHeaderSchemaURI is the unprotected header schema URI for a JAdES signature.
// Port of JAdES_UNPROTECTED_HEADER_SCHEMA_URI.
const jadesUnprotectedHeaderSchemaURI = "19182-unprotected-jsonSchema.json"

// JAdESUnprotectedHeaderUtils contains utils for JAdES Unprotected header validation against
// JSON schemas.
type JAdESUnprotectedHeaderUtils struct {
	abstractJAdESUtils
}

var (
	jadesUnprotectedHeaderUtilsSingleton     *JAdESUnprotectedHeaderUtils
	jadesUnprotectedHeaderUtilsSingletonOnce sync.Once
)

// JAdESUnprotectedHeaderUtilsInstance returns the singleton instance of
// JAdESUnprotectedHeaderUtils. Port of the private constructor + #getInstance lazy singleton.
func JAdESUnprotectedHeaderUtilsInstance() *JAdESUnprotectedHeaderUtils {
	jadesUnprotectedHeaderUtilsSingletonOnce.Do(func() {
		jadesUnprotectedHeaderUtilsSingleton = &JAdESUnprotectedHeaderUtils{}
	})
	return jadesUnprotectedHeaderUtilsSingleton
}

// SchemaURI ports the overridden #getSchemaURI.
func (u *JAdESUnprotectedHeaderUtils) SchemaURI() string {
	return jadesUnprotectedHeaderSchemaURI
}

// SchemaDefinitions ports the overridden #getSchemaDefinitions, adding this validator's own
// schema entry on top of the AbstractJAdESUtils base entries.
func (u *JAdESUnprotectedHeaderUtils) SchemaDefinitions() map[string]string {
	definitions := u.schemaDefinitions()
	definitions[jadesUnprotectedHeaderSchemaURI] = jadesUnprotectedHeaderSchemaURI
	return definitions
}

// ValidateAgainstSchema validates a JSON against the JAdES unprotected header schema per ETSI TS
// 119 182-1, returning a list of error messages - empty when the JSON is valid. Port of
// JSONSchemaAbstractUtils#validateAgainstSchema(String) as inherited here.
func (u *JAdESUnprotectedHeaderUtils) ValidateAgainstSchema(jsonString string) []string {
	return validateAgainstSchema(u.SchemaURI(), u.SchemaDefinitions(), jsonString)
}
