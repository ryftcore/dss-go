package pades

import (
	"math"
	"testing"
)

// TestPdfNumberToInt pins pdfNumberToInt's bounds. Go leaves a float-to-integer
// conversion whose operand is out of the target type's range implementation-defined
// (CodeQL go/incorrect-integer-conversion, alerts 38-40); this asserts the guard
// actually rejects everything outside int32 rather than merely compiling.
func TestPdfNumberToInt(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  int
		ok    bool
	}{
		{"zero", 0, 0, true},
		{"negative one", -1, -1, true},
		{"int32 max is in range", math.MaxInt32, math.MaxInt32, true},
		{"int32 min is in range", math.MinInt32, math.MinInt32, true},
		{"one past int32 max is rejected", math.MaxInt32 + 1, 0, false},
		{"one before int32 min is rejected", math.MinInt32 - 1, 0, false},
		{"far beyond range is rejected", 1e300, 0, false},
		{"negative infinity is rejected", math.Inf(-1), 0, false},
		{"positive infinity is rejected", math.Inf(1), 0, false},
		{"NaN is rejected", math.NaN(), 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := pdfNumberToInt(tc.value)
			if ok != tc.ok {
				t.Fatalf("pdfNumberToInt(%v) ok = %v, want %v", tc.value, ok, tc.ok)
			}
			if ok && got != tc.want {
				t.Errorf("pdfNumberToInt(%v) = %v, want %v", tc.value, got, tc.want)
			}
		})
	}
}

// fakeByteRangeArray is a minimal PdfArray that implements only what byteRange()
// touches (Size and Number). Everything else panics through the embedded nil
// PdfArray if it is ever called, which is deliberate: it would mean the test needs
// to model more of the interface.
type fakeByteRangeArray struct {
	PdfArray
	numbers []*float64
}

func (a *fakeByteRangeArray) Size() int             { return len(a.numbers) }
func (a *fakeByteRangeArray) Number(i int) *float64 { return a.numbers[i] }

// fakeSigFieldDict is a minimal PdfDict that only serves a /ByteRange array back
// from AsArray, so byteRange() can be exercised without a real PDF fixture.
type fakeSigFieldDict struct {
	PdfDict
	byteRangeArray PdfArray
}

func (d *fakeSigFieldDict) AsArray(name string) PdfArray {
	if name == PAdESConstantsByteRangeName {
		return d.byteRangeArray
	}
	return nil
}

// TestByteRangeRejectsOutOfRangeNumber is the alert-39 regression: a /ByteRange
// entry outside the int32 range used to be narrowed with the platform-dependent
// int(*float64) conversion. It must now surface as an error from byteRange(),
// not a corrupted *ByteRange.
func TestByteRangeRejectsOutOfRangeNumber(t *testing.T) {
	zero := 0.0
	huge := 1e300
	dict := &fakeSigFieldDict{
		byteRangeArray: &fakeByteRangeArray{
			numbers: []*float64{&zero, &huge, &zero, &zero},
		},
	}
	f := NewPdfSigDictWrapperFactory(dict)
	got, err := f.byteRange()
	if err == nil {
		t.Fatalf("byteRange() = %v, %v; want a non-nil error for an out-of-range entry", got, err)
	}
}

// TestByteRangeAcceptsInRangeNumbers is the control case: well-formed /ByteRange
// values must still parse into the expected *ByteRange.
func TestByteRangeAcceptsInRangeNumbers(t *testing.T) {
	a := 0.0
	b := 840.0
	c := 960.0
	d := 240.0
	dict := &fakeSigFieldDict{
		byteRangeArray: &fakeByteRangeArray{
			numbers: []*float64{&a, &b, &c, &d},
		},
	}
	f := NewPdfSigDictWrapperFactory(dict)
	got, err := f.byteRange()
	if err != nil {
		t.Fatalf("byteRange() error = %v, want nil", err)
	}
	if got.FirstPartStart() != 0 || got.FirstPartEnd() != 840 ||
		got.SecondPartStart() != 960 || got.SecondPartEnd() != 240 {
		t.Errorf("byteRange() = [%d %d %d %d], want [0 840 960 240]",
			got.FirstPartStart(), got.FirstPartEnd(), got.SecondPartStart(), got.SecondPartEnd())
	}
}
