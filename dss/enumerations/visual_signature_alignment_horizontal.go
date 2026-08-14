// Ported from dss-enumerations/.../VisualSignatureAlignmentHorizontal.java (DSS 6.5.RC1).
package enumerations

// VisualSignatureAlignmentHorizontal is the visual signature horizontal
// position on the pdf page.
type VisualSignatureAlignmentHorizontal string

const (
	// VisualSignatureAlignmentHorizontal_NONE is the default; the x axis is
	// the x coordinate.
	VisualSignatureAlignmentHorizontal_NONE VisualSignatureAlignmentHorizontal = "NONE"
	// VisualSignatureAlignmentHorizontal_LEFT means the x axis is left
	// padding.
	VisualSignatureAlignmentHorizontal_LEFT VisualSignatureAlignmentHorizontal = "LEFT"
	// VisualSignatureAlignmentHorizontal_CENTER means the x axis is
	// automatically calculated.
	VisualSignatureAlignmentHorizontal_CENTER VisualSignatureAlignmentHorizontal = "CENTER"
	// VisualSignatureAlignmentHorizontal_RIGHT means the x axis is right
	// padding.
	VisualSignatureAlignmentHorizontal_RIGHT VisualSignatureAlignmentHorizontal = "RIGHT"
)

// VisualSignatureAlignmentHorizontalValues returns all constants in
// declaration order.
func VisualSignatureAlignmentHorizontalValues() []VisualSignatureAlignmentHorizontal {
	return []VisualSignatureAlignmentHorizontal{
		VisualSignatureAlignmentHorizontal_NONE,
		VisualSignatureAlignmentHorizontal_LEFT,
		VisualSignatureAlignmentHorizontal_CENTER,
		VisualSignatureAlignmentHorizontal_RIGHT,
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
