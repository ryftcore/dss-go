package tsl

import "testing"

func TestCertificatePivotStatusJavaNames(t *testing.T) {
	cases := map[CertificatePivotStatus]string{
		CertificatePivotStatusAdded:      "ADDED",
		CertificatePivotStatusNotChanged: "NOT_CHANGED",
		CertificatePivotStatusRemoved:    "REMOVED",
	}
	for constant, javaName := range cases {
		if string(constant) != javaName {
			t.Fatalf("expected constant value %q to equal Java name %q", constant, javaName)
		}
	}
}
