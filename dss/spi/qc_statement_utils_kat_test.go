package spi

import (
	"fmt"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
	"golang.org/x/crypto/cryptobyte"
	cbasn1 "golang.org/x/crypto/cryptobyte/asn1"
)

// TestQcStatementUtilsKnownAnswers replays the QCStatements known answers produced with
// BouncyCastle; see certificate_extensions_utils_kat_test.go for how the fixture is generated.
func TestQcStatementUtilsKnownAnswers(t *testing.T) {
	for index, answers := range certificateExtensionsKATLoad(t) {
		name := answers.str(t, "file")
		if certificateExtensionsKATUnparseable[name] {
			continue
		}
		t.Run(fmt.Sprintf("%02d_%s", index, name), func(t *testing.T) {
			certificateToken := certificateExtensionsKATToken(t, name)
			qcStatements := QcStatementUtilsQcStatements(certificateToken)
			// CertificateExtensionsUtils#getQcStatements only adds a second checkCritical call.
			if viaExtensions := CertificateExtensionsUtilsQcStatements(certificateToken); (viaExtensions == nil) != (qcStatements == nil) {
				t.Fatalf("CertificateExtensionsUtilsQcStatements disagrees with QcStatementUtilsQcStatements")
			}
			if answers.isNull(t, "qc") {
				if qcStatements != nil {
					t.Fatalf("QcStatements: got a value, want nil")
				}
				return
			}
			if qcStatements == nil {
				t.Fatalf("QcStatements: got nil, want a value")
			}
			certificateExtensionsKATAssertHex(t, "qc.octets", qcStatements.Octets(), answers)
			if got, want := qcStatements.IsCritical(), answers.boolean(t, "qc.critical"); got != want {
				t.Errorf("QcStatements.IsCritical() = %v, want %v", got, want)
			}
			if got, want := qcStatements.IsQcCompliance(), answers.boolean(t, "qc.qcCompliance"); got != want {
				t.Errorf("QcStatements.IsQcCompliance() = %v, want %v", got, want)
			}
			if got, want := qcStatements.IsQcQSCD(), answers.boolean(t, "qc.qcQSCD"); got != want {
				t.Errorf("QcStatements.IsQcQSCD() = %v, want %v", got, want)
			}

			// QcLimitValue
			qcLimitValue := qcStatements.QcLimitValue()
			if answers.isNull(t, "qc.limit") {
				if qcLimitValue != nil {
					t.Errorf("QcLimitValue: got a value, want nil")
				}
			} else {
				if qcLimitValue == nil {
					t.Fatalf("QcLimitValue: got nil, want a value")
				}
				if got, want := qcLimitValue.Currency(), answers.orEmpty(t, "qc.limit.currency"); got != want {
					t.Errorf("QCLimitValue.Currency() = %q, want %q", got, want)
				}
				if got, want := qcLimitValue.Amount(), answers.integer(t, "qc.limit.amount"); got != want {
					t.Errorf("QCLimitValue.Amount() = %d, want %d", got, want)
				}
				if got, want := qcLimitValue.Exponent(), answers.integer(t, "qc.limit.exponent"); got != want {
					t.Errorf("QCLimitValue.Exponent() = %d, want %d", got, want)
				}
			}

			// QcEuRetentionPeriod
			retentionPeriod := qcStatements.QcEuRetentionPeriod()
			if answers.isNull(t, "qc.retentionPeriod") {
				if retentionPeriod != nil {
					t.Errorf("QcEuRetentionPeriod = %d, want nil", *retentionPeriod)
				}
			} else {
				if retentionPeriod == nil {
					t.Fatalf("QcEuRetentionPeriod: got nil, want a value")
				}
				if got, want := *retentionPeriod, answers.integer(t, "qc.retentionPeriod"); got != want {
					t.Errorf("QcEuRetentionPeriod = %d, want %d", got, want)
				}
			}

			// QcEuPDS
			pdsLocations := qcStatements.QcEuPDS()
			if got, want := len(pdsLocations), answers.integer(t, "qc.pds.count"); got != want {
				t.Fatalf("QcEuPDS has %d entries, want %d", got, want)
			}
			for i, pdsLocation := range pdsLocations {
				if got, want := pdsLocation.Url(), answers.orEmpty(t, fmt.Sprintf("qc.pds.%d.url", i)); got != want {
					t.Errorf("PdsLocation %d URL = %q, want %q", i, got, want)
				}
				if got, want := pdsLocation.Language(), answers.orEmpty(t, fmt.Sprintf("qc.pds.%d.language", i)); got != want {
					t.Errorf("PdsLocation %d language = %q, want %q", i, got, want)
				}
			}

			// QcTypes, identified by their OID
			qcTypeOids := make([]string, 0)
			for _, qcType := range qcStatements.QcTypes() {
				qcTypeOids = append(qcTypeOids, qcType.OID())
			}
			certificateExtensionsKATAssertStrings(t, "qc.types", qcTypeOids, answers)

			certificateExtensionsKATAssertStrings(t, "qc.ccLegislation", qcStatements.QcLegislationCountryCodes(), answers)
			certificateExtensionsKATAssertStrings(t, "qc.qscdLegislation", qcStatements.QcQSCDLegislationCountryCodes(), answers)
			certificateExtensionsKATAssertStrings(t, "qc.otherOids", qcStatements.OtherOids(), answers)

			// QcSemanticsIdentifier: the fixture records the OID, the model the enumeration constant.
			wantSemantics := enumerations.SemanticsIdentifier("")
			if !answers.isNull(t, "qc.semanticsIdentifier") {
				wantSemantics = enumerations.SemanticsIdentifierFromOID(answers.str(t, "qc.semanticsIdentifier"))
			}
			if got := qcStatements.QcSemanticsIdentifier(); got != wantSemantics {
				t.Errorf("QcSemanticsIdentifier = %q, want %q", got, wantSemantics)
			}

			// QcIdentMethod
			qcIdentMethod := qcStatements.QcIdentMethod()
			if answers.isNull(t, "qc.identMethod") {
				if qcIdentMethod != nil {
					t.Errorf("QcIdentMethod = %v, want nil", qcIdentMethod)
				}
			} else {
				if qcIdentMethod == nil {
					t.Fatalf("QcIdentMethod: got nil, want a value")
				}
				if got, want := qcIdentMethod.OID(), answers.str(t, "qc.identMethod"); got != want {
					t.Errorf("QcIdentMethod OID = %q, want %q", got, want)
				}
			}

			// PSD2QcType
			psd2QcType := qcStatements.Psd2QcType()
			if answers.isNull(t, "qc.psd2") {
				if psd2QcType != nil {
					t.Errorf("Psd2QcType: got a value, want nil")
				}
			} else {
				if psd2QcType == nil {
					t.Fatalf("Psd2QcType: got nil, want a value")
				}
				if got, want := psd2QcType.NcaName(), answers.orEmpty(t, "qc.psd2.ncaName"); got != want {
					t.Errorf("PSD2QcType.NcaName() = %q, want %q", got, want)
				}
				if got, want := psd2QcType.NcaId(), answers.orEmpty(t, "qc.psd2.ncaId"); got != want {
					t.Errorf("PSD2QcType.NcaId() = %q, want %q", got, want)
				}
				rolesOfPSP := psd2QcType.RolesOfPSP()
				if got, want := len(rolesOfPSP), answers.integer(t, "qc.psd2.roles.count"); got != want {
					t.Fatalf("PSD2QcType.RolesOfPSP() has %d entries, want %d", got, want)
				}
				for i, roleOfPSP := range rolesOfPSP {
					wantOid := enumerations.RoleOfPspOidFromOid(answers.str(t, fmt.Sprintf("qc.psd2.roles.%d.oid", i)))
					if got := roleOfPSP.PspOid(); got != wantOid {
						t.Errorf("RoleOfPSP %d OID = %q, want %q", i, got, wantOid)
					}
					if got, want := roleOfPSP.PspName(), answers.orEmpty(t, fmt.Sprintf("qc.psd2.roles.%d.name", i)); got != want {
						t.Errorf("RoleOfPSP %d name = %q, want %q", i, got, want)
					}
				}
			}

			// QCPSB
			qcPSB := qcStatements.QcPSB()
			if answers.isNull(t, "qc.psb") {
				if qcPSB != nil {
					t.Errorf("QcPSB: got a value, want nil")
				}
				return
			}
			if qcPSB == nil {
				t.Fatalf("QcPSB: got nil, want a value")
			}
			if got, want := qcPSB.CountryOfLegislation(), answers.orEmpty(t, "qc.psb.country"); got != want {
				t.Errorf("QCPSB.CountryOfLegislation() = %q, want %q", got, want)
			}
			if got, want := qcPSB.AuthSourceIdentification(), answers.orEmpty(t, "qc.psb.authSource"); got != want {
				t.Errorf("QCPSB.AuthSourceIdentification() = %q, want %q", got, want)
			}
			if got, want := qcPSB.LegislationIdentification(), answers.orEmpty(t, "qc.psb.legislation"); got != want {
				t.Errorf("QCPSB.LegislationIdentification() = %q, want %q", got, want)
			}
		})
	}
}

// TestQcStatementUtilsFromSequence covers the ASN1Sequence overload, which upstream exposes for the
// signed attributes carrying QCStatements outside a certificate.
func TestQcStatementUtilsFromSequence(t *testing.T) {
	if got := QcStatementUtilsQcStatementsFromSequence(nil); got != nil {
		t.Errorf("QcStatementsFromSequence(nil) = %v, want nil", got)
	}

	// SEQUENCE { SEQUENCE { OBJECT IDENTIFIER id-etsi-qcs-QcCompliance } }
	wellFormed := qcStatementUtilsTestSequence(t, func(builder *cryptobyte.Builder) {
		builder.AddASN1(cbasn1.SEQUENCE, func(statement *cryptobyte.Builder) {
			statement.AddASN1ObjectIdentifier([]int{0, 4, 0, 1862, 1, 1})
		})
	})
	qcStatements := QcStatementUtilsQcStatementsFromSequence(wellFormed)
	if qcStatements == nil {
		t.Fatalf("QcStatementsFromSequence: got nil, want a value")
	}
	if !qcStatements.IsQcCompliance() {
		t.Errorf("IsQcCompliance() = false, want true")
	}
	if got := utils.ToHex(qcStatements.Octets()); got != utils.ToHex(wellFormed) {
		t.Errorf("Octets() = %s, want %s", got, utils.ToHex(wellFormed))
	}

	// SEQUENCE { OBJECT IDENTIFIER } - the element is not a QCStatement and is skipped.
	bareOid := qcStatementUtilsTestSequence(t, func(builder *cryptobyte.Builder) {
		builder.AddASN1ObjectIdentifier([]int{0, 4, 0, 1862, 1, 1})
	})
	qcStatements = QcStatementUtilsQcStatementsFromSequence(bareOid)
	if qcStatements == nil {
		t.Fatalf("QcStatementsFromSequence: got nil, want a value")
	}
	if qcStatements.IsQcCompliance() {
		t.Errorf("IsQcCompliance() = true, want false")
	}
}

func qcStatementUtilsTestSequence(t *testing.T, content func(*cryptobyte.Builder)) []byte {
	t.Helper()
	builder := cryptobyte.NewBuilder(nil)
	builder.AddASN1(cbasn1.SEQUENCE, content)
	encoded, err := builder.Bytes()
	if err != nil {
		t.Fatalf("unable to build the test sequence: %v", err)
	}
	return encoded
}

// TestQcStatementUtilsPresence checks the lookup helpers against the certificate upstream's
// QcStatementsUtilsTest#certWithQcQSCDlegislationQcStatement and #cert1 use.
func TestQcStatementUtilsPresence(t *testing.T) {
	answers := certificateExtensionsKATLoad(t)
	var withQSCDLegislation, withPsd2 *model.CertificateToken
	for _, entry := range answers {
		name := entry.str(t, "file")
		if certificateExtensionsKATUnparseable[name] || entry.isNull(t, "qc") {
			continue
		}
		if entry.integer(t, "qc.qscdLegislation.count") > 0 && withQSCDLegislation == nil {
			withQSCDLegislation = certificateExtensionsKATToken(t, name)
		}
		if !entry.isNull(t, "qc.psd2") && withPsd2 == nil {
			withPsd2 = certificateExtensionsKATToken(t, name)
		}
	}
	if withQSCDLegislation == nil || withPsd2 == nil {
		t.Fatalf("the fixture corpus no longer covers QcQSCDlegislation and PSD2QcType")
	}

	qcStatements := QcStatementUtilsQcStatements(withQSCDLegislation)
	if !QcStatementUtilsIsQcQSCDlegislationPresent(qcStatements, "ZZ") {
		t.Errorf("IsQcQSCDlegislationPresent(ZZ) = false, want true")
	}
	if QcStatementUtilsIsQcQSCDlegislationPresent(qcStatements, "XX") {
		t.Errorf("IsQcQSCDlegislationPresent(XX) = true, want false")
	}
	if !QcStatementUtilsIsQcStatementPresent(qcStatements, enumerations.QCStatement_QC_QSCD_LEGISLATION.OID()) {
		t.Errorf("IsQcStatementPresent(QC_QSCD_LEGISLATION) = false, want true")
	}
	if !QcStatementUtilsIsQcStatementPresent(qcStatements, enumerations.QCStatement_QC_COMPLIANCE.OID()) {
		t.Errorf("IsQcStatementPresent(QC_COMPLIANCE) = false, want true")
	}
	if QcStatementUtilsIsQcStatementPresent(qcStatements, "1.2.3.4.5") {
		t.Errorf("IsQcStatementPresent(1.2.3.4.5) = true, want false")
	}

	qcStatements = QcStatementUtilsQcStatements(withPsd2)
	if !QcStatementUtilsIsQcStatementPresent(qcStatements, OID_psd2_qcStatement.String()) {
		t.Errorf("IsQcStatementPresent(psd2) = false, want true")
	}
	if !QcStatementUtilsIsQcTypePresent(qcStatements, enumerations.QCTypeEnum_QCT_WEB.OID()) {
		t.Errorf("IsQcTypePresent(qc-type-web) = false, want true")
	}
	if QcStatementUtilsIsQcTypePresent(qcStatements, enumerations.QCTypeEnum_QCT_ESEAL.OID()) {
		t.Errorf("IsQcTypePresent(qc-type-eseal) = true, want false")
	}
	if QcStatementUtilsIsQcLegislationPresent(qcStatements, "CZ") {
		t.Errorf("IsQcLegislationPresent(CZ) = true, want false")
	}
}

// TestQcStatementUtilsQcTypesForOIDs covers the public getQcTypes(List<String>) overload, blank
// OIDs included.
func TestQcStatementUtilsQcTypesForOIDs(t *testing.T) {
	qcTypes := QcStatementUtilsQcTypesForOIDs([]string{
		enumerations.QCTypeEnum_QCT_ESIGN.OID(), "", "   ", "1.2.3.4",
	})
	if len(qcTypes) != 2 {
		t.Fatalf("QcTypesForOIDs returned %d types, want 2", len(qcTypes))
	}
	if got := qcTypes[0].OID(); got != enumerations.QCTypeEnum_QCT_ESIGN.OID() {
		t.Errorf("qcTypes[0].OID() = %q, want %q", got, enumerations.QCTypeEnum_QCT_ESIGN.OID())
	}
	if got := qcTypes[1].OID(); got != "1.2.3.4" {
		t.Errorf("qcTypes[1].OID() = %q, want %q", got, "1.2.3.4")
	}
	if got := qcTypes[1].Description(); got != enumerations.QCType_UNKNOWN_TYPE {
		t.Errorf("qcTypes[1].Description() = %q, want %q", got, enumerations.QCType_UNKNOWN_TYPE)
	}
}
