// Tests for the package-level helper functions AbstractTimestampSource's Java counterpart
// exposes as protected instance methods, matching
// dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/timestamp/AbstractTimestampSource.java
// at /home/user/dss-upstream. See that file's own header comment for why these landed as
// package-level functions rather than methods.
package timestamp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
)

// timestampFixture reads a fixture from internal/cmscore/testdata, reused here rather than
// duplicated per PORTING.md's "port test vectors, not JUnit code" guidance - this package's own
// testdata directory does not (yet) carry timestamp binaries.
func timestampFixture(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "..", "internal", "cmscore", "testdata", name))
	if err != nil {
		t.Fatalf("unable to read fixture %s: %v", name, err)
	}
	return content
}

func loadFixtureTimestampToken(t *testing.T, timestampType enumerations.TimestampType) *validation.TimestampToken {
	t.Helper()
	token, err := validation.NewTimestampToken(timestampFixture(t, "timestamp-token.tst"), timestampType)
	if err != nil {
		t.Fatalf("NewTimestampToken() failed: %v", err)
	}
	return token
}

// ---- addReference / addReferences: dedup by TimestampedReference.Equals ----------------------

func TestAddReferenceDedupsByEquals(t *testing.T) {
	var list []*validation.TimestampedReference
	ref1 := validation.NewTimestampedReference("id-1", enumerations.TimestampedObjectType_CERTIFICATE)
	ref1Dup := validation.NewTimestampedReference("id-1", enumerations.TimestampedObjectType_CERTIFICATE)
	ref2 := validation.NewTimestampedReference("id-2", enumerations.TimestampedObjectType_CERTIFICATE)

	addReference(&list, ref1)
	addReference(&list, ref1Dup) // equal to ref1 (same id + category): must not be appended again
	addReference(&list, ref2)

	if len(list) != 2 {
		t.Fatalf("len(list) = %d, want 2 (ref1Dup deduped against ref1)", len(list))
	}
	if list[0] != ref1 || list[1] != ref2 {
		t.Fatalf("list = %v, want [ref1, ref2] preserving insertion order", list)
	}
}

func TestAddReferenceForIdentifier(t *testing.T) {
	var list []*validation.TimestampedReference
	id := fakeXMLIdentifiable{id: "cert-id"}
	addReferenceForIdentifier(&list, id, enumerations.TimestampedObjectType_CERTIFICATE)
	if len(list) != 1 {
		t.Fatalf("len(list) = %d, want 1", len(list))
	}
	if list[0].ObjectId() != "cert-id" || list[0].Category() != enumerations.TimestampedObjectType_CERTIFICATE {
		t.Fatalf("list[0] = %+v, want {cert-id CERTIFICATE}", list[0])
	}
}

type fakeXMLIdentifiable struct{ id string }

func (f fakeXMLIdentifiable) AsXmlID() string { return f.id }

// ---- mergeReferences / containsEqualReference / timestampAddReferences ------------------------

func TestMergeReferencesDoesNotMutateBase(t *testing.T) {
	base := []*validation.TimestampedReference{
		validation.NewTimestampedReference("a", enumerations.TimestampedObjectType_CERTIFICATE),
	}
	additional := []*validation.TimestampedReference{
		validation.NewTimestampedReference("a", enumerations.TimestampedObjectType_CERTIFICATE), // dup
		validation.NewTimestampedReference("b", enumerations.TimestampedObjectType_REVOCATION),
	}

	merged := mergeReferences(base, additional)

	if len(base) != 1 {
		t.Fatalf("mergeReferences must not mutate its base argument in place, len(base) = %d", len(base))
	}
	if len(merged) != 2 {
		t.Fatalf("len(merged) = %d, want 2 (dedup the shared reference)", len(merged))
	}
}

func TestContainsEqualReference(t *testing.T) {
	refs := []*validation.TimestampedReference{
		validation.NewTimestampedReference("a", enumerations.TimestampedObjectType_CERTIFICATE),
	}
	same := validation.NewTimestampedReference("a", enumerations.TimestampedObjectType_CERTIFICATE)
	different := validation.NewTimestampedReference("b", enumerations.TimestampedObjectType_CERTIFICATE)

	if !containsEqualReference(refs, same) {
		t.Fatal("containsEqualReference() = false, want true for an equal reference")
	}
	if containsEqualReference(refs, different) {
		t.Fatal("containsEqualReference() = true, want false for a differing reference")
	}
}

func TestTimestampAddReferences(t *testing.T) {
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	ref := validation.NewTimestampedReference("extra", enumerations.TimestampedObjectType_SIGNED_DATA)

	before := len(token.TimestampedReferences())
	timestampAddReferences(token, []*validation.TimestampedReference{ref})
	after := token.TimestampedReferences()

	if len(after) != before+1 {
		t.Fatalf("len(TimestampedReferences()) = %d, want %d", len(after), before+1)
	}
	found := false
	for _, r := range after {
		if r.Equals(ref) {
			found = true
		}
	}
	if !found {
		t.Fatal("the added reference is missing from TimestampedReferences() after timestampAddReferences")
	}

	// Idempotent: adding the same reference again must not duplicate it.
	timestampAddReferences(token, []*validation.TimestampedReference{ref})
	if len(token.TimestampedReferences()) != before+1 {
		t.Fatalf("timestampAddReferences duplicated an already-present reference: len = %d, want %d",
			len(token.TimestampedReferences()), before+1)
	}
}

// ---- ensureOnlyDataTimestampReferencesPresent --------------------------------------------------

func TestEnsureOnlyDataTimestampReferencesPresentDropsUncoveredSignedData(t *testing.T) {
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)

	covered := validation.NewTimestampedReference("covered-data", enumerations.TimestampedObjectType_SIGNED_DATA)
	uncovered := validation.NewTimestampedReference("uncovered-data", enumerations.TimestampedObjectType_SIGNED_DATA)
	nonData := validation.NewTimestampedReference("cert-ref", enumerations.TimestampedObjectType_CERTIFICATE)
	token.SetTimestampedReferences([]*validation.TimestampedReference{covered, uncovered, nonData})

	// referencesToCheck only vouches for "covered"; SIGNED_DATA entries absent from it must be
	// dropped, non-SIGNED_DATA entries (nonData) must survive regardless.
	ensureOnlyDataTimestampReferencesPresent(token, []*validation.TimestampedReference{covered})

	got := token.TimestampedReferences()
	if len(got) != 2 {
		t.Fatalf("len(TimestampedReferences()) = %d, want 2 (covered + nonData survive)", len(got))
	}
	for _, r := range got {
		if r.Equals(uncovered) {
			t.Fatal("uncovered SIGNED_DATA reference should have been removed")
		}
	}
}

// ---- CreateReferenceForCertificate / CreateReferencesForCertificates --------------------------

func generateTimestampSourceTestCertificate(t *testing.T) *model.CertificateToken {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %s", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(7),
		Subject:      pkix.Name{CommonName: "timestamp source test cert"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %s", err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse certificate: %s", err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("NewCertificateToken: %s", err)
	}
	return token
}

func TestCreateReferenceForCertificate(t *testing.T) {
	cert := generateTimestampSourceTestCertificate(t)
	ref := CreateReferenceForCertificate(cert)
	if ref.Category() != enumerations.TimestampedObjectType_CERTIFICATE {
		t.Fatalf("Category() = %v, want CERTIFICATE", ref.Category())
	}
	if ref.ObjectId() != cert.DSSID().AsXmlID() {
		t.Fatalf("ObjectId() = %q, want the certificate's DSS id %q", ref.ObjectId(), cert.DSSID().AsXmlID())
	}
}

func TestCreateReferencesForCertificatesDedups(t *testing.T) {
	cert1 := generateTimestampSourceTestCertificate(t)
	cert2 := generateTimestampSourceTestCertificate(t)
	refs := CreateReferencesForCertificates([]*model.CertificateToken{cert1, cert2, cert1})
	if len(refs) != 2 {
		t.Fatalf("len(refs) = %d, want 2 (cert1 deduped)", len(refs))
	}
}

// ---- SignerDataTimestampedReferences: recurses through Children() -----------------------------

// testSignatureScope is a minimal concrete SignatureScope, since no format-specific
// implementation (XAdES/CAdES manifest entries, ...) has landed yet to exercise this against.
type testSignatureScope struct {
	scope.SignatureScopeBase
}

func newTestSignatureScope(doc model.DSSDocument) *testSignatureScope {
	return &testSignatureScope{SignatureScopeBase: scope.NewSignatureScopeBase(doc)}
}

func (s *testSignatureScope) Description(model.TokenIdentifierProvider) string { return "" }
func (s *testSignatureScope) Type() enumerations.SignatureScopeType            { return "" }

var _ scope.SignatureScope = (*testSignatureScope)(nil)

func TestSignerDataTimestampedReferencesRecursesThroughChildren(t *testing.T) {
	child := newTestSignatureScope(model.NewInMemoryDocumentWithName([]byte("child"), "child.txt"))
	parent := newTestSignatureScope(model.NewInMemoryDocumentWithName([]byte("parent"), "parent.txt"))
	parent.AddChildSignatureScope(child)

	refs := SignerDataTimestampedReferences([]scope.SignatureScope{parent})

	if len(refs) != 2 {
		t.Fatalf("len(refs) = %d, want 2 (parent + child)", len(refs))
	}
	for _, r := range refs {
		if r.Category() != enumerations.TimestampedObjectType_SIGNED_DATA {
			t.Fatalf("Category() = %v, want SIGNED_DATA for every reference", r.Category())
		}
	}
	wantParentID := parent.DSSIDAsString()
	wantChildID := child.DSSIDAsString()
	if refs[0].ObjectId() != wantParentID {
		t.Fatalf("refs[0].ObjectId() = %q, want the parent's DSS id %q (parent added before recursing into children)", refs[0].ObjectId(), wantParentID)
	}
	if refs[1].ObjectId() != wantChildID {
		t.Fatalf("refs[1].ObjectId() = %q, want the child's DSS id %q", refs[1].ObjectId(), wantChildID)
	}
}

// ---- ReferencesFromTimestamp: a real fixture TimestampToken, empty merged sources -------------

func TestReferencesFromTimestampIncludesTheTokenItself(t *testing.T) {
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)

	refs, err := ReferencesFromTimestamp(token,
		spi.NewListCertificateSource(),
		spi.NewListRevocationSource[revocation.CRL](),
		spi.NewListRevocationSource[revocation.OCSP]())
	if err != nil {
		t.Fatalf("ReferencesFromTimestamp: %v", err)
	}

	found := false
	for _, r := range refs {
		if r.Category() == enumerations.TimestampedObjectType_TIMESTAMP && r.ObjectId() == token.DSSIDAsString() {
			found = true
		}
	}
	if !found {
		t.Fatal("ReferencesFromTimestamp() did not include a TIMESTAMP reference to the token itself")
	}
}
