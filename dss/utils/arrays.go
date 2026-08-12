// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1),
// array-related methods, matching org.apache.commons.lang3.ArrayUtils
// semantics.
//
// Java overloads isArrayEmpty/isArrayNotEmpty/arraySize for Object[],
// byte[] and char[]; Go generics collapse those three overloads into one
// implementation each, applicable to any slice type.

package utils

// IsArrayEmpty checks if the array is nil or empty.
func IsArrayEmpty[T any](array []T) bool {
	return len(array) == 0
}

// IsArrayNotEmpty checks if the array is not nil nor empty.
func IsArrayNotEmpty[T any](array []T) bool {
	return len(array) != 0
}

// ArraySize gets the size of the array. If nil, returns 0.
func ArraySize[T any](array []T) int {
	return len(array)
}

// Subarray returns a copy of array from index start (inclusive) to index
// end (exclusive), clipped to the array bounds.
//
// NOTE on naming: upstream IUtils/ApacheCommonsUtils names the third
// parameter "length" but passes it straight through to
// org.apache.commons.lang3.ArrayUtils.subarray(array, startIndexInclusive,
// endIndexExclusive) — i.e. despite its Java name it is actually an END
// INDEX, not a length. This port keeps the same behavior and names the
// parameter accordingly to avoid perpetuating the confusing upstream name.
//
// A nil array returns nil. start < 0 is treated as 0. end beyond the
// array length is clipped to the array length. If start >= end (after
// clipping) an empty (non-nil) slice is returned.
func Subarray(array []byte, start, end int) []byte {
	if array == nil {
		return nil
	}
	if start < 0 {
		start = 0
	}
	if end > len(array) {
		end = len(array)
	}
	if start >= end {
		return []byte{}
	}
	out := make([]byte, end-start)
	copy(out, array[start:end])
	return out
}

// Concat concatenates byte arrays into a single new byte array. The new
// array contains all bytes of each array followed by all bytes of the
// next array. A nil element is treated as empty (Java throws NPE on a nil
// element; this port takes the nil/empty-safe route consistent with the
// rest of the package).
func Concat(byteArrays ...[]byte) []byte {
	total := 0
	for _, b := range byteArrays {
		total += len(b)
	}
	result := make([]byte, 0, total)
	for _, b := range byteArrays {
		result = append(result, b...)
	}
	return result
}

// StartsWith checks if byteArray starts with prefixArray.
func StartsWith(byteArray, prefixArray []byte) bool {
	if byteArray == nil || prefixArray == nil {
		return false
	}
	if len(byteArray) < len(prefixArray) {
		return false
	}
	for i := range prefixArray {
		if byteArray[i] != prefixArray[i] {
			return false
		}
	}
	return true
}
