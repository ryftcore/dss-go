package lote

import (
	"testing"
	"time"
)

func TestNewTrustedEntityServiceStatusAndInformationExtensionsPanicsOnNilBuilder(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nil builder")
		}
		if r != "ServiceStatusAndInformationExtensionsBuilder cannot be null!" {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()
	NewTrustedEntityServiceStatusAndInformationExtensions(nil)
}

func TestServiceStatusAndInformationExtensionsBuilderRoundTrip(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	status := NewServiceStatusAndInformationExtensionsBuilder().
		SetNames(map[string][]string{"en": {"granted"}}).
		SetType("http://uri.etsi.org/TrstSvc/Svctype/CA/QC").
		SetStatus("http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted").
		SetServiceSupplyPoints([]string{"https://example.org/crl"}).
		SetStartDate(start).
		SetEndDate(end).
		Build()

	if status.Type() != "http://uri.etsi.org/TrstSvc/Svctype/CA/QC" {
		t.Fatalf("unexpected Type(): %s", status.Type())
	}
	if !status.StartDate().Equal(start) || !status.EndDate().Equal(end) {
		t.Fatalf("unexpected StartDate()/EndDate(): %v / %v", status.StartDate(), status.EndDate())
	}

	// A ServiceStatusAndInformationExtensions interface value must be satisfiable.
	var iface ServiceStatusAndInformationExtensions = status
	if iface.Status() != status.Status() {
		t.Fatal("expected interface Status() to match the concrete value")
	}

	rebuilt := NewServiceStatusAndInformationExtensionsBuilderFrom(status).Build()
	if rebuilt.Type() != status.Type() || rebuilt.Status() != status.Status() {
		t.Fatal("expected the copy-from builder to preserve fields")
	}
}
