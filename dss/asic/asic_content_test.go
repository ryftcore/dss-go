package asic

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

func asicContentTestDocument(name string) model.DSSDocument {
	return model.NewInMemoryDocumentWithName([]byte(name), name)
}

// TestASiCContentAllDocumentsOrder pins the concatenation order of getAllDocuments(), which is what
// ZipUtils.createZipArchive writes: "mimetype" first, then signed documents, signatures, manifests,
// archive manifests, evidence record manifests, timestamps, evidence records, unsupported files and
// finally folders.
func TestASiCContentAllDocumentsOrder(t *testing.T) {
	asicContent := NewASiCContent()
	asicContent.SetMimeTypeDocument(asicContentTestDocument("mimetype"))
	asicContent.SetSignedDocuments([]model.DSSDocument{asicContentTestDocument("test.txt")})
	asicContent.SetSignatureDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/signature001.p7s")})
	asicContent.SetManifestDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/ASiCManifest001.xml")})
	asicContent.SetArchiveManifestDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/ASiCArchiveManifest001.xml")})
	asicContent.SetEvidenceRecordManifestDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/ASiCEvidenceRecordManifest001.xml")})
	asicContent.SetTimestampDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/timestamp001.tst")})
	asicContent.SetEvidenceRecordDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/evidencerecord001.xml")})
	asicContent.SetUnsupportedDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/unknown.bin")})
	asicContent.SetFolders([]model.DSSDocument{asicContentTestDocument("folder/")})

	want := []string{
		"mimetype",
		"test.txt",
		"META-INF/signature001.p7s",
		"META-INF/ASiCManifest001.xml",
		"META-INF/ASiCArchiveManifest001.xml",
		"META-INF/ASiCEvidenceRecordManifest001.xml",
		"META-INF/timestamp001.tst",
		"META-INF/evidencerecord001.xml",
		"META-INF/unknown.bin",
		"folder/",
	}
	got := modelNames(asicContent.AllDocuments())
	if len(got) != len(want) {
		t.Fatalf("allDocuments = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("allDocuments[%d] = %q, want %q", i, got[i], want[i])
		}
	}

	// "package.zip" container documents are deliberately NOT part of getAllDocuments(): they
	// are the extracted contents of a signed archive, not container entries.
	asicContent.SetContainerDocuments([]model.DSSDocument{asicContentTestDocument("inner.txt")})
	if len(asicContent.AllDocuments()) != len(want) {
		t.Fatal("container documents must not be written into the archive")
	}
}

// TestASiCContentAllDocumentsSkipsEmptyGroups pins that an absent mimetype and empty groups simply
// contribute nothing.
func TestASiCContentAllDocumentsSkipsEmptyGroups(t *testing.T) {
	asicContent := NewASiCContent()
	if documents := asicContent.AllDocuments(); len(documents) != 0 {
		t.Fatalf("allDocuments = %v, want empty", modelNames(documents))
	}
	asicContent.SetSignedDocuments([]model.DSSDocument{asicContentTestDocument("a.txt")})
	if got := modelNames(asicContent.AllDocuments()); len(got) != 1 || got[0] != "a.txt" {
		t.Fatalf("allDocuments = %v", got)
	}
}

// TestASiCContentAllManifestDocuments pins the three manifest groups getAllManifestDocuments()
// concatenates, in order.
func TestASiCContentAllManifestDocuments(t *testing.T) {
	asicContent := NewASiCContent()
	asicContent.SetManifestDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/ASiCManifest001.xml")})
	asicContent.SetArchiveManifestDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/ASiCArchiveManifest001.xml")})
	asicContent.SetEvidenceRecordManifestDocuments([]model.DSSDocument{asicContentTestDocument("META-INF/ASiCEvidenceRecordManifest001.xml")})

	want := []string{
		"META-INF/ASiCManifest001.xml",
		"META-INF/ASiCArchiveManifest001.xml",
		"META-INF/ASiCEvidenceRecordManifest001.xml",
	}
	got := modelNames(asicContent.AllManifestDocuments())
	if len(got) != len(want) {
		t.Fatalf("allManifestDocuments = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("allManifestDocuments[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

// TestASiCContentRootLevelSignedDocuments pins the root-level filter, which - unlike
// ASiCUtils.getRootLevelDocuments - also rejects backslash-separated paths and keeps "mimetype".
func TestASiCContentRootLevelSignedDocuments(t *testing.T) {
	asicContent := NewASiCContent()
	asicContent.SetSignedDocuments([]model.DSSDocument{
		asicContentTestDocument("a.txt"),
		asicContentTestDocument("folder/b.txt"),
		asicContentTestDocument(`windows\c.txt`),
		asicContentTestDocument("d.txt"),
	})
	got := modelNames(asicContent.RootLevelSignedDocuments())
	if len(got) != 2 || got[0] != "a.txt" || got[1] != "d.txt" {
		t.Fatalf("rootLevelSignedDocuments = %v", got)
	}

	empty := NewASiCContent()
	if documents := empty.RootLevelSignedDocuments(); len(documents) != 0 {
		t.Fatalf("rootLevelSignedDocuments = %v, want empty", modelNames(documents))
	}
}

// TestASiCContentScalarAccessors pins the plain getters/setters other chunks build on.
func TestASiCContentScalarAccessors(t *testing.T) {
	asicContent := NewASiCContent()
	container := asicContentTestDocument("container.asice")
	asicContent.SetAsicContainer(container)
	asicContent.SetContainerType(enumerations.ASiCContainerType_ASiC_E)
	asicContent.SetZipComment("mimetype=application/vnd.etsi.asic-e+zip")

	if asicContent.AsicContainer() != container {
		t.Error("asicContainer round trip failed")
	}
	if asicContent.ContainerType() != enumerations.ASiCContainerType_ASiC_E {
		t.Errorf("containerType = %q", asicContent.ContainerType())
	}
	if asicContent.ZipComment() != "mimetype=application/vnd.etsi.asic-e+zip" {
		t.Errorf("zipComment = %q", asicContent.ZipComment())
	}
}
