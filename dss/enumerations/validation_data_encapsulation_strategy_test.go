package enumerations

import "testing"

func TestValidationDataEncapsulationStrategyValues(t *testing.T) {
	want := []ValidationDataEncapsulationStrategy{
		ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData,
		ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataLTSeparated,
		ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData,
		ValidationDataEncapsulationStrategyCertificateRevocationValuesAndAnyValidationData,
		ValidationDataEncapsulationStrategyAnyValidationDataOnly,
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
