package enumerations

import "testing"

func TestEndorsementTypeValue(t *testing.T) {
	cases := []struct {
		v    EndorsementType
		want string
	}{
		{EndorsementTypeCertified, "certified"},
		{EndorsementTypeClaimed, "claimed"},
		{EndorsementTypeSigned, "signed"},
	}
	for _, c := range cases {
		if got := c.v.Value(); got != c.want {
			t.Errorf("%v.Value() = %q, want %q", c.v, got, c.want)
		}
	}
}

func TestEndorsementTypeFromString(t *testing.T) {
	for _, c := range []struct {
		s    string
		want EndorsementType
	}{
		{"certified", EndorsementTypeCertified},
		{"claimed", EndorsementTypeClaimed},
		{"signed", EndorsementTypeSigned},
	} {
		got := EndorsementTypeFromString(c.s)
		if got != c.want {
			t.Errorf("EndorsementTypeFromString(%q) = %q, want %q", c.s, got, c.want)
		}
	}

	// Upstream returns null for an unknown value rather than throwing.
	if got := EndorsementTypeFromString("bogus"); got != "" {
		t.Errorf("EndorsementTypeFromString(\"bogus\") = %q, want \"\"", got)
	}
}

func TestEndorsementTypeValueOf(t *testing.T) {
	for _, v := range EndorsementTypeValues() {
		got, err := EndorsementTypeValueOf(string(v))
		if err != nil {
			t.Fatalf("EndorsementTypeValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("EndorsementTypeValueOf(%q) = %q, want %q", v, got, v)
		}
	}
	if _, err := EndorsementTypeValueOf("bogus"); err == nil {
		t.Error("EndorsementTypeValueOf(\"bogus\") expected error, got nil")
	}
}
