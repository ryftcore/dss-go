package tsl

import "testing"

func TestTrustServiceProviderRoundTrip(t *testing.T) {
	tsp := NewTrustServiceProvider()
	tsp.SetNames(map[string][]string{"en": {"Name"}})
	tsp.SetTradeNames(map[string][]string{"en": {"Trade"}})
	tsp.SetRegistrationIdentifiers([]string{"REG1"})
	tsp.SetPostalAddresses(map[string]string{"en": "123 Street"})
	tsp.SetElectronicAddresses(map[string][]string{"en": {"mailto:a@b.c"}})
	tsp.SetInformation(map[string]string{"en": "info"})
	tsp.SetServices([]*TrustService{NewTrustService(nil, nil)})
	tsp.SetTerritory("FR")

	if tsp.Names()["en"][0] != "Name" {
		t.Fatalf("unexpected Names: %v", tsp.Names())
	}
	if tsp.TradeNames()["en"][0] != "Trade" {
		t.Fatalf("unexpected TradeNames: %v", tsp.TradeNames())
	}
	if len(tsp.RegistrationIdentifiers()) != 1 || tsp.RegistrationIdentifiers()[0] != "REG1" {
		t.Fatalf("unexpected RegistrationIdentifiers: %v", tsp.RegistrationIdentifiers())
	}
	if tsp.PostalAddresses()["en"] != "123 Street" {
		t.Fatalf("unexpected PostalAddresses: %v", tsp.PostalAddresses())
	}
	if tsp.ElectronicAddresses()["en"][0] != "mailto:a@b.c" {
		t.Fatalf("unexpected ElectronicAddresses: %v", tsp.ElectronicAddresses())
	}
	if tsp.Information()["en"] != "info" {
		t.Fatalf("unexpected Information: %v", tsp.Information())
	}
	if len(tsp.Services()) != 1 {
		t.Fatalf("unexpected Services: %v", tsp.Services())
	}
	if tsp.Territory() != "FR" {
		t.Fatalf("unexpected Territory: %s", tsp.Territory())
	}
}
