package validation

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
)

// The known answers below come from testdata/bc-oracle.txt, which testdata/gen/TspFixtures.java
// produces by reading the same fixtures with BouncyCastle 1.84 on OpenJDK 21. The fixtures are
// two RFC 3161 tokens issued with exactly the wiring KeyEntityTSPSource uses:
//
//	timestamp-token.tst          SHA-256 message imprint, SHA-512 signature, TSA + CA embedded
//	timestamp-token-sha512.tst   SHA-512 message imprint, SHA-256 signature, TSA only

// timestampTokenKATOracle reads testdata/bc-oracle.txt into a map of its key=value lines.
func timestampTokenKATOracle(t *testing.T) map[string]string {
	t.Helper()
	file, err := os.Open(corpustest.Path(t, "bc-oracle.txt"))
	if err != nil {
		t.Fatalf("unable to read the BouncyCastle oracle: %v", err)
	}
	defer file.Close()

	oracle := map[string]string{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<20)
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), "=")
		if found {
			oracle[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("unable to read the BouncyCastle oracle: %v", err)
	}
	return oracle
}

// timestampTokenKATFile reads a fixture of testdata.
func timestampTokenKATFile(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("unable to read %s: %v", name, err)
	}
	return content
}

// timestampTokenKATCertificate loads one of the DER certificates of testdata.
func timestampTokenKATCertificate(t *testing.T, name string) *model.CertificateToken {
	t.Helper()
	certificate, err := x509.ParseCertificate(timestampTokenKATFile(t, name))
	if err != nil {
		t.Fatalf("unable to parse %s: %v", name, err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("unable to build a CertificateToken for %s: %v", name, err)
	}
	return token
}

// TestTimestampTokenKAT_Fields walks the TSTInfo surface of the SHA-256 fixture and compares
// every value with the one BouncyCastle reads from the same bytes.
func TestTimestampTokenKAT_Fields(t *testing.T) {
	oracle := timestampTokenKATOracle(t)
	binaries := timestampTokenKATFile(t, "timestamp-token.tst")

	token, err := NewTimestampToken(binaries, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}

	if got := token.TimeStampType(); got != enumerations.TimestampType_SIGNATURE_TIMESTAMP {
		t.Errorf("TimeStampType() = %s, want SIGNATURE_TIMESTAMP", got)
	}
	if got, want := token.GenerationTime().UnixMilli(), oracle["timestamp-token.tst.genTime"]; strconv.FormatInt(got, 10) != want {
		t.Errorf("GenerationTime() = %d, want %s", got, want)
	}
	if !token.CreationDate().Equal(token.GenerationTime()) {
		t.Errorf("CreationDate() = %s, want the generation time %s", token.CreationDate(), token.GenerationTime())
	}
	if got := token.DigestAlgorithm(); got != enumerations.DigestAlgorithm_SHA256 {
		t.Errorf("DigestAlgorithm() = %s, want SHA256", got)
	}
	messageImprint := token.MessageImprint()
	if got, want := hex.EncodeToString(messageImprint.Value()),
		oracle["timestamp-token.tst.messageImprintDigest"]; got != want {
		t.Errorf("MessageImprint().Value() = %s, want %s", got, want)
	}
	if got := messageImprint.Algorithm(); got != enumerations.DigestAlgorithm_SHA256 {
		t.Errorf("MessageImprint().Algorithm() = %s, want SHA256", got)
	}
	if got := token.TSTInfoTsa(); got != nil {
		t.Errorf("TSTInfoTsa() = %s, want nil (the fixture carries no tsa field)", got)
	}
	if got := len(token.Certificates()); got != 2 {
		t.Errorf("len(Certificates()) = %d, want 2", got)
	}
	if got := token.UnsignedAttributes(); len(got) != 0 {
		t.Errorf("UnsignedAttributes() = %v, want none", got)
	}
	if got := token.SignerInformation(); got == nil {
		t.Fatal("SignerInformation() = nil, want the TSA signer")
	}
	if got, want := token.SignerInformation().DigestAlgorithm.Algorithm.String(),
		oracle["timestamp-token.tst.signerDigestAlgOID"]; got != want {
		t.Errorf("SignerInformation().DigestAlgorithm = %s, want %s", got, want)
	}
	if got := len(token.SignerInformationStoreInfos()); got != 1 {
		t.Errorf("len(SignerInformationStoreInfos()) = %d, want 1", got)
	}
	if got := token.CandidatesForSigningCertificate(); got == nil || got.IsEmpty() {
		t.Error("CandidatesForSigningCertificate() is empty, want the embedded TSA certificate")
	}

	// getEncoded() DER-encodes the token; the fixture is already DER, so it round-trips.
	encoded := token.Encoded()
	if got, want := hex.EncodeToString(sha256Sum(encoded)), oracle["timestamp-token.tst.sha256"]; got != want {
		t.Errorf("Encoded() digest = %s, want %s", got, want)
	}
	digest, err := token.Digest(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatalf("Digest() failed: %v", err)
	}
	if got, want := hex.EncodeToString(digest), oracle["timestamp-token.tst.sha256"]; got != want {
		t.Errorf("Digest(SHA256) = %s, want %s", got, want)
	}

	// The identifier is the "T-" prefixed digest of the DER encoding, which makes it stable
	// across re-parsings of the same bytes.
	twin, err := NewTimestampToken(binaries, enumerations.TimestampType_ARCHIVE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	if token.DSSIDAsString() != twin.DSSIDAsString() {
		t.Errorf("DSSIDAsString() = %s and %s for the same binaries", token.DSSIDAsString(), twin.DSSIDAsString())
	}
	if !strings.HasPrefix(token.DSSIDAsString(), "T-") {
		t.Errorf("DSSIDAsString() = %s, want the T- prefix", token.DSSIDAsString())
	}
}

// TestTimestampTokenKAT_SignedBy validates the fixture against its TSA certificate, which
// exercises the whole checkSignedAndValid chain: SID matching, the ESSCertIDv2 certificate hash,
// the timeStamping extended key usage, the validity of the certificate at genTime and the CMS
// signature over the signed attributes.
func TestTimestampTokenKAT_SignedBy(t *testing.T) {
	token, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token.tst"),
		enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	tsa := timestampTokenKATCertificate(t, "tsa.crt")

	if !token.IsSignedByToken(tsa) {
		t.Fatalf("IsSignedByToken(tsa) = false, want true (reason: %s)", token.InvalidityReason())
	}
	if !token.IsSignatureIntact() {
		t.Error("IsSignatureIntact() = false, want true")
	}
	if got := token.SignatureValidity(); got != enumerations.SignatureValidity_VALID {
		t.Errorf("SignatureValidity() = %s, want VALID", got)
	}
	if got := token.SignatureAlgorithm(); got != enumerations.SignatureAlgorithm_RSA_SHA512 {
		t.Errorf("SignatureAlgorithm() = %s, want RSA_SHA512", got)
	}
	if got := token.IssuerX500Principal(); got == nil || !strings.Contains(got.String(), "Test TSA") {
		t.Errorf("IssuerX500Principal() = %v, want the TSA subject", got)
	}
	// A second call is answered from the remembered public key.
	if !token.IsSignedByToken(tsa) {
		t.Error("IsSignedByToken(tsa) = false on the second call, want true")
	}

	// The CA certificate is not the signer: its issuer and serial do not match the SID.
	fresh, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token.tst"),
		enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	if fresh.IsSignedByToken(timestampTokenKATCertificate(t, "ca.crt")) {
		t.Error("IsSignedByToken(ca) = true, want false")
	}
	if got := fresh.SignatureValidity(); got != enumerations.SignatureValidity_NOT_EVALUATED {
		t.Errorf("SignatureValidity() = %s, want NOT_EVALUATED for a certificate the SID rejects", got)
	}

	// plain.crt carries the TSA's public key under another serial number and without the
	// timeStamping extended key usage, so its key would verify the signature: it is the
	// SignerIdentifier gate, not the cryptography, that has to turn it down.
	plain := timestampTokenKATCertificate(t, "plain.crt")
	if !plain.PublicKey().Equals(tsa.PublicKey()) {
		t.Fatal("plain.crt does not share the TSA public key; the fixture lost its purpose")
	}
	again, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token.tst"),
		enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	if again.IsSignedByToken(plain) {
		t.Error("IsSignedByToken(plain) = true, want false: the SID names another certificate")
	}
}

// TestTimestampTokenKAT_MatchData exercises the three matchData families against the content the
// fixture was issued over.
func TestTimestampTokenKAT_MatchData(t *testing.T) {
	content := timestampTokenKATFile(t, "content.bin")
	token, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token.tst"),
		enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}

	if token.IsProcessed() {
		t.Error("IsProcessed() = true before any matchData call, want false")
	}
	expected := sha256Sum(content)
	if !token.MatchData(expected) {
		t.Error("MatchData(sha256(content)) = false, want true")
	}
	if !token.IsProcessed() {
		t.Error("IsProcessed() = false after matchData, want true")
	}
	if !token.IsMessageImprintDataFound() || !token.IsMessageImprintDataIntact() {
		t.Error("the message imprint was not reported as found and intact")
	}

	// A DSSDocument over the same content digests to the same value.
	document := model.NewInMemoryDocument(content)
	matched, err := token.MatchDataDocument(document)
	if err != nil {
		t.Fatalf("MatchDataDocument() failed: %v", err)
	}
	if !matched {
		t.Error("MatchDataDocument(content) = false, want true")
	}

	// So does a DSSMessageDigest carrying the same algorithm and value.
	if !token.MatchDataMessageDigest(model.NewDSSMessageDigestWithValue(
		enumerations.DigestAlgorithm_SHA256, expected)) {
		t.Error("MatchDataMessageDigest(SHA256) = false, want true")
	}
	// A message digest computed with another algorithm never matches.
	if token.MatchDataMessageDigest(model.NewDSSMessageDigestWithValue(
		enumerations.DigestAlgorithm_SHA512, expected)) {
		t.Error("MatchDataMessageDigest(SHA512) = true, want false: the imprint is a SHA-256 one")
	}
	// Neither does an empty one.
	if token.MatchDataMessageDigest(model.CreateEmptyDSSMessageDigest()) {
		t.Error("MatchDataMessageDigest(empty) = true, want false")
	}

	// Absent data is reported as "not found", not as "not intact".
	if token.MatchData(nil) {
		t.Error("MatchData(nil) = true, want false")
	}
	if token.IsMessageImprintDataFound() {
		t.Error("IsMessageImprintDataFound() = true after MatchData(nil), want false")
	}
	// Wrong data is found but not intact.
	if token.MatchData(sha256Sum([]byte("something else"))) {
		t.Error("MatchData(other) = true, want false")
	}
	if !token.IsMessageImprintDataFound() {
		t.Error("IsMessageImprintDataFound() = false for a provided digest, want true")
	}
	if token.IsMessageImprintDataIntact() {
		t.Error("IsMessageImprintDataIntact() = true for a wrong digest, want false")
	}
}

// TestTimestampTokenKAT_Sha512Fixture checks the second fixture, whose imprint algorithm and
// signature algorithm are swapped and whose chain holds the TSA certificate only.
func TestTimestampTokenKAT_Sha512Fixture(t *testing.T) {
	oracle := timestampTokenKATOracle(t)
	token, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token-sha512.tst"),
		enumerations.TimestampType_ARCHIVE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}

	if got := token.DigestAlgorithm(); got != enumerations.DigestAlgorithm_SHA512 {
		t.Errorf("DigestAlgorithm() = %s, want SHA512", got)
	}
	if got, want := hex.EncodeToString(token.MessageImprint().Value()),
		oracle["timestamp-token-sha512.tst.messageImprintDigest"]; got != want {
		t.Errorf("MessageImprint().Value() = %s, want %s", got, want)
	}
	if got := len(token.Certificates()); got != 1 {
		t.Errorf("len(Certificates()) = %d, want 1", got)
	}
	if !token.IsSignedByToken(timestampTokenKATCertificate(t, "tsa.crt")) {
		t.Fatalf("IsSignedByToken(tsa) = false, want true (reason: %s)", token.InvalidityReason())
	}
	if got := token.SignatureAlgorithm(); got != enumerations.SignatureAlgorithm_RSA_SHA256 {
		t.Errorf("SignatureAlgorithm() = %s, want RSA_SHA256", got)
	}
	if !token.MatchData(sha512Sum(timestampTokenKATFile(t, "content.bin"))) {
		t.Error("MatchData(sha512(content)) = false, want true")
	}
	if !token.IsValid() {
		t.Error("IsValid() = false, want true")
	}
}

// TestTimestampTokenKAT_TamperedContentInfo checks that a token whose eContent was altered no
// longer verifies: the message-digest signed attribute stops matching the TSTInfo, which is the
// CMSSignerDigestMismatchException branch of SignerInformation#verify.
func TestTimestampTokenKAT_TamperedContentInfo(t *testing.T) {
	binaries := timestampTokenKATFile(t, "timestamp-token.tst")
	// The message imprint digest is embedded in the TSTInfo; flipping one of its bytes leaves
	// the structure intact and only breaks the message-digest attribute.
	imprint, err := hex.DecodeString(timestampTokenKATOracle(t)["timestamp-token.tst.messageImprintDigest"])
	if err != nil {
		t.Fatalf("unable to read the oracle imprint: %v", err)
	}
	index := indexOf(binaries, imprint)
	if index < 0 {
		t.Fatal("the message imprint was not found in the fixture")
	}
	tampered := append([]byte{}, binaries...)
	tampered[index] ^= 0x01

	token, err := NewTimestampToken(tampered, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	if token.IsSignedByToken(timestampTokenKATCertificate(t, "tsa.crt")) {
		t.Error("IsSignedByToken(tsa) = true for a tampered token, want false")
	}
	if got := token.SignatureValidity(); got != enumerations.SignatureValidity_INVALID {
		t.Errorf("SignatureValidity() = %s, want INVALID", got)
	}
	if got := token.InvalidityReason(); !strings.Contains(got, "CMSSignerDigestMismatchException") {
		t.Errorf("InvalidityReason() = %q, want the message-digest mismatch", got)
	}
}

// TestTimestampTokenKAT_ToString renders a validated token, whose text has to name the TSA, the
// generation time and the two verdicts.
func TestTimestampTokenKAT_ToString(t *testing.T) {
	token, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token.tst"),
		enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	if got := token.String(); !strings.Contains(got, "TimestampToken[signedBy=null") {
		t.Errorf("String() = %q, want an unvalidated token", got)
	}
	if !token.IsSignedByToken(timestampTokenKATCertificate(t, "tsa.crt")) {
		t.Fatalf("IsSignedByToken(tsa) = false, want true (reason: %s)", token.InvalidityReason())
	}
	token.MatchData(sha256Sum(timestampTokenKATFile(t, "content.bin")))

	rendered := token.String()
	for _, expected := range []string{
		"Test TSA", "SIGNATURE_TIMESTAMP",
		"Timestamp's signature validity: VALID", "Timestamp MATCHES the signed data.",
	} {
		if !strings.Contains(rendered, expected) {
			t.Errorf("String() = %q, want it to contain %q", rendered, expected)
		}
	}
	if got := token.Abbreviation(); !strings.HasPrefix(got, "SIGNATURE_TIMESTAMP: T-") {
		t.Errorf("Abbreviation() = %q, want the type and the identifier", got)
	}
}

// TestTimestampTokenKAT_Rejected covers the two structural checks BouncyCastle's TimeStampToken
// constructor performs and this port performs in its place.
func TestTimestampTokenKAT_Rejected(t *testing.T) {
	if _, err := NewTimestampToken([]byte{0x30, 0x00}, enumerations.TimestampType_SIGNATURE_TIMESTAMP); err == nil {
		t.Error("NewTimestampToken() accepted a document that is not a CMS")
	}
	// A CMS whose eContentType is not id-ct-TSTInfo is not a time-stamp.
	binaries := timestampTokenKATFile(t, "timestamp-token.tst")
	notATimestamp := append([]byte{}, binaries...)
	// id-smime-ct-TSTInfo is 1.2.840.113549.1.9.16.1.4; its last arc identifies the content
	// type, so changing it makes the encapsulated content something else.
	tstInfoOID, _ := hex.DecodeString("060b2a864886f70d0109100104")
	index := indexOf(notATimestamp, tstInfoOID)
	if index < 0 {
		t.Fatal("the eContentType was not found in the fixture")
	}
	notATimestamp[index+12] = 0x05
	if _, err := NewTimestampToken(notATimestamp, enumerations.TimestampType_SIGNATURE_TIMESTAMP); err == nil {
		t.Error("NewTimestampToken() accepted a CMS that does not encapsulate a TSTInfo")
	}
}

// TestTimestampTokenIsMessageImprintDataIntactPanics reproduces the IllegalStateException upstream
// raises when the message imprint is read before any matchData call.
func TestTimestampTokenIsMessageImprintDataIntactPanics(t *testing.T) {
	token, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token.tst"),
		enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	defer func() {
		recovered := recover()
		if recovered != "Invoke matchData(byte[] data) method before!" {
			t.Errorf("recover() = %v, want the IllegalStateException message", recovered)
		}
	}()
	token.IsMessageImprintDataIntact()
	t.Error("IsMessageImprintDataIntact() returned, want a panic")
}

// TestTimestampTokenIsSignedByPublicKeyPanics reproduces the UnsupportedOperationException of the
// two PublicKey-based overloads.
func TestTimestampTokenIsSignedByPublicKeyPanics(t *testing.T) {
	token, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token.tst"),
		enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	for name, call := range map[string]func(){
		"isSignedBy":      func() { token.IsSignedBy(nil) },
		"checkIsSignedBy": func() { token.CheckIsSignedBy(nil) },
	} {
		func() {
			defer func() {
				if recovered := recover(); recovered == nil {
					t.Errorf("%s() returned, want a panic", name)
				} else if !strings.Contains(recovered.(string), "for a TimestampToken validation!") {
					t.Errorf("%s() panicked with %v, want the UnsupportedOperationException message", name, recovered)
				}
			}()
			call()
		}()
	}
}

// TestTimestampTokenSetters covers the plain state the validation layers set on a token.
func TestTimestampTokenSetters(t *testing.T) {
	token, err := NewTimestampTokenWithReferences(timestampTokenKATFile(t, "timestamp-token.tst"),
		enumerations.TimestampType_ARCHIVE_TIMESTAMP,
		[]*TimestampedReference{NewTimestampedReference("id", enumerations.TimestampedObjectType_SIGNATURE)})
	if err != nil {
		t.Fatalf("NewTimestampTokenWithReferences() failed: %v", err)
	}

	if got := len(token.TimestampedReferences()); got != 1 {
		t.Errorf("len(TimestampedReferences()) = %d, want 1", got)
	}
	token.SetFilename("timestamp.tst")
	if got := token.Filename(); got != "timestamp.tst" {
		t.Errorf("Filename() = %q, want timestamp.tst", got)
	}
	token.SetArchiveTimestampType(enumerations.ArchiveTimestampType_CAdES_V3)
	if got := token.ArchiveTimestampType(); got != enumerations.ArchiveTimestampType_CAdES_V3 {
		t.Errorf("ArchiveTimestampType() = %s, want CAdES_V3", got)
	}
	token.SetCanonicalizationMethod("http://www.w3.org/2001/10/xml-exc-c14n#")
	if got := token.CanonicalizationMethod(); got != "http://www.w3.org/2001/10/xml-exc-c14n#" {
		t.Errorf("CanonicalizationMethod() = %q, want the exclusive c14n URI", got)
	}
	includes := []*TimestampInclude{NewTimestampIncludeWithURI("#r-id-1", true)}
	token.SetTimestampIncludes(includes)
	if got := token.TimestampIncludes(); len(got) != 1 || got[0].URI() != "#r-id-1" || !got[0].IsReferencedData() {
		t.Errorf("TimestampIncludes() = %v, want the single include", got)
	}
	manifest := model.NewManifestFile()
	token.SetManifestFile(manifest)
	if token.ManifestFile() != manifest {
		t.Error("ManifestFile() did not return the manifest that was set")
	}
	status := NewArchiveTimestampHashIndexStatus()
	status.SetVersion(enumerations.ArchiveTimestampHashIndexVersion_ATS_HASH_INDEX_V3)
	status.AddErrorMessage("boom")
	token.SetAtsHashIndexStatus(status)
	if got := token.AtsHashIndexStatus(); got == nil ||
		got.Version() != enumerations.ArchiveTimestampHashIndexVersion_ATS_HASH_INDEX_V3 ||
		len(got.ErrorMessages()) != 1 {
		t.Errorf("AtsHashIndexStatus() = %v, want the status that was set", got)
	}
	if got := token.DetachedEvidenceRecords(); got == nil || len(got) != 0 {
		t.Errorf("DetachedEvidenceRecords() = %v, want an empty, non-nil list", got)
	}
}

// TestTimestampTokenAreReferenceValidationsValid covers the evidence-record reference rule: an
// orphan reference never invalidates the time-stamp, any other broken one does.
func TestTimestampTokenAreReferenceValidationsValid(t *testing.T) {
	token, err := NewTimestampToken(timestampTokenKATFile(t, "timestamp-token.tst"),
		enumerations.TimestampType_ARCHIVE_TIMESTAMP)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	if !token.AreReferenceValidationsValid() {
		t.Error("AreReferenceValidationsValid() = false without any reference, want true")
	}

	intact := model.NewReferenceValidation()
	intact.SetType(enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_OBJECT)
	intact.SetFound(true)
	intact.SetIntact(true)
	token.SetReferenceValidations([]*model.ReferenceValidation{intact})
	if !token.AreReferenceValidationsValid() {
		t.Error("AreReferenceValidationsValid() = false for an intact reference, want true")
	}

	orphan := model.NewReferenceValidation()
	orphan.SetType(enumerations.DigestMatcherType_EVIDENCE_RECORD_ORPHAN_REFERENCE)
	token.SetReferenceValidations([]*model.ReferenceValidation{intact, orphan})
	if !token.AreReferenceValidationsValid() {
		t.Error("AreReferenceValidationsValid() = false for an orphan reference, want true")
	}

	broken := model.NewReferenceValidation()
	broken.SetType(enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_OBJECT)
	broken.SetFound(true)
	token.SetReferenceValidations([]*model.ReferenceValidation{intact, broken})
	if token.AreReferenceValidationsValid() {
		t.Error("AreReferenceValidationsValid() = true for a broken reference, want false")
	}
}

// indexOf returns the index of needle in haystack, or -1.
func indexOf(haystack, needle []byte) int {
	return strings.Index(string(haystack), string(needle))
}

// sha256Sum digests with SHA-256.
func sha256Sum(content []byte) []byte {
	sum := sha256.Sum256(content)
	return sum[:]
}

// sha512Sum digests with SHA-512.
func sha512Sum(content []byte) []byte {
	sum := sha512.Sum512(content)
	return sum[:]
}
