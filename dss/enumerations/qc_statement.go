// Ported from dss-enumerations/.../QCStatement.java (DSS 6.5.RC1).
package enumerations

// QCStatement defines QCStatements based on ETSI EN 319 412-5. Implements
// OidDescription.
type QCStatement string

const (
	// QCStatement_QC_COMPLIANCE claims that the certificate is a EU
	// qualified certificate.
	QCStatement_QC_COMPLIANCE QCStatement = "QC_COMPLIANCE"
	// QCStatement_QC_LIMIT_VALUE regards limits on the value of
	// transactions.
	QCStatement_QC_LIMIT_VALUE QCStatement = "QC_LIMIT_VALUE"
	// QCStatement_QC_RETENTION_PERIOD indicates the duration of the
	// retention period of material information.
	QCStatement_QC_RETENTION_PERIOD QCStatement = "QC_RETENTION_PERIOD"
	// QCStatement_QC_SSCD claims that the private key related to the
	// certified public key resides in a QSCD.
	QCStatement_QC_SSCD QCStatement = "QC_SSCD"
	// QCStatement_QC_PDS regards location of PKI Disclosure Statements
	// (PDS).
	QCStatement_QC_PDS QCStatement = "QC_PDS"
	// QCStatement_QC_TYPE claims that the certificate is a EU qualified
	// certificate of a particular type.
	QCStatement_QC_TYPE QCStatement = "QC_TYPE"
	// QCStatement_QC_CCLEGISLATION is esi4-qcStatement-7.
	QCStatement_QC_CCLEGISLATION QCStatement = "QC_CCLEGISLATION"
	// QCStatement_QC_IDENT_METHOD is esi4-qcStatement-8.
	QCStatement_QC_IDENT_METHOD QCStatement = "QC_IDENT_METHOD"
	// QCStatement_QC_QSCD_LEGISLATION is id-etsi-qcs-QcQSCDlegislation.
	QCStatement_QC_QSCD_LEGISLATION QCStatement = "QC_QSCD_LEGISLATION"
	// QCStatement_QC_PSB is PSB certificate mandatory data
	// (esi4-qcStatement-10).
	QCStatement_QC_PSB QCStatement = "QC_PSB"
)

type qcStatementFields struct {
	description string
	oid         string
}

// qcStatementData holds the (description, oid) tuple for each constant.
var qcStatementData = map[QCStatement]qcStatementFields{
	QCStatement_QC_COMPLIANCE:       {"qc-compliance", "0.4.0.1862.1.1"},
	QCStatement_QC_LIMIT_VALUE:      {"qc-limit-value", "0.4.0.1862.1.2"},
	QCStatement_QC_RETENTION_PERIOD: {"qc-retention-period", "0.4.0.1862.1.3"},
	QCStatement_QC_SSCD:             {"qc-sscd", "0.4.0.1862.1.4"},
	QCStatement_QC_PDS:              {"qc-pds", "0.4.0.1862.1.5"},
	QCStatement_QC_TYPE:             {"qc-type", "0.4.0.1862.1.6"},
	QCStatement_QC_CCLEGISLATION:    {"qc-cclegislation", "0.4.0.1862.1.7"},
	QCStatement_QC_IDENT_METHOD:     {"qc-identMethod", "0.4.0.1862.1.8"},
	QCStatement_QC_QSCD_LEGISLATION: {"qc-qscdLegislation", "0.4.0.1862.1.9"},
	QCStatement_QC_PSB:              {"qc-psb", "0.4.0.194126.1.3"},
}

// QCStatementValues returns all constants in declaration order.
func QCStatementValues() []QCStatement {
	return []QCStatement{
		QCStatement_QC_COMPLIANCE,
		QCStatement_QC_LIMIT_VALUE,
		QCStatement_QC_RETENTION_PERIOD,
		QCStatement_QC_SSCD,
		QCStatement_QC_PDS,
		QCStatement_QC_TYPE,
		QCStatement_QC_CCLEGISLATION,
		QCStatement_QC_IDENT_METHOD,
		QCStatement_QC_QSCD_LEGISLATION,
		QCStatement_QC_PSB,
	}
}

// OID returns the object identifier. Implements OidBasedEnum.
func (q QCStatement) OID() string {
	return qcStatementData[q].oid
}

// Description returns the user-friendly identifier. Implements
// OidDescription.
func (q QCStatement) Description() string {
	return qcStatementData[q].description
}

// QCStatementForLabel gets a QCStatement for the given label description
// string. Returns "" (zero value) if description matches no known
// QCStatement, mirroring Java's null return.
func QCStatementForLabel(description string) QCStatement {
	for _, v := range QCStatementValues() {
		if qcStatementData[v].description == description {
			return v
		}
	}
	return ""
}

// QCStatementForOID gets a QCStatement for the given OID string. Returns ""
// (zero value) if oid matches no known QCStatement, mirroring Java's null
// return.
func QCStatementForOID(oid string) QCStatement {
	for _, v := range QCStatementValues() {
		if qcStatementData[v].oid == oid {
			return v
		}
	}
	return ""
}

// compile-time interface assertion.
var _ OidDescription = QCStatement("")
