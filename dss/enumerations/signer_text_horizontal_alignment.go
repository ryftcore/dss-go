// Ported from dss-enumerations/.../SignerTextHorizontalAlignment.java (DSS 6.5.RC1).
package enumerations

// SignerTextHorizontalAlignment defines the more line text horizontal
// alignment.
type SignerTextHorizontalAlignment string

const (
	// SignerTextHorizontalAlignmentLeft aligns the text to left.
	SignerTextHorizontalAlignmentLeft SignerTextHorizontalAlignment = "LEFT"
	// SignerTextHorizontalAlignmentCenter centers the text.
	SignerTextHorizontalAlignmentCenter SignerTextHorizontalAlignment = "CENTER"
	// SignerTextHorizontalAlignmentRight aligns the text to right.
	SignerTextHorizontalAlignmentRight SignerTextHorizontalAlignment = "RIGHT"
)

// SignerTextHorizontalAlignmentValues returns all constants in declaration
// order.
func SignerTextHorizontalAlignmentValues() []SignerTextHorizontalAlignment {
	return []SignerTextHorizontalAlignment{
		SignerTextHorizontalAlignmentLeft,
		SignerTextHorizontalAlignmentCenter,
		SignerTextHorizontalAlignmentRight,
	}
}
