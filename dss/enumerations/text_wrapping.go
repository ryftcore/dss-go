// Ported from dss-enumerations/.../TextWrapping.java (DSS 6.5.RC1).
//
// This enumeration defines a set of possibilities for text wrapping within a
// signature field with a fixed width and height for a PDF visual signature
// creation.
package enumerations

import "fmt"

// TextWrapping defines text wrapping options for a signature field.
type TextWrapping string

const (
	// TextWrappingFillBox: a font size is adapted in order to fill the
	// whole signature field's space, by keeping the defined whitespaces in
	// new lines by user. When using with a combination of image, the image
	// block is computed at first and the rest space is filled by text.
	TextWrappingFillBox TextWrapping = "FILL_BOX"
	// TextWrappingFillBoxAndLineBreak: the text is formatted, by
	// separating the provided text to multiple lines in order to find the
	// biggest font size in order to wrap the text to the defined signature
	// field's box. When using with a combination of image, the image block
	// is computed at first and the rest space is filled by text.
	TextWrappingFillBoxAndLineBreak TextWrapping = "FILL_BOX_AND_LINEBREAK"
	// TextWrappingFontBased: the text is generated based on the font
	// values provided within parameters. When using this value with
	// combination of image, the text is computed at first and the rest
	// space is filled by the image.
	TextWrappingFontBased TextWrapping = "FONT_BASED"
)

// TextWrappingValues returns all constants in declaration order.
func TextWrappingValues() []TextWrapping {
	return []TextWrapping{
		TextWrappingFillBox,
		TextWrappingFillBoxAndLineBreak,
		TextWrappingFontBased,
	}
}

// TextWrappingValueOf returns the TextWrapping matching the given Java enum name.
func TextWrappingValueOf(name string) (TextWrapping, error) {
	for _, v := range TextWrappingValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant TextWrapping.%s", name)
}
