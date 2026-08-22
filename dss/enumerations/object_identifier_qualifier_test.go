package enumerations

import "testing"

func TestObjectIdentifierQualifierValue(t *testing.T) {
	tests := []struct {
		v    ObjectIdentifierQualifier
		want string
	}{
		{ObjectIdentifierQualifierOIDAsURI, "OIDAsURI"},
		{ObjectIdentifierQualifierOIDAsURN, "OIDAsURN"},
	}
	for _, tt := range tests {
		if got := tt.v.Value(); got != tt.want {
			t.Errorf("%v.Value() = %q, want %q", tt.v, got, tt.want)
		}
		if got := ObjectIdentifierQualifierFromValue(tt.want); got != tt.v {
			t.Errorf("ObjectIdentifierQualifierFromValue(%q) = %q, want %q", tt.want, got, tt.v)
		}
	}
	if got := ObjectIdentifierQualifierFromValue("bogus"); got != "" {
		t.Errorf("ObjectIdentifierQualifierFromValue(bogus) = %q, want empty", got)
	}
}

func TestObjectIdentifierQualifierValueOf(t *testing.T) {
	for _, v := range ObjectIdentifierQualifierValues() {
		got, err := ObjectIdentifierQualifierValueOf(string(v))
		if err != nil {
			t.Errorf("ObjectIdentifierQualifierValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("ObjectIdentifierQualifierValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestObjectIdentifierQualifierValueOfUnknown(t *testing.T) {
	if _, err := ObjectIdentifierQualifierValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
