package enumerations

import "testing"

func TestExtendedKeyUsageFields(t *testing.T) {
	tests := []struct {
		v           ExtendedKeyUsage
		description string
		oid         string
	}{
		{ExtendedKeyUsageServerAuth, "serverAuth", "1.3.6.1.5.5.7.3.1"},
		{ExtendedKeyUsageClientAuth, "clientAuth", "1.3.6.1.5.5.7.3.2"},
		{ExtendedKeyUsageCodeSigning, "codeSigning", "1.3.6.1.5.5.7.3.3"},
		{ExtendedKeyUsageEmailProtection, "emailProtection", "1.3.6.1.5.5.7.3.4"},
		{ExtendedKeyUsageTimestamping, "timeStamping", "1.3.6.1.5.5.7.3.8"},
		{ExtendedKeyUsageOCSPSigning, "ocspSigning", "1.3.6.1.5.5.7.3.9"},
		{ExtendedKeyUsageTSLSigning, "tslSigning", "0.4.0.2231.3.0"},
		{ExtendedKeyUsageTSLBinding, "tslBinding", "0.4.0.194115.1.0"},
	}
	for _, tt := range tests {
		if got := tt.v.OID(); got != tt.oid {
			t.Errorf("%v.OID() = %q, want %q", tt.v, got, tt.oid)
		}
		if got := tt.v.Description(); got != tt.description {
			t.Errorf("%v.Description() = %q, want %q", tt.v, got, tt.description)
		}
	}
}

func TestExtendedKeyUsageValueOf(t *testing.T) {
	for _, v := range ExtendedKeyUsageValues() {
		got, err := ExtendedKeyUsageValueOf(string(v))
		if err != nil {
			t.Errorf("ExtendedKeyUsageValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("ExtendedKeyUsageValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestExtendedKeyUsageValueOfUnknown(t *testing.T) {
	if _, err := ExtendedKeyUsageValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
