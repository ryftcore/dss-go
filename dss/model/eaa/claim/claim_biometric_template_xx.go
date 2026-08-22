// Ported from dss-model/.../claim/ClaimBiometricTemplateXX.java (DSS 6.5.RC1).
package claim

// BiometricTemplateXX contains biometric data of an EAA holder. The
// type of data can be extracted using the Type method.
type BiometricTemplateXX interface {
	Claim

	// Type gets type of the biometric data. Ports
	// ClaimBiometricTemplateXX#getType.
	Type() string
}
