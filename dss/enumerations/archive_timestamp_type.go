// Ported from dss-enumerations/.../ArchiveTimestampType.java (DSS 6.5.RC1).
package enumerations

// ArchiveTimestampType represents the different types of archive timestamps.
type ArchiveTimestampType string

const (
	// ArchiveTimestampTypeXAdES is the default XAdES Archive Timestamp Type.
	ArchiveTimestampTypeXAdES ArchiveTimestampType = "XAdES"
	// ArchiveTimestampTypeXAdES141 is the ETSI TS 101 903 (XAdES 1.4.1)
	// ArchiveTimeStamp.
	ArchiveTimestampTypeXAdES141 ArchiveTimestampType = "XAdES_141"
	// ArchiveTimestampTypeCAdES is the default CAdES Archive Timestamp Type.
	ArchiveTimestampTypeCAdES ArchiveTimestampType = "CAdES"
	// ArchiveTimestampTypeCAdESV2 is the archive-time-stamp-v2.
	ArchiveTimestampTypeCAdESV2 ArchiveTimestampType = "CAdES_V2"
	// ArchiveTimestampTypeCAdESV3 is the archive-time-stamp-v3.
	ArchiveTimestampTypeCAdESV3 ArchiveTimestampType = "CAdES_V3"
	// ArchiveTimestampTypeCAdESDetached is the detached timestamp, used
	// for ASiC.
	ArchiveTimestampTypeCAdESDetached ArchiveTimestampType = "CAdES_DETACHED"
	// ArchiveTimestampTypeJAdES is the JAdES: arcTst.
	ArchiveTimestampTypeJAdES ArchiveTimestampType = "JAdES"
	// ArchiveTimestampTypeCBAdES is the CB-AdES: arcTst.
	ArchiveTimestampTypeCBAdES ArchiveTimestampType = "CB_AdES"
	// ArchiveTimestampTypePAdES is the DOCUMENT_TIMESTAMP covering a DSS
	// dictionary (revision).
	ArchiveTimestampTypePAdES ArchiveTimestampType = "PAdES"
	// ArchiveTimestampTypeXMLEvidenceRecord is the XML Evidence Record
	// time-stamp.
	ArchiveTimestampTypeXMLEvidenceRecord ArchiveTimestampType = "XML_EVIDENCE_RECORD"
	// ArchiveTimestampTypeASN1EvidenceRecord is the ASN.1 Evidence
	// Record time-stamp.
	ArchiveTimestampTypeASN1EvidenceRecord ArchiveTimestampType = "ASN1_EVIDENCE_RECORD"
)

// ArchiveTimestampTypeValues returns all constants in declaration order.
func ArchiveTimestampTypeValues() []ArchiveTimestampType {
	return []ArchiveTimestampType{
		ArchiveTimestampTypeXAdES,
		ArchiveTimestampTypeXAdES141,
		ArchiveTimestampTypeCAdES,
		ArchiveTimestampTypeCAdESV2,
		ArchiveTimestampTypeCAdESV3,
		ArchiveTimestampTypeCAdESDetached,
		ArchiveTimestampTypeJAdES,
		ArchiveTimestampTypeCBAdES,
		ArchiveTimestampTypePAdES,
		ArchiveTimestampTypeXMLEvidenceRecord,
		ArchiveTimestampTypeASN1EvidenceRecord,
	}
}
