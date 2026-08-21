// Tests for DetachedTimestampSource, matching
// dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/timestamp/DetachedTimestampSource.java
// upstream.
package timestamp

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

// fakeEvidenceRecord is a partial validation.EvidenceRecord double (same nil-embed technique as
// spi/validation's own determinism test fakes), overriding only what
// DetachedTimestampSource.AddExternalEvidenceRecord's call graph touches when the evidence
// record itself carries no timestamps/certificates/manifest.
type fakeEvidenceRecord struct {
	validation.EvidenceRecord
	id           string
	manifestFile *model.ManifestFile
}

func (f *fakeEvidenceRecord) Id() string                               { return f.id }
func (f *fakeEvidenceRecord) Timestamps() []*validation.TimestampToken { return nil }
func (f *fakeEvidenceRecord) DetachedEvidenceRecords() []validation.EvidenceRecord {
	return nil
}
func (f *fakeEvidenceRecord) TimestampedReferences() []*validation.TimestampedReference {
	return nil
}
func (f *fakeEvidenceRecord) SetTimestampedReferences([]*validation.TimestampedReference) {}
func (f *fakeEvidenceRecord) ManifestFile() *model.ManifestFile                           { return f.manifestFile }

var _ validation.EvidenceRecord = (*fakeEvidenceRecord)(nil)

func TestNewDetachedTimestampSourceDefaults(t *testing.T) {
	source := NewDetachedTimestampSource()
	if source.DetachedTimestamps() != nil {
		t.Fatalf("DetachedTimestamps() = %v, want nil for a fresh source", source.DetachedTimestamps())
	}
}

func TestNewDetachedTimestampSourceWithTimestamp(t *testing.T) {
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	source := NewDetachedTimestampSourceWithTimestamp(token)
	got := source.DetachedTimestamps()
	if len(got) != 1 || got[0] != token {
		t.Fatalf("DetachedTimestamps() = %v, want [token]", got)
	}
}

func TestAddExternalTimestampAppendsAndPopulatesCertificateSource(t *testing.T) {
	source := NewDetachedTimestampSource()
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)

	if err := source.AddExternalTimestamp(token); err != nil {
		t.Fatalf("AddExternalTimestamp: %v", err)
	}

	got := source.DetachedTimestamps()
	if len(got) != 1 || got[0] != token {
		t.Fatalf("DetachedTimestamps() = %v, want [token]", got)
	}
	// timestamp-token.tst embeds the TSA and CA certificates (see the KAT fixture header
	// comment in spi/validation/timestamp_token_kat_test.go); populateSources must have merged
	// them into the source's own certificate source.
	if source.certificateSource.NumberOfCertificates() == 0 {
		t.Fatal("expected AddExternalTimestamp to populate the merged certificate source from the token")
	}
}

func TestAddExternalTimestampAccumulates(t *testing.T) {
	source := NewDetachedTimestampSource()
	first := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	second := loadFixtureTimestampToken(t, enumerations.TimestampType_ARCHIVE_TIMESTAMP)

	if err := source.AddExternalTimestamp(first); err != nil {
		t.Fatalf("AddExternalTimestamp(first): %v", err)
	}
	if err := source.AddExternalTimestamp(second); err != nil {
		t.Fatalf("AddExternalTimestamp(second): %v", err)
	}

	got := source.DetachedTimestamps()
	if len(got) != 2 || got[0] != first || got[1] != second {
		t.Fatalf("DetachedTimestamps() = %v, want [first, second] in insertion order", got)
	}
}

func TestAddExternalEvidenceRecordAppends(t *testing.T) {
	source := NewDetachedTimestampSource()
	er := &fakeEvidenceRecord{id: "er-1"}

	if err := source.AddExternalEvidenceRecord(er); err != nil {
		t.Fatalf("AddExternalEvidenceRecord: %v", err)
	}
	if len(source.detachedEvidenceRecords) != 1 || source.detachedEvidenceRecords[0] != er {
		t.Fatalf("detachedEvidenceRecords = %v, want [er]", source.detachedEvidenceRecords)
	}
}

// ---- isCoveredTimestamp -----------------------------------------------------------------------

func TestIsCoveredTimestampNilManifestAlwaysCovered(t *testing.T) {
	source := NewDetachedTimestampSource()
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	er := &fakeEvidenceRecord{id: "er-1"} // ManifestFile() returns nil

	if !source.isCoveredTimestamp(er, token) {
		t.Fatal("isCoveredTimestamp() = false, want true when the evidence record carries no manifest")
	}
}

func TestIsCoveredTimestampManifestCoversByFilename(t *testing.T) {
	source := NewDetachedTimestampSource()
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	token.SetFilename("token.tst")

	entry := model.NewManifestEntry()
	entry.SetUri("token.tst")
	manifest := model.NewManifestFile()
	manifest.SetEntries([]*model.ManifestEntry{entry})
	er := &fakeEvidenceRecord{id: "er-1", manifestFile: manifest}

	if !source.isCoveredTimestamp(er, token) {
		t.Fatal("isCoveredTimestamp() = false, want true: the manifest lists an entry matching the token's filename")
	}
}

func TestIsCoveredTimestampManifestDoesNotCover(t *testing.T) {
	source := NewDetachedTimestampSource()
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	token.SetFilename("token.tst")

	entry := model.NewManifestEntry()
	entry.SetUri("other-file.bin")
	manifest := model.NewManifestFile()
	manifest.SetEntries([]*model.ManifestEntry{entry})
	er := &fakeEvidenceRecord{id: "er-1", manifestFile: manifest}

	if source.isCoveredTimestamp(er, token) {
		t.Fatal("isCoveredTimestamp() = true, want false: no manifest entry matches the token's filename")
	}
}

func TestIsCoveredTimestampManifestEmptyFilenameNeverMatches(t *testing.T) {
	source := NewDetachedTimestampSource()
	token := loadFixtureTimestampToken(t, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	// token.Filename() left empty on purpose.

	entry := model.NewManifestEntry()
	entry.SetUri("") // even an empty-URI entry must not match an empty filename
	manifest := model.NewManifestFile()
	manifest.SetEntries([]*model.ManifestEntry{entry})
	er := &fakeEvidenceRecord{id: "er-1", manifestFile: manifest}

	if source.isCoveredTimestamp(er, token) {
		t.Fatal("isCoveredTimestamp() = true, want false: an empty Filename() must never be considered covered")
	}
}
