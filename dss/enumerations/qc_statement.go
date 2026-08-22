// Ported from dss-enumerations/.../QCStatement.java (DSS 6.5.RC1).
package enumerations

// QCStatement defines QCStatements based on ETSI EN 319 412-5. Implements
// OidDescription.
type QCStatement string

const (
	// QCStatementQCCompliance claims that the certificate is a EU
	// qualified certificate.
	QCStatementQCCompliance QCStatement = "QC_COMPLIANCE"
	// QCStatementQCLimitValue regards limits on the value of
	// transactions.
	QCStatementQCLimitValue QCStatement = "QC_LIMIT_VALUE"
	// QCStatementQCRetentionPeriod indicates the duration of the
	// retention period of material information.
	QCStatementQCRetentionPeriod QCStatement = "QC_RETENTION_PERIOD"
	// QCStatementQCSSCD claims that the private key related to the
	// certified public key resides in a QSCD.
	QCStatementQCSSCD QCStatement = "QC_SSCD"
	// QCStatementQCPDS regards location of PKI Disclosure Statements
	// (PDS).
	QCStatementQCPDS QCStatement = "QC_PDS"
	// QCStatementQCType claims that the certificate is a EU qualified
	// certificate of a particular type.
	QCStatementQCType QCStatement = "QC_TYPE"
	// QCStatementQCCCLegislation is esi4-qcStatement-7.
	QCStatementQCCCLegislation QCStatement = "QC_CCLEGISLATION"
	// QCStatementQCIdentMethod is esi4-qcStatement-8.
	QCStatementQCIdentMethod QCStatement = "QC_IDENT_METHOD"
	// QCStatementQCQSCDLegislation is id-etsi-qcs-QcQSCDlegislation.
	QCStatementQCQSCDLegislation QCStatement = "QC_QSCD_LEGISLATION"
	// QCStatementQCPSB is PSB certificate mandatory data
	// (esi4-qcStatement-10).
	QCStatementQCPSB QCStatement = "QC_PSB"
)

type qcStatementFields struct {
	description string
	oid         string
}

// qcStatementData holds the (description, oid) tuple for each constant.
var qcStatementData = map[QCStatement]qcStatementFields{
	QCStatementQCCompliance:      {"qc-compliance", "0.4.0.1862.1.1"},
	QCStatementQCLimitValue:      {"qc-limit-value", "0.4.0.1862.1.2"},
	QCStatementQCRetentionPeriod: {"qc-retention-period", "0.4.0.1862.1.3"},
	QCStatementQCSSCD:            {"qc-sscd", "0.4.0.1862.1.4"},
	QCStatementQCPDS:             {"qc-pds", "0.4.0.1862.1.5"},
	QCStatementQCType:            {"qc-type", "0.4.0.1862.1.6"},
	QCStatementQCCCLegislation:   {"qc-cclegislation", "0.4.0.1862.1.7"},
	QCStatementQCIdentMethod:     {"qc-identMethod", "0.4.0.1862.1.8"},
	QCStatementQCQSCDLegislation: {"qc-qscdLegislation", "0.4.0.1862.1.9"},
	QCStatementQCPSB:             {"qc-psb", "0.4.0.194126.1.3"},
}

// QCStatementValues returns all constants in declaration order.
func QCStatementValues() []QCStatement {
	return []QCStatement{
		QCStatementQCCompliance,
		QCStatementQCLimitValue,
		QCStatementQCRetentionPeriod,
		QCStatementQCSSCD,
		QCStatementQCPDS,
		QCStatementQCType,
		QCStatementQCCCLegislation,
		QCStatementQCIdentMethod,
		QCStatementQCQSCDLegislation,
		QCStatementQCPSB,
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
