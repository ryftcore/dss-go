package enumerations

import "testing"

func TestArchiveTimestampHashIndexVersion(t *testing.T) {
	cases := []struct {
		v     ArchiveTimestampHashIndexVersion
		label string
		oid   string
	}{
		{ArchiveTimestampHashIndexVersion_ATS_HASH_INDEX, "ats-hash-index", "0.4.0.1733.2.5"},
		{ArchiveTimestampHashIndexVersion_ATS_HASH_INDEX_V2, "ats-hash-index-v2", "0.4.0.19122.1.4"},
		{ArchiveTimestampHashIndexVersion_ATS_HASH_INDEX_V3, "ats-hash-index-v3", "0.4.0.19122.1.5"},
	}
	if len(ArchiveTimestampHashIndexVersionValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(ArchiveTimestampHashIndexVersionValues()))
	}
	for _, c := range cases {
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
		if got := c.v.OID(); got != c.oid {
			t.Errorf("%v.OID() = %q, want %q", c.v, got, c.oid)
		}
		if got := ArchiveTimestampHashIndexVersionForLabel(c.label); got != c.v {
			t.Errorf("ArchiveTimestampHashIndexVersionForLabel(%q) = %v, want %v", c.label, got, c.v)
		}
		if got := ArchiveTimestampHashIndexVersionForOID(c.oid); got != c.v {
			t.Errorf("ArchiveTimestampHashIndexVersionForOID(%q) = %v, want %v", c.oid, got, c.v)
		}
	}
	if got := ArchiveTimestampHashIndexVersionForLabel("nope"); got != "" {
		t.Errorf("expected zero value for unknown label, got %v", got)
	}
	if got := ArchiveTimestampHashIndexVersionForOID("9.9.9"); got != "" {
		t.Errorf("expected zero value for unknown oid, got %v", got)
	}
}
