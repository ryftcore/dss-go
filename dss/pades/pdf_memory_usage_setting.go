// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfMemoryUsageSetting.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). Java's protected constructors plus package-visibility-free
// static factories collapse to a value struct plus free constructor functions, following the
// PAdESUtils<Method> naming convention already used for other flattened Java static members in
// this package (see pades_utils.go's header): PdfMemoryUsageSetting<Method>. Only
// PdfMemoryUsageSettingMemoryFull is referenced by a landed sibling (pades_utils.go's
// PAdESUtilsDefaultPdfMemoryUsageSetting); the rest are ported for completeness of the type, as
// native_pdf_document_reader.go's header documents this whole setting is accepted and ignored by
// the native PDF engine, which always works from an in-memory byte slice.
package pades

import "fmt"

// PdfMemoryUsageSettingMode is the memory usage setting's load option. Port of
// PdfMemoryUsageSetting.Mode.
type PdfMemoryUsageSettingMode int

const (
	// PdfMemoryUsageSettingModeMemoryFull loads the file to memory prior to reading.
	// Port of Mode.MEMORY_FULL.
	PdfMemoryUsageSettingModeMemoryFull PdfMemoryUsageSettingMode = iota

	// PdfMemoryUsageSettingModeMemoryBuffered loads the file to memory during reading.
	// Port of Mode.MEMORY_BUFFERED.
	PdfMemoryUsageSettingModeMemoryBuffered

	// PdfMemoryUsageSettingModeFile connects to the file in the filesystem during reading.
	// Port of Mode.FILE.
	PdfMemoryUsageSettingModeFile

	// PdfMemoryUsageSettingModeMixed loads a portion of the file to memory during reading, and
	// handles the rest in a temporary file. Port of Mode.MIXED.
	PdfMemoryUsageSettingModeMixed
)

// String ports the implicit Enum#name/toString.
func (m PdfMemoryUsageSettingMode) String() string {
	switch m {
	case PdfMemoryUsageSettingModeMemoryFull:
		return "MEMORY_FULL"
	case PdfMemoryUsageSettingModeMemoryBuffered:
		return "MEMORY_BUFFERED"
	case PdfMemoryUsageSettingModeFile:
		return "FILE"
	case PdfMemoryUsageSettingModeMixed:
		return "MIXED"
	default:
		return "UNKNOWN"
	}
}

// PdfMemoryUsageSetting represents the PDF document loading setting on signature creation or
// validation. Port of the PdfMemoryUsageSetting class.
type PdfMemoryUsageSetting struct {
	// Mode is the chosen PDF memory usage mode. Port of #getMode.
	Mode PdfMemoryUsageSettingMode

	// MaxMemoryBytes is the maximum number of bytes to be stored in memory, applicable only when
	// Mode is MemoryBuffered or Mixed; -1 means unrestricted. Port of #getMaxMemoryBytes.
	MaxMemoryBytes int64

	// MaxStorageBytes is the maximum number of bytes the in-memory and temporary file may have
	// together, applicable only when Mode is File or Mixed; -1 means unrestricted.
	// Port of #getMaxStorageBytes.
	MaxStorageBytes int64
}

// PdfMemoryUsageSettingMemoryFull represents a memory-forced type of handling: the file is read
// and loaded to a byte array before processing. Port of #memoryFull.
func PdfMemoryUsageSettingMemoryFull() PdfMemoryUsageSetting {
	return PdfMemoryUsageSetting{Mode: PdfMemoryUsageSettingModeMemoryFull, MaxMemoryBytes: -1, MaxStorageBytes: -1}
}

// PdfMemoryUsageSettingMemoryBuffered represents memory unrestricted allocation size load mode:
// the file is loaded to memory as needed. Port of #memoryBuffered().
func PdfMemoryUsageSettingMemoryBuffered() PdfMemoryUsageSetting {
	return PdfMemoryUsageSettingMemoryBufferedWithMaxBytes(-1)
}

// PdfMemoryUsageSettingMemoryBufferedWithMaxBytes represents memory allocation size load mode
// restricted to maxBytes. Port of #memoryBuffered(long).
func PdfMemoryUsageSettingMemoryBufferedWithMaxBytes(maxBytes int64) PdfMemoryUsageSetting {
	return PdfMemoryUsageSetting{Mode: PdfMemoryUsageSettingModeMemoryBuffered, MaxMemoryBytes: maxBytes, MaxStorageBytes: -1}
}

// PdfMemoryUsageSettingFileOnly represents file only unrestricted allocation size load mode: the
// content of a file is stored in a temporary file in the filesystem. Port of #fileOnly().
func PdfMemoryUsageSettingFileOnly() PdfMemoryUsageSetting {
	return PdfMemoryUsageSettingFileOnlyWithMaxStorageBytes(-1)
}

// PdfMemoryUsageSettingFileOnlyWithMaxStorageBytes represents file only allocation size load
// mode restricted to maxStorageBytes. Port of #fileOnly(long).
func PdfMemoryUsageSettingFileOnlyWithMaxStorageBytes(maxStorageBytes int64) PdfMemoryUsageSetting {
	return PdfMemoryUsageSetting{Mode: PdfMemoryUsageSettingModeFile, MaxMemoryBytes: -1, MaxStorageBytes: maxStorageBytes}
}

// PdfMemoryUsageSettingMixed represents mixed memory-first unrestricted file allocation size load
// mode. Port of #mixed(long).
func PdfMemoryUsageSettingMixed(maxMemoryBytes int64) PdfMemoryUsageSetting {
	return PdfMemoryUsageSettingMixedWithMaxStorageBytes(maxMemoryBytes, -1)
}

// PdfMemoryUsageSettingMixedWithMaxStorageBytes represents mixed memory-first restricted file
// allocation size load mode. Port of #mixed(long, long).
func PdfMemoryUsageSettingMixedWithMaxStorageBytes(maxMemoryBytes, maxStorageBytes int64) PdfMemoryUsageSetting {
	return PdfMemoryUsageSetting{Mode: PdfMemoryUsageSettingModeMixed, MaxMemoryBytes: maxMemoryBytes, MaxStorageBytes: maxStorageBytes}
}

// String ports #toString.
func (s PdfMemoryUsageSetting) String() string {
	return fmt.Sprintf("PdfMemoryUsageSetting[%s, %d, %d]", s.Mode, s.MaxMemoryBytes, s.MaxStorageBytes)
}
