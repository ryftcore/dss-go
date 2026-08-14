package enumerations

import "testing"

func TestValidationDataEncapsulationStrategyValues(t *testing.T) {
	want := []ValidationDataEncapsulationStrategy{
		ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA,
		ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA_LT_SEPARATED,
		ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA_AND_ANY_VALIDATION_DATA,
		ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_ANY_VALIDATION_DATA,
		ValidationDataEncapsulationStrategy_ANY_VALIDATION_DATA_ONLY,
	}
	got := ValidationDataEncapsulationStrategyValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
}
