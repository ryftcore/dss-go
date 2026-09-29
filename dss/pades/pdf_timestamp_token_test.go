// Cover for the PdfTimestampToken registry (see pdf_timestamp_token.go): it stands in for Java's
// `instanceof PdfTimestampToken`, and must neither lose a wrapper whose bare
// *validation.TimestampToken is still in use nor pin every token it has ever seen.
package pades

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

func pdfTimestampTokenRegistrySize() int {
	pdfTimestampTokenRegistryMu.Lock()
	defer pdfTimestampTokenRegistryMu.Unlock()
	return len(pdfTimestampTokenRegistry)
}

// settleRegistry runs the collector until the registry stops shrinking (cleanups run on their own
// goroutine, after the cycle that found the token unreachable) or the deadline passes.
func settleRegistry(want func(size int) bool) int {
	deadline := time.Now().Add(10 * time.Second)
	for {
		runtime.GC()
		if size := pdfTimestampTokenRegistrySize(); want(size) || time.Now().After(deadline) {
			return size
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// timestampTokensOf returns the bare token of every document time-stamp revision of the fixture.
func timestampTokensOf(t *testing.T, data []byte) []*validation.TimestampToken {
	t.Helper()
	service := NewDefaultPdfObjFactory().NewPAdESSignatureService()
	var tokens []*validation.TimestampToken
	for _, revision := range service.GetRevisions(model.NewInMemoryDocument(data), nil) {
		if timestampRevision, ok := revision.(*PdfDocTimestampRevision); ok {
			tokens = append(tokens, timestampRevision.TimestampToken().TimestampToken)
		}
	}
	if len(tokens) == 0 {
		t.Fatal("the fixture yielded no document time-stamp")
	}
	return tokens
}

func TestPdfTimestampTokenRegistryDoesNotPinTokens(t *testing.T) {
	data, err := os.ReadFile(padesFixturePath(t, "upstream/validation/pades-5-signatures-and-1-document-timestamp.pdf"))
	if err != nil {
		t.Fatal(err)
	}

	// Holding only the bare token (no revision, no wrapper) must keep the wrapper resolvable
	// through any number of collections: the revision is dropped when timestampTokensOf returns.
	kept := timestampTokensOf(t, data)[0]
	for i := 0; i < 5; i++ { // let the discarded revisions' tokens (and any earlier test's) go first
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
	}
	before := pdfTimestampTokenRegistrySize()
	wrapper, ok := PdfTimestampTokenOf(kept)
	if !ok || wrapper == nil || wrapper.TimestampToken != kept || wrapper.PdfRevision() == nil {
		t.Fatalf("PdfTimestampTokenOf lost the wrapper of a live token: %v, %v", wrapper, ok)
	}
	wrapper = nil

	// Validating many documents and dropping the results must not grow the registry.
	const documents = 20
	for i := 0; i < documents; i++ {
		timestampTokensOf(t, data)
	}
	after := settleRegistry(func(size int) bool { return size <= before })
	if after > before {
		t.Errorf("registry holds %d entries after %d discarded documents, want at most the %d live before them",
			after, documents, before)
	}

	runtime.KeepAlive(kept)
	if _, ok := PdfTimestampTokenOf(kept); !ok {
		t.Error("PdfTimestampTokenOf lost the wrapper of a live token after collection")
	}
}

func TestPdfTimestampTokenOfUnknownToken(t *testing.T) {
	if wrapper, ok := PdfTimestampTokenOf(nil); ok || wrapper != nil {
		t.Errorf("PdfTimestampTokenOf(nil) = %v, %v", wrapper, ok)
	}
	if wrapper, ok := PdfTimestampTokenOf(&validation.TimestampToken{}); ok || wrapper != nil {
		t.Errorf("PdfTimestampTokenOf(unregistered) = %v, %v", wrapper, ok)
	}
}
