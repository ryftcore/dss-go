// KAT test for asice_with_xades_manifest_builder.go: the expected XML below is dumped verbatim
// from a Java oracle run over dss-asic-xades 6.5.RC1 (+ dependencies, all installed from the
// upstream source tree at /home/user/dss-upstream via `mvn -o install`) against
// ASiCEWithXAdESManifestBuilder#build() for the fixture built below - see
// /home/user/esig/dss/asic/xades's S7_BRIEF.md task note and the CADSIGN chunk's
// manifest_kat_test.go for the established pattern. The manifest.xml this builder produces is
// signed by the XAdES signature, so the comparison is on exact bytes: element order, namespace
// declaration, attribute order, MimeType strings and the default-MimeType fallback all included.
package xades

import (
	"io"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

func TestASiCEWithXAdESManifestBuilder_KAT(t *testing.T) {
	documents := []model.DSSDocument{
		model.NewInMemoryDocumentWithMimeType([]byte("Hello World !"), "test.text", enumerations.MimeTypeEnum_TEXT),
		model.NewInMemoryDocumentWithMimeType([]byte("<root/>"), "data.xml", enumerations.MimeTypeEnum_XML),
		model.NewInMemoryDocumentWithMimeType([]byte{0x00, 0x01, 0x02, 0x03}, "binary.bin", enumerations.MimeTypeEnum_BINARY),
		// No MimeType at all - exercises the builder's "no MimeType defined" default-to-BINARY
		// fallback, matching the Java oracle input's explicit `(MimeType) null` third argument.
		model.NewInMemoryDocumentWithMimeType([]byte("no-mimetype"), "no-mime.dat", nil),
	}

	manifest, err := NewASiCEWithXAdESManifestBuilder().
		SetDocuments(documents).
		SetManifestFilename("META-INF/manifest.xml").
		Build()
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}

	const wantName = "META-INF/manifest.xml"
	if manifest.Name() != wantName {
		t.Errorf("manifest filename = %q, want %q", manifest.Name(), wantName)
	}

	const wantXML = `<?xml version="1.0" encoding="UTF-8" standalone="no"?><manifest:manifest xmlns:manifest="urn:oasis:names:tc:opendocument:xmlns:manifest:1.0" manifest:version="1.2"><manifest:file-entry manifest:full-path="/" manifest:media-type="application/vnd.etsi.asic-e+zip"/><manifest:file-entry manifest:full-path="test.text" manifest:media-type="text/plain"/><manifest:file-entry manifest:full-path="data.xml" manifest:media-type="text/xml"/><manifest:file-entry manifest:full-path="binary.bin" manifest:media-type="application/octet-stream"/><manifest:file-entry manifest:full-path="no-mime.dat" manifest:media-type="application/octet-stream"/></manifest:manifest>`

	stream, err := manifest.OpenStream()
	if err != nil {
		t.Fatalf("open manifest stream: %v", err)
	}
	defer stream.Close()
	content, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("read manifest stream: %v", err)
	}

	if got := string(content); got != wantXML {
		t.Errorf("manifest bytes differ from the Java oracle\n go   = %s\n java = %s", got, wantXML)
	}
}

// TestASiCEWithXAdESManifestBuilder_EntriesVsDocumentsConflict pins Java's
// getEntries()/IllegalArgumentException("Either DSSDocuments or ManifestEntries shall be
// provided!") for both the both-set and the neither-set case.
func TestASiCEWithXAdESManifestBuilder_EntriesVsDocumentsConflict(t *testing.T) {
	t.Run("neither set", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("Build() did not panic")
			}
			if msg, ok := r.(string); !ok || msg != "Either DSSDocuments or ManifestEntries shall be provided!" {
				t.Errorf("panic = %v, want %q", r, "Either DSSDocuments or ManifestEntries shall be provided!")
			}
		}()
		_, _ = NewASiCEWithXAdESManifestBuilder().SetManifestFilename("META-INF/manifest.xml").Build()
	})

	t.Run("both set", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("Build() did not panic")
			}
			if msg, ok := r.(string); !ok || msg != "Either DSSDocuments or ManifestEntries shall be provided!" {
				t.Errorf("panic = %v, want %q", r, "Either DSSDocuments or ManifestEntries shall be provided!")
			}
		}()
		documents := []model.DSSDocument{model.NewInMemoryDocumentWithMimeType([]byte("a"), "a.txt", enumerations.MimeTypeEnum_TEXT)}
		entries := []*model.ManifestEntry{{}}
		entries[0].SetUri("a.txt")
		_, _ = NewASiCEWithXAdESManifestBuilder().
			SetDocuments(documents).SetEntries(entries).SetManifestFilename("META-INF/manifest.xml").Build()
	})
}
