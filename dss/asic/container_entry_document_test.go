package asic

import (
	"archive/zip"
	"io"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TestNewContainerEntryDocument pins the constructor: the wrapper adopts the content's name and
// mime-type, and creates a DSSZipEntry bearing the same name.
func TestNewContainerEntryDocument(t *testing.T) {
	content := model.NewInMemoryDocumentWithName([]byte("<sig/>"), "META-INF/signatures001.xml")
	entryDocument := NewContainerEntryDocument(content)

	if entryDocument.Name() != "META-INF/signatures001.xml" {
		t.Errorf("name = %q", entryDocument.Name())
	}
	if entryDocument.MimeType() != content.MimeType() {
		t.Error("mimeType was not adopted from the content")
	}
	if entryDocument.ZipEntry().Name() != entryDocument.Name() {
		t.Errorf("zipEntry name = %q, want %q", entryDocument.ZipEntry().Name(), entryDocument.Name())
	}
	if entryDocument.ZipEntry().CompressionMethod() != int(zip.Deflate) {
		t.Errorf("compressionMethod = %d, want DEFLATED", entryDocument.ZipEntry().CompressionMethod())
	}
	if entryDocument.Content() != model.DSSDocument(content) {
		t.Error("the wrapped content is not reachable")
	}

	stream, err := entryDocument.OpenStream()
	if err != nil {
		t.Fatalf("openStream: %v", err)
	}
	data, _ := io.ReadAll(stream)
	_ = stream.Close()
	if string(data) != "<sig/>" {
		t.Errorf("content = %q", data)
	}
}

// TestContainerEntryDocumentSetNameKeepsZipEntryInSync pins the override upstream declares
// precisely so that renaming a document renames the ZIP entry it will be written as.
func TestContainerEntryDocumentSetNameKeepsZipEntryInSync(t *testing.T) {
	entryDocument := NewContainerEntryDocument(model.NewInMemoryDocumentWithName([]byte("x"), "a.txt"))
	entryDocument.SetName("META-INF/b.txt")
	if entryDocument.Name() != "META-INF/b.txt" || entryDocument.ZipEntry().Name() != "META-INF/b.txt" {
		t.Fatalf("name/zipEntry name = %q/%q", entryDocument.Name(), entryDocument.ZipEntry().Name())
	}
}

// TestNewContainerEntryDocumentWithZipEntryRejectsNameMismatch pins the IllegalArgumentException
// that keeps the document name and the ZIP entry name from drifting apart.
func TestNewContainerEntryDocumentWithZipEntryRejectsNameMismatch(t *testing.T) {
	content := model.NewInMemoryDocumentWithName([]byte("x"), "a.txt")

	if _, err := NewContainerEntryDocumentWithZipEntry(content, NewDSSZipEntry("a.txt")); err != nil {
		t.Fatalf("matching names must be accepted: %v", err)
	}
	_, err := NewContainerEntryDocumentWithZipEntry(content, NewDSSZipEntry("b.txt"))
	if err == nil {
		t.Fatal("a name mismatch must be rejected")
	}
	if err.Error() != "Name of the document shall match the name of ZipEntry!" {
		t.Errorf("error = %q, want the Java message", err.Error())
	}
}

// TestNewContainerEntryDocumentPanicsOnMissingInput pins the two requireNonNull messages.
func TestNewContainerEntryDocumentPanicsOnMissingInput(t *testing.T) {
	func() {
		defer func() {
			if recovered := recover(); recovered != "Document content cannot be null!" {
				t.Fatalf("panic = %v, want the Java message", recovered)
			}
		}()
		NewContainerEntryDocument(nil)
	}()

	defer func() {
		if recovered := recover(); recovered != "Document shall contain name!" {
			t.Fatalf("panic = %v, want the Java message", recovered)
		}
	}()
	NewContainerEntryDocument(model.NewInMemoryDocument([]byte("x")))
}

// TestContainerEntryDocumentDigest pins the CommonDocument digest plumbing this type carries its
// own copy of, including its caching.
func TestContainerEntryDocumentDigest(t *testing.T) {
	entryDocument := NewContainerEntryDocument(model.NewInMemoryDocumentWithName([]byte("hello"), "a.txt"))
	first, err := entryDocument.DigestValue(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatalf("digestValue: %v", err)
	}
	// SHA-256("hello")
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got := hexString(first); got != want {
		t.Fatalf("digest = %s, want %s", got, want)
	}
	second, err := entryDocument.Digest(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	if hexString(second.Value()) != want {
		t.Fatalf("cached digest = %s, want %s", hexString(second.Value()), want)
	}
}

func hexString(data []byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, 0, len(data)*2)
	for _, b := range data {
		out = append(out, digits[b>>4], digits[b&0xf])
	}
	return string(out)
}

// TestFileArchiveEntryReadsThroughTheArchive pins that FileArchiveEntry resolves its payload out of
// the archive file lazily, and that its stream can be opened repeatedly (the archive handle is
// owned by the returned stream, not by the document).
func TestFileArchiveEntryReadsThroughTheArchive(t *testing.T) {
	container := zipCoreFileDocument(t, "dss-asic-cades/src/test/resources/signable/test.zip")
	documents, err := NewSecureContainerHandler().ExtractContainerContent(container)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if len(documents) == 0 {
		t.Fatal("fixture produced no entries")
	}
	entry, ok := documents[0].(*FileArchiveEntry)
	if !ok {
		t.Fatalf("entry is %T, want *FileArchiveEntry", documents[0])
	}
	first := zipCoreSha256(t, entry)
	second := zipCoreSha256(t, entry)
	if first != second {
		t.Fatalf("re-reading the entry produced a different digest: %s vs %s", first, second)
	}

	entry.SetName("renamed.txt")
	if entry.ZipEntry().Name() != "renamed.txt" {
		t.Errorf("zipEntry name = %q, want the renamed one", entry.ZipEntry().Name())
	}
	// The payload still resolves: the lookup uses the original entry name, not the document's.
	if third := zipCoreSha256(t, entry); third != first {
		t.Fatalf("renaming the document broke the payload lookup: %s vs %s", third, first)
	}
}
