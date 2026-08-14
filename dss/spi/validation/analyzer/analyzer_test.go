// Tests for the analyzer package's dispatch logic: the DocumentAnalyzerFactory /
// EvidenceRecordAnalyzerFactory registries (Go's stand-in for Java's ServiceLoader) and
// DefaultDocumentAnalyzer's virtual-dispatch plumbing (InitDefaultDocumentAnalyzer /
// DefaultDocumentAnalyzerOverrides), matched against
// dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/{DocumentAnalyzerFactory,
// DefaultDocumentAnalyzer}.java and the evidencerecord/ counterparts at /home/user/dss-upstream.
package analyzer

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/signature"
	"github.com/utain/esig/dss/spi/validation"
)

// ---- DocumentAnalyzerFactory registry dispatch -----------------------------

type fakeDocumentAnalyzerFactory struct {
	supports func(model.DSSDocument) bool
	result   DocumentAnalyzer
}

func (f *fakeDocumentAnalyzerFactory) IsSupported(document model.DSSDocument) bool {
	return f.supports(document)
}

func (f *fakeDocumentAnalyzerFactory) Create(document model.DSSDocument) DocumentAnalyzer {
	return f.result
}

// fakeAnalyzer is a minimal concrete DocumentAnalyzer used as the marker returned by
// fakeDocumentAnalyzerFactory.Create, and as the DefaultDocumentAnalyzer test double below.
type fakeAnalyzer struct {
	DefaultDocumentAnalyzer
	name string

	isSupportedFunc              func(model.DSSDocument) bool
	originalDocumentsForSigFunc  func(validation.AdvancedSignature) []model.DSSDocument
	buildSignaturesFunc          func() []validation.AdvancedSignature
	buildDetachedTimestampsFunc  func() []*validation.TimestampToken
	buildDetachedEvidenceRecords func() []validation.EvidenceRecord
}

func newFakeAnalyzer(name string) *fakeAnalyzer {
	a := &fakeAnalyzer{DefaultDocumentAnalyzer: NewDefaultDocumentAnalyzerBase(), name: name}
	a.InitDefaultDocumentAnalyzer(a)
	return a
}

func (f *fakeAnalyzer) IsSupported(document model.DSSDocument) bool {
	if f.isSupportedFunc != nil {
		return f.isSupportedFunc(document)
	}
	return true
}

func (f *fakeAnalyzer) OriginalDocumentsForSignature(sig validation.AdvancedSignature) []model.DSSDocument {
	if f.originalDocumentsForSigFunc != nil {
		return f.originalDocumentsForSigFunc(sig)
	}
	return nil
}

func (f *fakeAnalyzer) BuildSignatures() []validation.AdvancedSignature {
	if f.buildSignaturesFunc != nil {
		return f.buildSignaturesFunc()
	}
	return f.DefaultDocumentAnalyzer.BuildSignatures()
}

func (f *fakeAnalyzer) BuildDetachedTimestamps() []*validation.TimestampToken {
	if f.buildDetachedTimestampsFunc != nil {
		return f.buildDetachedTimestampsFunc()
	}
	return f.DefaultDocumentAnalyzer.BuildDetachedTimestamps()
}

func (f *fakeAnalyzer) BuildDetachedEvidenceRecords() []validation.EvidenceRecord {
	if f.buildDetachedEvidenceRecords != nil {
		return f.buildDetachedEvidenceRecords()
	}
	return f.DefaultDocumentAnalyzer.BuildDetachedEvidenceRecords()
}

var _ DocumentAnalyzer = (*fakeAnalyzer)(nil)

// withRegistry saves and restores the package-level factory registries around a test, since
// they are shared mutable global state (mirroring the Java ServiceLoader-populated singleton
// these registries replace).
func withDocumentAnalyzerRegistry(t *testing.T, fn func()) {
	t.Helper()
	saved := documentAnalyzerFactoryRegistry
	documentAnalyzerFactoryRegistry = nil
	t.Cleanup(func() { documentAnalyzerFactoryRegistry = saved })
	fn()
}

func withEvidenceRecordAnalyzerRegistry(t *testing.T, fn func()) {
	t.Helper()
	saved := evidenceRecordAnalyzerFactoryRegistry
	evidenceRecordAnalyzerFactoryRegistry = nil
	t.Cleanup(func() { evidenceRecordAnalyzerFactoryRegistry = saved })
	fn()
}

func TestDocumentAnalyzerFromDocumentNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil document")
		}
	}()
	_, _ = DocumentAnalyzerFromDocument(nil)
}

func TestDocumentAnalyzerFromDocumentNoFactorySupports(t *testing.T) {
	withDocumentAnalyzerRegistry(t, func() {
		RegisterDocumentAnalyzerFactory(&fakeDocumentAnalyzerFactory{
			supports: func(model.DSSDocument) bool { return false },
		})
		_, err := DocumentAnalyzerFromDocument(model.NewInMemoryDocument([]byte("x")))
		if err == nil {
			t.Fatal("expected error when no factory supports the document")
		}
	})
}

func TestDocumentAnalyzerFromDocumentFirstMatchWins(t *testing.T) {
	withDocumentAnalyzerRegistry(t, func() {
		first := newFakeAnalyzer("first")
		second := newFakeAnalyzer("second")
		// Both factories claim support; registration order decides which wins, mirroring
		// ServiceLoader iteration order and the Java fromDocument loop's first-match return.
		RegisterDocumentAnalyzerFactory(&fakeDocumentAnalyzerFactory{
			supports: func(model.DSSDocument) bool { return true },
			result:   first,
		})
		RegisterDocumentAnalyzerFactory(&fakeDocumentAnalyzerFactory{
			supports: func(model.DSSDocument) bool { return true },
			result:   second,
		})

		got, err := DocumentAnalyzerFromDocument(model.NewInMemoryDocument([]byte("x")))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != DocumentAnalyzer(first) {
			t.Fatalf("DocumentAnalyzerFromDocument returned %v, want the first matching factory's result", got)
		}
	})
}

func TestDocumentAnalyzerFromDocumentSkipsUnsupportedFactories(t *testing.T) {
	withDocumentAnalyzerRegistry(t, func() {
		match := newFakeAnalyzer("match")
		RegisterDocumentAnalyzerFactory(&fakeDocumentAnalyzerFactory{
			supports: func(model.DSSDocument) bool { return false },
		})
		RegisterDocumentAnalyzerFactory(&fakeDocumentAnalyzerFactory{
			supports: func(model.DSSDocument) bool { return true },
			result:   match,
		})

		got, err := DocumentAnalyzerFromDocument(model.NewInMemoryDocument([]byte("x")))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != DocumentAnalyzer(match) {
			t.Fatalf("DocumentAnalyzerFromDocument returned %v, want match", got)
		}
	})
}

// ---- EvidenceRecordAnalyzerFactory registry dispatch -----------------------

type fakeEvidenceRecordAnalyzer struct {
	EvidenceRecordAnalyzer
}

func (fakeEvidenceRecordAnalyzer) EvidenceRecord() validation.EvidenceRecord { return nil }
func (fakeEvidenceRecordAnalyzer) EvidenceRecordType() enumerations.EvidenceRecordTypeEnum {
	return ""
}

type fakeEvidenceRecordAnalyzerFactory struct {
	supports func(model.DSSDocument) bool
	result   EvidenceRecordAnalyzer
}

func (f *fakeEvidenceRecordAnalyzerFactory) IsSupported(document model.DSSDocument) bool {
	return f.supports(document)
}

func (f *fakeEvidenceRecordAnalyzerFactory) Create(document model.DSSDocument) EvidenceRecordAnalyzer {
	return f.result
}

func TestEvidenceRecordAnalyzerIsSupportedDocumentNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil document")
		}
	}()
	EvidenceRecordAnalyzerIsSupportedDocument(nil)
}

func TestEvidenceRecordAnalyzerFromDocumentNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil document")
		}
	}()
	_, _ = EvidenceRecordAnalyzerFromDocument(nil)
}

func TestEvidenceRecordAnalyzerIsSupportedDocumentFalseWhenNoneMatch(t *testing.T) {
	withEvidenceRecordAnalyzerRegistry(t, func() {
		doc := model.NewInMemoryDocument([]byte("x"))
		if EvidenceRecordAnalyzerIsSupportedDocument(doc) {
			t.Fatal("expected false with no registered factories")
		}
	})
}

func TestEvidenceRecordAnalyzerFromDocumentDispatchesToMatchingFactory(t *testing.T) {
	withEvidenceRecordAnalyzerRegistry(t, func() {
		doc := model.NewInMemoryDocument([]byte("x"))

		RegisterEvidenceRecordAnalyzerFactory(&fakeEvidenceRecordAnalyzerFactory{
			supports: func(model.DSSDocument) bool { return false },
		})
		if EvidenceRecordAnalyzerIsSupportedDocument(doc) {
			t.Fatal("expected false: only registered factory rejects the document")
		}
		if _, err := EvidenceRecordAnalyzerFromDocument(doc); err == nil {
			t.Fatal("expected error: no factory supports the document")
		}

		var want EvidenceRecordAnalyzer = fakeEvidenceRecordAnalyzer{}
		RegisterEvidenceRecordAnalyzerFactory(&fakeEvidenceRecordAnalyzerFactory{
			supports: func(model.DSSDocument) bool { return true },
			result:   want,
		})

		if !EvidenceRecordAnalyzerIsSupportedDocument(doc) {
			t.Fatal("expected true once a supporting factory is registered")
		}
		got, err := EvidenceRecordAnalyzerFromDocument(doc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != want {
			t.Fatalf("EvidenceRecordAnalyzerFromDocument returned %v, want %v", got, want)
		}
	})
}

// ---- RegisterDocumentAnalyzerFactory / RegisterEvidenceRecordAnalyzerFactory append order -----

func TestRegisterDocumentAnalyzerFactoryAppendsInOrder(t *testing.T) {
	withDocumentAnalyzerRegistry(t, func() {
		if len(documentAnalyzerFactoryRegistry) != 0 {
			t.Fatalf("expected empty registry, got %d entries", len(documentAnalyzerFactoryRegistry))
		}
		f1 := &fakeDocumentAnalyzerFactory{supports: func(model.DSSDocument) bool { return false }}
		f2 := &fakeDocumentAnalyzerFactory{supports: func(model.DSSDocument) bool { return false }}
		RegisterDocumentAnalyzerFactory(f1)
		RegisterDocumentAnalyzerFactory(f2)
		if len(documentAnalyzerFactoryRegistry) != 2 {
			t.Fatalf("len(registry) = %d, want 2", len(documentAnalyzerFactoryRegistry))
		}
		if documentAnalyzerFactoryRegistry[0] != DocumentAnalyzerFactory(f1) || documentAnalyzerFactoryRegistry[1] != DocumentAnalyzerFactory(f2) {
			t.Fatal("registry entries not in registration order")
		}
	})
}

// ---- DefaultDocumentAnalyzer: base state and required-field panics ---------

func TestNewDefaultDocumentAnalyzerBaseDefaults(t *testing.T) {
	base := NewDefaultDocumentAnalyzerBase()
	if base.DetachedContents() == nil {
		t.Fatal("DetachedContents() should default to an empty (non-nil) slice, matching the Java constructor's ArrayList<>()")
	}
	if len(base.DetachedContents()) != 0 {
		t.Fatalf("DetachedContents() = %v, want empty", base.DetachedContents())
	}
	if base.TokenIdentifierProvider() == nil {
		t.Fatal("TokenIdentifierProvider() should default to an OriginalIdentifierProvider")
	}
}

func TestDefaultDocumentAnalyzerDocumentPanicsWhenUnset(t *testing.T) {
	a := newFakeAnalyzer("a")
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when Document() is called without SetDocument")
		}
	}()
	a.Document()
}

func TestDefaultDocumentAnalyzerSetDocumentThenDocument(t *testing.T) {
	a := newFakeAnalyzer("a")
	doc := model.NewInMemoryDocument([]byte("hello"))
	a.SetDocument(doc)
	if a.Document() != doc {
		t.Fatal("Document() did not return the document set via SetDocument")
	}
}

func TestDefaultDocumentAnalyzerSetCertificateVerifierNilPanics(t *testing.T) {
	a := newFakeAnalyzer("a")
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil certificateVerifier")
		}
	}()
	a.SetCertificateVerifier(nil)
}

func TestDefaultDocumentAnalyzerSetTokenIdentifierProviderNilPanics(t *testing.T) {
	a := newFakeAnalyzer("a")
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil tokenIdentifierProvider")
		}
	}()
	a.SetTokenIdentifierProvider(nil)
}

func TestDefaultDocumentAnalyzerValidatePanicsWithoutCertificateVerifier(t *testing.T) {
	a := newFakeAnalyzer("a")
	a.SetDocument(model.NewInMemoryDocument([]byte("x")))
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: CertificateVerifier is not defined")
		}
	}()
	a.Validate()
}

func TestDefaultDocumentAnalyzerValidatePanicsWithoutDocument(t *testing.T) {
	a := newFakeAnalyzer("a")
	a.SetCertificateVerifier(validation.NewCommonCertificateVerifier())
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: Document is not provided")
		}
	}()
	a.Validate()
}

func TestDefaultDocumentAnalyzerSignatureByIDEmptyIDPanics(t *testing.T) {
	a := newFakeAnalyzer("a")
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for empty signatureId")
		}
	}()
	a.SignatureByID("")
}

// ---- DefaultDocumentAnalyzer: caching of Signatures/DetachedTimestamps/DetachedEvidenceRecords

func TestDefaultDocumentAnalyzerSignaturesCachesBuildSignatures(t *testing.T) {
	a := newFakeAnalyzer("a")
	calls := 0
	sig := &fakeAdvancedSignature{id: "sig-1"}
	a.buildSignaturesFunc = func() []validation.AdvancedSignature {
		calls++
		return []validation.AdvancedSignature{sig}
	}

	got1 := a.Signatures()
	got2 := a.Signatures()

	if calls != 1 {
		t.Fatalf("BuildSignatures called %d times, want 1 (cached)", calls)
	}
	if len(got1) != 1 || got1[0] != validation.AdvancedSignature(sig) {
		t.Fatalf("Signatures() = %v", got1)
	}
	if len(got2) != 1 || got2[0] != validation.AdvancedSignature(sig) {
		t.Fatalf("second Signatures() call = %v", got2)
	}
}

func TestDefaultDocumentAnalyzerDetachedTimestampsCachesBuild(t *testing.T) {
	a := newFakeAnalyzer("a")
	calls := 0
	a.buildDetachedTimestampsFunc = func() []*validation.TimestampToken {
		calls++
		return []*validation.TimestampToken{{}}
	}
	a.DetachedTimestamps()
	a.DetachedTimestamps()
	if calls != 1 {
		t.Fatalf("BuildDetachedTimestamps called %d times, want 1 (cached)", calls)
	}
}

func TestDefaultDocumentAnalyzerDetachedEvidenceRecordsCachesBuild(t *testing.T) {
	a := newFakeAnalyzer("a")
	calls := 0
	a.buildDetachedEvidenceRecords = func() []validation.EvidenceRecord {
		calls++
		return nil
	}
	a.DetachedEvidenceRecords()
	a.DetachedEvidenceRecords()
	if calls != 1 {
		t.Fatalf("BuildDetachedEvidenceRecords called %d times, want 1 (cached)", calls)
	}
}

func TestDefaultDocumentAnalyzerBuildDetachedEvidenceRecordsDefaultEmpty(t *testing.T) {
	a := newFakeAnalyzer("a")
	if got := a.DetachedEvidenceRecords(); got != nil {
		t.Fatalf("DetachedEvidenceRecords() default = %v, want nil (no detached evidence record documents set)", got)
	}
}

// ---- DefaultDocumentAnalyzer: ValidationTime ------------------------------

func TestDefaultDocumentAnalyzerValidationTimeDefaultsToNow(t *testing.T) {
	a := newFakeAnalyzer("a")
	t1 := a.ValidationTime()
	t2 := a.ValidationTime()
	if !t1.Equal(t2) {
		t.Fatalf("ValidationTime() should be cached after first call: %v != %v", t1, t2)
	}
}

// ---- DefaultDocumentAnalyzer: GetValidationDataWithTimestamps empty-input error ----

func TestDefaultDocumentAnalyzerGetValidationDataWithTimestampsEmptyErrors(t *testing.T) {
	a := newFakeAnalyzer("a")
	got, err := a.GetValidationDataWithTimestamps(nil, nil)
	if err == nil {
		t.Fatal("expected error when neither signatures nor detachedTimestamps are provided")
	}
	if got != nil {
		t.Fatalf("expected nil result on error, got %v", got)
	}
}

// ---- DefaultDocumentAnalyzer: OriginalDocuments dispatch -------------------

func TestDefaultDocumentAnalyzerOriginalDocumentsUnknownIDReturnsNil(t *testing.T) {
	a := newFakeAnalyzer("a")
	if got := a.OriginalDocuments("does-not-exist"); got != nil {
		t.Fatalf("OriginalDocuments() = %v, want nil for unknown signature id", got)
	}
}

func TestDefaultDocumentAnalyzerOriginalDocumentsDispatchesToOverride(t *testing.T) {
	a := newFakeAnalyzer("a")
	sig := &fakeAdvancedSignature{id: "sig-1"}
	a.buildSignaturesFunc = func() []validation.AdvancedSignature {
		return []validation.AdvancedSignature{sig}
	}
	want := []model.DSSDocument{model.NewInMemoryDocument([]byte("orig"))}
	a.originalDocumentsForSigFunc = func(s validation.AdvancedSignature) []model.DSSDocument {
		if s != validation.AdvancedSignature(sig) {
			t.Fatalf("OriginalDocumentsForSignature called with %v, want sig", s)
		}
		return want
	}

	got := a.OriginalDocuments("sig-1")
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("OriginalDocuments() = %v, want %v", got, want)
	}
}

// ---- DefaultDocumentAnalyzer: GetAllSignatures + counter signatures --------

func TestDefaultDocumentAnalyzerGetAllSignaturesIncludesCounterSignatures(t *testing.T) {
	a := newFakeAnalyzer("a")
	a.SetCertificateVerifier(validation.NewCommonCertificateVerifier())

	counter := &fakeAdvancedSignature{id: "counter"}
	master := &fakeAdvancedSignature{id: "master", counterSignatures: []validation.AdvancedSignature{counter}}
	a.buildSignaturesFunc = func() []validation.AdvancedSignature {
		return []validation.AdvancedSignature{master}
	}

	all := a.GetAllSignatures()
	if len(all) != 2 {
		t.Fatalf("GetAllSignatures() returned %d signatures, want 2 (master + counter)", len(all))
	}
	if all[0].ID() != "master" || all[1].ID() != "counter" {
		t.Fatalf("GetAllSignatures() = [%s, %s], want [master, counter]", all[0].ID(), all[1].ID())
	}
}

// ---- fakeAdvancedSignature: a partial AdvancedSignature double, embedding the (nil)
// interface so every one of its methods type-checks (same technique as
// signatureStatusDeterminismFakeSignature in package validation), overriding only what the
// exercised analyzer code paths call. ----

type fakeAdvancedSignature struct {
	validation.AdvancedSignature
	id                string
	daIdentifier      string
	counterSignatures []validation.AdvancedSignature
}

func (f *fakeAdvancedSignature) ID() string           { return f.id }
func (f *fakeAdvancedSignature) DAIdentifier() string { return f.daIdentifier }
func (f *fakeAdvancedSignature) CounterSignatures() []validation.AdvancedSignature {
	return f.counterSignatures
}
func (f *fakeAdvancedSignature) AllTimestamps() []*validation.TimestampToken { return nil }
func (f *fakeAdvancedSignature) EmbeddedEvidenceRecords() []validation.EvidenceRecord {
	return nil
}
func (f *fakeAdvancedSignature) InitBaselineRequirementsChecker(validation.CertificateVerifier) {}
func (f *fakeAdvancedSignature) SignaturePolicy() *signature.SignaturePolicy                    { return nil }

var _ validation.AdvancedSignature = (*fakeAdvancedSignature)(nil)
