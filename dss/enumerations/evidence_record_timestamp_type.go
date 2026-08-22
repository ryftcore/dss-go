// Ported from dss-enumerations/.../EvidenceRecordTimestampType.java (DSS 6.5.RC1).
//
// Different types of Evidence Record time-stamps.
package enumerations

import "fmt"

// EvidenceRecordTimestampType represents the different types of Evidence
// Record time-stamps.
type EvidenceRecordTimestampType string

const (
	// EvidenceRecordTimestampTypeArchiveTimestamp is the initial archive
	// time-stamp.
	EvidenceRecordTimestampTypeArchiveTimestamp EvidenceRecordTimestampType = "ARCHIVE_TIMESTAMP"
	// EvidenceRecordTimestampTypeTimestampRenewalArchiveTimestamp is
	// the time-stamp used to renew a previous archive time-stamp within the
	// same ArchiveTimestampChain.
	EvidenceRecordTimestampTypeTimestampRenewalArchiveTimestamp EvidenceRecordTimestampType = "TIMESTAMP_RENEWAL_ARCHIVE_TIMESTAMP"
	// EvidenceRecordTimestampTypeHashTreeRenewalArchiveTimestamp is
	// the time-stamp used to renew a hash-tree, starting a new
	// ArchiveTimeStampSequence.
	EvidenceRecordTimestampTypeHashTreeRenewalArchiveTimestamp EvidenceRecordTimestampType = "HASH_TREE_RENEWAL_ARCHIVE_TIMESTAMP"
)

// EvidenceRecordTimestampTypeValues returns all constants in declaration order.
func EvidenceRecordTimestampTypeValues() []EvidenceRecordTimestampType {
	return []EvidenceRecordTimestampType{
		EvidenceRecordTimestampTypeArchiveTimestamp,
		EvidenceRecordTimestampTypeTimestampRenewalArchiveTimestamp,
		EvidenceRecordTimestampTypeHashTreeRenewalArchiveTimestamp,
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
