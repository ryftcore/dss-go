package pades

import "testing"

// TestByteRangeValidateRejectsMalformedRanges pins every rejection branch of ByteRange.Validate,
// the port of eu.europa.esig.dss.pades.validation.ByteRange#validate.
//
// It exists because a mutation check found the branch that matters most here to be completely
// unguarded: disabling the "second hash part must start after the first hash part" comparison -
// the check that rejects a /ByteRange whose two covered spans OVERLAP, i.e. the byte-range
// overlap attack - left the whole test suite green, including the
// validation/pdf-byterange-overlap.pdf fixture in the upstream cross-validation corpus. That
// fixture reaches its expected "not intact" verdict through the CMS message-digest comparison
// instead, so it can never fail if this check regresses. A defence-in-depth check with no test
// of its own is a check that silently rots, so each condition gets a direct case here.
func TestByteRangeValidateRejectsMalformedRanges(t *testing.T) {
	for _, testCase := range []struct {
		name    string
		input   []int
		wantErr string
	}{
		{
			name:  "well-formed range validates",
			input: []int{0, 840, 960, 240},
		},
		{
			name:  "adjacent parts are allowed (empty gap)",
			input: []int{0, 840, 840, 240},
		},
		{
			name:    "nil array",
			input:   nil,
			wantErr: "Incorrect ByteRange size",
		},
		{
			name:    "wrong arity",
			input:   []int{0, 840, 960},
			wantErr: "Incorrect ByteRange size",
		},
		{
			name:    "does not start at the beginning of the file",
			input:   []int{1, 840, 960, 240},
			wantErr: "The ByteRange must cover start of file",
		},
		{
			name:    "negative first part",
			input:   []int{0, -1, 960, 240},
			wantErr: "The first hash part doesn't cover anything",
		},
		{
			// The byte-range overlap attack: the second covered span starts INSIDE the first,
			// so the bytes in between are covered twice and the bytes the /Contents placeholder
			// actually occupies are never excluded.
			name:    "second part overlaps the first",
			input:   []int{0, 840, 839, 240},
			wantErr: "The second hash part must start after the first hash part",
		},
		{
			name:    "second part starts before the file",
			input:   []int{0, 0, -1, 240},
			wantErr: "The second hash part must start after the first hash part",
		},
		{
			name:    "negative second part",
			input:   []int{0, 840, 960, -1},
			wantErr: "The second hash part doesn't cover anything",
		},
	} {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			err := NewByteRange(testCase.input).Validate()
			switch {
			case testCase.wantErr == "" && err != nil:
				t.Fatalf("Validate() = %v, want no error", err)
			case testCase.wantErr != "" && err == nil:
				t.Fatalf("Validate() = nil, want %q", testCase.wantErr)
			case testCase.wantErr != "" && err.Error() != testCase.wantErr:
				t.Fatalf("Validate() = %q, want %q", err.Error(), testCase.wantErr)
			}
		})
	}
}
