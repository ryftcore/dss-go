package tsl

import "testing"

func TestCertificatePivotStatusJavaNames(t *testing.T) {
	cases := map[CertificatePivotStatus]string{
		CertificatePivotStatus_ADDED:       "ADDED",
		CertificatePivotStatus_NOT_CHANGED: "NOT_CHANGED",
		CertificatePivotStatus_REMOVED:     "REMOVED",
	}
	for constant, javaName := range cases {
		if string(constant) != javaName {
			t.Fatalf("expected constant value %q to equal Java name %q", constant, javaName)
		}
	}
}
