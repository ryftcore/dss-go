// Ported from dss-enumerations/.../SignerTextVerticalAlignment.java (DSS 6.5.RC1).
package enumerations

// SignerTextVerticalAlignment defines image from text vertical alignment
// in connection with the image.
type SignerTextVerticalAlignment string

const (
	// SignerTextVerticalAlignmentTop aligns the text with the top of the
	// picture.
	SignerTextVerticalAlignmentTop SignerTextVerticalAlignment = "TOP"
	// SignerTextVerticalAlignmentMiddle aligns the text with the center
	// of the picture.
	SignerTextVerticalAlignmentMiddle SignerTextVerticalAlignment = "MIDDLE"
	// SignerTextVerticalAlignmentBottom aligns the text with the bottom
	// of the picture.
	SignerTextVerticalAlignmentBottom SignerTextVerticalAlignment = "BOTTOM"
)

// SignerTextVerticalAlignmentValues returns all constants in declaration order.
func SignerTextVerticalAlignmentValues() []SignerTextVerticalAlignment {
	return []SignerTextVerticalAlignment{
		SignerTextVerticalAlignmentTop,
		SignerTextVerticalAlignmentMiddle,
		SignerTextVerticalAlignmentBottom,
	}
}
