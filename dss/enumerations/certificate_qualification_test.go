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
		{CertificateQualification_QCERT_FOR_ESIG_QSCD, "QC for eSig with QSCD", "Qualified Certificate for Electronic Signatures with private key on QSCD", true, CertificateType_ESIGN, true, false, true},
		{CertificateQualification_QCERT_FOR_ESEAL_QSCD, "QC for eSeal with QSCD", "Qualified Certificate for Electronic Seals with private key on QSCD", true, CertificateType_ESEAL, false, true, true},
		{CertificateQualification_QCERT_FOR_UNKNOWN_QSCD, "QC for unknown type with QSCD", "Qualified Certificate for unknown type with its private key residing in a QSCD", true, CertificateType_UNKNOWN, false, false, true},
		{CertificateQualification_QCERT_FOR_ESIG, "QC for eSig", "Qualified Certificate for Electronic Signatures", true, CertificateType_ESIGN, true, false, false},
		{CertificateQualification_QCERT_FOR_ESEAL, "QC for eSeal", "Qualified Certificate for Electronic Seals", true, CertificateType_ESEAL, false, true, false},
		{CertificateQualification_QCERT_FOR_WSA, "QC for WSA", "Qualified Certificate for Web Site Authentications", true, CertificateType_WSA, false, false, false},
		{CertificateQualification_QCERT_FOR_UNKNOWN, "QC for unknown type", "Qualified Certificate for unknown type", true, CertificateType_UNKNOWN, false, false, false},
		{CertificateQualification_CERT_FOR_ESIG, "Cert for eSig", "Certificate for Electronic Signatures", false, CertificateType_ESIGN, true, false, false},
		{CertificateQualification_CERT_FOR_ESEAL, "Cert for eSeal", "Certificate for Electronic Seals", false, CertificateType_ESEAL, false, true, false},
		{CertificateQualification_CERT_FOR_WSA, "Cert for WSA", "Certificate for Web Site Authentications", false, CertificateType_WSA, false, false, false},
		{CertificateQualification_CERT_FOR_UNKNOWN, "Cert for unknown type", "Certificate for unknown type", false, CertificateType_UNKNOWN, false, false, false},
		{CertificateQualification_NA, "N/A", "Not applicable", false, CertificateType_UNKNOWN, false, false, false},
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
	if got, err := CertificateQualificationForName(string(CertificateQualification_NA)); got != CertificateQualification_NA || err != nil {
		t.Errorf("CertificateQualificationForName(NA) = %v, %v; want NA, nil", got, err)
	}
	if got := CertificateQualificationFromReadable(""); got != "" {
		t.Errorf("CertificateQualificationFromReadable(\"\") = %v, want \"\"", got)
	}
	if got := CertificateQualificationFromReadable("nope"); got != "" {
		t.Errorf("CertificateQualificationFromReadable(nope) = %v, want \"\"", got)
	}
}
