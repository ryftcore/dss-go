package tsl

import "testing"

func TestOtherTSLPointerBuilderRoundTrip(t *testing.T) {
	mra := NewMRA()
	p := NewOtherTSLPointerBuilder().
		SetTslLocation("https://example.org/tl.xml").
		SetSchemeTerritory("FR").
		SetTslType("http://uri.etsi.org/TrstSvc/TrustedList/TSLType/generic").
		SetMimeType("application/vnd.etsi.tsl+xml").
		SetSchemeOperatorNames(map[string][]string{"en": {"Operator"}}).
		SetSchemeTypeCommunityRules(map[string][]string{"en": {"Rule"}}).
		SetMra(mra).
		Build()

	if p.Location() != "https://example.org/tl.xml" {
		t.Fatalf("unexpected Location: %s", p.Location())
	}
	if p.TSLLocation() != "https://example.org/tl.xml" {
		t.Fatalf("unexpected TSLLocation: %s", p.TSLLocation())
	}
	if p.SchemeTerritory() != "FR" {
		t.Fatalf("unexpected SchemeTerritory: %s", p.SchemeTerritory())
	}
	if p.TslType() != "http://uri.etsi.org/TrstSvc/TrustedList/TSLType/generic" {
		t.Fatalf("unexpected TslType: %s", p.TslType())
	}
	if p.MimeType() != "application/vnd.etsi.tsl+xml" {
		t.Fatalf("unexpected MimeType: %s", p.MimeType())
	}
	if p.SchemeOperatorNames()["en"][0] != "Operator" {
		t.Fatalf("unexpected SchemeOperatorNames: %v", p.SchemeOperatorNames())
	}
	if p.SchemeTypeCommunityRules()["en"][0] != "Rule" {
		t.Fatalf("unexpected SchemeTypeCommunityRules: %v", p.SchemeTypeCommunityRules())
	}
	if p.Mra() != mra {
		t.Fatalf("unexpected Mra: %v", p.Mra())
	}
	if p.SdiCertificates() != nil {
		t.Fatalf("expected nil SdiCertificates, got %v", p.SdiCertificates())
	}
}

func TestOtherTSLPointerEmptyConstructor(t *testing.T) {
	p := NewOtherTSLPointer()
	if p.Location() != "" {
		t.Fatalf("expected empty Location, got %s", p.Location())
	}
}
