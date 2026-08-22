package asic

import (
	"archive/zip"
	"encoding/binary"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TestASiCParametersAccessors pins the defaults and the round trip of every field.
func TestASiCParametersAccessors(t *testing.T) {
	parameters := NewASiCParameters()
	if parameters.IsZipComment() || parameters.MimeType() != "" || parameters.ContainerType() != "" {
		t.Fatalf("defaults are not the Java ones: %+v", parameters)
	}

	parameters.SetZipComment(true)
	parameters.SetMimeType("application/vnd.etsi.asic-s+zip")
	parameters.SetContainerType(enumerations.ASiCContainerTypeASiCS)
	if !parameters.IsZipComment() || parameters.MimeType() != "application/vnd.etsi.asic-s+zip" ||
		parameters.ContainerType() != enumerations.ASiCContainerTypeASiCS {
		t.Fatalf("accessors did not round trip: %+v", parameters)
	}
}

// TestASiCParametersEquals pins the three fields equals() compares.
func TestASiCParametersEquals(t *testing.T) {
	build := func() *Parameters {
		parameters := NewASiCParameters()
		parameters.SetZipComment(true)
		parameters.SetMimeType("application/vnd.etsi.asic-e+zip")
		parameters.SetContainerType(enumerations.ASiCContainerTypeASiCE)
		return parameters
	}
	first, second := build(), build()
	if !first.Equals(second) || !first.Equals(first) || first.Equals(nil) {
		t.Fatal("equals is not reflexive/consistent")
	}
	for _, mutate := range []func(*Parameters){
		func(p *Parameters) { p.SetZipComment(false) },
		func(p *Parameters) { p.SetMimeType("other") },
		func(p *Parameters) { p.SetContainerType(enumerations.ASiCContainerTypeASiCS) },
	} {
		mutated := build()
		mutate(mutated)
		if first.Equals(mutated) {
			t.Errorf("parameters differing in one field must not be equal: %+v", mutated)
		}
	}
}

// TestASiCParametersMimeTypeResolution pins getMimeType(Parameters): an explicit mimetype wins,
// otherwise the container type decides.
func TestASiCParametersMimeTypeResolution(t *testing.T) {
	parameters := NewASiCParameters()
	parameters.SetContainerType(enumerations.ASiCContainerTypeASiCS)
	if got := UtilsMimeTypeFromParameters(parameters); got != enumerations.MimeType(enumerations.MimeTypeEnumASiCS) {
		t.Errorf("mimeType = %v, want ASICS", got)
	}
	parameters.SetContainerType(enumerations.ASiCContainerTypeASiCE)
	if got := UtilsMimeTypeFromParameters(parameters); got != enumerations.MimeType(enumerations.MimeTypeEnumASiCE) {
		t.Errorf("mimeType = %v, want ASICE", got)
	}
	parameters.SetMimeType("application/vnd.oasis.opendocument.text")
	if got := UtilsMimeTypeFromParameters(parameters); got.MimeTypeString() != "application/vnd.oasis.opendocument.text" {
		t.Errorf("an explicit mimetype must win, got %q", got.MimeTypeString())
	}

	if !UtilsIsASiCE(parameters) || UtilsIsASiCS(parameters) {
		t.Error("isASiCE/isASiCS disagree with the container type")
	}
	defer func() {
		if recovered := recover(); recovered != "ASiCContainerType must be defined!" {
			t.Fatalf("panic = %v, want the Java message", recovered)
		}
	}()
	UtilsIsASiCE(NewASiCParameters())
}

// TestASiCContainerEvidenceRecordParametersEmbedsASiCParameters pins the inheritance: the base
// accessors are promoted and equals() compares both halves.
func TestASiCContainerEvidenceRecordParametersEmbedsASiCParameters(t *testing.T) {
	parameters := NewASiCContainerEvidenceRecordParameters()
	parameters.SetContainerType(enumerations.ASiCContainerTypeASiCE)
	if parameters.ContainerType() != enumerations.ASiCContainerTypeASiCE {
		t.Error("the promoted base accessor did not work")
	}

	manifest := model.NewInMemoryDocumentWithName([]byte("<m/>"), "META-INF/ASiCEvidenceRecordManifest001.xml")
	parameters.SetAsicEvidenceRecordManifest(manifest)
	if parameters.AsicEvidenceRecordManifest() != model.DSSDocument(manifest) {
		t.Error("the manifest did not round trip")
	}

	other := NewASiCContainerEvidenceRecordParameters()
	other.SetContainerType(enumerations.ASiCContainerTypeASiCE)
	if parameters.Equals(other) {
		t.Error("parameters differing in the manifest must not be equal")
	}
	other.SetAsicEvidenceRecordManifest(manifest)
	if !parameters.Equals(other) {
		t.Error("identical parameters must be equal")
	}
	other.SetContainerType(enumerations.ASiCContainerTypeASiCS)
	if parameters.Equals(other) {
		t.Error("equals must also compare the embedded base fields")
	}
}

// TestSecureContainerHandlerDeflatedEntriesUseADataDescriptor pins the other half of deviation 2 in
// zip_utils.go: DEFLATED entries keep java.util.zip.ZipOutputStream's layout, i.e. the
// data-descriptor flag is set and the local header's sizes are zero. The corpus round trip reads
// hundreds of such entries back, so this is the assertion that names the contract.
func TestSecureContainerHandlerDeflatedEntriesUseADataDescriptor(t *testing.T) {
	documents := []model.DSSDocument{
		model.NewInMemoryDocumentWithName([]byte("compress me, compress me, compress me"), "a.txt"),
	}
	archive, err := NewSecureContainerHandler().CreateZipArchive(documents, time.Now(), "")
	if err != nil {
		t.Fatalf("createZipArchive: %v", err)
	}
	raw := zipCoreRawBytes(t, archive)
	header := zipCoreParseFirstLocalHeader(t, raw)
	if header.Method != zip.Deflate {
		t.Fatalf("method = %d, want DEFLATED", header.Method)
	}
	if header.Flags&0x8 == 0 {
		t.Error("a DEFLATED entry must carry the data-descriptor flag, as ZipOutputStream does")
	}
	if header.CRC32 != 0 || header.CompressedSize != 0 || header.UncompressedSize != 0 {
		t.Error("a deferred entry must leave crc/sizes zero in the local header")
	}
	// And the descriptor really is there and readable: the sequential reader recovers the
	// sizes from it.
	extracted, err := NewSecureContainerHandler().ExtractContainerContent(archive)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	entry := extracted[0].(DSSZipEntryDocument).ZipEntry()
	if entry.Size() != int64(len("compress me, compress me, compress me")) {
		t.Errorf("size recovered from the data descriptor = %d", entry.Size())
	}
	if entry.Crc() == 0 {
		t.Error("crc was not recovered from the data descriptor")
	}
	// Sanity: the descriptor signature is present where the reader expects it.
	descriptorOffset := header.DataOffset + int(entry.CompressedSize())
	if binary.LittleEndian.Uint32(raw[descriptorOffset:]) != 0x08074b50 {
		t.Error("no data descriptor signature follows the entry payload")
	}
}
