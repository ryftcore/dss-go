package enumerations

import "testing"

type signatureQualificationCase struct {
	v        SignatureQualification
	readable string
	label    string
	uri      string
}

func signatureQualificationCases() []signatureQualificationCase {
	return []signatureQualificationCase{
		{SignatureQualification_QESIG, "QESig", "Qualified Electronic Signature", "urn:cef:dss:signatureQualification:QESig"},
		{SignatureQualification_QESEAL, "QESeal", "Qualified Electronic Seal", "urn:cef:dss:signatureQualification:QESeal"},
		{SignatureQualification_UNKNOWN_QC_QSCD, "Unknown-QC-QSCD", "Signature produced by a Qualified Certificate of unknown type with its private key residing in a QSCD", "urn:cef:dss:signatureQualification:UnknownQCwithQSCD"},
		{SignatureQualification_ADESIG_QC, "AdESig-QC", "Advanced Electronic Signature supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:AdESigQC"},
		{SignatureQualification_ADESEAL_QC, "AdESeal-QC", "Advanced Electronic Seal supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:AdESealQC"},
		{SignatureQualification_UNKNOWN_QC, "Unknown-QC", "Signature produced by a Qualified Certificate of unknown type", "urn:cef:dss:signatureQualification:UnknownQC"},
		{SignatureQualification_ADESIG, "AdESig", "Advanced Electronic Signature", "urn:cef:dss:signatureQualification:AdESig"},
		{SignatureQualification_ADESEAL, "AdESeal", "Advanced Electronic Seal", "urn:cef:dss:signatureQualification:AdESeal"},
		{SignatureQualification_UNKNOWN, "Unknown", "Signature produced by a Certificate of unknown type", "urn:cef:dss:signatureQualification:Unknown"},
		{SignatureQualification_INDETERMINATE_QESIG, "Indeterminate QESig", "Indeterminate Qualified Electronic Signature", "urn:cef:dss:signatureQualification:indeterminateQESig"},
		{SignatureQualification_INDETERMINATE_QESEAL, "Indeterminate QESeal", "Indeterminate Qualified Electronic Seal", "urn:cef:dss:signatureQualification:indeterminateQESeal"},
		{SignatureQualification_INDETERMINATE_UNKNOWN_QC_QSCD, "Indeterminate Unknown-QC-QSCD", "Indeterminate Signature produced by a Qualified Certificate of unknown type with its private key residing in a QSCD", "urn:cef:dss:signatureQualification:indeterminateUnknownQCwithQSCD"},
		{SignatureQualification_INDETERMINATE_ADESIG_QC, "Indeterminate AdESig-QC", "Indeterminate Advanced Electronic Signature supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:indeterminateAdESigQC"},
		{SignatureQualification_INDETERMINATE_ADESEAL_QC, "Indeterminate AdESeal-QC", "Indeterminate Advanced Electronic Seal supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:indeterminateAdESealQC"},
		{SignatureQualification_INDETERMINATE_UNKNOWN_QC, "Indeterminate Unknown-QC", "Indeterminate Signature produced by a Qualified Certificate of unknown type", "urn:cef:dss:signatureQualification:indeterminateUnknownQC"},
		{SignatureQualification_INDETERMINATE_ADESIG, "Indeterminate AdESig", "Indeterminate Advanced Electronic Signature", "urn:cef:dss:signatureQualification:indeterminateAdESig"},
		{SignatureQualification_INDETERMINATE_ADESEAL, "Indeterminate AdESeal", "Indeterminate Advanced Electronic Seal", "urn:cef:dss:signatureQualification:indeterminateAdESeal"},
		{SignatureQualification_INDETERMINATE_UNKNOWN, "Indeterminate Unknown", "Indeterminate Signature produced by a Certificate of unknown type", "urn:cef:dss:signatureQualification:indeterminateUnknown"},
		{SignatureQualification_NOT_ADES_QC_QSCD, "Not AdES but QC with QSCD", "Not Advanced Electronic Signature but supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:notAdESbutQCwithQSCD"},
		{SignatureQualification_NOT_ADES_QC, "Not AdES but QC", "Not Advanced Electronic Signature but supported by a Qualified Certificate", "urn:cef:dss:signatureQualification:notAdESbutQC"},
		{SignatureQualification_NOT_ADES, "Not AdES", "Not Advanced Electronic Signature", "urn:cef:dss:signatureQualification:notAdES"},
		{SignatureQualification_NA, "N/A", "Not applicable", "urn:cef:dss:signatureQualification:notApplicable"},
	}
}

func TestSignatureQualificationFields(t *testing.T) {
	cases := signatureQualificationCases()
	if len(SignatureQualificationValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(SignatureQualificationValues()))
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

func TestSignatureQualificationForName(t *testing.T) {
	for _, c := range signatureQualificationCases() {
		got, err := SignatureQualificationForName(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("SignatureQualificationForName(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
	}
	if got, err := SignatureQualificationForName(""); err != nil || got != "" {
		t.Errorf("SignatureQualificationForName(\"\") = %v, %v; want \"\", nil", got, err)
	}
	if _, err := SignatureQualificationForName("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}

func TestSignatureQualificationFromReadable(t *testing.T) {
	for _, c := range signatureQualificationCases() {
		if got := SignatureQualificationFromReadable(c.readable); got != c.v {
			t.Errorf("SignatureQualificationFromReadable(%q) = %v, want %v", c.readable, got, c.v)
		}
	}
	if got := SignatureQualificationFromReadable(""); got != "" {
		t.Errorf("SignatureQualificationFromReadable(\"\") = %v, want zero value", got)
	}
	if got := SignatureQualificationFromReadable("nope"); got != "" {
		t.Errorf("SignatureQualificationFromReadable(unknown) = %v, want zero value", got)
	}
}

func TestSignatureQualificationForURI(t *testing.T) {
	for _, c := range signatureQualificationCases() {
		if got := SignatureQualificationForURI(c.uri); got != c.v {
			t.Errorf("SignatureQualificationForURI(%q) = %v, want %v", c.uri, got, c.v)
		}
	}
	if got := SignatureQualificationForURI(""); got != "" {
		t.Errorf("SignatureQualificationForURI(\"\") = %v, want zero value", got)
	}
	if got := SignatureQualificationForURI("nope"); got != "" {
		t.Errorf("SignatureQualificationForURI(unknown) = %v, want zero value", got)
	}
}
