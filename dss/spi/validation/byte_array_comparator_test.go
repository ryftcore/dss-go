// Ported from dss-spi/src/test/java/eu/europa/esig/dss/spi/validation/evidencerecord/ByteArrayComparatorTest.java (DSS 6.5.RC1) -- no upstream test exists; this is a fresh table test per PORTING.md's "every registry-like table gets an exhaustive table test" guidance, applied here to the pure comparator function.
package validation

import "testing"

func TestByteArrayComparatorCompare(t *testing.T) {
	tests := []struct {
		name string
		a, b []byte
		want int
	}{
		{"equal", []byte{1, 2, 3}, []byte{1, 2, 3}, 0},
		{"less byte", []byte{1, 2, 3}, []byte{1, 2, 4}, -1},
		{"greater byte", []byte{1, 2, 4}, []byte{1, 2, 3}, 1},
		{"prefix shorter", []byte{1, 2}, []byte{1, 2, 3}, -1},
		{"prefix longer", []byte{1, 2, 3}, []byte{1, 2}, 1},
		{"empty vs empty", []byte{}, []byte{}, 0},
		{"empty vs non-empty", []byte{}, []byte{1}, -1},
		{"unsigned comparison", []byte{0x80}, []byte{0x7f}, 1}, // 0x80 as unsigned (128) > 0x7f (127)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByteArrayComparatorCompare(tt.a, tt.b); got != tt.want {
				t.Errorf("ByteArrayComparatorCompare(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
			// symmetry: compare(b, a) should have the opposite sign (or both zero).
			inverse := ByteArrayComparatorCompare(tt.b, tt.a)
			if (tt.want > 0 && inverse >= 0) || (tt.want < 0 && inverse <= 0) || (tt.want == 0 && inverse != 0) {
				t.Errorf("ByteArrayComparatorCompare(%v, %v) = %d not symmetric with reverse = %d", tt.a, tt.b, tt.want, inverse)
			}
		})
	}
}
