// Ported from dss-enumerations/.../VisualSignatureRotation.java (DSS 6.5.RC1).
package enumerations

// VisualSignatureRotation defines rotation support.
type VisualSignatureRotation string

const (
	// VisualSignatureRotationNone is the default, no rotate.
	VisualSignatureRotationNone VisualSignatureRotation = "NONE"
	// VisualSignatureRotationAutomatic automatically rotates.
	VisualSignatureRotationAutomatic VisualSignatureRotation = "AUTOMATIC"
	// VisualSignatureRotationRotate90 rotates by 90.
	VisualSignatureRotationRotate90 VisualSignatureRotation = "ROTATE_90"
	// VisualSignatureRotationRotate180 rotates by 180.
	VisualSignatureRotationRotate180 VisualSignatureRotation = "ROTATE_180"
	// VisualSignatureRotationRotate270 rotates by 270.
	VisualSignatureRotationRotate270 VisualSignatureRotation = "ROTATE_270"
)

// VisualSignatureRotationValues returns all constants in declaration order.
func VisualSignatureRotationValues() []VisualSignatureRotation {
	return []VisualSignatureRotation{
		VisualSignatureRotationNone,
		VisualSignatureRotationAutomatic,
		VisualSignatureRotationRotate90,
		VisualSignatureRotationRotate180,
		VisualSignatureRotationRotate270,
	}
}
