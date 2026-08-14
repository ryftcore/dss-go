// Ported from dss-enumerations/.../EvidenceRecordTimestampType.java (DSS 6.5.RC1).
//
// Different types of Evidence Record time-stamps.
package enumerations

import "fmt"

// EvidenceRecordTimestampType represents the different types of Evidence
// Record time-stamps.
type EvidenceRecordTimestampType string

const (
	// EvidenceRecordTimestampType_ARCHIVE_TIMESTAMP is the initial archive
	// time-stamp.
	EvidenceRecordTimestampType_ARCHIVE_TIMESTAMP EvidenceRecordTimestampType = "ARCHIVE_TIMESTAMP"
	// EvidenceRecordTimestampType_TIMESTAMP_RENEWAL_ARCHIVE_TIMESTAMP is
	// the time-stamp used to renew a previous archive time-stamp within the
	// same ArchiveTimestampChain.
	EvidenceRecordTimestampType_TIMESTAMP_RENEWAL_ARCHIVE_TIMESTAMP EvidenceRecordTimestampType = "TIMESTAMP_RENEWAL_ARCHIVE_TIMESTAMP"
	// EvidenceRecordTimestampType_HASH_TREE_RENEWAL_ARCHIVE_TIMESTAMP is
	// the time-stamp used to renew a hash-tree, starting a new
	// ArchiveTimeStampSequence.
	EvidenceRecordTimestampType_HASH_TREE_RENEWAL_ARCHIVE_TIMESTAMP EvidenceRecordTimestampType = "HASH_TREE_RENEWAL_ARCHIVE_TIMESTAMP"
)

// EvidenceRecordTimestampTypeValues returns all constants in declaration order.
func EvidenceRecordTimestampTypeValues() []EvidenceRecordTimestampType {
	return []EvidenceRecordTimestampType{
		EvidenceRecordTimestampType_ARCHIVE_TIMESTAMP,
		EvidenceRecordTimestampType_TIMESTAMP_RENEWAL_ARCHIVE_TIMESTAMP,
		EvidenceRecordTimestampType_HASH_TREE_RENEWAL_ARCHIVE_TIMESTAMP,
	}
}

// EvidenceRecordTimestampTypeValueOf returns the EvidenceRecordTimestampType matching the given Java enum name.
func EvidenceRecordTimestampTypeValueOf(name string) (EvidenceRecordTimestampType, error) {
	for _, v := range EvidenceRecordTimestampTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant EvidenceRecordTimestampType.%s", name)
}
