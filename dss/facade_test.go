// Regression tests for the facade: every format signs, round-trips through
// Validate, and comes back with the level and verdict it should. The examples
// in example_test.go document the API; this file is the coverage.
package dss_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
)

func testSigner(t *testing.T) *dss.Signer {
	t.Helper()
	signer, err := dss.OpenPKCS12("testdata/signer_rsa.p12", "testpassword")
	if err != nil {
		t.Fatalf("OpenPKCS12: %v", err)
	}
	t.Cleanup(signer.Close)
	return signer
}

func testTSA(t *testing.T) dss.TSPSource {
	t.Helper()
	tsa, err := spivalidation.NewKeyEntityTSPSourceFromKeyStorePath(
		"testdata/tsa_ec.p12", "PKCS12", "testpassword", "", "testpassword")
	if err != nil {
		t.Fatalf("KeyEntityTSPSource: %v", err)
	}
	tsa.SetTsaPolicy("1.2.3.4.5.6.7.8.9")
	return tsa
}

// TestSignValidateRoundTrip covers every format the facade offers, at level B
// and - where a local TSA is enough - at level T, and asserts that the port's
// own validator recognises the level and passes the signature once the test
// root is trusted.
func TestSignValidateRoundTrip(t *testing.T) {
	signer := testSigner(t)
	tsa := testTSA(t)

	pdf, err := dss.OpenDocument("testdata/sample.pdf")
	if err != nil {
		t.Fatalf("OpenDocument: %v", err)
	}
	xmlDocument := dss.NewDocument("invoice.xml", []byte("<invoice><total>42</total></invoice>"))
	binaryDocument := dss.NewDocument("payload.bin", []byte("payload"))
	jsonDocument := dss.NewDocument("payload.json", []byte(`{"amount":42}`))

	cases := []struct {
		name      string
		documents []dss.Document
		options   dss.SignOptions
		wantLevel dss.SignatureLevel
		wantName  string
	}{
		{"pades-b", []dss.Document{pdf},
			dss.SignOptions{Format: dss.FormatPAdES, Level: dss.LevelB},
			"PAdES_BASELINE_B", "sample-signed-pades-baseline-b.pdf"},
		{"xades-b-enveloped", []dss.Document{xmlDocument},
			dss.SignOptions{Format: dss.FormatXAdES, Level: dss.LevelB},
			"XAdES_BASELINE_B", "invoice-signed-xades-baseline-b.xml"},
		{"cades-b-enveloping", []dss.Document{binaryDocument},
			dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelB},
			"CAdES_BASELINE_B", "payload-signed-cades-baseline-b.p7m"},
		{"cades-b-detached", []dss.Document{binaryDocument},
			dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelB, Packaging: dss.PackagingDetached},
			"CAdES_BASELINE_B", "payload-signed-cades-baseline-b.p7s"},
		{"jades-b-compact", []dss.Document{jsonDocument},
			dss.SignOptions{Format: dss.FormatJAdES, Level: dss.LevelB},
			"JAdES_BASELINE_B", "payload-signed-jades-baseline-b.json"},
		{"asice-cades-b", []dss.Document{binaryDocument, jsonDocument},
			dss.SignOptions{Format: dss.FormatASiCWithCAdES, Level: dss.LevelB},
			"CAdES_BASELINE_B", "container-signed-cades-baseline-b.sce"},
		{"asice-xades-b", []dss.Document{binaryDocument, jsonDocument},
			dss.SignOptions{Format: dss.FormatASiCWithXAdES, Level: dss.LevelB},
			"XAdES_BASELINE_B", "container-signed-xades-baseline-b.sce"},
		{"asics-xades-b", []dss.Document{binaryDocument},
			dss.SignOptions{Format: dss.FormatASiCWithXAdES, Level: dss.LevelB},
			"XAdES_BASELINE_B", "container-signed-xades-baseline-b.scs"},
		{"pades-t", []dss.Document{pdf},
			dss.SignOptions{Format: dss.FormatPAdES, Level: dss.LevelT, TSPSource: tsa},
			"PAdES_BASELINE_T", "sample-signed-pades-baseline-t.pdf"},
		{"cades-t", []dss.Document{binaryDocument},
			dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelT, TSPSource: tsa},
			"CAdES_BASELINE_T", "payload-signed-cades-baseline-t.p7m"},
		{"xades-t", []dss.Document{xmlDocument},
			dss.SignOptions{Format: dss.FormatXAdES, Level: dss.LevelT, TSPSource: tsa},
			"XAdES_BASELINE_T", "invoice-signed-xades-baseline-t.xml"},
		{"jades-t-flattened", []dss.Document{jsonDocument},
			dss.SignOptions{Format: dss.FormatJAdES, Level: dss.LevelT, TSPSource: tsa,
				JWSSerialization: dss.JWSFlattenedJSON},
			"JAdES_BASELINE_T", "payload-signed-jades-baseline-t.json"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			signed, err := dss.SignMultiple(testCase.documents, signer, testCase.options)
			if err != nil {
				t.Fatalf("SignMultiple: %v", err)
			}
			if signed.Name() != testCase.wantName {
				t.Errorf("document name = %q, want %q", signed.Name(), testCase.wantName)
			}

			var detached []dss.Document
			if testCase.options.Packaging == dss.PackagingDetached {
				detached = testCase.documents
			}
			reports, err := dss.Validate(signed, dss.ValidateOptions{
				DetachedContents:    detached,
				TrustedCertificates: signer.CertificateChain(),
			})
			if err != nil {
				t.Fatalf("Validate: %v", err)
			}
			verdicts := reports.Verdicts()
			if len(verdicts) != 1 {
				t.Fatalf("got %d signatures, want 1", len(verdicts))
			}
			if verdicts[0].Indication != dss.IndicationTotalPassed {
				t.Errorf("indication = %s (%s), want TOTAL_PASSED",
					verdicts[0].Indication, verdicts[0].SubIndication)
			}
			if verdicts[0].SignatureLevel != testCase.wantLevel {
				t.Errorf("level = %s, want %s", verdicts[0].SignatureLevel, testCase.wantLevel)
			}
			if !reports.Valid() {
				t.Error("Reports.Valid() = false, want true")
			}
		})
	}
}

// TestExtend checks that a B-level signature can be raised to T without the
// signing key, and that the validator sees the new level.
func TestExtend(t *testing.T) {
	signer := testSigner(t)

	document := dss.NewDocument("payload.bin", []byte("payload"))
	signed, err := dss.Sign(document, signer, dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelB})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	extended, err := dss.Extend(signed, dss.ExtendOptions{
		Format:    dss.FormatCAdES,
		Level:     dss.LevelT,
		TSPSource: testTSA(t),
	})
	if err != nil {
		t.Fatalf("Extend: %v", err)
	}

	reports, err := dss.Validate(extended, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	verdicts := reports.Verdicts()
	if len(verdicts) != 1 || verdicts[0].SignatureLevel != "CAdES_BASELINE_T" {
		t.Fatalf("verdicts = %+v, want one CAdES_BASELINE_T", verdicts)
	}
	if verdicts[0].BestSignatureTime == nil {
		t.Error("best signature time is nil, but a T-level signature carries a time-stamp")
	}
}

// TestValidateUntrusted pins the verdict for a chain that reaches no trust
// anchor: INDETERMINATE, not a failure and not a pass.
func TestValidateUntrusted(t *testing.T) {
	signer := testSigner(t)
	document := dss.NewDocument("invoice.xml", []byte("<invoice/>"))
	signed, err := dss.Sign(document, signer, dss.SignOptions{Format: dss.FormatXAdES, Level: dss.LevelB})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	reports, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	verdict := reports.Verdicts()[0]
	if verdict.Indication != dss.IndicationIndeterminate {
		t.Errorf("indication = %s, want INDETERMINATE", verdict.Indication)
	}
	if verdict.Valid() || reports.Valid() {
		t.Error("an unanchored signature must not report as valid")
	}
	if len(verdict.Errors) == 0 {
		t.Error("expected at least one AdES validation error explaining the indication")
	}
}

// TestValidateReportsMarshal checks the four report marshallers all produce
// their document.
func TestValidateReportsMarshal(t *testing.T) {
	signer := testSigner(t)
	document := dss.NewDocument("invoice.xml", []byte("<invoice/>"))
	signed, err := dss.Sign(document, signer, dss.SignOptions{Format: dss.FormatXAdES, Level: dss.LevelB})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	reports, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}

	marshallers := map[string]func() (string, error){
		"simple":     reports.SimpleReportXML,
		"detailed":   reports.DetailedReportXML,
		"diagnostic": reports.DiagnosticDataXML,
		"etsi":       reports.ETSIValidationReportXML,
	}
	for name, marshal := range marshallers {
		xml, err := marshal()
		if err != nil {
			t.Errorf("%s report: %v", name, err)
			continue
		}
		if len(xml) == 0 {
			t.Errorf("%s report is empty", name)
		}
	}
}

// TestSignOptionValidation covers the rejections the facade performs before it
// reaches the ported services.
func TestSignOptionValidation(t *testing.T) {
	signer := testSigner(t)
	document := dss.NewDocument("payload.bin", []byte("payload"))

	cases := []struct {
		name      string
		documents []dss.Document
		options   dss.SignOptions
		want      error
	}{
		{"unknown format", []dss.Document{document},
			dss.SignOptions{Format: "PGP", Level: dss.LevelB}, dss.ErrUnsupportedFormat},
		{"unknown level", []dss.Document{document},
			dss.SignOptions{Format: dss.FormatCAdES, Level: "X"}, dss.ErrUnsupportedLevel},
		{"missing TSP source", []dss.Document{document},
			dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelT}, dss.ErrTSPSourceRequired},
		{"several documents, single-document format", []dss.Document{document, document},
			dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelB}, dss.ErrMultipleDocuments},
		{"no document", nil,
			dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelB}, dss.ErrNoDocument},
		{"password protection, non-PDF format", []dss.Document{document},
			dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelB, PasswordProtection: []byte("secret")},
			dss.ErrPasswordProtectionNotApplicable},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if _, err := dss.SignMultiple(testCase.documents, signer, testCase.options); !errors.Is(err, testCase.want) {
				t.Errorf("error = %v, want %v", err, testCase.want)
			}
		})
	}
}

// TestExtendOptionValidation covers Extend's own rejections.
func TestExtendOptionValidation(t *testing.T) {
	document := dss.NewDocument("payload.bin", []byte("payload"))

	if _, err := dss.Extend(document, dss.ExtendOptions{Format: "PGP", Level: dss.LevelT}); !errors.Is(err, dss.ErrUnsupportedFormat) {
		t.Errorf("error = %v, want ErrUnsupportedFormat", err)
	}
	if _, err := dss.Extend(document, dss.ExtendOptions{Format: dss.FormatCAdES, Level: dss.LevelT}); !errors.Is(err, dss.ErrTSPSourceRequired) {
		t.Errorf("error = %v, want ErrTSPSourceRequired", err)
	}
	if _, err := dss.Extend(nil, dss.ExtendOptions{Format: dss.FormatCAdES, Level: dss.LevelT}); !errors.Is(err, dss.ErrNoDocument) {
		t.Errorf("error = %v, want ErrNoDocument", err)
	}
	if _, err := dss.Extend(document, dss.ExtendOptions{
		Format: dss.FormatCAdES, Level: dss.LevelT, TSPSource: testTSA(t), PasswordProtection: []byte("secret"),
	}); !errors.Is(err, dss.ErrPasswordProtectionNotApplicable) {
		t.Errorf("error = %v, want ErrPasswordProtectionNotApplicable", err)
	}
}

// TestProtectedPDF covers PasswordProtection on all three option types
// against upstream's dss-pades/src/test/resources/protected/open_protected.pdf
// (AES-128, user password a single space): the password opens the document for
// signing, validation and extension; without it, or with the wrong one, each
// of them fails to open the document - an error, never a verdict. The signed
// and extended documents stay encrypted under the same password.
func TestProtectedPDF(t *testing.T) {
	signer := testSigner(t)
	password := []byte(" ")

	protected, err := dss.OpenDocument("testdata/upstream/protected/open_protected.pdf")
	if err != nil {
		t.Fatalf("OpenDocument: %v", err)
	}

	// wantInvalidPassword fails the test unless err is the "cannot open the
	// encrypted document" error - not any error, which a broken fixture or
	// a different rejection would also produce.
	wantInvalidPassword := func(step string, err error) {
		t.Helper()
		if err == nil || !strings.Contains(err.Error(), "invalid password") {
			t.Fatalf("%s: err = %v, want an invalid-password error", step, err)
		}
	}

	_, err = dss.Sign(protected, signer, dss.SignOptions{Format: dss.FormatPAdES, Level: dss.LevelB})
	wantInvalidPassword("Sign without the password", err)
	_, err = dss.Sign(protected, signer, dss.SignOptions{
		Format: dss.FormatPAdES, Level: dss.LevelB, PasswordProtection: []byte("wrong"),
	})
	wantInvalidPassword("Sign with the wrong password", err)

	signed, err := dss.Sign(protected, signer, dss.SignOptions{
		Format: dss.FormatPAdES, Level: dss.LevelB, PasswordProtection: password,
	})
	if err != nil {
		t.Fatalf("Sign with the password: %v", err)
	}
	if signed.Name() != "open-protected-signed-pades-baseline-b.pdf" {
		t.Errorf("signed name = %q, want open-protected-signed-pades-baseline-b.pdf", signed.Name())
	}

	// Still encrypted: the result does not open without the password.
	_, err = dss.Validate(signed, dss.ValidateOptions{})
	wantInvalidPassword("Validate without the password", err)
	reports, err := dss.Validate(signed, dss.ValidateOptions{
		PasswordProtection:  password,
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		t.Fatalf("Validate with the password: %v", err)
	}
	verdicts := reports.Verdicts()
	if len(verdicts) != 1 || verdicts[0].SignatureLevel != "PAdES_BASELINE_B" || !verdicts[0].Valid() {
		t.Fatalf("verdicts = %+v, want one valid PAdES_BASELINE_B", verdicts)
	}

	_, err = dss.Extend(signed, dss.ExtendOptions{
		Format: dss.FormatPAdES, Level: dss.LevelT, TSPSource: testTSA(t),
	})
	wantInvalidPassword("Extend without the password", err)
	extended, err := dss.Extend(signed, dss.ExtendOptions{
		Format: dss.FormatPAdES, Level: dss.LevelT, TSPSource: testTSA(t), PasswordProtection: password,
	})
	if err != nil {
		t.Fatalf("Extend with the password: %v", err)
	}
	_, err = dss.Validate(extended, dss.ValidateOptions{})
	wantInvalidPassword("Validate extended without the password", err)
	reports, err = dss.Validate(extended, dss.ValidateOptions{
		PasswordProtection:  password,
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		t.Fatalf("Validate extended with the password: %v", err)
	}
	verdicts = reports.Verdicts()
	if len(verdicts) != 1 || verdicts[0].SignatureLevel != "PAdES_BASELINE_T" || !verdicts[0].Valid() {
		t.Fatalf("extended verdicts = %+v, want one valid PAdES_BASELINE_T", verdicts)
	}

	// A password for anything but a PDF is refused, since Validate detects
	// the format itself and only the PDF validator can take one.
	xmlDocument := dss.NewDocument("invoice.xml", []byte("<invoice/>"))
	xmlSigned, err := dss.Sign(xmlDocument, signer, dss.SignOptions{Format: dss.FormatXAdES, Level: dss.LevelB})
	if err != nil {
		t.Fatalf("Sign XAdES: %v", err)
	}
	if _, err := dss.Validate(xmlSigned, dss.ValidateOptions{PasswordProtection: password}); !errors.Is(err, dss.ErrPasswordProtectionNotApplicable) {
		t.Errorf("Validate XML with a password: error = %v, want ErrPasswordProtectionNotApplicable", err)
	}

	// A PDF that is not encrypted opens whatever password it is given:
	// upstream pdfbox only consults the password when the trailer carries
	// /Encrypt, and so does the port. A password on a plain PDF is
	// therefore not an error.
	plain, err := dss.OpenDocument("testdata/sample.pdf")
	if err != nil {
		t.Fatalf("OpenDocument: %v", err)
	}
	if _, err := dss.Sign(plain, signer, dss.SignOptions{
		Format: dss.FormatPAdES, Level: dss.LevelB, PasswordProtection: []byte("ignored"),
	}); err != nil {
		t.Errorf("Sign of an unencrypted PDF with a password: %v, want it ignored as upstream does", err)
	}
}

// TestValidateRejectsSuiteWithoutPolicy pins the one option combination the
// underlying entry points cannot express.
func TestValidateRejectsSuiteWithoutPolicy(t *testing.T) {
	document := dss.NewDocument("payload.bin", []byte("payload"))
	_, err := dss.Validate(document, dss.ValidateOptions{
		CryptographicSuite: dss.NewDocument("suite.xml", []byte("<x/>")),
	})
	if err == nil {
		t.Fatal("expected an error when a cryptographic suite is given without a policy")
	}
}

// TestValidateUnsupportedDocument checks that an unrecognisable document is
// reported as an error rather than as a panic escaping the facade.
func TestValidateUnsupportedDocument(t *testing.T) {
	if _, err := dss.Validate(dss.NewDocument("noise.bin", []byte("not a signature")), dss.ValidateOptions{}); err == nil {
		t.Fatal("expected an error for a document carrying no signature")
	}
	if _, err := dss.Validate(nil, dss.ValidateOptions{}); !errors.Is(err, dss.ErrNoDocument) {
		t.Errorf("error = %v, want ErrNoDocument", err)
	}
}

// TestBaselineLevelTable checks the format/level mapping exhaustively, since
// every signature the facade produces depends on it.
func TestBaselineLevelTable(t *testing.T) {
	want := map[dss.Format]map[dss.Level]dss.SignatureLevel{
		dss.FormatCAdES: {
			dss.LevelB: "CAdES_BASELINE_B", dss.LevelT: "CAdES_BASELINE_T",
			dss.LevelLT: "CAdES_BASELINE_LT", dss.LevelLTA: "CAdES_BASELINE_LTA",
		},
		dss.FormatXAdES: {
			dss.LevelB: "XAdES_BASELINE_B", dss.LevelT: "XAdES_BASELINE_T",
			dss.LevelLT: "XAdES_BASELINE_LT", dss.LevelLTA: "XAdES_BASELINE_LTA",
		},
		dss.FormatPAdES: {
			dss.LevelB: "PAdES_BASELINE_B", dss.LevelT: "PAdES_BASELINE_T",
			dss.LevelLT: "PAdES_BASELINE_LT", dss.LevelLTA: "PAdES_BASELINE_LTA",
		},
		dss.FormatJAdES: {
			dss.LevelB: "JAdES_BASELINE_B", dss.LevelT: "JAdES_BASELINE_T",
			dss.LevelLT: "JAdES_BASELINE_LT", dss.LevelLTA: "JAdES_BASELINE_LTA",
		},
		dss.FormatASiCWithCAdES: {
			dss.LevelB: "CAdES_BASELINE_B", dss.LevelT: "CAdES_BASELINE_T",
			dss.LevelLT: "CAdES_BASELINE_LT", dss.LevelLTA: "CAdES_BASELINE_LTA",
		},
		dss.FormatASiCWithXAdES: {
			dss.LevelB: "XAdES_BASELINE_B", dss.LevelT: "XAdES_BASELINE_T",
			dss.LevelLT: "XAdES_BASELINE_LT", dss.LevelLTA: "XAdES_BASELINE_LTA",
		},
	}
	for format, levels := range want {
		for level, wantSignatureLevel := range levels {
			got, err := format.BaselineLevel(level)
			if err != nil {
				t.Errorf("%s/%s: %v", format, level, err)
				continue
			}
			if got != wantSignatureLevel {
				t.Errorf("%s/%s = %s, want %s", format, level, string(got), string(wantSignatureLevel))
			}
		}
	}
}

// TestFacadeEntryPointsDoNotPanic pins that inputs the ported code asserts
// on - a nil certificate buffer, nil document content, a zero Signer - come
// back as errors (or, for Close, as a no-op) rather than as a panic escaping
// the facade.
func TestFacadeEntryPointsDoNotPanic(t *testing.T) {
	if _, err := dss.LoadCertificateBytes(nil); err == nil {
		t.Error("LoadCertificateBytes(nil): expected an error")
	}
	if _, err := dss.Validate(dss.NewDocument("empty.bin", nil), dss.ValidateOptions{}); err == nil {
		t.Error("Validate of an empty document: expected an error")
	}
	document := dss.NewDocument("payload.xml", []byte("<a/>"))
	if _, err := dss.Sign(document, &dss.Signer{}, dss.SignOptions{Format: dss.FormatXAdES, Level: dss.LevelB}); err == nil {
		t.Error("Sign with a zero Signer: expected an error")
	}
	var signer *dss.Signer
	signer.Close()
}
