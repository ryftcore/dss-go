// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/QcStatements.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// QcStatements contains the QcStatement information based on ETSI EN 319 412-1/5 and ETSI
// TS 119 495.
type QcStatements struct {
	CertificateExtension

	// qcCompliance is:
	//
	//	esi4-qcStatement-1 QC-STATEMENT ::= { IDENTIFIED BY id-etsi-qcs-QcCompliance }
	//	id-etsi-qcs-QcCompliance OBJECT IDENTIFIER ::= { id-etsi-qcs 1 }
	qcCompliance bool

	// qcLimitValue is:
	//
	//	esi4-qcStatement-2 QC-STATEMENT ::= { SYNTAX QcEuLimitValue IDENTIFIED BY
	//	id-etsi-qcs-QcLimitValue }
	//	id-etsi-qcs-QcLimitValue OBJECT IDENTIFIER ::= { id-etsi-qcs 2 }
	qcLimitValue *QCLimitValue

	// qcEuRetentionPeriod is:
	//
	//	esi4-qcStatement-3 QC-STATEMENT ::= { SYNTAX QcEuRetentionPeriod IDENTIFIED BY
	//	id-etsi-qcs-QcRetentionPeriod }
	//	id-etsi-qcs-QcRetentionPeriod OBJECT IDENTIFIER ::= { id-etsi-qcs 3 }
	qcEuRetentionPeriod *int

	// qcQSCD is:
	//
	//	esi4-qcStatement-4 QC-STATEMENT ::= { IDENTIFIED BY id-etsi-qcs-QcSSCD }
	//	id-etsi-qcs-QcSSCD OBJECT IDENTIFIER ::= { id-etsi-qcs 4 }
	qcQSCD bool

	// qcEuPDS is:
	//
	//	esi4-qcStatement-5 QC-STATEMENT ::= { SYNTAX QcEuPDS IDENTIFIED BY
	//	id-etsi-qcs-QcPDS }
	//	id-etsi-qcs-QcPDS OBJECT IDENTIFIER ::= { id-etsi-qcs 5 }
	qcEuPDS []*PdsLocation

	// qcTypes is:
	//
	//	esi4-qcStatement-6 QC-STATEMENT ::= { SYNTAX QcType IDENTIFIED BY
	//	id-etsi-qcs-QcType }
	//	id-etsi-qcs-QcType OBJECT IDENTIFIER ::= { id-etsi-qcs 6 }
	qcTypes []enumerations.QCType

	// qcLegislationCountryCodes is:
	//
	//	esi4-qcStatement-7 QC-STATEMENT ::= { SYNTAX QcCClegislation IDENTIFIED BY
	//	id-etsi-qcsQcCClegislation }
	//	id-etsi-qcs-QcCClegislation OBJECT IDENTIFIER ::= { id-etsi-qcs 7 }
	//
	//	QcCClegislation ::= SEQUENCE OF CountryName
	//	CountryName ::= PrintableString (SIZE (2)) (CONSTRAINED BY { -- ISO 3166-1 [6]
	//	alpha-2 codes only -- })
	qcLegislationCountryCodes []string

	// qcSemanticsIdentifier is:
	//
	//	id-etsi-qcs-semantics-identifiers OBJECT IDENTIFIER ::= { itu-t(0)
	//	identified-organization(4) etsi(0) id-cert-profile(194121) 1 }
	qcSemanticsIdentifier enumerations.SemanticsIdentifier

	// psd2QcType is:
	//
	//	etsi-psd2-qcStatement QC-STATEMENT ::= {SYNTAX PSD2QcType IDENTIFIED BY
	//	id-etsi-psd2-qcStatement }
	//	id-etsi-psd2-qcStatement OBJECT IDENTIFIER ::= { itu-t(0)
	//	identified-organization(4) etsi(0) psd2(19495) qcstatement(2) }
	psd2QcType *PSD2QcType

	// qcIdentMethod is:
	//
	//	esi4-qcStatement-8 QC-STATEMENT ::= { SYNTAX QcIdentMethod IDENTIFIED
	//	BY id-etsi-qcs-QcIdentMethod }
	//	QcIdentMethod ::= SEQUENCE SIZE (1) OF OBJECT IDENTIFIER ( id-etsi-qct-eIDAS1-ab |
	//	id-etsi-qct-eIDAS1-cd | id-etsi-qct-eIDAS2-acd | id-etsi-qct-eIDAS2-b | ...)
	//	id-etsi-qcs-QcIdentMethod OBJECT IDENTIFIER ::= { id-etsi-qcs 8 }
	qcIdentMethod enumerations.QCIdentMethod

	// qcQSCDLegislationCountryCodes is:
	//
	//	esi4-qcStatement-9 QC-STATEMENT ::= { SYNTAX QcQSCDlegislation IDENTIFIED BY
	//	id-etsi-qcs-QcQCSDlegislation }
	//
	//	id-etsi-qcs-QcQSCDlegislation OBJECT IDENTIFIER ::= { id-etsi-qcs 9 }
	//	QcQSCDlegislation ::= SEQUENCE SIZE (1..MAX) OF CountryName
	//	CountryName ::= PrintableString (SIZE (2)) (CONSTRAINED BY { -- ISO 3166-1 [6]
	//	alpha-2 codes only -- })
	qcQSCDLegislationCountryCodes []string

	// qcPSB is:
	//
	//	 -- PSB certificate mandatory data
	//	esi4-qcStatement-10 QC-STATEMENT ::= { SYNTAX QcPSB IDENTIFIED
	//	BY id-etsi-qcs-QcPSB }
	qcPSB *QCPSB

	// otherOids contains OIDs defined in QcStatements which are not supported by the
	// current implementation.
	otherOids []string
}

// NewQcStatements instantiates the object with null values. Ports the default
// constructor (which initializes otherOids to an empty list).
func NewQcStatements() *QcStatements {
	return &QcStatements{
		// Java's no-arg constructor calls super(CertificateExtensionEnum.X.getOid()), the
		// OID-only CertificateExtension(String) constructor - NOT
		// CertificateExtension(CertificateExtensionEnum). The description therefore stays
		// null, and the diagnostic-data builder emits no description attribute for it.
		CertificateExtension: NewCertificateExtension(enumerations.CertificateExtensionEnum_QC_STATEMENTS.OID()),
		otherOids:            []string{},
	}
}

// IsQcCompliance gets if the certificate is qc-compliant: TRUE if the qc-compliance
// extension is present, FALSE otherwise.
func (q *QcStatements) IsQcCompliance() bool {
	return q.qcCompliance
}

// SetQcCompliance sets if the qc-compliance extension is present.
func (q *QcStatements) SetQcCompliance(qcCompliance bool) {
	q.qcCompliance = qcCompliance
}

// QcLimitValue gets the qc-limit-value.
func (q *QcStatements) QcLimitValue() *QCLimitValue {
	return q.qcLimitValue
}

// SetQcLimitValue sets the qc-limit-value.
func (q *QcStatements) SetQcLimitValue(qcLimitValue *QCLimitValue) {
	q.qcLimitValue = qcLimitValue
}

// QcEuRetentionPeriod gets the qc-eu-retention-period.
func (q *QcStatements) QcEuRetentionPeriod() *int {
	return q.qcEuRetentionPeriod
}

// SetQcEuRetentionPeriod sets the qc-eu-retention-period.
func (q *QcStatements) SetQcEuRetentionPeriod(qcEuRetentionPeriod *int) {
	q.qcEuRetentionPeriod = qcEuRetentionPeriod
}

// IsQcQSCD gets if the certificate is qc-qscd: TRUE if the qc-qscd extension is present,
// FALSE otherwise.
func (q *QcStatements) IsQcQSCD() bool {
	return q.qcQSCD
}

// SetQcQSCD sets if the certificate is qc-qscd.
func (q *QcStatements) SetQcQSCD(qcQSCD bool) {
	q.qcQSCD = qcQSCD
}

// QcEuPDS gets the qc-eu-pds.
func (q *QcStatements) QcEuPDS() []*PdsLocation {
	return q.qcEuPDS
}

// SetQcEuPDS sets the qc-eu-pds.
func (q *QcStatements) SetQcEuPDS(qcEuPDS []*PdsLocation) {
	q.qcEuPDS = qcEuPDS
}

// QcTypes gets a list of QCTypes.
func (q *QcStatements) QcTypes() []enumerations.QCType {
	return q.qcTypes
}

// SetQcTypes sets a list of QCTypes.
func (q *QcStatements) SetQcTypes(qcTypes []enumerations.QCType) {
	q.qcTypes = qcTypes
}

// QcLegislationCountryCodes gets the qc-legislation-country-codes.
func (q *QcStatements) QcLegislationCountryCodes() []string {
	return q.qcLegislationCountryCodes
}

// SetQcLegislationCountryCodes sets the qc-legislation-country-codes.
func (q *QcStatements) SetQcLegislationCountryCodes(qcLegislationCountryCodes []string) {
	q.qcLegislationCountryCodes = qcLegislationCountryCodes
}

// QcSemanticsIdentifier gets the qc-semantics-identifier.
func (q *QcStatements) QcSemanticsIdentifier() enumerations.SemanticsIdentifier {
	return q.qcSemanticsIdentifier
}

// SetQcSemanticsIdentifier sets the qc-semantics-identifier.
func (q *QcStatements) SetQcSemanticsIdentifier(qcSemanticsIdentifier enumerations.SemanticsIdentifier) {
	q.qcSemanticsIdentifier = qcSemanticsIdentifier
}

// Psd2QcType gets the psd2-qc-type.
func (q *QcStatements) Psd2QcType() *PSD2QcType {
	return q.psd2QcType
}

// SetPsd2QcType sets the psd2-qc-type.
func (q *QcStatements) SetPsd2QcType(psd2QcType *PSD2QcType) {
	q.psd2QcType = psd2QcType
}

// QcIdentMethod gets the qc-identification-method.
func (q *QcStatements) QcIdentMethod() enumerations.QCIdentMethod {
	return q.qcIdentMethod
}

// SetQcIdentMethod sets the qc-identification-method.
func (q *QcStatements) SetQcIdentMethod(qcIdentMethod enumerations.QCIdentMethod) {
	q.qcIdentMethod = qcIdentMethod
}

// QcQSCDLegislationCountryCodes gets the qc-qscd-legislation-country-codes.
func (q *QcStatements) QcQSCDLegislationCountryCodes() []string {
	return q.qcQSCDLegislationCountryCodes
}

// SetQcQSCDLegislationCountryCodes sets the qc-qscd-legislation-country-codes.
func (q *QcStatements) SetQcQSCDLegislationCountryCodes(qcQSCDLegislationCountryCodes []string) {
	q.qcQSCDLegislationCountryCodes = qcQSCDLegislationCountryCodes
}

// QcPSB gets information for a Public Sector Body's Electronic Attestation of Attributes
// (PSBEAA) provider certificate.
func (q *QcStatements) QcPSB() *QCPSB {
	return q.qcPSB
}

// SetQcPSB sets information for a Public Sector Body's Electronic Attestation of
// Attributes (PSBEAA) provider certificate.
func (q *QcStatements) SetQcPSB(qcPSB *QCPSB) {
	q.qcPSB = qcPSB
}

// OtherOids returns a list of found OIDs not supported by the current implementation.
func (q *QcStatements) OtherOids() []string {
	return q.otherOids
}

// AddOtherOid adds a found OID not supported by the implementation.
func (q *QcStatements) AddOtherOid(oid string) {
	q.otherOids = append(q.otherOids, oid)
}
