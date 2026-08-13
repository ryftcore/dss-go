package utils

import "testing"

func TestNamespaceContextMapRegisterAndResolve(t *testing.T) {
	m := NewNamespaceContextMap()

	if !m.RegisterNamespace("ds", "http://www.w3.org/2000/09/xmldsig#") {
		t.Error("expected first registration to report true (new prefix)")
	}
	if m.RegisterNamespace("ds", "http://www.w3.org/2000/09/xmldsig#") {
		t.Error("expected re-registration to report false (existed already)")
	}

	if got := m.NamespaceURI("ds"); got != "http://www.w3.org/2000/09/xmldsig#" {
		t.Errorf("NamespaceURI(ds) = %q", got)
	}
	if got := m.NamespaceURI("unregistered"); got != "" {
		t.Errorf("NamespaceURI(unregistered) = %q, want \"\"", got)
	}

	prefix, ok := m.Prefix("http://www.w3.org/2000/09/xmldsig#")
	if !ok || prefix != "ds" {
		t.Errorf("Prefix() = (%q, %v), want (ds, true)", prefix, ok)
	}

	if _, ok := m.Prefix("urn:unregistered"); ok {
		t.Error("expected Prefix() to report false for an unregistered URI")
	}
}

func TestNamespaceContextMapRebind(t *testing.T) {
	m := NewNamespaceContextMap()
	m.RegisterNamespace("p", "urn:a")
	m.RegisterNamespace("p", "urn:b")
	if got := m.NamespaceURI("p"); got != "urn:b" {
		t.Errorf("expected rebound prefix to resolve to the latest URI, got %q", got)
	}
}

func TestNamespaceContextMapPrefixesDeterministic(t *testing.T) {
	m := NewNamespaceContextMap()
	m.RegisterNamespace("b", "urn:shared")
	m.RegisterNamespace("a", "urn:shared")
	m.RegisterNamespace("c", "urn:shared")

	got := m.Prefixes("urn:shared")
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("Prefixes() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Prefixes() = %v, want %v", got, want)
		}
	}

	if empty := m.Prefixes("urn:never-registered"); len(empty) != 0 {
		t.Fatalf("expected an empty slice for an unregistered URI, got %v", empty)
	}
}

func TestNamespaceContextMapPrefixMapIsACopy(t *testing.T) {
	m := NewNamespaceContextMap()
	m.RegisterNamespace("p", "urn:a")

	pm := m.PrefixMap()
	pm["p"] = "urn:mutated"

	if got := m.NamespaceURI("p"); got != "urn:a" {
		t.Fatalf("expected PrefixMap() to return a copy, but the registry changed to %q", got)
	}
}

func TestNamespaceContextMapToXPath10(t *testing.T) {
	m := NewNamespaceContextMap()
	m.RegisterNamespace("ds", "urn:ds")

	ns := namespaceContextMapToXPath10(m)
	if got := ns.NamespaceURI("ds"); got != "urn:ds" {
		t.Errorf("namespaceContextMapToXPath10 conversion: NamespaceURI(ds) = %q", got)
	}

	if got := namespaceContextMapToXPath10(nil); got != nil {
		t.Errorf("expected nil conversion for nil map, got %v", got)
	}
}
