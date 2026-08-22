package pades

import "math"

// pdfNumberToInt narrows a PDF number to an int, reporting false when the value
// cannot be represented. Go leaves a float-to-integer conversion whose operand is
// out of range implementation-defined (spec, "Conversions between numeric types"),
// so an untrusted /ByteRange or /P entry must be range-checked before it is
// converted. The bound is Java's: upstream holds these values in an int[] / int,
// so anything outside the 32-bit range was never representable upstream either.
func pdfNumberToInt(value float64) (int, bool) {
	if math.IsNaN(value) || value < math.MinInt32 || value > math.MaxInt32 {
		return 0, false
	}
	return int(value), true
}
