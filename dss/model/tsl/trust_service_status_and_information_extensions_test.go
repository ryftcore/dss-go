package tsl

import (
	"testing"
	"time"
)

func TestTrustServiceStatusAndInformationExtensionsBuilderRoundTrip(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	expired := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	cfq := NewConditionForQualifiers(nil, []string{"q"})

	status := NewTrustServiceStatusAndInformationExtensionsBuilder().
		SetNames(map[string][]string{"en": {"Service"}}).
		SetType("http://uri.etsi.org/TrstSvc/Svctype/CA/QC").
		SetStatus("http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted").
		SetConditionsForQualifiers([]*ConditionForQualifiers{cfq}).
		SetAdditionalServiceInfoUris([]string{"http://uri.etsi.org/foo"}).
		SetServiceSupplyPoints([]string{"http://example.org/point"}).
		SetExpiredCertsRevocationInfo(expired).
		SetStartDate(start).
		SetEndDate(end).
		Build()

	if status.Type() != "http://uri.etsi.org/TrstSvc/Svctype/CA/QC" {
		t.Fatalf("unexpected Type: %s", status.Type())
	}
	if status.Status() != "http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted" {
		t.Fatalf("unexpected Status: %s", status.Status())
	}
	if status.Names()["en"][0] != "Service" {
		t.Fatalf("unexpected Names: %v", status.Names())
	}
	if len(status.ConditionsForQualifiers()) != 1 || status.ConditionsForQualifiers()[0] != cfq {
		t.Fatalf("unexpected ConditionsForQualifiers: %v", status.ConditionsForQualifiers())
	}
	if len(status.AdditionalServiceInfoUris()) != 1 {
		t.Fatalf("unexpected AdditionalServiceInfoUris: %v", status.AdditionalServiceInfoUris())
	}
	if len(status.ServiceSupplyPoints()) != 1 {
		t.Fatalf("unexpected ServiceSupplyPoints: %v", status.ServiceSupplyPoints())
	}
	if !status.ExpiredCertsRevocationInfo().Equal(expired) {
		t.Fatalf("unexpected ExpiredCertsRevocationInfo: %v", status.ExpiredCertsRevocationInfo())
	}
	if !status.StartDate().Equal(start) || !status.EndDate().Equal(end) {
		t.Fatalf("unexpected time range: start=%v end=%v", status.StartDate(), status.EndDate())
	}
}

func TestTrustServiceStatusAndInformationExtensionsBuilderFromExisting(t *testing.T) {
	original := NewTrustServiceStatusAndInformationExtensionsBuilder().
		SetType("type").
		SetStatus("status").
		Build()

	copyBuilder := NewTrustServiceStatusAndInformationExtensionsBuilderFrom(original)
	copied := copyBuilder.Build()

	if !original.Equals(copied) {
		t.Fatalf("expected the copy built from an existing status to be equal to the original")
	}
}

func TestTrustServiceStatusAndInformationExtensionsBuildPanicsOnNilBuilder(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected a panic for a nil builder")
		}
	}()
	NewTrustServiceStatusAndInformationExtensions(nil)
}

func TestTrustServiceStatusAndInformationExtensionsEquals(t *testing.T) {
	a := NewTrustServiceStatusAndInformationExtensionsBuilder().SetType("t").Build()
	b := NewTrustServiceStatusAndInformationExtensionsBuilder().SetType("t").Build()
	if !a.Equals(b) {
		t.Fatalf("expected equal values")
	}
	c := NewTrustServiceStatusAndInformationExtensionsBuilder().SetType("other").Build()
	if a.Equals(c) {
		t.Fatalf("expected unequal values")
	}
	if a.Equals(nil) {
		t.Fatalf("expected Equals(nil) to be false")
	}
}
