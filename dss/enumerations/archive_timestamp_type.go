// Ported from dss-enumerations/.../ArchiveTimestampType.java (DSS 6.5.RC1).
package enumerations

// ArchiveTimestampType represents the different types of archive timestamps.
type ArchiveTimestampType string

const (
	// ArchiveTimestampType_XAdES is the default XAdES Archive Timestamp Type.
	ArchiveTimestampType_XAdES ArchiveTimestampType = "XAdES"
	// ArchiveTimestampType_XAdES_141 is the ETSI TS 101 903 (XAdES 1.4.1)
	// ArchiveTimeStamp.
	ArchiveTimestampType_XAdES_141 ArchiveTimestampType = "XAdES_141"
	// ArchiveTimestampType_CAdES is the default CAdES Archive Timestamp Type.
	ArchiveTimestampType_CAdES ArchiveTimestampType = "CAdES"
	// ArchiveTimestampType_CAdES_V2 is the archive-time-stamp-v2.
	ArchiveTimestampType_CAdES_V2 ArchiveTimestampType = "CAdES_V2"
	// ArchiveTimestampType_CAdES_V3 is the archive-time-stamp-v3.
	ArchiveTimestampType_CAdES_V3 ArchiveTimestampType = "CAdES_V3"
	// ArchiveTimestampType_CAdES_DETACHED is the detached timestamp, used
	// for ASiC.
	ArchiveTimestampType_CAdES_DETACHED ArchiveTimestampType = "CAdES_DETACHED"
	// ArchiveTimestampType_JAdES is the JAdES: arcTst.
	ArchiveTimestampType_JAdES ArchiveTimestampType = "JAdES"
	// ArchiveTimestampType_CB_AdES is the CB-AdES: arcTst.
	ArchiveTimestampType_CB_AdES ArchiveTimestampType = "CB_AdES"
	// ArchiveTimestampType_PAdES is the DOCUMENT_TIMESTAMP covering a DSS
	// dictionary (revision).
	ArchiveTimestampType_PAdES ArchiveTimestampType = "PAdES"
	// ArchiveTimestampType_XML_EVIDENCE_RECORD is the XML Evidence Record
	// time-stamp.
	ArchiveTimestampType_XML_EVIDENCE_RECORD ArchiveTimestampType = "XML_EVIDENCE_RECORD"
	// ArchiveTimestampType_ASN1_EVIDENCE_RECORD is the ASN.1 Evidence
	// Record time-stamp.
	ArchiveTimestampType_ASN1_EVIDENCE_RECORD ArchiveTimestampType = "ASN1_EVIDENCE_RECORD"
)

// ArchiveTimestampTypeValues returns all constants in declaration order.
func ArchiveTimestampTypeValues() []ArchiveTimestampType {
	return []ArchiveTimestampType{
		ArchiveTimestampType_XAdES,
		ArchiveTimestampType_XAdES_141,
		ArchiveTimestampType_CAdES,
		ArchiveTimestampType_CAdES_V2,
		ArchiveTimestampType_CAdES_V3,
		ArchiveTimestampType_CAdES_DETACHED,
		ArchiveTimestampType_JAdES,
		ArchiveTimestampType_CB_AdES,
		ArchiveTimestampType_PAdES,
		ArchiveTimestampType_XML_EVIDENCE_RECORD,
		ArchiveTimestampType_ASN1_EVIDENCE_RECORD,
	}
}
