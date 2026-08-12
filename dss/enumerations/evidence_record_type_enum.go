// Ported from dss-enumerations/.../EvidenceRecordTypeEnum.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// EvidenceRecordTypeEnum defines supported Evidence Record types.
type EvidenceRecordTypeEnum string

const (
	// EvidenceRecordTypeEnum_XML_EVIDENCE_RECORD is an XML Evidence Record
	// according to RFC 6283.
	EvidenceRecordTypeEnum_XML_EVIDENCE_RECORD EvidenceRecordTypeEnum = "XML_EVIDENCE_RECORD"
	// EvidenceRecordTypeEnum_ASN1_EVIDENCE_RECORD is an ASN.1 Evidence
	// Record according to RFC 4998.
	EvidenceRecordTypeEnum_ASN1_EVIDENCE_RECORD EvidenceRecordTypeEnum = "ASN1_EVIDENCE_RECORD"
)

// evidenceRecordTypeEnumLabels holds the user-friendly label for each
// constant.
var evidenceRecordTypeEnumLabels = map[EvidenceRecordTypeEnum]string{
	EvidenceRecordTypeEnum_XML_EVIDENCE_RECORD:  "XML Evidence Record",
	EvidenceRecordTypeEnum_ASN1_EVIDENCE_RECORD: "ASN.1 Evidence Record",
}

// EvidenceRecordTypeEnumValues returns all constants in declaration order.
func EvidenceRecordTypeEnumValues() []EvidenceRecordTypeEnum {
	return []EvidenceRecordTypeEnum{
		EvidenceRecordTypeEnum_XML_EVIDENCE_RECORD,
		EvidenceRecordTypeEnum_ASN1_EVIDENCE_RECORD,
	}
}

// Label gets a user-friendly descriptor of an evidence record type.
func (e EvidenceRecordTypeEnum) Label() string {
	return evidenceRecordTypeEnumLabels[e]
}

// EvidenceRecordTypeEnumFromLabel gets an EvidenceRecordTypeEnum for the
// given label string value. Returns an error if label matches no known
// evidence record type, mirroring Java's UnsupportedOperationException.
func EvidenceRecordTypeEnumFromLabel(label string) (EvidenceRecordTypeEnum, error) {
	for _, v := range EvidenceRecordTypeEnumValues() {
		if evidenceRecordTypeEnumLabels[v] == label {
			return v, nil
		}
	}
	return "", fmt.Errorf("evidence record of type '%s' is not supported", label)
}

// compile-time interface assertion: none (no marker interface).
