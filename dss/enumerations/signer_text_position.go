// Ported from dss-enumerations/.../SignerTextPosition.java (DSS 6.5.RC1).
//
// Enum to define where to add a signer text inside a signature field
// relatively to an image.
package enumerations

import "fmt"

// SignerTextPosition defines where to add a signer text inside a signature
// field relative to an image.
type SignerTextPosition string

const (
	// SignerTextPosition_TOP: the text on the top of the picture.
	SignerTextPosition_TOP SignerTextPosition = "TOP"
	// SignerTextPosition_BOTTOM: the text on the bottom of the picture.
	SignerTextPosition_BOTTOM SignerTextPosition = "BOTTOM"
	// SignerTextPosition_RIGHT: the text on the right of the picture.
	SignerTextPosition_RIGHT SignerTextPosition = "RIGHT"
	// SignerTextPosition_LEFT: the text on the left of the picture.
	SignerTextPosition_LEFT SignerTextPosition = "LEFT"
)

// SignerTextPositionValues returns all constants in declaration order.
func SignerTextPositionValues() []SignerTextPosition {
	return []SignerTextPosition{
		SignerTextPosition_TOP,
		SignerTextPosition_BOTTOM,
		SignerTextPosition_RIGHT,
		SignerTextPosition_LEFT,
	}
}

// SignerTextPositionValueOf returns the SignerTextPosition matching the given Java enum name.
func SignerTextPositionValueOf(name string) (SignerTextPosition, error) {
	for _, v := range SignerTextPositionValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant SignerTextPosition.%s", name)
}
