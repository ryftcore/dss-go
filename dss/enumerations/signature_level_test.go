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
		{SignatureLevelXMLNotETSI, SignatureFormXAdES, false, SignatureProfileNotETSI, false},
		{SignatureLevelXAdESBES, SignatureFormXAdES, false, SignatureProfileExtendedBES, false},
		{SignatureLevelXAdESEPES, SignatureFormXAdES, false, SignatureProfileExtendedEPES, false},
		{SignatureLevelXAdEST, SignatureFormXAdES, false, SignatureProfileExtendedT, false},
		{SignatureLevelXAdESLT, SignatureFormXAdES, false, SignatureProfileExtendedLT, false},
		{SignatureLevelXAdESC, SignatureFormXAdES, false, SignatureProfileExtendedC, false},
		{SignatureLevelXAdESX, SignatureFormXAdES, false, SignatureProfileExtendedX, false},
		{SignatureLevelXAdESXL, SignatureFormXAdES, false, SignatureProfileExtendedXL, false},
		{SignatureLevelXAdESA, SignatureFormXAdES, false, SignatureProfileExtendedA, false},
		{SignatureLevelXAdESERS, SignatureFormXAdES, false, SignatureProfileExtendedERS, false},
		{SignatureLevelXAdESBaselineB, SignatureFormXAdES, false, SignatureProfileBaselineB, false},
		{SignatureLevelXAdESBaselineT, SignatureFormXAdES, false, SignatureProfileBaselineT, false},
		{SignatureLevelXAdESBaselineLT, SignatureFormXAdES, false, SignatureProfileBaselineLT, false},
		{SignatureLevelXAdESBaselineLTA, SignatureFormXAdES, false, SignatureProfileBaselineLTA, false},

		{SignatureLevelCMSNotETSI, SignatureFormCAdES, false, SignatureProfileNotETSI, false},
		{SignatureLevelCAdESBES, SignatureFormCAdES, false, SignatureProfileExtendedBES, false},
		{SignatureLevelCAdESEPES, SignatureFormCAdES, false, SignatureProfileExtendedEPES, false},
		{SignatureLevelCAdEST, SignatureFormCAdES, false, SignatureProfileExtendedT, false},
		{SignatureLevelCAdESLT, SignatureFormCAdES, false, SignatureProfileExtendedLT, false},
		{SignatureLevelCAdESC, SignatureFormCAdES, false, SignatureProfileExtendedC, false},
		{SignatureLevelCAdESX, SignatureFormCAdES, false, SignatureProfileExtendedX, false},
		{SignatureLevelCAdESXL, SignatureFormCAdES, false, SignatureProfileExtendedXL, false},
		{SignatureLevelCAdESA, SignatureFormCAdES, false, SignatureProfileExtendedA, false},
		{SignatureLevelCAdESERS, SignatureFormCAdES, false, SignatureProfileExtendedERS, false},
		{SignatureLevelCAdESBaselineB, SignatureFormCAdES, false, SignatureProfileBaselineB, false},
		{SignatureLevelCAdESBaselineT, SignatureFormCAdES, false, SignatureProfileBaselineT, false},
		{SignatureLevelCAdESBaselineLT, SignatureFormCAdES, false, SignatureProfileBaselineLT, false},
		{SignatureLevelCAdESBaselineLTA, SignatureFormCAdES, false, SignatureProfileBaselineLTA, false},

		{SignatureLevelPDFNotETSI, SignatureFormPAdES, false, SignatureProfileNotETSI, false},
		{SignatureLevelPKCS7B, SignatureFormPKCS7, false, SignatureProfileNotETSI, false},
		{SignatureLevelPKCS7T, SignatureFormPKCS7, false, SignatureProfileNotETSI, true},
		{SignatureLevelPKCS7LT, SignatureFormPKCS7, false, SignatureProfileNotETSI, true},
		{SignatureLevelPKCS7LTA, SignatureFormPKCS7, false, SignatureProfileNotETSI, true},
		{SignatureLevelPAdESBES, SignatureFormPAdES, false, SignatureProfileExtendedBES, false},
		{SignatureLevelPAdESEPES, SignatureFormPAdES, false, SignatureProfileExtendedEPES, false},
		{SignatureLevelPAdESLTV, SignatureFormPAdES, false, SignatureProfileExtendedLTV, false},
		{SignatureLevelPAdESBaselineB, SignatureFormPAdES, false, SignatureProfileBaselineB, false},
		{SignatureLevelPAdESBaselineT, SignatureFormPAdES, false, SignatureProfileBaselineT, false},
		{SignatureLevelPAdESBaselineLT, SignatureFormPAdES, false, SignatureProfileBaselineLT, false},
		{SignatureLevelPAdESBaselineLTA, SignatureFormPAdES, false, SignatureProfileBaselineLTA, false},

		{SignatureLevelJSONNotETSI, SignatureFormJAdES, false, SignatureProfileNotETSI, false},
		{SignatureLevelJAdES, SignatureFormJAdES, false, SignatureProfileAdES, false},
		{SignatureLevelJAdESBaselineB, SignatureFormJAdES, false, SignatureProfileBaselineB, false},
		{SignatureLevelJAdESBaselineT, SignatureFormJAdES, false, SignatureProfileBaselineT, false},
		{SignatureLevelJAdESBaselineLT, SignatureFormJAdES, false, SignatureProfileBaselineLT, false},
		{SignatureLevelJAdESBaselineLTA, SignatureFormJAdES, false, SignatureProfileBaselineLTA, false},

		{SignatureLevelCBORNotETSI, SignatureFormCBAdES, false, SignatureProfileNotETSI, false},
		{SignatureLevelCBAdES, SignatureFormCBAdES, false, SignatureProfileAdES, false},
		{SignatureLevelCBAdESBaselineB, SignatureFormCBAdES, false, SignatureProfileBaselineB, false},
		{SignatureLevelCBAdESBaselineT, SignatureFormCBAdES, false, SignatureProfileBaselineT, false},
		{SignatureLevelCBAdESBaselineLT, SignatureFormCBAdES, false, SignatureProfileBaselineLT, false},
		{SignatureLevelCBAdESBaselineLTA, SignatureFormCBAdES, false, SignatureProfileBaselineLTA, false},

		{SignatureLevelUnknown, "", true, SignatureProfileNotETSI, false},
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
			if got, err := GetSignatureLevel(c.form, c.prof); err != nil || got != SignatureLevelPKCS7B {
				t.Errorf("GetSignatureLevel(%v, %v) = %v, %v; want %v, nil (canonical first match)", c.form, c.prof, got, err, SignatureLevelPKCS7B)
			}
		}
	}
	if _, err := SignatureLevelValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
	// Upstream throws UnsupportedOperationException rather than returning null here: the
	// scan reaches UNKNOWN, whose null signature form makes getSignatureForm() throw.
	if got, err := GetSignatureLevel(SignatureFormXAdES, SignatureProfile("bogus")); err == nil {
		t.Errorf("GetSignatureLevel(unmatched) = %v, nil; want an error", got)
	}
	if got, err := GetSignatureLevel("", ""); err == nil {
		t.Errorf("GetSignatureLevel(\"\", \"\") = %v, nil; want an error", got)
	}
}

func TestSignatureLevel_ValueByNameAndString(t *testing.T) {
	got, err := SignatureLevelValueByName("XAdES-BASELINE-B")
	if err != nil || got != SignatureLevelXAdESBaselineB {
		t.Errorf("SignatureLevelValueByName(XAdES-BASELINE-B) = %v, %v; want %v, nil", got, err, SignatureLevelXAdESBaselineB)
	}
	if _, err := SignatureLevelValueByName("NOPE-NOPE"); err == nil {
		t.Error("expected error for unknown dashed name")
	}
	if got := SignatureLevelXAdESBaselineB.String(); got != "XAdES-BASELINE-B" {
		t.Errorf("SignatureLevelXAdESBaselineB.String() = %q, want %q", got, "XAdES-BASELINE-B")
	}
}
