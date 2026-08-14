package enumerations

import "testing"

func TestSubIndicationURI(t *testing.T) {
	tests := []struct {
		v    SubIndication
		want string
	}{
		{SubIndication_FORMAT_FAILURE, "urn:etsi:019102:subindication:FORMAT_FAILURE"},
		{SubIndication_HASH_FAILURE, "urn:etsi:019102:subindication:HASH_FAILURE"},
		{SubIndication_SIG_CRYPTO_FAILURE, "urn:etsi:019102:subindication:SIG_CRYPTO_FAILURE"},
		{SubIndication_REVOKED, "urn:etsi:019102:subindication:REVOKED"},
		{SubIndication_EXPIRED, "urn:etsi:019102:subindication:EXPIRED"},
		{SubIndication_NOT_YET_VALID, "urn:etsi:019102:subindication:NOT_YET_VALID"},
		{SubIndication_SIG_CONSTRAINTS_FAILURE, "urn:etsi:019102:subindication:SIG_CONSTRAINTS_FAILURE"},
		{SubIndication_CHAIN_CONSTRAINTS_FAILURE, "urn:etsi:019102:subindication:CHAIN_CONSTRAINTS_FAILURE"},
		{SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE, "urn:etsi:019102:subindication:CERTIFICATE_CHAIN_GENERAL_FAILURE"},
		{SubIndication_CRYPTO_CONSTRAINTS_FAILURE, "urn:etsi:019102:subindication:CRYPTO_CONSTRAINTS_FAILURE"},
		{SubIndication_POLICY_PROCESSING_ERROR, "urn:etsi:019102:subindication:POLICY_PROCESSING_ERROR"},
		{SubIndication_SIGNATURE_POLICY_NOT_AVAILABLE, "urn:etsi:019102:subindication:SIGNATURE_POLICY_NOT_AVAILABLE"},
		{SubIndication_TIMESTAMP_ORDER_FAILURE, "urn:etsi:019102:subindication:TIMESTAMP_ORDER_FAILURE"},
		{SubIndication_NO_SIGNING_CERTIFICATE_FOUND, "urn:etsi:019102:subindication:NO_SIGNING_CERTIFICATE_FOUND"},
		{SubIndication_NO_CERTIFICATE_CHAIN_FOUND, "urn:etsi:019102:subindication:NO_CERTIFICATE_CHAIN_FOUND"},
		{SubIndication_NO_CERTIFICATE_CHAIN_FOUND_NO_POE, "urn:etsi:019102:subindication:NO_CERTIFICATE_CHAIN_FOUND_NO_POE"},
		{SubIndication_REVOKED_NO_POE, "urn:etsi:019102:subindication:REVOKED_NO_POE"},
		{SubIndication_REVOKED_CA_NO_POE, "urn:etsi:019102:subindication:REVOKED_CA_NO_POE"},
		{SubIndication_OUT_OF_BOUNDS_NOT_REVOKED, "urn:etsi:019102:subindication:OUT_OF_BOUNDS_NOT_REVOKED"},
		{SubIndication_OUT_OF_BOUNDS_NO_POE, "urn:etsi:019102:subindication:OUT_OF_BOUNDS_NO_POE"},
		{SubIndication_REVOCATION_OUT_OF_BOUNDS_NO_POE, "urn:etsi:019102:subindication:REVOCATION_OUT_OF_BOUNDS_NO_POE"},
		{SubIndication_CRYPTO_CONSTRAINTS_FAILURE_NO_POE, "urn:etsi:019102:subindication:CRYPTO_CONSTRAINTS_FAILURE_NO_POE"},
		{SubIndication_NO_POE, "urn:etsi:019102:subindication:NO_POE"},
		{SubIndication_TRY_LATER, "urn:etsi:019102:subindication:TRY_LATER"},
		{SubIndication_SIGNED_DATA_NOT_FOUND, "urn:etsi:019102:subindication:SIGNED_DATA_NOT_FOUND"},
		{SubIndication_EAA_CONSTRAINTS_FAILURE, "urn:cef:dss:subindication:EAA_CONSTRAINTS_FAILURE"},
	}
	if len(tests) != len(SubIndicationValues()) {
		t.Fatalf("test table has %d entries, want %d (one per constant)", len(tests), len(SubIndicationValues()))
	}
	for _, tt := range tests {
		if got := tt.v.URI(); got != tt.want {
			t.Errorf("%v.URI() = %q, want %q", tt.v, got, tt.want)
		}
	}
}

func TestSubIndicationValueOf(t *testing.T) {
	for _, v := range SubIndicationValues() {
		got, err := SubIndicationValueOf(string(v))
		if err != nil {
			t.Errorf("SubIndicationValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("SubIndicationValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestSubIndicationForName(t *testing.T) {
	got, err := SubIndicationForName("")
	if err != nil || got != "" {
		t.Errorf("SubIndicationForName(\"\") = (%q, %v), want (\"\", nil)", got, err)
	}
	got, err = SubIndicationForName("REVOKED")
	if err != nil || got != SubIndication_REVOKED {
		t.Errorf("SubIndicationForName(REVOKED) = (%q, %v), want (%q, nil)", got, err, SubIndication_REVOKED)
	}
	if _, err := SubIndicationForName("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
