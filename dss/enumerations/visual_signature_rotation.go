// Ported from dss-enumerations/.../VisualSignatureRotation.java (DSS 6.5.RC1).
package enumerations

// VisualSignatureRotation defines rotation support.
type VisualSignatureRotation string

const (
	// VisualSignatureRotation_NONE is the default, no rotate.
	VisualSignatureRotation_NONE VisualSignatureRotation = "NONE"
	// VisualSignatureRotation_AUTOMATIC automatically rotates.
	VisualSignatureRotation_AUTOMATIC VisualSignatureRotation = "AUTOMATIC"
	// VisualSignatureRotation_ROTATE_90 rotates by 90.
	VisualSignatureRotation_ROTATE_90 VisualSignatureRotation = "ROTATE_90"
	// VisualSignatureRotation_ROTATE_180 rotates by 180.
	VisualSignatureRotation_ROTATE_180 VisualSignatureRotation = "ROTATE_180"
	// VisualSignatureRotation_ROTATE_270 rotates by 270.
	VisualSignatureRotation_ROTATE_270 VisualSignatureRotation = "ROTATE_270"
)

// VisualSignatureRotationValues returns all constants in declaration order.
func VisualSignatureRotationValues() []VisualSignatureRotation {
	return []VisualSignatureRotation{
		VisualSignatureRotation_NONE,
		VisualSignatureRotation_AUTOMATIC,
		VisualSignatureRotation_ROTATE_90,
		VisualSignatureRotation_ROTATE_180,
		VisualSignatureRotation_ROTATE_270,
	}
}
