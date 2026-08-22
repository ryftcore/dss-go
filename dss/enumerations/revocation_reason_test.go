package enumerations

import "testing"

func TestRevocationReason(t *testing.T) {
	cases := []struct {
		v         RevocationReason
		shortName string
		uri       string
		value     int
	}{
		{RevocationReasonUnspecified, "unspecified", "urn:etsi:019102:revocationReason:unspecified", 0},
		{RevocationReasonKeyCompromise, "keyCompromise", "urn:etsi:019102:revocationReason:keyCompromise", 1},
		{RevocationReasonCACompromise, "cACompromise", "urn:etsi:019102:revocationReason:cACompromise", 2},
		{RevocationReasonAffiliationChanged, "affiliationChanged", "urn:etsi:019102:revocationReason:affiliationChanged", 3},
		{RevocationReasonSuperseded, "superseded", "urn:etsi:019102:revocationReason:superseded", 4},
		{RevocationReasonCessationOfOperation, "cessationOfOperation", "urn:etsi:019102:revocationReason:cessationOfOperation", 5},
		{RevocationReasonCertificateHold, "certificateHold", "urn:etsi:019102:revocationReason:certificateHold", 6},
		{RevocationReasonRemoveFromCRL, "removeFromCRL", "urn:etsi:019102:revocationReason:removeFromCRL", 8},
		{RevocationReasonPrivilegeWithdrawn, "privilegeWithdrawn", "urn:etsi:019102:revocationReason:privilegeWithdrawn", 9},
		{RevocationReasonAACompromise, "aACompromise", "urn:etsi:019102:revocationReason:aACompromise", 10},
	}
	if len(RevocationReasonValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(RevocationReasonValues()))
	}
	for _, c := range cases {
		if got := c.v.ShortName(); got != c.shortName {
			t.Errorf("%v.ShortName() = %q, want %q", c.v, got, c.shortName)
		}
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
		if got := c.v.Value(); got != c.value {
			t.Errorf("%v.Value() = %d, want %d", c.v, got, c.value)
		}
		if got := RevocationReasonFromInt(c.value); got != c.v {
			t.Errorf("RevocationReasonFromInt(%d) = %v, want %v", c.value, got, c.v)
		}
		if got := RevocationReasonFromValue(c.shortName); got != c.v {
			t.Errorf("RevocationReasonFromValue(%q) = %v, want %v", c.shortName, got, c.v)
		}
	}
	if got := RevocationReasonFromInt(7); got != "" {
		t.Errorf("RevocationReasonFromInt(7) = %v, want \"\" (7 is unused)", got)
	}
	if got := RevocationReasonFromValue("nope"); got != "" {
		t.Errorf("RevocationReasonFromValue(nope) = %v, want \"\"", got)
	}
}
