// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/SignatureProperties.java (DSS 6.5.RC1).
package validation

// SignatureProperties defines a "signed-signature-element" or "unsigned-signature-element" of
// a signature. SA is the implementation of a signature attribute (signed or unsigned)
// corresponding to the current signature format; Java's `<SA extends SignatureAttribute>`
// bound becomes a Go generic type parameter.
type SignatureProperties[SA SignatureAttribute] interface {
	// IsExist checks if "unsigned-signature-properties" exists and can be processed. Port of
	// isExist().
	IsExist() bool

	// Attributes returns a list of children contained in the element. Port of getAttributes().
	Attributes() []SA
}
