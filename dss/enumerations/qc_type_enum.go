// Ported from dss-enumerations/.../QCTypeEnum.java (DSS 6.5.RC1).
package enumerations

// QCTypeEnum defines QC type identifiers based on ETSI EN 319 412-5.
// Implements QCType.
type QCTypeEnum string

const (
	// QCTypeEnumQCTESign is id-etsi-qct-esign OBJECT IDENTIFIER ::= {
	// id-etsi-qcs-QcType 1 } -- Certificate for electronic signatures as
	// defined in Regulation (EU) No 910/2014.
	QCTypeEnumQCTESign QCTypeEnum = "QCT_ESIGN"
	// QCTypeEnumQCTESeal is id-etsi-qct-eseal OBJECT IDENTIFIER ::= {
	// id-etsi-qcs-QcType 2 } -- Certificate for electronic seals as
	// defined in Regulation (EU) No 910/2014.
	QCTypeEnumQCTESeal QCTypeEnum = "QCT_ESEAL"
	// QCTypeEnumQCTWeb is id-etsi-qct-web OBJECT IDENTIFIER ::= {
	// id-etsi-qcs-QcType 3 } -- Certificate for website authentication as
	// defined in Regulation (EU) No 910/2014.
	QCTypeEnumQCTWeb QCTypeEnum = "QCT_WEB"
	// QCTypeEnumQCTPID is id-etsi-qct-pid OBJECT IDENTIFIER ::= {
	// id-etsi-eidas2-qct-extensions 1 } -- Certificate for PID provider
	// sign/seal certificate. TS 119 412-6 QcType.
	QCTypeEnumQCTPID QCTypeEnum = "QCT_PID"
	// QCTypeEnumQCTWAL is id-etsi-qct-wal OBJECT IDENTIFIER ::= {
	// id-etsi-eidas2-qct-extensions 2 } -- Certificate for Wallet provider
	// sign/seal certificate. TS 119 412-6 QcType.
	QCTypeEnumQCTWAL QCTypeEnum = "QCT_WAL"
)

type qcTypeEnumFields struct {
	description string
	oid         string
}

// qcTypeEnumData holds the (description, oid) pair for each constant.
var qcTypeEnumData = map[QCTypeEnum]qcTypeEnumFields{
	QCTypeEnumQCTESign: {"qc-type-esign", "0.4.0.1862.1.6.1"},
	QCTypeEnumQCTESeal: {"qc-type-eseal", "0.4.0.1862.1.6.2"},
	QCTypeEnumQCTWeb:   {"qc-type-web", "0.4.0.1862.1.6.3"},
	QCTypeEnumQCTPID:   {"qc-type-pid", "0.4.0.194126.1.1"},
	QCTypeEnumQCTWAL:   {"qc-type-wal", "0.4.0.194126.1.2"},
}

// QCTypeEnumValues returns all constants in declaration order.
func QCTypeEnumValues() []QCTypeEnum {
	return []QCTypeEnum{
		QCTypeEnumQCTESign,
		QCTypeEnumQCTESeal,
		QCTypeEnumQCTWeb,
		QCTypeEnumQCTPID,
		QCTypeEnumQCTWAL,
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
