package pdf

import (
	"errors"
	"math"
	"testing"
)

// TestIsValidGeneration pins the ISO 32000-1 Table 18 bound isValidGeneration
// enforces: a generation number is 0..65535, exactly what ObjectKey.Gen
// (uint16) can hold without wrapping into some other object's key (CodeQL
// go/incorrect-integer-conversion, alerts 31-35).
func TestIsValidGeneration(t *testing.T) {
	tests := []struct {
		name string
		gen  int64
		want bool
	}{
		{"zero is valid", 0, true},
		{"max uint16 is valid", 65535, true},
		{"negative one is invalid", -1, false},
		{"one past max uint16 is invalid", 65536, false},
		{"far beyond range is invalid", math.MaxInt64, false},
		{"far below range is invalid", math.MinInt64, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidGeneration(tc.gen); got != tc.want {
				t.Errorf("isValidGeneration(%d) = %v, want %v", tc.gen, got, tc.want)
			}
		})
	}
}

// TestXRefTableParsingDropsInvalidGenerations exercises parseXrefTableAt
// directly against a hand-built classic xref table: strconv.ParseInt(fields[1],
// 10, 32) accepts a negative literal, and uint16(-1) is 65535, so without the
// isValidGeneration guard a "0000000010 -1 n" entry would be recorded as
// generation 65535 - a wrapped value that can collide with a legitimate entry
// (CodeQL alert 32/33).
func TestXRefTableParsingDropsInvalidGenerations(t *testing.T) {
	data := []byte("xref\n" +
		"0 4\n" +
		"0000000000 65535 f\r\n" +
		"0000000010 -1    n\r\n" + // negative generation: must be dropped entirely
		"0000000020 65535 n\r\n" + // the maximum valid generation: must be kept
		"0000000030 65536 n\r\n" + // one past the maximum: must be dropped entirely
		"trailer\n<< /Size 4 /Root 1 0 R >>\n")

	var warn []Warning
	p := newParser(data, &warn, 512, 1<<20)
	d := &Document{}
	rs, err := d.parseXrefTableAt(p, 0, false)
	if err != nil {
		t.Fatalf("parseXrefTableAt: %v", err)
	}

	got := make(map[ObjectKey]bool, len(rs.entries))
	for _, e := range rs.entries {
		got[e.key] = true
	}
	if got[ObjectKey{Num: 1, Gen: 65535}] {
		t.Errorf("entry with a negative generation literal must be dropped, not wrapped into generation 65535: entries = %v", got)
	}
	if !got[ObjectKey{Num: 2, Gen: 65535}] {
		t.Errorf("entry with generation 65535 (the maximum ISO 32000-1 Table 18 allows) must be kept: entries = %v", got)
	}
	if got[ObjectKey{Num: 3, Gen: 0}] {
		t.Errorf("entry with generation 65536 must be dropped, not wrapped into generation 0: entries = %v", got)
	}
	for _, e := range rs.entries {
		if e.key.Num == 3 {
			t.Errorf("object 3's out-of-range-generation entry must not be recorded under any generation, got %v", e.key)
		}
	}
}

// TestParseIndirectAtRejectsOversizedGeneration is the parser.go half of the
// same guard (alert 31): an `N G obj` header whose G is a huge but syntactically
// valid integer literal (not clamped by the lexer, since it is well inside the
// int64 range) must not be truncated into a plausible ObjectKey.
func TestParseIndirectAtRejectsOversizedGeneration(t *testing.T) {
	var warn []Warning
	p := newParser([]byte("1 4294967296 obj\n<< >>\nendobj\n"), &warn, 512, 1<<20)
	if _, _, ok := p.parseIndirectAt(0, ObjectKey{Num: 1}); ok {
		t.Errorf("parseIndirectAt with generation 4294967296 (2^32) succeeded, want ok == false")
	}
}

// TestClampDictInt pins clampDictInt's saturation at the int32 bounds: every
// /Encrypt entry it feeds (/V, /R, /Length) has a small spec-defined range, so
// saturating rather than erroring cannot change the reading of a well-formed
// document; it only stops a hostile 64-bit value from wrapping into a plausible
// one where int is 32 bits wide (CodeQL alerts 22-27).
func TestClampDictInt(t *testing.T) {
	tests := []struct {
		name string
		in   int64
		want int
	}{
		{"zero passes through", 0, 0},
		{"in-range positive passes through", 5, 5},
		{"in-range negative passes through", -5, -5},
		{"int32 max passes through", math.MaxInt32, math.MaxInt32},
		{"int32 min passes through", math.MinInt32, math.MinInt32},
		{"one past int32 max saturates", math.MaxInt32 + 1, math.MaxInt32},
		{"far beyond int32 max saturates", 4294967300, math.MaxInt32},
		{"int64 max saturates", math.MaxInt64, math.MaxInt32},
		{"one before int32 min saturates", math.MinInt32 - 1, math.MinInt32},
		{"int64 min saturates", math.MinInt64, math.MinInt32},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := clampDictInt(tc.in); got != tc.want {
				t.Errorf("clampDictInt(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// TestEncryptionVSaturatesRatherThanWraps is the end-to-end regression for
// clampDictInt: on a 32-bit build, int(int64(4294967300)) wraps to a small
// positive int that could alias a supported /V. Saturating at MaxInt32 makes the
// unsupported-handler rejection architecture-independent instead.
func TestEncryptionVSaturatesRatherThanWraps(t *testing.T) {
	objs := catalogObjs(rdrObj{num: 5,
		body: "<< /Filter /Standard /V 4294967300 /R 2 /Length 40 >>"})
	data := buildReaderPDF("%PDF-1.6\n", objs, "/Encrypt 5 0 R\n")
	_, err := OpenBytes(data, nil)
	if !errors.Is(err, ErrUnsupportedSecurityHandler) {
		t.Errorf("err = %v, want ErrUnsupportedSecurityHandler", err)
	}
}

// TestPageRotationNormalisesHugeAndNonFiniteValues is the alert-29 regression:
// /Rotate used to be narrowed to (32-bit) int before the mod-360 normalisation,
// and the Real arm did an unchecked float->int conversion. The arithmetic now
// happens in int64, and non-finite/out-of-int32-range Real values are rejected
// up front rather than converted.
func TestPageRotationNormalisesHugeAndNonFiniteValues(t *testing.T) {
	t.Run("ordinary positive rotation", func(t *testing.T) {
		objs := catalogObjs(rdrObj{num: 3,
			body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] /Rotate 90 >>"})
		d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
		if got := d.PageRotation(1); got != 90 {
			t.Errorf("PageRotation = %d, want 90", got)
		}
	})
	t.Run("negative rotation normalises into range", func(t *testing.T) {
		objs := catalogObjs(rdrObj{num: 3,
			body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] /Rotate -90 >>"})
		d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
		if got := d.PageRotation(1); got != 270 {
			t.Errorf("PageRotation = %d, want 270", got)
		}
	})
	t.Run("huge integer literal still normalises without truncating", func(t *testing.T) {
		objs := catalogObjs(rdrObj{num: 3,
			body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] /Rotate 4294967386 >>"})
		d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
		// 4294967386 mod 360 normalises to 346, which rounds down to the 0
		// quadrant - the same value the "no usable /Rotate" default produces,
		// which is the point: a hostile /Rotate cannot pick an arbitrary
		// quadrant via 32-bit truncation.
		if got := d.PageRotation(1); got != 0 {
			t.Errorf("PageRotation = %d, want 0", got)
		}
	})
	t.Run("NaN Real is rejected, not converted", func(t *testing.T) {
		objs := catalogObjs(rdrObj{num: 3,
			body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] >>"})
		d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
		dict, _, err := d.Page(1)
		if err != nil {
			t.Fatal(err)
		}
		dict.Set("Rotate", Real{Val: math.NaN()})
		if got := d.PageRotation(1); got != 0 {
			t.Errorf("PageRotation = %d, want 0 for a NaN /Rotate", got)
		}
	})
	t.Run("Inf Real is rejected, not converted", func(t *testing.T) {
		objs := catalogObjs(rdrObj{num: 3,
			body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] >>"})
		d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
		dict, _, err := d.Page(1)
		if err != nil {
			t.Fatal(err)
		}
		dict.Set("Rotate", Real{Val: math.Inf(1)})
		if got := d.PageRotation(1); got != 0 {
			t.Errorf("PageRotation = %d, want 0 for an Inf /Rotate", got)
		}
	})
	t.Run("out-of-int32-range Real is rejected, not converted", func(t *testing.T) {
		objs := catalogObjs(rdrObj{num: 3,
			body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] >>"})
		d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", objs, ""))
		dict, _, err := d.Page(1)
		if err != nil {
			t.Fatal(err)
		}
		dict.Set("Rotate", Real{Val: 1e300})
		if got := d.PageRotation(1); got != 0 {
			t.Errorf("PageRotation = %d, want 0 for a 1e300 /Rotate", got)
		}
	})
}

// TestClampKeyLen pins the [0, digest] bound clampKeyLen enforces. clampDictInt
// saturates /Length's magnitude but not its sign, and the /V 4 and /V 5 arms of
// setupEncryption divide it by 8 without the floor /V 2 has, so both ends of the
// range used to reach the Algorithm 2/3 digest slicing unbounded.
func TestClampKeyLen(t *testing.T) {
	const digestLen = 16
	tests := []struct {
		name   string
		keyLen int
		want   int
	}{
		{"zero passes through", 0, 0},
		{"in-range passes through", 5, 5},
		{"digest length passes through", digestLen, digestLen},
		{"negative floors at zero", -1, 0},
		{"far negative floors at zero", math.MinInt32, 0},
		{"one past digest caps", digestLen + 1, digestLen},
		{"far beyond digest caps", 250000000, digestLen},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := clampKeyLen(tc.keyLen, digestLen); got != tc.want {
				t.Errorf("clampKeyLen(%d, %d) = %d, want %d", tc.keyLen, digestLen, got, tc.want)
			}
		})
	}
}

// TestHostileEncryptLengthDoesNotPanic is the end-to-end regression for
// clampKeyLen. /Length is attacker-controlled through pdf.OpenBytes, which
// pades.NativePdfDocumentReader calls on an untrusted document, and a /V 4 or
// /V 5 dictionary carries it into computeRC4Key and computeKeyRev234 unfloored:
// /Length -8 sliced digest[:-1] and /Length 2000000000 sliced digest[:250000000],
// both panicking rather than rejecting the file. Every case must now come back
// as an error.
func TestHostileEncryptLengthDoesNotPanic(t *testing.T) {
	tests := []struct {
		name string
		dict string
	}{
		{"negative length, R4", "<< /Filter /Standard /V 4 /R 4 /Length -8 /P -1 /O <00> /U <00> >>"},
		{"negative length, R3", "<< /Filter /Standard /V 4 /R 3 /Length -8 /P -1 /O <00> /U <00> >>"},
		{"negative length, R2", "<< /Filter /Standard /V 4 /R 2 /Length -8 /P -1 /O <00> /U <00> >>"},
		{"negative length, V5", "<< /Filter /Standard /V 5 /R 6 /Length -8 /P -1 /O <00> /U <00> >>"},
		{"huge length, R4", "<< /Filter /Standard /V 4 /R 4 /Length 2000000000 /P -1 /O <00> /U <00> >>"},
		{"huge length, R3", "<< /Filter /Standard /V 2 /R 3 /Length 2000000000 /P -1 /O <00> /U <00> >>"},
		{"int32-saturated length", "<< /Filter /Standard /V 4 /R 4 /Length 99999999999999999999 /P -1 /O <00> /U <00> >>"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := buildReaderPDF("%PDF-1.6\n", catalogObjs(rdrObj{num: 5, body: tc.dict}),
				"/Encrypt 5 0 R\n/ID [<30313233343536373839616263646566> <30313233343536373839616263646566>]\n")
			if _, err := OpenBytes(data, nil); err == nil {
				t.Errorf("OpenBytes accepted a hostile /Length, want an error")
			}
		})
	}
}
