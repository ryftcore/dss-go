package enumerations

import "testing"

func TestSemanticsIdentifierFields(t *testing.T) {
	cases := []struct {
		v           SemanticsIdentifier
		name        string
		oid         string
		description string
	}{
		{SemanticsIdentifierQcsSemanticsIdNatural, "qcs-semanticsId-Natural", "0.4.0.194121.1.1", "Semantics identifier for natural person"},
		{SemanticsIdentifierQcsSemanticsIdLegal, "qcs-SemanticsId-Legal", "0.4.0.194121.1.2", "Semantics identifier for legal person"},
		{SemanticsIdentifierQcsSemanticsIdEIDASNatural, "qcs-semanticsId-eIDASNatural", "0.4.0.194121.1.3", "Semantics identifier for eIDAS natural person"},
		{SemanticsIdentifierQcsSemanticsIdEIDASLegal, "qcs-SemanticsId-eIDASLegal", "0.4.0.194121.1.4", "Semantics identifier for eIDAS legal person"},
	}
	for _, c := range cases {
		if got := c.v.Name(); got != c.name {
			t.Errorf("%v.Name() = %q, want %q", c.v, got, c.name)
		}
		if got := c.v.OID(); got != c.oid {
			t.Errorf("%v.OID() = %q, want %q", c.v, got, c.oid)
		}
		if got := c.v.Description(); got != c.description {
			t.Errorf("%v.Description() = %q, want %q", c.v, got, c.description)
		}
	}
}

func TestSemanticsIdentifierFromName(t *testing.T) {
	for _, v := range SemanticsIdentifierValues() {
		got := SemanticsIdentifierFromName(v.Name())
		if got != v {
			t.Errorf("SemanticsIdentifierFromName(%q) = %q, want %q", v.Name(), got, v)
		}
	}
	// Upstream returns null for an unknown name rather than throwing.
	if got := SemanticsIdentifierFromName("bogus"); got != "" {
		t.Errorf("SemanticsIdentifierFromName(\"bogus\") = %q, want \"\"", got)
	}
}

func TestSemanticsIdentifierFromOID(t *testing.T) {
	for _, v := range SemanticsIdentifierValues() {
		got := SemanticsIdentifierFromOID(v.OID())
		if got != v {
			t.Errorf("SemanticsIdentifierFromOID(%q) = %q, want %q", v.OID(), got, v)
		}
	}
	// Upstream returns null for an unknown OID rather than throwing.
	if got := SemanticsIdentifierFromOID("9.9.9"); got != "" {
		t.Errorf("SemanticsIdentifierFromOID(\"9.9.9\") = %q, want \"\"", got)
	}
}

func TestSemanticsIdentifierValueOf(t *testing.T) {
	for _, v := range SemanticsIdentifierValues() {
		got, err := SemanticsIdentifierValueOf(string(v))
		if err != nil {
			t.Fatalf("SemanticsIdentifierValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("SemanticsIdentifierValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := SemanticsIdentifierValueOf("bogus"); err == nil {
		t.Error("SemanticsIdentifierValueOf(\"bogus\") expected error, got nil")
	}
}
