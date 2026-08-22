// Ported from dss-enumerations/.../ArchiveTimestampHashIndexVersion.java (DSS 6.5.RC1).
package enumerations

// ArchiveTimestampHashIndexVersion defines a version type of the
// ats-hash-index attribute from the archive-time-stamp-v3.
type ArchiveTimestampHashIndexVersion string

const (
	// ArchiveTimestampHashIndexVersionATSHashIndex is the deprecated
	// ats-hash-index Attribute. See TS 101-733, ch. "6.4.2 ats-hash-index
	// Attribute".
	ArchiveTimestampHashIndexVersionATSHashIndex ArchiveTimestampHashIndexVersion = "ATS_HASH_INDEX"

	// ArchiveTimestampHashIndexVersionATSHashIndexV2 is the deprecated
	// ats-hash-index-v2 Attribute. See ETSI EN 319 122-1 v1.0.0, ch. "5.5.2
	// The ats-hash-index-v2 attribute".
	ArchiveTimestampHashIndexVersionATSHashIndexV2 ArchiveTimestampHashIndexVersion = "ATS_HASH_INDEX_V2"

	// ArchiveTimestampHashIndexVersionATSHashIndexV3 is the
	// ats-hash-index-v3 Attribute. See ETSI EN 319 122-1 v1.1.0, ch. "5.5.2
	// The ats-hash-index-v3 attribute".
	ArchiveTimestampHashIndexVersionATSHashIndexV3 ArchiveTimestampHashIndexVersion = "ATS_HASH_INDEX_V3"
)

type archiveTimestampHashIndexVersionFields struct {
	label string
	oid   string
}

// archiveTimestampHashIndexVersionData holds the (label, oid) pair for each constant.
var archiveTimestampHashIndexVersionData = map[ArchiveTimestampHashIndexVersion]archiveTimestampHashIndexVersionFields{
	ArchiveTimestampHashIndexVersionATSHashIndex:   {"ats-hash-index", "0.4.0.1733.2.5"},
	ArchiveTimestampHashIndexVersionATSHashIndexV2: {"ats-hash-index-v2", "0.4.0.19122.1.4"},
	ArchiveTimestampHashIndexVersionATSHashIndexV3: {"ats-hash-index-v3", "0.4.0.19122.1.5"},
}

// ArchiveTimestampHashIndexVersionValues returns all constants in declaration order.
func ArchiveTimestampHashIndexVersionValues() []ArchiveTimestampHashIndexVersion {
	return []ArchiveTimestampHashIndexVersion{
		ArchiveTimestampHashIndexVersionATSHashIndex,
		ArchiveTimestampHashIndexVersionATSHashIndexV2,
		ArchiveTimestampHashIndexVersionATSHashIndexV3,
	}
}

// Label gets a user-friendly text label.
func (a ArchiveTimestampHashIndexVersion) Label() string {
	return archiveTimestampHashIndexVersionData[a].label
}

// OID gets the unique object identifier.
func (a ArchiveTimestampHashIndexVersion) OID() string {
	return archiveTimestampHashIndexVersionData[a].oid
}

// ArchiveTimestampHashIndexVersionForLabel gets the
// ArchiveTimestampHashIndexVersion for the given label, if found. Returns
// "" (the zero value) if not found, mirroring Java's null return.
func ArchiveTimestampHashIndexVersionForLabel(label string) ArchiveTimestampHashIndexVersion {
	for _, v := range ArchiveTimestampHashIndexVersionValues() {
		if v.Label() == label {
			return v
		}
	}
	return ""
}

// ArchiveTimestampHashIndexVersionForOID gets the
// ArchiveTimestampHashIndexVersion for the given OID, if found. Returns ""
// (the zero value) if not found, mirroring Java's null return.
func ArchiveTimestampHashIndexVersionForOID(oid string) ArchiveTimestampHashIndexVersion {
	for _, v := range ArchiveTimestampHashIndexVersionValues() {
		if v.OID() == oid {
			return v
		}
	}
	return ""
}
