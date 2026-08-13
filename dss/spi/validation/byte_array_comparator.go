// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/evidencerecord/ByteArrayComparator.java (DSS 6.5.RC1).
package validation

// ByteArrayComparator compares two byte arrays lexicographically, treating each byte as
// unsigned (matching Java's `& 0xff` masking), and falls back to length when one array is a
// byte-for-byte prefix of the other.
//
// Inspired by
// https://github.com/bcgit/bc-java/blob/main/pkix/src/main/java/org/bouncycastle/tsp/ers/ByteArrayComparator.java
//
// Java models this as a singleton (getInstance()) implementing Comparator<byte[]>; Go has no
// class identity to guard, so ByteArrayComparatorCompare is exposed directly as a function
// instead of manufacturing a singleton value around it.
func ByteArrayComparatorCompare(o1, o2 []byte) int {
	for i := 0; i < len(o1) && i < len(o2); i++ {
		a := o1[i]
		b := o2[i]
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
	}
	return len(o1) - len(o2)
}
