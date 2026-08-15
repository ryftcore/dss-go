package validationreport

import "testing"

// TestEnumReexports checks the alias re-exports (object_type.go,
// constraint_status.go, type_of_proof.go, signature_validation_process_id.go,
// uri_based_enum_parser.go - see doc.go's "Enum and parser re-exports") stay
// wired to jaxb's concrete implementations: every constant round-trips
// through this package's Parse functions exactly as it does through jaxb's.
func TestEnumReexports(t *testing.T) {
	for _, v := range ObjectTypeValues() {
		if got := ParseObjectType(Print(v)); got != v {
			t.Errorf("ParseObjectType(Print(%v)) = %v", v, got)
		}
	}
	for _, v := range ConstraintStatusValues() {
		if got := ParseConstraintStatus(Print(v)); got != v {
			t.Errorf("ParseConstraintStatus(Print(%v)) = %v", v, got)
		}
	}
	for _, v := range TypeOfProofValues() {
		if got := ParseTypeOfProof(Print(v)); got != v {
			t.Errorf("ParseTypeOfProof(Print(%v)) = %v", v, got)
		}
	}
	for _, v := range SignatureValidationProcessIDValues() {
		if got := ParseSignatureValidationProcessID(Print(v)); got != v {
			t.Errorf("ParseSignatureValidationProcessID(Print(%v)) = %v", v, got)
		}
	}

	if _, err := ObjectTypeValueOf("CERTIFICATE"); err != nil {
		t.Errorf("ObjectTypeValueOf(CERTIFICATE): %v", err)
	}
	if _, err := ObjectTypeValueOf("NOT_A_VALUE"); err == nil {
		t.Error("ObjectTypeValueOf(NOT_A_VALUE): expected error")
	}
}
