package asic

import (
	"io"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

// countingDocument counts how many times its content is opened.
type countingDocument struct {
	model.DSSDocument
	opens int
}

func (d *countingDocument) OpenStream() (io.ReadCloser, error) {
	d.opens++
	return d.DSSDocument.OpenStream()
}

const asicManifestParserTestManifest = `<?xml version="1.0" encoding="UTF-8"?>
<asic:ASiCManifest xmlns:asic="http://uri.etsi.org/02918/v1.2.1#" xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
  <asic:SigReference URI="META-INF/signature.p7s" MimeType="application/x-pkcs7-signature"/>
  <asic:DataObjectReference URI="test.txt" MimeType="text/plain">
    <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
    <ds:DigestValue>q1MKE+RZFJgrefm34/uplM/R8/si9xzqGvvwK0YMbR0=</ds:DigestValue>
  </asic:DataObjectReference>
</asic:ASiCManifest>`

// TestManifestParserGetManifestFile pins what a document has to look like to be parsed as an ASiC
// manifest (nil otherwise), and that the manifest is parsed only once per lookup (A14B-PERF-001).
func TestManifestParserGetManifestFile(t *testing.T) {
	newDocument := func(content string) *countingDocument {
		return &countingDocument{DSSDocument: model.NewInMemoryDocumentWithName([]byte(content), "META-INF/ASiCManifest.xml")}
	}

	t.Run("manifest", func(t *testing.T) {
		document := newDocument(asicManifestParserTestManifest)
		manifestFile := ManifestParserGetManifestFile(document)
		if manifestFile == nil {
			t.Fatal("expected a manifest file")
		}
		if got := manifestFile.SignatureFilename(); got != "META-INF/signature.p7s" {
			t.Errorf("SignatureFilename() = %q", got)
		}
		if got := len(manifestFile.Entries()); got != 1 {
			t.Errorf("got %d entries, want 1", got)
		}
		// one read for the XML preamble check and one for the (single) parse; the manifest used to
		// be parsed twice more (isDOM, then buildDOM)
		if document.opens > 2 {
			t.Errorf("the manifest document was opened %d times, want at most 2", document.opens)
		}
	})

	t.Run("linked manifest", func(t *testing.T) {
		other := newDocument(`<?xml version="1.0"?><root/>`)
		manifest := newDocument(asicManifestParserTestManifest)
		linked := ManifestParserGetLinkedManifest([]model.DSSDocument{other, manifest}, "META-INF/signature.p7s")
		if linked != model.DSSDocument(manifest) {
			t.Errorf("linked manifest = %v, want the manifest document", linked)
		}
		if ManifestParserGetLinkedManifest([]model.DSSDocument{other, manifest}, "META-INF/other.p7s") != nil {
			t.Error("no manifest expected for an unknown signature name")
		}
	})

	for name, content := range map[string]string{
		"not XML":                    "just some text",
		"empty":                      "",
		"XML that is not a manifest": `<?xml version="1.0"?><root><child/></root>`,
		"truncated XML":              asicManifestParserTestManifest[:len(asicManifestParserTestManifest)/2],
		"no XML preamble":            "<asic:ASiCManifest",
	} {
		t.Run(name, func(t *testing.T) {
			if manifestFile := ManifestParserGetManifestFile(newDocument(content)); manifestFile != nil {
				t.Errorf("expected nil, got %v", manifestFile)
			}
		})
	}
}
