// Ported from dss-enumerations/.../VisualSignatureAlignmentHorizontal.java (DSS 6.5.RC1).
package enumerations

// VisualSignatureAlignmentHorizontal is the visual signature horizontal
// position on the pdf page.
type VisualSignatureAlignmentHorizontal string

const (
	// VisualSignatureAlignmentHorizontalNone is the default; the x axis is
	// the x coordinate.
	VisualSignatureAlignmentHorizontalNone VisualSignatureAlignmentHorizontal = "NONE"
	// VisualSignatureAlignmentHorizontalLeft means the x axis is left
	// padding.
	VisualSignatureAlignmentHorizontalLeft VisualSignatureAlignmentHorizontal = "LEFT"
	// VisualSignatureAlignmentHorizontalCenter means the x axis is
	// automatically calculated.
	VisualSignatureAlignmentHorizontalCenter VisualSignatureAlignmentHorizontal = "CENTER"
	// VisualSignatureAlignmentHorizontalRight means the x axis is right
	// padding.
	VisualSignatureAlignmentHorizontalRight VisualSignatureAlignmentHorizontal = "RIGHT"
)

// VisualSignatureAlignmentHorizontalValues returns all constants in
// declaration order.
func VisualSignatureAlignmentHorizontalValues() []VisualSignatureAlignmentHorizontal {
	return []VisualSignatureAlignmentHorizontal{
		VisualSignatureAlignmentHorizontalNone,
		VisualSignatureAlignmentHorizontalLeft,
		VisualSignatureAlignmentHorizontalCenter,
		VisualSignatureAlignmentHorizontalRight,
	}
}

// VisualSignatureAlignmentHorizontalValueOf returns the constant matching
// the given Java enum name.
func VisualSignatureAlignmentHorizontalValueOf(name string) (VisualSignatureAlignmentHorizontal, error) {
	for _, v := range VisualSignatureAlignmentHorizontalValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &visualSignatureAlignmentHorizontalInvalidValueError{name}
}

type visualSignatureAlignmentHorizontalInvalidValueError struct {
	name string
}

func (e *visualSignatureAlignmentHorizontalInvalidValueError) Error() string {
	return "no enum constant VisualSignatureAlignmentHorizontal." + e.name
}
