package asic

import (
	"errors"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// defaultASiCContainerExtractorProbe is a concrete subclass with the CAdES-shaped predicates, used
// to exercise the grouping logic and the Overrides dispatch.
type defaultASiCContainerExtractorProbe struct {
	DefaultASiCContainerExtractor
	calls map[string]int
}

func newDefaultASiCContainerExtractorProbe(asicContainer model.DSSDocument) *defaultASiCContainerExtractorProbe {
	probe := &defaultASiCContainerExtractorProbe{calls: map[string]int{}}
	probe.InitDefaultASiCContainerExtractor(probe, asicContainer)
	return probe
}

func (p *defaultASiCContainerExtractorProbe) IsAllowedManifest(entryName string) bool {
	p.calls["manifest"]++
	return ASiCUtilsIsManifest(entryName)
}

func (p *defaultASiCContainerExtractorProbe) IsAllowedArchiveManifest(entryName string) bool {
	p.calls["archiveManifest"]++
	return ASiCUtilsIsArchiveManifest(entryName)
}

func (p *defaultASiCContainerExtractorProbe) IsAllowedEvidenceRecordManifest(entryName string) bool {
	p.calls["evidenceRecordManifest"]++
	return ASiCUtilsIsEvidenceRecordManifest(entryName)
}

func (p *defaultASiCContainerExtractorProbe) IsAllowedSignature(entryName string) bool {
	p.calls["signature"]++
	return ASiCUtilsIsCAdES(entryName)
}

func (p *defaultASiCContainerExtractorProbe) IsAllowedTimestamp(entryName string) bool {
	p.calls["timestamp"]++
	return ASiCUtilsIsTimestamp(entryName)
}

func (p *defaultASiCContainerExtractorProbe) IsAllowedEvidenceRecord(entryName string) bool {
	p.calls["evidenceRecord"]++
	return ASiCUtilsIsEvidenceRecord(entryName)
}

func (p *defaultASiCContainerExtractorProbe) IsSupportedContainerFormat() bool { return true }

// TestDefaultASiCContainerExtractorGroupsEntries pins zipParsing's grouping, and with it the
// virtual-dispatch contract: every isAllowed* decision must come from the subclass, not from this
// file's (non-existent) defaults.
func TestDefaultASiCContainerExtractorGroupsEntries(t *testing.T) {
	container := zipCoreFileDocument(t, "dss-asic-cades/src/test/resources/validation/evidencerecord/er-asn1-incorrect-hash.asice")
	probe := newDefaultASiCContainerExtractorProbe(container)

	asicContent, err := probe.Extract()
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if probe.calls["signature"] == 0 {
		t.Fatal("the subclass predicates were not reached - the Overrides dispatch is broken")
	}
	if asicContent.AsicContainer() != model.DSSDocument(container) {
		t.Error("the original container was not recorded")
	}
	if asicContent.ContainerType() != enumerations.ASiCContainerTypeASiCE {
		t.Errorf("containerType = %q, want ASiC_E", asicContent.ContainerType())
	}
	if asicContent.MimeTypeDocument() == nil || asicContent.MimeTypeDocument().Name() != ASiCUtilsMimeType {
		t.Error("the mimetype entry was not routed to MimeTypeDocument")
	}
	if len(asicContent.SignedDocuments()) == 0 {
		t.Error("no signed document was collected")
	}
	if len(asicContent.EvidenceRecordDocuments()) == 0 {
		t.Error("no evidence record was collected")
	}
	if len(asicContent.EvidenceRecordManifestDocuments()) == 0 {
		t.Error("no evidence record manifest was collected")
	}
	for _, doc := range asicContent.SignedDocuments() {
		if doc.Name() == ASiCUtilsMimeType {
			t.Error("the mimetype entry must not be a signed document")
		}
	}
	for _, doc := range asicContent.AllDocuments() {
		if doc.Name() == "" {
			t.Error("an entry lost its name during grouping")
		}
	}
}

// TestDefaultASiCContainerExtractorRejectsEmptyContainer pins the IllegalInputException upstream
// raises for an archive with no readable entries.
func TestDefaultASiCContainerExtractorRejectsEmptyContainer(t *testing.T) {
	container := zipCoreFileDocument(t, "dss-asic-cades/src/test/resources/validation/malformed-container.asics")
	probe := newDefaultASiCContainerExtractorProbe(container)

	_, err := probe.Extract()
	if err == nil {
		t.Fatal("expected an IllegalInputException for a container with no entries")
	}
	var illegalInput *exception.IllegalInputException
	want := "The provided file with name 'malformed-container.asics' does not contain documents inside. " +
		"Probably file has an unsupported format or has been corrupted. The signature validation is not possible"
	if !errors.As(err, &illegalInput) || illegalInput.Message != want {
		t.Fatalf("error = %v, want IllegalInputException %q", err, want)
	}
}

// defaultASiCContainerExtractorFactoryProbe registers itself in the ServiceLoader stand-in.
type defaultASiCContainerExtractorFactoryProbe struct {
	supported bool
}

func (f *defaultASiCContainerExtractorFactoryProbe) IsSupported(asicContainer model.DSSDocument) bool {
	return f.supported
}

func (f *defaultASiCContainerExtractorFactoryProbe) Create(asicContainer model.DSSDocument) ASiCContainerExtractor {
	return newDefaultASiCContainerExtractorProbe(asicContainer)
}

// TestDefaultASiCContainerExtractorFromDocument pins the registry that replaces Java's
// ServiceLoader: first supporting factory wins, and an unsupported document is refused with
// upstream's message.
func TestDefaultASiCContainerExtractorFromDocument(t *testing.T) {
	saved := asicContainerExtractorFactoryRegistry
	defer func() { asicContainerExtractorFactoryRegistry = saved }()
	asicContainerExtractorFactoryRegistry = nil

	doc := model.NewInMemoryDocumentWithName([]byte("x"), "x.asice")
	if _, err := DefaultASiCContainerExtractorFromDocument(doc); err == nil {
		t.Fatal("expected a failure with no registered factory")
	} else if err.Error() != "Document format not recognized/handled" {
		t.Fatalf("error = %q, want the Java message", err.Error())
	}

	RegisterASiCContainerExtractorFactory(&defaultASiCContainerExtractorFactoryProbe{supported: false})
	RegisterASiCContainerExtractorFactory(&defaultASiCContainerExtractorFactoryProbe{supported: true})
	extractor, err := DefaultASiCContainerExtractorFromDocument(doc)
	if err != nil {
		t.Fatalf("fromDocument: %v", err)
	}
	if !extractor.IsSupportedContainerFormat() {
		t.Error("the supporting factory was not selected")
	}

	defer func() {
		if recovered := recover(); recovered != "ASiC container cannot be null!" {
			t.Fatalf("panic = %v, want the Java message", recovered)
		}
	}()
	_, _ = DefaultASiCContainerExtractorFromDocument(nil)
}

// TestDefaultASiCContainerExtractorCollectsFolders pins that folder entries are routed to
// getFolders() and never mistaken for signed documents or unsupported files.
func TestDefaultASiCContainerExtractorCollectsFolders(t *testing.T) {
	container := zipCoreFileDocument(t, "dss-asic-cades/src/test/resources/validation/open-document-signed.odt")
	probe := newDefaultASiCContainerExtractorProbe(container)
	asicContent, err := probe.Extract()
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if len(asicContent.Folders()) == 0 {
		t.Fatal("no folder entry was collected from an OpenDocument container")
	}
	for _, folder := range asicContent.Folders() {
		if folder.Name()[len(folder.Name())-1] != '/' {
			t.Errorf("folder %q does not end with '/'", folder.Name())
		}
	}
	for _, doc := range asicContent.SignedDocuments() {
		if doc.Name()[len(doc.Name())-1] == '/' {
			t.Errorf("folder %q was collected as a signed document", doc.Name())
		}
	}
}
