package asic

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
)

// zipUtilsRecordingHandler is a ZipContainerHandler that records that it was asked to do the work,
// used to prove the ZipUtils singleton really delegates to the configured builder.
type zipUtilsRecordingHandler struct {
	built *int
}

func (h *zipUtilsRecordingHandler) ExtractContainerContent(zipArchive model.DSSDocument) ([]model.DSSDocument, error) {
	return []model.DSSDocument{model.NewInMemoryDocumentWithName([]byte("stub"), "stub.txt")}, nil
}

func (h *zipUtilsRecordingHandler) ExtractEntryNames(zipArchive model.DSSDocument) ([]string, error) {
	return []string{"stub.txt"}, nil
}

func (h *zipUtilsRecordingHandler) CreateZipArchive(containerEntries []model.DSSDocument, creationTime time.Time, zipComment string) (model.DSSDocument, error) {
	return model.NewInMemoryDocumentWithName([]byte("stub"), "stub.zip"), nil
}

type zipUtilsRecordingBuilder struct {
	built int
}

func (b *zipUtilsRecordingBuilder) Build() ZipContainerHandler {
	b.built++
	return &zipUtilsRecordingHandler{}
}

// TestZipUtilsInstanceIsASingleton pins the singleton and its default handler builder.
func TestZipUtilsInstanceIsASingleton(t *testing.T) {
	first, second := ZipUtilsInstance(), ZipUtilsInstance()
	if first != second {
		t.Fatal("ZipUtilsInstance returned two different instances")
	}
	if _, ok := first.zipContainerHandlerBuilder.(*SecureContainerHandlerBuilder); !ok {
		t.Fatalf("default handler builder is %T, want *SecureContainerHandlerBuilder", first.zipContainerHandlerBuilder)
	}
}

// TestZipUtilsSetZipContainerHandlerBuilder pins that a custom builder is consulted on EVERY call
// (upstream builds a fresh handler each time, because the handler carries per-run counters).
func TestZipUtilsSetZipContainerHandlerBuilder(t *testing.T) {
	zipUtils := ZipUtilsInstance()
	original := zipUtils.zipContainerHandlerBuilder
	defer zipUtils.SetZipContainerHandlerBuilder(original)

	recording := &zipUtilsRecordingBuilder{}
	zipUtils.SetZipContainerHandlerBuilder(recording)

	doc := model.NewInMemoryDocumentWithName([]byte("ignored"), "ignored.zip")
	if names, err := zipUtils.ExtractEntryNames(doc); err != nil || len(names) != 1 || names[0] != "stub.txt" {
		t.Fatalf("ExtractEntryNames = %v (err %v)", names, err)
	}
	if documents, err := zipUtils.ExtractContainerContent(doc); err != nil || len(documents) != 1 {
		t.Fatalf("ExtractContainerContent = %v (err %v)", documents, err)
	}
	if _, err := zipUtils.CreateZipArchiveFromEntries(nil); err != nil {
		t.Fatalf("CreateZipArchiveFromEntries: %v", err)
	}
	if recording.built != 3 {
		t.Fatalf("handler was built %d times, want 3 (one per call)", recording.built)
	}
}

// TestZipUtilsSetZipContainerHandlerBuilderPanicsOnNil pins the Objects.requireNonNull message.
func TestZipUtilsSetZipContainerHandlerBuilderPanicsOnNil(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != "ZipContainerHandlerBuilder shall be defined!" {
			t.Fatalf("panic = %v, want the Java message", recovered)
		}
	}()
	ZipUtilsInstance().SetZipContainerHandlerBuilder(nil)
}

// TestZipUtilsCreateZipArchiveFromASiCContent pins that the Content overload writes
// getAllDocuments() in order and carries the content's zip comment over.
func TestZipUtilsCreateZipArchiveFromASiCContent(t *testing.T) {
	asicContent := NewASiCContent()
	asicContent.SetMimeTypeDocument(asicUtilsCreateMimetypeDocument(zipCoreMimeType("application/vnd.etsi.asic-e+zip")))
	asicContent.SetSignedDocuments([]model.DSSDocument{model.NewInMemoryDocumentWithName([]byte("hello"), "test.txt")})
	asicContent.SetSignatureDocuments([]model.DSSDocument{
		model.NewInMemoryDocumentWithName([]byte("<sig/>"), "META-INF/signatures001.xml"),
	})
	asicContent.SetZipComment("mimetype=application/vnd.etsi.asic-e+zip")

	archive, err := ZipUtilsInstance().CreateZipArchiveAt(asicContent, time.Date(2021, time.June, 15, 10, 30, 44, 0, time.UTC))
	if err != nil {
		t.Fatalf("CreateZipArchiveAt: %v", err)
	}
	names, err := ZipUtilsInstance().ExtractEntryNames(archive)
	if err != nil {
		t.Fatalf("ExtractEntryNames: %v", err)
	}
	want := []string{"mimetype", "test.txt", "META-INF/signatures001.xml"}
	if len(names) != len(want) {
		t.Fatalf("entries = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("entry[%d] = %q, want %q", i, names[i], want[i])
		}
	}
	comment, err := UtilsZipCommentFromArchiveContainer(archive)
	if err != nil {
		t.Fatalf("zip comment: %v", err)
	}
	if comment != asicContent.ZipComment() {
		t.Errorf("zip comment = %q, want %q", comment, asicContent.ZipComment())
	}
}
