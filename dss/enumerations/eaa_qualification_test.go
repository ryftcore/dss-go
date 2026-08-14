package enumerations

import "testing"

type eaaQualificationCase struct {
	v        EAAQualification
	readable string
	label    string
	uri      string
}

func eaaQualificationCases() []eaaQualificationCase {
	return []eaaQualificationCase{
		{EAAQualification_QEAA, "QEAA", "Qualified Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:QEAA"},
		{EAAQualification_EAA, "EAA", "Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:EAA"},
		{EAAQualification_PUBEAA, "PuB-EAA", "Electronic Attestation of Attributes issued by or on behalf of a public sector body", "urn:cef:dss:eaaQualification:PUBEAA"},
		{EAAQualification_PID, "PID", "Personal Identification Data", "urn:cef:dss:eaaQualification:PID"},
		{EAAQualification_UNKNOWN, "Unknown", "Electronic Attestation of Attributes of unknown type", "urn:cef:dss:eaaQualification:Unknown"},
		{EAAQualification_INDETERMINATE_QEAA, "Indeterminate QEAA", "Indeterminate Qualified Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:indeterminateQEAA"},
		{EAAQualification_INDETERMINATE_EAA, "Indeterminate EAA", "Indeterminate Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:indeterminateEAA"},
		{EAAQualification_INDETERMINATE_PUBEAA, "Indeterminate Pub-EAA", "Indeterminate Electronic Attestation of Attributes issued by or on behalf of a public sector body", "urn:cef:dss:eaaQualification:indeterminatePUBEAA"},
		{EAAQualification_INDETERMINATE_PID, "Indeterminate PID", "Indeterminate Personal Identification Data", "urn:cef:dss:eaaQualification:indeterminatePID"},
		{EAAQualification_INDETERMINATE_UNKNOWN, "Indeterminate Unknown", "Indeterminate Electronic Attestation of Attributes of unknown type", "urn:cef:dss:eaaQualification:indeterminateUnknown"},
		{EAAQualification_NOT_EAA, "Not EAA", "Not Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:NOTEAA"},
		{EAAQualification_NA, "N/A", "Not applicable", "urn:cef:dss:eaaQualification:NA"},
	}
}

func TestEAAQualificationFields(t *testing.T) {
	cases := eaaQualificationCases()
	if len(EAAQualificationValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(EAAQualificationValues()))
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

func TestEAAQualificationForName(t *testing.T) {
	for _, c := range eaaQualificationCases() {
		got, err := EAAQualificationForName(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("EAAQualificationForName(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
	}
	if got, err := EAAQualificationForName(""); err != nil || got != "" {
		t.Errorf("EAAQualificationForName(\"\") = %v, %v; want \"\", nil", got, err)
	}
	if _, err := EAAQualificationForName("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}

func TestEAAQualificationFromReadable(t *testing.T) {
	for _, c := range eaaQualificationCases() {
		if got := EAAQualificationFromReadable(c.readable); got != c.v {
			t.Errorf("EAAQualificationFromReadable(%q) = %v, want %v", c.readable, got, c.v)
		}
	}
	if got := EAAQualificationFromReadable(""); got != "" {
		t.Errorf("EAAQualificationFromReadable(\"\") = %v, want zero value", got)
	}
	if got := EAAQualificationFromReadable("nope"); got != "" {
		t.Errorf("EAAQualificationFromReadable(unknown) = %v, want zero value", got)
	}
}
