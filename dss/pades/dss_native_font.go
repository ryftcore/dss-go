// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/DSSNativeFont.java (DSS 6.5.RC1).
//
// Java's `<F extends Object>` type parameter is unbounded generics; Go's `any` bound is the
// direct equivalent.
package pades

// DSSNativeFont is the native font used in PDF libraries, generic over F the class of the font
// instance.
type DSSNativeFont[F any] interface {
	// Font returns a native font for the given implementation. Port of #getFont.
	Font() F
}
