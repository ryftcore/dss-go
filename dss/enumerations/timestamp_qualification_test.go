package enumerations

import "testing"

func TestTimestampQualificationFields(t *testing.T) {
	cases := []struct {
		v        TimestampQualification
		readable string
		label    string
		uri      string
	}{
		{TimestampQualificationQTSA, "QTSA", "Qualified timestamp", "urn:cef:dss:timestampQualification:QTSA"},
		{TimestampQualificationTSA, "TSA", "Not qualified timestamp", "urn:cef:dss:timestampQualification:TSA"},
		{TimestampQualificationNA, "N/A", "Not applicable", "urn:cef:dss:timestampQualification:notApplicable"},
	}
	if len(TimestampQualificationValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(TimestampQualificationValues()))
	}
	for _, c := range cases {
		if got := c.v.Readable(); got != c.readable {
			t.Errorf("%v.Readable() = %q, want %q", c.v, got, c.readable)
		}
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
	}
}
