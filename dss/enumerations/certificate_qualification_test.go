package enumerations

import "testing"

func TestCertificateQualification(t *testing.T) {
	cases := []struct {
		v          CertificateQualification
		readable   string
		label      string
		isQc       bool
		certType   CertificateType
		isForEsig  bool
		isForEseal bool
		isQscd     bool
	}{
		{CertificateQualificationQCERTForESigQSCD, "QC for eSig with QSCD", "Qualified Certificate for Electronic Signatures with private key on QSCD", true, CertificateTypeESign, true, false, true},
		{CertificateQualificationQCERTForESealQSCD, "QC for eSeal with QSCD", "Qualified Certificate for Electronic Seals with private key on QSCD", true, CertificateTypeESeal, false, true, true},
		{CertificateQualificationQCERTForUnknownQSCD, "QC for unknown type with QSCD", "Qualified Certificate for unknown type with its private key residing in a QSCD", true, CertificateTypeUnknown, false, false, true},
		{CertificateQualificationQCERTForESig, "QC for eSig", "Qualified Certificate for Electronic Signatures", true, CertificateTypeESign, true, false, false},
		{CertificateQualificationQCERTForESeal, "QC for eSeal", "Qualified Certificate for Electronic Seals", true, CertificateTypeESeal, false, true, false},
		{CertificateQualificationQCERTForWSA, "QC for WSA", "Qualified Certificate for Web Site Authentications", true, CertificateTypeWSA, false, false, false},
		{CertificateQualificationQCERTForUnknown, "QC for unknown type", "Qualified Certificate for unknown type", true, CertificateTypeUnknown, false, false, false},
		{CertificateQualificationCertForESig, "Cert for eSig", "Certificate for Electronic Signatures", false, CertificateTypeESign, true, false, false},
		{CertificateQualificationCertForESeal, "Cert for eSeal", "Certificate for Electronic Seals", false, CertificateTypeESeal, false, true, false},
		{CertificateQualificationCertForWSA, "Cert for WSA", "Certificate for Web Site Authentications", false, CertificateTypeWSA, false, false, false},
		{CertificateQualificationCertForUnknown, "Cert for unknown type", "Certificate for unknown type", false, CertificateTypeUnknown, false, false, false},
		{CertificateQualificationNA, "N/A", "Not applicable", false, CertificateTypeUnknown, false, false, false},
	}
	if len(CertificateQualificationValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(CertificateQualificationValues()))
	}
	for _, c := range cases {
		if got := c.v.Readable(); got != c.readable {
			t.Errorf("%v.Readable() = %q, want %q", c.v, got, c.readable)
		}
		if got := c.v.Label(); got != c.label {
			t.Errorf("%v.Label() = %q, want %q", c.v, got, c.label)
		}
		if got := c.v.IsQc(); got != c.isQc {
			t.Errorf("%v.IsQc() = %v, want %v", c.v, got, c.isQc)
		}
		if got := c.v.Type(); got != c.certType {
			t.Errorf("%v.Type() = %v, want %v", c.v, got, c.certType)
		}
		if got := c.v.IsForEsig(); got != c.isForEsig {
			t.Errorf("%v.IsForEsig() = %v, want %v", c.v, got, c.isForEsig)
		}
		if got := c.v.IsForEseal(); got != c.isForEseal {
			t.Errorf("%v.IsForEseal() = %v, want %v", c.v, got, c.isForEseal)
		}
		if got := c.v.IsQscd(); got != c.isQscd {
			t.Errorf("%v.IsQscd() = %v, want %v", c.v, got, c.isQscd)
		}
		got, err := CertificateQualificationValueOf(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("CertificateQualificationValueOf(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
		if got := CertificateQualificationFromReadable(c.readable); got != c.v {
			t.Errorf("CertificateQualificationFromReadable(%q) = %v, want %v", c.readable, got, c.v)
		}
	}

	if _, err := CertificateQualificationValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
	if got, err := CertificateQualificationForName(""); got != "" || err != nil {
		t.Errorf("CertificateQualificationForName(\"\") = %v, %v; want \"\", nil", got, err)
	}
	if got, err := CertificateQualificationForName(string(CertificateQualificationNA)); got != CertificateQualificationNA || err != nil {
		t.Errorf("CertificateQualificationForName(NA) = %v, %v; want NA, nil", got, err)
	}
	if got := CertificateQualificationFromReadable(""); got != "" {
		t.Errorf("CertificateQualificationFromReadable(\"\") = %v, want \"\"", got)
	}
	if got := CertificateQualificationFromReadable("nope"); got != "" {
		t.Errorf("CertificateQualificationFromReadable(nope) = %v, want \"\"", got)
	}
}
