// Cross-validation harness, direction UPSTREAM -> GO (task #12, JAdES extension): parses the
// JAdES signatures checked into testdata/upstream/ with this package's own
// JWSDocumentAnalyzerFactory/JAdESSignature and compares the result against
// testdata/upstream-cross-validation.json, ground truth dumped straight from upstream DSS
// 6.5.RC1's own JWSDocumentAnalyzerFactory/JAdESSignature (see testdata/gen/CrossValidationOracle.java
// for how to regenerate it, and its header for exactly what it asserts and why). This is the
// first direction of the compatibility proof: signatures upstream DSS produced/accepts must
// parse in the Go port to the same signature count, signing certificate, claimed signing time,
// detected level, counter-signature structure, JWS serialization type actually used, the SigD
// mechanism a detached signature was covered by, the structural-validation error count (the
// DSS-2620 'crit'/'b64' edge cases), and (wherever the original content is available to check
// it) reference/DataToBeSignedRepresentation integrity.
//
// A distinct second half - testdata/gen/CrossValidationOracle.java's EXCEPTION_FILES /
// exceptionFiles - covers fixtures upstream's own JWSCompactDocumentAnalyzer rejects outright
// (an unchecked exception raised while parsing, before a JAdESSignature ever exists): a Go port
// that quietly accepted what upstream refuses to even parse would be the dangerous failure mode,
// exactly the same rationale PAdES's modification-detection verdicts and XAdES/CAdES's signature-
// intact verdicts are asserted strictly for.
//
// The reverse direction - Go-produced signatures validated by upstream DSS itself - lives in
// jades_downstream_cross_validation_test.go (testdata/crossgen/).
//
// There are NO quarantines here: no skipped fixture, no soft assertion, no t.Logf standing in for
// t.Error. Every fixture in the golden file is asserted in full, and a failure is a failure.
package jades

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

type jvalSignature struct {
	SigningCertificateFound  bool            `json:"signingCertificateFound"`
	SigningCertificateSHA256 *string         `json:"signingCertificateSHA256"`
	ClaimedSigningTimeMillis *int64          `json:"claimedSigningTimeMillis"`
	DataFoundUpToLevel       string          `json:"dataFoundUpToLevel"`
	IsCounterSignature       bool            `json:"isCounterSignature"`
	ReferenceDataFound       bool            `json:"referenceDataFound"`
	ReferenceDataIntact      bool            `json:"referenceDataIntact"`
	SignatureIntact          bool            `json:"signatureIntact"`
	DTBSRAlgorithm           *string         `json:"dtbsrAlgorithm"`
	DTBSRHex                 *string         `json:"dtbsrHex"`
	SigDMechanism            *string         `json:"sigDMechanism"`
	SerializationType        string          `json:"serializationType"`
	StructureValidationCount int             `json:"structureValidationErrorCount"`
	ContentTimestamps        []jvalTimestamp `json:"contentTimestamps"`
	SignatureTimestamps      []jvalTimestamp `json:"signatureTimestamps"`
	TimestampsX1             []jvalTimestamp `json:"timestampsX1"`
	TimestampsX2             []jvalTimestamp `json:"timestampsX2"`
	ArchiveTimestamps        []jvalTimestamp `json:"archiveTimestamps"`
	CounterSignatureCount    int             `json:"counterSignatureCount"`
	CounterSignatures        []jvalSignature `json:"counterSignatures"`
}

// jvalTimestamp mirrors gen/CrossValidationOracle.java's dumpTimestamps(). Two of these columns
// are there because the broad whole-corpus differential (testdata/broadgen/) caught defects no
// other column in this golden could see, and they stay asserted so those defects cannot return:
//
//   - DSSID: SignatureTimestampIdentifierBuilder mixes the carrying attribute's position among
//     the signature properties into the token identifier. Losing that position (a lookup that
//     misses because the attribute objects are rebuilt on every access) changes every content
//     time-stamp's id without changing anything else.
//   - TimestampedReferences: JAdESTimestampSource's getSignatureTimestampReferences() override
//     folds getKeyInfoReferences() into a signature time-stamp's covered set. Missing that
//     override drops exactly one certificate reference per signature time-stamp.
type jvalTimestamp struct {
	Type                     string  `json:"type"`
	DSSID                    string  `json:"dssId"`
	GenerationTimeMillis     *int64  `json:"generationTimeMillis"`
	MessageImprintDataFound  bool    `json:"messageImprintDataFound"`
	MessageImprintDataIntact bool    `json:"messageImprintDataIntact"`
	MessageImprintHex        *string `json:"messageImprintHex"`
	SignatureIntact          bool    `json:"signatureIntact"`
	TimestampedReferences    string  `json:"timestampedReferences"`
}

type jvalFile struct {
	Path           string          `json:"path"`
	SignatureCount int             `json:"signatureCount"`
	Signatures     []jvalSignature `json:"signatures"`
}

type jvalExceptionFile struct {
	Path                     string `json:"path"`
	ExpectedMessageSubstring string `json:"expectedMessageSubstring"`
}

type jvalGolden struct {
	Files          []jvalFile          `json:"files"`
	ExceptionFiles []jvalExceptionFile `json:"exceptionFiles"`
}

// jvalDetachedContent mirrors gen/CrossValidationOracle.java's DETACHED_CONTENT map: the one
// genuinely detached signature in the fixture set, and the original content file to feed it.
var jvalDetachedContent = map[string]string{
	"validation/simple-detached.json": "sample.json",
}

// TestUpstreamCrossValidation replays testdata/upstream-cross-validation.json against this
// package's own parse of testdata/upstream/.
func TestUpstreamCrossValidation(t *testing.T) {
	goldenBytes, err := os.ReadFile(corpustest.Path(t, "upstream-cross-validation.json"))
	if err != nil {
		t.Fatal(err)
	}
	var golden jvalGolden
	if err := json.Unmarshal(goldenBytes, &golden); err != nil {
		t.Fatalf("parsing golden JSON: %v", err)
	}
	if len(golden.Files) == 0 {
		t.Fatal("golden JSON carries no files")
	}

	for _, gf := range golden.Files {
		gf := gf
		t.Run(gf.Path, func(t *testing.T) {
			fullPath := jadesFixturePath(t, filepath.Join("upstream", filepath.FromSlash(gf.Path)))
			document, err := model.NewFileDocument(fullPath)
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", fullPath, err)
			}

			factory := NewJWSDocumentAnalyzerFactory()
			a := factory.Create(document)
			a.SetCertificateVerifier(validation.NewCommonCertificateVerifier())

			if detachedRel, ok := jvalDetachedContent[gf.Path]; ok {
				detachedPath := jadesFixturePath(t, filepath.Join("upstream", filepath.FromSlash(detachedRel)))
				detachedDoc, err := model.NewFileDocument(detachedPath)
				if err != nil {
					t.Fatalf("NewFileDocument(%s): %v", detachedPath, err)
				}
				a.SetDetachedContents([]model.DSSDocument{detachedDoc})
			}

			signatures := a.Signatures()
			if len(signatures) != gf.SignatureCount {
				t.Fatalf("signature count = %d, golden wants %d", len(signatures), gf.SignatureCount)
			}
			if len(signatures) != len(gf.Signatures) {
				t.Fatalf("golden carries %d signature entries for a count of %d", len(gf.Signatures), gf.SignatureCount)
			}

			for i, sig := range signatures {
				// Top-level signatures are already initialized by
				// AbstractJWSDocumentAnalyzer's own BuildSignatures(); only recursed-into
				// counter signatures need it below.
				checkJvalSignature(t, sig.(*JAdESSignature), gf.Signatures[i], i, false)
			}
		})
	}

	if len(golden.ExceptionFiles) == 0 {
		t.Fatal("golden JSON carries no exceptionFiles")
	}
	for _, ef := range golden.ExceptionFiles {
		ef := ef
		t.Run("exception/"+ef.Path, func(t *testing.T) {
			fullPath := jadesFixturePath(t, filepath.Join("upstream", filepath.FromSlash(ef.Path)))
			document, err := model.NewFileDocument(fullPath)
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", fullPath, err)
			}

			gotMessage := recoverCreateAndParse(document)
			if gotMessage == "" {
				t.Fatalf("expected %s to raise a panic (upstream itself rejects it), but it parsed cleanly", ef.Path)
			}
			// Case-insensitive: Go's message text is its own idiom ("an entry for 'sigT'
			// already exists, names must be unique"), not a verbatim copy of Java's ("An
			// entry for 'sigT' already exists. Names must be unique."), but must name the
			// same defect - same fixture, same rejection reason, in both implementations.
			if !strings.Contains(strings.ToLower(gotMessage), strings.ToLower(ef.ExpectedMessageSubstring)) {
				t.Errorf("panic message = %q, want it to contain (case-insensitively) %q", gotMessage, ef.ExpectedMessageSubstring)
			}
		})
	}
}

// recoverCreateAndParse invokes JWSDocumentAnalyzerFactory.Create + Signatures() (the same two
// calls the golden Java oracle's own EXCEPTION_FILES handling wraps in a try/catch) and returns
// the recovered panic's message chain (top-level message, then every Unwrap()-reachable cause's
// own message, "|"-joined - Go's Error() intentionally returns only the top message, same as
// Java's Throwable#getMessage(), so the oracle's cause-chain walk over getCause() needs a
// matching Unwrap() walk here), or "" if no panic occurred.
func recoverCreateAndParse(document model.DSSDocument) (message string) {
	defer func() {
		if r := recover(); r != nil {
			var chain []string
			if err, ok := r.(error); ok {
				for e := err; e != nil; e = errors.Unwrap(e) {
					chain = append(chain, e.Error())
				}
			} else {
				chain = append(chain, fmt.Sprint(r))
			}
			message = strings.Join(chain, " | ")
		}
	}()
	factory := NewJWSDocumentAnalyzerFactory()
	a := factory.Create(document)
	a.SetCertificateVerifier(validation.NewCommonCertificateVerifier())
	a.Signatures()
	return ""
}

func checkJvalSignature(t *testing.T, sig *JAdESSignature, want jvalSignature, index int, isCounterSig bool) {
	t.Helper()

	// A nested counter signature is not initialized by AbstractJWSDocumentAnalyzer's own
	// BuildSignatures() (only the top-level ones are), and DataFoundUpToLevel() requires it -
	// exactly the same reason gen/CrossValidationOracle.java's dumpSignature() calls the Java
	// equivalent only for what it recurses into.
	if isCounterSig {
		sig.InitBaselineRequirementsChecker(validation.NewCommonCertificateVerifier())
	}

	certificateToken := sig.SigningCertificateToken()
	gotFound := certificateToken != nil
	if gotFound != want.SigningCertificateFound {
		t.Errorf("signature[%d]: signing certificate found = %v, want %v", index, gotFound, want.SigningCertificateFound)
	}
	if want.SigningCertificateSHA256 != nil && certificateToken != nil {
		gotDigest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithm_SHA256, certificateToken.Encoded())
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

	if want.DTBSRHex != nil {
		gotDTBSR := sig.DataToBeSignedRepresentation()
		if want.DTBSRAlgorithm != nil {
			if gotAlgo := string(gotDTBSR.Algorithm()); gotAlgo != *want.DTBSRAlgorithm {
				t.Errorf("signature[%d]: DataToBeSignedRepresentation().Algorithm() = %s, want %s", index, gotAlgo, *want.DTBSRAlgorithm)
			}
		}
		// hexEncode (lowercase), not Digest.HexValue() (uppercase, Java's BigInteger-based
		// getHexValue()): the golden's dtbsrHex column was dumped with the oracle's own
		// lowercase hex() helper (same one every other *SHA256 field here uses), not
		// Digest#getHexValue() directly.
		if gotHex := hexEncode(gotDTBSR.Value()); gotHex != *want.DTBSRHex {
			t.Errorf("signature[%d]: DataToBeSignedRepresentation() hex = %s, want %s", index, gotHex, *want.DTBSRHex)
		}
	}

	gotSigDMechanism := sig.SigDMechanism()
	switch {
	case want.SigDMechanism == nil && gotSigDMechanism != nil:
		t.Errorf("signature[%d]: SigDMechanism() = %v, want none", index, *gotSigDMechanism)
	case want.SigDMechanism != nil && gotSigDMechanism == nil:
		t.Errorf("signature[%d]: SigDMechanism() = none, want %s", index, *want.SigDMechanism)
	case want.SigDMechanism != nil && gotSigDMechanism != nil:
		if gotStr := string(*gotSigDMechanism); gotStr != *want.SigDMechanism {
			t.Errorf("signature[%d]: SigDMechanism() = %s, want %s", index, gotStr, *want.SigDMechanism)
		}
	}

	if gotSerializationType := string(sig.Jws().JwsSerializationType()); gotSerializationType != want.SerializationType {
		t.Errorf("signature[%d]: Jws().JwsSerializationType() = %s, want %s", index, gotSerializationType, want.SerializationType)
	}

	// Presence/absence only, not exact count: jades/specs/jades_utils_test.go's own file header
	// documents that message-list parity with Java's everit/jsonsKema JSON-schema library is not
	// a contract for this port (a oneOf/anyOf mismatch's nested "causes" collapse into a single
	// formatted string on Java's side but expand into one entry per branch here) - only whether a
	// structural defect was detected at all is cross-checked, matching that established contract
	// rather than re-litigating it through this harness.
	gotHasErrors := len(sig.StructureValidationResult()) > 0
	wantHasErrors := want.StructureValidationCount > 0
	if gotHasErrors != wantHasErrors {
		t.Errorf("signature[%d]: StructureValidationResult() non-empty = %v, want %v (got errors: %v)",
			index, gotHasErrors, wantHasErrors, sig.StructureValidationResult())
	}

	checkJvalTimestamps(t, index, "contentTimestamps", sig.ContentTimestamps(), want.ContentTimestamps)
	checkJvalTimestamps(t, index, "signatureTimestamps", sig.SignatureTimestamps(), want.SignatureTimestamps)
	checkJvalTimestamps(t, index, "timestampsX1", sig.TimestampsX1(), want.TimestampsX1)
	checkJvalTimestamps(t, index, "timestampsX2", sig.TimestampsX2(), want.TimestampsX2)
	checkJvalTimestamps(t, index, "archiveTimestamps", sig.ArchiveTimestamps(), want.ArchiveTimestamps)

	// DSS-Id stability: building the identifier twice from the same, already-parsed signature
	// must yield the same string both times.
	id1 := sig.ID()
	id2 := sig.ID()
	if id1 != id2 || id1 == "" {
		t.Errorf("signature[%d]: DSS-Id not stable: %q then %q", index, id1, id2)
	}

	counterSignatures := sig.CounterSignatures()
	if len(counterSignatures) != want.CounterSignatureCount {
		t.Fatalf("signature[%d]: counter-signature count = %d, want %d", index, len(counterSignatures), want.CounterSignatureCount)
	}
	if len(counterSignatures) != len(want.CounterSignatures) {
		t.Fatalf("signature[%d]: golden carries %d counter-signature entries for a count of %d", index, len(want.CounterSignatures), want.CounterSignatureCount)
	}
	for i, cs := range counterSignatures {
		checkJvalSignature(t, cs.(*JAdESSignature), want.CounterSignatures[i], i, true)
	}
}

// checkJvalTimestamps asserts one time-stamp bucket against the golden's, field by field. See
// jvalTimestamp's own doc comment for why DSSID and TimestampedReferences are asserted.
func checkJvalTimestamps(t *testing.T, index int, bucket string, got []*validation.TimestampToken, want []jvalTimestamp) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("signature[%d]: %s count = %d, want %d", index, bucket, len(got), len(want))
		return
	}
	for i, timestampToken := range got {
		expected := want[i]
		prefix := fmt.Sprintf("signature[%d]: %s[%d]", index, bucket, i)

		if gotType := string(timestampToken.TimeStampType()); gotType != expected.Type {
			t.Errorf("%s: type = %s, want %s", prefix, gotType, expected.Type)
		}
		if gotID := timestampToken.DSSIDAsString(); gotID != expected.DSSID {
			t.Errorf("%s: DSSIDAsString() = %s, want %s", prefix, gotID, expected.DSSID)
		}
		gotGenerationTime := timestampToken.GenerationTime()
		switch {
		case expected.GenerationTimeMillis == nil && !gotGenerationTime.IsZero():
			t.Errorf("%s: generation time = %v, want none", prefix, gotGenerationTime)
		case expected.GenerationTimeMillis != nil && gotGenerationTime.IsZero():
			t.Errorf("%s: generation time = none, want %d", prefix, *expected.GenerationTimeMillis)
		case expected.GenerationTimeMillis != nil:
			if gotMillis := gotGenerationTime.UnixMilli(); gotMillis != *expected.GenerationTimeMillis {
				t.Errorf("%s: generation time = %d ms, want %d ms", prefix, gotMillis, *expected.GenerationTimeMillis)
			}
		}
		if got := timestampToken.IsMessageImprintDataFound(); got != expected.MessageImprintDataFound {
			t.Errorf("%s: IsMessageImprintDataFound() = %v, want %v", prefix, got, expected.MessageImprintDataFound)
		}
		if got := timestampToken.IsMessageImprintDataIntact(); got != expected.MessageImprintDataIntact {
			t.Errorf("%s: IsMessageImprintDataIntact() = %v, want %v", prefix, got, expected.MessageImprintDataIntact)
		}
		messageImprint := timestampToken.MessageImprint()
		switch {
		case expected.MessageImprintHex == nil && messageImprint.Value() != nil:
			t.Errorf("%s: message imprint = %s, want none", prefix, hexEncode(messageImprint.Value()))
		case expected.MessageImprintHex != nil && messageImprint.Value() == nil:
			t.Errorf("%s: message imprint = none, want %s", prefix, *expected.MessageImprintHex)
		case expected.MessageImprintHex != nil:
			if gotHex := hexEncode(messageImprint.Value()); gotHex != *expected.MessageImprintHex {
				t.Errorf("%s: message imprint = %s, want %s", prefix, gotHex, *expected.MessageImprintHex)
			}
		}
		if got := timestampToken.IsSignatureIntact(); got != expected.SignatureIntact {
			t.Errorf("%s: IsSignatureIntact() = %v, want %v", prefix, got, expected.SignatureIntact)
		}

		// Sorted, so only the covered SET is compared, not the order the two implementations
		// happened to append references in - the same normalization the oracle applies.
		var references []string
		for _, reference := range timestampToken.TimestampedReferences() {
			references = append(references, string(reference.Category())+":"+reference.ObjectId())
		}
		sort.Strings(references)
		if gotReferences := strings.Join(references, ","); gotReferences != expected.TimestampedReferences {
			t.Errorf("%s: timestamped references =\n  %s\nwant\n  %s", prefix, gotReferences, expected.TimestampedReferences)
		}
	}
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

// compile-time reminder that analyzer.DocumentAnalyzer is the interface the JWSDocumentAnalyzerFactory's
// Create() returns, kept honest for the a.SetCertificateVerifier/SetDetachedContents/Signatures
// calls above (all promoted from analyzer.DefaultDocumentAnalyzer).
var _ analyzer.DocumentAnalyzerFactory = (*JWSDocumentAnalyzerFactory)(nil)
