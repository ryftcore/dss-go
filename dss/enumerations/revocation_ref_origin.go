// Ported from dss-enumerations/.../RevocationRefOrigin.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// RevocationRefOrigin lists possible revocation reference origins.
type RevocationRefOrigin string

const (
	// RevocationRefOrigin_COMPLETE_REVOCATION_REFS: the revocation reference was found
	// in the signature 'complete-revocation-references' attribute (used in CAdES and XAdES).
	RevocationRefOrigin_COMPLETE_REVOCATION_REFS RevocationRefOrigin = "COMPLETE_REVOCATION_REFS"
	// RevocationRefOrigin_ATTRIBUTE_REVOCATION_REFS: the revocation reference was found
	// in the signature 'attribute-revocation-references' attribute (used in CAdES and XAdES).
	RevocationRefOrigin_ATTRIBUTE_REVOCATION_REFS RevocationRefOrigin = "ATTRIBUTE_REVOCATION_REFS"
)

// RevocationRefOriginValues returns all RevocationRefOrigin constants in declaration order.
func RevocationRefOriginValues() []RevocationRefOrigin {
	return []RevocationRefOrigin{
		RevocationRefOrigin_COMPLETE_REVOCATION_REFS,
		RevocationRefOrigin_ATTRIBUTE_REVOCATION_REFS,
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
