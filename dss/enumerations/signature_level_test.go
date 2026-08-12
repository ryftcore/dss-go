// Ported from dss-enumerations/.../SignatureLevel.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestSignatureLevel_FormAndProfile(t *testing.T) {
	cases := []struct {
		v      SignatureLevel
		form   SignatureForm
		unsupp bool
		prof   SignatureProfile
		// ambiguous marks (form, profile) pairs shared by more than one
		// SignatureLevel constant (upstream: PKCS7_B/_T/_LT/_LTA all
		// construct with (PKCS7, NOT_ETSI)). Java's getSignatureLevel scans
		// values() in declaration order and returns the first match, so
		// only the first-declared constant for a given pair round-trips
		// through GetSignatureLevel; later duplicates legitimately resolve
		// back to that first constant instead of themselves.
		ambiguous bool
	}{
		{SignatureLevel_XML_NOT_ETSI, SignatureForm_XAdES, false, SignatureProfile_NOT_ETSI, false},
		{SignatureLevel_XAdES_BES, SignatureForm_XAdES, false, SignatureProfile_EXTENDED_BES, false},
		{SignatureLevel_XAdES_EPES, SignatureForm_XAdES, false, SignatureProfile_EXTENDED_EPES, false},
		{SignatureLevel_XAdES_T, SignatureForm_XAdES, false, SignatureProfile_EXTENDED_T, false},
		{SignatureLevel_XAdES_LT, SignatureForm_XAdES, false, SignatureProfile_EXTENDED_LT, false},
		{SignatureLevel_XAdES_C, SignatureForm_XAdES, false, SignatureProfile_EXTENDED_C, false},
		{SignatureLevel_XAdES_X, SignatureForm_XAdES, false, SignatureProfile_EXTENDED_X, false},
		{SignatureLevel_XAdES_XL, SignatureForm_XAdES, false, SignatureProfile_EXTENDED_XL, false},
		{SignatureLevel_XAdES_A, SignatureForm_XAdES, false, SignatureProfile_EXTENDED_A, false},
		{SignatureLevel_XAdES_ERS, SignatureForm_XAdES, false, SignatureProfile_EXTENDED_ERS, false},
		{SignatureLevel_XAdES_BASELINE_B, SignatureForm_XAdES, false, SignatureProfile_BASELINE_B, false},
		{SignatureLevel_XAdES_BASELINE_T, SignatureForm_XAdES, false, SignatureProfile_BASELINE_T, false},
		{SignatureLevel_XAdES_BASELINE_LT, SignatureForm_XAdES, false, SignatureProfile_BASELINE_LT, false},
		{SignatureLevel_XAdES_BASELINE_LTA, SignatureForm_XAdES, false, SignatureProfile_BASELINE_LTA, false},

		{SignatureLevel_CMS_NOT_ETSI, SignatureForm_CAdES, false, SignatureProfile_NOT_ETSI, false},
		{SignatureLevel_CAdES_BES, SignatureForm_CAdES, false, SignatureProfile_EXTENDED_BES, false},
		{SignatureLevel_CAdES_EPES, SignatureForm_CAdES, false, SignatureProfile_EXTENDED_EPES, false},
		{SignatureLevel_CAdES_T, SignatureForm_CAdES, false, SignatureProfile_EXTENDED_T, false},
		{SignatureLevel_CAdES_LT, SignatureForm_CAdES, false, SignatureProfile_EXTENDED_LT, false},
		{SignatureLevel_CAdES_C, SignatureForm_CAdES, false, SignatureProfile_EXTENDED_C, false},
		{SignatureLevel_CAdES_X, SignatureForm_CAdES, false, SignatureProfile_EXTENDED_X, false},
		{SignatureLevel_CAdES_XL, SignatureForm_CAdES, false, SignatureProfile_EXTENDED_XL, false},
		{SignatureLevel_CAdES_A, SignatureForm_CAdES, false, SignatureProfile_EXTENDED_A, false},
		{SignatureLevel_CAdES_ERS, SignatureForm_CAdES, false, SignatureProfile_EXTENDED_ERS, false},
		{SignatureLevel_CAdES_BASELINE_B, SignatureForm_CAdES, false, SignatureProfile_BASELINE_B, false},
		{SignatureLevel_CAdES_BASELINE_T, SignatureForm_CAdES, false, SignatureProfile_BASELINE_T, false},
		{SignatureLevel_CAdES_BASELINE_LT, SignatureForm_CAdES, false, SignatureProfile_BASELINE_LT, false},
		{SignatureLevel_CAdES_BASELINE_LTA, SignatureForm_CAdES, false, SignatureProfile_BASELINE_LTA, false},

		{SignatureLevel_PDF_NOT_ETSI, SignatureForm_PAdES, false, SignatureProfile_NOT_ETSI, false},
		{SignatureLevel_PKCS7_B, SignatureForm_PKCS7, false, SignatureProfile_NOT_ETSI, false},
		{SignatureLevel_PKCS7_T, SignatureForm_PKCS7, false, SignatureProfile_NOT_ETSI, true},
		{SignatureLevel_PKCS7_LT, SignatureForm_PKCS7, false, SignatureProfile_NOT_ETSI, true},
		{SignatureLevel_PKCS7_LTA, SignatureForm_PKCS7, false, SignatureProfile_NOT_ETSI, true},
		{SignatureLevel_PAdES_BES, SignatureForm_PAdES, false, SignatureProfile_EXTENDED_BES, false},
		{SignatureLevel_PAdES_EPES, SignatureForm_PAdES, false, SignatureProfile_EXTENDED_EPES, false},
		{SignatureLevel_PAdES_LTV, SignatureForm_PAdES, false, SignatureProfile_EXTENDED_LTV, false},
		{SignatureLevel_PAdES_BASELINE_B, SignatureForm_PAdES, false, SignatureProfile_BASELINE_B, false},
		{SignatureLevel_PAdES_BASELINE_T, SignatureForm_PAdES, false, SignatureProfile_BASELINE_T, false},
		{SignatureLevel_PAdES_BASELINE_LT, SignatureForm_PAdES, false, SignatureProfile_BASELINE_LT, false},
		{SignatureLevel_PAdES_BASELINE_LTA, SignatureForm_PAdES, false, SignatureProfile_BASELINE_LTA, false},

		{SignatureLevel_JSON_NOT_ETSI, SignatureForm_JAdES, false, SignatureProfile_NOT_ETSI, false},
		{SignatureLevel_JAdES, SignatureForm_JAdES, false, SignatureProfile_AdES, false},
		{SignatureLevel_JAdES_BASELINE_B, SignatureForm_JAdES, false, SignatureProfile_BASELINE_B, false},
		{SignatureLevel_JAdES_BASELINE_T, SignatureForm_JAdES, false, SignatureProfile_BASELINE_T, false},
		{SignatureLevel_JAdES_BASELINE_LT, SignatureForm_JAdES, false, SignatureProfile_BASELINE_LT, false},
		{SignatureLevel_JAdES_BASELINE_LTA, SignatureForm_JAdES, false, SignatureProfile_BASELINE_LTA, false},

		{SignatureLevel_CBOR_NOT_ETSI, SignatureForm_CBAdES, false, SignatureProfile_NOT_ETSI, false},
		{SignatureLevel_CB_AdES, SignatureForm_CBAdES, false, SignatureProfile_AdES, false},
		{SignatureLevel_CB_AdES_BASELINE_B, SignatureForm_CBAdES, false, SignatureProfile_BASELINE_B, false},
		{SignatureLevel_CB_AdES_BASELINE_T, SignatureForm_CBAdES, false, SignatureProfile_BASELINE_T, false},
		{SignatureLevel_CB_AdES_BASELINE_LT, SignatureForm_CBAdES, false, SignatureProfile_BASELINE_LT, false},
		{SignatureLevel_CB_AdES_BASELINE_LTA, SignatureForm_CBAdES, false, SignatureProfile_BASELINE_LTA, false},

		{SignatureLevel_UNKNOWN, "", true, SignatureProfile_NOT_ETSI, false},
	}
	if len(SignatureLevelValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(SignatureLevelValues()))
	}
	for _, c := range cases {
		form, err := c.v.SignatureForm()
		if c.unsupp {
			if err == nil {
				t.Errorf("%v.SignatureForm() expected error, got %v", c.v, form)
			}
		} else if err != nil || form != c.form {
			t.Errorf("%v.SignatureForm() = %v, %v; want %v, nil", c.v, form, err, c.form)
		}
		if got, err := c.v.SignatureProfile(); err != nil || got != c.prof {
			t.Errorf("%v.SignatureProfile() = %v, %v; want %v, nil", c.v, got, err, c.prof)
		}
		got, err := SignatureLevelValueOf(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("SignatureLevelValueOf(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
		if !c.unsupp && !c.ambiguous {
			if got, err := GetSignatureLevel(c.form, c.prof); err != nil || got != c.v {
				t.Errorf("GetSignatureLevel(%v, %v) = %v, %v; want %v, nil", c.form, c.prof, got, err, c.v)
			}
		} else if c.ambiguous {
			// Declaration-order collision: the (form, profile) pair is shared
			// with an earlier constant, so GetSignatureLevel legitimately
			// resolves to that earlier constant (PKCS7_B), not c.v.
			if got, err := GetSignatureLevel(c.form, c.prof); err != nil || got != SignatureLevel_PKCS7_B {
				t.Errorf("GetSignatureLevel(%v, %v) = %v, %v; want %v, nil (canonical first match)", c.form, c.prof, got, err, SignatureLevel_PKCS7_B)
			}
		}
	}
	if _, err := SignatureLevelValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
	// Upstream throws UnsupportedOperationException rather than returning null here: the
	// scan reaches UNKNOWN, whose null signature form makes getSignatureForm() throw.
	if got, err := GetSignatureLevel(SignatureForm_XAdES, SignatureProfile("bogus")); err == nil {
		t.Errorf("GetSignatureLevel(unmatched) = %v, nil; want an error", got)
	}
	if got, err := GetSignatureLevel("", ""); err == nil {
		t.Errorf("GetSignatureLevel(\"\", \"\") = %v, nil; want an error", got)
	}
}

func TestSignatureLevel_ValueByNameAndString(t *testing.T) {
	got, err := SignatureLevelValueByName("XAdES-BASELINE-B")
	if err != nil || got != SignatureLevel_XAdES_BASELINE_B {
		t.Errorf("SignatureLevelValueByName(XAdES-BASELINE-B) = %v, %v; want %v, nil", got, err, SignatureLevel_XAdES_BASELINE_B)
	}
	if _, err := SignatureLevelValueByName("NOPE-NOPE"); err == nil {
		t.Error("expected error for unknown dashed name")
	}
	if got := SignatureLevel_XAdES_BASELINE_B.String(); got != "XAdES-BASELINE-B" {
		t.Errorf("SignatureLevel_XAdES_BASELINE_B.String() = %q, want %q", got, "XAdES-BASELINE-B")
	}
}
