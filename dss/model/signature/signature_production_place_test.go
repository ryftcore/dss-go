// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignatureProductionPlace.java (DSS 6.5.RC1).
package signature

import (
	"reflect"
	"testing"
)

func TestSignatureProductionPlace_RoundTrip(t *testing.T) {
	p := NewProductionPlace()
	if got := p.PostalAddress(); len(got) != 0 {
		t.Fatalf("PostalAddress() on a fresh instance = %v, want empty slice", got)
	}

	p.SetCity("Luxembourg")
	p.SetStateOrProvince("Luxembourg District")
	p.SetPostOfficeBoxNumber("123")
	p.SetPostalCode("L-1234")
	p.SetCountryName("LU")
	p.SetStreetAddress("1 rue du Fort")
	p.SetPostalAddress([]string{"1 rue du Fort", "L-1234 Luxembourg"})

	if got, want := p.City(), "Luxembourg"; got != want {
		t.Fatalf("City() = %q, want %q", got, want)
	}
	if got, want := p.StateOrProvince(), "Luxembourg District"; got != want {
		t.Fatalf("StateOrProvince() = %q, want %q", got, want)
	}
	if got, want := p.PostOfficeBoxNumber(), "123"; got != want {
		t.Fatalf("PostOfficeBoxNumber() = %q, want %q", got, want)
	}
	if got, want := p.PostalCode(), "L-1234"; got != want {
		t.Fatalf("PostalCode() = %q, want %q", got, want)
	}
	if got, want := p.CountryName(), "LU"; got != want {
		t.Fatalf("CountryName() = %q, want %q", got, want)
	}
	if got, want := p.StreetAddress(), "1 rue du Fort"; got != want {
		t.Fatalf("StreetAddress() = %q, want %q", got, want)
	}
	if got, want := p.PostalAddress(), []string{"1 rue du Fort", "L-1234 Luxembourg"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("PostalAddress() = %v, want %v", got, want)
	}
}
