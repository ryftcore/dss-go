// Ported from specs-jades/src/main/java/eu/europa/esig/jades/JAdESProtectedHeaderUtils.java (DSS 6.5.RC1).
//
// Contains utils for JAdES Protected header validation against JSON schemas.
package specs

import "sync"

// jadesProtectedHeaderSchemaURI is the protected header schema URI for a JAdES signature. Port
// of JAdES_PROTECTED_HEADER_SCHEMA_URI.
const jadesProtectedHeaderSchemaURI = "19182-protected-jsonSchema.json"

// JAdESProtectedHeaderUtils contains utils for JAdES Protected header validation against JSON
// schemas.
type JAdESProtectedHeaderUtils struct {
	abstractJAdESUtils
}

var (
	jadesProtectedHeaderUtilsSingleton     *JAdESProtectedHeaderUtils
	jadesProtectedHeaderUtilsSingletonOnce sync.Once
)

// JAdESProtectedHeaderUtilsInstance returns the singleton instance of JAdESProtectedHeaderUtils.
// Port of the private constructor + #getInstance lazy singleton.
func JAdESProtectedHeaderUtilsInstance() *JAdESProtectedHeaderUtils {
	jadesProtectedHeaderUtilsSingletonOnce.Do(func() {
		jadesProtectedHeaderUtilsSingleton = &JAdESProtectedHeaderUtils{}
	})
	return jadesProtectedHeaderUtilsSingleton
}

// SchemaURI ports the overridden #getSchemaURI.
func (u *JAdESProtectedHeaderUtils) SchemaURI() string {
	return jadesProtectedHeaderSchemaURI
}

// SchemaDefinitions ports the overridden #getSchemaDefinitions, adding this validator's own
// schema entry on top of the AbstractJAdESUtils base entries.
func (u *JAdESProtectedHeaderUtils) SchemaDefinitions() map[string]string {
	definitions := u.schemaDefinitions()
	definitions[jadesProtectedHeaderSchemaURI] = jadesProtectedHeaderSchemaURI
	return definitions
}

// ValidateAgainstSchema validates a JSON against the JAdES protected header schema per ETSI TS
// 119 182-1, returning a list of error messages - empty when the JSON is valid. Port of
// JSONSchemaAbstractUtils#validateAgainstSchema(String) as inherited here.
func (u *JAdESProtectedHeaderUtils) ValidateAgainstSchema(jsonString string) []string {
	return validateAgainstSchema(u.SchemaURI(), u.SchemaDefinitions(), jsonString)
}
