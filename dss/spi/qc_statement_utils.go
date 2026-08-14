// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/QcStatementUtils.java (DSS 6.5.RC1).
//
// Naming: the Java class holds static methods only and its getQcStatements() collides with
// CertificateExtensionsUtils#getQcStatements() once eu.europa.esig.dss.spi and spi.x509 flatten
// into the single Go package spi, so every exported function carries the owning Java class name
// as a prefix.
//
// BouncyCastle replacement: the QCStatements structures (RFC 3739, ETSI EN 319 412-5,
// ETSI TS 119 495) are decoded with cryptobyte; the helpers at the bottom of the file reproduce
// BouncyCastle's QCStatement, MonetaryValue and SemanticsInformation decoding. Attribute values go
// through DSSASN1UtilsString, the port of DSSASN1Utils#getString that upstream itself calls.
package spi

import (
	encasn1 "encoding/asn1"
	"math/big"

	"golang.org/x/crypto/cryptobyte"
	cbasn1 "golang.org/x/crypto/cryptobyte/asn1"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/extension"
	"github.com/utain/esig/dss/utils"
)

// qcStatementUtilsOIDIdQcsPkixQCSyntaxV2 is RFC3739QCObjectIdentifiers.id_qcs_pkixQCSyntax_v2.
// It has no counterpart in enumerations.QCStatement, so the literal is kept here.
const qcStatementUtilsOIDIdQcsPkixQCSyntaxV2 = "1.3.6.1.5.5.7.11.2"

// QcStatementUtilsQcStatements extracts the QCStatements from a certificate token, or nil when the
// extension is absent or cannot be read. Port of getQcStatements(CertificateToken).
func QcStatementUtilsQcStatements(certToken *model.CertificateToken) *extension.QcStatements {
	oid := enumerations.CertificateExtensionEnum_QC_STATEMENTS.OID()
	qcStatementsExtension := qcStatementUtilsExtensionContent(certToken, oid)
	if len(qcStatementsExtension) == 0 {
		return nil
	}
	qcStatements := QcStatementUtilsQcStatementsFromSequence(qcStatementsExtension)
	if qcStatements == nil {
		return nil
	}
	qcStatements.CheckCritical(certToken)
	return qcStatements
}

// QcStatementUtilsQcStatementsFromSequence extracts the QCStatements from the DER encoding of a
// QCStatements SEQUENCE. Port of getQcStatements(ASN1Sequence); a nil or empty encoding stands for
// the null sequence and yields nil.
//
// NOTE: does not check whether the extension is critical. Use QcStatementUtilsQcStatements if you
// need to know whether the certificate extension is critical.
//
// The octets recorded on the returned QcStatements are the SEQUENCE itself, NOT the OCTET STRING
// wrapping it as for every other certificate extension; that asymmetry is upstream's. Java
// re-encodes the parsed sequence with DSSASN1Utils#getDEREncoded, which for the DER input handled
// here returns the very same bytes, so the argument is stored verbatim.
func QcStatementUtilsQcStatementsFromSequence(qcStatementsSeq []byte) *extension.QcStatements {
	if len(qcStatementsSeq) == 0 {
		return nil
	}
	statements, ok := qcStatementUtilsSequenceElements(qcStatementsSeq)
	if !ok {
		return nil
	}

	result := extension.NewQcStatements()
	result.SetOctets(qcStatementsSeq)
	for _, element := range statements {
		oid, statementInfo, ok := qcStatementUtilsQCStatement(element)
		if !ok {
			continue
		}
		switch {
		case QcStatementUtilsIsQcCompliance(oid):
			result.SetQcCompliance(true)
		case QcStatementUtilsIsQcLimitValue(oid):
			result.SetQcLimitValue(qcStatementUtilsQcLimitValue(statementInfo))
		case QcStatementUtilsIsQcRetentionPeriod(oid):
			result.SetQcEuRetentionPeriod(qcStatementUtilsQcEuRetentionPeriod(statementInfo))
		case QcStatementUtilsIsQcSSCD(oid):
			result.SetQcQSCD(true)
		case QcStatementUtilsIsQcPds(oid):
			result.SetQcEuPDS(qcStatementUtilsQcEuPDS(statementInfo))
		case QcStatementUtilsIsQcType(oid):
			result.SetQcTypes(qcStatementUtilsQcTypes(statementInfo))
		case QcStatementUtilsIsQcCClegislation(oid):
			result.SetQcLegislationCountryCodes(qcStatementUtilsQcLegislationCountryCodes(statementInfo))
		case QcStatementUtilsIsQcSemanticsIdentifier(oid):
			result.SetQcSemanticsIdentifier(qcStatementUtilsQcSemanticsIdentifier(statementInfo))
		case QcStatementUtilsIsPsd2QcType(oid):
			result.SetPsd2QcType(qcStatementUtilsPsd2QcType(statementInfo))
		case QcStatementUtilsIsQcQSCDlegislation(oid):
			result.SetQcQSCDLegislationCountryCodes(qcStatementUtilsQcLegislationCountryCodes(statementInfo))
		case QcStatementUtilsIsQcIdentMethod(oid):
			result.SetQcIdentMethod(qcStatementUtilsQcIdentMethod(statementInfo))
		case QcStatementUtilsIsCertForPSB(oid):
			result.SetQcPSB(qcStatementUtilsQcPSB(statementInfo))
		default:
			result.AddOtherOid(oid)
		}
	}
	return result
}

// QcStatementUtilsIsQcCompliance reports whether the given OID is a QcCompliance statement.
// Port of isQcCompliance(String).
func QcStatementUtilsIsQcCompliance(oid string) bool {
	return enumerations.QCStatement_QC_COMPLIANCE.OID() == oid
}

// QcStatementUtilsIsQcLimitValue reports whether the given OID is a QcLimitValue statement.
// Port of isQcLimitValue(String).
func QcStatementUtilsIsQcLimitValue(oid string) bool {
	return enumerations.QCStatement_QC_LIMIT_VALUE.OID() == oid
}

// QcStatementUtilsIsQcRetentionPeriod reports whether the given OID is a QcRetentionPeriod
// statement. Port of isQcRetentionPeriod(String).
func QcStatementUtilsIsQcRetentionPeriod(oid string) bool {
	return enumerations.QCStatement_QC_RETENTION_PERIOD.OID() == oid
}

// QcStatementUtilsIsQcSSCD reports whether the given OID is a QcSSCD statement.
// Port of isQcSSCD(String).
func QcStatementUtilsIsQcSSCD(oid string) bool {
	return enumerations.QCStatement_QC_SSCD.OID() == oid
}

// QcStatementUtilsIsQcPds reports whether the given OID is a QcPds statement.
// Port of isQcPds(String).
func QcStatementUtilsIsQcPds(oid string) bool {
	return enumerations.QCStatement_QC_PDS.OID() == oid
}

// QcStatementUtilsIsQcType reports whether the given OID is a QcType statement.
// Port of isQcType(String).
func QcStatementUtilsIsQcType(oid string) bool {
	return enumerations.QCStatement_QC_TYPE.OID() == oid
}

// QcStatementUtilsIsQcCClegislation reports whether the given OID is a QcCClegislation statement.
// Port of isQcCClegislation(String).
func QcStatementUtilsIsQcCClegislation(oid string) bool {
	return enumerations.QCStatement_QC_CCLEGISLATION.OID() == oid
}

// QcStatementUtilsIsQcSemanticsIdentifier reports whether the given OID is a QcSemanticsIdentifier
// statement. Port of isQcSemanticsIdentifier(String).
func QcStatementUtilsIsQcSemanticsIdentifier(oid string) bool {
	return qcStatementUtilsOIDIdQcsPkixQCSyntaxV2 == oid
}

// QcStatementUtilsIsPsd2QcType reports whether the given OID is a Psd2QcType statement.
// Port of isPsd2QcType(String).
func QcStatementUtilsIsPsd2QcType(oid string) bool {
	return OID_psd2_qcStatement.String() == oid
}

// QcStatementUtilsIsQcIdentMethod reports whether the given OID is a QcIdentMethod statement.
// Port of isQcIdentMethod(String).
func QcStatementUtilsIsQcIdentMethod(oid string) bool {
	return enumerations.QCStatement_QC_IDENT_METHOD.OID() == oid
}

// QcStatementUtilsIsQcQSCDlegislation reports whether the given OID is a QcQSCDlegislation
// statement. Port of isQcQSCDlegislation(String).
func QcStatementUtilsIsQcQSCDlegislation(oid string) bool {
	return enumerations.QCStatement_QC_QSCD_LEGISLATION.OID() == oid
}

// QcStatementUtilsIsCertForPSB reports whether the given OID identifies a Public Sector Body's
// Electronic Attestation of Attributes (PSBEAA) provider certificate. Port of isCertForPSB(String).
func QcStatementUtilsIsCertForPSB(oid string) bool {
	return enumerations.QCStatement_QC_PSB.OID() == oid
}

// qcStatementUtilsQcLimitValue ports the private getQcLimitValue(ASN1Encodable) over
// BouncyCastle's MonetaryValue.
//
//	MonetaryValue ::= SEQUENCE { currency Iso4217CurrencyCode, amount INTEGER, exponent INTEGER }
//	Iso4217CurrencyCode ::= CHOICE { alphabetic PrintableString (SIZE 3), numeric INTEGER (1..999) }
func qcStatementUtilsQcLimitValue(statementInfo []byte) *extension.QCLimitValue {
	fields, ok := qcStatementUtilsSequenceElements(statementInfo)
	if !ok || len(fields) != 3 {
		return nil
	}
	// getCurrency().getAlphabetic() answers null - here the empty string - for a numeric code.
	currency := ""
	if text, isPrintableString := qcStatementUtilsPrintableString(fields[0]); isPrintableString {
		if len(text) > 3 {
			return nil // Iso4217CurrencyCode: IllegalArgumentException on an oversized code
		}
		currency = text
	} else if _, isNumeric := qcStatementUtilsIntegerElement(fields[0]); !isNumeric {
		return nil // Iso4217CurrencyCode: neither an INTEGER nor a PrintableString
	}
	amount, ok := qcStatementUtilsIntegerElement(fields[1])
	if !ok {
		return nil
	}
	exponent, ok := qcStatementUtilsIntegerElement(fields[2])
	if !ok {
		return nil
	}

	result := extension.NewQCLimitValue()
	result.SetCurrency(currency)
	result.SetAmount(qcStatementUtilsIntValue(amount))
	result.SetExponent(qcStatementUtilsIntValue(exponent))
	return result
}

// qcStatementUtilsQcEuRetentionPeriod ports the private getQcEuRetentionPeriod(ASN1Encodable).
// ASN1Integer#intValueExact() raises for a value outside the int range, which upstream swallows.
func qcStatementUtilsQcEuRetentionPeriod(statementInfo []byte) *int {
	value, ok := qcStatementUtilsIntegerElement(statementInfo)
	if !ok || !value.IsInt64() {
		return nil
	}
	exact := value.Int64()
	if exact < -2147483648 || exact > 2147483647 {
		return nil
	}
	retentionPeriod := int(exact)
	return &retentionPeriod
}

// qcStatementUtilsQcEuPDS ports the private getQcEuPDS(ASN1Encodable). Entries collected before a
// decoding failure are kept, as upstream builds the list outside its try block.
//
//	PdsLocations ::= SEQUENCE SIZE (1..MAX) OF PdsLocation
//	PdsLocation ::= SEQUENCE { url IA5String, language PrintableString (SIZE(2)) }
func qcStatementUtilsQcEuPDS(statementInfo []byte) []*extension.PdsLocation {
	result := make([]*extension.PdsLocation, 0)
	elements, ok := qcStatementUtilsSequenceElements(statementInfo)
	if !ok {
		return result
	}
	for _, element := range elements {
		fields, isSequence := qcStatementUtilsSequenceElements(element)
		if !isSequence {
			// Not an ASN1Sequence: upstream logs and skips the entry.
			continue
		}
		if len(fields) < 2 {
			// Upstream indexes the sequence blindly and aborts on the resulting exception.
			return result
		}
		pdsLocation := extension.NewPdsLocation()
		pdsLocation.SetUrl(DSSASN1UtilsString(fields[0]))
		pdsLocation.SetLanguage(DSSASN1UtilsString(fields[1]))
		result = append(result, pdsLocation)
	}
	return result
}

// qcStatementUtilsQcTypes ports the private getQcTypes(ASN1Encodable).
func qcStatementUtilsQcTypes(statementInfo []byte) []enumerations.QCType {
	oids := make([]string, 0)
	if elements, ok := qcStatementUtilsSequenceElements(statementInfo); ok {
		for _, element := range elements {
			if oid, isObjectIdentifier := qcStatementUtilsObjectIdentifierElement(element); isObjectIdentifier {
				oids = append(oids, oid)
			}
			// A non-ASN1ObjectIdentifier entry is logged and skipped upstream.
		}
	}
	return QcStatementUtilsQcTypesForOIDs(oids)
}

// QcStatementUtilsQcTypesForOIDs returns the QCTypes for the given QcType OIDs.
// Port of the public getQcTypes(List<String>).
func QcStatementUtilsQcTypesForOIDs(oids []string) []enumerations.QCType {
	result := make([]enumerations.QCType, 0, len(oids))
	for _, oid := range oids {
		if utils.IsStringNotBlank(oid) {
			result = append(result, enumerations.QCTypeFromOID(oid))
		}
	}
	return result
}

// qcStatementUtilsQcLegislationCountryCodes ports the private
// getQcLegislationCountryCodes(ASN1Encodable), shared by QcCClegislation and QcQSCDlegislation.
func qcStatementUtilsQcLegislationCountryCodes(statementInfo []byte) []string {
	result := make([]string, 0)
	elements, ok := qcStatementUtilsSequenceElements(statementInfo)
	if !ok {
		return result
	}
	for _, element := range elements {
		result = append(result, DSSASN1UtilsString(element))
	}
	return result
}

// qcStatementUtilsQcSemanticsIdentifier ports the private
// getQcSemanticsIdentifier(ASN1Encodable) over BouncyCastle's SemanticsInformation.
//
//	SemanticsInformation ::= SEQUENCE { semanticsIdentifier OBJECT IDENTIFIER OPTIONAL,
//	    nameRegistrationAuthorities NameRegistrationAuthorities OPTIONAL }
func qcStatementUtilsQcSemanticsIdentifier(statementInfo []byte) enumerations.SemanticsIdentifier {
	elements, ok := qcStatementUtilsSequenceElements(statementInfo)
	if !ok || len(elements) == 0 {
		return ""
	}
	oid, isObjectIdentifier := qcStatementUtilsObjectIdentifierElement(elements[0])
	if !isObjectIdentifier {
		return ""
	}
	return enumerations.SemanticsIdentifierFromOID(oid)
}

// qcStatementUtilsPsd2QcType ports the private getPsd2QcType(ASN1Encodable).
//
//	PSD2QcType ::= SEQUENCE { rolesOfPSP RolesOfPSP, nCAName UTF8String, nCAId UTF8String }
//	RolesOfPSP ::= SEQUENCE OF RoleOfPSP
//	RoleOfPSP ::= SEQUENCE { roleOfPspOid RoleOfPspOid, roleOfPspName RoleOfPspName }
func qcStatementUtilsPsd2QcType(statementInfo []byte) *extension.PSD2QcType {
	fields, ok := qcStatementUtilsSequenceElements(statementInfo)
	if !ok || len(fields) < 3 {
		return nil
	}
	roles, isSequence := qcStatementUtilsSequenceElements(fields[0])
	if !isSequence {
		return nil
	}
	rolesOfPSP := make([]*extension.RoleOfPSP, 0, len(roles))
	for _, role := range roles {
		roleFields, isSequence := qcStatementUtilsSequenceElements(role)
		if !isSequence || len(roleFields) < 2 {
			return nil
		}
		oid, isObjectIdentifier := qcStatementUtilsObjectIdentifierElement(roleFields[0])
		if !isObjectIdentifier {
			return nil // upstream casts to ASN1ObjectIdentifier and aborts on a ClassCastException
		}
		roleOfPSP := extension.NewRoleOfPSP()
		roleOfPSP.SetPspOid(enumerations.RoleOfPspOidFromOid(oid))
		roleOfPSP.SetPspName(DSSASN1UtilsString(roleFields[1]))
		rolesOfPSP = append(rolesOfPSP, roleOfPSP)
	}

	result := extension.NewPSD2QcType()
	result.SetRolesOfPSP(rolesOfPSP)
	result.SetNcaName(DSSASN1UtilsString(fields[1]))
	result.SetNcaId(DSSASN1UtilsString(fields[2]))
	return result
}

// qcStatementUtilsQcIdentMethod ports the private getQcIdentMethod(ASN1Encodable) together with the
// private getQcIdentMethod(String).
func qcStatementUtilsQcIdentMethod(statementInfo []byte) enumerations.QCIdentMethod {
	elements, ok := qcStatementUtilsSequenceElements(statementInfo)
	if !ok {
		return nil
	}
	if len(elements) != 1 {
		// The sequence size of QCIdentMethod shall be equal to 1; the value is skipped otherwise.
		return nil
	}
	oid, isObjectIdentifier := qcStatementUtilsObjectIdentifierElement(elements[0])
	if !isObjectIdentifier {
		return nil
	}
	if !utils.IsStringNotBlank(oid) {
		return nil
	}
	return enumerations.QCIdentMethodFromOID(oid)
}

// qcStatementUtilsQcPSB ports the private getQcPSB(ASN1Encodable).
func qcStatementUtilsQcPSB(statementInfo []byte) *extension.QCPSB {
	fields, ok := qcStatementUtilsSequenceElements(statementInfo)
	if !ok {
		return nil
	}
	if len(fields) != 3 {
		// The sequence size of QCPSB shall be equal to 3; the value is skipped otherwise.
		return nil
	}
	qcPSB := extension.NewQCPSB()
	qcPSB.SetCountryOfLegislation(DSSASN1UtilsString(fields[0]))
	qcPSB.SetAuthSourceIdentification(DSSASN1UtilsString(fields[1]))
	qcPSB.SetLegislationIdentification(DSSASN1UtilsString(fields[2]))
	return qcPSB
}

// QcStatementUtilsIsQcStatementPresent reports whether a QCStatement with the given OID is present
// within the QcStatements. Port of isQcStatementPresent(QcStatements, String).
func QcStatementUtilsIsQcStatementPresent(qcStatements *extension.QcStatements, qcStatementOid string) bool {
	switch {
	case QcStatementUtilsIsQcCompliance(qcStatementOid):
		return qcStatements.IsQcCompliance()
	case QcStatementUtilsIsQcLimitValue(qcStatementOid):
		return qcStatements.QcLimitValue() != nil
	case QcStatementUtilsIsQcRetentionPeriod(qcStatementOid):
		return qcStatements.QcEuRetentionPeriod() != nil
	case QcStatementUtilsIsQcSSCD(qcStatementOid):
		return qcStatements.IsQcQSCD()
	case QcStatementUtilsIsQcPds(qcStatementOid):
		return utils.IsCollectionNotEmpty(qcStatements.QcEuPDS())
	case QcStatementUtilsIsQcType(qcStatementOid):
		return utils.IsCollectionNotEmpty(qcStatements.QcTypes())
	case QcStatementUtilsIsQcCClegislation(qcStatementOid):
		return utils.IsCollectionNotEmpty(qcStatements.QcLegislationCountryCodes())
	case QcStatementUtilsIsQcSemanticsIdentifier(qcStatementOid):
		return qcStatements.QcSemanticsIdentifier() != ""
	case QcStatementUtilsIsPsd2QcType(qcStatementOid):
		return qcStatements.Psd2QcType() != nil
	case QcStatementUtilsIsQcQSCDlegislation(qcStatementOid):
		return utils.IsCollectionNotEmpty(qcStatements.QcQSCDLegislationCountryCodes())
	case QcStatementUtilsIsQcIdentMethod(qcStatementOid):
		return qcStatements.QcIdentMethod() != nil
	case QcStatementUtilsIsCertForPSB(qcStatementOid):
		return qcStatements.QcPSB() != nil
	default:
		for _, otherOid := range qcStatements.OtherOids() {
			if otherOid == qcStatementOid {
				return true
			}
		}
		return false
	}
}

// QcStatementUtilsIsQcTypePresent reports whether a QCType with the given OID is present within the
// QcStatements. Port of isQcTypePresent(QcStatements, String).
func QcStatementUtilsIsQcTypePresent(qcStatements *extension.QcStatements, qcTypeOid string) bool {
	for _, qcType := range qcStatements.QcTypes() {
		if qcTypeOid == qcType.OID() {
			return true
		}
	}
	return false
}

// QcStatementUtilsIsQcLegislationPresent reports whether a QCLegislation code is present within the
// QcStatements. Port of isQcLegislationPresent(QcStatements, String).
func QcStatementUtilsIsQcLegislationPresent(qcStatements *extension.QcStatements, qcLegislation string) bool {
	for _, countryCode := range qcStatements.QcLegislationCountryCodes() {
		if countryCode == qcLegislation {
			return true
		}
	}
	return false
}

// QcStatementUtilsIsQcQSCDlegislationPresent reports whether a QcQSCDLegislation code is present
// within the QcStatements. Port of isQcQSCDlegislationPresent(QcStatements, String).
func QcStatementUtilsIsQcQSCDlegislationPresent(qcStatements *extension.QcStatements, qcQSCDlegislation string) bool {
	for _, countryCode := range qcStatements.QcQSCDLegislationCountryCodes() {
		if countryCode == qcQSCDlegislation {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------------------------
// ASN.1 helpers replacing BouncyCastle. They are private to this file per the porting conventions.
// ---------------------------------------------------------------------------------------------

// qcStatementUtilsExtensionContent returns the extension value octets, unwrapped, standing in for
// BouncyCastle's X509CertificateImpl#getExtensionBytes(Certificate, String) followed by
// DSSASN1Utils#getAsn1SequenceFromDerOctetString.
func qcStatementUtilsExtensionContent(certToken *model.CertificateToken, oid string) []byte {
	for _, ext := range certToken.Certificate().Extensions {
		if ext.Id.String() == oid {
			return ext.Value
		}
	}
	return nil
}

// qcStatementUtilsQCStatement ports BouncyCastle's QCStatement decoding:
// QCStatement ::= SEQUENCE { statementId OBJECT IDENTIFIER, statementInfo ANY DEFINED BY statementId OPTIONAL }.
// A sequence of any other size raises an IllegalArgumentException there, which upstream logs and skips.
func qcStatementUtilsQCStatement(element []byte) (oid string, statementInfo []byte, ok bool) {
	fields, isSequence := qcStatementUtilsSequenceElements(element)
	if !isSequence || len(fields) < 1 || len(fields) > 2 {
		return "", nil, false
	}
	statementID, isObjectIdentifier := qcStatementUtilsObjectIdentifierElement(fields[0])
	if !isObjectIdentifier {
		return "", nil, false
	}
	if len(fields) == 2 {
		statementInfo = fields[1]
	}
	return statementID, statementInfo, true
}

// qcStatementUtilsSequenceElements returns the complete DER encoding of every element of the
// SEQUENCE starting at der. Trailing bytes are ignored, as ASN1InputStream#readObject() does.
func qcStatementUtilsSequenceElements(der []byte) ([][]byte, bool) {
	input := cryptobyte.String(der)
	var sequence cryptobyte.String
	if !input.ReadASN1(&sequence, cbasn1.SEQUENCE) {
		return nil, false
	}
	elements := make([][]byte, 0)
	for !sequence.Empty() {
		var element cryptobyte.String
		var tag cbasn1.Tag
		if !sequence.ReadAnyASN1Element(&element, &tag) {
			return nil, false
		}
		elements = append(elements, element)
	}
	return elements, true
}

// qcStatementUtilsObjectIdentifierElement decodes a complete OBJECT IDENTIFIER element.
func qcStatementUtilsObjectIdentifierElement(der []byte) (string, bool) {
	input := cryptobyte.String(der)
	var objectIdentifier encasn1.ObjectIdentifier
	if !input.ReadASN1ObjectIdentifier(&objectIdentifier) {
		return "", false
	}
	return objectIdentifier.String(), true
}

// qcStatementUtilsIntegerElement decodes a complete INTEGER element.
func qcStatementUtilsIntegerElement(der []byte) (*big.Int, bool) {
	input := cryptobyte.String(der)
	var content cryptobyte.String
	if !input.ReadASN1(&content, cbasn1.INTEGER) || len(content) == 0 {
		return nil, false
	}
	value := new(big.Int).SetBytes(content)
	if content[0]&0x80 != 0 {
		value.Sub(value, new(big.Int).Lsh(big.NewInt(1), uint(len(content))*8))
	}
	return value, true
}

// qcStatementUtilsIntValue ports java.math.BigInteger#intValue(), which keeps the low 32 bits.
func qcStatementUtilsIntValue(value *big.Int) int {
	low := new(big.Int).And(value, big.NewInt(0xFFFFFFFF))
	return int(int32(uint32(low.Uint64())))
}

// qcStatementUtilsPrintableString decodes a complete PrintableString element.
func qcStatementUtilsPrintableString(der []byte) (string, bool) {
	input := cryptobyte.String(der)
	var content cryptobyte.String
	if !input.ReadASN1(&content, cbasn1.PrintableString) {
		return "", false
	}
	runes := make([]rune, len(content))
	for i, b := range content {
		runes[i] = rune(b)
	}
	return string(runes), true
}
