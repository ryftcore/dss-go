package enumerations

import "testing"

func TestSubIndicationURI(t *testing.T) {
	tests := []struct {
		v    SubIndication
		want string
	}{
		{SubIndicationFormatFailure, "urn:etsi:019102:subindication:FORMAT_FAILURE"},
		{SubIndicationHashFailure, "urn:etsi:019102:subindication:HASH_FAILURE"},
		{SubIndicationSigCryptoFailure, "urn:etsi:019102:subindication:SIG_CRYPTO_FAILURE"},
		{SubIndicationRevoked, "urn:etsi:019102:subindication:REVOKED"},
		{SubIndicationExpired, "urn:etsi:019102:subindication:EXPIRED"},
		{SubIndicationNotYetValid, "urn:etsi:019102:subindication:NOT_YET_VALID"},
		{SubIndicationSigConstraintsFailure, "urn:etsi:019102:subindication:SIG_CONSTRAINTS_FAILURE"},
		{SubIndicationChainConstraintsFailure, "urn:etsi:019102:subindication:CHAIN_CONSTRAINTS_FAILURE"},
		{SubIndicationCertificateChainGeneralFailure, "urn:etsi:019102:subindication:CERTIFICATE_CHAIN_GENERAL_FAILURE"},
		{SubIndicationCryptoConstraintsFailure, "urn:etsi:019102:subindication:CRYPTO_CONSTRAINTS_FAILURE"},
		{SubIndicationPolicyProcessingError, "urn:etsi:019102:subindication:POLICY_PROCESSING_ERROR"},
		{SubIndicationSignaturePolicyNotAvailable, "urn:etsi:019102:subindication:SIGNATURE_POLICY_NOT_AVAILABLE"},
		{SubIndicationTimestampOrderFailure, "urn:etsi:019102:subindication:TIMESTAMP_ORDER_FAILURE"},
		{SubIndicationNoSigningCertificateFound, "urn:etsi:019102:subindication:NO_SIGNING_CERTIFICATE_FOUND"},
		{SubIndicationNoCertificateChainFound, "urn:etsi:019102:subindication:NO_CERTIFICATE_CHAIN_FOUND"},
		{SubIndicationNoCertificateChainFoundNoPOE, "urn:etsi:019102:subindication:NO_CERTIFICATE_CHAIN_FOUND_NO_POE"},
		{SubIndicationRevokedNoPOE, "urn:etsi:019102:subindication:REVOKED_NO_POE"},
		{SubIndicationRevokedCANoPOE, "urn:etsi:019102:subindication:REVOKED_CA_NO_POE"},
		{SubIndicationOutOfBoundsNotRevoked, "urn:etsi:019102:subindication:OUT_OF_BOUNDS_NOT_REVOKED"},
		{SubIndicationOutOfBoundsNoPOE, "urn:etsi:019102:subindication:OUT_OF_BOUNDS_NO_POE"},
		{SubIndicationRevocationOutOfBoundsNoPOE, "urn:etsi:019102:subindication:REVOCATION_OUT_OF_BOUNDS_NO_POE"},
		{SubIndicationCryptoConstraintsFailureNoPOE, "urn:etsi:019102:subindication:CRYPTO_CONSTRAINTS_FAILURE_NO_POE"},
		{SubIndicationNoPOE, "urn:etsi:019102:subindication:NO_POE"},
		{SubIndicationTryLater, "urn:etsi:019102:subindication:TRY_LATER"},
		{SubIndicationSignedDataNotFound, "urn:etsi:019102:subindication:SIGNED_DATA_NOT_FOUND"},
		{SubIndicationEAAConstraintsFailure, "urn:cef:dss:subindication:EAA_CONSTRAINTS_FAILURE"},
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
	if err != nil || got != SubIndicationRevoked {
		t.Errorf("SubIndicationForName(REVOKED) = (%q, %v), want (%q, nil)", got, err, SubIndicationRevoked)
	}
	if _, err := SubIndicationForName("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
