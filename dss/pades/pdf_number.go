package pades

import "math"

// pdfNumberToInt narrows a PDF number to an int, reporting false when the value
// cannot be represented. Go leaves a float-to-integer conversion whose operand is
// out of range implementation-defined, so an untrusted /ByteRange or /P entry must
// be range-checked before it is converted. The bound is Java's: upstream holds
// these values in an int[] / int, so anything outside the 32-bit range was never
// representable upstream either.
func pdfNumberToInt(value float64) (int, bool) {
	// Reject NaN, ±Inf and anything beyond 2^53 first: within that range a float64
	// is an exact integer and the int64 conversion below is lossless.
	if math.IsNaN(value) || value < -(1<<53) || value > 1<<53 {
		return 0, false
	}
	wide := int64(value)
	if wide < math.MinInt32 || wide > math.MaxInt32 {
		return 0, false
	}
	return int(wide), true
}
