package tsl

import (
	"testing"
	"time"
)

func TestCertificateTrustTimeNotTrusted(t *testing.T) {
	ctt := NewCertificateTrustTime(false)
	if ctt.IsTrusted() {
		t.Fatalf("expected not trusted")
	}
	if ctt.IsTrustedAtTime(time.Now()) {
		t.Fatalf("a not-trusted entry must never be trusted at any time")
	}
}

func TestCertificateTrustTimeIndefinitelyTrusted(t *testing.T) {
	ctt := NewCertificateTrustTime(true)
	if !ctt.IsTrusted() {
		t.Fatalf("expected trusted")
	}
	// Zero start/end date stand in for Java's null (unbounded) range.
	if !ctt.IsTrustedAtTime(time.Now()) {
		t.Fatalf("an indefinitely trusted entry must be trusted at any time")
	}
}

func TestCertificateTrustTimeIsTrustedAtTime(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	ctt := NewCertificateTrustTimeWithRange(start, end)

	if !ctt.IsTrustedAtTime(start) {
		t.Fatalf("expected trusted at the range start (inclusive via dateBefore/dateAfter equality)")
	}
	if !ctt.IsTrustedAtTime(end) {
		t.Fatalf("expected trusted at the range end (inclusive via dateBefore/dateAfter equality)")
	}
	if !ctt.IsTrustedAtTime(start.Add(time.Hour)) {
		t.Fatalf("expected trusted within the range")
	}
	if ctt.IsTrustedAtTime(start.Add(-time.Hour)) {
		t.Fatalf("expected not trusted before the range")
	}
	if ctt.IsTrustedAtTime(end.Add(time.Hour)) {
		t.Fatalf("expected not trusted after the range")
	}
}

func TestCertificateTrustTimeJointTrustTime(t *testing.T) {
	start := time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC)
	ctt := NewCertificateTrustTimeWithRange(start, end)

	otherStart := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	otherEnd := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	joint := ctt.JointTrustTime(otherStart, otherEnd)

	if !joint.StartDate().Equal(otherStart) {
		t.Fatalf("expected joint start to widen to the earlier date, got %v", joint.StartDate())
	}
	if !joint.EndDate().Equal(otherEnd) {
		t.Fatalf("expected joint end to widen to the later date, got %v", joint.EndDate())
	}
	// The original entry must not be mutated.
	if !ctt.StartDate().Equal(start) || !ctt.EndDate().Equal(end) {
		t.Fatalf("expected original CertificateTrustTime to be unchanged")
	}
}

func TestCertificateTrustTimeEquals(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	a := NewCertificateTrustTimeWithRange(start, end)
	b := NewCertificateTrustTimeWithRange(start, end)

	if !a.Equals(b) {
		t.Fatalf("expected equal CertificateTrustTime values")
	}
	if a.Equals(nil) {
		t.Fatalf("expected Equals(nil) to be false")
	}
	if a.Equals(NewCertificateTrustTime(false)) {
		t.Fatalf("expected unequal CertificateTrustTime values")
	}
}
