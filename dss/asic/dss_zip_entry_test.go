package asic

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"testing"
	"time"
)

// TestNewDSSZipEntryDefaults pins the DEFLATED default and the empty metadata a freshly named entry
// carries.
func TestNewDSSZipEntryDefaults(t *testing.T) {
	entry := NewDSSZipEntry("META-INF/signature001.p7s")
	if entry.Name() != "META-INF/signature001.p7s" {
		t.Errorf("name = %q", entry.Name())
	}
	if entry.CompressionMethod() != int(zip.Deflate) {
		t.Errorf("compressionMethod = %d, want DEFLATED (8)", entry.CompressionMethod())
	}
	if entry.Comment() != "" || entry.Extra() != nil {
		t.Errorf("comment = %q / extra = %x, want empty", entry.Comment(), entry.Extra())
	}
	if !entry.CreationTime().IsZero() || !entry.ModificationTime().IsZero() || !entry.LastAccessTime().IsZero() {
		t.Error("a fresh entry must carry no timestamps")
	}
	if entry.Size() != 0 || entry.CompressedSize() != 0 || entry.Crc() != 0 {
		t.Error("a fresh entry must carry no sizes")
	}
}

// TestDSSZipEntryCreateZipEntryCopiesOnlyStableFields pins createZipEntry's documented contract:
// name, comment, method and extra are copied; sizes, CRC and the modification time are not, because
// the writer recomputes them.
func TestDSSZipEntryCreateZipEntryCopiesOnlyStableFields(t *testing.T) {
	entry := NewDSSZipEntry("data.bin")
	entry.SetComment("a comment")
	entry.SetCompressionMethod(int(zip.Store))
	entry.SetExtra([]byte{0x55, 0x54, 0x05, 0x00, 0x01, 0x01, 0x02, 0x03, 0x04})

	fileHeader := entry.CreateZipEntry()
	if fileHeader.Name != "data.bin" || fileHeader.Comment != "a comment" {
		t.Errorf("name/comment = %q/%q", fileHeader.Name, fileHeader.Comment)
	}
	if fileHeader.Method != zip.Store {
		t.Errorf("method = %d, want STORED", fileHeader.Method)
	}
	if !bytes.Equal(fileHeader.Extra, entry.Extra()) {
		t.Errorf("extra = %x, want %x", fileHeader.Extra, entry.Extra())
	}
	if fileHeader.CRC32 != 0 || fileHeader.UncompressedSize64 != 0 || fileHeader.CompressedSize64 != 0 {
		t.Error("createZipEntry must not carry sizes or CRC over")
	}
	if !fileHeader.Modified.IsZero() || fileHeader.ModifiedDate != 0 || fileHeader.ModifiedTime != 0 {
		t.Error("createZipEntry must not carry a modification time over")
	}
}

// dssZipEntryExtendedTimestampExtra builds an Info-ZIP extended timestamp (0x5455) extra field.
func dssZipEntryExtendedTimestampExtra(flags byte, times ...int32) []byte {
	payload := []byte{flags}
	for _, value := range times {
		payload = binary.LittleEndian.AppendUint32(payload, uint32(value))
	}
	extra := make([]byte, 0, 4+len(payload))
	extra = binary.LittleEndian.AppendUint16(extra, 0x5455)
	extra = binary.LittleEndian.AppendUint16(extra, uint16(len(payload)))
	return append(extra, payload...)
}

// TestDSSZipEntryDecodesExtendedTimestampExtra pins the port of java.util.zip.ZipEntry#setExtra0's
// EXTID_EXTT branch, which is where getCreationTime()/getLastAccessTime() come from.
func TestDSSZipEntryDecodesExtendedTimestampExtra(t *testing.T) {
	const mtime, atime, ctime = 1600000000, 1600000001, 1600000002
	fileHeader := &zip.FileHeader{
		Name:  "data.bin",
		Extra: dssZipEntryExtendedTimestampExtra(0x1|0x2|0x4, mtime, atime, ctime),
	}
	entry := NewDSSZipEntryFromFileHeader(fileHeader)
	if got := entry.ModificationTime().Unix(); got != mtime {
		t.Errorf("modificationTime = %d, want %d", got, mtime)
	}
	if got := entry.LastAccessTime().Unix(); got != atime {
		t.Errorf("lastAccessTime = %d, want %d", got, atime)
	}
	if got := entry.CreationTime().Unix(); got != ctime {
		t.Errorf("creationTime = %d, want %d", got, ctime)
	}

	// Only the modification-time bit set: the other two stay absent, and the trailing bytes
	// are not misread as further timestamps.
	entry = NewDSSZipEntryFromFileHeader(&zip.FileHeader{
		Name:  "data.bin",
		Extra: dssZipEntryExtendedTimestampExtra(0x1, mtime),
	})
	if got := entry.ModificationTime().Unix(); got != mtime {
		t.Errorf("modificationTime = %d, want %d", got, mtime)
	}
	if !entry.LastAccessTime().IsZero() || !entry.CreationTime().IsZero() {
		t.Error("access/creation times must stay absent when their flag bits are clear")
	}
}

// TestDSSZipEntryDecodesNTFSExtra pins the EXTID_NTFS branch, including its
// WINDOWS_TIME_NOT_AVAILABLE skip. The vector is taken from a real corpus entry
// (dss-asic-cades/src/test/resources/signable/test.zip's central directory).
func TestDSSZipEntryDecodesNTFSExtra(t *testing.T) {
	extra, err := hex.DecodeString("0a0020000000000001001800f6a5d35f804bd501f6a5d35f804bd501f6a5d35f804bd501")
	if err != nil {
		t.Fatalf("decode fixture extra: %v", err)
	}
	entry := NewDSSZipEntryFromFileHeader(&zip.FileHeader{Name: "test.txt", Extra: extra})
	// 0x01d54b805fd3a5f6 FILETIME ticks (100ns since 1601-01-01 UTC) == 2019-08-05T11:24:41.2900854Z
	want := time.Date(2019, time.August, 5, 11, 24, 41, 290085400, time.UTC)
	for _, tc := range []struct {
		label string
		got   time.Time
	}{
		{"modificationTime", entry.ModificationTime()},
		{"lastAccessTime", entry.LastAccessTime()},
		{"creationTime", entry.CreationTime()},
	} {
		if !tc.got.Equal(want) {
			t.Errorf("%s = %s, want %s", tc.label, tc.got.UTC(), want)
		}
	}

	// A zeroed FILETIME is WINDOWS_TIME_NOT_AVAILABLE and must be skipped, not decoded as
	// year 1601.
	zeroed := append([]byte(nil), extra...)
	for i := 8; i < 8+24; i++ {
		zeroed[i] = 0
	}
	entry = NewDSSZipEntryFromFileHeader(&zip.FileHeader{Name: "test.txt", Extra: zeroed})
	if !entry.CreationTime().IsZero() || !entry.LastAccessTime().IsZero() {
		t.Error("WINDOWS_TIME_NOT_AVAILABLE must leave the timestamps absent")
	}
}

// TestDSSZipEntryIgnoresMalformedExtra pins the length guards: a truncated or empty extra field is
// skipped rather than read past.
func TestDSSZipEntryIgnoresMalformedExtra(t *testing.T) {
	for _, extra := range [][]byte{
		{0x55},                               // shorter than a header
		{0x55, 0x54, 0x20, 0x00, 0x01},       // declared size exceeds the payload
		{0x55, 0x54, 0x00, 0x00},             // zero-length EXTT payload
		{0x0a, 0x00, 0x04, 0x00, 0, 0, 0, 0}, // NTFS shorter than 32 bytes
	} {
		entry := NewDSSZipEntryFromFileHeader(&zip.FileHeader{Name: "x", Extra: extra})
		if !entry.CreationTime().IsZero() || !entry.LastAccessTime().IsZero() {
			t.Errorf("extra %x produced timestamps", extra)
		}
	}
}

// TestDSSZipEntryFromFileHeaderCopiesSizes pins that the sizes/CRC a fully-read entry carries land
// on the DSSZipEntry, which is what the round-trip assertions rely on.
func TestDSSZipEntryFromFileHeaderCopiesSizes(t *testing.T) {
	modified := time.Date(2021, time.June, 15, 10, 30, 44, 0, time.UTC)
	entry := NewDSSZipEntryFromFileHeader(&zip.FileHeader{
		Name:               "data.bin",
		Comment:            "c",
		Method:             zip.Deflate,
		Modified:           modified,
		CRC32:              0xdeadbeef,
		CompressedSize64:   11,
		UncompressedSize64: 22,
	})
	if entry.Crc() != 0xdeadbeef || entry.CompressedSize() != 11 || entry.Size() != 22 {
		t.Errorf("crc/csize/size = %d/%d/%d", entry.Crc(), entry.CompressedSize(), entry.Size())
	}
	if entry.Comment() != "c" || entry.CompressionMethod() != int(zip.Deflate) {
		t.Errorf("comment/method = %q/%d", entry.Comment(), entry.CompressionMethod())
	}
	// With no extra field, getLastModifiedTime() falls back to the MS-DOS timestamp.
	if !entry.ModificationTime().Equal(modified) {
		t.Errorf("modificationTime = %s, want %s", entry.ModificationTime(), modified)
	}
}

// TestDSSZipEntryEquals pins the field set equals() compares.
func TestDSSZipEntryEquals(t *testing.T) {
	build := func() *DSSZipEntry {
		entry := NewDSSZipEntry("a.txt")
		entry.SetComment("c")
		entry.SetExtra([]byte{1, 2, 3})
		entry.SetCreationTime(time.Unix(1600000000, 0).UTC())
		return entry
	}
	first, second := build(), build()
	if !first.Equals(second) {
		t.Fatal("identical entries must be equal")
	}
	if !first.Equals(first) {
		t.Fatal("an entry must equal itself")
	}
	if first.Equals(nil) {
		t.Fatal("an entry must not equal nil")
	}

	for _, mutate := range []func(*DSSZipEntry){
		func(e *DSSZipEntry) { e.SetName("b.txt") },
		func(e *DSSZipEntry) { e.SetComment("other") },
		func(e *DSSZipEntry) { e.SetCompressionMethod(int(zip.Store)) },
		func(e *DSSZipEntry) { e.SetExtra([]byte{9}) },
		func(e *DSSZipEntry) { e.SetCreationTime(time.Unix(1, 0).UTC()) },
	} {
		mutated := build()
		mutate(mutated)
		if first.Equals(mutated) {
			t.Errorf("entries differing in one field must not be equal (%+v)", mutated)
		}
	}
}

// TestNewDSSZipEntryFromFileHeaderPanicsOnNil pins the Objects.requireNonNull message.
func TestNewDSSZipEntryFromFileHeaderPanicsOnNil(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != "ZipEntry cannot be null!" {
			t.Fatalf("panic = %v, want the Java message", recovered)
		}
	}()
	NewDSSZipEntryFromFileHeader(nil)
}
