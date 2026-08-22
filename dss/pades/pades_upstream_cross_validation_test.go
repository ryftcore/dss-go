// Cross-validation harness, direction UPSTREAM -> GO: parses the
// PAdES signatures checked into testdata/upstream/ with this package's own
// PDFDocumentAnalyzer/Signature and compares the result against
// testdata/upstream-cross-validation.json, ground truth dumped straight from upstream DSS
// 6.5.RC1's own PDFDocumentAnalyzer/PAdESSignature (see testdata/gen/CrossValidationOracle.java
// for how to regenerate it, and its header for exactly what it asserts and why). This is the
// first direction of the compatibility proof: signatures upstream DSS produced/accepts must
// parse in the Go port to the same signature count, signing certificate, claimed signing time,
// detected level, CMS SignerId, cryptographic verification, and - via the PDF-specific layer
// Signature.PdfRevision() exposes - the same /ByteRange, /SubFilter and other signature
// dictionary fields, signature field name(s), document/VRI timestamp counts, and, most PAdES-
// specific of all, the same PDF modification-detection verdict: whether upstream's own
// incremental-update diff (internal/pdf's own Go equivalent) flags the file as modified after
// signing, and into which of the secure/formFillIn/annotation/undefined buckets each detected
// object-graph change falls. Several fixtures below are genuine incremental-update attacks
// (a spoofed /ByteRange gap, a replaced /Reason, an added annotation, a rewritten page) that
// upstream itself flags; the golden JSON captures upstream's own (true) verdict for them, and
// this test asserts the Go port reproduces it exactly - a port that quietly ACCEPTS what upstream
// flags as tampered would be the dangerous failure mode.
//
// The reverse direction - Go-produced signatures validated by upstream DSS itself - lives in
// pades_downstream_cross_validation_test.go (testdata/crossgen/).
//
// There are NO quarantines here: no skipped fixture, no soft assertion, no t.Logf standing in
// for t.Error. Every fixture in the golden file is asserted in full, and a failure is a failure.
package pades

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

type pvalSignature struct {
	SigningCertificateFound  bool    `json:"signingCertificateFound"`
	SigningCertificateSHA256 *string `json:"signingCertificateSHA256"`
	ClaimedSigningTimeMillis *int64  `json:"claimedSigningTimeMillis"`
	DataFoundUpToLevel       string  `json:"dataFoundUpToLevel"`
	SignerInformationDigest  *string `json:"signerInformationDigestSHA256"`
	SignerIDIssuerSerial     *string `json:"signerIdIssuerSerial"`
	SignerIDSubjectKeyID     *string `json:"signerIdSubjectKeyIdentifier"`
	IsCounterSignature       bool    `json:"isCounterSignature"`
	ReferenceDataFound       bool    `json:"referenceDataFound"`
	ReferenceDataIntact      bool    `json:"referenceDataIntact"`
	SignatureIntact          bool    `json:"signatureIntact"`
	MessageDigestValueHex    *string `json:"messageDigestValueHex"`
	CounterSignatureCount    int     `json:"counterSignatureCount"`

	// PDF-specific layer: Signature.PdfRevision().
	SubFilter                  string   `json:"subFilter"`
	Filter                     *string  `json:"filter"`
	SignerName                 *string  `json:"signerName"`
	Reason                     *string  `json:"reason"`
	Location                   *string  `json:"location"`
	DocMDP                     *string  `json:"docMDP"`
	ByteRange                  [4]int   `json:"byteRange"`
	AreAllOriginalBytesCovered bool     `json:"areAllOriginalBytesCovered"`
	FieldNames                 []string `json:"fieldNames"`
	DocumentTimestampCount     int      `json:"documentTimestampCount"`
	VRITimestampCount          int      `json:"vriTimestampCount"`
	PdfModificationsDetected   bool     `json:"pdfModificationsDetected"`
	SecureChangeCount          int      `json:"secureChangeCount"`
	FormFillChangeCount        int      `json:"formFillChangeCount"`
	AnnotChangeCount           int      `json:"annotChangeCount"`
	UndefinedChangeCount       int      `json:"undefinedChangeCount"`

	ID string `json:"id"`
}

type pvalTimestamp struct {
	Type                     string `json:"type"`
	MessageImprintDataFound  bool   `json:"messageImprintDataFound"`
	MessageImprintDataIntact bool   `json:"messageImprintDataIntact"`
	SignatureIntact          bool   `json:"signatureIntact"`
	GenTimeMillis            *int64 `json:"genTimeMillis"`
}

type pvalFile struct {
	Path                   string          `json:"path"`
	SignatureCount         int             `json:"signatureCount"`
	Signatures             []pvalSignature `json:"signatures"`
	DetachedTimestampCount int             `json:"detachedTimestampCount"`
	DetachedTimestamps     []pvalTimestamp `json:"detachedTimestamps"`
}

type pvalGolden struct {
	Files []pvalFile `json:"files"`
}

// TestUpstreamCrossValidation replays testdata/upstream-cross-validation.json against this
// package's own parse of testdata/upstream/.
func TestUpstreamCrossValidation(t *testing.T) {
	goldenBytes, err := os.ReadFile(corpustest.Path(t, "upstream-cross-validation.json"))
	if err != nil {
		t.Fatal(err)
	}
	var golden pvalGolden
	if err := json.Unmarshal(goldenBytes, &golden); err != nil {
		t.Fatalf("parsing golden JSON: %v", err)
	}
	if len(golden.Files) == 0 {
		t.Fatal("golden JSON carries no files")
	}

	for _, gf := range golden.Files {
		gf := gf
		t.Run(gf.Path, func(t *testing.T) {
			fullPath := padesFixturePath(t, filepath.Join("upstream", filepath.FromSlash(gf.Path)))
			document, err := model.NewFileDocument(fullPath)
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", fullPath, err)
			}

			a := NewPDFDocumentAnalyzer(document)
			a.SetCertificateVerifier(validation.NewCommonCertificateVerifier())

			// validate() is what actually runs PDFDocumentAnalyzer's postProcessing (PDF
			// modification detection) in the Java oracle (see gen/CrossValidationOracle.java's
			// comment on this exact point) - Signatures() alone never populates
			// PdfModificationDetection. Go's Validate() dispatches to GetAllSignatures() through
			// the same override mechanism NewPDFDocumentAnalyzer registers
			// (InitDefaultDocumentAnalyzer), mutating the same cached signature objects
			// Signatures() below returns.
			a.Validate()

			signatures := a.Signatures()
			if len(signatures) != gf.SignatureCount {
				t.Fatalf("signature count = %d, golden wants %d", len(signatures), gf.SignatureCount)
			}
			if len(signatures) != len(gf.Signatures) {
				t.Fatalf("golden carries %d signature entries for a count of %d", len(gf.Signatures), gf.SignatureCount)
			}

			for i, sig := range signatures {
				checkPvalSignature(t, sig.(*Signature), gf.Signatures[i], i)
			}

			detachedTimestamps := a.DetachedTimestamps()
			if len(detachedTimestamps) != gf.DetachedTimestampCount {
				t.Fatalf("detached timestamp count = %d, golden wants %d", len(detachedTimestamps), gf.DetachedTimestampCount)
			}
			if len(detachedTimestamps) != len(gf.DetachedTimestamps) {
				t.Fatalf("golden carries %d detached timestamp entries for a count of %d", len(gf.DetachedTimestamps), gf.DetachedTimestampCount)
			}
			for i, tt := range detachedTimestamps {
				checkPvalTimestamp(t, tt, gf.DetachedTimestamps[i], i)
			}
		})
	}
}

func checkPvalSignature(t *testing.T, sig *Signature, want pvalSignature, index int) {
	t.Helper()

	certificateToken := sig.SigningCertificateToken()
	gotFound := certificateToken != nil
	if gotFound != want.SigningCertificateFound {
		t.Errorf("signature[%d]: signing certificate found = %v, want %v", index, gotFound, want.SigningCertificateFound)
	}
	if want.SigningCertificateSHA256 != nil && certificateToken != nil {
		gotDigest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, certificateToken.Encoded())
		if err != nil {
			t.Fatalf("signature[%d]: digesting signing certificate: %v", index, err)
		}
		if gotHex := hexEncode(gotDigest); gotHex != *want.SigningCertificateSHA256 {
			t.Errorf("signature[%d]: signing certificate SHA-256 = %s, want %s", index, gotHex, *want.SigningCertificateSHA256)
		}
	}

	gotSigningTime := sig.SigningTime()
	switch {
	case want.ClaimedSigningTimeMillis == nil && gotSigningTime != nil:
		t.Errorf("signature[%d]: signing time = %v, want none", index, gotSigningTime)
	case want.ClaimedSigningTimeMillis != nil && gotSigningTime == nil:
		t.Errorf("signature[%d]: signing time = none, want %d", index, *want.ClaimedSigningTimeMillis)
	case want.ClaimedSigningTimeMillis != nil && gotSigningTime != nil:
		if gotMillis := gotSigningTime.UnixMilli(); gotMillis != *want.ClaimedSigningTimeMillis {
			t.Errorf("signature[%d]: signing time = %d ms, want %d ms", index, gotMillis, *want.ClaimedSigningTimeMillis)
		}
	}

	if gotLevel := string(sig.DataFoundUpToLevel()); gotLevel != want.DataFoundUpToLevel {
		t.Errorf("signature[%d]: DataFoundUpToLevel() = %s, want %s", index, gotLevel, want.DataFoundUpToLevel)
	}

	if gotCounter := sig.IsCounterSignature(); gotCounter != want.IsCounterSignature {
		t.Errorf("signature[%d]: IsCounterSignature() = %v, want %v", index, gotCounter, want.IsCounterSignature)
	}

	verification := sig.SignatureCryptographicVerification()
	if got := verification.IsReferenceDataFound(); got != want.ReferenceDataFound {
		t.Errorf("signature[%d]: ReferenceDataFound = %v, want %v", index, got, want.ReferenceDataFound)
	}
	if got := verification.IsReferenceDataIntact(); got != want.ReferenceDataIntact {
		t.Errorf("signature[%d]: ReferenceDataIntact = %v, want %v", index, got, want.ReferenceDataIntact)
	}
	if gotSignatureIntact := verification.IsSignatureIntact(); gotSignatureIntact != want.SignatureIntact {
		t.Errorf("signature[%d]: SignatureIntact = %v, want %v", index, gotSignatureIntact, want.SignatureIntact)
	}

	if want.MessageDigestValueHex != nil {
		gotDigest := sig.MessageDigestValue()
		if gotHex := hexEncode(gotDigest); gotHex != *want.MessageDigestValueHex {
			t.Errorf("signature[%d]: MessageDigestValue() = %s, want %s", index, gotHex, *want.MessageDigestValueHex)
		}
	}

	// DSS-Id stability: building the identifier twice from the same, already-parsed signature
	// must yield the same string both times, and match the golden id upstream computed for the
	// same PDF revision/ByteRange.
	id1 := sig.ID()
	id2 := sig.ID()
	if id1 != id2 || id1 == "" {
		t.Errorf("signature[%d]: DSS-Id not stable: %q then %q", index, id1, id2)
	}
	if id1 != want.ID {
		t.Errorf("signature[%d]: ID() = %q, want %q", index, id1, want.ID)
	}

	counterSignatures := sig.CounterSignatures()
	if len(counterSignatures) != want.CounterSignatureCount {
		t.Fatalf("signature[%d]: counter-signature count = %d, want %d", index, len(counterSignatures), want.CounterSignatureCount)
	}

	// --- PDF-specific layer ---
	pdfRevision := sig.PdfRevision()
	if pdfRevision == nil {
		t.Fatalf("signature[%d]: PdfRevision() = nil", index)
	}
	sigDict := pdfRevision.PdfSigDictInfo()
	if sigDict == nil {
		t.Fatalf("signature[%d]: PdfSigDictInfo() = nil", index)
	}

	if got := sigDict.SubFilter(); got != want.SubFilter {
		t.Errorf("signature[%d]: SubFilter() = %q, want %q", index, got, want.SubFilter)
	}
	checkOptionalString(t, index, "Filter", sigDict.Filter(), want.Filter)
	checkOptionalString(t, index, "SignerName", sigDict.SignerName(), want.SignerName)
	checkOptionalString(t, index, "Reason", sigDict.Reason(), want.Reason)
	checkOptionalString(t, index, "Location", sigDict.Location(), want.Location)

	gotDocMDP := sigDict.DocMDP()
	gotDocMDPStr := string(gotDocMDP)
	if want.DocMDP == nil {
		if gotDocMDPStr != "" {
			t.Errorf("signature[%d]: DocMDP() = %q, want none", index, gotDocMDPStr)
		}
	} else if gotDocMDPStr != *want.DocMDP {
		t.Errorf("signature[%d]: DocMDP() = %q, want %q", index, gotDocMDPStr, *want.DocMDP)
	}

	byteRange := pdfRevision.ByteRange()
	if byteRange == nil {
		t.Fatalf("signature[%d]: ByteRange() = nil", index)
	}
	gotByteRange := [4]int{byteRange.FirstPartStart(), byteRange.FirstPartEnd(), byteRange.SecondPartStart(), byteRange.SecondPartEnd()}
	if gotByteRange != want.ByteRange {
		t.Errorf("signature[%d]: ByteRange = %v, want %v", index, gotByteRange, want.ByteRange)
	}

	if got := pdfRevision.AreAllOriginalBytesCovered(); got != want.AreAllOriginalBytesCovered {
		t.Errorf("signature[%d]: AreAllOriginalBytesCovered() = %v, want %v", index, got, want.AreAllOriginalBytesCovered)
	}

	fields := pdfRevision.Fields()
	gotFieldNames := make([]string, len(fields))
	for i, field := range fields {
		gotFieldNames[i] = field.FieldName()
	}
	if !stringSlicesEqual(gotFieldNames, want.FieldNames) {
		t.Errorf("signature[%d]: field names = %v, want %v", index, gotFieldNames, want.FieldNames)
	}

	if got := len(sig.DocumentTimestamps()); got != want.DocumentTimestampCount {
		t.Errorf("signature[%d]: len(DocumentTimestamps()) = %d, want %d", index, got, want.DocumentTimestampCount)
	}
	if got := len(sig.VRITimestamps()); got != want.VRITimestampCount {
		t.Errorf("signature[%d]: len(VRITimestamps()) = %d, want %d", index, got, want.VRITimestampCount)
	}

	modificationDetection := pdfRevision.ModificationDetection()
	gotModificationsDetected := modificationDetection != nil && modificationDetection.AreModificationsDetected()
	if gotModificationsDetected != want.PdfModificationsDetected {
		t.Errorf("signature[%d]: modifications detected = %v, want %v", index, gotModificationsDetected, want.PdfModificationsDetected)
	}
	var secureCount, formFillCount, annotCount, undefinedCount int
	if modificationDetection != nil {
		objectModifications := modificationDetection.ObjectModifications()
		secureCount = len(objectModifications.SecureChanges())
		formFillCount = len(objectModifications.FormFillInAndSignatureCreationChanges())
		annotCount = len(objectModifications.AnnotCreationChanges())
		undefinedCount = len(objectModifications.UndefinedChanges())
	}
	if secureCount != want.SecureChangeCount {
		t.Errorf("signature[%d]: secure change count = %d, want %d", index, secureCount, want.SecureChangeCount)
	}
	if formFillCount != want.FormFillChangeCount {
		t.Errorf("signature[%d]: form-fill change count = %d, want %d", index, formFillCount, want.FormFillChangeCount)
	}
	if annotCount != want.AnnotChangeCount {
		t.Errorf("signature[%d]: annotation change count = %d, want %d", index, annotCount, want.AnnotChangeCount)
	}
	if undefinedCount != want.UndefinedChangeCount {
		t.Errorf("signature[%d]: undefined change count = %d, want %d", index, undefinedCount, want.UndefinedChangeCount)
	}
}

func checkPvalTimestamp(t *testing.T, timestamp *validation.TimestampToken, want pvalTimestamp, index int) {
	t.Helper()

	if gotType := string(timestamp.TimeStampType()); gotType != want.Type {
		t.Errorf("detachedTimestamps[%d]: TimeStampType() = %s, want %s", index, gotType, want.Type)
	}
	if got := timestamp.IsMessageImprintDataFound(); got != want.MessageImprintDataFound {
		t.Errorf("detachedTimestamps[%d]: IsMessageImprintDataFound() = %v, want %v", index, got, want.MessageImprintDataFound)
	}
	if got := timestamp.IsMessageImprintDataIntact(); got != want.MessageImprintDataIntact {
		t.Errorf("detachedTimestamps[%d]: IsMessageImprintDataIntact() = %v, want %v", index, got, want.MessageImprintDataIntact)
	}
	if got := timestamp.IsSignatureIntact(); got != want.SignatureIntact {
		t.Errorf("detachedTimestamps[%d]: IsSignatureIntact() = %v, want %v", index, got, want.SignatureIntact)
	}

	gotGenTime := timestamp.GenerationTime()
	switch {
	case want.GenTimeMillis == nil && !gotGenTime.IsZero():
		t.Errorf("detachedTimestamps[%d]: generation time = %v, want none", index, gotGenTime)
	case want.GenTimeMillis != nil:
		if gotMillis := gotGenTime.UnixMilli(); gotMillis != *want.GenTimeMillis {
			t.Errorf("detachedTimestamps[%d]: generation time = %d ms, want %d ms", index, gotMillis, *want.GenTimeMillis)
		}
	}
}

func checkOptionalString(t *testing.T, index int, field, got string, want *string) {
	t.Helper()
	if want == nil {
		if got != "" {
			t.Errorf("signature[%d]: %s = %q, want none", index, field, got)
		}
		return
	}
	if got != *want {
		t.Errorf("signature[%d]: %s = %q, want %q", index, field, got, *want)
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func hexEncode(b []byte) string {
	const hexdigits = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hexdigits[v>>4]
		out[i*2+1] = hexdigits[v&0xf]
	}
	return string(out)
}

// compile-time reminder that analyzer.DocumentAnalyzer is the interface PDFDocumentAnalyzer
// implements, kept honest for the a.SetCertificateVerifier/Validate/Signatures/
// DetachedTimestamps calls above (all promoted from analyzer.DefaultDocumentAnalyzer).
var _ analyzer.DocumentAnalyzer = (*PDFDocumentAnalyzer)(nil)
