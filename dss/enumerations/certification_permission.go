// Ported from dss-enumerations/.../CertificationPermission.java (DSS 6.5.RC1).
//
// Refers to ISO 32000 DocMDP chapter.
package enumerations

import "fmt"

// CertificationPermission is used to set the allowed level of permission for
// PDF modifications.
type CertificationPermission string

const (
	// CertificationPermissionNoChangePermitted: no changes to the
	// document are permitted; any change to the document shall invalidate
	// the signature.
	CertificationPermissionNoChangePermitted CertificationPermission = "NO_CHANGE_PERMITTED"
	// CertificationPermissionMinimalChangesPermitted: permitted changes
	// shall be filling in forms, instantiating page templates, and
	// signing; other changes shall invalidate the signature.
	CertificationPermissionMinimalChangesPermitted CertificationPermission = "MINIMAL_CHANGES_PERMITTED"
	// CertificationPermissionChangesPermitted: permitted changes are the
	// same as for 2, as well as annotation creation, deletion, and
	// modification; other changes shall invalidate the signature.
	CertificationPermissionChangesPermitted CertificationPermission = "CHANGES_PERMITTED"
)

// certificationPermissionCodes holds the /DocMDP code for each constant.
var certificationPermissionCodes = map[CertificationPermission]int{
	CertificationPermissionNoChangePermitted:       1,
	CertificationPermissionMinimalChangesPermitted: 2,
	CertificationPermissionChangesPermitted:        3,
}

// CertificationPermissionValues returns all constants in declaration order.
func CertificationPermissionValues() []CertificationPermission {
	return []CertificationPermission{
		CertificationPermissionNoChangePermitted,
		CertificationPermissionMinimalChangesPermitted,
		CertificationPermissionChangesPermitted,
	}
}

// CertificationPermissionValueOf returns the CertificationPermission matching the given Java enum name.
func CertificationPermissionValueOf(name string) (CertificationPermission, error) {
	for _, v := range CertificationPermissionValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant CertificationPermission.%s", name)
}

// Code gets the value of the /DocMDP dictionary.
func (c CertificationPermission) Code() int {
	return certificationPermissionCodes[c]
}

// CertificationPermissionFromCode returns a CertificationPermission
// corresponding to the given code value.
func CertificationPermissionFromCode(code int) (CertificationPermission, error) {
	for _, v := range CertificationPermissionValues() {
		if code == v.Code() {
			return v, nil
		}
	}
	return "", fmt.Errorf("not supported /DocMDP code value : %d", code)
}
