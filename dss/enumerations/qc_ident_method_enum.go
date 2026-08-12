// Ported from dss-enumerations/.../QCIdentMethodEnum.java (DSS 6.5.RC1).
package enumerations

// QCIdentMethodEnum contains QcQSCDlegislation statements defined within
// the ETSI EN 319 412-5 "4.3.5 QCStatement stating the used eIDAS/eIDAS2
// Article 24. identification method".
type QCIdentMethodEnum string

const (
	// QCIdentMethodEnum_QCT_EIDAS2_ACD is:
	//
	//	id-etsi-qct-eIDAS2-acd OBJECT IDENTIFIER ::= { id-etsi-qcs-QcIdentMethod 3 }
	//
	// -- Identification according to eIDAS2 Article 24. paragraph 1a a) or
	// c) or d).
	QCIdentMethodEnum_QCT_EIDAS2_ACD QCIdentMethodEnum = "QCT_EIDAS2_ACD"
	// QCIdentMethodEnum_QCT_EIDAS2_B is:
	//
	//	id-etsi-qct-eIDAS2-b OBJECT IDENTIFIER ::= { id-etsi-qcs-QcIdentMethod 4 }
	//
	// -- Identification according to eIDAS2 Article 24. paragraph 1a b).
	QCIdentMethodEnum_QCT_EIDAS2_B QCIdentMethodEnum = "QCT_EIDAS2_B"
)

type qcIdentMethodEnumFields struct {
	description string
	oid         string
}

// qcIdentMethodEnumData holds the (description, oid) tuple for each
// constant. NOTE: QCT_EIDAS2_B's description is copied verbatim from the
// upstream Java source, which reuses QCT_EIDAS2_ACD's description string
// ("qc-ident-method-eIDAS2-acd") rather than a "-b" variant; this looks
// like an upstream copy-paste bug but is preserved per the porting brief's
// "copy verbatim" rule.
var qcIdentMethodEnumData = map[QCIdentMethodEnum]qcIdentMethodEnumFields{
	QCIdentMethodEnum_QCT_EIDAS2_ACD: {"qc-ident-method-eIDAS2-acd", "0.4.0.1862.1.8.3"},
	QCIdentMethodEnum_QCT_EIDAS2_B:   {"qc-ident-method-eIDAS2-acd", "0.4.0.1862.1.8.4"},
}

// QCIdentMethodEnumValues returns all constants in declaration order.
func QCIdentMethodEnumValues() []QCIdentMethodEnum {
	return []QCIdentMethodEnum{
		QCIdentMethodEnum_QCT_EIDAS2_ACD,
		QCIdentMethodEnum_QCT_EIDAS2_B,
	}
}

// OID returns the QCIdentMethod OID.
func (q QCIdentMethodEnum) OID() string {
	return qcIdentMethodEnumData[q].oid
}

// Description returns the QCIdentMethod description.
func (q QCIdentMethodEnum) Description() string {
	return qcIdentMethodEnumData[q].description
}

// QCIdentMethodEnumForLabel gets a QCIdentMethodEnum for the given label
// description string. Panics if description is empty, mirroring Java's
// Objects.requireNonNull check.
//
// NOTE: because QCT_EIDAS2_ACD and QCT_EIDAS2_B share the same description
// string upstream, this lookup always resolves to QCT_EIDAS2_ACD (the first
// match in declaration order) for that shared description — matching
// Java's values() iteration returning the first match.
func QCIdentMethodEnumForLabel(description string) QCIdentMethodEnum {
	if description == "" {
		panic("Description label cannot be null!")
	}
	for _, v := range QCIdentMethodEnumValues() {
		if qcIdentMethodEnumData[v].description == description {
			return v
		}
	}
	return ""
}

// Compile-time assertion that QCIdentMethodEnum implements QCIdentMethod.
var _ QCIdentMethod = QCIdentMethodEnum_QCT_EIDAS2_ACD
