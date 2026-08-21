package lote

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model/tsl"
)

func TestOtherListPointerBuilderRoundTrip(t *testing.T) {
	mra := tsl.NewMRA()
	builder := NewOtherListPointerBuilder().
		SetLocationUrl("https://example.org/lote.xml").
		SetSchemeTerritory("BE").
		SetTslType("http://uri.etsi.org/19602/StatusDetermination/independent").
		SetMimeType("application/vnd.etsi.tsl+xml").
		SetSchemeOperatorNames(map[string][]string{"en": {"Example"}}).
		SetSchemeTypeCommunityRules(map[string][]string{"en": {"rules"}}).
		SetMra(mra)

	pointer := builder.Build()

	if pointer.LocationUrl() != "https://example.org/lote.xml" {
		t.Fatalf("unexpected LocationUrl(): %s", pointer.LocationUrl())
	}
	if pointer.Type() != "http://uri.etsi.org/19602/StatusDetermination/independent" {
		t.Fatalf("unexpected Type(): %s", pointer.Type())
	}
	if pointer.SchemeTerritory() != "BE" {
		t.Fatalf("unexpected SchemeTerritory(): %s", pointer.SchemeTerritory())
	}
	// DEVIATION kept verbatim from upstream Java: the builder's mra is never copied onto the
	// built OtherListPointer, which has no Mra() accessor at all - only the builder does.
	if builder.Mra() != mra {
		t.Fatal("expected builder.Mra() to return the value set via SetMra")
	}
}

func TestNewOtherListPointerIsEmpty(t *testing.T) {
	p := NewOtherListPointer()
	if p.LocationUrl() != "" || p.SdiCertificates() != nil {
		t.Fatalf("expected a zero-value OtherListPointer, got %+v", p)
	}
}
