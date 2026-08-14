// Ported from dss-enumerations/.../QCTypeEnum.java (DSS 6.5.RC1).
package enumerations

// QCTypeEnum defines QC type identifiers based on ETSI EN 319 412-5.
// Implements QCType.
type QCTypeEnum string

const (
	// QCTypeEnum_QCT_ESIGN is id-etsi-qct-esign OBJECT IDENTIFIER ::= {
	// id-etsi-qcs-QcType 1 } -- Certificate for electronic signatures as
	// defined in Regulation (EU) No 910/2014.
	QCTypeEnum_QCT_ESIGN QCTypeEnum = "QCT_ESIGN"
	// QCTypeEnum_QCT_ESEAL is id-etsi-qct-eseal OBJECT IDENTIFIER ::= {
	// id-etsi-qcs-QcType 2 } -- Certificate for electronic seals as
	// defined in Regulation (EU) No 910/2014.
	QCTypeEnum_QCT_ESEAL QCTypeEnum = "QCT_ESEAL"
	// QCTypeEnum_QCT_WEB is id-etsi-qct-web OBJECT IDENTIFIER ::= {
	// id-etsi-qcs-QcType 3 } -- Certificate for website authentication as
	// defined in Regulation (EU) No 910/2014.
	QCTypeEnum_QCT_WEB QCTypeEnum = "QCT_WEB"
	// QCTypeEnum_QCT_PID is id-etsi-qct-pid OBJECT IDENTIFIER ::= {
	// id-etsi-eidas2-qct-extensions 1 } -- Certificate for PID provider
	// sign/seal certificate. TS 119 412-6 QcType.
	QCTypeEnum_QCT_PID QCTypeEnum = "QCT_PID"
	// QCTypeEnum_QCT_WAL is id-etsi-qct-wal OBJECT IDENTIFIER ::= {
	// id-etsi-eidas2-qct-extensions 2 } -- Certificate for Wallet provider
	// sign/seal certificate. TS 119 412-6 QcType.
	QCTypeEnum_QCT_WAL QCTypeEnum = "QCT_WAL"
)

type qcTypeEnumFields struct {
	description string
	oid         string
}

// qcTypeEnumData holds the (description, oid) pair for each constant.
var qcTypeEnumData = map[QCTypeEnum]qcTypeEnumFields{
	QCTypeEnum_QCT_ESIGN: {"qc-type-esign", "0.4.0.1862.1.6.1"},
	QCTypeEnum_QCT_ESEAL: {"qc-type-eseal", "0.4.0.1862.1.6.2"},
	QCTypeEnum_QCT_WEB:   {"qc-type-web", "0.4.0.1862.1.6.3"},
	QCTypeEnum_QCT_PID:   {"qc-type-pid", "0.4.0.194126.1.1"},
	QCTypeEnum_QCT_WAL:   {"qc-type-wal", "0.4.0.194126.1.2"},
}

// QCTypeEnumValues returns all constants in declaration order.
func QCTypeEnumValues() []QCTypeEnum {
	return []QCTypeEnum{
		QCTypeEnum_QCT_ESIGN,
		QCTypeEnum_QCT_ESEAL,
		QCTypeEnum_QCT_WEB,
		QCTypeEnum_QCT_PID,
		QCTypeEnum_QCT_WAL,
	}
}

// OID returns the QCType OID. Implements QCType (via OidDescription).
func (q QCTypeEnum) OID() string {
	return qcTypeEnumData[q].oid
}

// Description returns the QCType description. Implements QCType (via
// OidDescription).
func (q QCTypeEnum) Description() string {
	return qcTypeEnumData[q].description
}

// QCTypeEnumForLabel gets a QCTypeEnum for the given label description
// string. Returns "" (the zero value) if not found, mirroring Java's null
// return.
func QCTypeEnumForLabel(description string) QCTypeEnum {
	for _, v := range QCTypeEnumValues() {
		if qcTypeEnumData[v].description == description {
			return v
		}
	}
	return ""
}
