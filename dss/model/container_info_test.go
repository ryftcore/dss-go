package model

import (
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestContainerInfoRoundTrip(t *testing.T) {
	ci := NewContainerInfo()
	ci.SetContainerType(enumerations.ASiCContainerType_ASiC_S)
	ci.SetZipComment("comment")
	ci.SetMimeTypeContent("application/vnd.etsi.asic-s+zip")
	ci.SetSignedDocumentFilenames([]string{"a.txt", "b.txt"})
	mf := NewManifestFile()
	ci.SetManifestFiles([]*ManifestFile{mf})

	if ci.ContainerType() != enumerations.ASiCContainerType_ASiC_S {
		t.Fatalf("ContainerType() = %v", ci.ContainerType())
	}
	if !ci.IsMimeTypeFilePresent() {
		t.Fatal("expected IsMimeTypeFilePresent() to be true")
	}
	if !reflect.DeepEqual(ci.SignedDocumentFilenames(), []string{"a.txt", "b.txt"}) {
		t.Fatalf("SignedDocumentFilenames() = %v", ci.SignedDocumentFilenames())
	}
	if len(ci.ManifestFiles()) != 1 || ci.ManifestFiles()[0] != mf {
		t.Fatal("ManifestFiles() did not round-trip")
	}
}

func TestContainerInfoMimeTypeFileAbsent(t *testing.T) {
	ci := NewContainerInfo()
	if ci.IsMimeTypeFilePresent() {
		t.Fatal("expected IsMimeTypeFilePresent() to be false when unset")
	}
}
