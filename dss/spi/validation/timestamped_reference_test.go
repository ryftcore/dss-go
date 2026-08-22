package validation

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// TestTimestampedReference covers the value semantics the timestamp source relies on: two
// references are the same reference when their id and category agree.
func TestTimestampedReference(t *testing.T) {
	reference := NewTimestampedReference("S-1234", enumerations.TimestampedObjectTypeSignature)

	if got := reference.ObjectId(); got != "S-1234" {
		t.Errorf("ObjectId() = %q, want S-1234", got)
	}
	if got := reference.Category(); got != enumerations.TimestampedObjectTypeSignature {
		t.Errorf("Category() = %s, want SIGNATURE", got)
	}
	if got, want := reference.String(),
		"TimestampedReference with Id [S-1234] and type [SIGNATURE]"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	cases := map[string]struct {
		other *TimestampedReference
		equal bool
	}{
		"same":          {NewTimestampedReference("S-1234", enumerations.TimestampedObjectTypeSignature), true},
		"other id":      {NewTimestampedReference("S-5678", enumerations.TimestampedObjectTypeSignature), false},
		"other type":    {NewTimestampedReference("S-1234", enumerations.TimestampedObjectTypeCertificate), false},
		"nil reference": {nil, false},
	}
	for name, testCase := range cases {
		if got := reference.Equals(testCase.other); got != testCase.equal {
			t.Errorf("Equals(%s) = %t, want %t", name, got, testCase.equal)
		}
	}
	if !reference.Equals(reference) {
		t.Error("Equals(self) = false, want true")
	}
}

// TestTimestampInclude covers the XAdES Include value object.
func TestTimestampInclude(t *testing.T) {
	empty := NewTimestampInclude()
	if empty.URI() != "" || empty.IsReferencedData() {
		t.Errorf("NewTimestampInclude() = {%q, %t}, want the zero value", empty.URI(), empty.IsReferencedData())
	}
	empty.SetURI("#r-id-1")
	empty.SetReferencedData(true)
	if empty.URI() != "#r-id-1" || !empty.IsReferencedData() {
		t.Errorf("the setters did not take: {%q, %t}", empty.URI(), empty.IsReferencedData())
	}

	include := NewTimestampIncludeWithURI("#r-id-2", false)
	if include.URI() != "#r-id-2" || include.IsReferencedData() {
		t.Errorf("NewTimestampIncludeWithURI() = {%q, %t}, want {#r-id-2, false}",
			include.URI(), include.IsReferencedData())
	}
}

// TestArchiveTimestampHashIndexStatus covers the ats-hash-index status, whose error list is
// created on demand.
func TestArchiveTimestampHashIndexStatus(t *testing.T) {
	status := NewArchiveTimestampHashIndexStatus()

	if got := status.Version(); got != "" {
		t.Errorf("Version() = %s, want the unset value", got)
	}
	if got := status.ErrorMessages(); got == nil || len(got) != 0 {
		t.Errorf("ErrorMessages() = %v, want an empty, non-nil list", got)
	}

	status.SetVersion(enumerations.ArchiveTimestampHashIndexVersionATSHashIndexV2)
	status.AddErrorMessage("first")
	status.AddErrorMessage("second")

	if got := status.Version(); got != enumerations.ArchiveTimestampHashIndexVersionATSHashIndexV2 {
		t.Errorf("Version() = %s, want ATS_HASH_INDEX_V2", got)
	}
	if got := status.ErrorMessages(); len(got) != 2 || got[0] != "first" || got[1] != "second" {
		t.Errorf("ErrorMessages() = %v, want [first second]", got)
	}
}
