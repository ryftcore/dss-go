package tsl

import "testing"

func TestServiceTypeASiRoundTrip(t *testing.T) {
	s := NewServiceTypeASi()
	s.SetType("http://uri.etsi.org/TrstSvc/Svctype/CA/QC")
	s.SetAsi("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/RootCA-QC")

	if s.Type() != "http://uri.etsi.org/TrstSvc/Svctype/CA/QC" {
		t.Fatalf("unexpected Type: %s", s.Type())
	}
	if s.Asi() != "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/RootCA-QC" {
		t.Fatalf("unexpected Asi: %s", s.Asi())
	}
}

func TestServiceTypeASiEquals(t *testing.T) {
	a := NewServiceTypeASi()
	a.SetType("t")
	a.SetAsi("a")
	b := NewServiceTypeASi()
	b.SetType("t")
	b.SetAsi("a")

	if !a.Equals(b) {
		t.Fatal("expected equal ServiceTypeASi to be Equals()")
	}
	b.SetAsi("other")
	if a.Equals(b) {
		t.Fatal("expected different Asi to not be Equals()")
	}
}

func TestServiceTypeASiComparableAsMapKey(t *testing.T) {
	// ServiceEquivalence.typeAsiEquivalence keys a map by value ServiceTypeASi, matching the
	// Java Map<ServiceTypeASi, ServiceTypeASi> keyed by structural equals/hashCode.
	m := map[ServiceTypeASi]ServiceTypeASi{
		{typ: "t1", asi: "a1"}: {typ: "t2", asi: "a2"},
	}
	if len(m) != 1 {
		t.Fatalf("expected ServiceTypeASi to be usable as a map key, got len=%d", len(m))
	}
}
