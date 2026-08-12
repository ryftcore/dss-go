package enumerations

import "testing"

func TestCommitmentTypeEnumFields(t *testing.T) {
	tests := []struct {
		v   CommitmentTypeEnum
		uri string
		oid string
	}{
		{CommitmentTypeEnum_ProofOfOrigin, "http://uri.etsi.org/01903/v1.2.2#ProofOfOrigin", "1.2.840.113549.1.9.16.6.1"},
		{CommitmentTypeEnum_ProofOfReceipt, "http://uri.etsi.org/01903/v1.2.2#ProofOfReceipt", "1.2.840.113549.1.9.16.6.2"},
		{CommitmentTypeEnum_ProofOfDelivery, "http://uri.etsi.org/01903/v1.2.2#ProofOfDelivery", "1.2.840.113549.1.9.16.6.3"},
		{CommitmentTypeEnum_ProofOfSender, "http://uri.etsi.org/01903/v1.2.2#ProofOfSender", "1.2.840.113549.1.9.16.6.4"},
		{CommitmentTypeEnum_ProofOfApproval, "http://uri.etsi.org/01903/v1.2.2#ProofOfApproval", "1.2.840.113549.1.9.16.6.5"},
		{CommitmentTypeEnum_ProofOfCreation, "http://uri.etsi.org/01903/v1.2.2#ProofOfCreation", "1.2.840.113549.1.9.16.6.6"},
	}
	for _, tt := range tests {
		if got := tt.v.URI(); got != tt.uri {
			t.Errorf("%v.URI() = %q, want %q", tt.v, got, tt.uri)
		}
		if got := tt.v.OID(); got != tt.oid {
			t.Errorf("%v.OID() = %q, want %q", tt.v, got, tt.oid)
		}
		if got := tt.v.Description(); got != string(tt.v) {
			t.Errorf("%v.Description() = %q, want %q", tt.v, got, string(tt.v))
		}
		if got := tt.v.Qualifier(); got != "" {
			t.Errorf("%v.Qualifier() = %q, want empty", tt.v, got)
		}
		if got := tt.v.DocumentationReferences(); got != nil {
			t.Errorf("%v.DocumentationReferences() = %v, want nil", tt.v, got)
		}
	}
}

func TestCommitmentTypeEnumValueOf(t *testing.T) {
	for _, v := range CommitmentTypeEnumValues() {
		got, err := CommitmentTypeEnumValueOf(string(v))
		if err != nil {
			t.Errorf("CommitmentTypeEnumValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("CommitmentTypeEnumValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestCommitmentTypeEnumValueOfUnknown(t *testing.T) {
	if _, err := CommitmentTypeEnumValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}

// Compile-time assertion that CommitmentTypeEnum implements CommitmentType.
var _ CommitmentType = CommitmentTypeEnum_ProofOfOrigin
