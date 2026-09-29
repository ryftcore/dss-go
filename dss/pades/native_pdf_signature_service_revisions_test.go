// Regression cover for AbstractPDFSignatureService#getRevisions on a hostile document.
//
// Upstream wraps the whole per-signature-dictionary body of getRevisions in
// try { ... } catch (Exception e) { LOG.warn("Unable to parse signature ..."); }: a revision
// that cannot be built is logged and skipped, and the signatures after it are still analysed.
// PdfDocTimestampRevision's constructor throws a DSSException when /Contents is not an RFC 3161
// token, which is what a signature dictionary re-labelled /SubFilter /ETSI.RFC3161 carries; the
// port let the equivalent panic escape GetRevisions and abort the analysis of every signature in
// the document.
package pades

import (
	"bytes"
	"os"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

func countRevisions(revisions []PdfRevision) (signatures, timestamps int) {
	for _, revision := range revisions {
		switch revision.(type) {
		case *PdfSignatureRevision:
			signatures++
		case *PdfDocTimestampRevision:
			timestamps++
		}
	}
	return signatures, timestamps
}

func TestGetRevisionsSkipsRevisionThatCannotBeBuilt(t *testing.T) {
	data, err := os.ReadFile(padesFixturePath(t, "upstream/validation/pades-5-signatures-and-1-document-timestamp.pdf"))
	if err != nil {
		t.Fatal(err)
	}
	service := NewDefaultPdfObjFactory().NewPAdESSignatureService()

	signatures, timestamps := countRevisions(service.GetRevisions(model.NewInMemoryDocument(data), nil))
	if signatures != 5 || timestamps != 1 {
		t.Fatalf("intact document: %d signature(s) and %d document time-stamp(s), want 5 and 1", signatures, timestamps)
	}

	// Turn the first signature dictionary into something that claims to be a document
	// time-stamp while /Contents still holds its CAdES signature. The replacements keep every
	// byte offset: /Type is blanked out (an absent /Type is accepted for a time-stamp) and the
	// /SubFilter value is padded with spaces.
	subFilter := []byte("/SubFilter /ETSI.CAdES.detached")
	relabelled := []byte("/SubFilter /ETSI.RFC3161       ")
	if len(subFilter) != len(relabelled) {
		t.Fatal("test bug: replacement changes the length")
	}
	corrupted := append([]byte(nil), data...)
	at := bytes.Index(corrupted, subFilter)
	typeAt := bytes.Index(corrupted[at:], []byte("/Type /Sig"))
	if at < 0 || typeAt < 0 || typeAt > 200 {
		t.Fatalf("unexpected fixture layout: /SubFilter at %d, /Type at +%d", at, typeAt)
	}
	copy(corrupted[at:], relabelled)
	copy(corrupted[at+typeAt:], bytes.Repeat([]byte{' '}, len("/Type /Sig")))

	var revisions []PdfRevision
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("GetRevisions panicked on a document time-stamp that is not an RFC 3161 token: %v", r)
			}
		}()
		revisions = service.GetRevisions(model.NewInMemoryDocument(corrupted), nil)
	}()
	signatures, timestamps = countRevisions(revisions)
	if signatures != 4 || timestamps != 1 {
		t.Errorf("corrupted document: %d signature(s) and %d document time-stamp(s), "+
			"want the 4 intact signatures and the genuine time-stamp kept, the malformed one skipped",
			signatures, timestamps)
	}
}
