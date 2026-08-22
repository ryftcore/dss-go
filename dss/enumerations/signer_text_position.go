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
	// SignerTextPositionTop: the text on the top of the picture.
	SignerTextPositionTop SignerTextPosition = "TOP"
	// SignerTextPositionBottom: the text on the bottom of the picture.
	SignerTextPositionBottom SignerTextPosition = "BOTTOM"
	// SignerTextPositionRight: the text on the right of the picture.
	SignerTextPositionRight SignerTextPosition = "RIGHT"
	// SignerTextPositionLeft: the text on the left of the picture.
	SignerTextPositionLeft SignerTextPosition = "LEFT"
)

// SignerTextPositionValues returns all constants in declaration order.
func SignerTextPositionValues() []SignerTextPosition {
	return []SignerTextPosition{
		SignerTextPositionTop,
		SignerTextPositionBottom,
		SignerTextPositionRight,
		SignerTextPositionLeft,
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
