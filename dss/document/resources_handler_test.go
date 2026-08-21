// Tests for AbstractResourcesHandler, InMemoryResourcesHandler(Builder) and
// TempFileResourcesHandler(Builder), matching
// dss-document/src/main/java/eu/europa/esig/dss/signature/resources/{AbstractResourcesHandler,
// InMemoryResourcesHandler,InMemoryResourcesHandlerBuilder,TempFileResourcesHandler,
// TempFileResourcesHandlerBuilder}.java upstream.
package document

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

// ---- AbstractResourcesHandler -----------------------------------------------

func TestAbstractResourcesHandlerOutputStreamPanicsBeforeCreate(t *testing.T) {
	h := &AbstractResourcesHandler{}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: OutputStream() before CreateOutputStream()")
		}
	}()
	h.OutputStream()
}

func TestAbstractResourcesHandlerCreateOutputStreamTwicePanics(t *testing.T) {
	h := &AbstractResourcesHandler{BuildOutputStream: func() (io.Writer, error) { return io.Discard, nil }}
	if _, err := h.CreateOutputStream(); err != nil {
		t.Fatalf("CreateOutputStream: %s", err)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: OutputStream already created")
		}
	}()
	_, _ = h.CreateOutputStream()
}

func TestAbstractResourcesHandlerCloseNoOpBeforeCreate(t *testing.T) {
	h := &AbstractResourcesHandler{}
	if err := h.Close(); err != nil {
		t.Fatalf("Close() before create: %s", err)
	}
}

// ---- InMemoryResourcesHandler ------------------------------------------------

func TestInMemoryResourcesHandlerRoundTrip(t *testing.T) {
	h := NewInMemoryResourcesHandler()
	w, err := h.CreateOutputStream()
	if err != nil {
		t.Fatalf("CreateOutputStream: %s", err)
	}
	if _, err := w.Write([]byte("hello")); err != nil {
		t.Fatalf("Write: %s", err)
	}

	doc, err := h.WriteToDSSDocument()
	if err != nil {
		t.Fatalf("WriteToDSSDocument: %s", err)
	}
	reader, err := doc.OpenStream()
	if err != nil {
		t.Fatalf("OpenStream: %s", err)
	}
	defer reader.Close()
	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll: %s", err)
	}
	if string(got) != "hello" {
		t.Fatalf("document content = %q, want %q", got, "hello")
	}
}

func TestInMemoryResourcesHandlerBuilderProducesInMemoryResourcesHandler(t *testing.T) {
	b := NewInMemoryResourcesHandlerBuilder()
	h := b.CreateResourcesHandler()
	if _, ok := h.(*InMemoryResourcesHandler); !ok {
		t.Fatalf("CreateResourcesHandler() = %T, want *InMemoryResourcesHandler", h)
	}
}

// ---- TempFileResourcesHandler -------------------------------------------------

func TestTempFileResourcesHandlerRoundTrip(t *testing.T) {
	dir := t.TempDir()
	h, err := NewTempFileResourcesHandler("prefix-", ".tmp", dir)
	if err != nil {
		t.Fatalf("NewTempFileResourcesHandler: %s", err)
	}

	w, err := h.CreateOutputStream()
	if err != nil {
		t.Fatalf("CreateOutputStream: %s", err)
	}
	if _, err := w.Write([]byte("temp file contents")); err != nil {
		t.Fatalf("Write: %s", err)
	}

	doc, err := h.WriteToDSSDocument()
	if err != nil {
		t.Fatalf("WriteToDSSDocument: %s", err)
	}
	reader, err := doc.OpenStream()
	if err != nil {
		t.Fatalf("OpenStream: %s", err)
	}
	defer reader.Close()
	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll: %s", err)
	}
	if string(got) != "temp file contents" {
		t.Fatalf("document content = %q, want %q", got, "temp file contents")
	}

	// WriteToDSSDocument clears toBeDeleted, so the backing file must still exist on disk
	// (the caller owns cleanup of the returned FileDocument from here on).
	fileDoc, ok := doc.(*model.FileDocument)
	if !ok {
		t.Fatalf("WriteToDSSDocument() returned %T, want *model.FileDocument", doc)
	}
	if _, err := os.Stat(fileDoc.Path()); err != nil {
		t.Fatalf("expected the file backing the returned document to survive WriteToDSSDocument: %s", err)
	}
	t.Cleanup(func() { _ = os.Remove(fileDoc.Path()) })
}

func TestTempFileResourcesHandlerCloseWithoutWriteRemovesFile(t *testing.T) {
	dir := t.TempDir()
	h, err := NewTempFileResourcesHandler("prefix-", ".tmp", dir)
	if err != nil {
		t.Fatalf("NewTempFileResourcesHandler: %s", err)
	}
	path := h.tempFilePath
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected the temp file to exist right after creation: %s", err)
	}

	if err := h.Close(); err != nil {
		t.Fatalf("Close: %s", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected the temp file to be removed by Close() (toBeDeleted default true), stat err = %v", err)
	}
}

func TestTempFileResourcesHandlerBuilderDefaults(t *testing.T) {
	b := NewTempFileResourcesHandlerBuilder()
	dir := t.TempDir()
	b.SetTempFileDirectory(dir).SetFileNamePrefix("myprefix-").SetFileNameSuffix(".dat")

	h := b.CreateResourcesHandler()
	tfh, ok := h.(*TempFileResourcesHandler)
	if !ok {
		t.Fatalf("CreateResourcesHandler() = %T, want *TempFileResourcesHandler", h)
	}
	base := filepath.Base(tfh.tempFilePath)
	if filepath.Dir(tfh.tempFilePath) != dir {
		t.Fatalf("temp file created in %q, want %q", filepath.Dir(tfh.tempFilePath), dir)
	}
	if len(base) < len("myprefix-") || base[:len("myprefix-")] != "myprefix-" {
		t.Fatalf("temp file name %q does not start with the configured prefix", base)
	}

	b.Clear()
	if _, err := os.Stat(tfh.tempFilePath); !os.IsNotExist(err) {
		t.Fatalf("expected Clear() to remove the handler's temp file, stat err = %v", err)
	}
}

func TestTempFileResourcesHandlerBuilderCreatesMissingDirectory(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b", "c")

	b := NewTempFileResourcesHandlerBuilder().SetTempFileDirectory(nested)
	h := b.CreateResourcesHandler()
	tfh := h.(*TempFileResourcesHandler)
	if filepath.Dir(tfh.tempFilePath) != nested {
		t.Fatalf("temp file created in %q, want %q", filepath.Dir(tfh.tempFilePath), nested)
	}
	if _, err := os.Stat(nested); err != nil {
		t.Fatalf("expected the missing directory tree to be created: %s", err)
	}
}
