package diagnostic

import (
	"testing"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
)

func strp(s string) *string { return &s }
func boolp(b bool) *bool    { return &b }

func TestClaimWrapper_Text(t *testing.T) {
	c := NewClaimWrapper(&jaxb.XmlClaim{
		XmlClaimContent: jaxb.XmlClaimContent{Text: strp("hello")},
		XmlClaimAttrs:   jaxb.XmlClaimAttrs{Name: strp("Foo")},
	})
	if !c.IsText() || c.Text() != "hello" {
		t.Fatalf("expected text 'hello', got IsText=%v Text=%q", c.IsText(), c.Text())
	}
	if c.Name() != "Foo" {
		t.Fatalf("expected name Foo, got %q", c.Name())
	}
	if c.IsNumber() || c.IsBoolean() || c.IsBinary() || c.IsDateTime() || c.IsList() || c.IsMap() {
		t.Fatalf("expected only IsText true")
	}
	if c.DisplayValue() != "hello" {
		t.Fatalf("expected display value 'hello', got %q", c.DisplayValue())
	}
}

func TestClaimWrapper_ListAndMap(t *testing.T) {
	item1 := &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("a")}}
	item2 := &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("b")}}
	c := NewClaimWrapper(&jaxb.XmlClaim{
		XmlClaimContent: jaxb.XmlClaimContent{Item: []*jaxb.XmlClaim{item1, item2}},
	})
	if !c.IsList() {
		t.Fatalf("expected IsList true")
	}
	list := c.List()
	if len(list) != 2 || list[0].Text() != "a" || list[1].Text() != "b" {
		t.Fatalf("unexpected list contents: %+v", list)
	}
	if list[0].Parent() != c {
		t.Fatalf("expected list item parent to be c")
	}

	entry := &jaxb.XmlClaim{
		XmlClaimContent: jaxb.XmlClaimContent{Text: strp("v")},
		XmlClaimAttrs:   jaxb.XmlClaimAttrs{Name: strp("k")},
	}
	m := NewClaimWrapper(&jaxb.XmlClaim{
		XmlClaimContent: jaxb.XmlClaimContent{Entry: []*jaxb.XmlClaim{entry}},
	})
	if !m.IsMap() {
		t.Fatalf("expected IsMap true")
	}
	mp := m.Map()
	if len(mp) != 1 || mp["k"] == nil || mp["k"].Text() != "v" {
		t.Fatalf("unexpected map contents: %+v", mp)
	}
}

func TestClaimWrapper_Equals(t *testing.T) {
	a := NewClaimWrapper(&jaxb.XmlClaim{
		XmlClaimContent: jaxb.XmlClaimContent{Text: strp("x")},
		XmlClaimAttrs:   jaxb.XmlClaimAttrs{Name: strp("N")},
	})
	b := NewClaimWrapper(&jaxb.XmlClaim{
		XmlClaimContent: jaxb.XmlClaimContent{Text: strp("x")},
		XmlClaimAttrs:   jaxb.XmlClaimAttrs{Name: strp("N")},
	})
	if !a.Equals(b) {
		t.Fatalf("expected a.Equals(b) true")
	}
	c := NewClaimWrapper(&jaxb.XmlClaim{
		XmlClaimContent: jaxb.XmlClaimContent{Text: strp("y")},
		XmlClaimAttrs:   jaxb.XmlClaimAttrs{Name: strp("N")},
	})
	if a.Equals(c) {
		t.Fatalf("expected a.Equals(c) false")
	}
	if a.Equals(nil) {
		t.Fatalf("expected a.Equals(nil) false")
	}
}

// TestAddressClaimWrapper_AsClaimNoRecursion exercises the AsClaim()->Map()->getters cycle that
// previously could recurse infinitely (see the package note in claim_wrapper.go); it must
// terminate and produce the expected map.
func TestAddressClaimWrapper_AsClaimNoRecursion(t *testing.T) {
	wrapped := &jaxb.XmlAddressClaim{
		City:        &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("Lisbon")}, XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("City")}},
		CountryName: &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("PT")}, XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("Country")}},
	}
	addr := NewAddressClaimWrapper(wrapped)
	if !addr.IsMap() {
		t.Fatalf("expected AddressClaimWrapper.IsMap() true")
	}
	m := addr.Map()
	if len(m) != 2 {
		t.Fatalf("expected 2 map entries, got %d: %+v", len(m), m)
	}
	if m["City"] == nil || m["City"].Text() != "Lisbon" {
		t.Fatalf("expected City=Lisbon, got %+v", m["City"])
	}

	asClaim := addr.AsClaim()
	if !asClaim.IsMap() {
		t.Fatalf("expected AsClaim().IsMap() true")
	}
	m2 := asClaim.Map()
	if len(m2) != 2 || m2["Country"] == nil || m2["Country"].Text() != "PT" {
		t.Fatalf("unexpected AsClaim().Map() contents: %+v", m2)
	}
}

func TestDrivingPrivilegesClaimWrapper_ListOverride(t *testing.T) {
	dp := &jaxb.XmlDrivingPrivilegeClaim{
		VehicleCategoryCode: &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("B")}, XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("VehicleCategoryCode")}},
	}
	wrapped := &jaxb.XmlDrivingPrivilegesClaim{DrivingPrivilege: []*jaxb.XmlDrivingPrivilegeClaim{dp}}
	dps := NewDrivingPrivilegesClaimWrapper(wrapped)
	if !dps.IsList() {
		t.Fatalf("expected IsList true")
	}
	list := dps.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 list entry, got %d", len(list))
	}
	if !list[0].IsMap() {
		t.Fatalf("expected the wrapped DrivingPrivilegeClaimWrapper's AsClaim() to report IsMap true")
	}
	asClaim := dps.AsClaim()
	if !asClaim.IsList() || len(asClaim.List()) != 1 {
		t.Fatalf("expected AsClaim().List() to carry the same single entry")
	}
}

// TestCredentialSubjectClaimWrapper_NestedOverrides exercises the mixed concrete/generic
// PlaceOfBirth typing (see the file note in place_of_birth_claim_wrapper.go) alongside the
// concretely-typed Birthdate/Address fields.
func TestCredentialSubjectClaimWrapper_NestedOverrides(t *testing.T) {
	wrapped := &jaxb.XmlCredentialSubjectClaim{
		FullName: &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("Jane Doe")}, XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("FullName")}},
		Birthdate: &jaxb.XmlBirthdateClaim{
			Birthdate:     &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("1990-01-01")}, XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("Birthdate")}},
			XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("Birthdate")},
		},
		Address: &jaxb.XmlAddressClaim{
			City:          &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("Porto")}, XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("City")}},
			XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("Address")},
		},
		// Generic XmlClaim PlaceOfBirth (see the file note in place_of_birth_claim_wrapper.go):
		// exercises the false-instanceof path (no City/Region/Country available).
		PlaceOfBirth: &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("Porto, PT")}, XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("PlaceOfBirth")}},
	}
	cs := NewCredentialSubjectClaimWrapper(wrapped)
	if !cs.IsMap() {
		t.Fatalf("expected IsMap true")
	}

	fullName := cs.FullName()
	if fullName == nil || fullName.Text() != "Jane Doe" {
		t.Fatalf("unexpected FullName: %+v", fullName)
	}

	birthdate := cs.Birthdate()
	if birthdate == nil {
		t.Fatalf("expected Birthdate non-nil")
	}
	if bd := birthdate.Birthdate(); bd == nil || bd.Text() != "1990-01-01" {
		t.Fatalf("unexpected nested Birthdate.Birthdate(): %+v", bd)
	}

	address := cs.Address()
	if address == nil || address.City() == nil || address.City().Text() != "Porto" {
		t.Fatalf("unexpected Address: %+v", address)
	}

	placeOfBirth := cs.PlaceOfBirth()
	if placeOfBirth == nil {
		t.Fatalf("expected PlaceOfBirth non-nil")
	}
	if placeOfBirth.City() != nil {
		t.Fatalf("expected City() nil for a generically-typed PlaceOfBirth claim, got %+v", placeOfBirth.City())
	}
	if placeOfBirth.IsMap() {
		t.Fatalf("expected IsMap() false for a generically-typed PlaceOfBirth claim")
	}
	if placeOfBirth.Text() != "Porto, PT" {
		t.Fatalf("expected the generic claim's own Text() to remain readable, got %q", placeOfBirth.Text())
	}

	// Map() must recurse through the overridden children without infinite looping.
	m := cs.Map()
	if m["FullName"] == nil || m["Birthdate"] == nil || m["Address"] == nil || m["PlaceOfBirth"] == nil {
		t.Fatalf("expected all four claims present in Map(): %+v", m)
	}
	if !m["Birthdate"].IsMap() {
		t.Fatalf("expected the Birthdate entry in the map to itself report IsMap true (AsClaim() override baked in)")
	}
}

func TestPlaceOfBirthClaimWrapper_ConcretePath(t *testing.T) {
	wrapped := &jaxb.XmlPlaceOfBirthClaim{
		City:    &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("Faro")}, XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("City")}},
		Country: &jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp("PT")}, XmlClaimAttrs: jaxb.XmlClaimAttrs{Name: strp("Country")}},
	}
	pob := NewPlaceOfBirthClaimWrapper(wrapped)
	if !pob.IsMap() {
		t.Fatalf("expected IsMap true for a concretely-typed PlaceOfBirthClaim")
	}
	if pob.City() == nil || pob.City().Text() != "Faro" {
		t.Fatalf("unexpected City(): %+v", pob.City())
	}
	if pob.Wrapped() != wrapped {
		t.Fatalf("expected Wrapped() to return the concrete value")
	}
}

func TestClaimWrapper_DateTimeDisplayValue(t *testing.T) {
	ts := jaxb.NewXSDateTime(time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC))
	c := NewClaimWrapper(&jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{DateTime: ts}})
	if !c.IsDateTime() {
		t.Fatalf("expected IsDateTime true")
	}
	if got, want := c.DisplayValue(), "2024-03-15T10:30:00Z"; got != want {
		t.Fatalf("DisplayValue() = %q, want %q", got, want)
	}
}

func TestClaimWrapper_NumberAndBoolean(t *testing.T) {
	n := NewClaimWrapper(&jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Number: jaxb.NewBigIntegerFromInt64(42)}})
	if !n.IsNumber() || n.Number().Int64() != 42 {
		t.Fatalf("unexpected number claim: %+v", n.Number())
	}
	if n.DisplayValue() != "42" {
		t.Fatalf("expected display value 42, got %q", n.DisplayValue())
	}

	b := NewClaimWrapper(&jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Boolean: boolp(true)}})
	if !b.IsBoolean() || !*b.Boolean() {
		t.Fatalf("unexpected boolean claim")
	}
	if b.DisplayValue() != "true" {
		t.Fatalf("expected display value true, got %q", b.DisplayValue())
	}
}

// TestClaimMapKeyOrderIsSorted pins the stable order claimMapDisplayValue renders a map-valued
// claim in; see claimMapKeyOrder for why Java's own HashMap order cannot be reproduced.
func TestClaimMapKeyOrderIsSorted(t *testing.T) {
	keys := []string{"formatted", "street_address", "locality", "region", "postal_code"}
	want := []string{"formatted", "locality", "postal_code", "region", "street_address"}
	entries := make(map[string]*ClaimWrapper, len(keys))
	for _, key := range keys {
		entries[key] = NewClaimWrapper(&jaxb.XmlClaim{})
	}
	got := claimMapKeyOrder(entries)
	if len(got) != len(want) {
		t.Fatalf("claimMapKeyOrder() = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("claimMapKeyOrder() = %v, want %v", got, want)
		}
	}
}

// TestClaimMapDisplayValueIsStable guards the determinism the order function buys: the same map
// must render the same string on every call, which ranging over a Go map would not give.
func TestClaimMapDisplayValueIsStable(t *testing.T) {
	entries := map[string]*ClaimWrapper{}
	for _, key := range []string{"formatted", "street_address", "locality", "region", "postal_code"} {
		entries[key] = NewClaimWrapper(&jaxb.XmlClaim{XmlClaimContent: jaxb.XmlClaimContent{Text: strp(key)}})
	}
	first := claimMapDisplayValue(entries)
	for i := 0; i < 64; i++ {
		if got := claimMapDisplayValue(entries); got != first {
			t.Fatalf("claimMapDisplayValue is not stable:\n %s\n %s", first, got)
		}
	}
}
