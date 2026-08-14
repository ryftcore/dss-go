// Ported from dss-model/.../claim/ClaimBiometricTemplateXX.java (DSS 6.5.RC1).
package claim

// ClaimBiometricTemplateXX contains biometric data of an EAA holder. The
// type of data can be extracted using the Type method.
type ClaimBiometricTemplateXX interface {
	Claim

	// Type gets type of the biometric data. Ports
	// ClaimBiometricTemplateXX#getType.
	Type() string
}
