package jades

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// openCountingDocument counts how many times its content is read.
type openCountingDocument struct {
	model.DSSDocument
	opens *int
}

func (d openCountingDocument) OpenStream() (io.ReadCloser, error) {
	*d.opens++
	return d.DSSDocument.OpenStream()
}

// TestReferenceValidationsByUriHashMechanismDigestsEachDocumentOnce checks that a sigD with many
// 'hashV' entries digests every detached document once rather than once per entry (J17A-PERF-001)
// while the verdicts stay those of the per-entry evaluation: documents are matched by digest first
// and by name as a fallback.
func TestReferenceValidationsByUriHashMechanismDigestsEachDocumentOnce(t *testing.T) {
	const entries = 12

	contents := []string{"first document", "second document", "third document"}
	opens := make([]int, len(contents))
	var detached []model.DSSDocument
	for i, content := range contents {
		document := model.NewInMemoryDocumentWithName([]byte(content), fmt.Sprintf("doc%d.txt", i))
		detached = append(detached, openCountingDocument{DSSDocument: document, opens: &opens[i]})
	}

	// b64=true (the default): the digest is computed over the base64url encoding of the document
	digestOfEncoded := func(document model.DSSDocument) string {
		encoded, err := DSSJsonUtilsToBase64UrlDocument(document)
		if err != nil {
			t.Fatal(err)
		}
		value, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, []byte(encoded))
		if err != nil {
			t.Fatal(err)
		}
		return jose.Base64URLEncode(value)
	}

	var pars, hashes []string
	// entry 0: matches detached[1] by digest although named after nothing
	pars = append(pars, "unknown0")
	hashes = append(hashes, digestOfEncoded(detached[1]))
	// entry 1: wrong digest, but named after detached[2] (name fallback), so found but not intact
	pars = append(pars, "doc2.txt")
	hashes = append(hashes, jose.Base64URLEncode([]byte("wrong digest value 0123456789abcd")))
	// remaining entries: neither a digest nor a name match anything
	for i := 2; i < entries; i++ {
		pars = append(pars, fmt.Sprintf("missing%d", i))
		hashes = append(hashes, jose.Base64URLEncode([]byte(fmt.Sprintf("no such digest %02d 0123456789abcdef", i))))
	}
	for i := range opens {
		opens[i] = 0 // computing the expected digest above read the document
	}

	quote := func(values []string) string { return `["` + strings.Join(values, `","`) + `"]` }
	header := `{"alg":"ES256","sigD":{"mId":"http://uri.etsi.org/19182/ObjectIdByURIHash","pars":` + quote(pars) +
		`,"hashM":"S256","hashV":` + quote(hashes) + `}}`
	jws, err := NewJWSFromCompactSerializationParts([]string{jose.Base64URLEncode([]byte(header)), "", jose.Base64URLEncode([]byte("signature"))})
	if err != nil {
		t.Fatal(err)
	}
	signature := NewSignature(jws)
	signature.SetDetachedContents(detached)

	validations := signature.referenceValidationsByUriHashMechanism()
	if len(validations) != entries {
		t.Fatalf("got %d reference validations, want %d", len(validations), entries)
	}

	if !validations[0].IsFound() || !validations[0].IsIntact() ||
		validations[0].Document().Name() != "doc1.txt" {
		t.Errorf("entry 0: found=%v intact=%v document=%v, want the digest match doc1.txt",
			validations[0].IsFound(), validations[0].IsIntact(), validations[0].Document())
	}
	if !validations[1].IsFound() || validations[1].IsIntact() ||
		validations[1].Document().Name() != "doc2.txt" {
		t.Errorf("entry 1: found=%v intact=%v document=%v, want the name match doc2.txt, not intact",
			validations[1].IsFound(), validations[1].IsIntact(), validations[1].Document())
	}
	for i := 2; i < entries; i++ {
		if validations[i].IsFound() || validations[i].IsIntact() {
			t.Errorf("entry %d: found=%v intact=%v, want neither", i, validations[i].IsFound(), validations[i].IsIntact())
		}
	}

	for i, count := range opens {
		if count > 1 {
			t.Errorf("detached document %d was read %d times, want at most once", i, count)
		}
	}
}
