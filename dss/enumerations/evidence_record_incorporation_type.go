// Ported from dss-enumerations/.../EvidenceRecordIncorporationType.java (DSS 6.5.RC1).
//
// Defines an unsigned attribute type within a CAdES signature for
// incorporation of an evidence record (i.e. internal vs external).
package enumerations

import "fmt"

// EvidenceRecordIncorporationType represents the incorporation type of an
// evidence record within a CAdES signature.
type EvidenceRecordIncorporationType string

const (
	// EvidenceRecordIncorporationTypeInternalEvidenceRecord defines the
	// internal-evidence-records attribute (clause 5.2) protecting the whole
	// SignedData instance and used in cases of attached signatures.
	EvidenceRecordIncorporationTypeInternalEvidenceRecord EvidenceRecordIncorporationType = "INTERNAL_EVIDENCE_RECORD"
	// EvidenceRecordIncorporationTypeExternalEvidenceRecord defines the
	// external-evidence-records attribute (clause 5.3) also protecting the
	// whole SignedData instance not containing an eContent element within
	// encapContentInfo (a detached signature), and the external signed
	// data.
	EvidenceRecordIncorporationTypeExternalEvidenceRecord EvidenceRecordIncorporationType = "EXTERNAL_EVIDENCE_RECORD"
)

// EvidenceRecordIncorporationTypeValues returns all constants in declaration order.
func EvidenceRecordIncorporationTypeValues() []EvidenceRecordIncorporationType {
	return []EvidenceRecordIncorporationType{
		EvidenceRecordIncorporationTypeInternalEvidenceRecord,
		EvidenceRecordIncorporationTypeExternalEvidenceRecord,
	}
}

// EvidenceRecordIncorporationTypeValueOf returns the EvidenceRecordIncorporationType matching the given Java enum name.
func EvidenceRecordIncorporationTypeValueOf(name string) (EvidenceRecordIncorporationType, error) {
	for _, v := range EvidenceRecordIncorporationTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant EvidenceRecordIncorporationType.%s", name)
}
