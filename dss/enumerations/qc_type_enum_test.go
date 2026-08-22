package enumerations

import "testing"

func TestQCTypeEnum(t *testing.T) {
	cases := []struct {
		v           QCTypeEnum
		description string
		oid         string
	}{
		{QCTypeEnumQCTESign, "qc-type-esign", "0.4.0.1862.1.6.1"},
		{QCTypeEnumQCTESeal, "qc-type-eseal", "0.4.0.1862.1.6.2"},
		{QCTypeEnumQCTWeb, "qc-type-web", "0.4.0.1862.1.6.3"},
		{QCTypeEnumQCTPID, "qc-type-pid", "0.4.0.194126.1.1"},
		{QCTypeEnumQCTWAL, "qc-type-wal", "0.4.0.194126.1.2"},
	}
	if len(QCTypeEnumValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(QCTypeEnumValues()))
	}
	for _, c := range cases {
		if got := c.v.Description(); got != c.description {
			t.Errorf("%v.Description() = %q, want %q", c.v, got, c.description)
		}
		if got := c.v.OID(); got != c.oid {
			t.Errorf("%v.OID() = %q, want %q", c.v, got, c.oid)
		}
		if got := QCTypeEnumForLabel(c.description); got != c.v {
			t.Errorf("QCTypeEnumForLabel(%q) = %v, want %v", c.description, got, c.v)
		}
	}
	if got := QCTypeEnumForLabel("nope"); got != "" {
		t.Errorf("expected zero value for unknown label, got %v", got)
	}
}
