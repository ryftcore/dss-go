// Ported from dss-enumerations/.../SignerTextVerticalAlignment.java (DSS 6.5.RC1).
package enumerations

// SignerTextVerticalAlignment defines image from text vertical alignment
// in connection with the image.
type SignerTextVerticalAlignment string

const (
	// SignerTextVerticalAlignment_TOP aligns the text with the top of the
	// picture.
	SignerTextVerticalAlignment_TOP SignerTextVerticalAlignment = "TOP"
	// SignerTextVerticalAlignment_MIDDLE aligns the text with the center
	// of the picture.
	SignerTextVerticalAlignment_MIDDLE SignerTextVerticalAlignment = "MIDDLE"
	// SignerTextVerticalAlignment_BOTTOM aligns the text with the bottom
	// of the picture.
	SignerTextVerticalAlignment_BOTTOM SignerTextVerticalAlignment = "BOTTOM"
)

// SignerTextVerticalAlignmentValues returns all constants in declaration order.
func SignerTextVerticalAlignmentValues() []SignerTextVerticalAlignment {
	return []SignerTextVerticalAlignment{
		SignerTextVerticalAlignment_TOP,
		SignerTextVerticalAlignment_MIDDLE,
		SignerTextVerticalAlignment_BOTTOM,
	}
}
