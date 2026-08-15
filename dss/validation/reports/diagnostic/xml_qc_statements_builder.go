// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/XmlQcStatementsBuilder.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model/x509/extension"
	"github.com/utain/esig/dss/utils"
)

// XmlQcStatementsBuilder is used to build a XmlQcStatements object and enveloped objects.
type XmlQcStatementsBuilder struct {
}

// NewXmlQcStatementsBuilder is the port of the default constructor.
func NewXmlQcStatementsBuilder() *XmlQcStatementsBuilder {
	return &XmlQcStatementsBuilder{}
}

// Build builds the XmlQcStatements. Port of build(QcStatements).
func (b *XmlQcStatementsBuilder) Build(qcStatements *extension.QcStatements) *jaxb.XmlQcStatements {
	result := &jaxb.XmlQcStatements{}
	oid := qcStatements.OID()
	result.OID = &oid
	critical := qcStatements.IsCritical()
	result.Critical = &critical
	result.QcCompliance = b.BuildXmlQcCompliance(qcStatements.IsQcCompliance())
	result.QcSSCD = b.BuildXmlQcSSCD(qcStatements.IsQcQSCD())
	if qcStatements.QcEuRetentionPeriod() != nil {
		result.QcEuRetentionPeriod = qcStatements.QcEuRetentionPeriod()
	}
	if qcStatements.QcLimitValue() != nil {
		result.QcEuLimitValue = b.BuildQcEuLimitValue(qcStatements.QcLimitValue())
	}
	if utils.IsCollectionNotEmpty(qcStatements.QcTypes()) {
		result.QcTypes = &jaxb.QcTypesWrapper{Items: b.BuildXmlQcTypes(qcStatements.QcTypes())}
	}
	if utils.IsCollectionNotEmpty(qcStatements.QcEuPDS()) {
		result.QcEuPDS = &jaxb.QcEuPDSWrapper{Items: b.BuildXmlQcEuPSD(qcStatements.QcEuPDS())}
	}
	if qcStatements.QcSemanticsIdentifier() != "" {
		result.SemanticsIdentifier = b.BuildSemanticsIdentifier(qcStatements.QcSemanticsIdentifier())
	}
	if utils.IsCollectionNotEmpty(qcStatements.QcLegislationCountryCodes()) {
		result.QcCClegislation = &jaxb.QcCClegislationWrapper{Items: qcStatements.QcLegislationCountryCodes()}
	}
	if qcStatements.Psd2QcType() != nil {
		result.PSD2QcInfo = b.BuildPSD2QcInfo(qcStatements.Psd2QcType())
	}
	if utils.IsCollectionNotEmpty(qcStatements.QcQSCDLegislationCountryCodes()) {
		result.QcQSCDlegislation = &jaxb.QcQSCDlegislationWrapper{Items: qcStatements.QcQSCDLegislationCountryCodes()}
	}
	if qcStatements.QcIdentMethod() != nil {
		result.QcIdentMethod = b.getXmlOid(qcStatements.QcIdentMethod())
	}
	if qcStatements.QcPSB() != nil {
		result.QcPSB = b.BuildXmlQcPSB(qcStatements.QcPSB())
	}
	if utils.IsCollectionNotEmpty(qcStatements.OtherOids()) {
		result.OtherOIDs = &jaxb.OtherOIDsWrapper{Items: b.buildXmlOIDs(qcStatements.OtherOids())}
	}
	return result
}

// BuildXmlQcEuPSD builds a list of XML QcEuPSDs. Port of buildXmlQcEuPSD(List<PdsLocation>).
func (b *XmlQcStatementsBuilder) BuildXmlQcEuPSD(qcEuPDS []*extension.PdsLocation) []*jaxb.XmlLangAndValue {
	result := make([]*jaxb.XmlLangAndValue, 0, len(qcEuPDS))
	for _, pdsLocation := range qcEuPDS {
		xmlPdsLocation := &jaxb.XmlLangAndValue{}
		lang := pdsLocation.Language()
		xmlPdsLocation.Lang = &lang
		xmlPdsLocation.Value = pdsLocation.Url()
		result = append(result, xmlPdsLocation)
	}
	return result
}

// BuildXmlQcSSCD builds a XmlQcSSCD. Port of buildXmlQcSSCD(boolean).
func (b *XmlQcStatementsBuilder) BuildXmlQcSSCD(present bool) *jaxb.XmlQcSSCD {
	return &jaxb.XmlQcSSCD{Present: present}
}

// BuildXmlQcCompliance builds a XmlQcCompliance. Port of buildXmlQcCompliance(boolean).
func (b *XmlQcStatementsBuilder) BuildXmlQcCompliance(present bool) *jaxb.XmlQcCompliance {
	return &jaxb.XmlQcCompliance{Present: present}
}

// BuildPSD2QcInfo builds a XmlPSD2QcInfo. Port of buildPSD2QcInfo(PSD2QcType).
func (b *XmlQcStatementsBuilder) BuildPSD2QcInfo(psd2QcStatement *extension.PSD2QcType) *jaxb.XmlPSD2QcInfo {
	xmlInfo := &jaxb.XmlPSD2QcInfo{}
	ncaID := psd2QcStatement.NcaId()
	xmlInfo.NcaId = &ncaID
	ncaName := psd2QcStatement.NcaName()
	xmlInfo.NcaName = &ncaName
	rolesOfPSP := psd2QcStatement.RolesOfPSP()
	psd2Roles := make([]*jaxb.XmlRoleOfPSP, 0, len(rolesOfPSP))
	for _, roleOfPSP := range rolesOfPSP {
		xmlRole := &jaxb.XmlRoleOfPSP{}
		role := roleOfPSP.PspOid()
		xmlRole.Oid = b.getXmlOid(role)
		name := roleOfPSP.PspName()
		xmlRole.Name = &name
		psd2Roles = append(psd2Roles, xmlRole)
	}
	xmlInfo.RolesOfPSP = &jaxb.RolesOfPSPWrapper{Items: psd2Roles}
	return xmlInfo
}

// BuildSemanticsIdentifier builds a Semantics Identifier XmlOID. Port of
// buildSemanticsIdentifier(OidDescription).
func (b *XmlQcStatementsBuilder) BuildSemanticsIdentifier(semanticsIdentifier enumerations.OidDescription) *jaxb.XmlOID {
	return b.getXmlOid(semanticsIdentifier)
}

func (b *XmlQcStatementsBuilder) getXmlOid(oidDescription enumerations.OidDescription) *jaxb.XmlOID {
	if oidDescription == nil {
		return nil
	}
	xmlOID := &jaxb.XmlOID{}
	xmlOID.Value = oidDescription.OID()
	description := oidDescription.Description()
	xmlOID.Description = &description
	return xmlOID
}

// BuildXmlQcTypes builds a list of XML QcTypes. Port of buildXmlQcTypes(List<QCType>).
func (b *XmlQcStatementsBuilder) BuildXmlQcTypes(qcTypes []enumerations.QCType) []*jaxb.XmlOID {
	result := make([]*jaxb.XmlOID, 0, len(qcTypes))
	if utils.IsCollectionNotEmpty(qcTypes) {
		for _, qcType := range qcTypes {
			xmlOID := &jaxb.XmlOID{}
			xmlOID.Value = qcType.OID()
			description := qcType.Description()
			xmlOID.Description = &description
			result = append(result, xmlOID)
		}
	}
	return result
}

// buildXmlOIDs builds a list of XmlOIDs from a list of Strings. Port of buildXmlOIDs(List<String>).
func (b *XmlQcStatementsBuilder) buildXmlOIDs(oids []string) []*jaxb.XmlOID {
	result := make([]*jaxb.XmlOID, 0, len(oids))
	if utils.IsCollectionNotEmpty(oids) {
		for _, oid := range oids {
			xmlOID := &jaxb.XmlOID{}
			v := oid
			xmlOID.Value = v
			result = append(result, xmlOID)
		}
	}
	return result
}

// BuildQcEuLimitValue builds a XmlQcEuLimitValue. Port of buildQcEuLimitValue(QCLimitValue).
func (b *XmlQcStatementsBuilder) BuildQcEuLimitValue(qcLimitValue *extension.QCLimitValue) *jaxb.XmlQcEuLimitValue {
	xmlQcEuLimitValue := &jaxb.XmlQcEuLimitValue{}
	currency := qcLimitValue.Currency()
	xmlQcEuLimitValue.Currency = &currency
	xmlQcEuLimitValue.Amount = qcLimitValue.Amount()
	xmlQcEuLimitValue.Exponent = qcLimitValue.Exponent()
	return xmlQcEuLimitValue
}

// BuildXmlCertForPID builds a XmlCertForPID. Port of buildXmlCertForPID(boolean).
func (b *XmlQcStatementsBuilder) BuildXmlCertForPID(present bool) *jaxb.XmlCertForPID {
	if present {
		return &jaxb.XmlCertForPID{Present: true}
	}
	return nil
}

// BuildXmlCertForWallet builds a XmlCertForWallet. Port of buildXmlCertForWallet(boolean).
func (b *XmlQcStatementsBuilder) BuildXmlCertForWallet(present bool) *jaxb.XmlCertForWallet {
	if present {
		return &jaxb.XmlCertForWallet{Present: true}
	}
	return nil
}

// BuildXmlQcPSB builds a XmlQcPSB. Port of buildXmlQcPSB(QCPSB).
func (b *XmlQcStatementsBuilder) BuildXmlQcPSB(qcPSB *extension.QCPSB) *jaxb.XmlQcPSB {
	xmlQcPSB := &jaxb.XmlQcPSB{}
	countryOfLegislation := qcPSB.CountryOfLegislation()
	xmlQcPSB.CountryOfLegislation = &countryOfLegislation
	authSourceIdentification := qcPSB.AuthSourceIdentification()
	xmlQcPSB.AuthSourceIdentification = &authSourceIdentification
	legislationIdentification := qcPSB.LegislationIdentification()
	xmlQcPSB.LegislationIdentification = &legislationIdentification
	return xmlQcPSB
}

// Copy builds a deep copy of XmlQcStatements. NOTE: does not copy MRA content.
// Port of copy(XmlQcStatements).
func (b *XmlQcStatementsBuilder) Copy(xmlQcStatements *jaxb.XmlQcStatements) *jaxb.XmlQcStatements {
	copyVal := &jaxb.XmlQcStatements{}
	copyVal.OID = xmlQcStatements.OID
	copyVal.Description = xmlQcStatements.Description
	copyVal.Critical = xmlQcStatements.Critical
	if xmlQcStatements.QcCompliance != nil {
		copyVal.QcCompliance = &jaxb.XmlQcCompliance{Present: xmlQcStatements.QcCompliance.Present}
	}
	if xmlQcStatements.QcEuLimitValue != nil {
		copyVal.QcEuLimitValue = &jaxb.XmlQcEuLimitValue{
			Amount:   xmlQcStatements.QcEuLimitValue.Amount,
			Currency: xmlQcStatements.QcEuLimitValue.Currency,
			Exponent: xmlQcStatements.QcEuLimitValue.Exponent,
		}
	}
	copyVal.QcEuRetentionPeriod = xmlQcStatements.QcEuRetentionPeriod
	if xmlQcStatements.QcSSCD != nil {
		copyVal.QcSSCD = &jaxb.XmlQcSSCD{Present: xmlQcStatements.QcSSCD.Present}
	}
	if xmlQcStatements.SemanticsIdentifier != nil {
		copyVal.SemanticsIdentifier = &jaxb.XmlOID{
			XmlOIDContent: jaxb.XmlOIDContent{Value: xmlQcStatements.SemanticsIdentifier.Value},
			XmlOIDAttrs:   jaxb.XmlOIDAttrs{Description: xmlQcStatements.SemanticsIdentifier.Description},
		}
	}
	if xmlQcStatements.PSD2QcInfo != nil {
		xmlPSD2QcInfo := &jaxb.XmlPSD2QcInfo{
			NcaId:   xmlQcStatements.PSD2QcInfo.NcaId,
			NcaName: xmlQcStatements.PSD2QcInfo.NcaName,
		}
		roles := make([]*jaxb.XmlRoleOfPSP, 0, len(xmlQcStatements.PSD2QcInfo.RolesOfPSP.All()))
		for _, roleOfPSP := range xmlQcStatements.PSD2QcInfo.RolesOfPSP.All() {
			xmlRoleOfPSP := &jaxb.XmlRoleOfPSP{
				Name: roleOfPSP.Name,
				Oid:  roleOfPSP.Oid,
			}
			roles = append(roles, xmlRoleOfPSP)
		}
		xmlPSD2QcInfo.RolesOfPSP = &jaxb.RolesOfPSPWrapper{Items: roles}
		copyVal.PSD2QcInfo = xmlPSD2QcInfo
	}
	qcEuPDS := make([]*jaxb.XmlLangAndValue, 0)
	for _, xmlLangAndValue := range xmlQcStatements.QcEuPDS.All() {
		qcEuPDS = append(qcEuPDS, &jaxb.XmlLangAndValue{Lang: xmlLangAndValue.Lang, Value: xmlLangAndValue.Value})
	}
	copyVal.QcEuPDS = &jaxb.QcEuPDSWrapper{Items: qcEuPDS}
	qcTypes := make([]*jaxb.XmlOID, 0)
	for _, xmlOID := range xmlQcStatements.QcTypes.All() {
		qcTypes = append(qcTypes, &jaxb.XmlOID{
			XmlOIDContent: jaxb.XmlOIDContent{Value: xmlOID.Value},
			XmlOIDAttrs:   jaxb.XmlOIDAttrs{Description: xmlOID.Description},
		})
	}
	copyVal.QcTypes = &jaxb.QcTypesWrapper{Items: qcTypes}
	copyVal.QcCClegislation = &jaxb.QcCClegislationWrapper{Items: append([]string{}, xmlQcStatements.QcCClegislation.All()...)}
	copyVal.QcQSCDlegislation = &jaxb.QcQSCDlegislationWrapper{Items: append([]string{}, xmlQcStatements.QcQSCDlegislation.All()...)}
	if xmlQcStatements.QcIdentMethod != nil {
		copyVal.QcIdentMethod = &jaxb.XmlOID{
			XmlOIDContent: jaxb.XmlOIDContent{Value: xmlQcStatements.QcIdentMethod.Value},
			XmlOIDAttrs:   jaxb.XmlOIDAttrs{Description: xmlQcStatements.QcIdentMethod.Description},
		}
	}
	if xmlQcStatements.QcPSB != nil {
		copyVal.QcPSB = &jaxb.XmlQcPSB{
			CountryOfLegislation:      xmlQcStatements.QcPSB.CountryOfLegislation,
			AuthSourceIdentification:  xmlQcStatements.QcPSB.AuthSourceIdentification,
			LegislationIdentification: xmlQcStatements.QcPSB.LegislationIdentification,
		}
	}
	otherOIDs := make([]*jaxb.XmlOID, 0)
	for _, xmlOID := range xmlQcStatements.OtherOIDs.All() {
		otherOIDs = append(otherOIDs, &jaxb.XmlOID{
			XmlOIDContent: jaxb.XmlOIDContent{Value: xmlOID.Value},
			XmlOIDAttrs:   jaxb.XmlOIDAttrs{Description: xmlOID.Description},
		})
	}
	copyVal.OtherOIDs = &jaxb.OtherOIDsWrapper{Items: otherOIDs}
	return copyVal
}
