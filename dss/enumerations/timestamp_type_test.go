package enumerations

import "testing"

func TestTimestampTypeOrderPredicates(t *testing.T) {
	tests := []struct {
		v          TimestampType
		content    bool
		signature  bool
		validation bool
		container  bool
		document   bool
		archival   bool
		evidence   bool
		covers     bool
	}{
		{TimestampTypeContentTimestamp, true, false, false, false, false, false, false, false},
		{TimestampTypeAllDataObjectsTimestamp, true, false, false, false, false, false, false, false},
		{TimestampTypeIndividualDataObjectsTimestamp, true, false, false, false, false, false, false, false},
		{TimestampTypeSignatureTimestamp, false, true, false, false, false, false, false, true},
		{TimestampTypeVRITimestamp, false, true, false, false, false, false, false, true},
		{TimestampTypeValidationDataRefsOnlyTimestamp, false, false, true, false, false, false, false, false},
		{TimestampTypeValidationDataTimestamp, false, false, true, false, false, false, false, true},
		{TimestampTypeContainerTimestamp, false, false, false, true, false, false, false, true},
		{TimestampTypeDocumentTimestamp, false, false, false, false, true, false, false, true},
		{TimestampTypeArchiveTimestamp, false, false, false, false, false, true, false, true},
		{TimestampTypeEvidenceRecordTimestamp, false, false, false, false, false, false, true, true},
	}
	for _, tt := range tests {
		if got := tt.v.IsContentTimestamp(); got != tt.content {
			t.Errorf("%v.IsContentTimestamp() = %v, want %v", tt.v, got, tt.content)
		}
		if got := tt.v.IsSignatureTimestamp(); got != tt.signature {
			t.Errorf("%v.IsSignatureTimestamp() = %v, want %v", tt.v, got, tt.signature)
		}
		if got := tt.v.IsValidationDataTimestamp(); got != tt.validation {
			t.Errorf("%v.IsValidationDataTimestamp() = %v, want %v", tt.v, got, tt.validation)
		}
		if got := tt.v.IsContainerTimestamp(); got != tt.container {
			t.Errorf("%v.IsContainerTimestamp() = %v, want %v", tt.v, got, tt.container)
		}
		if got := tt.v.IsDocumentTimestamp(); got != tt.document {
			t.Errorf("%v.IsDocumentTimestamp() = %v, want %v", tt.v, got, tt.document)
		}
		if got := tt.v.IsArchivalTimestamp(); got != tt.archival {
			t.Errorf("%v.IsArchivalTimestamp() = %v, want %v", tt.v, got, tt.archival)
		}
		if got := tt.v.IsEvidenceRecordTimestamp(); got != tt.evidence {
			t.Errorf("%v.IsEvidenceRecordTimestamp() = %v, want %v", tt.v, got, tt.evidence)
		}
		if got := tt.v.CoversSignature(); got != tt.covers {
			t.Errorf("%v.CoversSignature() = %v, want %v", tt.v, got, tt.covers)
		}
	}
}

func TestTimestampTypeCompare(t *testing.T) {
	if TimestampTypeContentTimestamp.Compare(TimestampTypeSignatureTimestamp) >= 0 {
		t.Error("CONTENT_TIMESTAMP should compare before SIGNATURE_TIMESTAMP")
	}
	if TimestampTypeArchiveTimestamp.Compare(TimestampTypeValidationDataTimestamp) <= 0 {
		t.Error("ARCHIVE_TIMESTAMP should compare after VALIDATION_DATA_TIMESTAMP")
	}
	if TimestampTypeSignatureTimestamp.Compare(TimestampTypeVRITimestamp) != 0 {
		t.Error("SIGNATURE_TIMESTAMP and VRI_TIMESTAMP share the same order")
	}
}

func TestTimestampTypeValueOf(t *testing.T) {
	for _, v := range TimestampTypeValues() {
		got, err := TimestampTypeValueOf(string(v))
		if err != nil {
			t.Errorf("TimestampTypeValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("TimestampTypeValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestTimestampTypeValueOfUnknown(t *testing.T) {
	if _, err := TimestampTypeValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
