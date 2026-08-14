package lote

import "testing"

func TestTrustedEntityRoundTrip(t *testing.T) {
	e := NewTrustedEntity()
	e.SetNames(map[string][]string{"en": {"Example CA"}})
	e.SetTradeNames(map[string][]string{"en": {"Example"}})
	e.SetRegistrationIdentifiers([]string{"VATBE-1234"})
	e.SetPostalAddresses(map[string]string{"en": "1 Rue Example"})
	e.SetElectronicAddresses(map[string][]string{"en": {"mailto:test@example.org"}})
	e.SetInformation(map[string][]string{"en": {"info"}})
	e.SetTerritory("BE")

	service := NewTrustedEntityService(nil, nil)
	e.SetServices([]*TrustedEntityService{service})

	if e.Territory() != "BE" {
		t.Fatalf("unexpected Territory(): %s", e.Territory())
	}
	if len(e.Services()) != 1 || e.Services()[0] != service {
		t.Fatalf("unexpected Services(): %v", e.Services())
	}
	if len(e.RegistrationIdentifiers()) != 1 || e.RegistrationIdentifiers()[0] != "VATBE-1234" {
		t.Fatalf("unexpected RegistrationIdentifiers(): %v", e.RegistrationIdentifiers())
	}
}

func TestTrustedEntityServiceBuilderRoundTrip(t *testing.T) {
	b := NewTrustEntityServiceBuilder().
		SetCertificates(nil).
		SetStatusAndInformationExtensions(nil)
	service := b.Build()
	if service == nil {
		t.Fatal("expected Build() to return a non-nil TrustedEntityService")
	}
}
