// Ported from dss-enumerations/.../RevocationRefOrigin.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// RevocationRefOrigin lists possible revocation reference origins.
type RevocationRefOrigin string

const (
	// RevocationRefOriginCompleteRevocationRefs: the revocation reference was found
	// in the signature 'complete-revocation-references' attribute (used in CAdES and XAdES).
	RevocationRefOriginCompleteRevocationRefs RevocationRefOrigin = "COMPLETE_REVOCATION_REFS"
	// RevocationRefOriginAttributeRevocationRefs: the revocation reference was found
	// in the signature 'attribute-revocation-references' attribute (used in CAdES and XAdES).
	RevocationRefOriginAttributeRevocationRefs RevocationRefOrigin = "ATTRIBUTE_REVOCATION_REFS"
)

// RevocationRefOriginValues returns all RevocationRefOrigin constants in declaration order.
func RevocationRefOriginValues() []RevocationRefOrigin {
	return []RevocationRefOrigin{
		RevocationRefOriginCompleteRevocationRefs,
		RevocationRefOriginAttributeRevocationRefs,
	}
}

// RevocationRefOriginValueOf returns the RevocationRefOrigin matching the given Java enum name.
func RevocationRefOriginValueOf(name string) (RevocationRefOrigin, error) {
	for _, v := range RevocationRefOriginValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant RevocationRefOrigin.%s", name)
}
