// Ported from dss-enumerations/.../VisualSignatureAlignmentVertical.java (DSS 6.5.RC1).
package enumerations

// VisualSignatureAlignmentVertical is the visual signature vertical
// position on the pdf page.
type VisualSignatureAlignmentVertical string

const (
	// VisualSignatureAlignmentVertical_NONE is the default: y axis is the y
	// coordinate.
	VisualSignatureAlignmentVertical_NONE VisualSignatureAlignmentVertical = "NONE"
	// VisualSignatureAlignmentVertical_TOP: y axis is the top padding.
	VisualSignatureAlignmentVertical_TOP VisualSignatureAlignmentVertical = "TOP"
	// VisualSignatureAlignmentVertical_MIDDLE: y axis automatically
	// calculated.
	VisualSignatureAlignmentVertical_MIDDLE VisualSignatureAlignmentVertical = "MIDDLE"
	// VisualSignatureAlignmentVertical_BOTTOM: y axis is the bottom
	// padding.
	VisualSignatureAlignmentVertical_BOTTOM VisualSignatureAlignmentVertical = "BOTTOM"
)

// VisualSignatureAlignmentVerticalValues returns all constants in
// declaration order.
func VisualSignatureAlignmentVerticalValues() []VisualSignatureAlignmentVertical {
	return []VisualSignatureAlignmentVertical{
		VisualSignatureAlignmentVertical_NONE,
		VisualSignatureAlignmentVertical_TOP,
		VisualSignatureAlignmentVertical_MIDDLE,
		VisualSignatureAlignmentVertical_BOTTOM,
	}
}
