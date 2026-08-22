// Ported from dss-enumerations/.../VisualSignatureAlignmentVertical.java (DSS 6.5.RC1).
package enumerations

// VisualSignatureAlignmentVertical is the visual signature vertical
// position on the pdf page.
type VisualSignatureAlignmentVertical string

const (
	// VisualSignatureAlignmentVerticalNone is the default: y axis is the y
	// coordinate.
	VisualSignatureAlignmentVerticalNone VisualSignatureAlignmentVertical = "NONE"
	// VisualSignatureAlignmentVerticalTop: y axis is the top padding.
	VisualSignatureAlignmentVerticalTop VisualSignatureAlignmentVertical = "TOP"
	// VisualSignatureAlignmentVerticalMiddle: y axis automatically
	// calculated.
	VisualSignatureAlignmentVerticalMiddle VisualSignatureAlignmentVertical = "MIDDLE"
	// VisualSignatureAlignmentVerticalBottom: y axis is the bottom
	// padding.
	VisualSignatureAlignmentVerticalBottom VisualSignatureAlignmentVertical = "BOTTOM"
)

// VisualSignatureAlignmentVerticalValues returns all constants in
// declaration order.
func VisualSignatureAlignmentVerticalValues() []VisualSignatureAlignmentVertical {
	return []VisualSignatureAlignmentVertical{
		VisualSignatureAlignmentVerticalNone,
		VisualSignatureAlignmentVerticalTop,
		VisualSignatureAlignmentVerticalMiddle,
		VisualSignatureAlignmentVerticalBottom,
	}
}
