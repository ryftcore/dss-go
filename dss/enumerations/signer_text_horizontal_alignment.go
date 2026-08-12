// Ported from dss-enumerations/.../SignerTextHorizontalAlignment.java (DSS 6.5.RC1).
package enumerations

// SignerTextHorizontalAlignment defines the more line text horizontal
// alignment.
type SignerTextHorizontalAlignment string

const (
	// SignerTextHorizontalAlignment_LEFT aligns the text to left.
	SignerTextHorizontalAlignment_LEFT SignerTextHorizontalAlignment = "LEFT"
	// SignerTextHorizontalAlignment_CENTER centers the text.
	SignerTextHorizontalAlignment_CENTER SignerTextHorizontalAlignment = "CENTER"
	// SignerTextHorizontalAlignment_RIGHT aligns the text to right.
	SignerTextHorizontalAlignment_RIGHT SignerTextHorizontalAlignment = "RIGHT"
)

// SignerTextHorizontalAlignmentValues returns all constants in declaration
// order.
func SignerTextHorizontalAlignmentValues() []SignerTextHorizontalAlignment {
	return []SignerTextHorizontalAlignment{
		SignerTextHorizontalAlignment_LEFT,
		SignerTextHorizontalAlignment_CENTER,
		SignerTextHorizontalAlignment_RIGHT,
	}
}
